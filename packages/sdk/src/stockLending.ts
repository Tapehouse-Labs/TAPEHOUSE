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

/**
 * `asset`'s lending vault at one block: the share of its tokens lent, with 18 decimals; what its borrower pays and its
 * lenders earn, a year with 18 decimals, the owner's fee out; what it holds for recalls, holds free and lent, in the
 * token's raw units; what the borrower may still take, at most `maxUtilization` and never what recalls wait for; its
 * tokens in all; the ticket next in turn; the share it lends at most, in basis points; how long a recall's notice runs,
 * in seconds; and the owner's share of the fee, in basis points.
 */
export async function terms(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const address = entry(deployments.stockLending, asset, '.tapehouse.StockLending')
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: stockLendingVaultAbi, blockNumber } as const
  const [
    utilization,
    borrowRate,
    supplyRate,
    locked,
    idle,
    debt,
    borrowable,
    totalAssets,
    head,
    maxUtilization,
    notice,
    feeShare,
  ] = await Promise.all([
    readContract(client, { ...read, functionName: 'utilization' }),
    readContract(client, { ...read, functionName: 'borrowRate' }),
    readContract(client, { ...read, functionName: 'supplyRate' }),
    readContract(client, { ...read, functionName: 'locked' }),
    readContract(client, { ...read, functionName: 'idle' }),
    readContract(client, { ...read, functionName: 'debt' }),
    readContract(client, { ...read, functionName: 'borrowable' }),
    readContract(client, { ...read, functionName: 'totalAssets' }),
    readContract(client, { ...read, functionName: 'head' }),
    readContract(client, { ...read, functionName: 'MAX_UTILIZATION' }),
    readContract(client, { ...read, functionName: 'NOTICE' }),
    readContract(client, { ...read, functionName: 'feeShare' }),
  ])
  return {
    utilization,
    borrowRate,
    supplyRate,
    locked,
    idle,
    debt,
    borrowable,
    totalAssets,
    head,
    maxUtilization,
    notice,
    feeShare,
  }
}
