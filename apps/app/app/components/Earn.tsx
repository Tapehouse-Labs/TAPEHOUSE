// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useTranslations } from "next-intl";
import { Fragment, lazy, type ReactNode, Suspense, useState } from "react";
import type { Address } from "viem";
import { useChains } from "wagmi";
import type { AccountState } from "@/lib/accountRead";
import { backstopNotes, coverNotes, exposures, ledger, lendingNotes, span, supplyNotes } from "@/lib/earn";
import { useEarnState } from "@/lib/earnHook";
import type { EarnState, Holding } from "@/lib/earnRead";
import { useWallet } from "@/lib/hooks";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import type { Vault } from "@/lib/intent";
import type { Note } from "@/lib/notes";
import { Wallets } from "./Wallets";

const EarnDesk = lazy(() => import("./EarnDesk").then((m) => ({ default: m.EarnDesk })));

/** What a connected account can do in each open row: the composers, which load only once a wallet connects. */
export type Actions = {
  vault: (vault: Vault, holding: Holding | undefined) => ReactNode;
  backstop: (backstop: NonNullable<EarnState["backstop"]>, timestamp: bigint) => ReactNode;
  lend: (lending: EarnState["lending"][number]) => ReactNode;
};

/** The Earn surface: the ledger of every place money can work, with the account's composers once a wallet connects. */
export function Earn() {
  const wallet = useWallet();
  const row = useState<string | undefined>("supply");
  if (wallet.state !== "connected") return <EarnView row={row} />;
  return (
    <Suspense fallback={<EarnView row={row} />}>
      <EarnDesk row={row} />
    </Suspense>
  );
}

/** The open row of the ledger and its setter, kept above the composers' chunk so a row opened while it loads stays open. */
export type OpenRow = ReturnType<typeof useState<string | undefined>>;

/**
 * What the account's shares are worth under the bar, then one ledger of every place money can work, read at one block.
 * `row` is the open row, `holder` the connected account, `positions` its margin positions, `actions` its composers and `account` its smart
 * account's block; a visitor gets the terms alone.
 */
export function EarnView({
  row: [open, setOpen],
  holder,
  positions,
  actions,
  account,
}: {
  row: OpenRow;
  holder?: Address;
  positions?: AccountState;
  actions?: Actions;
  account?: ReactNode;
}) {
  const t = useTranslations("earn.page");
  const e = useTranslations("earn");
  const w = useTranslations("wallet");
  const say = useMsg();
  const f = useFormat();
  const explain = useExplain();
  const wallet = useWallet();
  const [chain] = useChains();
  const connected = wallet.state === "connected";
  const state = useEarnState(holder);
  const data = state.data;

  const lent: Record<string, bigint> = {};
  for (const p of positions?.positions ?? [])
    for (const s of p.stocks) lent[s.asset] = (lent[s.asset] ?? 0n) + s.lent;
  const rows = data ? ledger(f, data, positions ? lent : {}) : [];
  const worth = data
    ? [data.supply?.lender, data.backstop?.depositor, data.cover?.writer].reduce((sum, h) => sum + (h?.assets ?? 0n), 0n)
    : 0n;
  const lentWords = Object.entries(lent)
    .filter(([, amount]) => amount > 0n)
    .map(([asset, amount]) => `${f.tokens(amount)} ${asset}`);

  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">{t("title")}</h1>
          {connected && data && worth > 0n ? (
            <p key={String(worth)} className="figure figure-earn num swap mt-3" aria-live="polite">
              ${f.usdg(worth)}
            </p>
          ) : (
            <p className="figure mt-3" data-empty>
              {t("none")}
            </p>
          )}
          <div key={wallet.state} className="swap" aria-live="polite">
            <p className="mt-6 max-w-[40ch] text-[19px] leading-snug text-body">
              {wallet.state === "wrong-network" ? (
                t.rich("wrongNetwork", {
                  chain: chain.name,
                  b: (name) => <span className="text-strong">{name}</span>,
                })
              ) : !connected ? (
                t("connect")
              ) : worth > 0n ? (
                lentWords.length > 0 ? (
                t("worthLent", { lent: f.list(lentWords) })
              ) : (
                t("worth")
              )
              ) : (
                t("open")
              )}
            </p>
            {!connected && (
              <div className="mt-8 flex min-h-[44px] flex-wrap items-center gap-3">
                {wallet.state === "wrong-network" ? (
                  <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
                    {wallet.pending ? w("waiting") : t("switch", { chain: chain.name })}
                  </button>
                ) : (
                  (wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />
                )}
              </div>
            )}
          </div>
          {data && <p className="label mt-3 num">{t("block", { block: new Intl.NumberFormat(f.locale).format(data.blockNumber) })}</p>}
        </div>

        {account}

        <section aria-labelledby="ledger-title" className="mt-10">
          <h2 id="ledger-title" className="label">
            {t("ledgerTitle")}
          </h2>
          {!data ? (
            <div aria-busy="true" className="mt-4">
              {state.error ? (
                <p role="alert" className="text-[14px] text-[var(--error)]">
                  {explain(state.error)}
                </p>
              ) : (
                <span className="skeleton block h-[18px] w-[18rem]" />
              )}
            </div>
          ) : (
            <table className="ledger num">
              <caption className="sr-only">
                {t("caption")}
              </caption>
              <thead>
                <tr>
                  <th scope="col">{t("place")}</th>
                  <th scope="col">{t("earns")}</th>
                  <th scope="col">{t("leaves")}</th>
                  <th scope="col">{t("bears")}</th>
                  <th scope="col" className="r">
                    {t("hold")}
                  </th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => {
                  const isOpen = open === row.id;
                  return (
                    <Fragment key={row.id}>
                      <tr data-open={isOpen ? "" : undefined}>
                        <th scope="row">
                          <button
                            type="button"
                            aria-expanded={isOpen}
                            aria-controls={`panel-${row.id}`}
                            onClick={() => setOpen(isOpen ? undefined : row.id)}
                          >
                            {say(row.place)}
                          </button>
                          <small>{say(row.what)}</small>
                        </th>
                        <td data-label={t("earns")}>
                          {say(row.earns)}
                          <small>{say(row.earnsNote)}</small>
                        </td>
                        <td data-label={t("leaves")} data-closed={row.closed ? "" : undefined}>
                          {say(row.leaves)}
                          <small>{say(row.leavesNote)}</small>
                        </td>
                        <td data-label={t("bears")}>
                          {say(row.bears)}
                          <small>{say(row.bearsNote)}</small>
                        </td>
                        <td data-label={t("hold")} className="r">
                          {row.holds ? say(row.holds) : "—"}
                        </td>
                      </tr>
                      {isOpen && (
                        <tr className="ledger-panel">
                          <td colSpan={5} id={`panel-${row.id}`}>
                            <Panel id={row.id} state={data} actions={actions} />
                          </td>
                        </tr>
                      )}
                    </Fragment>
                  );
                })}
                {!data.cover && (
                  <tr data-absent="">
                    <th scope="row">
                      <span>{e("ledger.cover.place")}</span>
                      <small>{e("ledger.cover.what")}</small>
                    </th>
                    <td colSpan={4}>{t("absent")}</td>
                  </tr>
                )}
              </tbody>
            </table>
          )}
        </section>
      </div>
    </main>
  );
}

function Panel({ id, state, actions }: { id: string; state: EarnState; actions: Actions | undefined }) {
  const t = useTranslations("earn");
  const p = useTranslations("earn.page");
  const say = useMsg();
  const f = useFormat();
  const amount = (value: bigint) => t("usdgAmount", { amount: f.usdg(value) });
  if (id === "supply" && state.supply) {
    const { market, lender } = state.supply;
    return (
      <div className="panel">
        <dl className="facts">
          <Fact term={p("borrowersPay")} value={p("aYear", { rate: f.percent(market.borrowRate) })} />
          <Fact term={p("lent")} value={amount(market.debt)} />
          <Fact term={p("notLent")} value={amount(market.idle)} />
          <Fact term={p("vaultHolds")} value={amount(market.totalAssets)} />
        </dl>
        <Notes notes={supplyNotes(f, state.supply)} />
        {actions?.vault("supply", lender)}
      </div>
    );
  }
  if (id === "cover" && state.cover) {
    const { writers, writer } = state.cover;
    return (
      <div className="panel">
        <dl className="facts">
          <Fact term={p("freeForCovers")} value={amount(writers.capacity)} />
          <Fact term={p("reserved")} value={amount(writers.reserved)} />
          <Fact term={p("premiumsWaiting")} value={amount(writers.premiums)} />
          <Fact term={p("coversOutstanding")} value={new Intl.NumberFormat(f.locale).format(writers.outstanding)} />
        </dl>
        <Notes notes={coverNotes(f, state.cover)} />
        {actions?.vault("cover", writer)}
      </div>
    );
  }
  if (id === "backstop" && state.backstop) {
    const b = state.backstop.state;
    return (
      <div className="panel">
        <table className="exposure num">
          <caption className="label">{p("exposureCaption")}</caption>
          <thead>
            <tr>
              <th scope="col">{p("positions")}</th>
              <th scope="col" className="r">
                {p("limit")}
              </th>
              <th scope="col" className="r">
                {p("covered")}
              </th>
              <th scope="col" className="r">
                {p("left")}
              </th>
            </tr>
          </thead>
          <tbody>
            {exposures(f, b, state.timestamp).map((e, i) => (
              <tr key={i}>
                <th scope="row">
                  {say(e.name)}
                  {e.change && <small>{p("then", { change: say(e.change) })}</small>}
                </th>
                <td className="r">{e.limit}</td>
                <td className="r">{e.covered}</td>
                <td className="r">{e.left}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <p className="mt-2 text-[13px] text-muted num">{p("held", { amount: f.usdg(b.held) })}</p>
        <Notes notes={backstopNotes(f, state.backstop)} />
        {actions?.backstop(state.backstop, state.timestamp)}
      </div>
    );
  }
  const lending = state.lending.find((l) => `lend-${l.asset}` === id);
  if (!lending) return null;
  const { asset, terms } = lending;
  return (
    <div className="panel">
      <dl className="facts">
        <Fact term={p("recallNotice")} value={say(span(terms.notice))} />
        <Fact term={p("heldForRecalls")} value={t("tokenAmount", { amount: f.tokens(terms.locked), asset })} />
        <Fact term={p("lentOfVault")} value={f.percent(terms.utilization)} />
        <Fact term={p("borrowable")} value={t("tokenAmount", { amount: f.tokens(terms.borrowable), asset })} />
      </dl>
      <Notes notes={lendingNotes(f, asset, terms)} />
      {actions?.lend(lending)}
    </div>
  );
}

function Notes({ notes }: { notes: Note[] }) {
  const t = useTranslations("earn.page");
  const say = useMsg();
  return (
    <ul aria-label={t("beforeDeposit")} className="notes mt-5 max-w-[62ch]">
      {notes.map((n) => (
        <li key={n.message.key} data-tone={n.tone}>
          {n.term && <b className="font-semibold text-strong">{say(n.term)}. </b>}
          {say(n.message)}
        </li>
      ))}
    </ul>
  );
}

function Fact({ term, value }: { term: string; value: string }) {
  return (
    <div>
      <dt className="label">{term}</dt>
      <dd className="mt-1 text-strong">{value}</dd>
    </div>
  );
}
