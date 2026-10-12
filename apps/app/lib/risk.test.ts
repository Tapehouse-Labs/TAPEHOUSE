// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { formatters } from "./format";
import { exposureUsed, premiumWords, roundStatus, twapWords } from "./risk";
import { say } from "./say";

const f = formatters("en");

describe("exposureUsed", () => {
  test("is the share of the limit the closure covered", () => {
    expect(say(exposureUsed(f, 25_000_000n, 100_000_000n))).toBe("25.00%");
    expect(say(exposureUsed(f, 1n, 3n))).toBe("33.33%");
  });
  test("reads zero covered as unused and an absent limit as none", () => {
    expect(say(exposureUsed(f, 0n, 100_000_000n))).toBe("0.00%");
    expect(say(exposureUsed(f, 0n, 0n))).toBe("No limit set");
    expect(say(exposureUsed(f, 0n, 0n), "es")).toBe("Sin límite definido");
  });
  test("caps at the limit", () => {
    expect(say(exposureUsed(f, 200n, 100n))).toBe("100.00%");
  });
  test("writes the percentage with the language's separator", () => {
    expect(say(exposureUsed(formatters("pt"), 1n, 3n), "pt")).toBe("33,33%");
  });
});

describe("roundStatus", () => {
  const open = { sealMs: 1n, cleared: false };
  test("is empty where no round was opened", () => {
    expect(roundStatus({ sealMs: 0n, cleared: false }, { openMs: 0n, revealing: false }, 0n)).toBe("none");
  });
  test("follows the phase of the current round", () => {
    expect(roundStatus(open, { openMs: 5n, revealing: false }, 5n)).toBe("committing");
    expect(roundStatus(open, { openMs: 5n, revealing: true }, 5n)).toBe("revealing");
  });
  test("separates a cleared round from one awaiting its clearing", () => {
    expect(roundStatus({ sealMs: 1n, cleared: true }, { openMs: 0n, revealing: false }, 5n)).toBe("cleared");
    expect(roundStatus(open, { openMs: 0n, revealing: false }, 5n)).toBe("awaiting");
  });
});

describe("twapWords", () => {
  test("shows the average once valid and says otherwise", () => {
    expect(say(twapWords(f, true, 23_492_917_988n))).toBe("234.93");
    expect(say(twapWords(f, false, 0n))).toBe("Not yet valid");
  });
});

describe("premiumWords", () => {
  test("signs the premium in percent", () => {
    expect(premiumWords(f, 0n)).toBe("0%");
    expect(premiumWords(f, 25n)).toBe("+0.25%");
    expect(premiumWords(f, -150n)).toBe("−1.5%");
  });
});
