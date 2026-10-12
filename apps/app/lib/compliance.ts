// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Msg, msg } from "./messages";
import sanctions from "./sanctions.json";

/** The countries under comprehensive US sanctions programs, by ISO 3166-1 alpha-2 code: Cuba, Iran, North Korea, Syria. */
export const RESTRICTED_COUNTRIES: readonly string[] = ["CU", "IR", "KP", "SY"];

/** The regions of Ukraine under comprehensive sanctions, by the subdivision code Vercel reports: Crimea, Sevastopol, Donetsk, Luhansk. */
export const RESTRICTED_REGIONS: Readonly<Record<string, readonly string[]>> = { UA: ["43", "40", "14", "09"] };

/** Whether a visitor's country and region, as the platform's geolocation reports them, are restricted. */
export function restrictedLocation(country: string | null | undefined, region: string | null | undefined): boolean {
  const code = country?.toUpperCase();
  if (!code) return false;
  return RESTRICTED_COUNTRIES.includes(code) || (RESTRICTED_REGIONS[code]?.includes(region ?? "") ?? false);
}

const LISTED: ReadonlySet<string> = new Set(sanctions.addresses);

/** Where the list of sanctioned addresses comes from, and when it was taken. */
export const SANCTIONS = { source: sanctions.source, retrievedAt: sanctions.retrievedAt } as const;

/** Whether `address` is among the digital currency addresses of OFAC's Specially Designated Nationals list. */
export function isSanctioned(address: string | null | undefined): boolean {
  return address != null && LISTED.has(address.toLowerCase());
}

export type Notice = { message: Msg; tone: "warning" | "loss" };

/**
 * What the visitor must be told on connecting: that an address is sanctioned, and which Stock Tokens the issuer has
 * blocked it from moving. `blocks` maps each asset to the addresses its issuer has blocked.
 */
export function notices(
  holders: readonly (string | undefined)[],
  blocks: readonly { asset: string; blocked: Readonly<Record<string, boolean>> }[],
): Notice[] {
  const out: Notice[] = [];
  if (holders.some(isSanctioned)) {
    out.push({ message: msg("compliance.sanctioned"), tone: "loss" });
  }
  const mine = new Set(holders.filter((h): h is string => h !== undefined).map((h) => h.toLowerCase()));
  for (const { asset, blocked } of blocks) {
    if (Object.entries(blocked).some(([address, is]) => is && mine.has(address.toLowerCase()))) {
      out.push({ message: msg("compliance.blocked", { asset }), tone: "warning" });
    }
  }
  return out;
}
