// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, encodeAbiParameters, type Hex, keccak256 } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry, toBytes32 } from './deployments.js'
import * as band from './band.js'
import { liquidatorAbi, reopeningAuctionAbi } from './generated.js'

/** A revealed bid of a round, as the auction keeps it. */
export type Bid = { bidder: Address; quantity: bigint; price: bigint; escrow: bigint; claimed: boolean }

/**
 * The hash a bid is committed under: `keccak256(abi.encode(bidder, symbol, openMs, quantity, price, salt))`, which
 * binds the bid to its bidder. `quantity` has 18 decimals and `price` is USD with 8 decimals per token.
 */
export function commitment({
  bidder,
  asset,
  openMs,
  quantity,
  price,
  salt,
}: {
  bidder: Address
  asset: string
  openMs: bigint
  quantity: bigint
  price: bigint
  salt: Hex
}): Hex {
  return keccak256(
    encodeAbiParameters(
      [
        { type: 'address' },
        { type: 'bytes32' },
        { type: 'uint64' },
        { type: 'uint256' },
        { type: 'uint256' },
        { type: 'bytes32' },
      ],
      [bidder, toBytes32(asset), openMs, quantity, price, salt],
    ),
  )
}

/**
 * A round's only clearing price: the highest bid price at which the bids at or above it take the whole supply; when
 * all of them together take less, the lowest bid price; with no bid, the floor. Undefined when no price clears it, as
 * for a round with bids and no supply.
 */
export function clearingPrice(bids: readonly Pick<Bid, 'quantity' | 'price'>[], supply: bigint, floor: bigint) {
  if (bids.length === 0) return floor
  if (supply === 0n) return undefined
  let lowest: bigint | undefined
  let best: bigint | undefined
  for (const candidate of bids) {
    if (lowest === undefined || candidate.price < lowest) lowest = candidate.price
    let above = 0n
    let at = 0n
    for (const bid of bids) {
      if (bid.price > candidate.price) above += bid.quantity
      else if (bid.price === candidate.price) at += bid.quantity
    }
    if (above < supply && above + at >= supply && (best === undefined || candidate.price > best)) best = candidate.price
  }
  return best ?? lowest
}

/**
 * Commits the caller to a hidden bid in `asset`'s round with `deposit` of its USDG, at least the bond and at least
 * what the bid will escrow, so the deposit may hide its size. Only while bids are committed.
 */
export function commit(deployments: Deployments, { asset, commitment, deposit }: { asset: string; commitment: Hex; deposit: bigint }) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'commit',
    args: [toBytes32(asset), commitment, deposit],
  } as const
}

/** Reveals the caller's bid in the half hour before the regular open `openMs`. */
export function reveal(
  deployments: Deployments,
  { asset, openMs, quantity, price, salt }: { asset: string; openMs: bigint; quantity: bigint; price: bigint; salt: Hex },
) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'reveal',
    args: [toBytes32(asset), openMs, quantity, price, salt],
  } as const
}

/** Sends the bidder of the bid at `index` its tokens and the rest of its escrow, once the round is cleared or lapsed. */
export function claim(deployments: Deployments, { asset, openMs, index }: { asset: string; openMs: bigint; index: bigint }) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'claim',
    args: [toBytes32(asset), openMs, index],
  } as const
}

/** Sends an unrevealed commitment's forfeit to the accounts' reserve and the rest of its deposit back. Anyone may. */
export function forfeit(deployments: Deployments, { asset, openMs, commitment }: { asset: string; openMs: bigint; commitment: Hex }) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'forfeit',
    args: [toBytes32(asset), openMs, commitment],
  } as const
}

/** Clears `asset`'s round at `price`, which `clearingPrice` computes, in the first hour of regular trading. Anyone may. */
export function clear(deployments: Deployments, { asset, openMs, price }: { asset: string; openMs: bigint; price: bigint }) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'clear',
    args: [toBytes32(asset), openMs, price],
  } as const
}

/** Enrolls `account`'s `position`, short at its bands' low edges, in `asset`'s round while bids are committed. Anyone may. */
export function enroll(
  deployments: Deployments,
  { account, position, asset }: { account: Address; position: Hex; asset: string },
) {
  return {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'enroll',
    args: [account, position, toBytes32(asset)],
  } as const
}

/**
 * The auction's constants: the least deposit and the least bid worth, in USDG with 6 decimals, and the lengths of the
 * reveal and clearing windows, in milliseconds.
 */
export async function constants(client: Client, deployments: Deployments, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: auction(deployments), abi: reopeningAuctionAbi, blockNumber } as const
  const [bond, minBid, revealMs, clearMs] = await Promise.all([
    readContract(client, { ...read, functionName: 'BOND' }),
    readContract(client, { ...read, functionName: 'MIN_BID' }),
    readContract(client, { ...read, functionName: 'REVEAL_MS' }),
    readContract(client, { ...read, functionName: 'CLEAR_MS' }),
  ])
  return { bond, minBid, revealMs, clearMs }
}

/** The regular open, in milliseconds, that bids are committed or revealed for, and whether they are being revealed; zero outside. */
export async function phase(client: Client, deployments: Deployments, at: At = {}) {
  const [openMs, revealing] = await readContract(client, {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'phase',
    blockNumber: at.blockNumber,
  })
  return { openMs, revealing }
}

/** The regular open of the last round opened, in milliseconds. */
export function lastOpenMs(client: Client, deployments: Deployments, at: At = {}) {
  return readContract(client, {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'lastOpenMs',
    blockNumber: at.blockNumber,
  })
}

/** `asset`'s round for the regular open `openMs`, its lots and its revealed bids, at one block. */
export async function round(client: Client, deployments: Deployments, { asset, openMs }: { asset: string; openMs: bigint }, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: auction(deployments), abi: reopeningAuctionAbi, blockNumber, args: [toBytes32(asset), openMs] } as const
  const [state, lots, bids] = await Promise.all([
    readContract(client, { ...read, functionName: 'round' }),
    readContract(client, { ...read, functionName: 'lots' }),
    readContract(client, { ...read, functionName: 'bids' }),
  ])
  return { ...state, lots, bids }
}

/**
 * The floor `asset`'s round for the regular open `openMs` will open with, before anything opens it: the low edge of the
 * band sealed before the reopen, or the live band's where none was sealed, less the liquidator's greatest discount.
 * Zero where the band has no price.
 */
export async function previewFloor(
  client: Client,
  deployments: Deployments,
  { asset, openMs }: { asset: string; openMs: bigint },
  at: At = {},
) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: auction(deployments), abi: reopeningAuctionAbi, blockNumber } as const
  const [lead, liquidator] = await Promise.all([
    readContract(client, { ...read, functionName: 'REOPEN_LEAD_MS' }),
    readContract(client, { ...read, functionName: 'liquidator' }),
  ])
  const [discount, seal, live] = await Promise.all([
    readContract(client, { address: liquidator, abi: liquidatorAbi, functionName: 'MAX_DISCOUNT_BPS', blockNumber }),
    band.sealed(client, deployments, asset, openMs - lead, { blockNumber }),
    band.latestBand(client, deployments, asset, { blockNumber }),
  ])
  const low = seal.state !== 0 ? seal.low : live.low
  return (low * (10_000n - discount)) / 10_000n
}

/** A commitment of a round: its bidder and its deposit, zero once revealed or forfeited. */
export async function commitmentOf(
  client: Client,
  deployments: Deployments,
  { asset, openMs, commitment }: { asset: string; openMs: bigint; commitment: Hex },
  at: At = {},
) {
  const [bidder, deposit] = await readContract(client, {
    address: auction(deployments),
    abi: reopeningAuctionAbi,
    functionName: 'commitments',
    args: [toBytes32(asset), openMs, commitment],
    blockNumber: at.blockNumber,
  })
  return { bidder, deposit }
}

function auction(deployments: Deployments): Address {
  return entry(deployments.tapehouse, 'ReopeningAuction', '.tapehouse')
}
