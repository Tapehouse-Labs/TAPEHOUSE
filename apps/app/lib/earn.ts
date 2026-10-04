// SPDX-License-Identifier: MIT OR Apache-2.0
import { CROSS } from "@tapehouse/sdk";
import { type Hex, hexToString } from "viem";
import { tokens, usdg } from "./amounts";
import type { EarnState } from "./earnRead";
import type { Note } from "./notes";

const WAD = 10n ** 18n;
const BPS = 10_000n;
const DAY = 86_400n;
const utc = new Intl.DateTimeFormat("en-GB", {
  timeZone: "UTC",
  day: "numeric",
  month: "short",
  hour: "2-digit",
  minute: "2-digit",
});

/** A rate or a share with 18 decimals as a percentage with two decimals, rounded half up. */
export function percent(wad: bigint): string {
  const hundredths = (wad * BPS + WAD / 2n) / WAD;
  return `${hundredths / 100n}.${(hundredths % 100n).toString().padStart(2, "0")}%`;
}

/** Basis points as a percentage, without trailing zeros. */
export function bpsPercent(bps: bigint | number): string {
  const value = BigInt(bps);
  const fraction = (value % 100n).toString().padStart(2, "0").replace(/0+$/, "");
  return `${value / 100n}${fraction ? `.${fraction}` : ""}%`;
}

/** A span of seconds in days and whole hours, or hours, rounded up, below a day. */
export function span(seconds: bigint): string {
  if (seconds >= DAY) {
    const days = seconds / DAY;
    const hours = (seconds % DAY) / 3_600n;
    return `${days} day${days === 1n ? "" : "s"}${hours > 0n ? `, ${hours} hour${hours === 1n ? "" : "s"}` : ""}`;
  }
  const hours = (seconds + 3_599n) / 3_600n;
  return `${hours} hour${hours === 1n ? "" : "s"}`;
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
  place: string;
  what: string;
  earns: string;
  earnsNote: string;
  leaves: string;
  leavesNote: string;
  bears: string;
  bearsNote: string;
  holds: string | undefined;
  closed?: boolean;
};

/** The position an exposure limit binds, in words. */
export function exposureName(position: Hex): string {
  return position === CROSS ? "Cross positions" : `${hexToString(position, { size: 32 })} isolated`;
}

/** The ledger's rows: the supply vault, the backstop, each Stock Token's lending vault, then the gap cover's writers. */
export function ledger(state: EarnState, lent: Record<string, bigint>): Row[] {
  const rows: Row[] = [];
  if (state.supply) {
    const { market, lender } = state.supply;
    rows.push({
      id: "supply",
      place: "Supply",
      what: "USDG",
      earns: percent(market.supplyRate),
      earnsNote: `a year now, ${percent(market.utilization)} lent`,
      leaves: "At once",
      leavesNote: `up to ${usdg(market.idle)} USDG not lent`,
      bears: "The last loss",
      bearsNote: "after the collateral and the backstop",
      holds: lender && `${usdg(lender.assets)} USDG`,
    });
  }
  if (state.backstop) {
    const { state: b, depositor } = state.backstop;
    rows.push({
      id: "backstop",
      place: "Backstop",
      what: "USDG",
      earns: "Weekend premium",
      earnsNote: `${bpsPercent(BPS - BigInt(b.premium.reserveShare))} of ${bpsPercent(b.premium.rate)} a year, as loans repay`,
      leaves: `${span(b.terms.cooldown)}, then ${span(b.terms.window)}`,
      leavesNote: "never while the market is closed",
      bears: "The second loss",
      bearsNote: "up to its limit per closure",
      holds: depositor && `${usdg(depositor.assets)} USDG`,
      closed: depositor !== undefined && depositor.maxDeposit === 0n,
    });
  }
  for (const { asset, terms } of state.lending) {
    rows.push({
      id: `lend-${asset}`,
      place: `Lend ${asset}`,
      what: "from your account",
      earns: terms.debt === 0n ? "Nothing now" : percent(terms.supplyRate),
      earnsNote: terms.debt === 0n ? `no one shorts ${asset}` : "a year, only while shorts pay",
      leaves: `${span(terms.notice)} notice`,
      leavesNote: "the recall queue, then a buy-in",
      bears: "Shorts' shortfall",
      bearsNote: "beyond their USDG, and issuer burns",
      holds: lent[asset] === undefined ? undefined : `${tokens(lent[asset])} ${asset}`,
    });
  }
  if (state.cover) {
    const { writers, writer } = state.cover;
    rows.push({
      id: "cover",
      place: "Write cover",
      what: "USDG, apart from the backstop",
      earns: "Premiums",
      earnsNote: "at release, less what covers pay",
      leaves: writers.outstanding === 0n ? "At once" : "After settlement",
      leavesNote: "locked while a cover is out",
      bears: "Covers' payouts",
      bearsNote: "each reserved in full until release",
      holds: writer && `${usdg(writer.assets)} USDG`,
      closed: writer !== undefined && writer.maxDeposit === 0n,
    });
  }
  return rows;
}

/** What lending USDG to the supply vault says before a deposit. */
export function supplyNotes(state: NonNullable<EarnState["supply"]>): Note[] {
  const { market } = state;
  const notes: Note[] = [];
  if (!market.open)
    notes.push({
      term: "Closed",
      text: "The vault takes no deposits now: deposits wait for its borrower, and deposits and withdrawals stop while USDG is paused, the vault frozen or its USDG wiped.",
      tone: "warning",
    });
  notes.push(
    {
      term: "Loss order",
      text: "Lenders are last in the loss cascade: what a liquidation leaves unpaid falls on the borrower's collateral, then on the gap backstop, and only what the backstop cannot cover is written off against every share. A wipe of the vault's USDG by its issuer is a loss to every share too, once anyone syncs it.",
    },
    {
      term: "Timing",
      text: "A loss is written off when it happens: a lender who withdrew before it loses nothing, and those who stay bear all of it. The reserve's liquidation fees are not your cover.",
    },
    {
      term: "Exit",
      text: `You may withdraw at most the USDG the vault holds, ${usdg(market.idle)} now; the ${usdg(market.debt)} lent returns as borrowers repay.`,
    },
  );
  return notes;
}

/** One exposure limit of the backstop, as the ledger lays it out. */
export type Exposure = { name: string; limit: string; covered: string; left: string; change?: string };

/**
 * The backstop's exposure limits for the current closure, the cross positions first: each limit, what it covered in
 * the closure it last counted, what is left of it, at most the USDG held, and a change that binds closures closing
 * from its date on, while one is pending at `timestamp`, in seconds.
 */
export function exposures(b: NonNullable<EarnState["backstop"]>["state"], timestamp: bigint): Exposure[] {
  return b.exposures.map((e) => ({
    name: exposureName(e.position),
    limit: usdg(e.limit),
    covered: usdg(e.covered),
    left: usdg(e.left),
    change:
      e.next !== e.current && timestamp * 1000n < e.fromMs
        ? `${usdg(e.next)} for closures from ${utc.format(new Date(Number(e.fromMs)))} UTC`
        : undefined,
  }));
}

/** What depositing with the gap backstop says before a deposit: the cascade, its limits, its gains, its premium and how to leave. */
export function backstopNotes(state: NonNullable<EarnState["backstop"]>): Note[] {
  const { state: b, depositor } = state;
  const notes: Note[] = [];
  if (depositor && depositor.maxDeposit === 0n)
    notes.push({
      term: "Closed",
      text: "Deposits are closed now: from the moment a close is recorded until a day after its reopening, and while USDG is paused, the backstop frozen or its USDG wiped.",
      tone: "warning",
    });
  notes.push({
    term: "Loss order",
      text: "The backstop is junior to lenders: when a liquidation leaves a position empty and still owing, it repays what it can before anything is written off to them, and forgives that position's premium, its own loss.",
  });
  notes.push(
    {
      term: "Limit",
      text: "It covers at most its limit a closure for the cross positions together and for each Stock Token's isolated positions, and never more than the USDG it holds; past either, lenders bear the rest. The same limit bounds what it buys at the reopening auction, counted as a cover is.",
    },
    {
      term: "Auction",
      text: "At the auction's clearing it buys what the bids leave of each lot at the clearing price. Those Stock Tokens are its depositors', in proportion to their shares at the purchase, not counted in the shares' value and never sold again; each depositor claims its own.",
    },
    {
      term: "Premium",
      text: `While the market is closed borrowers owe a premium of ${bpsPercent(b.premium.rate)} a year; ${bpsPercent(BPS - BigInt(b.premium.reserveShare))} of each premium paid comes here and ${bpsPercent(b.premium.reserveShare)} to the fee reserve. It arrives only as borrowers repay, in lumps, and joins the shares at once when anyone claims it, not streamed: ${usdg(b.premium.waiting)} USDG is waiting.`,
    },
    {
      term: "Exit",
      text: `Shares leave only through a cooldown: ${span(b.terms.cooldown)} after you start it you may redeem them for ${span(b.terms.window)}, never while the market is closed nor from the moment a coming close is recorded until ${span(b.terms.settlement)} after its reopening. Shares sent away leave the cooldown. Deposits and redemptions stop while USDG is paused, the backstop frozen or its USDG wiped.`,
    },
    {
      term: "Seed",
      text: `The team's seed is an ordinary deposit: the owner's shares are worth ${usdg(b.seed.assets)} USDG and leave the same way.`,
    },
  );
  return notes;
}

/** What lending a Stock Token says before it is lent: the vault's terms, the recall order and who bears what. */
export function lendingNotes(asset: string, terms: EarnState["lending"][number]["terms"]): Note[] {
  return [
    {
      term: "Recalls",
      text: `A recall takes the vault's free tokens at once; then new lending and repayments meet the queue's tickets in their order before anything is free; then the borrower within ${span(terms.notice)}, in calendar time; then a buy-in anyone may call. No one takes idle tokens first come, first served, and a ticket next in turn that no one settles holds the ones behind it until someone does, which anyone may.`,
    },
    {
      term: "Lent",
      text: `${percent(terms.utilization)} of the vault's ${asset} is lent, ${tokens(terms.locked)} is held for recalls and ${tokens(terms.borrowable)} more may be borrowed: at most ${bpsPercent(terms.maxUtilization)}, never what recalls wait for.`,
    },
    {
      term: "Earnings",
      text:
        terms.debt === 0n
          ? `The borrower is Tapehouse's short positions, so lenders earn only what the shorts pay: no one shorts ${asset} now, so lending it earns nothing.`
          : `The borrower is Tapehouse's short positions, so lenders earn only what the shorts pay: ${percent(terms.supplyRate)} a year now, after the owner's ${bpsPercent(terms.feeShare)} of the fee. The fee is paid in ${asset}; dividends and splits pass through.`,
    },
    {
      term: "Losses",
      text: `A short's loss beyond its USDG is written off on the vault, a loss to its lenders through their shares, recalled or not, as far as ${asset}'s shorts still owe it. A buy-in buys through ${asset}'s pool at no more than 1% above the band's centre; one the pool cannot meet waits while the notice runs.`,
    },
    {
      term: "Issuer",
      text: `A burn from the vault falls on its lenders alone, in proportion: the free tokens first, then those held for the latest tickets. While ${asset} is paused or the vault blocklisted nothing leaves and no buy-in runs, though recalls may be made and notices keep running. If the issuer fails, take back and withdraw to your own address, where you claim from the issuer yourself.`,
    },
  ];
}

/** What writing gap cover says before a deposit: what covers reserve, when premiums arrive and when writers are locked. */
export function coverNotes(state: NonNullable<EarnState["cover"]>): Note[] {
  const { writers, writer } = state;
  const notes: Note[] = [];
  if (writer && writer.maxDeposit === 0n)
    notes.push({
      term: "Closed",
      text: "Deposits are closed now: while a cover is out over the coming closure they stay open only while the session is open and before its sales end, never while another closure's covers are out, and never while USDG is paused, the cover frozen or its USDG wiped.",
      tone: "warning",
    });
  notes.push(
    {
      term: "Reserve",
      text: `This is the gap cover's own vault, apart from the backstop. Every cover reserves its whole payout from the writers' USDG until it is released: ${usdg(writers.reserved)} reserved, ${usdg(writers.capacity)} free.`,
    },
    {
      term: "Premiums",
      text: `Premiums join the writers' USDG only at release, less what each cover pays: ${usdg(writers.premiums)} USDG waits on ${writers.outstanding} cover${writers.outstanding === 1n ? "" : "s"}.`,
    },
    {
      term: "Exit",
      text: "Redemptions open only once every cover is released, so from the first cover sold you are locked until the closure's series settle or void; a series no one settles voids after a week, and a void cover refunds its premium to its buyer, so writers earn nothing on it. A wipe of the cover's USDG falls on the writers after the premiums and credits it holds.",
    },
  );
  return notes;
}
