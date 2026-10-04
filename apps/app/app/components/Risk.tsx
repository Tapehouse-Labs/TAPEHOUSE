// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useChains } from "wagmi";
import { usdg } from "@/lib/amounts";
import { exposureName } from "@/lib/earn";
import { explain } from "@/lib/errors";
import { price } from "@/lib/format";
import { exposureUsed, premiumWords, roundStatus, twapWords } from "@/lib/risk";
import { useRiskState } from "@/lib/riskHook";
import type { RiskState } from "@/lib/riskRead";
import { Bands } from "./Bands";

/** The risk surface: each band, the backstop's exposure and the latest reopening rounds, read at one block. */
export function Risk() {
  const [chain] = useChains();
  const state = useRiskState();
  const data = state.data;
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">Risk · bands, backstop and reopening</h1>
          <p className="mt-6 max-w-[40ch] text-[19px] leading-snug text-body">
            What stands between a weekend gap and a loss on {chain.name}, every figure read from the chain.
          </p>
          {data && <p className="label mt-3 num">Read at block {data.blockNumber.toLocaleString("en-US")}</p>}
        </div>
        <Bands />
        {!data ? (
          <div aria-busy="true" className="mt-10">
            {state.error ? (
              <p role="alert" className="text-[14px] text-[var(--error)]">
                {explain(state.error)}
              </p>
            ) : (
              <span className="skeleton block h-[18px] w-[18rem]" />
            )}
          </div>
        ) : (
          <>
            <Reference bands={data.bands} />
            <Exposure state={data.backstop} />
            <Reopening auction={data.auction} />
          </>
        )}
      </div>
    </main>
  );
}

function Reference({ bands }: { bands: RiskState["bands"] }) {
  return (
    <section aria-labelledby="reference-title" className="mt-10">
      <h2 id="reference-title" className="label">
        Reference price
      </h2>
      <table className="ledger num">
        <caption className="sr-only">Each asset&apos;s TWAP and the premium of its token&apos;s market over the band</caption>
        <thead>
          <tr>
            <th scope="col">Asset</th>
            <th scope="col" className="r">
              TWAP
            </th>
            <th scope="col" className="r">
              Premium
            </th>
          </tr>
        </thead>
        <tbody>
          {bands.map(({ asset, band }) => (
            <tr key={asset}>
              <th scope="row">{asset}</th>
              <td className="r">{twapWords(band.twapValid, band.twap)}</td>
              <td className="r">{band.mid > 0n ? premiumWords(band.premiumBps) : "—"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

function Exposure({ state }: { state: RiskState["backstop"] }) {
  return (
    <section aria-labelledby="exposure-title" className="mt-10">
      <h2 id="exposure-title" className="label">
        Backstop exposure
      </h2>
      {!state ? (
        <p className="mt-4 text-[15px] text-body">The gap backstop is not deployed on this chain yet.</p>
      ) : (
        <>
          <table className="ledger num">
            <caption className="sr-only">The gap backstop&apos;s exposure limit, what the current closure covered and what is left</caption>
            <thead>
              <tr>
                <th scope="col">Position</th>
                <th scope="col" className="r">
                  Limit
                </th>
                <th scope="col" className="r">
                  Covered
                </th>
                <th scope="col" className="r">
                  Used
                </th>
                <th scope="col" className="r">
                  Left
                </th>
              </tr>
            </thead>
            <tbody>
              {state.exposures.map((e) => (
                <tr key={e.position}>
                  <th scope="row">{exposureName(e.position)}</th>
                  <td className="r">{usdg(e.limit)}</td>
                  <td className="r">{usdg(e.covered)}</td>
                  <td className="r">{exposureUsed(e.covered, e.limit)}</td>
                  <td className="r">{usdg(e.left)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p className="mt-2 text-[13px] text-muted num">
            The backstop holds {usdg(state.held)} USDG; what is left is never more than that.
          </p>
        </>
      )}
    </section>
  );
}

function Reopening({ auction }: { auction: RiskState["auction"] }) {
  const rounds = auction?.rounds.filter((r) => r.round.sealMs > 0n) ?? [];
  return (
    <section aria-labelledby="reopening-title" className="mt-10 pb-16">
      <h2 id="reopening-title" className="label">
        Reopening auction
      </h2>
      {!auction ? (
        <p className="mt-4 text-[15px] text-body">The reopening auction is not deployed on this chain yet.</p>
      ) : rounds.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">No reopening round yet.</p>
      ) : (
        <table className="ledger num">
          <caption className="sr-only">The latest reopening round of each asset</caption>
          <thead>
            <tr>
              <th scope="col">Asset</th>
              <th scope="col">Status</th>
              <th scope="col" className="r">
                Floor
              </th>
              <th scope="col" className="r">
                Lots
              </th>
              <th scope="col" className="r">
                Clearing
              </th>
            </tr>
          </thead>
          <tbody>
            {rounds.map(({ asset, openMs, round }) => (
              <tr key={asset}>
                <th scope="row">{asset}</th>
                <td>{roundStatus(round, auction.phase, openMs)}</td>
                <td className="r">{price(round.floor)}</td>
                <td className="r">{round.lots}</td>
                <td className="r">{round.cleared ? price(round.price) : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}
