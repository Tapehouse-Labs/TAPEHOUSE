// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, hexToString } from 'viem'
import { readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry } from './deployments.js'
import { basketAbi } from './generated.js'

/**
 * Mints `shares` of the registry's basket `basket` to `receiver`, for at most `maxAssets` of each of its Stock Tokens,
 * in its order: `previewMint` gives what it takes, its part of every token rounded up.
 */
export function mint(
  deployments: Deployments,
  {
    basket,
    shares,
    receiver,
    maxAssets,
  }: { basket: string; shares: bigint; receiver: Address; maxAssets: readonly bigint[] },
) {
  return {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'mint',
    args: [shares, receiver, maxAssets],
  } as const
}

/** Redeems `shares` of `owner` in the registry's basket `basket` for its Stock Tokens, to `receiver`. */
export function redeem(
  deployments: Deployments,
  { basket, shares, receiver, owner }: { basket: string; shares: bigint; receiver: Address; owner: Address },
) {
  return {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'redeem',
    args: [shares, receiver, owner],
  } as const
}

/** The assets the basket holds and their Stock Tokens, in the order of every list of amounts. */
export async function components(client: Client, deployments: Deployments, basket: string, at: At = {}) {
  const [symbols, tokens] = await readContract(client, {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'components',
    ...at,
  })
  return { assets: symbols.map((symbol) => hexToString(symbol, { size: 32 })), tokens }
}

/** What minting `shares` takes of each Stock Token, in raw units, rounded up. */
export function previewMint(client: Client, deployments: Deployments, basket: string, shares: bigint, at: At = {}) {
  return readContract(client, {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'previewMint',
    args: [shares],
    ...at,
  })
}

/** What redeeming `shares` gives of each Stock Token, in raw units, rounded down. */
export function previewRedeem(client: Client, deployments: Deployments, basket: string, shares: bigint, at: At = {}) {
  return readContract(client, {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'previewRedeem',
    args: [shares],
    ...at,
  })
}

/** The target in effect: raw units of each Stock Token per share. */
export function target(client: Client, deployments: Deployments, basket: string, at: At = {}) {
  return readContract(client, {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'target',
    ...at,
  })
}

/** The target proposed and not yet in effect, and when it takes effect, in seconds; empty and zero if none. */
export async function pendingTarget(client: Client, deployments: Deployments, basket: string, at: At = {}) {
  const [units, effectiveAt] = await readContract(client, {
    address: basketAddress(deployments, basket),
    abi: basketAbi,
    functionName: 'pendingTarget',
    ...at,
  })
  return { units, effectiveAt }
}

function basketAddress(deployments: Deployments, basket: string) {
  return entry(deployments.baskets, basket, '.tapehouse.Baskets')
}
