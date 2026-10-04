// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from 'node:fs'
import {
  BaseError,
  checksumAddress,
  ContractFunctionRevertedError,
  createClient,
  custom,
  decodeFunctionData,
  encodeAbiParameters,
  encodeErrorResult,
  encodeFunctionData,
  encodeFunctionResult,
  erc20Abi,
  type Hex,
  hashDomain,
  hexToBigInt,
  padHex,
  parseAbi,
  parseAbiParameters,
  parseEventLogs,
  toEventSelector,
  toFunctionSelector,
  recoverTypedDataAddress,
  zeroAddress,
} from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { describe, expect, test } from 'vitest'
import {
  accounts,
  band,
  bandAbi,
  basketAbi,
  baskets,
  bandFeedAbi,
  CROSS,
  decodeRevert,
  decodeRevertData,
  type Deployments,
  engine,
  gapCover,
  gapCoverAbi,
  liquidatorAbi,
  marginAbi,
  marginAccountsAbi,
  morpho,
  morphoBandOracleAbi,
  type PackageSource,
  parseDeployments,
  quoterV2Abi,
  sharePriceFeed,
  shortPositionsAbi,
  shorts,
  sponsorPaymasterAbi,
  sponsorship,
  stockLending,
  stockLendingVaultAbi,
  stockTokenAbi,
  stockTokenRegistryAbi,
  supplyVaultAbi,
  toBytes32,
  tokenPriceFeed,
  usdgAbi,
} from '../src/index.ts'

const registry = (chainId: number): unknown =>
  JSON.parse(readFileSync(new URL(`../../../deployments/${chainId}.json`, import.meta.url), 'utf8'))
const robinhood = parseDeployments(registry(4663))
const arbitrum = parseDeployments(registry(42161))
const testnet = parseDeployments(registry(46630))
const alice = '0x70997970C51812dc3A010C7d01b50e0d17dc79C8'
const bob = '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC'

describe('deployments', () => {
  test('Chainlink feeds price the Stock Token on Robinhood Chain and the share on Arbitrum One', () => {
    expect(tokenPriceFeed(robinhood, 'NVDA_USD')).toEqual({
      address: '0x379EC4f7C378F34a1B47E4F3cbeBCbAC3E8E9F15',
      kind: 'token',
    })
    expect(sharePriceFeed(arbitrum, 'NVDA_USD')).toEqual({
      address: '0x4881A4418b5F2460B21d6F08CD5aA0678a7f262F',
      kind: 'share',
    })
    expect(() => sharePriceFeed(robinhood, 'NVDA_USD')).toThrow(
      '.chainlink.NVDA_USD prices the Stock Token on chain 4663.',
    )
    expect(() => tokenPriceFeed(arbitrum, 'NVDA_USD')).toThrow('.chainlink.NVDA_USD prices the share on chain 42161.')
  })

  test('every group is read, and a missing one is empty', () => {
    expect(testnet.chainId).toBe(46630)
    const raw = registry(46630) as { tapehouse: Record<string, unknown>; bandFeeds: Record<string, string> }
    for (const address of [...Object.values(raw.tapehouse), ...Object.values(raw.bandFeeds)])
      if (typeof address === 'string') expect(address).toBe(checksumAddress(address as Hex))
    expect(robinhood.uniswapV3.QuoterV2).toBe('0x33e885eD0Ec9bF04EcfB19341582aADCb4c8A9E7')
    expect(arbitrum.tokens).toEqual({})
  })

  test('addresses are checksummed, and a wrong checksum or a missing chain is refused', () => {
    const lower = parseDeployments({
      chainId: 412346,
      tapehouse: { Band: '0xa70118d3324d90532e7d2854627b13cace305641' },
    })
    expect(lower.tapehouse.Band).toBe('0xa70118d3324D90532E7D2854627b13CacE305641')
    expect(() =>
      parseDeployments({ chainId: 1, tokens: { USDG: '0xa70118D3324D90532E7D2854627b13CacE305641' } }),
    ).toThrow('.tokens.USDG is not an address.')
    expect(() => parseDeployments({ tokens: {} })).toThrow('The registry has no chainId.')
    expect(() => parseDeployments({ chainId: 1, tokens: { USDG: 1 } })).toThrow('.tokens.USDG is not an address.')
  })

  test("Morpho Blue and the band's Morpho oracles are read from .morpho and .morphoOracles", async () => {
    expect(robinhood.morpho).toEqual({
      Blue: '0x9D53d5E3bd5E8d4Cbfa6DB1ca238AEA02E651010',
      AdaptiveCurveIrm: '0x2BD3d5965B26B51814AC95127B2b80dD6CcC0fa1',
    })
    expect(parseDeployments({ chainId: 412346, morphoOracles: { NVDA: alice } }).morphoOracles).toEqual({ NVDA: alice })
    expect(testnet.morphoOracles).toEqual({})
    await expect(
      morpho.price(createClient({ transport: custom({ request: async () => '0x7' }) }), testnet, 'NVDA'),
    ).rejects.toThrow('The registry has no .morphoOracles.NVDA.')
  })

  test("Morpho Blue's markets are read from .morpho.Markets as 32-byte ids", () => {
    const id = '0x3a85e619751152991742810df6ec69ce473daef99e28a64ab2340d7b7ccfee49'
    const d = parseDeployments({
      chainId: 412346,
      morpho: { Blue: alice, Markets: { NVDA_USDG: `0x${id.slice(2).toUpperCase()}` } },
    })
    expect(d.morpho).toEqual({ Blue: alice })
    expect(d.morphoMarkets).toEqual({ NVDA_USDG: id })
    expect(robinhood.morphoMarkets).toEqual({})
    for (const wrong of [alice, `${id}00`, id.slice(0, 65), 7])
      expect(() => parseDeployments({ chainId: 1, morpho: { Markets: { NVDA_USDG: wrong } } })).toThrow(
        '.morpho.Markets.NVDA_USDG is not a 32-byte id.',
      )
  })

  test('a name every object inherits is not an entry of the registry', () => {
    for (const name of ['constructor', '__proto__', 'toString']) {
      expect(() => tokenPriceFeed(robinhood, name)).toThrow(`The registry has no .chainlink.${name}.`)
      expect(() => band.seal(robinhood, name)).toThrow(`The registry has no .bandFeeds.${name}.`)
    }
  })

  test('the stock lending vaults are read from .tapehouse.StockLending', () => {
    const d = parseDeployments({
      chainId: 412346,
      tapehouse: { MarginAccounts: alice, StockLending: { SPY: bob } },
    })
    expect(d.tapehouse).toEqual({ MarginAccounts: alice })
    expect(d.stockLending).toEqual({ SPY: bob })
  })

  test('the baskets are read from .tapehouse.Baskets', () => {
    const d = parseDeployments({
      chainId: 412346,
      tapehouse: { MarginAccounts: alice, Baskets: { PAIR: bob.toLowerCase() } },
    })
    expect(d.tapehouse).toEqual({ MarginAccounts: alice })
    expect(d.baskets).toEqual({ PAIR: bob })
    expect(testnet.baskets).toEqual({})
    expect(() => parseDeployments({ chainId: 1, tapehouse: { Baskets: { PAIR: 1 } } })).toThrow(
      '.tapehouse.Baskets.PAIR is not an address.',
    )
  })
})

describe('transactions', () => {
  test('a sale and a cover carry the limits a quote gives', () => {
    const sale = shorts.sell(testnetWithShorts(), {
      asset: 'SPY',
      amount: 10n ** 18n,
      minProceeds: 770n,
      account: alice,
    })
    const data = encodeFunctionData(sale)
    expect(data.slice(0, 10)).toBe(toFunctionSelector('sell(bytes32,uint256,uint256,address)'))
    expect(decodeFunctionData({ abi: shortPositionsAbi, data }).args).toEqual([
      toBytes32('SPY'),
      10n ** 18n,
      770n,
      alice,
    ])
    expect(sale.address).toBe(bob)
    const cover = shorts.cover(testnetWithShorts(), { asset: 'SPY', amount: 1n, maxCost: 2n, account: alice })
    expect(encodeFunctionData(cover).slice(0, 10)).toBe(toFunctionSelector('cover(bytes32,uint256,uint256,address)'))
  })

  test('an authorization lets one address act for the caller', () => {
    const call = accounts.setAuthorization(testnet, bob, true)
    expect(call.address).toBe(testnet.tapehouse.MarginAccounts)
    const data = encodeFunctionData(call)
    expect(data.slice(0, 10)).toBe(toFunctionSelector('setAuthorization(address,bool)'))
    expect(decodeFunctionData({ abi: marginAccountsAbi, data }).args).toEqual([bob, true])
  })

  test('a repayment and a withdrawal name the position', () => {
    const repay = accounts.repay(testnet, { position: CROSS, assets: 5n, account: alice })
    expect(decodeFunctionData({ abi: marginAccountsAbi, data: encodeFunctionData(repay) }).args).toEqual([
      CROSS,
      5n,
      alice,
    ])
    const withdraw = accounts.withdraw(testnet, {
      position: toBytes32('NVDA'),
      token: zeroAddress,
      amount: 1n,
      account: alice,
      receiver: bob,
    })
    expect(encodeFunctionData(withdraw).slice(0, 10)).toBe(
      toFunctionSelector('withdraw(bytes32,address,uint256,address,address)'),
    )
  })

  test('a basket mints and redeems in kind, and the accounts unwrap one', () => {
    const d = withBasket()
    const mint = baskets.mint(d, { basket: 'PAIR', shares: 10n ** 18n, receiver: alice, maxAssets: [3n, 4n] })
    expect(mint.address).toBe(bob)
    const minted = encodeFunctionData(mint)
    expect(minted.slice(0, 10)).toBe(toFunctionSelector('mint(uint256,address,uint256[])'))
    expect(decodeFunctionData({ abi: basketAbi, data: minted }).args).toEqual([10n ** 18n, alice, [3n, 4n]])
    const redeem = baskets.redeem(d, { basket: 'PAIR', shares: 5n, receiver: bob, owner: alice })
    expect(decodeFunctionData({ abi: basketAbi, data: encodeFunctionData(redeem) }).args).toEqual([5n, bob, alice])
    const unwrap = accounts.unwrap(d, { account: alice, basket: 'PAIR', shares: 5n })
    expect(unwrap.address).toBe(testnet.tapehouse.MarginAccounts)
    const data = encodeFunctionData(unwrap)
    expect(data.slice(0, 10)).toBe(toFunctionSelector('unwrap(address,address,uint256)'))
    expect(decodeFunctionData({ abi: marginAccountsAbi, data }).args).toEqual([alice, bob, 5n])
    expect(() => baskets.redeem(testnet, { basket: 'PAIR', shares: 1n, receiver: bob, owner: alice })).toThrow(
      'The registry has no .tapehouse.Baskets.PAIR.',
    )
  })

  test('a contract missing from the registry is named', () => {
    const empty = parseDeployments({ chainId: 1 })
    expect(() => shorts.mark(empty, 'SPY')).toThrow('The registry has no .tapehouse.ShortPositions.')
    expect(() => band.seal(empty, 'NVDA')).toThrow('The registry has no .bandFeeds.NVDA.')
  })

  test('signed RedStone packages come from the source the integrator plugs in', async () => {
    const packages = readFileSync(
      new URL('../../../stylus/contracts/band/testdata/nvda-24_7.hex', import.meta.url),
      'utf8',
    ).trim() as Hex
    const asked: string[][] = []
    const source: PackageSource = {
      payload: async (feedIds) => {
        asked.push([...feedIds])
        return packages
      },
    }
    const call = await band.writePrices(testnet, source, ['NVDA---24_7'])
    expect(asked).toEqual([['NVDA---24_7']])
    expect(call.address).toBe(testnet.tapehouse.Band)
    const { functionName, args } = decodeFunctionData({ abi: bandAbi, data: encodeFunctionData(call) })
    expect(functionName).toBe('writePrices')
    expect(args).toEqual([[toBytes32('NVDA---24_7')], packages])
  })
})

describe('errors', () => {
  test("Tapehouse's, the band's, the issuer's and the router's reverts are decoded", () => {
    const unauthorized = encodeErrorResult({ abi: shortPositionsAbi, errorName: 'Unauthorized', args: [bob, alice] })
    expect(decodeRevertData(unauthorized)).toEqual({ errorName: 'Unauthorized', args: [bob, alice] })
    const window = encodeErrorResult({ abi: bandFeedAbi, errorName: 'NotSealWindow', args: [2, 1790200000000n] })
    expect(decodeRevertData(window)).toEqual({ errorName: 'NotSealWindow', args: [2, 1790200000000n] })
    expect(decodeRevertData(encodeErrorResult({ abi: stockTokenAbi, errorName: 'IsPaused' }))).toEqual({
      errorName: 'IsPaused',
      args: [],
    })
    const old = encodeErrorResult({ abi: bandAbi, errorName: 'TimestampIsTooOld', args: [1790119750n, 1790300000n] })
    expect(decodeRevertData(old)).toEqual({ errorName: 'TimestampIsTooOld', args: [1790119750n, 1790300000n] })
    const past = encodeErrorResult({ abi: basketAbi, errorName: 'PastTarget', args: [toBytes32('SPY')] })
    expect(decodeRevertData(past)).toEqual({ errorName: 'PastTarget', args: [toBytes32('SPY')] })
    const frozen = encodeErrorResult({ abi: marginAccountsAbi, errorName: 'BasketFrozen', args: [bob] })
    expect(decodeRevertData(frozen)).toEqual({ errorName: 'BasketFrozen', args: [bob] })
    const capped = encodeErrorResult({ abi: marginAccountsAbi, errorName: 'DebtCapExceeded', args: [2n, 1n] })
    expect(decodeRevertData(capped)).toEqual({ errorName: 'DebtCapExceeded', args: [2n, 1n] })
    const sound = encodeErrorResult({ abi: liquidatorAbi, errorName: 'NotLiquidatable', args: [bob, CROSS] })
    expect(decodeRevertData(sound)).toEqual({ errorName: 'NotLiquidatable', args: [bob, CROSS] })
    const dry = encodeErrorResult({ abi: stockLendingVaultAbi, errorName: 'InsufficientLiquidity', args: [2n, 1n] })
    expect(decodeRevertData(dry)).toEqual({ errorName: 'InsufficientLiquidity', args: [2n, 1n] })
    const router = encodeErrorResult({
      abi: [{ type: 'error', name: 'Error', inputs: [{ type: 'string' }] }],
      errorName: 'Error',
      args: ['Too little received'],
    })
    expect(decodeRevertData(router)).toEqual({ errorName: 'Error', args: ['Too little received'] })
    expect(decodeRevertData('0xdeadbeef')).toBeUndefined()
  })

  test('the revert behind a failed call is decoded', () => {
    const error = new BaseError('the call failed', {
      cause: new ContractFunctionRevertedError({
        abi: [],
        functionName: 'deposit',
        data: encodeErrorResult({ abi: marginAccountsAbi, errorName: 'Unauthorized', args: [bob, alice] }),
      }),
    })
    expect(decodeRevert(error)).toEqual({ errorName: 'Unauthorized', args: [bob, alice] })
    expect(decodeRevert(new Error('not a revert'))).toBeUndefined()
  })
})

describe('quotes', () => {
  const quoter = '0x33e885eD0Ec9bF04EcfB19341582aADCb4c8A9E7'
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
        const [{ to, data }] = params as [{ to: string; data: Hex }]
        if (to === bob) return encodeFunctionResult({ abi: shortPositionsAbi, functionName: 'fee', result: 500 })
        expect(to).toBe(quoter)
        const {
          functionName,
          args: [quote],
        } = decodeFunctionData({ abi: quoterV2Abi, data })
        expect(quote?.fee).toBe(500)
        expect(quote?.tokenIn).toBe(
          functionName === 'quoteExactInputSingle' ? robinhood.tokens.SPY : robinhood.tokens.USDG,
        )
        return encodeFunctionResult({ abi: quoterV2Abi, functionName, result: [77_232_802n, 0n, 0, 0n] })
      },
    }),
  })
  const d: Deployments = { ...robinhood, tapehouse: { ShortPositions: bob } }

  test('a sale takes at least its quote less the slippage, rounded down', async () => {
    expect(await shorts.quoteSale(client, d, { asset: 'SPY', amount: 10n ** 17n, slippageBps: 50n })).toEqual({
      proceeds: 77_232_802n,
      minProceeds: 76_846_637n,
    })
  })

  test('a cover pays at most its quote plus the slippage, rounded up', async () => {
    expect(await shorts.quoteCover(client, d, { asset: 'SPY', amount: 10n ** 17n, slippageBps: 50n })).toEqual({
      cost: 77_232_802n,
      maxCost: 77_618_967n,
    })
  })

  test('a slippage outside 0 to 10000 basis points is refused', async () => {
    await expect(shorts.quoteSale(client, d, { asset: 'SPY', amount: 1n, slippageBps: 10_001n })).rejects.toThrow(
      'slippageBps 10001 is outside 0 to 10000.',
    )
    await expect(shorts.quoteCover(client, d, { asset: 'SPY', amount: 1n, slippageBps: -1n })).rejects.toThrow(
      'slippageBps -1 is outside 0 to 10000.',
    )
  })

  test('a quote that fails is an error, never a zero limit', async () => {
    const failing = createClient({
      transport: custom({
        async request({ params }) {
          const [{ to }] = params as [{ to: string }]
          if (to === bob) return encodeFunctionResult({ abi: shortPositionsAbi, functionName: 'fee', result: 0 })
          return '0x'
        },
      }),
    })
    await expect(shorts.quoteSale(failing, d, { asset: 'SPY', amount: 1n, slippageBps: 50n })).rejects.toThrow(
      'returned no data',
    )
  })
})

describe('account reads', () => {
  const blockTags: string[] = []
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method === 'eth_blockNumber') return '0x7'
        if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
        const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
        expect(to).toBe(testnet.tapehouse.MarginAccounts)
        blockTags.push(block)
        const { functionName, args } = decodeFunctionData({ abi: marginAccountsAbi, data })
        if (functionName === 'liquidationPrice') {
          expect(args).toEqual([alice, CROSS, toBytes32('SPY'), 5n])
          return encodeFunctionResult({ abi: marginAccountsAbi, functionName, result: 69_412_000_000n })
        }
        const result = { debt: 100n, premium: 3n, collateral: 2n * 10n ** 18n, leverage: 25_000n }
        if (!(functionName in result)) throw new Error(`unexpected ${functionName}`)
        return encodeFunctionResult({
          abi: marginAccountsAbi,
          functionName: functionName as keyof typeof result,
          result: result[functionName as keyof typeof result],
        })
      },
    }),
  })

  test('a repayment reads the debt and the premium at one block', async () => {
    expect(await accounts.repayment(client, testnet, alice, CROSS)).toEqual({ debt: 100n, premium: 3n, assets: 103n })
    expect(blockTags.splice(0).map((tag) => hexToBigInt(tag as Hex))).toEqual([7n, 7n])
  })

  test("a position's collateral, leverage and liquidation price are read", async () => {
    expect(await accounts.collateral(client, testnet, alice, CROSS, bob)).toBe(2n * 10n ** 18n)
    expect(await accounts.leverage(client, testnet, alice, CROSS)).toBe(25_000n)
    expect(await accounts.liquidationPrice(client, testnet, alice, CROSS, 'SPY', 5n)).toBe(69_412_000_000n)
  })
})

describe('Morpho oracles', () => {
  const bandAddress = '0xa70118d3324D90532E7D2854627b13CacE305641'
  const d: Deployments = { ...testnet, morphoOracles: { NVDA: bob } }
  const symbol = toBytes32('NVDA')
  const blocks: bigint[] = []
  const chain = ({
    price,
    halt = [false, 0n, 0n, false],
    step = 0,
  }: {
    price: bigint | Hex
    halt?: readonly [boolean, bigint, bigint, boolean]
    step?: number
  }) =>
    createClient({
      transport: custom(
        {
          async request({ method, params }) {
            if (method === 'eth_blockNumber') return '0x7'
            if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
            const [{ to, data }, block] = params as [{ to: string; data: Hex }, Hex | 'latest']
            if (block !== 'latest') blocks.push(hexToBigInt(block))
            if (to === bandAddress) {
              expect(decodeFunctionData({ abi: bandAbi, data })).toEqual({
                functionName: 'corporateAction',
                args: [symbol],
              })
              return encodeFunctionResult({ abi: bandAbi, functionName: 'corporateAction', result: [step, 0n, 0n, 0n] })
            }
            expect(to).toBe(bob)
            const { functionName } = decodeFunctionData({ abi: morphoBandOracleAbi, data })
            if (functionName === 'price' && typeof price !== 'bigint')
              throw Object.assign(new Error('execution reverted'), { code: 3, data: price })
            const result = {
              price,
              halt,
              band: bandAddress,
              symbol,
              collateralToken: alice,
              loanToken: testnet.tokens.USDG,
              scaleFactor: 10n ** 16n,
              owner: alice,
            } as const
            if (!(functionName in result)) throw new Error(`unexpected ${functionName}`)
            return encodeFunctionResult({
              abi: morphoBandOracleAbi,
              functionName: functionName as keyof typeof result,
              result: result[functionName as keyof typeof result],
            })
          },
        },
        { retryCount: 0 },
      ),
    })
  const noAnswer = encodeErrorResult({ abi: morphoBandOracleAbi, errorName: 'NoAnswer', args: [symbol] })

  test('an oracle answers its price at one block', async () => {
    expect(await morpho.price(chain({ price: 765_517_754_420_000_000_000_000_000n }), d, 'NVDA')).toEqual({
      price: 765_517_754_420_000_000_000_000_000n,
    })
    expect(blocks.splice(0)).toEqual([7n])
  })

  test('NoAnswer with no halt, pause or unconfirmed step is a stale band, never a price of zero', async () => {
    expect(await morpho.price(chain({ price: noAnswer }), d, 'NVDA')).toEqual({
      price: undefined,
      noPrice: 'stale',
      revert: { errorName: 'NoAnswer', args: [symbol] },
    })
    expect(blocks.splice(0)).toEqual([7n, 7n, 7n, 7n, 7n])
  })

  test("NoAnswer under a signed halt, the issuer's pause or an unconfirmed step is a halt", async () => {
    for (const halted of [
      chain({ price: noAnswer, halt: [true, 1_790_200_000n, 1_790_196_400n, false] }),
      chain({ price: noAnswer, halt: [false, 0n, 0n, true] }),
      chain({ price: noAnswer, step: 2 }),
    ])
      expect(await morpho.price(halted, d, 'NVDA')).toMatchObject({ price: undefined, noPrice: 'halted' })
    blocks.length = 0
  })

  test('SequencerNotSettled is no price, and any other revert an error', async () => {
    const sequencer = encodeErrorResult({ abi: morphoBandOracleAbi, errorName: 'SequencerNotSettled' })
    expect(await morpho.price(chain({ price: sequencer }), d, 'NVDA')).toEqual({
      price: undefined,
      noPrice: 'sequencerNotSettled',
      revert: { errorName: 'SequencerNotSettled', args: [] },
    })
    await expect(morpho.price(chain({ price: '0x' }), d, 'NVDA')).rejects.toThrow('reverted')
    blocks.length = 0
  })

  test("an oracle's band, symbol, tokens, scale, owner and halt are read", async () => {
    const client = chain({ price: 1n, halt: [true, 1_790_200_000n, 1_790_196_400n, false] })
    expect(await morpho.oracle(client, d, 'NVDA')).toEqual({
      address: bob,
      band: bandAddress,
      symbol,
      collateralToken: alice,
      loanToken: testnet.tokens.USDG,
      scaleFactor: 10n ** 16n,
      owner: alice,
    })
    expect(blocks.splice(0)).toEqual([7n, 7n, 7n, 7n, 7n, 7n])
    expect(await morpho.halt(client, d, 'NVDA')).toEqual({
      signedHalt: true,
      until: 1_790_200_000n,
      issuedAt: 1_790_196_400n,
      oraclePaused: false,
    })
    blocks.length = 0
  })

  test("a re-point's AssetMismatch and BandSet are decoded", () => {
    const mismatch = encodeErrorResult({
      abi: morphoBandOracleAbi,
      errorName: 'AssetMismatch',
      args: [bandAddress, bob],
    })
    expect(decodeRevertData(mismatch)).toEqual({ errorName: 'AssetMismatch', args: [bandAddress, bob] })
    const [set] = parseEventLogs({
      abi: morphoBandOracleAbi,
      logs: [
        {
          address: bob,
          topics: [
            toEventSelector('BandSet(address,address)'),
            padHex(alice, { size: 32 }),
            padHex(bandAddress, { size: 32 }),
          ],
          data: '0x',
          blockHash: null,
          blockNumber: null,
          logIndex: null,
          transactionHash: null,
          transactionIndex: null,
          removed: false,
        },
      ],
    })
    expect(set).toMatchObject({ eventName: 'BandSet', args: { previousBand: alice, newBand: bandAddress } })
  })
})

describe('basket reads', () => {
  const blockTags: string[] = []
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method === 'eth_blockNumber') return '0x9'
        if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
        const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
        blockTags.push(block)
        if (to === testnet.tapehouse.MarginAccounts) {
          const { functionName, args } = decodeFunctionData({ abi: marginAccountsAbi, data })
          if (functionName === 'stocks')
            return encodeFunctionResult({
              abi: marginAccountsAbi,
              functionName,
              result: [
                [toBytes32('NVDA'), toBytes32('SPY')],
                [alice, bob],
              ],
            })
          expect(functionName).toBe('inBaskets')
          expect(args).toEqual([alice, CROSS])
          return encodeFunctionResult({ abi: marginAccountsAbi, functionName: 'inBaskets', result: [10n, 5n] })
        }
        expect(to).toBe(bob)
        const { functionName, args } = decodeFunctionData({ abi: basketAbi, data })
        if (functionName === 'components')
          return encodeFunctionResult({
            abi: basketAbi,
            functionName,
            result: [
              [toBytes32('NVDA'), toBytes32('SPY')],
              [alice, bob],
            ],
          })
        if (functionName === 'pendingTarget')
          return encodeFunctionResult({ abi: basketAbi, functionName, result: [[1n, 2n], 1_790_604_800n] })
        if (functionName === 'target') return encodeFunctionResult({ abi: basketAbi, functionName, result: [7n, 8n] })
        expect(args).toEqual([3n])
        if (functionName === 'previewMint')
          return encodeFunctionResult({ abi: basketAbi, functionName, result: [4n, 2n] })
        if (functionName === 'previewRedeem')
          return encodeFunctionResult({ abi: basketAbi, functionName, result: [3n, 1n] })
        throw new Error(`unexpected ${functionName}`)
      },
    }),
  })
  const d = withBasket()

  test("a basket's components, previews and targets are read", async () => {
    expect(await baskets.components(client, d, 'PAIR')).toEqual({ assets: ['NVDA', 'SPY'], tokens: [alice, bob] })
    expect(await baskets.previewMint(client, d, 'PAIR', 3n)).toEqual([4n, 2n])
    expect(await baskets.previewRedeem(client, d, 'PAIR', 3n)).toEqual([3n, 1n])
    expect(await baskets.target(client, d, 'PAIR')).toEqual([7n, 8n])
    expect(await baskets.pendingTarget(client, d, 'PAIR', { blockNumber: 4n })).toEqual({
      units: [1n, 2n],
      effectiveAt: 1_790_604_800n,
    })
    expect(blockTags.splice(0).at(-1)).toBe('0x4')
  })

  test("what a position holds through its baskets is read by asset, at one block", async () => {
    expect(await accounts.inBaskets(client, d, alice, CROSS)).toEqual({ NVDA: 10n, SPY: 5n })
    expect(blockTags.splice(0).map((tag) => hexToBigInt(tag as Hex))).toEqual([9n, 9n])
  })
})

function withBasket(): Deployments {
  return { ...testnet, baskets: { PAIR: bob } }
}

describe('gap cover', () => {
  const d: Deployments = { ...robinhood, tapehouse: { GapCover: bob } }
  const layer = { notional: 10_000_000_000n, deductibleBps: 218n, limitBps: 1_218n }

  test('a purchase carries its layer, its premium limit and its holder', () => {
    const call = gapCover.buy(d, { asset: 'NVDA', ...layer, maxPremium: 13_908_962n, holder: alice })
    expect(call.address).toBe(bob)
    const data = encodeFunctionData(call)
    expect(data.slice(0, 10)).toBe(toFunctionSelector('buy(bytes32,uint256,uint256,uint256,uint256,address)'))
    expect(decodeFunctionData({ abi: gapCoverAbi, data }).args).toEqual([
      toBytes32('NVDA'),
      10_000_000_000n,
      218n,
      1_218n,
      13_908_962n,
      alice,
    ])
  })

  test('a settlement names the last rounds before the close and the reopen', () => {
    const call = gapCover.settle(d, {
      asset: 'SPY',
      closesMs: 1_790_985_600_000n,
      referenceRound: (1n << 64n) | 7n,
      lastRound: (1n << 64n) | 8n,
    })
    expect(decodeFunctionData({ abi: gapCoverAbi, data: encodeFunctionData(call) }).args).toEqual([
      toBytes32('SPY'),
      1_790_985_600_000n,
      (1n << 64n) | 7n,
      (1n << 64n) | 8n,
    ])
    expect(encodeFunctionData(gapCover.observe(d, { asset: 'SPY', closesMs: 1n })).slice(0, 10)).toBe(
      toFunctionSelector('observe(bytes32,uint64)'),
    )
    expect(encodeFunctionData(gapCover.release(d, 3n)).slice(0, 10)).toBe(toFunctionSelector('release(uint256)'))
    expect(encodeFunctionData(gapCover.claim(d, alice)).slice(0, 10)).toBe(toFunctionSelector('claim(address)'))
    expect(encodeFunctionData(gapCover.record(d))).toBe(toFunctionSelector('record()'))
    const measure = gapCover.measure(d, 'SPY')
    expect(measure.address).toBe(bob)
    expect(decodeFunctionData({ abi: gapCoverAbi, data: encodeFunctionData(measure) }).args).toEqual([toBytes32('SPY')])
    expect(encodeFunctionData(gapCover.deposit(d, { assets: 5n, receiver: alice })).slice(0, 10)).toBe(
      toFunctionSelector('deposit(uint256,address)'),
    )
    expect(encodeFunctionData(gapCover.redeem(d, { shares: 5n, receiver: alice, owner: alice })).slice(0, 10)).toBe(
      toFunctionSelector('redeem(uint256,address,address)'),
    )
    expect(() => gapCover.record(parseDeployments({ chainId: 1 }))).toThrow(
      'The registry has no .tapehouse.GapCover.',
    )
  })

  const client = (status: number) =>
    createClient({
      transport: custom({
        async request({ method, params }) {
          if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
          const [{ to, data }] = params as [{ to: string; data: Hex }]
          expect(to).toBe(bob)
          const { functionName, args } = decodeFunctionData({ abi: gapCoverAbi, data })
          if (functionName === 'quote') {
            expect(args).toEqual([toBytes32('NVDA'), 10_000_000_000n, 218n, 1_218n])
            return encodeFunctionResult({ abi: gapCoverAbi, functionName, result: 13_908_962n })
          }
          if (functionName === 'pricingGap')
            return encodeFunctionResult({ abi: gapCoverAbi, functionName, result: [279_862n, 111_945n] })
          if (functionName === 'series')
            return encodeFunctionResult({
              abi: gapCoverAbi,
              functionName,
              result: [10_000_000_000n, 22_244_729_849n, 22_280_257_368n, 0, status, false],
            })
          if (functionName === 'sales')
            return encodeFunctionResult({
              abi: gapCoverAbi,
              functionName,
              result: status === 0 ? [1_790_985_600_000n, 1_790_985_600_000n] : [0n, 0n],
            })
          if (functionName === 'reopenOf')
            return encodeFunctionResult({ abi: gapCoverAbi, functionName, result: 1_789_948_800_000n })
          if (functionName === 'covers')
            return encodeFunctionResult({
              abi: gapCoverAbi,
              functionName,
              result: [zeroAddress, 0n, 0, 0, `0x${'0'.repeat(64)}`, 0n, 0n],
            })
          throw new Error(`unexpected ${functionName}`)
        },
      }),
    })

  test('a quote reads the premium and gives the writers USDG it reserves', async () => {
    expect(await gapCover.quote(client(0), d, 'NVDA', layer)).toEqual({
      premium: 13_908_962n,
      reserve: 1_000_000_000n,
    })
    expect(await gapCover.pricingGap(client(0), d, 'NVDA')).toEqual({ gap: 279_862n, weekMove: 111_945n })
  })

  test('the sales read the closure on sale and when they end, and nothing while none is sold', async () => {
    expect(await gapCover.sales(client(0), d)).toEqual({
      closesMs: 1_790_985_600_000n,
      endsMs: 1_790_985_600_000n,
    })
    expect(await gapCover.sales(client(1), d)).toBeUndefined()
  })

  test('a series reads where it stands and its settlement once settled', async () => {
    expect(await gapCover.series(client(1), d, 'NVDA', 1_789_776_000_000n)).toEqual({
      notional: 10_000_000_000n,
      status: 'settled',
      settlement: { referencePrice: 22_244_729_849n, price: 22_280_257_368n },
      flagged: false,
      reopenMs: 1_789_948_800_000n,
    })
    const open = await gapCover.series(client(0), d, 'NVDA', 1_789_776_000_000n)
    expect(open.status).toBe('open')
    expect(open.settlement).toBeUndefined()
    expect((await gapCover.series(client(2), d, 'NVDA', 1_789_776_000_000n)).status).toBe('void')
    await expect(gapCover.series(client(3), d, 'NVDA', 1_789_776_000_000n)).rejects.toThrow('Unknown series status 3.')
    expect(await gapCover.cover(client(0), d, 1n)).toBeUndefined()
  })

  test('a cover pays the fall beyond its deductible up to its limit', () => {
    const settled = { referencePrice: 200_00000000n, price: 185_00000000n }
    const fivePercent = { notional: 10_000_000_000n, deductibleBps: 500n, limitBps: 1_500n }
    expect(gapCover.payout(fivePercent, settled)).toBe(250_000_000n)
    expect(gapCover.payout(fivePercent, { ...settled, price: 190_00000000n })).toBe(0n)
    expect(gapCover.payout(fivePercent, { ...settled, price: 210_00000000n })).toBe(0n)
    expect(gapCover.payout(fivePercent, { ...settled, price: 100_00000000n })).toBe(1_000_000_000n)
  })

  test("the cover's reverts are decoded", () => {
    expect(decodeRevertData(encodeErrorResult({ abi: gapCoverAbi, errorName: 'SalesClosed' }))).toEqual({
      errorName: 'SalesClosed',
      args: [],
    })
    const full = encodeErrorResult({ abi: gapCoverAbi, errorName: 'NoCapacity', args: [2n, 1n] })
    expect(decodeRevertData(full)).toEqual({ errorName: 'NoCapacity', args: [2n, 1n] })
    const stale = encodeErrorResult({ abi: gapCoverAbi, errorName: 'StaleReference', args: [267n, 268n] })
    expect(decodeRevertData(stale)).toEqual({ errorName: 'StaleReference', args: [267n, 268n] })
    const early = encodeErrorResult({ abi: gapCoverAbi, errorName: 'TooEarlyToMeasure', args: [1_790_899_200_000n] })
    expect(decodeRevertData(early)).toEqual({ errorName: 'TooEarlyToMeasure', args: [1_790_899_200_000n] })
  })
})

function testnetWithShorts(): Deployments {
  return { ...testnet, tapehouse: { ...testnet.tapehouse, ShortPositions: bob } }
}

describe('sponsorship', () => {
  const d: Deployments = { ...testnet, tapehouse: { ...testnet.tapehouse, SponsorPaymaster: bob } }

  test('smart accounts are SimpleAccount v0.7 on the registry\'s EntryPoint', () => {
    expect(sponsorship.entryPoint(robinhood)).toEqual({
      address: '0x0000000071727De22E5E9d8BAf0edAc6f37da032',
      version: '0.7',
    })
    expect(sponsorship.accountFactory(testnet)).toBe('0x91E60e0613810449d098b0b5Ec8b51A0FE8c8985')
    expect(() => sponsorship.entryPoint(arbitrum)).toThrow('The registry has no .erc4337.EntryPoint.')
  })

  test('the paymaster is read from .tapehouse.SponsorPaymaster', () => {
    expect(sponsorship.paymaster(d)).toBe(bob)
    expect(() => sponsorship.paymaster(robinhood)).toThrow('The registry has no .tapehouse.SponsorPaymaster.')
  })

  test('a transaction becomes a call of a user operation', () => {
    const call = sponsorship.asCall(accounts.setAuthorization(d, alice, true))
    expect(call.to).toBe(d.tapehouse.MarginAccounts)
    expect(encodeFunctionData(call)).toBe(encodeFunctionData(accounts.setAuthorization(d, alice, true)))
  })

  test("an account's calls are estimated cold, the L1 gas taken off", async () => {
    const client = createClient({
      transport: custom({
        async request({ method, params }) {
          if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
          const [{ from, to, data }] = params as [{ from: string; to: string; data: Hex }]
          expect(from.toLowerCase()).toBe(sponsorship.entryPoint(robinhood).address.toLowerCase())
          expect(to.toLowerCase()).toBe('0x00000000000000000000000000000000000000c8')
          expect(data.slice(0, 10)).toBe(toFunctionSelector('gasEstimateComponents(address,bool,bytes)'))
          return encodeAbiParameters(parseAbiParameters('uint64, uint64, uint256, uint256'), [3_068_532n, 2_444_200n, 1n, 1n])
        },
      }),
    })
    expect(await sponsorship.callGas(client, robinhood, alice, '0xb61d27f6')).toBe(624_332n)
  })

  test("an owner's account is read from the factory, salt 0", async () => {
    const client = createClient({
      transport: custom({
        async request({ params }) {
          const [{ to, data }] = params as [{ to: string; data: Hex }]
          expect(to).toBe(sponsorship.accountFactory(testnet))
          expect(data).toBe(
            encodeFunctionData({
              abi: parseAbi(['function getAddress(address,uint256) view returns (address)']),
              functionName: 'getAddress',
              args: [alice, 0n],
            }),
          )
          return padHex(bob, { size: 32 })
        },
      }),
    })
    expect(await sponsorship.accountAddress(client, testnet, alice)).toBe(bob)
  })

  test('the free operations left are read from the paymaster', async () => {
    const client = createClient({
      transport: custom({
        async request({ method, params }) {
          if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
          const [{ to, data }] = params as [{ to: string; data: Hex }]
          expect(to).toBe(bob)
          const { functionName, args } = decodeFunctionData({ abi: sponsorPaymasterAbi, data })
          expect(functionName).toBe('freeOperationsLeft')
          expect(args).toEqual([alice])
          return encodeFunctionResult({ abi: sponsorPaymasterAbi, functionName: 'freeOperationsLeft', result: 2n })
        },
      }),
    })
    expect(await sponsorship.freeOperationsLeft(client, d, alice)).toBe(2n)
  })
})

const nvda = '0x90F79bf6EB2c4f870365E785982E1f101E93b906'
const withStocks: Deployments = { ...testnet, tokens: { ...testnet.tokens, NVDA: nvda }, stockLending: { NVDA: bob } }

describe('borrowing', () => {
  const testnet = withStocks
  const signature =
    '0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c' as Hex

  test('a deposit with a permit carries its signature as v, r and s', () => {
    const call = accounts.depositWithPermit(testnet, {
      position: CROSS,
      token: nvda,
      amount: 3n,
      account: alice,
      deadline: 9n,
      signature,
    })
    const { functionName, args } = decodeFunctionData({ abi: marginAccountsAbi, data: encodeFunctionData(call) })
    expect(functionName).toBe('depositWithPermit')
    expect(args).toEqual([
      CROSS,
      nvda,
      3n,
      alice,
      9n,
      28,
      '0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf',
      '0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8',
    ])
  })

  test('ether is wrapped into WETH for the account that deposits it', () => {
    const call = accounts.wrap(testnet, { account: alice, value: 7n })
    expect(call.address).toBe(testnet.tokens.WETH)
    expect(call.value).toBe(7n)
    expect(encodeFunctionData(call)).toBe(
      `${toFunctionSelector('depositTo(address)')}${alice.slice(2).toLowerCase().padStart(64, '0')}`,
    )
  })

  test('a position lends, takes back, recalls and settles its Stock Tokens', () => {
    const at = { position: toBytes32('NVDA'), token: nvda, amount: 5n, account: alice } as const
    for (const [call, signature] of [
      [accounts.lend(testnet, at), 'lend(bytes32,address,uint256,address)'],
      [accounts.unlend(testnet, at), 'unlend(bytes32,address,uint256,address)'],
      [accounts.recall(testnet, at), 'recall(bytes32,address,uint256,address)'],
    ] as const) {
      const data = encodeFunctionData(call)
      expect(data.slice(0, 10)).toBe(toFunctionSelector(signature))
      expect(decodeFunctionData({ abi: marginAccountsAbi, data }).args).toEqual([toBytes32('NVDA'), nvda, 5n, alice])
    }
    const settle = accounts.settle(testnet, { account: alice, position: CROSS, token: nvda })
    expect(decodeFunctionData({ abi: marginAccountsAbi, data: encodeFunctionData(settle) }).args).toEqual([
      alice,
      CROSS,
      nvda,
    ])
  })

  const blockTags: string[] = []
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method === 'eth_blockNumber') return '0xb'
        if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
        const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
        blockTags.push(block)
        if (to === testnet.tapehouse.SupplyVault) {
          return encodeFunctionResult({ abi: supplyVaultAbi, functionName: 'debt', result: 412n })
        }
        if (to === nvda) {
          expect(decodeFunctionData({ abi: erc20Abi, data }).args).toEqual([testnet.tapehouse.MarginAccounts])
          return encodeFunctionResult({ abi: erc20Abi, functionName: 'balanceOf', result: 14n * 10n ** 18n })
        }
        if (to === testnet.tapehouse.Margin) {
          const { functionName, args } = decodeFunctionData({ abi: marginAbi, data })
          if (functionName === 'requirement') {
            expect(args).toEqual([[10n ** 18n], [17_120_000_000n], 172_800n, true])
            return encodeFunctionResult({ abi: marginAbi, functionName, result: [608n, 1] })
          }
          expect(functionName).toBe('currentRequirement')
          return encodeFunctionResult({ abi: marginAbi, functionName: 'currentRequirement', result: [500n, 0, 3] })
        }
        if (to === bob) {
          const { functionName, args } = decodeFunctionData({ abi: stockLendingVaultAbi, data })
          if (functionName === 'head') return encodeFunctionResult({ abi: stockLendingVaultAbi, functionName, result: 4n })
          expect(args).toEqual([6n])
          if (functionName === 'claimable')
            return encodeFunctionResult({ abi: stockLendingVaultAbi, functionName, result: 7n })
          expect(functionName).toBe('ticket')
          return encodeFunctionResult({
            abi: stockLendingVaultAbi,
            functionName: 'ticket',
            result: [10n, 20n, 3n, 1_790_000_000n],
          })
        }
        expect(to).toBe(testnet.tapehouse.MarginAccounts)
        const { functionName, args } = decodeFunctionData({ abi: marginAccountsAbi, data })
        const result = {
          stocks: [
            [toBytes32('NVDA'), toBytes32('TSLA'), toBytes32('SPY')],
            [nvda, zeroAddress, bob],
          ],
          lent: 2n,
          sellable: 9n,
          claim: 1n,
          recalls: [6n],
          holding: [16n * 10n ** 18n, 9n * 10n ** 17n, 100n * 10n ** 18n],
          debtCap: 1_000n,
          weekendDebtCap: 500n,
          weekendLeverage: 50_000,
          borrowingPaused: false,
          premiumRate: 500,
        }
        if (functionName === 'holding') expect(args).toEqual([toBytes32('NVDA')])
        if (['lent', 'sellable', 'claim', 'recalls'].includes(functionName)) expect(args).toEqual([alice, CROSS, nvda])
        if (!(functionName in result)) throw new Error(`unexpected ${functionName}`)
        return encodeFunctionResult({
          abi: marginAccountsAbi,
          functionName: functionName as keyof typeof result,
          result: result[functionName as keyof typeof result],
        } as never)
      },
    }),
  })

  test("an account's positions are the cross one and each Stock Token's isolated one", async () => {
    expect(await accounts.stocks(client, testnet)).toEqual([
      { asset: 'NVDA', token: nvda },
      { asset: 'TSLA', token: zeroAddress },
      { asset: 'SPY', token: bob },
    ])
    expect(await accounts.positions(client, testnet)).toEqual([CROSS, toBytes32('NVDA'), toBytes32('SPY')])
    blockTags.splice(0)
  })

  test('what a position lent, may sell, recalled and waits for is read at one block', async () => {
    expect(await accounts.lending(client, testnet, alice, CROSS, nvda)).toEqual({
      lent: 2n,
      sellable: 9n,
      claim: 1n,
      recalls: [6n],
    })
    expect(blockTags.splice(0).map((tag) => hexToBigInt(tag as Hex))).toEqual([11n, 11n, 11n, 11n])
  })

  test("an asset's holding reads its quantity, its cap and the tokens the accounts hold", async () => {
    expect(await accounts.holding(client, testnet, 'NVDA', nvda)).toEqual({
      units: 16n * 10n ** 18n,
      scale: 9n * 10n ** 17n,
      quantity: 144n * 10n ** 17n,
      cap: 100n * 10n ** 18n,
      balance: 14n * 10n ** 18n,
    })
    expect(new Set(blockTags.splice(0))).toEqual(new Set(['0xb']))
  })

  test("the accounts' limits on borrowing are read at one block", async () => {
    expect(await accounts.limits(client, testnet)).toEqual({
      debt: 412n,
      debtCap: 1_000n,
      weekendDebtCap: 500n,
      weekendLeverage: 50_000,
      borrowingPaused: false,
      premiumRate: 500,
    })
    expect(blockTags.splice(0)).toHaveLength(6)
  })

  test('the engine margins any portfolio, over two days by default', async () => {
    const portfolio = { quantities: [10n ** 18n], prices: [17_120_000_000n] }
    expect(await engine.requirement(client, testnet, portfolio, { spansClosure: true })).toEqual({
      margin: 608n,
      missing: 1,
    })
    expect(await engine.currentRequirement(client, testnet, portfolio)).toEqual({ margin: 500n, missing: 0, regime: 3 })
    blockTags.splice(0)
  })

  test('a recall ticket reads its span, its notice, what waits for it and the head, at one block', async () => {
    expect(await stockLending.ticket(client, testnet, 'NVDA', 6n)).toEqual({
      start: 10n,
      end: 20n,
      taken: 3n,
      dueAt: 1_790_000_000n,
      claimable: 7n,
      head: 4n,
    })
    expect(blockTags.splice(0)).toEqual(['0xb', '0xb', '0xb'])
  })
})

describe('issuer', () => {
  test("a Stock Token's pause and its registry's blocks are read at one block", async () => {
    const registry = '0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65'
    const tags: string[] = []
    const client = createClient({
      transport: custom({
        async request({ method, params }) {
          if (method === 'eth_blockNumber') return '0xc'
          const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
          tags.push(block)
          if (to === registry) {
            const { args } = decodeFunctionData({ abi: stockTokenRegistryAbi, data })
            return encodeFunctionResult({ abi: stockTokenRegistryAbi, functionName: 'isBlocked', result: args[0] === bob })
          }
          const { functionName } = decodeFunctionData({ abi: stockTokenAbi, data })
          if (functionName === 'paused') return encodeFunctionResult({ abi: stockTokenAbi, functionName, result: true })
          expect(functionName).toBe('ACCESS_CONTROLLED_REGISTRY')
          return encodeFunctionResult({ abi: stockTokenAbi, functionName: 'ACCESS_CONTROLLED_REGISTRY', result: registry })
        },
      }),
    })
    expect(await accounts.issuer(client, nvda, [alice, bob])).toEqual({ paused: true, blocked: { [alice]: false, [bob]: true } })
    expect(new Set(tags)).toEqual(new Set(['0xc']))
  })
})

describe('permits', () => {
  const owner = privateKeyToAccount('0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d')
  const token = nvda
  const usdg = testnet.tokens.USDG as `0x${string}`
  const domainOf = (address: `0x${string}`) =>
    ({ name: address === token ? 'NVIDIA • Robinhood Token' : 'Global Dollar', version: '1', chainId: 46630, verifyingContract: address }) as const
  const chain = (separator?: Hex) =>
    createClient({
      transport: custom({
        async request({ method, params }) {
          if (method === 'eth_chainId') return '0xb626'
          const [{ to, data }] = params as [{ to: `0x${string}`; data: Hex }]
          const selector = data.slice(0, 10)
          if (selector === toFunctionSelector('nonces(address)'))
            return encodeFunctionResult({ abi: usdgAbi, functionName: 'nonces', result: 4n })
          if (selector === toFunctionSelector('eip712Domain()')) {
            if (to !== token) return '0x'
            const d = domainOf(to)
            return encodeFunctionResult({
              abi: stockTokenAbi,
              functionName: 'eip712Domain',
              result: ['0x0f', d.name, d.version, 46630n, to, padHex('0x', { size: 32 }), []],
            })
          }
          if (selector === toFunctionSelector('name()'))
            return encodeFunctionResult({ abi: erc20Abi, functionName: 'name', result: domainOf(to).name })
          return encodeFunctionResult({
            abi: usdgAbi,
            functionName: 'DOMAIN_SEPARATOR',
            result:
              separator ??
              hashDomain({
                domain: { ...domainOf(to), chainId: 46630n },
                types: {
                  EIP712Domain: [
                    { name: 'name', type: 'string' },
                    { name: 'version', type: 'string' },
                    { name: 'chainId', type: 'uint256' },
                    { name: 'verifyingContract', type: 'address' },
                  ],
                },
              }),
          })
        },
      }),
    })
  const request = { owner: owner.address, account: bob, value: 5n, deadline: 99n } as const

  test("a Stock Token's permit is signed in the domain its eip712Domain reads", async () => {
    const typedData = await sponsorship.permit(chain(), { token, ...request })
    expect(typedData.domain).toEqual(domainOf(token))
    expect(typedData.message).toEqual({ owner: owner.address, spender: bob, value: 5n, nonce: 4n, deadline: 99n })
    const signature = await owner.signTypedData(typedData)
    expect(await recoverTypedDataAddress({ ...typedData, signature })).toBe(owner.address)
  })

  test("USDG's permit is signed in its name and version 1, checked against its domain separator", async () => {
    expect((await sponsorship.permit(chain(), { token: usdg, ...request })).domain).toEqual(domainOf(usdg))
    await expect(sponsorship.permit(chain(padHex('0x1', { size: 32 })), { token: usdg, ...request })).rejects.toThrow(
      'permit domain does not match its DOMAIN_SEPARATOR on chain 46630.',
    )
  })

  test("a Stock Token whose eip712Domain disagrees with its DOMAIN_SEPARATOR gets no permit", async () => {
    await expect(sponsorship.permit(chain(padHex('0x2', { size: 32 })), { token, ...request })).rejects.toThrow(
      'permit domain does not match its DOMAIN_SEPARATOR on chain 46630.',
    )
  })

  test('the account pulls what the permit allows: the permit, then the transfer', async () => {
    const typedData = await sponsorship.permit(chain(), { token, ...request })
    const signature = await owner.signTypedData(typedData)
    const [permit, transfer] = sponsorship.pull({ token, ...request, signature })
    expect(permit.to).toBe(token)
    expect(decodeFunctionData({ abi: usdgAbi, data: encodeFunctionData(permit) }).args?.slice(0, 4)).toEqual([
      owner.address,
      bob,
      5n,
      99n,
    ])
    expect(decodeFunctionData({ abi: erc20Abi, data: encodeFunctionData(transfer) }).args).toEqual([
      owner.address,
      bob,
      5n,
    ])
  })
})
