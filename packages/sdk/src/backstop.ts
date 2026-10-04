// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, type Hex } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import { stocks } from './accounts.js'
import type { At } from './band.js'
import { CROSS, type Deployments, entry, toBytes32 } from './deployments.js'
import { gapBackstopAbi, marginAccountsAbi } from './generated.js'
import { holder } from './vault.js'

/** Deposits `assets` of the caller's USDG with the backstop, for shares to `receiver`. */
export function deposit(deployments: Deployments, { assets, receiver }: { assets: bigint; receiver: Address }) {
  return { address: vault(deployments), abi: gapBackstopAbi, functionName: 'deposit', args: [assets, receiver] } as const
}

/**
 * Starts the caller's cooldown for all its shares: they may be redeemed from a week later, for six days, while
 * redemptions are open. A new cooldown replaces the last; shares sent away leave it.
 */
export function startCooldown(deployments: Deployments) {
  return { address: vault(deployments), abi: gapBackstopAbi, functionName: 'startCooldown' } as const
}

/** Sends `assets` of `owner`'s USDG to `receiver`, within its cooldown's window. */
export function withdraw(
  deployments: Deployments,
  { assets, receiver, owner }: { assets: bigint; receiver: Address; owner: Address },
) {
  return {
    address: vault(deployments),
    abi: gapBackstopAbi,
    functionName: 'withdraw',
    args: [assets, receiver, owner],
  } as const
}

/** Redeems `owner`'s `shares` for USDG to `receiver`, within its cooldown's window. */
export function redeem(
  deployments: Deployments,
  { shares, receiver, owner }: { shares: bigint; receiver: Address; owner: Address },
) {
  return {
    address: vault(deployments),
    abi: gapBackstopAbi,
    functionName: 'redeem',
    args: [shares, receiver, owner],
  } as const
}

/** Sends the caller the Stock Token `token` it gained from what the backstop bought at the reopening auction. */
export function claimGains(deployments: Deployments, token: Address) {
  return { address: vault(deployments), abi: gapBackstopAbi, functionName: 'claimGains', args: [token] } as const
}

/** Claims the premium the accounts set aside for the backstop, for its shares at once. Anyone may. */
export function claim(deployments: Deployments) {
  return { address: vault(deployments), abi: gapBackstopAbi, functionName: 'claim' } as const
}

/**
 * The backstop at one block: the USDG it holds and its shares; the close, in milliseconds, of the closure it last
 * counted a cover against; for the cross positions and each asset's isolated positions, the exposure limit of the
 * current closure, the limits before and from `fromMs`, what is left of it, at most the USDG held, and what it covered
 * in the closure from `closureMs`; the accounts' weekend premium, in basis points a year, the reserve's part of each
 * premium paid, in basis points, and the backstop's part waiting to be claimed; its cooldown, withdrawal window and
 * settlement, in seconds; and the owner's seed, as the USDG its shares are worth.
 */
export async function state(client: Client, deployments: Deployments, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: vault(deployments), abi: gapBackstopAbi, blockNumber } as const
  const premium = {
    address: entry(deployments.tapehouse, 'MarginAccounts', '.tapehouse'),
    abi: marginAccountsAbi,
    blockNumber,
  } as const
  const [assets, held, shares, closureMs, cooldown, window, settlement, owner, rate, reserveShare, waiting] =
    await Promise.all([
      stocks(client, deployments, { blockNumber }),
      readContract(client, { ...read, functionName: 'totalAssets' }),
      readContract(client, { ...read, functionName: 'totalSupply' }),
      readContract(client, { ...read, functionName: 'closureMs' }),
      readContract(client, { ...read, functionName: 'COOLDOWN' }),
      readContract(client, { ...read, functionName: 'WITHDRAWAL_WINDOW' }),
      readContract(client, { ...read, functionName: 'SETTLEMENT' }),
      readContract(client, { ...read, functionName: 'owner' }),
      readContract(client, { ...premium, functionName: 'premiumRate' }),
      readContract(client, { ...premium, functionName: 'reserveShare' }),
      readContract(client, { ...premium, functionName: 'backstopPremium' }),
    ])
  const positions: Hex[] = [CROSS, ...assets.map(({ asset }) => toBytes32(asset))]
  const [exposures, seedShares] = await Promise.all([
    Promise.all(
      positions.map(async (position) => {
        const [limit, [current, next, fromMs], left, covered] = await Promise.all([
          readContract(client, { ...read, functionName: 'exposureLimit', args: [position] }),
          readContract(client, { ...read, functionName: 'exposureLimits', args: [position] }),
          readContract(client, { ...read, functionName: 'exposureLeft', args: [position] }),
          readContract(client, { ...read, functionName: 'covered', args: [position, closureMs] }),
        ])
        return { position, limit, current, next, fromMs, left, covered }
      }),
    ),
    readContract(client, { ...read, functionName: 'balanceOf', args: [owner] }),
  ])
  const seed = await readContract(client, { ...read, functionName: 'convertToAssets', args: [seedShares] })
  return {
    held,
    shares,
    closureMs,
    exposures,
    premium: { rate, reserveShare, waiting },
    terms: { cooldown, window, settlement },
    seed: { owner, assets: seed },
  }
}

/**
 * What `owner` holds in the backstop, at one block: its shares, with 12 decimals, the USDG they are worth, and what
 * it may withdraw, redeem and deposit now; its cooldown, with when its window opens and closes, in seconds, once it
 * started one; and the units of each of `tokens` it gained and has not claimed.
 */
export async function depositor(
  client: Client,
  deployments: Deployments,
  owner: Address,
  tokens: readonly Address[],
  at: At = {},
) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: vault(deployments), abi: gapBackstopAbi, blockNumber } as const
  const [held, [cooling, startedAt], terms, gains] = await Promise.all([
    holder(client, read.address, owner, blockNumber),
    readContract(client, { ...read, functionName: 'cooldowns', args: [owner] }),
    Promise.all([
      readContract(client, { ...read, functionName: 'COOLDOWN' }),
      readContract(client, { ...read, functionName: 'WITHDRAWAL_WINDOW' }),
    ]),
    Promise.all(tokens.map((token) => readContract(client, { ...read, functionName: 'gains', args: [owner, token] }))),
  ])
  const [cooldown, window] = terms
  return {
    ...held,
    cooldown:
      startedAt === 0n
        ? undefined
        : { shares: cooling, from: startedAt + cooldown, until: startedAt + cooldown + window },
    gains: Object.fromEntries(tokens.map((token, i) => [token, gains[i] ?? 0n])) as Record<Address, bigint>,
  }
}

function vault(deployments: Deployments) {
  return entry(deployments.tapehouse, 'GapBackstop', '.tapehouse')
}
