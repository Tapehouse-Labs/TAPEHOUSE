// SPDX-License-Identifier: MIT OR Apache-2.0
import { decodeFunctionData, encodeFunctionData, padHex } from 'viem'
import { describe, expect, test } from 'vitest'
import { auction, parseDeployments, reopeningAuctionAbi, toBytes32 } from '../src/index.js'

const AUCTION = '0x58F39a807ea6CFbeBCBd71995b1dA2E0608BF642'
const deployments = parseDeployments({ chainId: 46630, tapehouse: { ReopeningAuction: AUCTION } })
const bidder = '0x70997970C51812dc3A010C7d01b50e0d17dc79C8'
const bid = (quantity: bigint, price: bigint) => ({ quantity, price })

describe('auction', () => {
  test('a bid is committed under the hash the Go SDK computes', () => {
    expect(
      auction.commitment({
        bidder,
        asset: 'SPY',
        openMs: 1_791_000_000_000n,
        quantity: 5n * 10n ** 17n,
        price: 60_000_000_000n,
        salt: padHex('0x010203', { dir: 'right', size: 32 }),
      }),
    ).toBe('0xe7e7807e44240c1f4e3f32b45152769b116802df7bb762bc47a7865c748e5bba')
  })

  test.each([
    ['no bid clears at the floor', [], 10n, 90n],
    ['the highest price the bids at or above it take the supply at', [bid(4n, 120n), bid(5n, 110n), bid(3n, 100n), bid(9n, 95n)], 10n, 100n],
    ['the bids above it may fill exactly', [bid(4n, 120n), bid(6n, 110n), bid(3n, 100n)], 10n, 110n],
    ['several bids at one price', [bid(2n, 120n), bid(4n, 105n), bid(5n, 105n), bid(1n, 99n)], 10n, 105n],
    ['undersubscribed clears at the lowest bid', [bid(2n, 120n), bid(3n, 101n)], 10n, 101n],
    ['bids and no supply cannot clear', [bid(2n, 120n)], 0n, undefined],
  ])('the clearing price, as the Go SDK computes it: %s', (_name, bids, supply, want) => {
    expect(auction.clearingPrice(bids, supply, 90n)).toBe(want)
  })

  test('each call names the auction, the asset as bytes32 and its arguments', () => {
    const salt = padHex('0x07', { size: 32 })
    const calls = [
      auction.commit(deployments, { asset: 'SPY', commitment: salt, deposit: 100_000_000n }),
      auction.reveal(deployments, { asset: 'SPY', openMs: 1n, quantity: 2n, price: 3n, salt }),
      auction.claim(deployments, { asset: 'SPY', openMs: 1n, index: 0n }),
      auction.forfeit(deployments, { asset: 'SPY', openMs: 1n, commitment: salt }),
      auction.clear(deployments, { asset: 'SPY', openMs: 1n, price: 3n }),
    ]
    for (const call of calls) {
      expect(call.address).toBe(AUCTION)
      expect(call.args[0]).toBe(toBytes32('SPY'))
    }
    const reveal = auction.reveal(deployments, { asset: 'SPY', openMs: 1n, quantity: 2n, price: 3n, salt })
    expect(decodeFunctionData({ abi: reopeningAuctionAbi, data: encodeFunctionData(reveal) }).functionName).toBe('reveal')
  })
})
