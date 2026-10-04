// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Client } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry } from './deployments.js'
import { stockLendingVaultAbi } from './generated.js'

/**
 * The recall ticket `id` of `asset`'s lending vault, at one block: the span of the queue it waits for, what it has taken,
 * when its notice runs out in seconds, what the vault holds for it now, and the ticket next in turn.
 */
export async function ticket(client: Client, deployments: Deployments, asset: string, id: bigint, at: At = {}) {
  const address = entry(deployments.stockLending, asset, '.tapehouse.StockLending')
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: stockLendingVaultAbi, blockNumber } as const
  const [[start, end, taken, dueAt], claimable, head] = await Promise.all([
    readContract(client, { ...read, functionName: 'ticket', args: [id] }),
    readContract(client, { ...read, functionName: 'claimable', args: [id] }),
    readContract(client, { ...read, functionName: 'head' }),
  ])
  return { start, end, taken, dueAt, claimable, head }
}
