// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { exposureUsed, premiumWords, roundStatus, twapWords } from "./risk";

describe("exposureUsed", () => {
  test("is the share of the limit the closure covered", () => {
    expect(exposureUsed(25_000_000n, 100_000_000n)).toBe("25.00%");
    expect(exposureUsed(1n, 3n)).toBe("33.33%");
  });
  test("reads zero covered as unused and an absent limit as none", () => {
    expect(exposureUsed(0n, 100_000_000n)).toBe("0.00%");
    expect(exposureUsed(0n, 0n)).toBe("No limit set");
  });
  test("caps at the limit", () => {
    expect(exposureUsed(200n, 100n)).toBe("100.00%");
  });
});

describe("roundStatus", () => {
  const open = { sealMs: 1n, cleared: false };
  test("is empty where no round was opened", () => {
    expect(roundStatus({ sealMs: 0n, cleared: false }, { openMs: 0n, revealing: false }, 0n)).toBe("No round");
  });
  test("follows the phase of the current round", () => {
    expect(roundStatus(open, { openMs: 5n, revealing: false }, 5n)).toBe("Committing");
    expect(roundStatus(open, { openMs: 5n, revealing: true }, 5n)).toBe("Revealing");
  });
  test("separates a cleared round from one awaiting its clearing", () => {
    expect(roundStatus({ sealMs: 1n, cleared: true }, { openMs: 0n, revealing: false }, 5n)).toBe("Cleared");
    expect(roundStatus(open, { openMs: 0n, revealing: false }, 5n)).toBe("Awaiting clearing");
  });
});

describe("twapWords", () => {
  test("shows the average once valid and says otherwise", () => {
    expect(twapWords(true, 23_492_917_988n)).toBe("234.93");
    expect(twapWords(false, 0n)).toBe("Not yet valid");
  });
});

describe("premiumWords", () => {
  test("signs the premium in percent", () => {
    expect(premiumWords(0n)).toBe("0%");
    expect(premiumWords(25n)).toBe("+0.25%");
    expect(premiumWords(-150n)).toBe("−1.5%");
  });
});
