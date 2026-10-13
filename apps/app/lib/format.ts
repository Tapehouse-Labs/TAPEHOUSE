// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Locale } from "@/i18n/config";

/**
 * How amounts are written and read in one language. The arithmetic is the same in every language and exact, on bigints
 * rounded half up; only the separators change, as `Intl` gives them for the language. A typed amount is read with the
 * language's decimal separator alone, and no group separator: `1.500` is never read as `1,5`, nor the reverse.
 */
export type Formatters = ReturnType<typeof formatters>;

export function formatters(locale: Locale) {
  const grouped = new Intl.NumberFormat(locale);
  const decimal = grouped.formatToParts(1.5).find((part) => part.type === "decimal")?.value ?? ".";
  const group = (value: bigint) => grouped.format(value);
  const utc = new Intl.DateTimeFormat(locale === "en" ? "en-GB" : locale, {
    timeZone: "UTC",
    day: "numeric",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });

  const day = new Intl.DateTimeFormat(locale === "en" ? "en-GB" : locale, {
    timeZone: "UTC",
    day: "numeric",
    month: "short",
    year: "numeric",
  });
  const list = new Intl.ListFormat(locale, { style: "long", type: "conjunction" });

  /** `units` over `scale` as a decimal with two places, rounded half up, grouped by the language's separator. */
  const fixed = (units: bigint, scale: bigint): string => {
    const cents = (units * 100n + scale / 2n) / scale;
    return `${group(cents / 100n)}${decimal}${(cents % 100n).toString().padStart(2, "0")}`;
  };

  return {
    locale,
    /** The language's decimal separator. */
    decimal,

    /** A band price, USD with 8 decimals, as dollars and cents rounded half up. */
    price: (value: bigint): string => fixed(value, 100_000_000n),

    /** A band's half-width, in basis points, as a signed percentage. */
    bandPercent: (halfBps: bigint): string =>
      `±${halfBps / 100n}${decimal}${(halfBps % 100n).toString().padStart(2, "0")}%`,

    /** A rate or a share with 18 decimals as a percentage with two decimals, rounded half up. */
    percent: (wad: bigint): string => {
      const hundredths = (wad * 10_000n + 10n ** 18n / 2n) / 10n ** 18n;
      return `${hundredths / 100n}${decimal}${(hundredths % 100n).toString().padStart(2, "0")}%`;
    },

    /** Basis points as a percentage, without trailing zeros. */
    bpsPercent: (bps: bigint | number): string => {
      const value = BigInt(bps);
      const fraction = (value % 100n).toString().padStart(2, "0").replace(/0+$/, "");
      return `${value / 100n}${fraction ? `${decimal}${fraction}` : ""}%`;
    },

    /** Leverage in basis points as a multiple with one decimal, rounded half up and without a trailing zero: `2.5`. */
    multiple: (bps: bigint | number): string => {
      const tenths = (BigInt(bps) + 500n) / 1_000n;
      return `${tenths / 10n}${tenths % 10n ? `${decimal}${tenths % 10n}` : ""}`;
    },

    /** A moment as its day, month and time in UTC, without the zone. */
    utc: (date: Date): string => utc.format(date),

    /** A day, as its date in UTC with the year. */
    date: (date: Date): string => day.format(date),

    /** Names joined as the language joins a list: `a, b and c`. */
    list: (items: string[]): string => list.format(items),

    /** An address as its first four and last four hex digits. */
    shortAddress: (address: string): string => `${address.slice(0, 6)}…${address.slice(-4)}`,

    /** A USD value with 18 decimals, as the engine and the accounts report it, in dollars and cents rounded half up. */
    usd: (value: bigint): string =>
      value < 0n ? `−$${fixed(-value, 10n ** 18n)}` : `$${fixed(value, 10n ** 18n)}`,

    /** A USDG amount, 6 decimals, with two decimals rounded half up. */
    usdg: (value: bigint): string => fixed(value, 10n ** 6n),

    /** A Stock Token or WETH amount, 18 decimals, with up to four decimals rounded down. */
    tokens: (value: bigint): string => {
      const tenThousandths = value / 10n ** 14n;
      const fraction = (tenThousandths % 10_000n).toString().padStart(4, "0").replace(/0+$/, "");
      return `${group(tenThousandths / 10_000n)}${fraction ? `${decimal}${fraction}` : ""}`;
    },

    /**
     * What a visitor types as it is kept: where the language's decimal separator is a comma, the point of a numeric
     * keypad is that separator too.
     */
    normalize: (typed: string): string => (decimal === "," ? typed.replaceAll(".", ",") : typed),

    /**
     * A typed amount at `decimals`; undefined unless it is a positive amount. Spaces around it are ignored, and in
     * English so are commas, which only group there.
     */
    parseAmount: (typed: string, decimals: number): bigint | undefined => {
      const trimmed = typed.trim();
      if (decimal !== "." && trimmed.includes(".")) return undefined;
      const plain = trimmed.replaceAll(decimal === "." ? "," : "", "").replace(decimal, ".");
      const match = /^(\d*)(?:\.(\d*))?$/.exec(plain);
      if (!match || (match[1] === "" && !match[2]) || (match[2]?.length ?? 0) > decimals) return undefined;
      const value = BigInt(`${match[1] || "0"}${(match[2] ?? "").padEnd(decimals, "0")}`);
      return value > 0n ? value : undefined;
    },

    /** An amount at `decimals` as it would be typed: every significant decimal, no separators. */
    unparse: (value: bigint, decimals: number): string => {
      const unit = 10n ** BigInt(decimals);
      const fraction = (value % unit).toString().padStart(decimals, "0").replace(/0+$/, "");
      return `${value / unit}${fraction ? `${decimal}${fraction}` : ""}`;
    },
  };
}
