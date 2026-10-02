// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Client, Hex } from 'viem'
import { readContract } from 'viem/actions'
import { type Deployments, entry, type PriceFeed, type PriceKind, toBytes32 } from './deployments.js'
import { aggregatorAbi, bandAbi, bandFeedAbi } from './generated.js'

/** The block a read is taken at: the latest by default. */
export type At = { blockNumber?: bigint }

/**
 * Where signed RedStone data packages come from: the integrator's own gateway client, cache or relay. It returns the
 * payload `writePrices` verifies, the packages for `feedIds` serialised as RedStone's EVM connector appends them. The
 * SDK holds no API key and calls no gateway.
 */
export interface PackageSource {
  payload(feedIds: readonly string[]): Promise<Hex>
}

/** The band of `asset`: its state (0 halted, 1 degraded, 2 closed, 3 open), live legs, centre, half-width and edges. */
export async function quote(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [state, live, mid, halfBps, low, high] = await readContract(client, {
    address: band(deployments),
    abi: bandAbi,
    functionName: 'quote',
    args: [toBytes32(asset)],
    ...at,
  })
  return { state, live, mid, halfBps, low, high }
}

/** Chainlink's 24/5 session (0 not known, 1 closed, 2 open), NYSE's state and next state, and their boundaries. */
export async function session(client: Client, deployments: Deployments, at: At = {}) {
  const [state, nyse, nyseNext, changeMs, boundaryMs] = await readContract(client, {
    address: band(deployments),
    abi: bandAbi,
    functionName: 'session',
    ...at,
  })
  return { state, nyse, nyseNext, changeMs, boundaryMs }
}

/** The trading halt of `asset`: whether a signed halt holds, until when, when it was issued, and the issuer's pause. */
export async function halt(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [signedHalt, until, issuedAt, oraclePaused] = await readContract(client, {
    address: band(deployments),
    abi: bandAbi,
    functionName: 'halt',
    args: [toBytes32(asset)],
    ...at,
  })
  return { signedHalt, until, issuedAt, oraclePaused }
}

/** The RedStone value the band stores for `feedId`, its package's timestamp in milliseconds and when it was written. */
export async function price(client: Client, deployments: Deployments, feedId: string, at: At = {}) {
  const [value, packageTimestampMs, writtenAt] = await readContract(client, {
    address: band(deployments),
    abi: bandAbi,
    functionName: 'price',
    args: [toBytes32(feedId)],
    ...at,
  })
  return { value, packageTimestampMs, writtenAt }
}

/** The band's `writePrices` of `feedIds`, with the payload `source` signs for them. */
export async function writePrices(deployments: Deployments, source: PackageSource, feedIds: readonly string[]) {
  const payload = await source.payload(feedIds)
  return {
    address: band(deployments),
    abi: bandAbi,
    functionName: 'writePrices',
    args: [feedIds.map(toBytes32), payload],
  } as const
}

/**
 * The low side of `asset`'s band from its `BandFeed`, as Chainlink's `latestRoundData`: its round and `updatedAt` are
 * the block time.
 */
export async function latestRound(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [roundId, answer, startedAt, updatedAt, answeredInRound] = await readContract(client, {
    address: entry(deployments.bandFeeds, asset, '.bandFeeds'),
    abi: bandFeedAbi,
    functionName: 'latestRoundData',
    ...at,
  })
  return { roundId, answer, startedAt, updatedAt, answeredInRound }
}

/** The whole band of `asset` from its `BandFeed`, the token's market beside it, and both halt sources. */
export function latestBand(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  return readContract(client, {
    address: entry(deployments.bandFeeds, asset, '.bandFeeds'),
    abi: bandFeedAbi,
    functionName: 'latestBand',
    ...at,
  })
}

/** The band of `asset` its `BandFeed` sealed before the reopen at `reopenMs`; all zero where none was sealed. */
export async function sealed(client: Client, deployments: Deployments, asset: string, reopenMs: bigint, at: At = {}) {
  const [state, live, mid, halfBps, low, high, sealedAt] = await readContract(client, {
    address: entry(deployments.bandFeeds, asset, '.bandFeeds'),
    abi: bandFeedAbi,
    functionName: 'seals',
    args: [reopenMs],
    ...at,
  })
  return { state, live, mid, halfBps, low, high, sealedAt }
}

/** `asset`'s `BandFeed.seal()`, which anyone may call in the ten minutes before the session reopens. */
export function seal(deployments: Deployments, asset: string) {
  return {
    address: entry(deployments.bandFeeds, asset, '.bandFeeds'),
    abi: bandFeedAbi,
    functionName: 'seal',
  } as const
}

/** The latest round of a Chainlink feed of `.chainlink`, carrying what it prices. */
export async function chainlinkRound<kind extends PriceKind>(client: Client, feed: PriceFeed<kind>, at: At = {}) {
  const [roundId, answer, startedAt, updatedAt, answeredInRound] = await readContract(client, {
    address: feed.address,
    abi: aggregatorAbi,
    functionName: 'latestRoundData',
    ...at,
  })
  return { kind: feed.kind, roundId, answer, startedAt, updatedAt, answeredInRound }
}

function band(deployments: Deployments) {
  return entry(deployments.tapehouse, 'Band', '.tapehouse')
}
