// SPDX-License-Identifier: MIT OR Apache-2.0
import { CROSS } from "@tapehouse/sdk";
import { type Hex, hexToString } from "viem";
import type { EarnState } from "./earnRead";
import type { Formatters } from "./format";
import { type Msg, msg, plain } from "./messages";
import type { Note } from "./notes";

const BPS = 10_000n;
const DAY = 86_400n;

/** A span of seconds in days and whole hours, or hours, rounded up, below a day. */
export function span(seconds: bigint): Msg {
  if (seconds >= DAY) {
    const days = Number(seconds / DAY);
    const hours = Number((seconds % DAY) / 3_600n);
    return hours > 0 ? msg("earn.span.daysHours", { days, hours }) : msg("earn.span.days", { days });
  }
  return msg("earn.span.hours", { hours: Number((seconds + 3_599n) / 3_600n) });
}

/** Where a depositor's backstop cooldown stands at `now`, in seconds. */
export type Cooldown =
  | { kind: "none" }
  | { kind: "cooling"; opensIn: bigint }
  | { kind: "window"; closesIn: bigint }
  | { kind: "lapsed" };

export function cooldown(c: { from: bigint; until: bigint } | undefined, now: bigint): Cooldown {
  if (!c) return { kind: "none" };
  if (now < c.from) return { kind: "cooling", opensIn: c.from - now };
  if (now <= c.until) return { kind: "window", closesIn: c.until - now };
  return { kind: "lapsed" };
}

/** A place money can work, and its answers to the ledger's four questions. */
export type Row = {
  id: string;
  place: Msg;
  what: Msg;
  earns: Msg;
  earnsNote: Msg;
  leaves: Msg;
  leavesNote: Msg;
  bears: Msg;
  bearsNote: Msg;
  holds: Msg | undefined;
  closed?: boolean;
};

/** The position an exposure limit binds, in words. */
export function exposureName(position: Hex): Msg {
  return position === CROSS ? msg("earn.exposure.cross") : msg("earn.exposure.isolated", { asset: hexToString(position, { size: 32 }) });
}

const usdgAmount = (f: Formatters, value: bigint): Msg => msg("earn.usdgAmount", { amount: f.usdg(value) });

/** The ledger's rows: the supply vault, the backstop, each Stock Token's lending vault, then the gap cover's writers. */
export function ledger(f: Formatters, state: EarnState, lent: Record<string, bigint>): Row[] {
  const rows: Row[] = [];
  if (state.supply) {
    const { market, lender } = state.supply;
    rows.push({
      id: "supply",
      place: msg("earn.ledger.supply.place"),
      what: plain("USDG"),
      earns: plain(f.percent(market.supplyRate)),
      earnsNote: msg("earn.ledger.supply.earnsNote", { utilization: f.percent(market.utilization) }),
      leaves: msg("earn.ledger.atOnce"),
      leavesNote: msg("earn.ledger.supply.leavesNote", { idle: f.usdg(market.idle) }),
      bears: msg("earn.ledger.supply.bears"),
      bearsNote: msg("earn.ledger.supply.bearsNote"),
      holds: lender && usdgAmount(f, lender.assets),
    });
  }
  if (state.backstop) {
    const { state: b, depositor } = state.backstop;
    rows.push({
      id: "backstop",
      place: msg("earn.ledger.backstop.place"),
      what: plain("USDG"),
      earns: msg("earn.ledger.backstop.earns"),
      earnsNote: msg("earn.ledger.backstop.earnsNote", {
        share: f.bpsPercent(BPS - BigInt(b.premium.reserveShare)),
        rate: f.bpsPercent(b.premium.rate),
      }),
      leaves: msg("earn.ledger.backstop.leaves", { cooldown: span(b.terms.cooldown), window: span(b.terms.window) }),
      leavesNote: msg("earn.ledger.backstop.leavesNote"),
      bears: msg("earn.ledger.backstop.bears"),
      bearsNote: msg("earn.ledger.backstop.bearsNote"),
      holds: depositor && usdgAmount(f, depositor.assets),
      closed: depositor !== undefined && depositor.maxDeposit === 0n,
    });
  }
  for (const { asset, terms } of state.lending) {
    rows.push({
      id: `lend-${asset}`,
      place: msg("earn.ledger.lend.place", { asset }),
      what: msg("earn.ledger.lend.what"),
      earns: terms.debt === 0n ? msg("earn.ledger.lend.nothingNow") : plain(f.percent(terms.supplyRate)),
      earnsNote: terms.debt === 0n ? msg("earn.ledger.lend.earnsNoteNone", { asset }) : msg("earn.ledger.lend.earnsNote"),
      leaves: msg("earn.ledger.lend.leaves", { notice: span(terms.notice) }),
      leavesNote: msg("earn.ledger.lend.leavesNote"),
      bears: msg("earn.ledger.lend.bears"),
      bearsNote: msg("earn.ledger.lend.bearsNote"),
      holds: lent[asset] === undefined ? undefined : msg("earn.tokenAmount", { amount: f.tokens(lent[asset]), asset }),
    });
  }
  if (state.cover) {
    const { writers, writer } = state.cover;
    rows.push({
      id: "cover",
      place: msg("earn.ledger.cover.place"),
      what: msg("earn.ledger.cover.what"),
      earns: msg("earn.ledger.cover.earns"),
      earnsNote: msg("earn.ledger.cover.earnsNote"),
      leaves: writers.outstanding === 0n ? msg("earn.ledger.atOnce") : msg("earn.ledger.cover.afterSettlement"),
      leavesNote: msg("earn.ledger.cover.leavesNote"),
      bears: msg("earn.ledger.cover.bears"),
      bearsNote: msg("earn.ledger.cover.bearsNote"),
      holds: writer && usdgAmount(f, writer.assets),
      closed: writer !== undefined && writer.maxDeposit === 0n,
    });
  }
  return rows;
}

const term = (name: "closed" | "lossOrder" | "timing" | "exit" | "limit" | "auction" | "premium" | "seed" | "recalls" | "lent" | "earnings" | "losses" | "issuer" | "reserve" | "premiums"): Msg =>
  msg(`earn.term.${name}`);

/** What lending USDG to the supply vault says before a deposit. */
export function supplyNotes(f: Formatters, state: NonNullable<EarnState["supply"]>): Note[] {
  const { market } = state;
  const notes: Note[] = [];
  if (!market.open) notes.push({ term: term("closed"), message: msg("earn.notes.supply.closed"), tone: "warning" });
  notes.push(
    { term: term("lossOrder"), message: msg("earn.notes.supply.lossOrder") },
    { term: term("timing"), message: msg("earn.notes.supply.timing") },
    { term: term("exit"), message: msg("earn.notes.supply.exit", { idle: f.usdg(market.idle), debt: f.usdg(market.debt) }) },
  );
  return notes;
}

/** One exposure limit of the backstop, as the ledger lays it out. */
export type Exposure = { name: Msg; limit: string; covered: string; left: string; change?: Msg };

/**
 * The backstop's exposure limits for the current closure, the cross positions first: each limit, what it covered in
 * the closure it last counted, what is left of it, at most the USDG held, and a change that binds closures closing
 * from its date on, while one is pending at `timestamp`, in seconds.
 */
export function exposures(f: Formatters, b: NonNullable<EarnState["backstop"]>["state"], timestamp: bigint): Exposure[] {
  return b.exposures.map((e) => ({
    name: exposureName(e.position),
    limit: f.usdg(e.limit),
    covered: f.usdg(e.covered),
    left: f.usdg(e.left),
    change:
      e.next !== e.current && timestamp * 1000n < e.fromMs
        ? msg("earn.exposure.change", { limit: f.usdg(e.next), date: f.utc(new Date(Number(e.fromMs))) })
        : undefined,
  }));
}

/** What depositing with the gap backstop says before a deposit: the cascade, its limits, its gains, its premium and how to leave. */
export function backstopNotes(f: Formatters, state: NonNullable<EarnState["backstop"]>): Note[] {
  const { state: b, depositor } = state;
  const notes: Note[] = [];
  if (depositor && depositor.maxDeposit === 0n)
    notes.push({ term: term("closed"), message: msg("earn.notes.backstop.closed"), tone: "warning" });
  notes.push(
    { term: term("lossOrder"), message: msg("earn.notes.backstop.lossOrder") },
    { term: term("limit"), message: msg("earn.notes.backstop.limit") },
    { term: term("auction"), message: msg("earn.notes.backstop.auction") },
    {
      term: term("premium"),
      message: msg("earn.notes.backstop.premium", {
        rate: f.bpsPercent(b.premium.rate),
        share: f.bpsPercent(BPS - BigInt(b.premium.reserveShare)),
        reserve: f.bpsPercent(b.premium.reserveShare),
        waiting: f.usdg(b.premium.waiting),
      }),
    },
    {
      term: term("exit"),
      message: msg("earn.notes.backstop.exit", {
        cooldown: span(b.terms.cooldown),
        window: span(b.terms.window),
        settlement: span(b.terms.settlement),
      }),
    },
    { term: term("seed"), message: msg("earn.notes.backstop.seed", { assets: f.usdg(b.seed.assets) }) },
  );
  return notes;
}

/** What lending a Stock Token says before it is lent: the vault's terms, the recall order and who bears what. */
export function lendingNotes(f: Formatters, asset: string, terms: EarnState["lending"][number]["terms"]): Note[] {
  return [
    { term: term("recalls"), message: msg("earn.notes.lending.recalls", { notice: span(terms.notice) }) },
    {
      term: term("lent"),
      message: msg("earn.notes.lending.lent", {
        utilization: f.percent(terms.utilization),
        asset,
        locked: f.tokens(terms.locked),
        borrowable: f.tokens(terms.borrowable),
        max: f.bpsPercent(terms.maxUtilization),
      }),
    },
    {
      term: term("earnings"),
      message:
        terms.debt === 0n
          ? msg("earn.notes.lending.earningsNone", { asset })
          : msg("earn.notes.lending.earnings", { rate: f.percent(terms.supplyRate), share: f.bpsPercent(terms.feeShare), asset }),
    },
    { term: term("losses"), message: msg("earn.notes.lending.losses", { asset }) },
    { term: term("issuer"), message: msg("earn.notes.lending.issuer", { asset }) },
  ];
}

/** What writing gap cover says before a deposit: what covers reserve, when premiums arrive and when writers are locked. */
export function coverNotes(f: Formatters, state: NonNullable<EarnState["cover"]>): Note[] {
  const { writers, writer } = state;
  const notes: Note[] = [];
  if (writer && writer.maxDeposit === 0n)
    notes.push({ term: term("closed"), message: msg("earn.notes.cover.closed"), tone: "warning" });
  notes.push(
    {
      term: term("reserve"),
      message: msg("earn.notes.cover.reserve", { reserved: f.usdg(writers.reserved), free: f.usdg(writers.capacity) }),
    },
    {
      term: term("premiums"),
      message: msg("earn.notes.cover.premiums", { premiums: f.usdg(writers.premiums), covers: Number(writers.outstanding) }),
    },
    { term: term("exit"), message: msg("earn.notes.cover.exit") },
  );
  return notes;
}
