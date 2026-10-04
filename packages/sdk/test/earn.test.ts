// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from 'node:fs'
import {
  type Abi,
  type Address,
  createClient,
  custom,
  decodeFunctionData,
  encodeFunctionResult,
  erc4626Abi,
  type Hex,
} from 'viem'
import { describe, expect, test } from 'vitest'
import {
  backstop,
  CROSS,
  gapBackstopAbi,
  gapCover,
  gapCoverAbi,
  marginAccountsAbi,
  parseDeployments,
  stockLending,
  stockLendingVaultAbi,
  supply,
  supplyVaultAbi,
  toBytes32,
} from '../src/index.ts'

const testnet = parseDeployments(
  JSON.parse(readFileSync(new URL('../../../deployments/46630.json', import.meta.url), 'utf8')),
)
const alice = '0x70997970C51812dc3A010C7d01b50e0d17dc79C8'
const nvda = '0x90F79bf6EB2c4f870365E785982E1f101E93b906'
const spy = '0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65'
const cover = '0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc'
const lendingVault = '0x976EA74026E726554dB657fA54763abd0C3a0aa9'
const deployments = {
  ...testnet,
  tapehouse: { ...testnet.tapehouse, GapCover: cover },
  stockLending: { NVDA: lendingVault },
} as const
const signature =
  '0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c'

/** A chain that answers each contract's views from `answers`, keyed by address and then function, and logs every read. */
function chain(answers: Record<string, { abi: Abi; results: Record<string, unknown | ((args: readonly unknown[]) => unknown)> }>) {
  const reads: { to: string; functionName: string; args: readonly unknown[]; block: string }[] = []
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method === 'eth_blockNumber') return '0x2a'
        if (method !== 'eth_call') throw new Error(`unexpected ${method}`)
        const [{ to, data }, block] = params as [{ to: string; data: Hex }, string]
        const contract = Object.entries(answers).find(([address]) => address.toLowerCase() === to.toLowerCase())?.[1]
        if (!contract) throw new Error(`unexpected call to ${to}`)
        const { functionName, args = [] } = decodeFunctionData({ abi: contract.abi, data })
        reads.push({ to, functionName, args, block })
        const answer = contract.results[functionName]
        if (answer === undefined) throw new Error(`unexpected ${functionName} on ${to}`)
        const result = typeof answer === 'function' ? answer(args) : answer
        return encodeFunctionResult({ abi: contract.abi, functionName, result } as never)
      },
    }),
  })
  return { client, reads }
}

describe('supply', () => {
  const vault = testnet.tapehouse.SupplyVault as Address

  test("a lender's transactions deposit, deposit with a permit, withdraw and redeem", () => {
    expect(supply.deposit(testnet, { assets: 5n, receiver: alice })).toMatchObject({
      address: vault,
      functionName: 'deposit',
      args: [5n, alice],
    })
    const { args } = supply.depositWithPermit(testnet, { assets: 5n, receiver: alice, deadline: 9n, signature })
    expect(args).toEqual([
      5n,
      alice,
      9n,
      28,
      '0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf',
      '0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8',
    ])
    expect(supply.withdraw(testnet, { assets: 3n, receiver: alice, owner: nvda }).args).toEqual([3n, alice, nvda])
    expect(supply.redeem(testnet, { shares: 4n, receiver: alice, owner: nvda })).toMatchObject({
      functionName: 'redeem',
      args: [4n, alice, nvda],
    })
  })

  test('the market reads its rates, its utilisation, what it holds and lent, and whether it takes deposits, at one block', async () => {
    const { client, reads } = chain({
      [vault]: {
        abi: supplyVaultAbi,
        results: {
          supplyRate: 3n * 10n ** 16n,
          borrowRate: 5n * 10n ** 16n,
          utilization: 6n * 10n ** 17n,
          idle: 400n,
          debt: 600n,
          totalAssets: 1_000n,
          maxDeposit: 0n,
        },
      },
    })
    expect(await supply.market(client, testnet)).toEqual({
      supplyRate: 3n * 10n ** 16n,
      borrowRate: 5n * 10n ** 16n,
      utilization: 6n * 10n ** 17n,
      idle: 400n,
      debt: 600n,
      totalAssets: 1_000n,
      open: false,
    })
    expect(new Set(reads.map((r) => r.block))).toEqual(new Set(['0x2a']))
  })

  test("a lender's shares are read with what they are worth and what may leave now", async () => {
    const { client, reads } = chain({
      [vault]: {
        abi: erc4626Abi,
        results: {
          balanceOf: 7n * 10n ** 12n,
          convertToAssets: 7n,
          maxWithdraw: 4n,
          maxRedeem: 4n * 10n ** 12n,
          maxDeposit: 2n ** 256n - 1n,
        },
      },
    })
    expect(await supply.lender(client, testnet, alice)).toEqual({
      shares: 7n * 10n ** 12n,
      assets: 7n,
      maxWithdraw: 4n,
      maxRedeem: 4n * 10n ** 12n,
      maxDeposit: 2n ** 256n - 1n,
    })
    expect(reads.find((r) => r.functionName === 'convertToAssets')?.args).toEqual([7n * 10n ** 12n])
    expect(reads.filter((r) => r.functionName !== 'convertToAssets').every((r) => r.args[0] === alice)).toBe(true)
  })
})

describe('backstop', () => {
  const vault = testnet.tapehouse.GapBackstop as Address
  const accounts = testnet.tapehouse.MarginAccounts as Address

  test("a depositor's transactions deposit, start the cooldown, redeem and claim gains, and anyone claims the premium", () => {
    expect(backstop.deposit(testnet, { assets: 5n, receiver: alice })).toMatchObject({
      address: vault,
      abi: gapBackstopAbi,
      functionName: 'deposit',
      args: [5n, alice],
    })
    expect(backstop.startCooldown(testnet)).toMatchObject({ address: vault, functionName: 'startCooldown' })
    expect(backstop.redeem(testnet, { shares: 4n, receiver: alice, owner: nvda }).args).toEqual([4n, alice, nvda])
    expect(backstop.withdraw(testnet, { assets: 3n, receiver: alice, owner: nvda }).args).toEqual([3n, alice, nvda])
    expect(backstop.claimGains(testnet, nvda)).toMatchObject({ functionName: 'claimGains', args: [nvda] })
    expect(backstop.claim(testnet)).toMatchObject({ address: vault, functionName: 'claim' })
  })

  const results = {
    totalAssets: 9_000n,
    totalSupply: 9_000n * 10n ** 6n,
    closureMs: 1_790_000_000_000n,
    COOLDOWN: 604_800n,
    WITHDRAWAL_WINDOW: 518_400n,
    SETTLEMENT: 86_400n,
    owner: alice,
    balanceOf: (args: readonly unknown[]) => (args[0] === alice ? 1_000n * 10n ** 6n : 300n * 10n ** 6n),
    convertToAssets: (args: readonly unknown[]) => (args[0] as bigint) / 10n ** 6n,
    exposureLimit: (args: readonly unknown[]) => (args[0] === CROSS ? 5_000n : 1_000n),
    exposureLimits: (args: readonly unknown[]) =>
      args[0] === CROSS ? [5_000n, 2_000n, 1_791_000_000_000n] : [1_000n, 1_000n, 0n],
    exposureLeft: (args: readonly unknown[]) => (args[0] === CROSS ? 4_500n : 1_000n),
    covered: (args: readonly unknown[]) => (args[0] === CROSS ? 500n : 0n),
    cooldowns: [300n * 10n ** 6n, 1_789_000_000n],
    maxRedeem: 0n,
    maxDeposit: 0n,
    maxWithdraw: 0n,
    gains: (args: readonly unknown[]) => (args[1] === nvda ? 2n * 10n ** 17n : 0n),
  }
  const answers = {
    [vault]: { abi: gapBackstopAbi, results },
    [accounts]: {
      abi: marginAccountsAbi,
      results: {
        stocks: [
          [toBytes32('NVDA'), toBytes32('SPY')],
          [nvda, spy],
        ],
        premiumRate: 500,
        reserveShare: 1_000,
        backstopPremium: 42n,
      },
    },
  }

  test('the backstop reads its exposure per closure, its premium, its terms and the seed at one block', async () => {
    const { client, reads } = chain(answers)
    expect(await backstop.state(client, testnet)).toEqual({
      held: 9_000n,
      shares: 9_000n * 10n ** 6n,
      closureMs: 1_790_000_000_000n,
      exposures: [
        { position: CROSS, limit: 5_000n, current: 5_000n, next: 2_000n, fromMs: 1_791_000_000_000n, left: 4_500n, covered: 500n },
        { position: toBytes32('NVDA'), limit: 1_000n, current: 1_000n, next: 1_000n, fromMs: 0n, left: 1_000n, covered: 0n },
        { position: toBytes32('SPY'), limit: 1_000n, current: 1_000n, next: 1_000n, fromMs: 0n, left: 1_000n, covered: 0n },
      ],
      premium: { rate: 500, reserveShare: 1_000, waiting: 42n },
      terms: { cooldown: 604_800n, window: 518_400n, settlement: 86_400n },
      seed: { owner: alice, assets: 1_000n },
    })
    expect(reads.find((r) => r.functionName === 'covered')?.args).toEqual([CROSS, 1_790_000_000_000n])
    expect(new Set(reads.map((r) => r.block))).toEqual(new Set(['0x2a']))
  })

  test("a depositor's shares, cooldown and gains per Stock Token are read at one block", async () => {
    const { client, reads } = chain(answers)
    expect(await backstop.depositor(client, testnet, nvda, [nvda, spy])).toEqual({
      shares: 300n * 10n ** 6n,
      assets: 300n,
      maxWithdraw: 0n,
      maxRedeem: 0n,
      maxDeposit: 0n,
      cooldown: { shares: 300n * 10n ** 6n, from: 1_789_604_800n, until: 1_790_123_200n },
      gains: { [nvda]: 2n * 10n ** 17n, [spy]: 0n },
    })
    expect(reads.filter((r) => r.functionName === 'gains').map((r) => r.args)).toEqual([
      [nvda, nvda],
      [nvda, spy],
    ])
    expect(new Set(reads.map((r) => r.block))).toEqual(new Set(['0x2a']))
  })

  test('a depositor that never started a cooldown has none', async () => {
    const { client } = chain({ ...answers, [vault]: { abi: gapBackstopAbi, results: { ...results, cooldowns: [0n, 0n] } } })
    expect((await backstop.depositor(client, testnet, nvda, [])).cooldown).toBeUndefined()
  })
})

describe('stock lending terms', () => {
  test("a vault's rates, what it lent and holds for recalls, and its fixed terms are read at one block", async () => {
    const { client, reads } = chain({
      [lendingVault]: {
        abi: stockLendingVaultAbi,
        results: {
          utilization: 5n * 10n ** 17n,
          borrowRate: 4n * 10n ** 16n,
          supplyRate: 17n * 10n ** 15n,
          locked: 3n,
          idle: 10n,
          debt: 13n,
          borrowable: 9n,
          totalAssets: 26n,
          head: 2n,
          MAX_UTILIZATION: 9_000,
          NOTICE: 86_400n,
          feeShare: 1_500,
        },
      },
    })
    expect(await stockLending.terms(client, deployments, 'NVDA')).toEqual({
      utilization: 5n * 10n ** 17n,
      borrowRate: 4n * 10n ** 16n,
      supplyRate: 17n * 10n ** 15n,
      locked: 3n,
      idle: 10n,
      debt: 13n,
      borrowable: 9n,
      totalAssets: 26n,
      head: 2n,
      maxUtilization: 9_000,
      notice: 86_400n,
      feeShare: 1_500,
    })
    expect(new Set(reads.map((r) => r.block))).toEqual(new Set(['0x2a']))
    await expect(stockLending.terms(client, deployments, 'TSLA')).rejects.toThrow(
      'The registry has no .tapehouse.StockLending.TSLA.',
    )
  })
})

describe('gap cover writers', () => {
  test('a writer withdraws USDG as well as redeeming shares', () => {
    expect(gapCover.withdraw(deployments, { assets: 3n, receiver: alice, owner: nvda })).toMatchObject({
      address: cover,
      functionName: 'withdraw',
      args: [3n, alice, nvda],
    })
  })

  test("the writers' USDG, what covers reserve and what they wait on are read at one block", async () => {
    const { client, reads } = chain({
      [cover]: {
        abi: gapCoverAbi,
        results: {
          totalAssets: 1_000n,
          totalSupply: 1_000n * 10n ** 6n,
          reserved: 400n,
          capacity: 600n,
          premiums: 12n,
          owed: 3n,
          outstanding: 2n,
          sales: [1_790_000_000_000n, 1_789_990_000_000n],
        },
      },
    })
    expect(await gapCover.writers(client, deployments)).toEqual({
      held: 1_000n,
      shares: 1_000n * 10n ** 6n,
      reserved: 400n,
      capacity: 600n,
      premiums: 12n,
      owed: 3n,
      outstanding: 2n,
      sales: { closesMs: 1_790_000_000_000n, endsMs: 1_789_990_000_000n },
    })
    expect(new Set(reads.map((r) => r.block))).toEqual(new Set(['0x2a']))
  })

  test("a writer's shares are read with what they are worth and what may leave now", async () => {
    const { client } = chain({
      [cover]: {
        abi: erc4626Abi,
        results: { balanceOf: 5n * 10n ** 6n, convertToAssets: 5n, maxWithdraw: 0n, maxRedeem: 0n, maxDeposit: 0n },
      },
    })
    expect(await gapCover.writer(client, deployments, alice)).toEqual({
      shares: 5n * 10n ** 6n,
      assets: 5n,
      maxWithdraw: 0n,
      maxRedeem: 0n,
      maxDeposit: 0n,
    })
  })
})
