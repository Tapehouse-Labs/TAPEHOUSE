// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Address, Client } from 'viem'
import { readContract, simulateContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry, toBytes32 } from './deployments.js'
import { quoterV2Abi, shortPositionsAbi } from './generated.js'

const BPS = 10_000n

/** Adds `amount` of the caller's USDG to `account`'s short of `asset`. */
export function deposit(
  deployments: Deployments,
  { asset, amount, account }: { asset: string; amount: bigint; account: Address },
) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'deposit',
    args: [toBytes32(asset), amount, account],
  } as const
}

/** Sends `amount` of USDG from `account`'s short of `asset` to `receiver`. */
export function withdraw(
  deployments: Deployments,
  { asset, amount, account, receiver }: { asset: string; amount: bigint; account: Address; receiver: Address },
) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'withdraw',
    args: [toBytes32(asset), amount, account, receiver],
  } as const
}

/** Borrows `amount` of `asset` and sells it for at least `minProceeds` of USDG, from `quoteSale`. */
export function sell(
  deployments: Deployments,
  { asset, amount, minProceeds, account }: { asset: string; amount: bigint; minProceeds: bigint; account: Address },
) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'sell',
    args: [toBytes32(asset), amount, minProceeds, account],
  } as const
}

/** Buys back `amount` of `asset` for at most `maxCost` of USDG, from `quoteCover`, and repays it. */
export function cover(
  deployments: Deployments,
  { asset, amount, maxCost, account }: { asset: string; amount: bigint; maxCost: bigint; account: Address },
) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'cover',
    args: [toBytes32(asset), amount, maxCost, account],
  } as const
}

/** Buys back all `account`'s short of `asset` owes once it falls short. Anyone may. */
export function liquidate(deployments: Deployments, { account, asset }: { account: Address; asset: string }) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'liquidate',
    args: [account, toBytes32(asset)],
  } as const
}

/** Records `asset`'s band centre for the restriction after a 10% fall. Anyone may. */
export function mark(deployments: Deployments, asset: string) {
  return {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'mark',
    args: [toBytes32(asset)],
  } as const
}

/** `account`'s short of `asset`: its USDG, negative in deficit, what it owes in the token, its shares and its book. */
export async function position(client: Client, deployments: Deployments, account: Address, asset: string, at: At = {}) {
  const [usdgHeld, debt, shares, epoch] = await readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'position',
    args: [account, toBytes32(asset)],
    ...at,
  })
  return { usdgHeld, debt, shares, epoch }
}

/** `account`'s short of `asset` at the band's high edge: equity and requirement in USD with 18 decimals. */
export async function health(client: Client, deployments: Deployments, account: Address, asset: string, at: At = {}) {
  const [equity, requirement, missing, regime] = await readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'health',
    args: [account, toBytes32(asset)],
    ...at,
  })
  return { equity, requirement, missing, regime }
}

/** The epoch of `asset`'s current book. */
export function epoch(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  return readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'epoch',
    args: [toBytes32(asset)],
    ...at,
  })
}

/** The borrow shares, USDG and cost index of `asset`'s shorts in book `epoch`. */
export async function book(client: Client, deployments: Deployments, asset: string, epoch: bigint, at: At = {}) {
  const [shares, usdgHeld, costIndex] = await readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'book',
    args: [toBytes32(asset), epoch],
    ...at,
  })
  return { shares, usdgHeld, costIndex }
}

/** Whether a short of `asset` sells at no less than the band's centre now, and the close it measures from. */
export async function restriction(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [restricted, close] = await readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'restriction',
    args: [toBytes32(asset)],
    ...at,
  })
  return { restricted, close }
}

/** The fee tier of `asset`'s pool with USDG; zero where it may not be shorted. */
export function fee(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  return readContract(client, {
    address: shorts(deployments),
    abi: shortPositionsAbi,
    functionName: 'fee',
    args: [toBytes32(asset)],
    ...at,
  })
}

/**
 * What selling `amount` of `asset` pays through its pool now, from Uniswap's QuoterV2, and that less `slippageBps`,
 * from 0 to 10000.
 */
export async function quoteSale(
  client: Client,
  deployments: Deployments,
  { asset, amount, slippageBps }: { asset: string; amount: bigint; slippageBps: bigint },
) {
  checkSlippage(slippageBps)
  const { result } = await simulateContract(client, {
    address: entry(deployments.uniswapV3, 'QuoterV2', '.uniswapV3'),
    abi: quoterV2Abi,
    functionName: 'quoteExactInputSingle',
    args: [
      {
        tokenIn: entry(deployments.tokens, asset, '.tokens'),
        tokenOut: entry(deployments.tokens, 'USDG', '.tokens'),
        amountIn: amount,
        fee: await fee(client, deployments, asset),
        sqrtPriceLimitX96: 0n,
      },
    ],
  })
  const proceeds = result[0]
  return { proceeds, minProceeds: (proceeds * (BPS - slippageBps)) / BPS }
}

/**
 * What buying back `amount` of `asset` costs through its pool now, from Uniswap's QuoterV2, and that plus
 * `slippageBps`, from 0 to 10000, rounded up.
 */
export async function quoteCover(
  client: Client,
  deployments: Deployments,
  { asset, amount, slippageBps }: { asset: string; amount: bigint; slippageBps: bigint },
) {
  checkSlippage(slippageBps)
  const { result } = await simulateContract(client, {
    address: entry(deployments.uniswapV3, 'QuoterV2', '.uniswapV3'),
    abi: quoterV2Abi,
    functionName: 'quoteExactOutputSingle',
    args: [
      {
        tokenIn: entry(deployments.tokens, 'USDG', '.tokens'),
        tokenOut: entry(deployments.tokens, asset, '.tokens'),
        amount,
        fee: await fee(client, deployments, asset),
        sqrtPriceLimitX96: 0n,
      },
    ],
  })
  const cost = result[0]
  return { cost, maxCost: (cost * (BPS + slippageBps) + BPS - 1n) / BPS }
}

function checkSlippage(slippageBps: bigint) {
  if (slippageBps < 0n || slippageBps > BPS) throw new RangeError(`slippageBps ${slippageBps} is outside 0 to 10000.`)
}

function shorts(deployments: Deployments) {
  return entry(deployments.tapehouse, 'ShortPositions', '.tapehouse')
}
