// SPDX-License-Identifier: MIT OR Apache-2.0
import { BaseError, ContractFunctionRevertedError, decodeErrorResult, type Hex } from 'viem'
import {
  bandAbi,
  bandFeedAbi,
  basketAbi,
  marginAccountsAbi,
  morphoBandOracleAbi,
  shortPositionsAbi,
  stockTokenAbi,
  usdgAbi,
} from './generated.js'

const items = [
  ...shortPositionsAbi,
  ...marginAccountsAbi,
  ...basketAbi,
  ...bandFeedAbi,
  ...morphoBandOracleAbi,
  ...bandAbi,
  ...stockTokenAbi,
  ...usdgAbi,
]

/** Every error Tapehouse's contracts, its baskets, the band, the Stock Tokens and USDG revert with. */
export const errorsAbi = items.filter(
  (item): item is Extract<(typeof items)[number], { type: 'error' }> => item.type === 'error',
)

/** A revert, decoded: a custom error of `errorsAbi`, or Solidity's `Error(string)` and `Panic(uint256)`. */
export type Revert = { errorName: string; args: readonly unknown[] }

/** Decodes the revert behind a viem error, such as a failed `simulateContract`, `writeContract` or `readContract`. */
export function decodeRevert(error: unknown): Revert | undefined {
  if (!(error instanceof BaseError)) return undefined
  const reverted = error.walk((cause) => cause instanceof ContractFunctionRevertedError)
  if (!(reverted instanceof ContractFunctionRevertedError) || !reverted.raw) return undefined
  return decodeRevertData(reverted.raw)
}

/** Decodes raw revert data. */
export function decodeRevertData(data: Hex): Revert | undefined {
  try {
    const { errorName, args } = decodeErrorResult({ abi: errorsAbi, data })
    return { errorName, args: args ?? [] }
  } catch {
    return undefined
  }
}
