// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { LOCALES } from "@/i18n/config";
import { formatters } from "./format";

const en = formatters("en");
const es = formatters("es");
const pt = formatters("pt");

describe("price", () => {
  test("renders 8-decimal USD with two decimals and the language's separators", () => {
    expect(en.price(18_250_000_000n)).toBe("182.50");
    expect(en.price(123_456_789_000_000n)).toBe("1,234,567.89");
    expect(es.price(123_456_789_000_000n)).toBe("1.234.567,89");
    expect(pt.price(123_456_789_000_000n)).toBe("1.234.567,89");
    expect(es.price(76_614_000_000n)).toBe("766,14");
  });

  test("rounds half up on the third decimal, carrying into the units", () => {
    expect(en.price(18_249_500_000n)).toBe("182.50");
    expect(en.price(18_249_499_999n)).toBe("182.49");
    expect(en.price(99_999_500_000n)).toBe("1,000.00");
    expect(pt.price(99_999_500_000n)).toBe("1.000,00");
  });

  test("renders zero and sub-cent prices", () => {
    expect(en.price(0n)).toBe("0.00");
    expect(en.price(499_999n)).toBe("0.00");
    expect(en.price(500_000n)).toBe("0.01");
  });
});

describe("bandPercent", () => {
  test("renders a half-width in basis points as a signed percentage", () => {
    expect(en.bandPercent(110n)).toBe("±1.10%");
    expect(en.bandPercent(5n)).toBe("±0.05%");
    expect(en.bandPercent(10_000n)).toBe("±100.00%");
    expect(es.bandPercent(30n)).toBe("±0,30%");
  });
});

describe("shortAddress", () => {
  test("keeps the first four and last four hex digits", () => {
    expect(en.shortAddress("0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E")).toBe("0x3f1E…2d0E");
  });
});

describe("usd", () => {
  test("renders 18-decimal USD as dollars and cents, rounded half up", () => {
    expect(en.usd(9_812_404_999_999_999_999_999n)).toBe("$9,812.40");
    expect(en.usd(9_812_405_000_000_000_000_000n)).toBe("$9,812.41");
    expect(en.usd(0n)).toBe("$0.00");
    expect(es.usd(9_812_405_000_000_000_000_000n)).toBe("$9812,41");
    expect(es.usd(98_124_050_000_000_000_000_000n)).toBe("$98.124,05");
    expect(pt.usd(1_500_000_000_000_000_000n)).toBe("$1,50");
  });

  test("renders a negative equity with a minus sign", () => {
    expect(en.usd(-1_500_000_000_000_000_000n)).toBe("−$1.50");
    expect(es.usd(-1_500_000_000_000_000_000n)).toBe("−$1,50");
  });
});

describe("usdg", () => {
  test("renders 6-decimal USDG with two decimals, rounded half up", () => {
    expect(en.usdg(2_000_000_000n)).toBe("2,000.00");
    expect(en.usdg(554_795n)).toBe("0.55");
    expect(en.usdg(5_000n)).toBe("0.01");
    expect(en.usdg(4_999n)).toBe("0.00");
    expect(es.usdg(2_000_000_000n)).toBe("2000,00");
    expect(es.usdg(12_000_000_000n)).toBe("12.000,00");
  });
});

describe("tokens", () => {
  test("renders 18-decimal token amounts with up to four decimals, rounded down, zeros trimmed", () => {
    expect(en.tokens(12n * 10n ** 18n)).toBe("12");
    expect(en.tokens(1_234_567_000_000_000_000n)).toBe("1.2345");
    expect(en.tokens(10n ** 14n - 1n)).toBe("0");
    expect(en.tokens(1_500_000_000_000_000_000_000n)).toBe("1,500");
    expect(es.tokens(1_234_567_000_000_000_000n)).toBe("1,2345");
    expect(pt.tokens(1_500_000_000_000_000_000_000n)).toBe("1.500");
  });
});

describe("parseAmount", () => {
  test("reads a typed amount at the token's decimals", () => {
    expect(en.parseAmount("2000", 6)).toBe(2_000_000_000n);
    expect(en.parseAmount("1,250.5", 6)).toBe(1_250_500_000n);
    expect(en.parseAmount(" 0.000001 ", 6)).toBe(1n);
    expect(en.parseAmount(".5", 18)).toBe(5n * 10n ** 17n);
    expect(es.parseAmount("1250,5", 6)).toBe(1_250_500_000n);
    expect(pt.parseAmount(",5", 18)).toBe(5n * 10n ** 17n);
  });

  test("refuses what is not a positive amount at those decimals", () => {
    for (const typed of ["", "0", "0.00", "-1", "1e3", "abc", "1.2.3", "0.0000001"]) {
      expect(en.parseAmount(typed, 6)).toBeUndefined();
    }
    for (const typed of ["", "0", "0,00", "-1", "1,2,3", "0,0000001"]) {
      expect(es.parseAmount(typed, 6)).toBeUndefined();
    }
  });

  test("in a comma language a point is never a group separator, so it is never read as one", () => {
    expect(es.parseAmount("1.500", 6)).toBeUndefined();
    expect(es.parseAmount("1.500,25", 6)).toBeUndefined();
    expect(pt.parseAmount("1.5", 6)).toBeUndefined();
  });
});

describe("normalize", () => {
  test("turns the point of a numeric keypad into the comma of a comma language, and leaves English alone", () => {
    expect(es.normalize("1.5")).toBe("1,5");
    expect(pt.normalize("0.25")).toBe("0,25");
    expect(en.normalize("1.5")).toBe("1.5");
    expect(es.parseAmount(es.normalize("1.5"), 6)).toBe(1_500_000n);
  });
});

describe("unparse", () => {
  test("writes an amount back as it would be typed, in each language", () => {
    expect(en.unparse(1_204_100_000n, 6)).toBe("1204.1");
    expect(es.unparse(1_204_100_000n, 6)).toBe("1204,1");
    expect(en.unparse(2n * 10n ** 18n, 18)).toBe("2");
    for (const locale of LOCALES) {
      const f = formatters(locale);
      expect(f.parseAmount(f.unparse(123_456_789n, 6), 6)).toBe(123_456_789n);
    }
  });
});
