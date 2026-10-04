// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { Fragment, lazy, type ReactNode, Suspense, useState } from "react";
import type { Address } from "viem";
import { useChains } from "wagmi";
import type { AccountState } from "@/lib/accountRead";
import { tokens, usdg } from "@/lib/amounts";
import { backstopNotes, coverNotes, exposures, ledger, lendingNotes, percent, span, supplyNotes } from "@/lib/earn";
import { useEarnState } from "@/lib/earnHook";
import type { EarnState, Holding } from "@/lib/earnRead";
import { explain } from "@/lib/errors";
import { useWallet } from "@/lib/hooks";
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
  const wallet = useWallet();
  const [chain] = useChains();
  const connected = wallet.state === "connected";
  const state = useEarnState(holder);
  const data = state.data;

  const lent: Record<string, bigint> = {};
  for (const p of positions?.positions ?? [])
    for (const s of p.stocks) lent[s.asset] = (lent[s.asset] ?? 0n) + s.lent;
  const rows = data ? ledger(data, positions ? lent : {}) : [];
  const worth = data
    ? [data.supply?.lender, data.backstop?.depositor, data.cover?.writer].reduce((sum, h) => sum + (h?.assets ?? 0n), 0n)
    : 0n;
  const lentWords = Object.entries(lent)
    .filter(([, amount]) => amount > 0n)
    .map(([asset, amount]) => `${tokens(amount)} ${asset}`);

  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">Earning · supply, backstop and cover</h1>
          {connected && data && worth > 0n ? (
            <p key={String(worth)} className="figure figure-earn num swap mt-3" aria-live="polite">
              ${usdg(worth)}
            </p>
          ) : (
            <p className="figure mt-3" data-empty>
              None yet
            </p>
          )}
          <div key={wallet.state} className="swap" aria-live="polite">
            <p className="mt-6 max-w-[40ch] text-[19px] leading-snug text-body">
              {wallet.state === "wrong-network" ? (
                <>
                  Your wallet is on another network. Tapehouse runs on <span className="text-strong">{chain.name}</span>.
                </>
              ) : !connected ? (
                "Connect a wallet to deposit. Every term below is read from the chain, before you need it."
              ) : worth > 0n ? (
                `What your account's shares are worth now${lentWords.length > 0 ? `, besides ${lentWords.join(" and ")} lent from it` : ""}. Withdrawals and claimed Stock Tokens go to your wallet.`
              ) : (
                "Open a place below and read its terms before you deposit. Withdrawals go back to your wallet."
              )}
            </p>
            {!connected && (
              <div className="mt-8 flex min-h-[44px] flex-wrap items-center gap-3">
                {wallet.state === "wrong-network" ? (
                  <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
                    {wallet.pending ? "Waiting for your wallet…" : `Switch to ${chain.name}`}
                  </button>
                ) : (
                  (wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />
                )}
              </div>
            )}
          </div>
          {data && <p className="label mt-3 num">Read at block {data.blockNumber.toLocaleString("en-US")}</p>}
        </div>

        {account}

        <section aria-labelledby="ledger-title" className="mt-10">
          <h2 id="ledger-title" className="label">
            Where your money can work
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
                Each place your money can work: what it earns, how it leaves, what it bears, and what your account holds there
              </caption>
              <thead>
                <tr>
                  <th scope="col">Place</th>
                  <th scope="col">Earns</th>
                  <th scope="col">Leaves</th>
                  <th scope="col">Bears</th>
                  <th scope="col" className="r">
                    You hold
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
                            {row.place}
                          </button>
                          <small>{row.what}</small>
                        </th>
                        <td data-label="Earns">
                          {row.earns}
                          <small>{row.earnsNote}</small>
                        </td>
                        <td data-label="Leaves" data-closed={row.closed ? "" : undefined}>
                          {row.leaves}
                          <small>{row.leavesNote}</small>
                        </td>
                        <td data-label="Bears">
                          {row.bears}
                          <small>{row.bearsNote}</small>
                        </td>
                        <td data-label="You hold" className="r">
                          {row.holds ?? "—"}
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
                      <span>Write cover</span>
                      <small>USDG, apart from the backstop</small>
                    </th>
                    <td colSpan={4}>Gap cover is not deployed on this chain yet.</td>
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
  if (id === "supply" && state.supply) {
    const { market, lender } = state.supply;
    return (
      <div className="panel">
        <dl className="facts">
          <Fact term="Borrowers pay" value={`${percent(market.borrowRate)} a year`} />
          <Fact term="Lent" value={`${usdg(market.debt)} USDG`} />
          <Fact term="Not lent" value={`${usdg(market.idle)} USDG`} />
          <Fact term="Vault holds" value={`${usdg(market.totalAssets)} USDG`} />
        </dl>
        <Notes notes={supplyNotes(state.supply)} />
        {actions?.vault("supply", lender)}
      </div>
    );
  }
  if (id === "cover" && state.cover) {
    const { writers, writer } = state.cover;
    return (
      <div className="panel">
        <dl className="facts">
          <Fact term="Free for covers" value={`${usdg(writers.capacity)} USDG`} />
          <Fact term="Reserved" value={`${usdg(writers.reserved)} USDG`} />
          <Fact term="Premiums waiting" value={`${usdg(writers.premiums)} USDG`} />
          <Fact term="Covers outstanding" value={writers.outstanding.toString()} />
        </dl>
        <Notes notes={coverNotes(state.cover)} />
        {actions?.vault("cover", writer)}
      </div>
    );
  }
  if (id === "backstop" && state.backstop) {
    const b = state.backstop.state;
    return (
      <div className="panel">
        <table className="exposure num">
          <caption className="label">Exposure this closure, USDG</caption>
          <thead>
            <tr>
              <th scope="col">Positions</th>
              <th scope="col" className="r">
                Limit
              </th>
              <th scope="col" className="r">
                Covered
              </th>
              <th scope="col" className="r">
                Left
              </th>
            </tr>
          </thead>
          <tbody>
            {exposures(b, state.timestamp).map((e) => (
              <tr key={e.name}>
                <th scope="row">
                  {e.name}
                  {e.change && <small>then {e.change}</small>}
                </th>
                <td className="r">{e.limit}</td>
                <td className="r">{e.covered}</td>
                <td className="r">{e.left}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <p className="mt-2 text-[13px] text-muted num">What is left is never more than the {usdg(b.held)} USDG the backstop holds.</p>
        <Notes notes={backstopNotes(state.backstop)} />
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
        <Fact term="Recall notice" value={span(terms.notice)} />
        <Fact term="Held for recalls" value={`${tokens(terms.locked)} ${asset}`} />
        <Fact term="Lent of the vault" value={percent(terms.utilization)} />
        <Fact term="Borrowable" value={`${tokens(terms.borrowable)} ${asset}`} />
      </dl>
      <Notes notes={lendingNotes(asset, terms)} />
      {actions?.lend(lending)}
    </div>
  );
}

function Notes({ notes }: { notes: Note[] }) {
  return (
    <ul aria-label="Before you deposit" className="notes mt-5 max-w-[62ch]">
      {notes.map((n) => (
        <li key={n.text} data-tone={n.tone}>
          {n.term && <b className="font-semibold text-strong">{n.term}. </b>}
          {n.text}
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
