// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useChains } from "wagmi";
import { explain } from "@/lib/errors";
import { bandPercent, price } from "@/lib/format";
import { useBands } from "@/lib/hooks";
import { bandState } from "@/lib/session";
import { Seal } from "./Seal";

/** Half-width, in basis points, that fills a row's track: wide enough for a degraded weekend band. */
const TRACK_BPS = 400n;

/** Every asset's band, live from the chain, each row hanging off the stem. */
export function Bands() {
  const bands = useBands();
  const [chain] = useChains();
  return (
    <section aria-labelledby="bands-title" className="mt-16">
      <div className="flex items-center justify-between gap-4">
        <h2 id="bands-title" className="label">
          Live bands<span className="hidden sm:inline"> · {chain.name}</span>
        </h2>
        <span className="sm:hidden">
          <Seal />
        </span>
      </div>
      <table className="mt-3 w-full border-collapse text-[14px]">
        <caption className="sr-only">Each asset&apos;s band: its state, live sources, centre and half-width</caption>
        <thead className="sr-only">
          <tr>
            <th scope="col">Asset</th>
            <th scope="col">Band</th>
            <th scope="col">State</th>
            <th scope="col">Centre</th>
          </tr>
        </thead>
        <tbody>
          {bands.map(({ asset, query }) => (
            <tr key={asset} className="row h-[60px] border-t border-rule">
              <th scope="row" className="w-[72px] text-left font-semibold text-strong">
                {asset}
              </th>
              <Row query={query} />
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

function Row({ query }: { query: ReturnType<typeof useBands>[number]["query"] }) {
  if (query.isPending) {
    return (
      <>
        <td className="px-3">
          <span className="skeleton block h-[18px]" />
        </td>
        <td className="hidden w-[118px] sm:table-cell" />
        <td className="w-[112px] text-right">
          <span className="skeleton ml-auto block h-4 w-16" />
        </td>
      </>
    );
  }
  if (query.isError) {
    return (
      <td colSpan={3} className="px-3 text-right text-[13px] text-muted">
        {explain(query.error)}
      </td>
    );
  }
  const q = query.data;
  const state = bandState(q.state);
  const priced = q.mid > 0n;
  const half = q.halfBps > TRACK_BPS ? TRACK_BPS : q.halfBps;
  const inset = `${(50 - Number((half * 50n) / TRACK_BPS)).toFixed(2)}%`;
  return (
    <>
      <td className="px-3">
        {priced ? (
          <span className="band-track block" aria-hidden="true">
            <span className="band-span" style={{ left: inset, right: inset }} />
            <span className="band-mid" />
          </span>
        ) : (
          <span className="block h-px bg-rule" aria-hidden="true" />
        )}
      </td>
      <td className="hidden w-[118px] text-right sm:table-cell">
        <span className={`label ${state === "HALTED" || state === "DEGRADED" ? "!text-warning" : ""}`}>
          {state}
          <span className="text-muted"> · {q.live} live</span>
        </span>
      </td>
      <td className="w-[112px] text-right">
        {priced ? (
          <span className="num block font-semibold text-strong">{price(q.mid)}</span>
        ) : (
          <span className="hidden font-semibold text-strong sm:block">—</span>
        )}
        {priced ? (
          <span className="num block text-[12px] text-muted">{bandPercent(q.halfBps)}</span>
        ) : (
          <span className={`label block sm:hidden ${state === "HALTED" || state === "DEGRADED" ? "!text-warning" : ""}`}>{state}</span>
        )}
      </td>
    </>
  );
}
