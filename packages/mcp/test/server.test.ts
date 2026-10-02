// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from 'node:fs'
import { Client, InMemoryTransport, StreamableHTTPClientTransport } from '@modelcontextprotocol/client'
import { createMcpHandler } from '@modelcontextprotocol/server'
import {
  type Abi,
  type Address,
  type Client as ViemClient,
  createClient,
  custom,
  decodeFunctionData,
  encodeErrorResult,
  encodeFunctionResult,
  erc20Abi,
  type Hex,
  http,
  maxUint256,
  toFunctionSelector,
  zeroAddress,
} from 'viem'
import { afterEach, describe, expect, test } from 'vitest'
import {
  bandAbi,
  basketAbi,
  CROSS,
  type Deployments,
  gapCoverAbi,
  marginAccountsAbi,
  morphoBandOracleAbi,
  parseDeployments,
  quoterV2Abi,
  shortPositionsAbi,
  toBytes32,
} from '@tapehouse/sdk'
import { configFromEnv } from '../src/config.ts'
import { createServer, RateLimit } from '../src/server.ts'

const robinhood = parseDeployments(
  JSON.parse(readFileSync(new URL('../../../deployments/4663.json', import.meta.url), 'utf8')),
)
const alice: Address = '0x70997970C51812dc3A010C7d01b50e0d17dc79C8'
const bob: Address = '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC'
const band: Address = '0xa70118d3324D90532E7D2854627b13CacE305641'
const accounts: Address = '0x5FbDB2315678afecb367f032d93F642f64180aa3'
const shorts: Address = '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
const oracle: Address = '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0'
const pair: Address = '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9'
const cover: Address = '0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9'
const nvda = robinhood.tokens.NVDA as Address
const spy = robinhood.tokens.SPY as Address
const deployments: Deployments = {
  ...robinhood,
  tapehouse: { Band: band, MarginAccounts: accounts, ShortPositions: shorts, GapCover: cover },
  morphoOracles: { NVDA: oracle, SPY: oracle },
  baskets: { PAIR: pair },
}
const abis: Record<string, Abi> = {
  [band]: bandAbi,
  [accounts]: marginAccountsAbi,
  [shorts]: shortPositionsAbi,
  [oracle]: morphoBandOracleAbi,
  [pair]: basketAbi,
  [cover]: gapCoverAbi,
  [robinhood.uniswapV3.QuoterV2 as string]: quoterV2Abi,
}

type Answer = unknown | { revert: Hex } | { fail: string }
type Call = { to: string; functionName: string; args: readonly unknown[]; block: string }

function chain(answers: Record<string, Answer | ((args: readonly unknown[], to: string) => Answer)>) {
  const calls: Call[] = []
  const methods: string[] = []
  const client = createClient({
    transport: custom(
      {
        async request({ method, params }) {
          methods.push(method)
          if (method === 'eth_blockNumber') return '0x7'
          if (method === 'eth_chainId') return '0x123f'
          if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
          const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
          const abi = abis[to] ?? erc20Abi
          const { functionName, args = [] } = decodeFunctionData({ abi, data })
          calls.push({ to, functionName, args, block })
          const answer = answers[functionName]
          if (answer === undefined) throw new Error(`unexpected ${functionName} on ${to}`)
          const result = typeof answer === 'function' ? answer(args, to) : answer
          if (typeof result === 'object' && result !== null && 'revert' in result)
            throw Object.assign(new Error('execution reverted'), { code: 3, data: result.revert })
          if (typeof result === 'object' && result !== null && 'fail' in result)
            throw Object.assign(new Error(result.fail), { code: 3 })
          return encodeFunctionResult({ abi, functionName, result } as Parameters<typeof encodeFunctionResult>[0])
        },
      },
      { retryCount: 0 },
    ),
  })
  return { client, calls, methods }
}

const open: Client[] = []

async function legacy(client: ViemClient, maxSlippageBps = 100, rateLimit = new RateLimit(120)) {
  const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair()
  await createServer({ client, deployments, maxSlippageBps, rateLimit }).connect(serverTransport)
  const mcp = new Client({ name: 'test', version: '1.0.0' })
  await mcp.connect(clientTransport)
  open.push(mcp)
  return mcp
}

async function modern(client: ViemClient) {
  const rateLimit = new RateLimit(120)
  const handler = createMcpHandler(() => createServer({ client, deployments, maxSlippageBps: 100, rateLimit }))
  const mcp = new Client({ name: 'test', version: '1.0.0' }, { versionNegotiation: { mode: 'auto' } })
  await mcp.connect(
    new StreamableHTTPClientTransport(new URL('http://mcp.test/mcp'), {
      fetch: (url, init) => handler.fetch(new Request(url, init)),
    }),
  )
  open.push(mcp)
  return mcp
}

afterEach(async () => {
  await Promise.all(open.splice(0).map((mcp) => mcp.close()))
})

async function call(mcp: Client, name: string, args: Record<string, unknown>) {
  const result = await mcp.callTool({ name, arguments: args })
  const [content] = result.content as [{ type: string; text: string }]
  return { result, text: content.text, structured: result.structuredContent as Record<string, unknown> }
}

const quote = { quote: [2, 1, 76_551_775_442n, 55n, 76_130_740_177n, 76_972_810_707n] }

describe('tools', () => {
  test('every tool has a title, read-only annotations and an output schema, in both protocol eras', async () => {
    const { client } = chain({})
    const older = await legacy(client)
    const newer = await modern(client)
    expect(older.getProtocolEra()).toBe('legacy')
    expect(newer.getProtocolEra()).toBe('modern')
    const [{ tools }, { tools: same }] = await Promise.all([older.listTools(), newer.listTools()])
    expect(same.map((tool) => tool.name)).toEqual(tools.map((tool) => tool.name))
    expect(tools.map((tool) => tool.name)).toEqual([
      'registry',
      'band_quote',
      'band_session',
      'band_halt',
      'band_latest_band',
      'band_sealed',
      'chainlink_round',
      'morpho_oracle',
      'baskets_components',
      'baskets_preview_mint',
      'baskets_preview_redeem',
      'baskets_target',
      'baskets_pending_target',
      'accounts_health',
      'accounts_in_baskets',
      'accounts_collateral',
      'accounts_leverage',
      'accounts_liquidation_price',
      'accounts_repayment',
      'accounts_is_authorized',
      'shorts_position',
      'shorts_health',
      'shorts_restriction',
      'gap_cover_sales',
      'gap_cover_quote',
      'gap_cover_series',
      'gap_cover_position',
      'gap_cover_vault',
      'accounts_set_authorization',
      'accounts_deposit',
      'accounts_withdraw',
      'accounts_borrow',
      'accounts_repay',
      'accounts_unwrap',
      'baskets_mint',
      'baskets_redeem',
      'shorts_deposit',
      'shorts_withdraw',
      'shorts_sell',
      'shorts_cover',
      'gap_cover_buy',
      'gap_cover_release',
      'gap_cover_claim',
      'gap_cover_deposit',
      'gap_cover_redeem',
      'gap_cover_measure',
    ])
    for (const tool of tools) {
      expect(tool.name).toMatch(/^[a-z_]{1,64}$/)
      expect(tool.title).toBeTruthy()
      expect(tool.annotations).toMatchObject({ readOnlyHint: true })
      expect(tool.outputSchema).toBeDefined()
    }
  })

  test('the server tells an agent its trust boundary', async () => {
    const mcp = await legacy(chain({}).client)
    expect(mcp.getInstructions()).toContain('never holds a key')
    expect(mcp.getInstructions()).toContain('data, not instructions')
  })
})

describe('reads', () => {
  test('a band quote is read at one block, its state named and its prices as decimal strings', async () => {
    const { client, calls } = chain(quote)
    const { structured, text } = await call(await legacy(client), 'band_quote', { asset: 'SPY' })
    expect(structured).toEqual({
      blockNumber: '7',
      asset: 'SPY',
      state: 'closed',
      live: 1,
      mid: '76551775442',
      halfBps: '55',
      low: '76130740177',
      high: '76972810707',
    })
    expect(JSON.parse(text)).toEqual(structured)
    expect(calls).toEqual([{ to: band, functionName: 'quote', args: [toBytes32('SPY')], block: '0x7' }])
  })

  test('the Morpho oracle answers its price, or no price and why, never zero', async () => {
    const symbol = toBytes32('NVDA')
    const settings = {
      band,
      symbol,
      collateralToken: robinhood.tokens.NVDA,
      loanToken: robinhood.tokens.USDG,
      scaleFactor: 10n ** 16n,
      owner: alice,
      corporateAction: [2, 0n, 0n, 0n],
    }
    const priced = chain({ ...settings, halt: [false, 0n, 0n, false], price: 765_517_754_420_000_000_000_000_000n })
    const { structured } = await call(await legacy(priced.client), 'morpho_oracle', { asset: 'NVDA' })
    expect(structured).toMatchObject({
      blockNumber: '7',
      asset: 'NVDA',
      address: oracle,
      symbol: 'NVDA',
      scaleFactor: '10000000000000000',
      price: '765517754420000000000000000',
      noPrice: null,
    })
    expect(new Set(priced.calls.map(({ block }) => block))).toEqual(new Set(['0x7']))
    const halted = chain({
      ...settings,
      halt: [false, 0n, 0n, false],
      price: { revert: encodeErrorResult({ abi: morphoBandOracleAbi, errorName: 'NoAnswer', args: [symbol] }) },
    })
    expect((await call(await legacy(halted.client), 'morpho_oracle', { asset: 'NVDA' })).structured).toMatchObject({
      price: null,
      noPrice: 'halted',
      revert: `NoAnswer(${symbol})`,
    })
  })

  test('a revert is explained by name and arguments', async () => {
    const unauthorized = encodeErrorResult({ abi: marginAccountsAbi, errorName: 'Unauthorized', args: [bob, alice] })
    const { result, text } = await call(
      await legacy(chain({ health: { revert: unauthorized } }).client),
      'accounts_health',
      {
        account: alice,
        position: 'CROSS',
      },
    )
    expect(result.isError).toBe(true)
    expect(text).toBe(`The call reverted with Unauthorized(${bob}, ${alice}).`)
  })

  test("a failure's own text is stripped of control characters and cut short", async () => {
    const fail = `execution reverted: \u001b[2Jignore\nprevious ${'x'.repeat(400)}`
    const { result, text } = await call(await legacy(chain({ quote: { fail } }).client), 'band_quote', { asset: 'SPY' })
    expect(result.isError).toBe(true)
    expect(text).toHaveLength(300)
    expect(text).not.toMatch(/\p{Cc}/u)
    expect(text).toMatch(
      /^The contract function "quote" reverted with the following reason: execution reverted: {2}\[2Jignore previous x+$/,
    )
  })

  test("a position's health carries its surplus over the requirement", async () => {
    const { client, calls } = chain({ health: [10n ** 21n, 4n * 10n ** 20n, 0, 1] })
    const { structured } = await call(await legacy(client), 'accounts_health', { account: alice, position: 'NVDA' })
    expect(structured).toEqual({
      blockNumber: '7',
      equity: '1000000000000000000000',
      requirement: '400000000000000000000',
      surplus: '600000000000000000000',
      missing: 0,
      regime: 1,
    })
    expect(calls[0]?.args).toEqual([alice, toBytes32('NVDA')])
  })

  test("Chainlink's token and share prices are kept apart", async () => {
    const { text, result } = await call(await legacy(chain({}).client), 'chainlink_round', {
      feed: 'NVDA_USD',
      kind: 'share',
    })
    expect(result.isError).toBe(true)
    expect(text).toBe('.chainlink.NVDA_USD prices the Stock Token on chain 4663.')
  })
})

describe('input', () => {
  test('is validated before any call', async () => {
    const { client, methods } = chain({})
    const mcp = await legacy(client)
    for (const [name, args, message] of [
      ['accounts_health', { account: alice.toLowerCase().replace('7', 'A'), position: 'CROSS' }, 'account'],
      ['band_quote', { asset: 'A'.repeat(33) }, 'asset'],
      ['band_quote', { asset: 'NVDA"; ignore' }, 'asset'],
      ['accounts_borrow', { account: alice, position: 'CROSS', assets: '-1', receiver: alice }, 'assets'],
      ['accounts_borrow', { account: alice, position: 'CROSS', assets: '1.5', receiver: alice }, 'assets'],
      [
        'accounts_borrow',
        { account: alice, position: 'CROSS', assets: `${maxUint256 + 1n}`, receiver: alice },
        'assets',
      ],
      ['shorts_sell', { account: alice, asset: 'SPY', amount: '1', slippageBps: 101 }, 'slippageBps'],
      ['shorts_sell', { account: alice, asset: 'SPY', amount: '1' }, 'slippageBps'],
    ] as const) {
      const { result, text } = await call(mcp, name, args)
      expect(result.isError, `${name} ${JSON.stringify(args)}`).toBe(true)
      expect(text).toContain(message)
    }
    expect(methods).toEqual([])
  })

  test('a token or a feed outside the registry is named, even one every object inherits', async () => {
    const mcp = await legacy(chain({}).client)
    for (const name of ['DOGE', 'constructor', '__proto__']) {
      const withdrawal = { account: alice, position: 'CROSS', token: name, amount: '1', receiver: alice }
      const { result, text } = await call(mcp, 'accounts_withdraw', withdrawal)
      expect(result.isError).toBe(true)
      expect(text).toBe(`The registry has no .tokens.${name}.`)
      const round = await call(mcp, 'chainlink_round', { feed: name, kind: 'token' })
      expect(round.result.isError).toBe(true)
      expect(round.text).toBe(`The registry has no .chainlink.${name}.`)
    }
  })

  test("calls beyond the operator's limit a minute are refused before any call", async () => {
    const { client, methods } = chain(quote)
    const mcp = await legacy(client, 100, new RateLimit(2))
    for (let i = 0; i < 2; i++) expect((await call(mcp, 'band_quote', { asset: 'SPY' })).result.isError).toBeFalsy()
    const read = methods.length
    const { result, text } = await call(mcp, 'band_quote', { asset: 'SPY' })
    expect(result.isError).toBe(true)
    expect(text).toBe('At most 2 tool calls a minute: try again later.')
    expect(methods).toHaveLength(read)
    const window = new RateLimit(2)
    expect([window.take(0), window.take(1), window.take(59_999), window.take(60_000), window.take(60_001)]).toEqual([
      true,
      true,
      false,
      true,
      true,
    ])
  })

  test('an RPC failure never shows the RPC URL', async () => {
    const client = createClient({ transport: http('http://127.0.0.1:9/secret-key', { retryCount: 0 }) })
    const { result, text } = await call(await legacy(client), 'band_session', {})
    expect(result.isError).toBe(true)
    expect(text).not.toContain('secret-key')
    expect(text).not.toContain('127.0.0.1')
  })
})

describe('prepared transactions', () => {
  test('an authorization is prepared for the account to sign, and nothing is signed', async () => {
    const { client, methods } = chain({})
    const { structured } = await call(await legacy(client), 'accounts_set_authorization', {
      operator: bob,
      allowed: true,
    })
    expect(structured).toEqual({
      chainId: 4663,
      signer: 'the account, which lets the operator act for it',
      grants: expect.stringContaining(`${bob} can withdraw`),
      calls: [
        {
          to: accounts,
          data: expect.stringMatching(/^0x/),
          value: '0',
          functionName: 'setAuthorization',
          args: { authorized: bob, allowed: true },
        },
      ],
    })
    const [prepared] = structured.calls as [{ data: Hex }]
    expect(decodeFunctionData({ abi: marginAccountsAbi, data: prepared.data }).args).toEqual([bob, true])
    expect(methods.filter((method) => method.startsWith('eth_send') || method.startsWith('eth_sign'))).toEqual([])
  })

  test('an authorization states everything it grants', async () => {
    const mcp = await legacy(chain({}).client)
    const { structured } = await call(mcp, 'accounts_set_authorization', { operator: bob, allowed: true })
    expect(structured.grants).toBe(
      `${bob} can withdraw every collateral token of every margin position and the USDG of every short to any address, borrow USDG to any address, and sell or buy back the account's shorts, until the account signs setAuthorization(${bob}, false).`,
    )
    const { tools } = await mcp.listTools()
    expect(tools.find((tool) => tool.name === 'accounts_set_authorization')?.description).toContain(
      "operator can withdraw every collateral token of every margin position and the USDG of every short to any address, borrow USDG to any address, and sell or buy back the account's shorts, until the account signs setAuthorization(operator, false).",
    )
    const revoked = await call(mcp, 'accounts_set_authorization', { operator: bob, allowed: false })
    expect(revoked.structured.grants).toBe(`Nothing: ${bob} can no longer act for the account.`)
  })

  test('a deposit is preceded by an approval only where the allowance falls short', async () => {
    const nvda = robinhood.tokens.NVDA as string
    const args = { from: bob, account: alice, position: 'NVDA', token: 'NVDA', amount: '1000000000000000000' }
    const short = chain({ allowance: 0n })
    const { structured } = await call(await legacy(short.client), 'accounts_deposit', args)
    expect(structured.signer).toBe(bob)
    expect(structured.calls).toMatchObject([
      { to: nvda, functionName: 'approve', args: { spender: accounts, amount: '1000000000000000000' } },
      {
        to: accounts,
        functionName: 'deposit',
        args: { position: toBytes32('NVDA'), token: nvda, amount: '1000000000000000000', account: alice },
      },
    ])
    expect(short.calls).toEqual([{ to: nvda, functionName: 'allowance', args: [bob, accounts], block: '0x7' }])
    const enough = chain({ allowance: maxUint256 })
    const { structured: direct } = await call(await legacy(enough.client), 'accounts_deposit', args)
    expect((direct.calls as { functionName: string }[]).map(({ functionName }) => functionName)).toEqual(['deposit'])
  })

  test("a sale's limit is QuoterV2's quote less the slippage the caller states, within the server's cap", async () => {
    const { client, calls } = chain({ fee: 500, quoteExactInputSingle: [77_232_802n, 0n, 0, 0n] })
    const { structured } = await call(await legacy(client, 50), 'shorts_sell', {
      account: alice,
      asset: 'SPY',
      amount: '100000000000000000',
      slippageBps: 50,
    })
    expect(structured).toMatchObject({
      signer: `${alice}, or an address it authorized`,
      quote: { proceeds: '77232802', minProceeds: '76846637', slippageBps: 50 },
      calls: [
        {
          to: shorts,
          functionName: 'sell',
          args: { symbol: toBytes32('SPY'), amount: '100000000000000000', minProceeds: '76846637', account: alice },
        },
      ],
    })
    expect(calls.map(({ functionName }) => functionName)).toEqual(['fee', 'quoteExactInputSingle'])
    const { result, text } = await call(await legacy(chain({}).client, 50), 'shorts_sell', {
      account: alice,
      asset: 'SPY',
      amount: '1',
      slippageBps: 51,
    })
    expect(result.isError).toBe(true)
    expect(text).toContain('slippageBps')
  })

  test('a cover of everything quotes what the short owes and buys back all of it', async () => {
    const { client } = chain({
      position: [10n ** 21n, 1_000_300_000_000_000_000n, 10n ** 18n, 0n],
      fee: 500,
      quoteExactOutputSingle: (args: readonly unknown[]) => {
        expect(args).toMatchObject([{ amount: 1_000_300_000_000_000_000n }])
        return [773_933_000n, 0n, 0, 0n]
      },
    })
    const { structured } = await call(await legacy(client), 'shorts_cover', {
      account: alice,
      asset: 'SPY',
      amount: 'max',
      slippageBps: 100,
    })
    expect(structured).toMatchObject({
      quote: { amount: '1000300000000000000', cost: '773933000', maxCost: '781672330', slippageBps: 100 },
      calls: [{ functionName: 'cover', args: { amount: maxUint256.toString(), maxCost: '781672330' } }],
    })
    const [{ data }] = structured.calls as [{ data: Hex }]
    expect(data.slice(0, 10)).toBe(toFunctionSelector('cover(bytes32,uint256,uint256,address)'))
  })
})

const components = {
  components: [
    [toBytes32('NVDA'), toBytes32('SPY')],
    [nvda, spy],
  ],
}

describe('baskets', () => {
  test("a basket's components, previews and targets are read at one block, by asset", async () => {
    const { client, calls } = chain({
      ...components,
      previewMint: [4n, 2n],
      previewRedeem: [3n, 1n],
      target: [10n ** 18n, 5n * 10n ** 17n],
      pendingTarget: [[10n ** 18n, 0n], 1_790_604_800n],
    })
    const mcp = await legacy(client)
    expect((await call(mcp, 'baskets_components', { basket: 'PAIR' })).structured).toEqual({
      blockNumber: '7',
      basket: 'PAIR',
      address: pair,
      components: [
        { asset: 'NVDA', token: nvda },
        { asset: 'SPY', token: spy },
      ],
    })
    expect((await call(mcp, 'baskets_preview_mint', { basket: 'PAIR', shares: '3' })).structured).toEqual({
      blockNumber: '7',
      basket: 'PAIR',
      shares: '3',
      assets: [
        { asset: 'NVDA', token: nvda, amount: '4' },
        { asset: 'SPY', token: spy, amount: '2' },
      ],
    })
    expect((await call(mcp, 'baskets_preview_redeem', { basket: 'PAIR', shares: '3' })).structured).toMatchObject({
      assets: [
        { asset: 'NVDA', amount: '3' },
        { asset: 'SPY', amount: '1' },
      ],
    })
    expect((await call(mcp, 'baskets_target', { basket: 'PAIR' })).structured).toEqual({
      blockNumber: '7',
      basket: 'PAIR',
      target: [
        { asset: 'NVDA', token: nvda, units: '1000000000000000000' },
        { asset: 'SPY', token: spy, units: '500000000000000000' },
      ],
    })
    expect((await call(mcp, 'baskets_pending_target', { basket: 'PAIR' })).structured).toMatchObject({
      target: [
        { asset: 'NVDA', units: '1000000000000000000' },
        { asset: 'SPY', units: '0' },
      ],
      effectiveAt: '1790604800',
    })
    expect(new Set(calls.map(({ to, block }) => `${to}@${block}`))).toEqual(new Set([`${pair}@0x7`]))
    const none = chain({ ...components, pendingTarget: [[], 0n] })
    expect(
      (await call(await legacy(none.client), 'baskets_pending_target', { basket: 'PAIR' })).structured,
    ).toMatchObject({ target: [], effectiveAt: '0' })
  })

  test('what a position holds through its baskets is read by asset, at one block', async () => {
    const { client, calls } = chain({
      stocks: [
        [toBytes32('NVDA'), toBytes32('SPY')],
        [nvda, spy],
      ],
      inBaskets: [10n, 5n],
    })
    const { structured } = await call(await legacy(client), 'accounts_in_baskets', {
      account: alice,
      position: 'CROSS',
    })
    expect(structured).toEqual({ blockNumber: '7', assets: { NVDA: '10', SPY: '5' } })
    expect(calls).toEqual([
      { to: accounts, functionName: 'stocks', args: [], block: '0x7' },
      { to: accounts, functionName: 'inBaskets', args: [alice, CROSS], block: '0x7' },
    ])
  })

  test('a basket outside the registry is named, even one every object inherits', async () => {
    const { client, methods } = chain({})
    const mcp = await legacy(client)
    for (const name of ['NONE', 'constructor', '__proto__']) {
      const { result, text } = await call(mcp, 'baskets_mint', { from: bob, basket: name, shares: '1', receiver: bob })
      expect(result.isError).toBe(true)
      expect(text).toBe(`The registry has no .tapehouse.Baskets.${name}.`)
    }
    expect(methods).toEqual([])
  })

  test("a mint's limit is previewMint, each token approved for exactly its part where the allowance falls short", async () => {
    const { client, calls } = chain({
      ...components,
      previewMint: (args: readonly unknown[]) => {
        expect(args).toEqual([2n * 10n ** 18n])
        return [2n * 10n ** 18n, 10n ** 18n]
      },
      allowance: (_: readonly unknown[], to: string) => (to === nvda ? 10n ** 18n : maxUint256),
    })
    const { structured } = await call(await legacy(client), 'baskets_mint', {
      from: bob,
      basket: 'PAIR',
      shares: '2000000000000000000',
      receiver: alice,
    })
    expect(structured).toMatchObject({
      chainId: 4663,
      signer: bob,
      assets: [
        { asset: 'NVDA', token: nvda, amount: '2000000000000000000' },
        { asset: 'SPY', token: spy, amount: '1000000000000000000' },
      ],
      calls: [
        { to: nvda, functionName: 'approve', args: { spender: pair, amount: '2000000000000000000' } },
        {
          to: pair,
          functionName: 'mint',
          args: {
            shares: '2000000000000000000',
            receiver: alice,
            maxAssets: ['2000000000000000000', '1000000000000000000'],
          },
        },
      ],
    })
    expect(structured.calls).toHaveLength(2)
    const [, minted] = structured.calls as [unknown, { data: Hex }]
    expect(decodeFunctionData({ abi: basketAbi, data: minted.data }).args).toEqual([
      2n * 10n ** 18n,
      alice,
      [2n * 10n ** 18n, 10n ** 18n],
    ])
    expect(new Set(calls.map(({ block }) => block))).toEqual(new Set(['0x7']))
    expect(calls.filter(({ functionName }) => functionName === 'allowance')).toMatchObject([
      { to: nvda, args: [bob, pair] },
      { to: spy, args: [bob, pair] },
    ])
  })

  test('a redemption and an unwrap name what they give and who signs', async () => {
    const { client } = chain({ ...components, previewRedeem: [3n, 1n] })
    const mcp = await legacy(client)
    const { structured } = await call(mcp, 'baskets_redeem', {
      basket: 'PAIR',
      shares: '5',
      receiver: bob,
      owner: alice,
    })
    expect(structured).toMatchObject({
      signer: `${alice}, or an address it allowed to spend the shares`,
      assets: [
        { asset: 'NVDA', amount: '3' },
        { asset: 'SPY', amount: '1' },
      ],
      calls: [{ to: pair, functionName: 'redeem', args: { shares: '5', receiver: bob, owner: alice } }],
    })
    const unwrap = await call(mcp, 'accounts_unwrap', { account: alice, basket: 'PAIR', shares: '5' })
    expect(unwrap.structured).toMatchObject({
      signer: `${alice}, an address it authorized, or anyone once the position falls short`,
      calls: [{ to: accounts, functionName: 'unwrap', args: { account: alice, basket: pair, shares: '5' } }],
    })
    const [{ data }] = unwrap.structured.calls as [{ data: Hex }]
    expect(decodeFunctionData({ abi: marginAccountsAbi, data }).args).toEqual([alice, pair, 5n])
  })

  test("a basket's shares go in and out of the cross position alone", async () => {
    const { client, calls } = chain({ allowance: 0n, collateral: 7n })
    const mcp = await legacy(client)
    const deposit = { from: bob, account: alice, position: 'CROSS', token: 'PAIR', amount: '7' }
    expect((await call(mcp, 'accounts_deposit', deposit)).structured.calls).toMatchObject([
      { to: pair, functionName: 'approve', args: { spender: accounts, amount: '7' } },
      { to: accounts, functionName: 'deposit', args: { position: CROSS, token: pair, amount: '7', account: alice } },
    ])
    const withdrawal = { account: alice, position: 'CROSS', token: 'PAIR', amount: '7', receiver: bob }
    expect((await call(mcp, 'accounts_withdraw', withdrawal)).structured.calls).toMatchObject([
      { to: accounts, functionName: 'withdraw', args: { token: pair, amount: '7', receiver: bob } },
    ])
    const held = await call(mcp, 'accounts_collateral', { account: alice, position: 'CROSS', token: 'PAIR' })
    expect(held.structured).toEqual({ blockNumber: '7', token: 'PAIR', amount: '7' })
    const read = calls.length
    for (const [name, args] of [
      ['accounts_deposit', { ...deposit, position: 'NVDA' }],
      ['accounts_withdraw', { ...withdrawal, position: 'NVDA' }],
      ['accounts_collateral', { account: alice, position: 'NVDA', token: 'PAIR' }],
    ] as const) {
      const { result, text } = await call(mcp, name, args)
      expect(result.isError).toBe(true)
      expect(text).toBe("A basket's shares are held in the cross position alone: position must be CROSS.")
    }
    expect(calls).toHaveLength(read)
  })

  test("a basket's reverts and the accounts' BasketFrozen are explained by name and arguments", async () => {
    const above = encodeErrorResult({ abi: basketAbi, errorName: 'AboveMaximum', args: [nvda, 2n, 1n] })
    const frozen = encodeErrorResult({ abi: marginAccountsAbi, errorName: 'BasketFrozen', args: [pair] })
    const mcp = await legacy(
      chain({ ...components, previewMint: { revert: above }, health: { revert: frozen } }).client,
    )
    expect((await call(mcp, 'baskets_preview_mint', { basket: 'PAIR', shares: '1' })).text).toBe(
      `The call reverted with AboveMaximum(${nvda}, 2, 1).`,
    )
    expect((await call(mcp, 'accounts_health', { account: alice, position: 'CROSS' })).text).toBe(
      `The call reverted with BasketFrozen(${pair}).`,
    )
  })

  test('the registry lists the baskets, and its output schema declares them', async () => {
    const mcp = await legacy(chain({}).client)
    const { structured } = await call(mcp, 'registry', {})
    expect(structured.baskets).toEqual({ PAIR: pair })
    const { tools } = await mcp.listTools()
    expect(tools.find((tool) => tool.name === 'registry')?.outputSchema?.properties).toHaveProperty('baskets')
  })
})

const usdg = robinhood.tokens.USDG as Address
const onSale = { sales: [1_790_985_600_000n, 1_790_985_600_000n] }
const closed = { sales: [0n, 0n] }
const coverOne = {
  covers: [bob, 1_789_776_000_000n, 218, 1_218, toBytes32('NVDA'), 10_000_000_000n, 13_908_962n],
  reopenOf: 1_789_948_800_000n,
}
const settledSeries = (status: number) => ({
  series: [10_000_000_000n, 200_00000000n, 185_00000000n, 0, status, true],
})

describe('gap cover', () => {
  test('the sales, a quote, a series, a cover and the vault are read at one block', async () => {
    const { client, calls } = chain({
      ...onSale,
      quote: 21_416_957n,
      minDeductible: 513n,
      pricingGap: [279_862n, 111_945n],
      capacity: 99_000_000_000n,
      ...coverOne,
      ...settledSeries(1),
      payouts: 5_000_000n,
      held: 100_000_000_000n,
      reserved: 1_000_000_000n,
      premiums: 21_416_957n,
      owed: 5_000_000n,
      outstanding: 1n,
      totalSupply: 100_000_000_000_000_000n,
      balanceOf: 40_000_000_000_000_000n,
      maxDeposit: maxUint256,
      maxRedeem: 0n,
      convertToAssets: 40_000_000_000n,
    })
    const mcp = await legacy(client)
    expect((await call(mcp, 'gap_cover_sales', {})).structured).toEqual({
      blockNumber: '7',
      onSale: true,
      closesMs: '1790985600000',
      endsMs: '1790985600000',
    })
    const layer = { asset: 'NVDA', notional: '10000000000', deductibleBps: 600, limitBps: 1600 }
    expect((await call(mcp, 'gap_cover_quote', layer)).structured).toEqual({
      blockNumber: '7',
      asset: 'NVDA',
      premium: '21416957',
      reserve: '1000000000',
      minDeductibleBps: '513',
      gap: '279862',
      weekMove: '111945',
      capacity: '99000000000',
      onSale: true,
    })
    expect((await call(mcp, 'gap_cover_series', { asset: 'NVDA', closesMs: '1789776000000' })).structured).toEqual({
      blockNumber: '7',
      asset: 'NVDA',
      closesMs: '1789776000000',
      notional: '10000000000',
      status: 'settled',
      referencePrice: '20000000000',
      price: '18500000000',
      flagged: true,
      reopenMs: '1789948800000',
    })
    expect((await call(mcp, 'gap_cover_position', { id: '1' })).structured).toEqual({
      blockNumber: '7',
      id: '1',
      holder: bob,
      asset: 'NVDA',
      closesMs: '1789776000000',
      notional: '10000000000',
      deductibleBps: '218',
      limitBps: '1218',
      premium: '13908962',
      series: {
        notional: '10000000000',
        status: 'settled',
        referencePrice: '20000000000',
        price: '18500000000',
        flagged: true,
        reopenMs: '1789948800000',
      },
      payout: '532000000',
      refund: null,
      credited: '5000000',
    })
    expect((await call(mcp, 'gap_cover_vault', { owner: alice })).structured).toEqual({
      blockNumber: '7',
      held: '100000000000',
      reserved: '1000000000',
      capacity: '99000000000',
      premiums: '21416957',
      owed: '5000000',
      outstanding: '1',
      totalShares: '100000000000000000',
      writer: {
        owner: alice,
        shares: '40000000000000000',
        assets: '40000000000',
        maxDeposit: maxUint256.toString(),
        maxRedeem: '0',
      },
    })
    expect((await call(mcp, 'gap_cover_vault', {})).structured.writer).toBeNull()
    expect(new Set(calls.map(({ to, block }) => `${to}@${block}`))).toEqual(new Set([`${cover}@0x7`]))
    expect(calls.find(({ functionName }) => functionName === 'pricingGap')?.args).toEqual([toBytes32('NVDA')])
  })

  test('a cover released or never bought, and a layer the cover refuses, are explained', async () => {
    const mcp = await legacy(
      chain({
        ...closed,
        covers: [zeroAddress, 0n, 0, 0, toBytes32(''), 0n, 0n],
        quote: {
          revert: encodeErrorResult({ abi: gapCoverAbi, errorName: 'InvalidLayer', args: [500n, 1500n, 513n] }),
        },
        minDeductible: 513n,
        pricingGap: [279_862n, 111_945n],
        capacity: 0n,
      }).client,
    )
    expect((await call(mcp, 'gap_cover_position', { id: '9' })).text).toBe(
      'Cover 9 does not exist or has been released.',
    )
    const refused = await call(mcp, 'gap_cover_quote', {
      asset: 'NVDA',
      notional: '1',
      deductibleBps: 500,
      limitBps: 1500,
    })
    expect(refused.result.isError).toBe(true)
    expect(refused.text).toBe('The call reverted with InvalidLayer(500, 1500, 513).')
    const wide = await call(mcp, 'gap_cover_quote', {
      asset: 'NVDA',
      notional: '1',
      deductibleBps: 0,
      limitBps: 10_001,
    })
    expect(wide.result.isError).toBe(true)
  })

  test('a purchase pays at most the quote, its USDG approved for exactly the premium where the allowance falls short', async () => {
    const args = {
      from: bob,
      asset: 'NVDA',
      notional: '10000000000',
      deductibleBps: 600,
      limitBps: 1600,
      holder: alice,
    }
    const short = chain({ ...onSale, quote: 21_416_957n, allowance: 0n })
    const { structured } = await call(await legacy(short.client), 'gap_cover_buy', args)
    expect(structured).toMatchObject({
      chainId: 4663,
      signer: bob,
      quote: { premium: '21416957', reserve: '1000000000', closesMs: '1790985600000', endsMs: '1790985600000' },
      calls: [
        { to: usdg, functionName: 'approve', args: { spender: cover, amount: '21416957' } },
        {
          to: cover,
          functionName: 'buy',
          args: {
            symbol: toBytes32('NVDA'),
            notional: '10000000000',
            deductibleBps: '600',
            limitBps: '1600',
            maxPremium: '21416957',
            holder: alice,
          },
        },
      ],
    })
    const [, bought] = structured.calls as [unknown, { data: Hex }]
    expect(bought.data.slice(0, 10)).toBe(toFunctionSelector('buy(bytes32,uint256,uint256,uint256,uint256,address)'))
    expect(new Set(short.calls.map(({ block }) => block))).toEqual(new Set(['0x7']))
    const enough = chain({ ...onSale, quote: 21_416_957n, allowance: 21_416_957n })
    const direct = (await call(await legacy(enough.client), 'gap_cover_buy', args)).structured
    expect((direct.calls as { functionName: string }[]).map(({ functionName }) => functionName)).toEqual(['buy'])
  })

  test('outside the sales a purchase is refused with SalesClosed before anything else is read', async () => {
    const { client, calls } = chain(closed)
    const { result, text } = await call(await legacy(client), 'gap_cover_buy', {
      from: bob,
      asset: 'NVDA',
      notional: '1',
      deductibleBps: 600,
      limitBps: 1600,
      holder: bob,
    })
    expect(result.isError).toBe(true)
    expect(text).toBe(
      'No cover is on sale at block 7: the purchase would revert with SalesClosed(). gap_cover_sales reads whether it is.',
    )
    expect(calls.map(({ functionName }) => functionName)).toEqual(['sales'])
  })

  test('a release names what it credits, and a claim what it sends', async () => {
    const settled = await legacy(chain({ ...coverOne, ...settledSeries(1) }).client)
    expect((await call(settled, 'gap_cover_release', { id: '1' })).structured).toMatchObject({
      signer: 'anyone',
      holder: bob,
      payout: '532000000',
      refund: '0',
      calls: [{ to: cover, functionName: 'release', args: { id: '1' } }],
    })
    const voided = await legacy(chain({ ...coverOne, ...settledSeries(2) }).client)
    expect((await call(voided, 'gap_cover_release', { id: '1' })).structured).toMatchObject({
      payout: '0',
      refund: '13908962',
    })
    const open = await legacy(chain({ ...coverOne, ...settledSeries(0) }).client)
    expect((await call(open, 'gap_cover_release', { id: '1' })).text).toBe(
      "Cover 1's series has not settled: the release would revert with NotSettled(1).",
    )
    const gone = await legacy(chain({ covers: [zeroAddress, 0n, 0, 0, toBytes32(''), 0n, 0n] }).client)
    expect((await call(gone, 'gap_cover_release', { id: '2' })).text).toBe(
      'Cover 2 does not exist or has been released: the release would revert with NoCover(2).',
    )
    const owed = await legacy(chain({ payouts: 532_000_000n }).client)
    expect((await call(owed, 'gap_cover_claim', { holder: bob, receiver: alice })).structured).toMatchObject({
      signer: bob,
      amount: '532000000',
      calls: [{ to: cover, functionName: 'claim', args: { receiver: alice } }],
    })
    const none = await legacy(chain({ payouts: 0n }).client)
    expect((await call(none, 'gap_cover_claim', { holder: bob, receiver: bob })).text).toBe(
      `Nothing is credited to ${bob}.`,
    )
  })

  test('writers deposit with an exact approval and redeem, each within what the vault takes now', async () => {
    const open = chain({ maxDeposit: maxUint256, previewDeposit: 2_000_000_000_000_000n, allowance: 0n })
    const deposit = { from: bob, assets: '2000000000', receiver: bob }
    expect((await call(await legacy(open.client), 'gap_cover_deposit', deposit)).structured).toMatchObject({
      signer: bob,
      shares: '2000000000000000',
      calls: [
        { to: usdg, functionName: 'approve', args: { spender: cover, amount: '2000000000' } },
        { to: cover, functionName: 'deposit', args: { assets: '2000000000', receiver: bob } },
      ],
    })
    const shut = await legacy(chain({ maxDeposit: 0n, previewDeposit: 1n }).client)
    expect((await call(shut, 'gap_cover_deposit', deposit)).text).toBe(
      `The vault takes at most 0 now: the deposit would revert with ERC4626ExceededMaxDeposit(${bob}, 2000000000, 0).`,
    )
    const redeem = { shares: '5', receiver: alice, owner: bob }
    const free = await legacy(chain({ maxRedeem: 5n, previewRedeem: 4n }).client)
    expect((await call(free, 'gap_cover_redeem', redeem)).structured).toMatchObject({
      signer: `${bob}, or an address it allowed to spend the shares`,
      assets: '4',
      calls: [{ to: cover, functionName: 'redeem', args: { shares: '5', receiver: alice, owner: bob } }],
    })
    const locked = await legacy(chain({ maxRedeem: 0n, previewRedeem: 4n }).client)
    expect((await call(locked, 'gap_cover_redeem', redeem)).text).toBe(
      `${bob} may redeem at most 0 shares now: the redemption would revert with ERC4626ExceededMaxRedeem(${bob}, 5, 0).`,
    )
  })

  test('a measurement names the week it keeps, and is refused outside the sales and before the week has passed', async () => {
    const ready = chain({ ...onSale, measure: 111_945n })
    expect((await call(await legacy(ready.client), 'gap_cover_measure', { asset: 'NVDA' })).structured).toMatchObject({
      chainId: 4663,
      signer: 'anyone',
      closesMs: '1790985600000',
      weekMove: '111945',
      calls: [{ to: cover, functionName: 'measure', args: { symbol: toBytes32('NVDA') } }],
    })
    expect(ready.calls.map(({ functionName, block }) => [functionName, block])).toEqual([
      ['sales', '0x7'],
      ['measure', '0x7'],
    ])
    const early = encodeErrorResult({ abi: gapCoverAbi, errorName: 'TooEarlyToMeasure', args: [1_790_899_200_000n] })
    const soon = await legacy(chain({ ...onSale, measure: { revert: early } }).client)
    expect((await call(soon, 'gap_cover_measure', { asset: 'NVDA' })).text).toBe(
      'The call reverted with TooEarlyToMeasure(1790899200000).',
    )
    const shut = chain(closed)
    const { result, text } = await call(await legacy(shut.client), 'gap_cover_measure', { asset: 'NVDA' })
    expect(result.isError).toBe(true)
    expect(text).toBe(
      'No cover is on sale at block 7: the measurement would revert with SalesClosed(). gap_cover_sales reads whether it is.',
    )
    expect(shut.calls.map(({ functionName }) => functionName)).toEqual(['sales'])
  })

  test('a registry without the gap cover is named before any call', async () => {
    const { client, methods } = chain({})
    const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair()
    const without = { ...deployments, tapehouse: { Band: band } }
    await createServer({ client, deployments: without, maxSlippageBps: 100, rateLimit: new RateLimit(120) }).connect(
      serverTransport,
    )
    const mcp = new Client({ name: 'test', version: '1.0.0' })
    await mcp.connect(clientTransport)
    open.push(mcp)
    for (const [name, args] of [
      ['gap_cover_vault', {}],
      ['gap_cover_deposit', { from: bob, assets: '1', receiver: bob }],
      ['gap_cover_series', { asset: 'NVDA', closesMs: '1' }],
    ] as const) {
      expect((await call(mcp, name, args)).text).toBe('The registry has no .tapehouse.GapCover.')
    }
    expect(methods).toEqual([])
  })
})

describe('configuration', () => {
  test('comes from the environment, and what it cannot use is refused', () => {
    expect(configFromEnv({ TAPEHOUSE_RPC_URL: 'http://127.0.0.1:8547', TAPEHOUSE_DEPLOYMENTS: 'd.json' })).toEqual({
      rpcUrl: 'http://127.0.0.1:8547',
      deployments: 'd.json',
      maxSlippageBps: 100,
      maxCallsPerMinute: 120,
    })
    expect(
      configFromEnv({
        TAPEHOUSE_RPC_URL: 'https://rpc',
        TAPEHOUSE_DEPLOYMENTS: 'd.json',
        TAPEHOUSE_MAX_SLIPPAGE_BPS: '0',
      }).maxSlippageBps,
    ).toBe(0)
    expect(() => configFromEnv({ TAPEHOUSE_DEPLOYMENTS: 'd.json' })).toThrow('Set TAPEHOUSE_RPC_URL')
    expect(() => configFromEnv({ TAPEHOUSE_RPC_URL: 'ftp://rpc', TAPEHOUSE_DEPLOYMENTS: 'd.json' })).toThrow(
      'TAPEHOUSE_RPC_URL is not an http or https URL.',
    )
    expect(() => configFromEnv({ TAPEHOUSE_RPC_URL: 'https://rpc' })).toThrow('Set TAPEHOUSE_DEPLOYMENTS')
    for (const cap of ['10001', '-1', '1.5', 'x'])
      expect(() =>
        configFromEnv({
          TAPEHOUSE_RPC_URL: 'https://rpc',
          TAPEHOUSE_DEPLOYMENTS: 'd.json',
          TAPEHOUSE_MAX_SLIPPAGE_BPS: cap,
        }),
      ).toThrow('TAPEHOUSE_MAX_SLIPPAGE_BPS is not a whole number of basis points from 0 to 10000.')
    const env = { TAPEHOUSE_RPC_URL: 'https://rpc', TAPEHOUSE_DEPLOYMENTS: 'd.json' }
    expect(configFromEnv({ ...env, TAPEHOUSE_MAX_CALLS_PER_MINUTE: '100000' }).maxCallsPerMinute).toBe(100_000)
    for (const cap of ['0', '100001', '1.5', '-1', 'x'])
      expect(() => configFromEnv({ ...env, TAPEHOUSE_MAX_CALLS_PER_MINUTE: cap })).toThrow(
        'TAPEHOUSE_MAX_CALLS_PER_MINUTE is not a whole number of calls from 1 to 100000.',
      )
  })
})
