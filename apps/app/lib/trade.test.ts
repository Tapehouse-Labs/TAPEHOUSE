// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import {
  basketNotes,
  countdown,
  coverPurchaseNotes,
  DEFAULT_SLIPPAGE,
  feeTier,
  gapPercent,
  parseSlippage,
  saleFloor,
  shortNotes,
  shortStanding,
  slippageWords,
  withSlippage,
} from "./trade";

const DAY = 86_400;

describe("slippage", () => {
  test("starts at half a percent, is typed in percent and stays within 0 to 10,000 bps", () => {
    expect(DEFAULT_SLIPPAGE).toBe("0.5");
    expect(parseSlippage("0.5")).toBe(50n);
    expect(parseSlippage("0")).toBe(0n);
    expect(parseSlippage("1.25")).toBe(125n);
    expect(parseSlippage("100")).toBe(10_000n);
    expect(parseSlippage("100.01")).toBeUndefined();
    expect(parseSlippage("0.005")).toBeUndefined();
    expect(parseSlippage("")).toBeUndefined();
    expect(slippageWords(50n)).toBe("0.5%");
    expect(slippageWords(0n)).toBe("0%");
  });

  test("the most an order pays adds the slippage and rounds up", () => {
    expect(withSlippage(10_000n, 50n)).toBe(10_050n);
    expect(withSlippage(1n, 50n)).toBe(2n);
    expect(withSlippage(0n, 50n)).toBe(0n);
    expect(withSlippage(7n, 0n)).toBe(7n);
  });
});

describe("figures", () => {
  test("a countdown in days, hours and minutes", () => {
    expect(countdown(6n * 3_600_000n + 12n * 60_000n + 30_000n, 30_000)).toBe("6h 12m");
    expect(countdown(BigInt((2 * DAY + 4 * 3600) * 1000), 0)).toBe("2d 4h");
    expect(countdown(59_000n, 0)).toBe("under a minute");
    expect(countdown(0n, 1_000)).toBe("now");
  });

  test("a fee tier in hundredths of a basis point, and a gap in millionths, as percentages", () => {
    expect(feeTier(500)).toBe("0.05%");
    expect(feeTier(3_000)).toBe("0.3%");
    expect(gapPercent(124_000n)).toBe("12.4%");
    expect(gapPercent(55_000n)).toBe("5.5%");
  });
});

describe("shorts", () => {
  test("a sale's floor is the band's low edge, or its centre under the restriction", () => {
    const quote = { low: 16_820_000_000n, mid: 17_000_000_000n, high: 17_180_000_000n };
    expect(saleFloor(quote, false)).toEqual({ price: 16_820_000_000n, words: "the band's low edge" });
    expect(saleFloor(quote, true)).toEqual({ price: 17_000_000_000n, words: "the band's centre" });
  });

  test("a short of an earlier book shows only its USDG; one in deficit closes only by liquidation", () => {
    expect(shortStanding({ usdgHeld: 5n, debt: 1n, shares: 1n, epoch: 2n }, 2n)).toBe("open");
    expect(shortStanding({ usdgHeld: 0n, debt: 0n, shares: 0n, epoch: 0n }, 2n)).toBe("none");
    expect(shortStanding({ usdgHeld: 7n, debt: 0n, shares: 3n, epoch: 1n }, 2n)).toBe("earlier");
    expect(shortStanding({ usdgHeld: -4n, debt: 1n, shares: 1n, epoch: 2n }, 2n)).toBe("deficit");
  });

  test("before a sale: the lending fee, the pool's tier, recalls, liquidation and the restriction", () => {
    const notes = shortNotes({ asset: "NVDA", borrowRate: 41n * 10n ** 15n, fee: 500, restricted: false });
    const text = notes.map((n) => n.text);
    expect(notes.map((n) => n.term)).toEqual(["Isolated", "Lending fee", "Pool", "Recalls", "Liquidation"]);
    expect(text[1]).toBe("4.10% a year on what you borrow, paid in NVDA to its lenders.");
    expect(text[2]).toContain("0.05% fee tier");
    expect(text[4]).toContain("for 0.5% of what it cost");
    expect(text[4]).toContain("while the market is not open, only once it falls short at both edges");
    const restricted = shortNotes({ asset: "NVDA", borrowRate: 0n, fee: 500, restricted: true });
    expect(restricted[0]).toMatchObject({ term: "Restricted", tone: "warning" });
    expect(restricted[0]!.text).toContain("no less than the band's centre until the next UTC day ends");
  });
});

describe("gap cover", () => {
  test("before a purchase: what it pays, from which round to which, and what it is not", () => {
    const text = coverPurchaseNotes("SPY").map((n) => n.text);
    expect(text[0]).toContain("last round before the market closes to its first round once it reopens");
    expect(text[0]).toContain("void and refunds its premium");
    expect(text[1]).toContain("falls only");
    expect(text[2]).toContain("not transferable");
  });
});

describe("baskets", () => {
  test("before a mint: the target, a pending one with its day, the rebalance rule and the frozen states", () => {
    const notes = basketNotes({
      basket: "PAIR",
      assets: ["NVDA", "SPY"],
      target: [5n * 10n ** 17n, 2n * 10n ** 17n],
      pending: { units: [4n * 10n ** 17n, 3n * 10n ** 17n], effectiveAt: 1_791_000_000n },
      now: 1_790_000_000n,
    });
    expect(notes[0]!.text).toContain("The target is 0.5 NVDA and 0.2 SPY per share");
    expect(notes[1]).toMatchObject({ term: "Pending", tone: "warning" });
    expect(notes[1]!.text).toContain("0.4 NVDA and 0.3 SPY from 3 Oct 2026");
    const text = notes.map((n) => n.text).join(" ");
    expect(text).toContain("worth at least what goes out at their high edges");
    expect(text).toContain("BasketFrozen");
    expect(text).toContain("AssetCapExceeded");
    expect(text).toContain("cross position only");
    const settled = basketNotes({ basket: "PAIR", assets: ["NVDA"], target: [10n ** 18n], pending: undefined, now: 0n });
    expect(settled.map((n) => n.term)).not.toContain("Pending");
  });
});
