// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { isSanctioned, notices, restrictedLocation, SANCTIONS } from "./compliance";
import { msg } from "./messages";
import sanctions from "./sanctions.json";

const alice = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
const smart = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC";

describe("restrictedLocation", () => {
  test("the embargoed countries and regions are restricted, whatever the case of the code", () => {
    for (const country of ["CU", "IR", "KP", "SY", "ir"]) expect(restrictedLocation(country, null)).toBe(true);
    for (const region of ["43", "40", "14", "09"]) expect(restrictedLocation("UA", region)).toBe(true);
  });

  test("the rest of Ukraine, other countries and an unknown location are not", () => {
    expect(restrictedLocation("UA", "30")).toBe(false);
    expect(restrictedLocation("UA", null)).toBe(false);
    expect(restrictedLocation("CO", "43")).toBe(false);
    expect(restrictedLocation("US", "CA")).toBe(false);
    expect(restrictedLocation(null, null)).toBe(false);
    expect(restrictedLocation(undefined, undefined)).toBe(false);
  });
});

describe("isSanctioned", () => {
  test("every address of the list is sanctioned, in any case, and others are not", () => {
    expect(sanctions.addresses.length).toBeGreaterThan(50);
    for (const address of sanctions.addresses) {
      expect(address).toMatch(/^0x[0-9a-f]{40}$/);
      expect(isSanctioned(address)).toBe(true);
      expect(isSanctioned(address.toUpperCase().replace("0X", "0x"))).toBe(true);
    }
    expect(isSanctioned(alice)).toBe(false);
    expect(isSanctioned(undefined)).toBe(false);
    expect(isSanctioned(null)).toBe(false);
  });

  test("says where the list came from", () => {
    expect(SANCTIONS.source).toBe("https://www.treasury.gov/ofac/downloads/sdn.csv");
    expect(SANCTIONS.retrievedAt).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });
});

describe("notices", () => {
  const blocks: { asset: string; blocked: Record<string, boolean> }[] = [
    { asset: "TSLA", blocked: { [smart]: true, [alice]: false } },
    { asset: "AMZN", blocked: { [alice]: false } },
  ];

  test("an issuer's block on the owner or the account is named by asset", () => {
    expect(notices([alice, smart], blocks)).toEqual([
      { message: msg("compliance.blocked", { asset: "TSLA" }), tone: "warning" },
    ]);
    expect(notices([alice, smart.toLowerCase()], blocks)).toHaveLength(1);
    expect(notices([alice], blocks)).toEqual([]);
  });

  test("a sanctioned owner or account is told first", () => {
    const listed = sanctions.addresses[0]!;
    const out = notices([alice, listed], blocks);
    expect(out[0]).toMatchObject({ tone: "loss" });
    expect(notices([listed], [])).toHaveLength(1);
    expect(notices([undefined, undefined], [])).toEqual([]);
  });
});
