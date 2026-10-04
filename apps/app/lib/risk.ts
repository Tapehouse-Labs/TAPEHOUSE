// SPDX-License-Identifier: MIT OR Apache-2.0
import { bpsPercent } from "./earn";
import { price } from "./format";

/** The share of an exposure limit the current closure covered, to two decimals; at most the whole limit. */
export function exposureUsed(covered: bigint, limit: bigint): string {
  if (limit === 0n) return "No limit set";
  const hundredths = ((covered > limit ? limit : covered) * 10_000n + limit / 2n) / limit;
  return `${hundredths / 100n}.${(hundredths % 100n).toString().padStart(2, "0")}%`;
}

export type RoundStatus = "No round" | "Committing" | "Revealing" | "Awaiting clearing" | "Cleared";

/** Where a round stands: `phase` is the auction's current round, `openMs` the regular open this round is for. */
export function roundStatus(
  round: { sealMs: bigint; cleared: boolean },
  phase: { openMs: bigint; revealing: boolean },
  openMs: bigint,
): RoundStatus {
  if (round.sealMs === 0n) return "No round";
  if (round.cleared) return "Cleared";
  if (phase.openMs === openMs && openMs !== 0n) return phase.revealing ? "Revealing" : "Committing";
  return "Awaiting clearing";
}

/** The band feed's TWAP in dollars, or that it is not yet valid. */
export function twapWords(valid: boolean, twap: bigint): string {
  return valid ? price(twap) : "Not yet valid";
}

/** The premium of the token's market over its band, in basis points, as a signed percentage. */
export function premiumWords(bps: bigint): string {
  if (bps === 0n) return "0%";
  return bps > 0n ? `+${bpsPercent(bps)}` : `−${bpsPercent(-bps)}`;
}
