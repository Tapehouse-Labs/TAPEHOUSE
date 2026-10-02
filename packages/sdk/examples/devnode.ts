// SPDX-License-Identifier: MIT OR Apache-2.0
// Usage: PRIVATE_KEY=0x… node examples/devnode.ts RPC_URL DEPLOYMENTS_JSON PAYLOAD_FILE
// Reads SPY's band and its BandFeed on a dev node, sends RedStone packages from a file through a PackageSource and
// decodes the band's refusal of their age, then authorizes a fresh address in the margin accounts, which sells 1 SPY
// short for the account at a QuoterV2 quote and buys it back.
import { readFileSync } from 'node:fs'
import {
  type Address,
  createWalletClient,
  defineChain,
  type Hex,
  http,
  maxUint256,
  parseEther,
  parseEventLogs,
  publicActions,
} from 'viem'
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts'
import {
  accounts,
  band,
  CROSS,
  decodeRevert,
  marginAccountsAbi,
  type PackageSource,
  parseDeployments,
  shortPositionsAbi,
  shorts,
  usdgAbi,
} from '@tapehouse/sdk'

const [rpc, registry, payloadFile] = process.argv.slice(2) as [string, string, string]
const key = process.env.PRIVATE_KEY as Hex | undefined
if (key === undefined) throw new Error('Set PRIVATE_KEY to the account owner\'s key.')
const deployments = parseDeployments(JSON.parse(readFileSync(registry, 'utf8')))
const chain = defineChain({
  id: deployments.chainId,
  name: 'dev node',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: [rpc] } },
})
const owner = createWalletClient({ account: privateKeyToAccount(key), chain, transport: http(rpc) }).extend(
  publicActions,
)
const operator = createWalletClient({
  account: privateKeyToAccount(generatePrivateKey()),
  chain,
  transport: http(rpc),
}).extend(publicActions)
const account = owner.account.address
const one = 10n ** 18n

function check(condition: boolean, message: string): asserts condition {
  if (!condition) throw new Error(`FAIL: ${message}`)
}

async function refusal(promise: Promise<unknown>) {
  const revert = await promise.then(
    () => undefined,
    (error: unknown) => decodeRevert(error),
  )
  check(revert !== undefined, 'the call did not revert with a decodable error')
  return `${revert.errorName}(${revert.args.join(', ')})`
}

async function send(client: typeof owner, call: Parameters<typeof owner.simulateContract>[0]) {
  const { request } = await client.simulateContract({ ...call, account: client.account })
  const receipt = await client.waitForTransactionReceipt({ hash: await client.writeContract(request) })
  check(receipt.status === 'success', `${call.functionName} reverted`)
  return receipt
}

check((await owner.getChainId()) === deployments.chainId, 'the registry is for another chain')
const block = await owner.getBlock()
const at = { blockNumber: block.number }
const quote = await band.quote(owner, deployments, 'SPY', at)
const session = await band.session(owner, deployments, at)
const halt = await band.halt(owner, deployments, 'SPY', at)
check(quote.state !== 0 && !halt.signedHalt, "SPY's band is halted")
const round = await band.latestRound(owner, deployments, 'SPY', at)
check(round.answer === BigInt(quote.low), "the feed's answer is not the band's low edge")
check(round.roundId === block.timestamp && round.updatedAt === block.timestamp, 'the round is not the block time')
const whole = await band.latestBand(owner, deployments, 'SPY', at)
check(whole.state === quote.state && whole.mid === quote.mid, 'latestBand is not the band')
const sealedBand = await band.sealed(owner, deployments, 'SPY', session.boundaryMs, at)
const sealing = await owner.simulateContract({ ...band.seal(deployments, 'SPY'), account }).then(
  () => 'within its window',
  (error: unknown) => {
    const revert = decodeRevert(error)
    check(
      revert?.errorName === 'NotSealWindow' &&
        revert.args[0] === session.state &&
        revert.args[1] === session.boundaryMs,
      'seal() did not revert with NotSealWindow for the session',
    )
    return `NotSealWindow(${revert.args.join(', ')})`
  },
)

const fromFile: PackageSource = { payload: async () => readFileSync(payloadFile, 'utf8').trim() as Hex }
const write = await band.writePrices(deployments, fromFile, ['NVDA---24_7'])
const stale = await refusal(owner.simulateContract({ ...write, account }))
check(stale.startsWith('TimestampIsTooOld(1790119750, '), `the band took the file's packages: ${stale}`)

await owner.waitForTransactionReceipt({
  hash: await owner.sendTransaction({ to: operator.account.address, value: parseEther('0.01') }),
})
const usdg = deployments.tokens.USDG as Address
const margin = 1_000_000_000n
const refused = await refusal(
  operator.simulateContract({
    ...shorts.deposit(deployments, { asset: 'SPY', amount: margin, account }),
    account: operator.account,
  }),
)
check(refused === `Unauthorized(${operator.account.address}, ${account})`, `the shorts did not refuse: ${refused}`)
const withdrawal = accounts.withdraw(deployments, {
  position: CROSS,
  token: usdg,
  amount: 1n,
  account,
  receiver: account,
})
check(
  (await refusal(operator.simulateContract({ ...withdrawal, account: operator.account }))) === refused,
  'the margin accounts did not refuse the operator',
)

const authorized = await send(owner, accounts.setAuthorization(deployments, operator.account.address, true))
const [set] = parseEventLogs({ abi: marginAccountsAbi, eventName: 'AuthorizationSet', logs: authorized.logs })
check(set?.args.authorized === operator.account.address && set.args.allowed, 'no AuthorizationSet')
check(
  await accounts.isAuthorized(owner, deployments, account, operator.account.address),
  'the operator is not authorized',
)

await send(owner, {
  address: usdg,
  abi: usdgAbi,
  functionName: 'approve',
  args: [deployments.tapehouse.ShortPositions as Address, margin],
})
await send(owner, shorts.deposit(deployments, { asset: 'SPY', amount: margin, account }))
const sale = await shorts.quoteSale(owner, deployments, { asset: 'SPY', amount: one, slippageBps: 50n })
const sold = await send(
  operator,
  shorts.sell(deployments, { asset: 'SPY', amount: one, minProceeds: sale.minProceeds, account }),
)
const [sell] = parseEventLogs({ abi: shortPositionsAbi, eventName: 'Sell', logs: sold.logs })
check(sell?.args.proceeds === sale.proceeds, "the sale did not pay QuoterV2's quote")
const open = await shorts.position(owner, deployments, account, 'SPY')
check(open.debt >= one, 'the short does not owe 1 SPY and its fee')
const health = await shorts.health(owner, deployments, account, 'SPY')
check(health.equity > BigInt(health.requirement), 'the short falls short')
const buyBack = await shorts.quoteCover(owner, deployments, { asset: 'SPY', amount: open.debt, slippageBps: 50n })
const covered = await send(
  operator,
  shorts.cover(deployments, { asset: 'SPY', amount: maxUint256, maxCost: buyBack.maxCost, account }),
)
const [cover] = parseEventLogs({ abi: shortPositionsAbi, eventName: 'Cover', logs: covered.logs })
check(
  cover !== undefined && cover.args.cost >= buyBack.cost && cover.args.cost <= buyBack.maxCost,
  "the buy-back did not cost QuoterV2's quote, and the fee accrued since, within the slippage",
)
const closed = await shorts.position(owner, deployments, account, 'SPY')
check(closed.shares === 0n && closed.debt === 0n, 'the short is still open')
await send(
  operator,
  shorts.withdraw(deployments, { asset: 'SPY', amount: BigInt(closed.usdgHeld), account, receiver: account }),
)
await send(owner, accounts.setAuthorization(deployments, operator.account.address, false))
check(
  !(await accounts.isAuthorized(owner, deployments, account, operator.account.address)),
  'the operator is still authorized',
)

console.log(
  `typescript sdk: SPY band state ${quote.state} at ${quote.mid}, ` +
    `feed round ${round.roundId} answers ${round.answer}, ` +
    `seal at ${session.boundaryMs} ${sealing} (sealed ${sealedBand.sealedAt}); writePrices ${stale}; ` +
    `before authorization ${refused}; sold 1 SPY for ${sell.args.proceeds} and bought it back for ${cover.args.cost}`,
)
console.log('PASS')
