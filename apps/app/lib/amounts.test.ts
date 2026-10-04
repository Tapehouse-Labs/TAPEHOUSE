// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { parseAmount, tokens, unparse, usd, usdg } from "./amounts";

describe("usd", () => {
  test("renders 18-decimal USD as dollars and cents, rounded half up", () => {
    expect(usd(9_812_404_999_999_999_999_999n)).toBe("$9,812.40");
    expect(usd(9_812_405_000_000_000_000_000n)).toBe("$9,812.41");
    expect(usd(0n)).toBe("$0.00");
  });

  test("renders a negative equity with a minus sign", () => {
    expect(usd(-1_500_000_000_000_000_000n)).toBe("−$1.50");
  });
});

describe("usdg", () => {
  test("renders 6-decimal USDG with two decimals, rounded half up", () => {
    expect(usdg(2_000_000_000n)).toBe("2,000.00");
    expect(usdg(554_795n)).toBe("0.55");
    expect(usdg(5_000n)).toBe("0.01");
    expect(usdg(4_999n)).toBe("0.00");
  });
});

describe("tokens", () => {
  test("renders 18-decimal token amounts with up to four decimals, rounded down, zeros trimmed", () => {
    expect(tokens(12n * 10n ** 18n)).toBe("12");
    expect(tokens(1_234_567_000_000_000_000n)).toBe("1.2345");
    expect(tokens(10n ** 14n - 1n)).toBe("0");
    expect(tokens(1_500_000_000_000_000_000_000n)).toBe("1,500");
  });
});

describe("parseAmount", () => {
  test("reads a typed amount at the token's decimals", () => {
    expect(parseAmount("2000", 6)).toBe(2_000_000_000n);
    expect(parseAmount("1,250.5", 6)).toBe(1_250_500_000n);
    expect(parseAmount(" 0.000001 ", 6)).toBe(1n);
    expect(parseAmount(".5", 18)).toBe(5n * 10n ** 17n);
  });

  test("refuses what is not a positive amount at those decimals", () => {
    for (const typed of ["", "0", "0.00", "-1", "1e3", "abc", "1.2.3", "0.0000001"]) {
      expect(parseAmount(typed, 6)).toBeUndefined();
    }
  });
});

describe("unparse", () => {
  test("writes an amount back as it would be typed", () => {
    expect(unparse(1_204_100_000n, 6)).toBe("1204.1");
    expect(unparse(2n * 10n ** 18n, 18)).toBe("2");
    expect(parseAmount(unparse(123_456_789n, 6), 6)).toBe(123_456_789n);
  });
});
