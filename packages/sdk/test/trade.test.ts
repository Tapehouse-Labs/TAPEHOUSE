// SPDX-License-Identifier: MIT OR Apache-2.0
import { createClient, custom, decodeFunctionData, encodeAbiParameters, encodeEventTopics, encodeFunctionResult, type Hex, zeroAddress } from 'viem'
import { expect, test } from 'vitest'
import { gapCover, gapCoverAbi, parseDeployments, toBytes32 } from '../src/index.ts'

const cover = '0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc'
const holder = '0x70997970C51812dc3A010C7d01b50e0d17dc79C8'
const deployments = parseDeployments({ chainId: 412346, tapehouse: { GapCover: cover } })

test("a holder's covers are read from its purchases, those released left out, at one block", async () => {
  const filters: unknown[] = []
  const tags: string[] = []
  const client = createClient({
    transport: custom({
      async request({ method, params }) {
        if (method === 'eth_blockNumber') return '0x2a'
        if (method === 'eth_getLogs') {
          filters.push(params)
          return [3n, 5n].map((id, i) => ({
            address: cover,
            topics: encodeEventTopics({ abi: gapCoverAbi, eventName: 'Bought', args: { id, holder, symbol: toBytes32('SPY') } }),
            data: encodeAbiParameters(
              [{ type: 'uint64' }, { type: 'uint256' }, { type: 'uint256' }, { type: 'uint256' }, { type: 'uint256' }],
              [1_790_000_000_000n, 10n ** 10n, 300n, 1_000n, 4n * 10n ** 6n],
            ),
            blockNumber: `0x${(10 + i).toString(16)}`,
            blockHash: `0x${'1'.repeat(64)}`,
            transactionHash: `0x${'2'.repeat(64)}`,
            transactionIndex: '0x0',
            logIndex: `0x${i}`,
            removed: false,
          }))
        }
        const [{ data }, block] = params as [{ data: Hex }, string]
        tags.push(block)
        const { args } = decodeFunctionData({ abi: gapCoverAbi, data })
        const released = args?.[0] === 3n
        return encodeFunctionResult({
          abi: gapCoverAbi,
          functionName: 'covers',
          result: [released ? zeroAddress : holder, 1_790_000_000_000n, 300, 1_000, toBytes32('SPY'), 10n ** 10n, 4n * 10n ** 6n],
        })
      },
    }),
  })
  const held = await gapCover.holdings(client, deployments, holder)
  expect(held).toEqual([
    {
      id: 5n,
      holder,
      closesMs: 1_790_000_000_000n,
      symbol: toBytes32('SPY'),
      notional: 10n ** 10n,
      deductibleBps: 300n,
      limitBps: 1_000n,
      premium: 4n * 10n ** 6n,
    },
  ])
  expect(filters[0]).toMatchObject([{ address: cover, fromBlock: '0x0', toBlock: '0x2a' }])
  expect(new Set(tags)).toEqual(new Set(['0x2a']))
})
