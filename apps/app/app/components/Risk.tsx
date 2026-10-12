// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useTranslations } from "next-intl";
import { useChains } from "wagmi";
import { exposureName } from "@/lib/earn";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import { exposureUsed, premiumWords, roundStatus, twapWords } from "@/lib/risk";
import { useRecord, useRiskState } from "@/lib/riskHook";
import type { RiskState } from "@/lib/riskRead";
import { Bands } from "./Bands";

/** The risk surface: each band, the backstop's exposure and the latest reopening rounds, read at one block. */
export function Risk() {
  const t = useTranslations("risk.page");
  const f = useFormat();
  const explain = useExplain();
  const [chain] = useChains();
  const state = useRiskState();
  const data = state.data;
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">{t("title")}</h1>
          <p className="mt-6 max-w-[40ch] text-[19px] leading-snug text-body">
            {t("intro", { chain: chain.name })}
          </p>
          {data && <p className="label mt-3 num">{t("block", { block: new Intl.NumberFormat(f.locale).format(data.blockNumber) })}</p>}
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
            <Record />
          </>
        )}
      </div>
    </main>
  );
}

function Reference({ bands }: { bands: RiskState["bands"] }) {
  const t = useTranslations("risk.reference");
  const f = useFormat();
  const say = useMsg();
  return (
    <section aria-labelledby="reference-title" className="mt-10">
      <h2 id="reference-title" className="label">
        {t("title")}
      </h2>
      <table className="ledger num">
        <caption className="sr-only">{t("caption")}</caption>
        <thead>
          <tr>
            <th scope="col">{t("asset")}</th>
            <th scope="col" className="r">
              {t("twap")}
            </th>
            <th scope="col" className="r">
              {t("premium")}
            </th>
          </tr>
        </thead>
        <tbody>
          {bands.map(({ asset, band }) => (
            <tr key={asset}>
              <th scope="row">{asset}</th>
              <td className="r">{say(twapWords(f, band.twapValid, band.twap))}</td>
              <td className="r">{band.mid > 0n ? premiumWords(f, band.premiumBps) : "—"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

function Exposure({ state }: { state: RiskState["backstop"] }) {
  const t = useTranslations("risk.exposure");
  const f = useFormat();
  const say = useMsg();
  return (
    <section aria-labelledby="exposure-title" className="mt-10">
      <h2 id="exposure-title" className="label">
        {t("title")}
      </h2>
      {!state ? (
        <p className="mt-4 text-[15px] text-body">{t("absent")}</p>
      ) : (
        <>
          <table className="ledger num">
            <caption className="sr-only">{t("caption")}</caption>
            <thead>
              <tr>
                <th scope="col">{t("position")}</th>
                <th scope="col" className="r">
                  {t("limit")}
                </th>
                <th scope="col" className="r">
                  {t("covered")}
                </th>
                <th scope="col" className="r">
                  {t("used")}
                </th>
                <th scope="col" className="r">
                  {t("left")}
                </th>
              </tr>
            </thead>
            <tbody>
              {state.exposures.map((e) => (
                <tr key={e.position}>
                  <th scope="row">{say(exposureName(e.position))}</th>
                  <td className="r">{f.usdg(e.limit)}</td>
                  <td className="r">{f.usdg(e.covered)}</td>
                  <td className="r">{say(exposureUsed(f, e.covered, e.limit))}</td>
                  <td className="r">{f.usdg(e.left)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p className="mt-2 text-[13px] text-muted num">
            {t("held", { amount: f.usdg(state.held) })}
          </p>
        </>
      )}
    </section>
  );
}

function Reopening({ auction }: { auction: RiskState["auction"] }) {
  const t = useTranslations("risk.reopening");
  const r = useTranslations("risk.round");
  const f = useFormat();
  const rounds = auction?.rounds.filter((r) => r.round.sealMs > 0n) ?? [];
  return (
    <section aria-labelledby="reopening-title" className="mt-10">
      <h2 id="reopening-title" className="label">
        {t("title")}
      </h2>
      {!auction ? (
        <p className="mt-4 text-[15px] text-body">{t("absent")}</p>
      ) : rounds.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">{t("none")}</p>
      ) : (
        <table className="ledger num">
          <caption className="sr-only">{t("caption")}</caption>
          <thead>
            <tr>
              <th scope="col">{t("asset")}</th>
              <th scope="col">{t("status")}</th>
              <th scope="col" className="r">
                {t("floor")}
              </th>
              <th scope="col" className="r">
                {t("lots")}
              </th>
              <th scope="col" className="r">
                {t("clearing")}
              </th>
            </tr>
          </thead>
          <tbody>
            {rounds.map(({ asset, openMs, round }) => (
              <tr key={asset}>
                <th scope="row">{asset}</th>
                <td>{r(roundStatus(round, auction.phase, openMs))}</td>
                <td className="r">{f.price(round.floor)}</td>
                <td className="r">{round.lots}</td>
                <td className="r">{round.cleared ? f.price(round.price) : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}

function Record() {
  const t = useTranslations("risk.record");
  const f = useFormat();
  const explain = useExplain();
  const [chain] = useChains();
  const record = useRecord();
  const explorer = chain.blockExplorers?.default.url;
  const link = (tx: string, finalized: boolean) => (
    <>
      {explorer ? (
        <a href={`${explorer}/tx/${tx}`} target="_blank" rel="noreferrer" className="underline">
          {f.shortAddress(tx)}
        </a>
      ) : (
        f.shortAddress(tx)
      )}
      {!finalized && <span className="text-muted"> · {t("unfinal")}</span>}
    </>
  );
  if (record.isPending && !record.isError && record.fetchStatus === "idle") return null;
  return (
    <section aria-labelledby="record-title" className="mt-10 pb-16">
      <h2 id="record-title" className="label">
        {t("title")}
      </h2>
      <p className="mt-4 max-w-[60ch] text-[15px] text-body">{t("intro")}</p>
      {record.isError ? (
        <p role="alert" className="mt-4 text-[14px] text-[var(--error)]">
          {explain(record.error)}
        </p>
      ) : !record.data ? (
        <span className="skeleton mt-4 block h-[18px] w-[18rem]" />
      ) : record.data.rounds.length === 0 && record.data.actions.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">{t("none")}</p>
      ) : (
        <>
          {record.data.rounds.length > 0 && (
            <>
              <h3 className="mt-6 text-[15px] font-semibold text-strong">{t("roundsTitle")}</h3>
              <table className="ledger num">
                <caption className="sr-only">{t("roundsCaption")}</caption>
                <thead>
                  <tr>
                    <th scope="col">{t("open")}</th>
                    <th scope="col">{t("asset")}</th>
                    <th scope="col" className="r">{t("price")}</th>
                    <th scope="col" className="r">{t("sold")}</th>
                    <th scope="col" className="r">{t("paid")}</th>
                    <th scope="col" className="r">{t("taken")}</th>
                    <th scope="col">Tx</th>
                  </tr>
                </thead>
                <tbody>
                  {record.data.rounds.map((r) => (
                    <tr key={`${r.tx}-${r.asset}`}>
                      <th scope="row">{f.utc(new Date(Number(r.openMs)))}</th>
                      <td>{r.asset}</td>
                      <td className="r">{f.price(r.price)}</td>
                      <td className="r">{f.tokens(r.sold)}</td>
                      <td className="r">{f.usdg(r.paid)}</td>
                      <td className="r">{f.tokens(r.taken)}</td>
                      <td>{link(r.tx, r.finalized)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}
          {record.data.actions.length > 0 && (
            <>
              <h3 className="mt-6 text-[15px] font-semibold text-strong">{t("backstopTitle")}</h3>
              <table className="ledger num">
                <caption className="sr-only">{t("backstopCaption")}</caption>
                <thead>
                  <tr>
                    <th scope="col">{t("closure")}</th>
                    <th scope="col">{t("action")}</th>
                    <th scope="col" className="r">{t("amount")}</th>
                    <th scope="col">Tx</th>
                  </tr>
                </thead>
                <tbody>
                  {record.data.actions.map((a) => (
                    <tr key={`${a.tx}-${a.kind}-${a.position}`}>
                      <th scope="row">{f.utc(new Date(Number(a.closesMs)))}</th>
                      <td>{t(a.kind)}</td>
                      <td className="r">{f.usdg(a.amount)}</td>
                      <td>{link(a.tx, a.finalized)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}
        </>
      )}
    </section>
  );
}
