// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, type Hex, zeroAddress } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry, toBytes32 } from './deployments.js'
import { gapCoverAbi } from './generated.js'
import { holder } from './vault.js'

const BPS = 10_000n

/**
 * Where a series stands: still to settle, settled at its reopening price, or void and refunding its premiums, because
 * it could not settle on a price the band vouched for or no one settled it within a week.
 */
export type SeriesStatus = 'open' | 'settled' | 'void'

/** The layer a cover pays: the fall beyond `deductibleBps`, up to `limitBps`, on `notional` USDG. */
export type Layer = { notional: bigint; deductibleBps: bigint; limitBps: bigint }

/** A series' settlement: its reference and reopening prices, in USD with 8 decimals, once it has settled. */
export type Settlement = { referencePrice: bigint; price: bigint }

/** Buys cover on `asset` for `holder` over the coming closure, paying at most `maxPremium`, from `quote`. */
export function buy(
  deployments: Deployments,
  { asset, maxPremium, holder, ...layer }: Layer & { asset: string; maxPremium: bigint; holder: Address },
) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'buy',
    args: [toBytes32(asset), layer.notional, layer.deductibleBps, layer.limitBps, maxPremium, holder],
  } as const
}

/** Deposits `assets` of the caller's USDG with the writers, for shares to `receiver`. */
export function deposit(deployments: Deployments, { assets, receiver }: { assets: bigint; receiver: Address }) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'deposit',
    args: [assets, receiver],
  } as const
}

/** Sends `assets` of `owner`'s USDG with the writers to `receiver`, once no cover is outstanding. */
export function withdraw(
  deployments: Deployments,
  { assets, receiver, owner }: { assets: bigint; receiver: Address; owner: Address },
) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'withdraw',
    args: [assets, receiver, owner],
  } as const
}

/** Redeems `owner`'s `shares` of the writers' USDG to `receiver`, once no cover is outstanding. */
export function redeem(
  deployments: Deployments,
  { shares, receiver, owner }: { shares: bigint; receiver: Address; owner: Address },
) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'redeem',
    args: [shares, receiver, owner],
  } as const
}

/** Records the session's next close, or the reopen while it is closed. Anyone may. */
export function record(deployments: Deployments) {
  return { address: gapCover(deployments), abi: gapCoverAbi, functionName: 'record' } as const
}

/**
 * Keeps the realised move of the week before the close on sale for `asset`'s series, so that its first buyer does not
 * pay for reading it from the feed: during the sales, once the week's last day, a day before the close, has passed.
 * Anyone may.
 */
export function measure(deployments: Deployments, asset: string) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'measure',
    args: [toBytes32(asset)],
  } as const
}

/** Records `asset`'s band in this minute's slot of the window after the closure from `closesMs`. Anyone may. */
export function observe(deployments: Deployments, { asset, closesMs }: { asset: string; closesMs: bigint }) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'observe',
    args: [toBytes32(asset), closesMs],
  } as const
}

/**
 * Settles `asset`'s series over the closure from `closesMs`, naming its feed's last rounds started before the sales
 * ended and before the reopen. Anyone may.
 */
export function settle(
  deployments: Deployments,
  {
    asset,
    closesMs,
    referenceRound,
    lastRound,
  }: { asset: string; closesMs: bigint; referenceRound: bigint; lastRound: bigint },
) {
  return {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'settle',
    args: [toBytes32(asset), closesMs, referenceRound, lastRound],
  } as const
}

/** Credits cover `id`'s payout or refund to its holder once its series has settled or is void. Anyone may. */
export function release(deployments: Deployments, id: bigint) {
  return { address: gapCover(deployments), abi: gapCoverAbi, functionName: 'release', args: [id] } as const
}

/** Sends `receiver` the USDG credited to the caller. */
export function claim(deployments: Deployments, receiver: Address) {
  return { address: gapCover(deployments), abi: gapCoverAbi, functionName: 'claim', args: [receiver] } as const
}

/**
 * The premium of cover on `asset` at the weekend gap of `pricingGap` and the band's centre now, and the writers'
 * USDG it reserves.
 */
export async function quote(client: Client, deployments: Deployments, asset: string, layer: Layer, at: At = {}) {
  const premium = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'quote',
    args: [toBytes32(asset), layer.notional, layer.deductibleBps, layer.limitBps],
    ...at,
  })
  return { premium, reserve: (layer.notional * (layer.limitBps - layer.deductibleBps) + BPS - 1n) / BPS }
}

/**
 * The smallest deductible of cover on `asset` now, in basis points: the start of the fitted tail at the weekend gap of
 * `pricingGap`, plus how far the feed's last answer stands above the band's centre.
 */
export function minDeductible(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  return readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'minDeductible',
    args: [toBytes32(asset)],
    ...at,
  })
}

/**
 * The weekend gap the cover prices `asset` at now, in millionths, and the realised move of the trading week before the
 * close on sale that it follows: during the sales twice that move, at least 55% of the engine's gap and at most eight
 * times it; outside them the engine's gap and a move of zero. A move of 2^64 - 1 is a week the feed cannot show, which
 * prices at the most.
 */
export async function pricingGap(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [gap, weekMove] = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'pricingGap',
    args: [toBytes32(asset)],
    ...at,
  })
  return { gap, weekMove }
}

/** The writers' USDG not reserved for a cover: the most a new cover may reserve. */
export function capacity(client: Client, deployments: Deployments, at: At = {}) {
  return readContract(client, { address: gapCover(deployments), abi: gapCoverAbi, functionName: 'capacity', ...at })
}

/** The close that keys the series of the coming or the last recorded closure, in milliseconds. */
export function lastCloseMs(client: Client, deployments: Deployments, at: At = {}) {
  return readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'lastCloseMs',
    ...at,
  })
}

/**
 * The closure on sale now: the close that keys its series and when its sales end, in milliseconds; undefined while no
 * cover is sold.
 */
export async function sales(client: Client, deployments: Deployments, at: At = {}) {
  const [closesMs, endsMs] = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'sales',
    ...at,
  })
  return endsMs === 0n ? undefined : { closesMs, endsMs }
}

/**
 * `asset`'s series over the closure from `closesMs`: the notional covered, where it stands, its settlement once it has
 * settled, whether it settled on the band's median centre rather than the first round, and when its reopen is.
 */
export async function series(
  client: Client,
  deployments: Deployments,
  asset: string,
  closesMs: bigint,
  at: At = {},
) {
  const [notional, referencePrice, price, , status, flagged] = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'series',
    args: [toBytes32(asset), closesMs],
    ...at,
  })
  const reopenMs = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'reopenOf',
    args: [closesMs],
    ...at,
  })
  const statuses: SeriesStatus[] = ['open', 'settled', 'void']
  const state = statuses[status]
  if (state === undefined) throw new Error(`Unknown series status ${status}.`)
  return {
    notional,
    status: state,
    settlement: state === 'settled' ? { referencePrice, price } : undefined,
    flagged,
    reopenMs,
  }
}

/** Cover `id`: its holder, closure, asset's symbol and layer, and its premium; undefined once released. */
export async function cover(client: Client, deployments: Deployments, id: bigint, at: At = {}) {
  const [holder, closesMs, deductibleBps, limitBps, symbol, notional, premium] = await readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'covers',
    args: [id],
    ...at,
  })
  if (holder === zeroAddress) return undefined
  return {
    holder,
    closesMs,
    symbol: symbol as Hex,
    notional,
    deductibleBps: BigInt(deductibleBps),
    limitBps: BigInt(limitBps),
    premium,
  }
}

/** The USDG credited to `holder` and not yet claimed. */
export function payouts(client: Client, deployments: Deployments, holder: Address, at: At = {}) {
  return readContract(client, {
    address: gapCover(deployments),
    abi: gapCoverAbi,
    functionName: 'payouts',
    args: [holder],
    ...at,
  })
}

/**
 * What a cover pays once its series settled: its notional times the fall beyond its deductible, up to its limit, before
 * any shortfall a wipe leaves.
 */
export function payout(layer: Layer, { referencePrice, price }: Settlement): bigint {
  if (price >= referencePrice) return 0n
  const limit = referencePrice * layer.limitBps
  const fall = (referencePrice - price) * BPS < limit ? (referencePrice - price) * BPS : limit
  const deductible = referencePrice * layer.deductibleBps
  if (fall <= deductible) return 0n
  return (layer.notional * (fall - deductible)) / (referencePrice * BPS)
}

/**
 * The writers at one block: the USDG they hold and their shares; what the outstanding covers reserve and what is left
 * for new ones; the premiums of the covers not yet released, which join their USDG only at release, less what each
 * pays; the payouts and refunds owed to holders; how many covers are not released, any of which stops redemptions; and
 * the closure on sale, undefined while none is.
 */
export async function writers(client: Client, deployments: Deployments, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: gapCover(deployments), abi: gapCoverAbi, blockNumber } as const
  const [held, shares, reserved, capacity, premiums, owed, outstanding, [closesMs, endsMs]] = await Promise.all([
    readContract(client, { ...read, functionName: 'totalAssets' }),
    readContract(client, { ...read, functionName: 'totalSupply' }),
    readContract(client, { ...read, functionName: 'reserved' }),
    readContract(client, { ...read, functionName: 'capacity' }),
    readContract(client, { ...read, functionName: 'premiums' }),
    readContract(client, { ...read, functionName: 'owed' }),
    readContract(client, { ...read, functionName: 'outstanding' }),
    readContract(client, { ...read, functionName: 'sales' }),
  ])
  return {
    held,
    shares,
    reserved,
    capacity,
    premiums,
    owed,
    outstanding,
    sales: endsMs === 0n ? undefined : { closesMs, endsMs },
  }
}

/**
 * What `owner` holds with the writers, at one block: its shares, with 12 decimals, the USDG they are worth, and what
 * it may withdraw, redeem and deposit now. While covers are out, deposits stay open only for the coming closure's, while
 * the session is open and before its sales end; redemptions open once every cover is released.
 */
export async function writer(client: Client, deployments: Deployments, owner: Address, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  return holder(client, gapCover(deployments), owner, blockNumber)
}

function gapCover(deployments: Deployments) {
  return entry(deployments.tapehouse, 'GapCover', '.tapehouse')
}
