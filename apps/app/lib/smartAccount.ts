// SPDX-License-Identifier: MIT OR Apache-2.0

import { type Msg, msg } from "./messages";

/** How many operations the sponsor paymaster pays for each account, besides its creation. */
export const FREE_OPERATIONS = 3;

/** The notches drawn across the stem, one per free action, lit while it is left. */
export function notches(left: bigint): boolean[] {
  return Array.from({ length: FREE_OPERATIONS }, (_, i) => BigInt(i) < left);
}

/** The line under the smart account's seal. */
export function freeActions({ deployed, left }: { deployed: boolean; left: bigint }): Msg {
  if (!deployed) return msg("wallet.smart.freeFirst", { left: Number(left) });
  if (left === 0n) return msg("wallet.smart.freeUsed");
  return msg("wallet.smart.freeLeft", { left: Number(left) });
}
