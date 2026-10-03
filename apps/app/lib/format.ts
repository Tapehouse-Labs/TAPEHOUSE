// SPDX-License-Identifier: MIT OR Apache-2.0

const grouped = new Intl.NumberFormat("en-US");

/** A band price, USD with 8 decimals, as dollars and cents rounded half up. */
export function price(value: bigint): string {
  const cents = (value + 500_000n) / 1_000_000n;
  return `${grouped.format(cents / 100n)}.${(cents % 100n).toString().padStart(2, "0")}`;
}

/** A band's half-width, in basis points, as a signed percentage. */
export function bandPercent(halfBps: bigint): string {
  return `±${halfBps / 100n}.${(halfBps % 100n).toString().padStart(2, "0")}%`;
}

/** An address as its first four and last four hex digits. */
export function shortAddress(address: string): string {
  return `${address.slice(0, 6)}…${address.slice(-4)}`;
}
