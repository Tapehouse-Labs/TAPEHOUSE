// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { bandState, seal } from "./session";

const HOUR = 3_600_000n;
const now = 1_790_000_000_000;

describe("seal", () => {
  test("a closed session counts whole hours to its reopening", () => {
    expect(seal({ state: 1, boundaryMs: BigInt(now) + 41n * HOUR + 59n * 60_000n }, now)).toEqual({
      state: "CLOSED",
      hours: 41,
      label: "Market closed, reopens in 41 hours",
    });
  });

  test("an open session counts whole hours to its close", () => {
    expect(seal({ state: 2, boundaryMs: BigInt(now) + 3n * HOUR }, now)).toEqual({
      state: "OPEN",
      hours: 3,
      label: "Market open, closes in 3 hours",
    });
  });

  test("a boundary under an hour away reads as zero hours, and one already past as zero", () => {
    expect(seal({ state: 2, boundaryMs: BigInt(now) + 59n * 60_000n }, now).hours).toBe(0);
    expect(seal({ state: 1, boundaryMs: BigInt(now) - HOUR }, now).hours).toBe(0);
  });

  test("an unknown session has no countdown", () => {
    expect(seal({ state: 0, boundaryMs: 0n }, now)).toEqual({
      state: "UNKNOWN",
      hours: undefined,
      label: "Market session not known, every band degraded",
    });
  });
});

describe("bandState", () => {
  test("names each state the band reports", () => {
    expect([0, 1, 2, 3].map(bandState)).toEqual(["HALTED", "DEGRADED", "CLOSED", "OPEN"]);
  });

  test("treats a state it does not know as halted", () => {
    expect(bandState(9)).toBe("HALTED");
  });
});
