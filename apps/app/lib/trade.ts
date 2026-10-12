// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Formatters } from "./format";
import { type Msg, msg, plain } from "./messages";
import type { Note } from "./notes";

/** The slippage a ticket starts at, in percent, as the trader would type it. */
export const DEFAULT_SLIPPAGE = "0.5";

const BPS = 10_000n;

/** A slippage typed in percent, to two decimals, in basis points from 0 to 10,000; undefined otherwise. */
export function parseSlippage(f: Formatters, typed: string): bigint | undefined {
  const trimmed = typed.trim();
  if (f.decimal !== "." && trimmed.includes(".")) return undefined;
  const match = /^(\d+)(?:\.(\d{0,2}))?$/.exec(trimmed.replace(f.decimal, "."));
  if (!match) return undefined;
  const bps = BigInt(match[1]!) * 100n + BigInt((match[2] ?? "").padEnd(2, "0"));
  return bps <= BPS ? bps : undefined;
}

/** `amount` and `bps` of it more, rounded up: the most an order pays when its price may move by that much. */
export function withSlippage(amount: bigint, bps: bigint): bigint {
  return amount + (amount * bps + BPS - 1n) / BPS;
}

/** The time from `nowMs` to `untilMs` in days and hours, or hours and minutes. */
export function countdown(untilMs: bigint, nowMs: number): Msg {
  const left = Number(untilMs) - nowMs;
  if (left <= 0) return msg("trade.countdown.now");
  const minutes = Math.floor(left / 60_000);
  if (minutes < 1) return msg("trade.countdown.underMinute");
  const days = Math.floor(minutes / 1_440);
  const hours = Math.floor((minutes % 1_440) / 60);
  if (days > 0) return msg("trade.countdown.daysHours", { days, hours });
  return hours > 0 ? msg("trade.countdown.hoursMinutes", { hours, minutes: minutes % 60 }) : msg("trade.countdown.minutes", { minutes });
}

/** A Uniswap fee tier, in hundredths of a basis point, as a percentage. */
export function feeTier(f: Formatters, fee: number): string {
  return f.bpsPercent(fee / 100);
}

/** A gap in millionths as a percentage, to a tenth. */
export function gapPercent(f: Formatters, ppm: bigint): string {
  const tenths = (ppm + 500n) / 1_000n;
  return `${tenths / 10n}${f.decimal}${tenths % 10n}%`;
}

/** The least a short's sale may fetch per token: the band's low edge, or its centre while the restriction holds. */
export function saleFloor(quote: { low: bigint; mid: bigint }, restricted: boolean) {
  return restricted
    ? { price: quote.mid, words: msg("trade.saleFloor.centre") }
    : { price: quote.low, words: msg("trade.saleFloor.lowEdge") };
}

/**
 * Where `account`'s short of a token stands against its book `current`: none; open; of an earlier book, holding only
 * its USDG to take out; or in deficit, which only a liquidation closes.
 */
export function shortStanding(
  position: { usdgHeld: bigint; debt: bigint; shares: bigint; epoch: bigint },
  current: bigint,
): "none" | "open" | "earlier" | "deficit" {
  if (position.usdgHeld < 0n) return "deficit";
  if (position.shares > 0n && position.epoch < current) return "earlier";
  if (position.shares === 0n && position.usdgHeld === 0n) return "none";
  return "open";
}

type Term = "restricted" | "isolated" | "lendingFee" | "pool" | "recalls" | "liquidation" | "pays" | "layer" | "holding" | "target" | "pending" | "rebalance" | "frozen" | "margin";
const term = (name: Term): Msg => msg(`trade.term.${name}`);

/** What a short sale says before it is signed. */
export function shortNotes(
  f: Formatters,
  { asset, borrowRate, fee, restricted }: { asset: string; borrowRate: bigint; fee: number; restricted: boolean },
): Note[] {
  const notes: Note[] = [];
  if (restricted) notes.push({ term: term("restricted"), message: msg("trade.notes.short.restricted", { asset }), tone: "warning" });
  notes.push(
    { term: term("isolated"), message: msg("trade.notes.short.isolated") },
    { term: term("lendingFee"), message: msg("trade.notes.short.lendingFee", { rate: f.percent(borrowRate), asset }) },
    { term: term("pool"), message: msg("trade.notes.short.pool", { asset, tier: feeTier(f, fee) }) },
    { term: term("recalls"), message: msg("trade.notes.short.recalls", { asset }) },
    { term: term("liquidation"), message: msg("trade.notes.short.liquidation", { asset }) },
  );
  return notes;
}

/** What gap cover on `asset` says before it is bought. */
export function coverPurchaseNotes(asset: string): Note[] {
  return [
    { term: term("pays"), message: msg("trade.notes.cover.pays", { asset }) },
    { term: term("layer"), message: msg("trade.notes.cover.layer") },
    { term: term("holding"), message: msg("trade.notes.cover.holding") },
  ];
}

/** What a basket says before a mint or a deposit: its target, a pending one, its rebalance rule and what freezes it. */
export function basketNotes(
  f: Formatters,
  {
    basket,
    assets,
    target,
    pending,
    now,
  }: {
    basket: string;
    assets: readonly string[];
    target: readonly bigint[];
    pending: { units: readonly bigint[]; effectiveAt: bigint } | undefined;
    now: bigint;
  },
): Note[] {
  const holds = (units: readonly bigint[]) => plain(f.list(assets.map((a, i) => `${f.tokens(units[i] ?? 0n)} ${a}`)));
  const notes: Note[] = [{ term: term("target"), message: msg("trade.notes.basket.target", { holds: holds(target), basket }) }];
  if (pending && pending.effectiveAt > now && pending.units.length > 0)
    notes.push({
      term: term("pending"),
      message: msg("trade.notes.basket.pending", { holds: holds(pending.units), date: f.date(new Date(Number(pending.effectiveAt) * 1000)) }),
      tone: "warning",
    });
  notes.push(
    { term: term("rebalance"), message: msg("trade.notes.basket.rebalance", { basket }) },
    { term: term("frozen"), message: msg("trade.notes.basket.frozen", { basket }) },
    { term: term("margin"), message: msg("trade.notes.basket.margin") },
  );
  return notes;
}
