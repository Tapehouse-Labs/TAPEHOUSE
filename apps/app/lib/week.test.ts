// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { capacity, premiumOver, regimeWord, week, weekendPlan } from "./week";

const D = 10n ** 18n;
const HOUR = 3_600_000;
const now = 1_790_000_000_000;
const base = { open: 3_200n * D, closed: 5_000n * D, gross: 20_000n * D, weekendLeverage: 50_000 };

describe("regimeWord", () => {
  test("names each regime the engine reports, an unknown one counted as closed", () => {
    expect([0, 1, 2, 3, 9].map(regimeWord)).toEqual([
      "Session unknown, counted as closed",
      "Closed",
      "Open, with a 25% buffer",
      "Closing",
      "Session unknown, counted as closed",
    ]);
  });
});

describe("week", () => {
  test("before the last trading day the close has no hour yet, and the week still runs to the reopen", () => {
    const marks = week({ ...base, regime: 2, current: 4_000n * D, boundaryMs: 0n, nowMs: now });
    expect(marks.map((m) => [m.key, m.hours, m.requirement])).toEqual([
      ["now", undefined, 4_000n * D],
      ["close", undefined, 5_000n * D],
      ["reopen", undefined, 4_000n * D],
      ["regular", undefined, 4_000n * D],
    ]);
  });

  test("the last trading day runs now, the ramp, the close, the 24/5 reopen and the regular open", () => {
    const boundaryMs = BigInt(now + 11 * HOUR);
    const marks = week({ ...base, regime: 2, current: 4_000n * D, boundaryMs, nowMs: now });
    expect(marks.map((m) => [m.key, m.hours, m.requirement])).toEqual([
      ["now", undefined, 4_000n * D],
      ["ramp", 4, 4_000n * D],
      ["close", 11, 5_000n * D],
      ["reopen", undefined, 4_000n * D],
      ["regular", undefined, 4_000n * D],
    ]);
  });

  test("the close needs at least the weekend leverage floor when it exceeds the closed requirement", () => {
    const marks = week({
      ...base,
      gross: 30_000n * D,
      regime: 2,
      current: 4_000n * D,
      boundaryMs: BigInt(now + 9 * HOUR),
      nowMs: now,
    });
    expect(marks.find((m) => m.key === "close")?.requirement).toBe(6_000n * D);
  });

  test("inside the ramp the next mark is the close", () => {
    const marks = week({ ...base, regime: 3, current: 4_400n * D, boundaryMs: BigInt(now + 3 * HOUR), nowMs: now });
    expect(marks.map((m) => [m.key, m.hours])).toEqual([
      ["now", undefined],
      ["close", 3],
      ["reopen", undefined],
      ["regular", undefined],
    ]);
  });

  test("a closed or unknown session runs now, then the 24/5 reopen, back to the buffered open, and the regular open", () => {
    for (const regime of [0, 1]) {
      const marks = week({ ...base, regime, current: 5_000n * D, boundaryMs: 0n, nowMs: now });
      expect(marks.map((m) => [m.key, m.requirement])).toEqual([
        ["now", 5_000n * D],
        ["reopen", 4_000n * D],
        ["regular", 4_000n * D],
      ]);
    }
    const closed = week({ ...base, regime: 1, current: 5_000n * D, boundaryMs: BigInt(now + 19 * HOUR + 1), nowMs: now });
    expect(closed.find((m) => m.key === "reopen")?.hours).toBe(19);
  });
});

describe("capacity", () => {
  const limits = { debt: 400_000_000_000n, debtCap: 1_000_000_000_000n, weekendDebtCap: 500_000_000_000n, weekendLeverage: 50_000 };

  test("is what the equity holds above the requirement, in USDG", () => {
    expect(capacity({ equity: 9_812n * D, requirement: 3_940n * D, gross: 20_000n * D, weekendAhead: false, ...limits })).toBe(
      5_872_000_000n,
    );
  });

  test("before a closure it also holds the gross exposure within 5× and the weekend debt cap", () => {
    expect(capacity({ equity: 9_812n * D, requirement: 3_940n * D, gross: 40_000n * D, weekendAhead: true, ...limits })).toBe(
      1_812_000_000n,
    );
    expect(
      capacity({ equity: 9_812n * D, requirement: 3_940n * D, gross: 0n, weekendAhead: true, ...limits, debt: 499_000_000_000n }),
    ).toBe(1_000_000_000n);
  });

  test("is never below zero, and stops at the debt cap", () => {
    expect(capacity({ equity: 1_000n * D, requirement: 3_940n * D, gross: 0n, weekendAhead: false, ...limits })).toBe(0n);
    expect(
      capacity({ equity: 9_812n * D, requirement: 0n, gross: 0n, weekendAhead: false, ...limits, debt: 999_999_000_000n }),
    ).toBe(1_000_000n);
  });
});

describe("weekendPlan", () => {
  test("says what to add or reduce for the gross exposure to fit 5× the equity", () => {
    expect(weekendPlan({ gross: 60_000n * D, equity: 10_000n * D, weekendLeverage: 50_000 })).toEqual({
      add: 2_000n * D,
      reduce: 10_000n * D,
    });
    expect(weekendPlan({ gross: 40_000n * D, equity: 10_000n * D, weekendLeverage: 50_000 })).toBeUndefined();
  });

  test("a position without equity must add all the floor", () => {
    expect(weekendPlan({ gross: 5_000n * D, equity: -1n * D, weekendLeverage: 50_000 })).toEqual({
      add: 1_001n * D,
      reduce: 5_000n * D,
    });
  });
});

describe("premiumOver", () => {
  test("is what a loan owes over a closure at the premium rate", () => {
    expect(premiumOver(2_000_000_000n, 500, 48)).toBe(547_945n);
    expect(premiumOver(0n, 500, 48)).toBe(0n);
  });
});
