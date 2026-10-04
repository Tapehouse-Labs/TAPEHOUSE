// SPDX-License-Identifier: MIT OR Apache-2.0

const grouped = new Intl.NumberFormat("en-US");

function fixed(units: bigint, scale: bigint): string {
  const cents = (units * 100n + scale / 2n) / scale;
  return `${grouped.format(cents / 100n)}.${(cents % 100n).toString().padStart(2, "0")}`;
}

/** A USD value with 18 decimals, as the engine and the accounts report it, in dollars and cents rounded half up. */
export function usd(value: bigint): string {
  return value < 0n ? `−$${fixed(-value, 10n ** 18n)}` : `$${fixed(value, 10n ** 18n)}`;
}

/** A USDG amount, 6 decimals, with two decimals rounded half up. */
export function usdg(value: bigint): string {
  return fixed(value, 10n ** 6n);
}

/** A Stock Token or WETH amount, 18 decimals, with up to four decimals rounded down. */
export function tokens(value: bigint): string {
  const tenThousandths = value / 10n ** 14n;
  const fraction = (tenThousandths % 10_000n).toString().padStart(4, "0").replace(/0+$/, "");
  return `${grouped.format(tenThousandths / 10_000n)}${fraction ? `.${fraction}` : ""}`;
}

/** A typed amount at `decimals`, thousands separators allowed; undefined unless it is a positive amount. */
export function parseAmount(typed: string, decimals: number): bigint | undefined {
  const match = /^(\d*)(?:\.(\d*))?$/.exec(typed.trim().replaceAll(",", ""));
  if (!match || (match[1] === "" && !match[2]) || (match[2]?.length ?? 0) > decimals) return undefined;
  const value = BigInt(`${match[1] || "0"}${(match[2] ?? "").padEnd(decimals, "0")}`);
  return value > 0n ? value : undefined;
}

/** An amount at `decimals` as it would be typed: every significant decimal, no separators. */
export function unparse(value: bigint, decimals: number): string {
  const unit = 10n ** BigInt(decimals);
  const fraction = (value % unit).toString().padStart(decimals, "0").replace(/0+$/, "");
  return `${value / unit}${fraction ? `.${fraction}` : ""}`;
}
