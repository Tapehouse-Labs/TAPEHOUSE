// SPDX-License-Identifier: MIT OR Apache-2.0
import { CROSS, toBytes32 } from "@tapehouse/sdk";
import { describe, expect, test } from "vitest";
import {
  backstopNotes,
  cooldown,
  coverNotes,
  exposureName,
  exposures,
  ledger,
  lendingNotes,
  span,
  supplyNotes,
} from "./earn";
import type { EarnState } from "./earnRead";
import { formatters } from "./format";
import { say } from "./say";

const f = formatters("en");
const D = 10n ** 18n;
const NVDA = "0x525c2aBA45F66987217323E8a05EA400C65D06DC";
const holding = (assets: bigint, maxDeposit = 2n ** 255n) => ({
  shares: assets * 10n ** 6n,
  assets,
  maxWithdraw: assets,
  maxRedeem: assets * 10n ** 6n,
  maxDeposit,
});
const terms = {
  utilization: 41n * 10n ** 16n,
  borrowRate: 4n * 10n ** 16n,
  supplyRate: 1_394n * 10n ** 13n,
  locked: 3n * D,
  idle: 10n * D,
  debt: 13n * D,
  borrowable: 9n * D,
  totalAssets: 26n * D,
  head: 2n,
  maxUtilization: 9_000,
  notice: 86_400n,
  feeShare: 1_500,
};

function state(over: Partial<EarnState> = {}): EarnState {
  return {
    blockNumber: 42n,
    timestamp: 1_790_000_000n,
    supply: {
      market: {
        supplyRate: 312n * 10n ** 14n,
        borrowRate: 5n * 10n ** 16n,
        utilization: 61n * 10n ** 16n,
        idle: 2_140_000_000n,
        debt: 3_348_000_000n,
        totalAssets: 5_488_000_000n,
        open: true,
      },
      lender: holding(1_204_100_000n),
    },
    backstop: {
      state: {
        held: 9_000_000_000n,
        shares: 9_000n * 10n ** 12n,
        closureMs: 1_789_000_000_000n,
        exposures: [
          { position: CROSS, limit: 5_000_000_000n, current: 5_000_000_000n, next: 2_000_000_000n, fromMs: 1_791_000_000_000n, left: 4_500_000_000n, covered: 500_000_000n },
          { position: toBytes32("NVDA"), limit: 1_000_000_000n, current: 1_000_000_000n, next: 1_000_000_000n, fromMs: 0n, left: 1_000_000_000n, covered: 0n },
        ],
        premium: { rate: 500, reserveShare: 1_000, waiting: 42_000_000n },
        terms: { cooldown: 604_800n, window: 518_400n, settlement: 86_400n },
        seed: { owner: NVDA, assets: 1_000_000_000n },
      },
      depositor: { ...holding(300_000_000n), cooldown: undefined, gains: { [NVDA]: 2n * 10n ** 17n } },
      tokens: [{ asset: "NVDA", token: NVDA }],
    },
    lending: [{ asset: "NVDA", token: NVDA, terms }],
    cover: {
      writers: {
        held: 1_000_000_000n,
        shares: 1_000n * 10n ** 12n,
        reserved: 400_000_000n,
        capacity: 600_000_000n,
        premiums: 12_000_000n,
        owed: 0n,
        outstanding: 2n,
        sales: undefined,
      },
      writer: holding(0n, 0n),
    },
    ...over,
  };
}

describe("figures", () => {
  test("rates, basis points and spans read as words and percentages", () => {
    expect(f.percent(312n * 10n ** 14n)).toBe("3.12%");
    expect(f.percent(5n * 10n ** 15n)).toBe("0.50%");
    expect(f.percent(0n)).toBe("0.00%");
    expect(f.bpsPercent(9_000)).toBe("90%");
    expect(f.bpsPercent(1_250n)).toBe("12.5%");
    expect(say(span(604_800n))).toBe("7 days");
    expect(say(span(86_400n))).toBe("1 day");
    expect(say(span(7_200n))).toBe("2 hours");
    expect(say(span(604_799n))).toBe("6 days, 23 hours");
    expect(say(exposureName(CROSS))).toBe("Cross positions");
    expect(say(exposureName(toBytes32("NVDA")))).toBe("NVDA isolated");
  });

  test("spans and exposures read in the other languages", () => {
    expect(say(span(86_400n), "es")).toBe("1 día");
    expect(say(span(604_799n), "pt")).toBe("6 dias, 23 horas");
    expect(say(exposureName(toBytes32("NVDA")), "es")).toBe("NVDA aislado");
    expect(say(exposureName(CROSS), "pt")).toBe("Posições cruzadas");
  });

  test("a cooldown is cooling, then its window, then lapsed", () => {
    const c = { from: 100n, until: 200n };
    expect(cooldown(undefined, 0n)).toEqual({ kind: "none" });
    expect(cooldown(c, 40n)).toEqual({ kind: "cooling", opensIn: 60n });
    expect(cooldown(c, 100n)).toEqual({ kind: "window", closesIn: 100n });
    expect(cooldown(c, 201n)).toEqual({ kind: "lapsed" });
  });
});

const said = (row: ReturnType<typeof ledger>[number], locale: "en" | "es" | "pt" = "en") =>
  Object.fromEntries(Object.entries(row).map(([k, v]) => [k, typeof v === "object" && v && "key" in v ? say(v, locale) : v]));

describe("the ledger", () => {
  test("each place answers what it earns, how it leaves and what you hold", () => {
    const rows = ledger(f, state(), { NVDA: 2n * D }).map((r) => said(r));
    expect(rows.map((r) => r.id)).toEqual(["supply", "backstop", "lend-NVDA", "cover"]);
    expect(rows[0]).toMatchObject({ bears: "The last loss", earns: "3.12%", earnsNote: "a year now, 61.00% lent", leaves: "At once", holds: "1,204.10 USDG" });
    expect(rows[1]).toMatchObject({ bears: "The second loss", earnsNote: "90% of 5% a year, as loans repay", leaves: "7 days, then 6 days", holds: "300.00 USDG" });
    expect(rows[2]).toMatchObject({ earns: "1.39%", leaves: "1 day notice", holds: "2 NVDA" });
    expect(rows[3]).toMatchObject({ leaves: "After settlement", holds: "0.00 USDG", closed: true });
  });

  test("a vault no one borrows from names no yield, and a disconnected visitor holds nothing", () => {
    const quiet = state({ lending: [{ asset: "NVDA", token: NVDA, terms: { ...terms, debt: 0n } }] });
    const rows = ledger(f, { ...quiet, supply: { market: quiet.supply!.market } }, {});
    expect(said(rows.find((r) => r.id === "lend-NVDA")!)).toMatchObject({ earns: "Nothing now", earnsNote: "no one shorts NVDA", holds: undefined });
    expect(rows[0]!.holds).toBeUndefined();
  });

  test("only what the registry deploys is a row", () => {
    expect(ledger(f, state({ cover: undefined, backstop: undefined, lending: [] }), {}).map((r) => r.id)).toEqual(["supply"]);
  });

  test("a row reads in Spanish with the language's separators", () => {
    const es = formatters("es");
    const row = said(ledger(es, state(), {})[0]!, "es");
    expect(row).toMatchObject({ bears: "La última pérdida", leaves: "De inmediato", holds: "1204,10 USDG" });
  });
});

const texts = (notes: ReturnType<typeof supplyNotes>) => notes.map((n) => say(n.message));

describe("notes before a deposit", () => {
  test("supply says lenders are last and a withdrawal is capped by what the vault holds", () => {
    const notes = supplyNotes(f, state().supply!);
    const words = texts(notes);
    expect(words[0]).toContain("Lenders are last in the loss cascade");
    expect(say(notes[0]!.term!)).toBe("Loss order");
    expect(words[1]).toContain("a lender who withdrew before it loses nothing");
    expect(words[2]).toBe("You may withdraw at most the USDG the vault holds, 2,140.00 now; the 3,348.00 lent returns as borrowers repay.");
    const stopped = supplyNotes(f, { market: { ...state().supply!.market, open: false } });
    expect(stopped[0]).toMatchObject({ tone: "warning" });
  });

  test("the backstop lays out each limit with its pending change, and states the premium's share and the cooldown", () => {
    const view = (timestamp: bigint) =>
      exposures(f, state().backstop!.state, timestamp).map((e) => ({ ...e, name: say(e.name), change: e.change && say(e.change) }));
    expect(view(1_790_000_000n)).toEqual([
      { name: "Cross positions", limit: "5,000.00", covered: "500.00", left: "4,500.00", change: "2,000.00 for closures from 3 Oct, 04:00 UTC" },
      { name: "NVDA isolated", limit: "1,000.00", covered: "0.00", left: "1,000.00", change: undefined },
    ]);
    expect(view(1_791_000_000n)[0]!.change).toBeUndefined();
    const words = texts(backstopNotes(f, state().backstop!));
    expect(words.find((n) => n.includes("past either, lenders bear the rest"))).toContain("counted as a cover is");
    expect(words.find((n) => n.includes("premium of 5% a year"))).toContain("90% of each premium paid comes here and 10% to the fee reserve");
    expect(words.find((n) => n.includes("cooldown"))).toContain("7 days after you start it you may redeem them for 6 days");
    expect(words.at(-1)).toBe("The team's seed is an ordinary deposit: the owner's shares are worth 1,000.00 USDG and leave the same way.");
    const closed = backstopNotes(f, { ...state().backstop!, depositor: { ...state().backstop!.depositor!, maxDeposit: 0n } });
    expect(closed[0]).toMatchObject({ tone: "warning" });
  });

  test("lending names a yield only while shorts pay, and states the recall order", () => {
    const notes = lendingNotes(f, "NVDA", terms);
    const words = texts(notes);
    expect(say(notes[0]!.term!)).toBe("Recalls");
    expect(words[0]).toContain("then the borrower within 1 day, in calendar time; then a buy-in anyone may call");
    expect(words[1]).toBe("41.00% of the vault's NVDA is lent, 3 is held for recalls and 9 more may be borrowed: at most 90%, never what recalls wait for.");
    expect(words[2]).toContain("1.39% a year now, after the owner's 15% of the fee");
    expect(texts(lendingNotes(f, "NVDA", { ...terms, debt: 0n }))[2]).toContain("no one shorts NVDA now, so lending it earns nothing");
  });

  test("writing cover states what covers reserve, when premiums arrive and when deposits are closed", () => {
    const notes = coverNotes(f, state().cover!);
    const words = texts(notes);
    expect(notes[0]).toMatchObject({ tone: "warning" });
    expect(words[1]).toContain("400.00 reserved, 600.00 free");
    expect(words[2]).toContain("12.00 USDG waits on 2 covers");
  });
});
