// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { wrapPlan } from "./wrapped";

const D = 10n ** 18n;

describe("wrapPlan", () => {
  test("a wallet holding enough of the wrapper wraps nothing", () => {
    expect(wrapPlan({ held: 5n * D, allowance: 0n, amount: 5n * D })).toBeUndefined();
    expect(wrapPlan({ held: 6n * D, allowance: 0n, amount: 5n * D })).toBeUndefined();
  });

  test("wraps what the wrapper is short of, and approves when the allowance does not cover it", () => {
    expect(wrapPlan({ held: 2n * D, allowance: 0n, amount: 5n * D })).toEqual({ amount: 3n * D, approve: true });
    expect(wrapPlan({ held: 2n * D, allowance: 2n * D, amount: 5n * D })).toEqual({ amount: 3n * D, approve: true });
    expect(wrapPlan({ held: 2n * D, allowance: 3n * D, amount: 5n * D })).toEqual({ amount: 3n * D, approve: false });
  });

  test("a wallet without any wrapper wraps all of it", () => {
    expect(wrapPlan({ held: 0n, allowance: 9n * D, amount: 4n * D })).toEqual({ amount: 4n * D, approve: false });
  });
});
