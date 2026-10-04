// SPDX-License-Identifier: MIT OR Apache-2.0

/** How long before a weekend or holiday close the requirement starts rising to the closed one, in milliseconds. */
export const RAMP_MS = 25_200_000n;

const BPS = 10_000n;
const HOUR_MS = 3_600_000n;
const YEAR_HOURS = 8_760n;

/** The regime the engine reports, in words: 0 unknown, 1 closed, 2 open, 3 closing. */
export function regimeWord(regime: number): string {
  return ["Session unknown, counted as closed", "Closed", "Open, with a 25% buffer", "Closing"][regime] ??
    "Session unknown, counted as closed";
}

export type Mark = {
  key: "now" | "ramp" | "close" | "reopen" | "regular";
  hours: number | undefined;
  requirement: bigint;
};

/**
 * The account's week from now, each mark with the requirement it needs there, in USD with 18 decimals. `current` is
 * the engine's requirement now; `open` and `closed` its requirement over two days without and with the weekend-gap
 * scenarios; `gross` the gross exposure. The open market holds a 25% buffer over `open`; a closure needs the largest
 * of that, `closed` and the gross exposure over `weekendLeverage`. The band sets `boundaryMs` only on the last trading
 * day before a weekend or holiday close, the end of its post-market; the ramp runs the seven hours before it. Before that
 * day the close has no hour yet. While closed, `boundaryMs` is the 24/5 reopen, which brings the buffered open
 * requirement back; the regular open is when the reopening auction sells.
 */
export function week({
  regime,
  current,
  open,
  closed,
  gross,
  weekendLeverage,
  boundaryMs,
  nowMs,
}: {
  regime: number;
  current: bigint;
  open: bigint;
  closed: bigint;
  gross: bigint;
  weekendLeverage: number;
  boundaryMs: bigint;
  nowMs: number;
}): Mark[] {
  const buffered = ceil(open * 5n, 4n);
  const across = max(buffered, closed, ceil(gross * BPS, BigInt(weekendLeverage)));
  const marks: Mark[] = [{ key: "now", hours: undefined, requirement: current }];
  const reopensIn = regime === 1 && boundaryMs > BigInt(nowMs) ? Number((boundaryMs - BigInt(nowMs)) / HOUR_MS) : undefined;
  const after: Mark[] = [
    { key: "reopen", hours: reopensIn, requirement: buffered },
    { key: "regular", hours: undefined, requirement: buffered },
  ];
  if (regime === 2 || regime === 3) {
    if (boundaryMs === 0n) return [...marks, { key: "close", hours: undefined, requirement: across }, ...after];
    const ahead = boundaryMs - BigInt(nowMs);
    if (regime === 2 && ahead > RAMP_MS) {
      marks.push({ key: "ramp", hours: Number((ahead - RAMP_MS) / HOUR_MS), requirement: buffered });
    }
    marks.push({ key: "close", hours: Number((ahead > 0n ? ahead : 0n) / HOUR_MS), requirement: across });
  }
  return [...marks, ...after];
}

/**
 * The most USDG a position may borrow now: what its equity holds above its requirement and, before a closure, above a
 * fifth of its gross exposure, within the debt cap and, before a closure, the weekend debt cap. Never below zero.
 */
export function capacity({
  equity,
  requirement,
  gross,
  weekendAhead,
  weekendLeverage,
  debt,
  debtCap,
  weekendDebtCap,
}: {
  equity: bigint;
  requirement: bigint;
  gross: bigint;
  weekendAhead: boolean;
  weekendLeverage: number;
  debt: bigint;
  debtCap: bigint;
  weekendDebtCap: bigint;
}): bigint {
  const room = [(equity - requirement) / 10n ** 12n, debtCap - debt];
  if (weekendAhead) {
    room.push((equity - ceil(gross * BPS, BigInt(weekendLeverage))) / 10n ** 12n, weekendDebtCap - debt);
  }
  const least = room.reduce((a, b) => (b < a ? b : a));
  return least > 0n ? least : 0n;
}

/**
 * What a position must add to its equity, or take off its gross exposure, for the gross exposure to stay within
 * `weekendLeverage` times its equity at the close, in USD with 18 decimals; undefined when it fits already.
 */
export function weekendPlan({
  gross,
  equity,
  weekendLeverage,
}: {
  gross: bigint;
  equity: bigint;
  weekendLeverage: number;
}): { add: bigint; reduce: bigint } | undefined {
  const floor = ceil(gross * BPS, BigInt(weekendLeverage));
  if (equity >= floor) return undefined;
  const allowed = equity > 0n ? (equity * BigInt(weekendLeverage)) / BPS : 0n;
  return { add: floor - equity, reduce: gross - allowed };
}

/** What `debt`, in USDG, owes in premium over `hours` of closure at `rate` basis points a year, rounded down. */
export function premiumOver(debt: bigint, rate: number, hours: number): bigint {
  return (debt * BigInt(rate) * BigInt(hours)) / (YEAR_HOURS * BPS);
}

function ceil(a: bigint, b: bigint): bigint {
  return a <= 0n ? 0n : (a + b - 1n) / b;
}

function max(...values: bigint[]): bigint {
  return values.reduce((a, b) => (b > a ? b : a));
}
