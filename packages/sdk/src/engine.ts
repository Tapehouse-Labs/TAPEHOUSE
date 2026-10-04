// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Client } from 'viem'
import { readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry } from './deployments.js'
import { marginAbi } from './generated.js'

/** The margin period of risk the engine margins every regime over, in seconds: two days. */
export const HORIZON = 172_800n

/** A portfolio as the engine takes it: signed token amounts with 18 decimals and USD prices with 8, in its order. */
export type Portfolio = { quantities: readonly bigint[]; prices: readonly bigint[] }

/**
 * The margin `portfolio` needs over `horizon` seconds, in USD with 18 decimals, without the open market's buffer, and
 * a bit for each asset whose liquidity input is missing. `spansClosure` counts the weekend-gap scenarios.
 */
export async function requirement(
  client: Client,
  deployments: Deployments,
  { quantities, prices }: Portfolio,
  { horizon = HORIZON, spansClosure }: { horizon?: bigint; spansClosure: boolean },
  at: At = {},
) {
  const [margin, missing] = await readContract(client, {
    address: engine(deployments),
    abi: marginAbi,
    functionName: 'requirement',
    args: [quantities, prices, horizon, spansClosure],
    ...at,
  })
  return { margin, missing }
}

/**
 * The margin `portfolio` needs now, in USD with 18 decimals, the missing bits, and the regime the band's session puts
 * it in: 0 unknown, 1 closed, 2 open, 3 closing.
 */
export async function currentRequirement(
  client: Client,
  deployments: Deployments,
  { quantities, prices }: Portfolio,
  at: At = {},
) {
  const [margin, missing, regime] = await readContract(client, {
    address: engine(deployments),
    abi: marginAbi,
    functionName: 'currentRequirement',
    args: [quantities, prices],
    ...at,
  })
  return { margin, missing, regime }
}

function engine(deployments: Deployments) {
  return entry(deployments.tapehouse, 'Margin', '.tapehouse')
}
