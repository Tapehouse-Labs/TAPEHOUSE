// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Formatters } from "./format";
import { type Msg, msg, plain } from "./messages";

/** The share of an exposure limit the current closure covered, to two decimals; at most the whole limit. */
export function exposureUsed(f: Formatters, covered: bigint, limit: bigint): Msg {
  if (limit === 0n) return msg("risk.noLimit");
  const hundredths = ((covered > limit ? limit : covered) * 10_000n + limit / 2n) / limit;
  return plain(`${hundredths / 100n}${f.decimal}${(hundredths % 100n).toString().padStart(2, "0")}%`);
}

export type RoundStatus = "none" | "committing" | "revealing" | "awaiting" | "cleared";

/** Where a round stands: `phase` is the auction's current round, `openMs` the regular open this round is for. */
export function roundStatus(
  round: { sealMs: bigint; cleared: boolean },
  phase: { openMs: bigint; revealing: boolean },
  openMs: bigint,
): RoundStatus {
  if (round.sealMs === 0n) return "none";
  if (round.cleared) return "cleared";
  if (phase.openMs === openMs && openMs !== 0n) return phase.revealing ? "revealing" : "committing";
  return "awaiting";
}

/** The band feed's TWAP in dollars, or that it is not yet valid. */
export function twapWords(f: Formatters, valid: boolean, twap: bigint): Msg {
  return valid ? plain(f.price(twap)) : msg("risk.twapInvalid");
}

/** The premium of the token's market over its band, in basis points, as a signed percentage. */
export function premiumWords(f: Formatters, bps: bigint): string {
  if (bps === 0n) return "0%";
  return bps > 0n ? `+${f.bpsPercent(bps)}` : `−${f.bpsPercent(-bps)}`;
}
