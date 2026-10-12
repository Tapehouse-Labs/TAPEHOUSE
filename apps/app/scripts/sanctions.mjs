// SPDX-License-Identifier: MIT OR Apache-2.0
// Usage: node scripts/sanctions.mjs
// Writes lib/sanctions.json: every EVM address in the digital currency addresses of OFAC's Specially Designated
// Nationals list, from the Treasury's own file, in lower case. Run it, and commit the result, when the list changes.
import { writeFile } from "node:fs/promises";

const SOURCE = "https://www.treasury.gov/ofac/downloads/sdn.csv";
const response = await fetch(SOURCE);
if (!response.ok) throw new Error(`${SOURCE} answered ${response.status}`);
const text = await response.text();

const addresses = new Set();
for (const line of text.split("\n")) {
  if (!line.includes("Digital Currency Address")) continue;
  for (const match of line.matchAll(/Digital Currency Address - [A-Z0-9]+ (0x[0-9a-fA-F]{40})/g)) {
    addresses.add(match[1].toLowerCase());
  }
}
if (addresses.size === 0) throw new Error("the list holds no EVM address; its format changed");

const file = new URL("../lib/sanctions.json", import.meta.url);
await writeFile(
  file,
  `${JSON.stringify({ source: SOURCE, retrievedAt: new Date().toISOString().slice(0, 10), addresses: [...addresses].sort() }, null, 2)}\n`,
);
console.log(`${addresses.size} addresses`);
