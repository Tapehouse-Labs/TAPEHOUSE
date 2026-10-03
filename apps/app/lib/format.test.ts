// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { bandPercent, price, shortAddress } from "./format";

describe("price", () => {
  test("renders 8-decimal USD with two decimals and thousands separators", () => {
    expect(price(18_250_000_000n)).toBe("182.50");
    expect(price(76_614_000_000n)).toBe("766.14");
    expect(price(123_456_789_000_000n)).toBe("1,234,567.89");
  });

  test("rounds half up on the third decimal, carrying into the units", () => {
    expect(price(18_249_500_000n)).toBe("182.50");
    expect(price(18_249_499_999n)).toBe("182.49");
    expect(price(99_999_500_000n)).toBe("1,000.00");
  });

  test("renders zero and sub-cent prices", () => {
    expect(price(0n)).toBe("0.00");
    expect(price(499_999n)).toBe("0.00");
    expect(price(500_000n)).toBe("0.01");
  });
});

describe("bandPercent", () => {
  test("renders a half-width in basis points as a signed percentage", () => {
    expect(bandPercent(110n)).toBe("±1.10%");
    expect(bandPercent(5n)).toBe("±0.05%");
    expect(bandPercent(10_000n)).toBe("±100.00%");
  });
});

describe("shortAddress", () => {
  test("keeps the first four and last four hex digits", () => {
    expect(shortAddress("0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E")).toBe("0x3f1E…2d0E");
  });
});
