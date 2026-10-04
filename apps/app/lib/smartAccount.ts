// SPDX-License-Identifier: MIT OR Apache-2.0

/** How many operations the sponsor paymaster pays for each account, besides its creation. */
export const FREE_OPERATIONS = 3;

/** The notches drawn across the stem, one per free action, lit while it is left. */
export function notches(left: bigint): boolean[] {
  return Array.from({ length: FREE_OPERATIONS }, (_, i) => BigInt(i) < left);
}

/** The line under the smart account's seal. */
export function freeActions({ deployed, left }: { deployed: boolean; left: bigint }): string {
  if (!deployed) return `Creating it costs you nothing, and so do your first ${left} actions.`;
  if (left === 0n) return "Free actions used. You pay gas from here.";
  return `${left} free action${left === 1n ? "" : "s"} left`;
}
