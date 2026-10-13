// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import path from "node:path";

const EXPLORER: Record<string, string> = {
  "46630": "https://explorer.testnet.chain.robinhood.com/address/",
};

type Row = { group: string; name: string; address: string };

function rows(node: unknown, group: string[] = [], out: Row[] = []): Row[] {
  if (node && typeof node === "object") {
    for (const [k, v] of Object.entries(node)) {
      if (typeof v === "string" && /^0x[0-9a-fA-F]{40}$/.test(v)) out.push({ group: group.join("."), name: k, address: v });
      else if (v && typeof v === "object") rows(v, [...group, k], out);
    }
  }
  return out;
}

/** A chain's registry as a table, read from `deployments/<chain>.json` when the page is built. */
export function Deployments({ chain }: { chain: string }) {
  const file = path.join(process.cwd(), "..", "..", "deployments", `${chain}.json`);
  const registry = JSON.parse(readFileSync(file, "utf8")) as Record<string, unknown>;
  const all = rows(registry);
  const groups = [...new Set(all.map((r) => r.group))];
  const explorer = EXPLORER[chain];
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Entry</th>
            <th>Address</th>
          </tr>
        </thead>
        <tbody>
          {groups.flatMap((g) =>
            all
              .filter((r) => r.group === g)
              .map((r) => (
                <tr key={`${g}.${r.name}`}>
                  <td>
                    <code>{`.${g}.${r.name}`}</code>
                  </td>
                  <td>
                    {explorer ? (
                      <a href={`${explorer}${r.address}`} rel="noreferrer" className="addr">
                        <code>{r.address}</code>
                      </a>
                    ) : (
                      <code>{r.address}</code>
                    )}
                  </td>
                </tr>
              )),
          )}
        </tbody>
      </table>
    </div>
  );
}
