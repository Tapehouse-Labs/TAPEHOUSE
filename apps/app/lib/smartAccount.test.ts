// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { freeActions, notches } from "./smartAccount";

describe("notches", () => {
  test("one notch per free action, lit while it is left", () => {
    expect(notches(3n)).toEqual([true, true, true]);
    expect(notches(2n)).toEqual([true, true, false]);
    expect(notches(0n)).toEqual([false, false, false]);
  });
});

describe("freeActions", () => {
  test("before the account exists, its creation and the first three actions are free", () => {
    expect(freeActions({ deployed: false, left: 3n })).toBe("Creating it costs you nothing, and so do your first 3 actions.");
  });

  test("a live account counts what is left", () => {
    expect(freeActions({ deployed: true, left: 2n })).toBe("2 free actions left");
    expect(freeActions({ deployed: true, left: 1n })).toBe("1 free action left");
    expect(freeActions({ deployed: true, left: 0n })).toBe("Free actions used. You pay gas from here.");
  });
});
