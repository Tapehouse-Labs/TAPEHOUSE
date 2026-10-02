// SPDX-License-Identifier: MIT OR Apache-2.0
// Usage: PRIVATE_KEY=0x… node examples/devnode.ts RPC_URL DEPLOYMENTS_JSON
// Starts the built server over Streamable HTTP on a dev node, without the key, and drives it with the official MCP
// client on the 2026-07-28 protocol: checks that the listener refuses another site's Host and Origin, reads SPY's and
// NVDA's bands and Morpho oracles against the SDK, then has the account authorize a fresh address through a prepared
// call, fund its SPY short, and the address sell 1 SPY short and buy it back, every call prepared by the server and
// signed here.
import { spawn } from 'node:child_process'
import { once } from 'node:events'
import { readFileSync } from 'node:fs'
import { request } from 'node:http'
import { Client, StreamableHTTPClientTransport } from '@modelcontextprotocol/client'
import { band, marginAccountsAbi, morpho, parseDeployments, shortPositionsAbi } from '@tapehouse/sdk'
import {
  type Address,
  createWalletClient,
  defineChain,
  type Hex,
  http,
  parseEther,
  parseEventLogs,
  publicActions,
} from 'viem'
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts'

const [rpc, registry] = process.argv.slice(2) as [string, string]
const key = process.env.PRIVATE_KEY as Hex | undefined
if (key === undefined) throw new Error("Set PRIVATE_KEY to the account owner's key.")
const deployments = parseDeployments(JSON.parse(readFileSync(registry, 'utf8')))
const chain = defineChain({
  id: deployments.chainId,
  name: 'dev node',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: [rpc] } },
})
const wallet = (secret: Hex) =>
  createWalletClient({ account: privateKeyToAccount(secret), chain, transport: http(rpc) }).extend(publicActions)
const owner = wallet(key)
const operator = wallet(generatePrivateKey())
const account = owner.account.address

function check(condition: boolean, message: string): asserts condition {
  if (!condition) throw new Error(`FAIL: ${message}`)
}

const server = spawn(process.execPath, [new URL('../dist/main.js', import.meta.url).pathname, '--http', '0'], {
  env: { PATH: process.env.PATH, TAPEHOUSE_RPC_URL: rpc, TAPEHOUSE_DEPLOYMENTS: registry },
  stdio: ['ignore', 'inherit', 'pipe'],
})
let banner = ''
server.stderr.setEncoding('utf8')
const url = await new Promise<URL>((resolve, reject) => {
  server.stderr.on('data', (chunk: string) => {
    banner += chunk
    const found = /http:\/\/127\.0\.0\.1:\d+\/mcp/.exec(banner)
    if (found) resolve(new URL(found[0]))
  })
  server.once('exit', () => reject(new Error(`FAIL: the server exited: ${banner}`)))
})
const mcp = new Client({ name: 'tapehouse-devnode', version: '0.1.0' }, { versionNegotiation: { mode: 'auto' } })

type Prepared = { signer: string; calls: { to: Address; data: Hex; functionName: string }[] }

async function use(name: string, args: Record<string, unknown> = {}) {
  const result = await mcp.callTool({ name, arguments: args })
  const [content] = result.content as [{ type: string; text: string }]
  check(result.isError !== true, `${name} failed: ${content.text}`)
  return result.structuredContent as Record<string, unknown>
}

/** The status the listener answers a request with `headers`, which `fetch` cannot set for `Host`. */
function status(headers: Record<string, string>) {
  return new Promise<number>((resolve, reject) => {
    const accept = { 'content-type': 'application/json', accept: 'application/json, text/event-stream' }
    request(url, { method: 'POST', headers: { ...accept, ...headers } }, (response) => {
      response.resume()
      resolve(response.statusCode ?? 0)
    })
      .on('error', reject)
      .end(JSON.stringify({ jsonrpc: '2.0', id: 1, method: 'ping' }))
  })
}

async function sign(client: typeof owner, prepared: Record<string, unknown>) {
  const receipts = []
  for (const call of (prepared as Prepared).calls) {
    const hash = await client.sendTransaction({ to: call.to, data: call.data })
    const receipt = await client.waitForTransactionReceipt({ hash })
    check(receipt.status === 'success', `${call.functionName} reverted`)
    receipts.push(receipt)
  }
  return receipts
}

try {
  for (const [header, value] of [
    ['host', 'attacker.example'],
    ['origin', 'https://attacker.example'],
  ] as const)
    check((await status({ [header]: value })) === 403, `the listener answered ${header}: ${value}`)
  await mcp.connect(new StreamableHTTPClientTransport(url))
  check(mcp.getProtocolEra() === 'modern', 'the client did not reach the 2026-07-28 protocol')
  const { tools } = await mcp.listTools()
  check(tools.length === 26 && tools.every((tool) => tool.annotations?.readOnlyHint), 'the tools are not all listed')

  const spy = await use('band_quote', { asset: 'SPY' })
  const at = { blockNumber: BigInt(spy.blockNumber as string) }
  const direct = await band.quote(owner, deployments, 'SPY', at)
  check(spy.low === direct.low.toString() && spy.state !== 'halted', "band_quote is not the band's quote")
  const priced = await use('morpho_oracle', { asset: 'SPY' })
  const oracleAt = { blockNumber: BigInt(priced.blockNumber as string) }
  const low = (await band.quote(owner, deployments, 'SPY', oracleAt)).low
  check(
    priced.price === (BigInt(low) * BigInt(priced.scaleFactor as string)).toString(),
    "SPY's oracle is not its low edge",
  )
  const nvda = await use('morpho_oracle', { asset: 'NVDA' })
  const nvdaAt = { blockNumber: BigInt(nvda.blockNumber as string) }
  const sdkAnswer = await morpho.price(owner, deployments, 'NVDA', nvdaAt)
  check(
    nvda.price === null && sdkAnswer.price === undefined && nvda.noPrice === sdkAnswer.noPrice,
    `NVDA's oracle answered ${String(nvda.price)}, or another reason than the SDK's`,
  )

  await owner.waitForTransactionReceipt({
    hash: await owner.sendTransaction({ to: operator.account.address, value: parseEther('0.01') }),
  })
  const before = await use('accounts_is_authorized', { account, operator: operator.account.address })
  check(before.authorized === false, 'a fresh address is authorized')
  const [authorized] = await sign(
    owner,
    await use('accounts_set_authorization', { operator: operator.account.address, allowed: true }),
  )
  const [set] = parseEventLogs({ abi: marginAccountsAbi, eventName: 'AuthorizationSet', logs: authorized?.logs ?? [] })
  check(set?.args.authorized === operator.account.address && set.args.allowed, 'no AuthorizationSet')

  const margin = '1000000000'
  await sign(owner, await use('shorts_deposit', { from: account, account, asset: 'SPY', amount: margin }))
  const sale = await use('shorts_sell', { account, asset: 'SPY', amount: parseEther('1').toString(), slippageBps: 50 })
  const [sold] = await sign(operator, sale)
  const [sell] = parseEventLogs({ abi: shortPositionsAbi, eventName: 'Sell', logs: sold?.logs ?? [] })
  const quoted = sale.quote as { proceeds: string; minProceeds: string }
  check(sell?.args.proceeds.toString() === quoted.proceeds, "the sale did not pay QuoterV2's quote")
  const open = await use('shorts_position', { account, asset: 'SPY' })
  check(BigInt(open.debt as string) >= parseEther('1'), 'the short does not owe 1 SPY')
  const health = await use('shorts_health', { account, asset: 'SPY' })
  check(BigInt(health.surplus as string) > 0n, 'the short falls short')

  const buyBack = await use('shorts_cover', { account, asset: 'SPY', amount: 'max', slippageBps: 50 })
  const [covered] = await sign(operator, buyBack)
  const [cover] = parseEventLogs({ abi: shortPositionsAbi, eventName: 'Cover', logs: covered?.logs ?? [] })
  const limits = buyBack.quote as { cost: string; maxCost: string }
  check(
    cover !== undefined && cover.args.cost >= BigInt(limits.cost) && cover.args.cost <= BigInt(limits.maxCost),
    "the buy-back did not cost QuoterV2's quote, and the fee accrued since, within the slippage",
  )
  const closed = await use('shorts_position', { account, asset: 'SPY' })
  check(closed.shares === '0' && closed.debt === '0', 'the short is still open')
  await sign(
    operator,
    await use('shorts_withdraw', { account, asset: 'SPY', amount: closed.usdgHeld, receiver: account }),
  )
  await sign(owner, await use('accounts_set_authorization', { operator: operator.account.address, allowed: false }))
  const after = await use('accounts_is_authorized', { account, operator: operator.account.address })
  check(after.authorized === false, 'the operator is still authorized')

  console.log(
    `typescript client over Streamable HTTP, ${mcp.getProtocolEra()} era: another site's Host and Origin refused; ` +
      `SPY band ${spy.state} at ${spy.mid}, ` +
      `its Morpho oracle ${priced.price}; NVDA's oracle no price (${nvda.noPrice}); ` +
      `sold 1 SPY for ${sell.args.proceeds} and bought it back for ${cover.args.cost}, every call prepared unsigned`,
  )
  console.log('PASS')
} finally {
  await mcp.close()
  server.kill()
  await once(server, 'exit')
}
