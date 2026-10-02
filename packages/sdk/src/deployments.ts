// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, getAddress, type Hex, isAddress, stringToHex, zeroHash } from 'viem'

/** What a Chainlink feed of `.chainlink` prices: the Stock Token, its multiplier included, or the share. */
export type PriceKind = 'token' | 'share'

/** A Chainlink feed of `.chainlink`, with what it prices. */
export type PriceFeed<kind extends PriceKind = PriceKind> = { address: Address; kind: kind }

/** A feed that prices the Stock Token, as Robinhood Chain's do. */
export type TokenPriceFeed = PriceFeed<'token'>

/** A feed that prices the share, as Arbitrum One's do. */
export type SharePriceFeed = PriceFeed<'share'>

/** A chain's address registry, `deployments/<chainId>.json`, every address checksummed. */
export type Deployments = {
  chainId: number
  tokens: Record<string, Address>
  chainlink: Record<string, PriceFeed>
  bandFeeds: Record<string, Address>
  tapehouse: Record<string, Address>
  stockLending: Record<string, Address>
  baskets: Record<string, Address>
  uniswapV3: Record<string, Address>
  morpho: Record<string, Address>
  morphoMarkets: Record<string, Hex>
  morphoOracles: Record<string, Address>
}

/** The chains whose Chainlink feeds price the share rather than the Stock Token. */
export const SHARE_PRICE_CHAINS: readonly number[] = [42161]

/** The margin accounts' cross position. */
export const CROSS: Hex = zeroHash

/** An asset's symbol, or a RedStone feed ID, as the contracts take it: its ASCII bytes, right-padded to 32. */
export function toBytes32(name: string): Hex {
  return stringToHex(name, { size: 32 })
}

/**
 * Parses a chain's registry as read from `deployments/<chainId>.json`. An address in mixed case must carry its
 * checksum; every address is returned checksummed by `getAddress`.
 */
export function parseDeployments(json: unknown): Deployments {
  const registry = record(json, 'the registry')
  const chainId = registry.chainId
  if (typeof chainId !== 'number' || !Number.isSafeInteger(chainId) || chainId <= 0)
    throw new Error('The registry has no chainId.')
  const kind: PriceKind = SHARE_PRICE_CHAINS.includes(chainId) ? 'share' : 'token'
  const { StockLending, Baskets, ...tapehouse } = record(registry.tapehouse ?? {}, '.tapehouse')
  const { Markets, ...morpho } = record(registry.morpho ?? {}, '.morpho')
  return {
    chainId,
    tokens: addresses(registry.tokens, '.tokens'),
    chainlink: Object.fromEntries(
      Object.entries(addresses(registry.chainlink, '.chainlink')).map(([name, address]) => [name, { address, kind }]),
    ),
    bandFeeds: addresses(registry.bandFeeds, '.bandFeeds'),
    tapehouse: addresses(tapehouse, '.tapehouse'),
    stockLending: addresses(StockLending, '.tapehouse.StockLending'),
    baskets: addresses(Baskets, '.tapehouse.Baskets'),
    uniswapV3: addresses(registry.uniswapV3, '.uniswapV3'),
    morpho: addresses(morpho, '.morpho'),
    morphoMarkets: ids(Markets, '.morpho.Markets'),
    morphoOracles: addresses(registry.morphoOracles, '.morphoOracles'),
  }
}

/** The Chainlink feed `name` of a chain whose feeds price the Stock Token. */
export function tokenPriceFeed(deployments: Deployments, name: string): TokenPriceFeed {
  const feed = entry(deployments.chainlink, name, '.chainlink')
  if (feed.kind !== 'token') throw new Error(`.chainlink.${name} prices the share on chain ${deployments.chainId}.`)
  return { address: feed.address, kind: feed.kind }
}

/** The Chainlink feed `name` of a chain whose feeds price the share. */
export function sharePriceFeed(deployments: Deployments, name: string): SharePriceFeed {
  const feed = entry(deployments.chainlink, name, '.chainlink')
  if (feed.kind !== 'share')
    throw new Error(`.chainlink.${name} prices the Stock Token on chain ${deployments.chainId}.`)
  return { address: feed.address, kind: feed.kind }
}

/** The address of `name` in the registry group `group`. */
export function entry<value>(group: Record<string, value>, name: string, path: string): value {
  const value = Object.hasOwn(group, name) ? group[name] : undefined
  if (value === undefined) throw new Error(`The registry has no ${path}.${name}.`)
  return value
}

function record(value: unknown, path: string): Record<string, unknown> {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) throw new Error(`${path} is not an object.`)
  return value as Record<string, unknown>
}

function addresses(value: unknown, path: string): Record<string, Address> {
  return Object.fromEntries(
    Object.entries(record(value ?? {}, path)).map(([name, address]) => {
      if (typeof address !== 'string' || !isAddress(address)) throw new Error(`${path}.${name} is not an address.`)
      return [name, getAddress(address)]
    }),
  )
}

function ids(value: unknown, path: string): Record<string, Hex> {
  return Object.fromEntries(
    Object.entries(record(value ?? {}, path)).map(([name, id]) => {
      if (typeof id !== 'string' || !/^0x[0-9a-fA-F]{64}$/.test(id))
        throw new Error(`${path}.${name} is not a 32-byte id.`)
      return [name, id.toLowerCase() as Hex]
    }),
  )
}
