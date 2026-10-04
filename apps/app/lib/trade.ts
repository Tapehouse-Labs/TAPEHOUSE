// SPDX-License-Identifier: MIT OR Apache-2.0
import { tokens } from "./amounts";
import { percent } from "./earn";
import type { Note } from "./notes";

/** The slippage a ticket starts at, in percent, as the trader would type it. */
export const DEFAULT_SLIPPAGE = "0.5";

const BPS = 10_000n;
const date = new Intl.DateTimeFormat("en-GB", { timeZone: "UTC", day: "numeric", month: "short", year: "numeric" });

/** A slippage typed in percent, to two decimals, in basis points from 0 to 10,000; undefined otherwise. */
export function parseSlippage(typed: string): bigint | undefined {
  const match = /^(\d+)(?:\.(\d{0,2}))?$/.exec(typed.trim());
  if (!match) return undefined;
  const bps = BigInt(match[1]!) * 100n + BigInt((match[2] ?? "").padEnd(2, "0"));
  return bps <= BPS ? bps : undefined;
}

/** `amount` and `bps` of it more, rounded up: the most an order pays when its price may move by that much. */
export function withSlippage(amount: bigint, bps: bigint): bigint {
  return amount + (amount * bps + BPS - 1n) / BPS;
}

/** A slippage in basis points as a percentage, without trailing zeros. */
export function slippageWords(bps: bigint): string {
  const fraction = (bps % 100n).toString().padStart(2, "0").replace(/0+$/, "");
  return `${bps / 100n}${fraction ? `.${fraction}` : ""}%`;
}

/** The time from `nowMs` to `untilMs` in days and hours, or hours and minutes. */
export function countdown(untilMs: bigint, nowMs: number): string {
  const left = Number(untilMs) - nowMs;
  if (left <= 0) return "now";
  const minutes = Math.floor(left / 60_000);
  if (minutes < 1) return "under a minute";
  const days = Math.floor(minutes / 1_440);
  const hours = Math.floor((minutes % 1_440) / 60);
  if (days > 0) return `${days}d ${hours}h`;
  return hours > 0 ? `${hours}h ${minutes % 60}m` : `${minutes}m`;
}

/** A Uniswap fee tier, in hundredths of a basis point, as a percentage. */
export function feeTier(fee: number): string {
  return `${fee / 10_000}%`;
}

/** A gap in millionths as a percentage, to a tenth. */
export function gapPercent(ppm: bigint): string {
  const tenths = (ppm + 500n) / 1_000n;
  return `${tenths / 10n}.${tenths % 10n}%`;
}

/** The least a short's sale may fetch per token: the band's low edge, or its centre while the restriction holds. */
export function saleFloor(quote: { low: bigint; mid: bigint }, restricted: boolean) {
  return restricted
    ? { price: quote.mid, words: "the band's centre" }
    : { price: quote.low, words: "the band's low edge" };
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

/** What a short sale says before it is signed. */
export function shortNotes({
  asset,
  borrowRate,
  fee,
  restricted,
}: {
  asset: string;
  borrowRate: bigint;
  fee: number;
  restricted: boolean;
}): Note[] {
  const notes: Note[] = [];
  if (restricted)
    notes.push({
      term: "Restricted",
      text: `${asset} fell 10% from its last close: a short sells only at no less than the band's centre until the next UTC day ends.`,
      tone: "warning",
    });
  notes.push(
    {
      term: "Isolated",
      text: `One short per account and Stock Token, on USDG alone, margined at the band's high edge; it nets against nothing in your margin account.`,
    },
    { term: "Lending fee", text: `${percent(borrowRate)} a year on what you borrow, paid in ${asset} to its lenders.` },
    {
      term: "Pool",
      text: `The sale and the buy-back trade through ${asset}'s pool with USDG, at its ${feeTier(fee)} fee tier.`,
    },
    {
      term: "Recalls",
      text: `When ${asset}'s lenders recall what the shorts borrowed and the notice runs out, every short of ${asset} buys its part back through the pool, in proportion.`,
    },
    {
      term: "Liquidation",
      text: `Once a short falls short at the band's high edge, anyone may buy it back, at no more than 1% above the band's centre and never above its high edge, for 0.5% of what it cost; while the market is not open, only once it falls short at both edges. What its USDG cannot buy back falls on ${asset}'s lenders.`,
    },
  );
  return notes;
}

/** What gap cover on `asset` says before it is bought. */
export function coverPurchaseNotes(asset: string): Note[] {
  return [
    {
      term: "Pays",
      text: `The fall of ${asset} from its Chainlink feed's last round before the market closes to its first round once it reopens, when that round lands inside the band within fifteen minutes; otherwise the median of the band's centre over those minutes. A series that cannot settle so is void and refunds its premium.`,
    },
    {
      term: "Layer",
      text: "It pays falls only, the layer between the deductible and the limit, on a notional in USDG, and nothing if the price rises or falls less than the deductible.",
    },
    {
      term: "Holding",
      text: "The cover is not transferable and is settled in cash: you hold no token. Your account claims a payout once its series settles, and the premium back if it is void.",
    },
  ];
}

/** What a basket says before a mint or a deposit: its target, a pending one, its rebalance rule and what freezes it. */
export function basketNotes({
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
}): Note[] {
  const holds = (units: readonly bigint[]) => assets.map((a, i) => `${tokens(units[i] ?? 0n)} ${a}`).join(" and ");
  const notes: Note[] = [{ term: "Target", text: `The target is ${holds(target)} per share: it prices the first mint and steers rebalances, while a share holds its part of what ${basket} holds.` }];
  if (pending && pending.effectiveAt > now && pending.units.length > 0)
    notes.push({
      term: "Pending",
      text: `The owner proposed ${holds(pending.units)} from ${date.format(new Date(Number(pending.effectiveAt) * 1000))}, seven days after the proposal, so a holder who disagrees redeems first.`,
      tone: "warning",
    });
  notes.push(
    {
      term: "Rebalance",
      text: `Anyone may trade with ${basket} toward its target, only when what comes in at the bands' low edges is worth at least what goes out at their high edges.`,
    },
    {
      term: "Frozen",
      text: `A paused token, or a blocklist on the accounts or on ${basket} itself, freezes ${basket} in the accounts (BasketFrozen); a blocklisted address can neither send nor receive its shares; a burn by the issuer falls on every holder at once.`,
    },
    {
      term: "Margin",
      text: `Shares enter the cross position only, never an isolated one nor while the account holds one of its tokens in an isolated position, and are margined as the Stock Tokens they hold, which count against their assets' caps (AssetCapExceeded). A margin position that falls short has its baskets unwrapped and their tokens sold as any other.`,
    },
  );
  return notes;
}
