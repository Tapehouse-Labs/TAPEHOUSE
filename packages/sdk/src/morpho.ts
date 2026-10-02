// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Client } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry } from './deployments.js'
import { decodeRevert, type Revert } from './errors.js'
import { bandAbi, morphoBandOracleAbi } from './generated.js'

/**
 * Why a Morpho oracle has no price: its band has no live leg, which a fresh `writePrices` ends (`stale`); the band holds
 * a signed halt, the issuer's pause or an unconfirmed multiplier step for the asset (`halted`); or the L2 sequencer is
 * not settled (`sequencerNotSettled`).
 */
export type NoPrice = 'stale' | 'halted' | 'sequencerNotSettled'

/**
 * A Morpho oracle's answer: the price of one whole collateral token in the loan token, with Morpho's
 * 36 + loan decimals − collateral decimals, or no price and why. Never a price of zero.
 */
export type OraclePrice = { price: bigint } | { price: undefined; noPrice: NoPrice; revert: Revert }

/**
 * The price `asset`'s Morpho oracle, `.morphoOracles.<asset>`, answers Morpho Blue. A revert with `NoAnswer` or
 * `SequencerNotSettled` is no price, never zero; the oracle's halt and the band's corporate action, read at the same
 * block, tell a stale band from a halt.
 */
export async function price(
  client: Client,
  deployments: Deployments,
  asset: string,
  at: At = {},
): Promise<OraclePrice> {
  const address = entry(deployments.morphoOracles, asset, '.morphoOracles')
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: morphoBandOracleAbi, blockNumber } as const
  try {
    return { price: await readContract(client, { ...read, functionName: 'price' }) }
  } catch (error) {
    const revert = decodeRevert(error)
    if (revert?.errorName === 'SequencerNotSettled') return { price: undefined, noPrice: 'sequencerNotSettled', revert }
    if (revert?.errorName !== 'NoAnswer') throw error
    const [[signedHalt, , , oraclePaused], band, symbol] = await Promise.all([
      readContract(client, { ...read, functionName: 'halt' }),
      readContract(client, { ...read, functionName: 'band' }),
      readContract(client, { ...read, functionName: 'symbol' }),
    ])
    const [step] = await readContract(client, {
      address: band,
      abi: bandAbi,
      functionName: 'corporateAction',
      args: [symbol],
      blockNumber,
    })
    return { price: undefined, noPrice: signedHalt || oraclePaused || step === 2 ? 'halted' : 'stale', revert }
  }
}

/**
 * `asset`'s trading halt as its Morpho oracle reports it from the band: whether a halt signed by Tapehouse's halt
 * signer holds, until when, when it was issued, and whether the issuer has paused the Stock Token's oracle.
 */
export async function halt(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const [signedHalt, until, issuedAt, oraclePaused] = await readContract(client, {
    address: entry(deployments.morphoOracles, asset, '.morphoOracles'),
    abi: morphoBandOracleAbi,
    functionName: 'halt',
    ...at,
  })
  return { signedHalt, until, issuedAt, oraclePaused }
}

/**
 * `asset`'s Morpho oracle: the band it reads, the symbol it prices, the collateral and loan tokens of its markets, the
 * factor on the band's 8-decimal price, and its owner, who may re-point it to another band. All read at one block.
 */
export async function oracle(client: Client, deployments: Deployments, asset: string, at: At = {}) {
  const address = entry(deployments.morphoOracles, asset, '.morphoOracles')
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: morphoBandOracleAbi, blockNumber } as const
  const [band, symbol, collateralToken, loanToken, scaleFactor, owner] = await Promise.all([
    readContract(client, { ...read, functionName: 'band' }),
    readContract(client, { ...read, functionName: 'symbol' }),
    readContract(client, { ...read, functionName: 'collateralToken' }),
    readContract(client, { ...read, functionName: 'loanToken' }),
    readContract(client, { ...read, functionName: 'scaleFactor' }),
    readContract(client, { ...read, functionName: 'owner' }),
  ])
  return { address, band, symbol, collateralToken, loanToken, scaleFactor, owner }
}
