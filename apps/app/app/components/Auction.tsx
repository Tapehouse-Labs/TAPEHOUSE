// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { auction } from "@tapehouse/sdk";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { useChains } from "wagmi";
import { useDeployments } from "@/app/providers";
import { auctionNotes, type BidRow, bidMove, bidProblem, bidRows, defaultDeposit, escrowOf } from "@/lib/auction";
import { AUCTION_PHASE, useAuctionSend, useAuctionState } from "@/lib/auctionHook";
import type { AuctionState, RoundView } from "@/lib/auctionRead";
import { useWallet } from "@/lib/hooks";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import { Wallets } from "./Wallets";

type Send = ReturnType<typeof useAuctionSend>;

const utc = (ms: bigint) => new Date(Number(ms));

/** The reopening auction: where it stands, each asset's round and its bids, and the visitor's own bids once a wallet connects. */
export function Auction() {
  const t = useTranslations("auction.page");
  const f = useFormat();
  const explain = useExplain();
  const deployments = useDeployments();
  const [chain] = useChains();
  const state = useAuctionState();
  const data = state.data;
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">{t("title")}</h1>
          <p className="mt-6 max-w-[40ch] text-[19px] leading-snug text-body">{t("intro", { chain: chain.name })}</p>
          {data && <p className="label mt-3 num">{t("block", { block: new Intl.NumberFormat(f.locale).format(data.blockNumber) })}</p>}
        </div>
        {deployments.tapehouse.ReopeningAuction === undefined ? (
          <p className="mt-10 pb-16 text-[15px] text-body">{t("absent")}</p>
        ) : !data ? (
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
          <Screen state={data} />
        )}
      </div>
    </main>
  );
}

function Screen({ state }: { state: AuctionState }) {
  const shown = state.rounds.filter((r) => r.openMs === state.now.openMs);
  return (
    <>
      <Stage state={state} />
      <Rounds state={state} rounds={shown} />
      <Bids rounds={shown} />
      <Wallet state={state} />
      <Clear state={state} rounds={shown} />
    </>
  );
}

function Stage({ state }: { state: AuctionState }) {
  const t = useTranslations("auction");
  const say = useMsg();
  const f = useFormat();
  const { now, constants } = state;
  return (
    <section aria-labelledby="stage-title" className="mt-10">
      <h2 id="stage-title" className="label">
        {t("stage.title")}
      </h2>
      <p role="status" className="mt-4 max-w-[62ch] text-[17px] leading-snug text-strong">
        {t(`stage.${now.stage}`, {
          open: f.utc(utc(now.openMs)),
          reveal: f.utc(utc(now.openMs - constants.revealMs)),
          closes: f.utc(utc(now.endsMs ?? now.openMs)),
        })}
      </p>
      <ul aria-label={t("notes.label")} className="notes mt-5 max-w-[62ch]">
        {auctionNotes().map((n) => (
          <li key={n.message.key} data-tone={n.tone}>
            {n.term && <b className="font-semibold text-strong">{say(n.term)}. </b>}
            {say(n.message)}
          </li>
        ))}
      </ul>
    </section>
  );
}

function proposedPrice({ round }: RoundView): bigint | undefined {
  return auction.clearingPrice(round.bids, round.supply, round.floor);
}

function Rounds({ state, rounds }: { state: AuctionState; rounds: RoundView[] }) {
  const t = useTranslations("auction.round");
  const f = useFormat();
  const cleared = rounds.filter((r) => r.round.cleared);
  return (
    <section aria-labelledby="round-title" className="mt-10">
      <h2 id="round-title" className="label">
        {t("title")}
      </h2>
      {state.now.stage === "idle" || rounds.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">{t("none")}</p>
      ) : (
        <>
          <div className="mt-3 overflow-x-auto">
            <table className="exposure num">
              <caption className="sr-only">{t("caption")}</caption>
              <thead>
                <tr>
                  <th scope="col">{t("asset")}</th>
                  <th scope="col" className="r">
                    {t("floor")}
                  </th>
                  <th scope="col" className="r">
                    {t("lots")}
                  </th>
                  <th scope="col" className="r">
                    {t("supply")}
                  </th>
                  <th scope="col" className="r">
                    {t("deposits")}
                  </th>
                  <th scope="col" className="r">
                    {t("bids")}
                  </th>
                  <th scope="col" className="r">
                    {t("price")}
                  </th>
                </tr>
              </thead>
              <tbody>
                {rounds.map((view) => {
                  const { asset, round } = view;
                  const proposed = proposedPrice(view);
                  const known = round.floor > 0n;
                  return (
                    <tr key={asset}>
                      <th scope="row">{asset}</th>
                      <td className="r">
                        {known ? f.price(round.floor) : "—"}
                        {known && <small className="block whitespace-nowrap">{round.sealMs > 0n ? t("fromSeal") : t("fromQuote")}</small>}
                      </td>
                      <td className="r">{round.lots.length}</td>
                      <td className="r">{f.tokens(round.supply)}</td>
                      <td className="r">{f.usdg(round.deposits)}</td>
                      <td className="r">{round.bids.length}</td>
                      <td className="r">
                        {round.cleared ? f.price(round.price) : known && proposed !== undefined ? f.price(proposed) : "—"}
                        {known && (round.cleared || proposed !== undefined) && <small className="block whitespace-nowrap">{round.cleared ? t("settled") : t("proposed")}</small>}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          {cleared.length > 0 && (
            <div className="mt-8 overflow-x-auto">
              <table className="exposure num">
                <caption className="label">{t("resultTitle")}</caption>
                <thead>
                  <tr>
                    <th scope="col">{t("asset")}</th>
                    <th scope="col" className="r">
                      {t("price")}
                    </th>
                    <th scope="col" className="r">
                      {t("sold")}
                    </th>
                    <th scope="col" className="r">
                      {t("paid")}
                    </th>
                    <th scope="col" className="r">
                      {t("taken")}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {cleared.map(({ asset, round }) => (
                    <tr key={asset}>
                      <th scope="row">{asset}</th>
                      <td className="r">{f.price(round.price)}</td>
                      <td className="r">{f.tokens(round.sold)}</td>
                      <td className="r">{f.usdg(round.paid)}</td>
                      <td className="r">{f.tokens(round.taken)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </section>
  );
}

function Bids({ rounds }: { rounds: RoundView[] }) {
  const t = useTranslations("auction.bids");
  const f = useFormat();
  const { address } = useWallet();
  const withBids = rounds.filter((r) => r.round.bids.length > 0);
  return (
    <section aria-labelledby="bids-title" className="mt-10">
      <h2 id="bids-title" className="label">
        {t("title")}
      </h2>
      {withBids.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">{t("none")}</p>
      ) : (
        withBids.map(({ asset, round }) => (
          <details key={asset} className="mt-3 border-t border-rule pt-3">
            <summary className="cursor-pointer text-[15px] text-strong">{t("summary", { asset, count: round.bids.length })}</summary>
            <div className="mt-3 overflow-x-auto">
              <table className="exposure num">
                <caption className="sr-only">{t("caption", { asset })}</caption>
                <thead>
                  <tr>
                    <th scope="col">{t("bidder")}</th>
                    <th scope="col" className="r">
                      {t("price")}
                    </th>
                    <th scope="col" className="r">
                      {t("quantity")}
                    </th>
                    <th scope="col" className="r">
                      {t("escrow")}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {round.bids
                    .toSorted((a, b) => (a.price === b.price ? 0 : a.price > b.price ? -1 : 1))
                    .map((bid, i) => (
                      <tr key={i}>
                        <th scope="row">
                          {bid.bidder.toLowerCase() === address?.toLowerCase() ? (
                            <span className="text-strong">{t("you")}</span>
                          ) : (
                            f.shortAddress(bid.bidder)
                          )}
                        </th>
                        <td className="r">{f.price(bid.price)}</td>
                        <td className="r">{f.tokens(bid.quantity)}</td>
                        <td className="r">{f.usdg(bid.escrow)}</td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>
          </details>
        ))
      )}
    </section>
  );
}

function Wallet({ state }: { state: AuctionState }) {
  const t = useTranslations("auction.page");
  const w = useTranslations("wallet");
  const wallet = useWallet();
  const [chain] = useChains();
  if (wallet.state === "connected" && state.wallet)
    return (
      <>
        <Form state={state} />
        <Mine state={state} />
      </>
    );
  if (wallet.state === "connected") return null;
  return (
    <section className="mt-10">
      <p className="max-w-[40ch] text-[15px] text-body" aria-live="polite">
        {wallet.state === "wrong-network"
          ? t.rich("wrongNetwork", { chain: chain.name, b: (name) => <span className="text-strong">{name}</span> })
          : t("connect")}
      </p>
      <div className="mt-6 flex min-h-[44px] flex-wrap items-center gap-3">
        {wallet.state === "wrong-network" ? (
          <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
            {wallet.pending ? w("waiting") : t("switch", { chain: chain.name })}
          </button>
        ) : (
          (wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />
        )}
      </div>
    </section>
  );
}

function Form({ state }: { state: AuctionState }) {
  const t = useTranslations("auction.bid");
  const f = useFormat();
  const say = useMsg();
  const send = useAuctionSend();
  const { constants, now, wallet } = state;
  const [asset, setAsset] = useState(state.assets[0] ?? "");
  const [typedQuantity, setQuantity] = useState("");
  const [typedPrice, setPrice] = useState("");
  const [typedDeposit, setDeposit] = useState<string>();
  const quantity = f.parseAmount(typedQuantity, 18);
  const price = f.parseAmount(typedPrice, 8);
  const depositText = typedDeposit ?? f.unparse(defaultDeposit({ quantity, price, bond: constants.bond }), 6);
  const deposit = f.parseAmount(depositText, 6);
  const opened = state.rounds.find((r) => r.asset === asset && r.openMs === now.openMs)?.round.floor ?? 0n;
  const floor = opened > 0n ? opened : (state.previews[asset] ?? 0n);
  const open = now.stage === "committing";
  const problem =
    quantity !== undefined && price !== undefined && deposit !== undefined
      ? bidProblem({ quantity, price, deposit, floor: floor === 0n ? price : floor, bond: constants.bond, minBid: constants.minBid })
      : undefined;
  const ready = open && quantity !== undefined && price !== undefined && deposit !== undefined && !problem;
  const edit = (set: (value: string) => void) => (e: React.ChangeEvent<HTMLInputElement>) => {
    set(f.normalize(e.target.value));
    if (!send.pending) send.reset();
  };
  if (!open)
    return (
      <section aria-labelledby="bid-title" className="mt-10">
        <h2 id="bid-title" className="label">
          {t("title")}
        </h2>
        <p className="mt-4 text-[15px] text-body">{t("closed")}</p>
      </section>
    );
  return (
    <section aria-labelledby="bid-title" className="mt-10">
      <h2 id="bid-title" className="label">
        {t("title")}
      </h2>
      <form
        className="mt-5 grid max-w-[40rem] gap-4 text-[15px] text-strong"
        onSubmit={(e) => {
          e.preventDefault();
          if (ready) send.send({ kind: "commit", asset, openMs: now.openMs, quantity, price, deposit });
        }}
      >
        <label className="flex flex-wrap items-center gap-3">
          <span className="min-w-[8rem]">{t("asset")}</span>
          <select className="picker" aria-label={t("assetLabel")} value={asset} onChange={(e) => setAsset(e.target.value)}>
            {state.assets.map((a) => (
              <option key={a} value={a}>
                {a}
              </option>
            ))}
          </select>
          {floor > 0n && <span className="text-[13px] text-muted num">{t("floor", { floor: f.price(floor) })}</span>}
        </label>
        <label className="flex flex-wrap items-center gap-3">
          <span className="min-w-[8rem]">{t("quantity")}</span>
          <input
            className="amount"
            inputMode="decimal"
            autoComplete="off"
            placeholder="0"
            aria-label={t("quantityLabel")}
            aria-invalid={typedQuantity !== "" && quantity === undefined}
            value={typedQuantity}
            onChange={edit(setQuantity)}
          />
          <span>{asset}</span>
        </label>
        <label className="flex flex-wrap items-center gap-3">
          <span className="min-w-[8rem]">{t("price")}</span>
          <input
            className="amount"
            inputMode="decimal"
            autoComplete="off"
            placeholder="0"
            aria-label={t("priceLabel")}
            aria-invalid={typedPrice !== "" && price === undefined}
            value={typedPrice}
            onChange={edit(setPrice)}
          />
          <span>USD</span>
        </label>
        <label className="flex flex-wrap items-center gap-3">
          <span className="min-w-[8rem]">{t("deposit")}</span>
          <input
            className="amount"
            inputMode="decimal"
            autoComplete="off"
            aria-label={t("depositLabel")}
            aria-invalid={deposit === undefined || problem !== undefined}
            value={depositText}
            onChange={edit(setDeposit)}
          />
          <span>USDG</span>
        </label>
        {quantity !== undefined && price !== undefined && (
          <p className="max-w-[60ch] text-[13px] text-body num">{t("escrow", { escrow: f.usdg(escrowOf(quantity, price)) })}</p>
        )}
        {problem && (
          <p role="alert" className="max-w-[60ch] text-[14px] text-[var(--error)]">
            {say(problem)}
          </p>
        )}
        <ul className="notes max-w-[62ch]">
          {floor === 0n ? (
            <li data-tone="warning">{t("floorUnset")}</li>
          ) : (
            opened === 0n && <li>{t("floorProjected")}</li>
          )}
          <li>{t("hides")}</li>
          <li>{t("saved")}</li>
          <li>{t("flow")}</li>
        </ul>
        <p className="text-[13px] text-muted num">{t("balance", { amount: f.usdg(wallet?.usdg ?? 0n), bond: f.usdg(constants.bond) })}</p>
        <div>
          <button type="submit" className="btn" disabled={!ready || send.pending}>
            {send.phase ? say(AUCTION_PHASE[send.phase]) : t("submit")}
          </button>
        </div>
      </form>
      <Outcome send={send} />
    </section>
  );
}

function Mine({ state }: { state: AuctionState }) {
  const t = useTranslations("auction.mine");
  const f = useFormat();
  const send = useAuctionSend();
  const rows = bidRows(state.wallet?.saved ?? [], state.wallet?.listed ?? []);
  return (
    <section aria-labelledby="mine-title" className="mt-10">
      <h2 id="mine-title" className="label">
        {t("title")}
      </h2>
      {rows.length === 0 ? (
        <p className="mt-4 text-[15px] text-body">{t("none")}</p>
      ) : (
        <div className="mt-3 overflow-x-auto">
          <table className="exposure num">
            <caption className="sr-only">{t("caption")}</caption>
            <thead>
              <tr>
                <th scope="col">{t("round")}</th>
                <th scope="col">{t("asset")}</th>
                <th scope="col">{t("bid")}</th>
                <th scope="col">{t("status")}</th>
                <th scope="col">{t("next")}</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={`${row.asset}:${row.openMs}:${row.saved?.commitment ?? row.index}`}>
                  <td>{f.utc(utc(row.openMs))}</td>
                  <th scope="row">{row.asset}</th>
                  <td>{t("offer", { quantity: f.tokens(row.quantity), price: f.price(row.price) })}</td>
                  <td>
                    {t(`state.${row.status}`)}
                    {!row.saved && <small>{t("unsaved")}</small>}
                  </td>
                  <td>
                    <Move row={row} state={state} send={send} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <Outcome send={send} />
    </section>
  );
}

function Move({ row, state, send }: { row: BidRow; state: AuctionState; send: Send }) {
  const t = useTranslations("auction.mine");
  const say = useMsg();
  const cleared = state.rounds.find((r) => r.asset === row.asset && r.openMs === row.openMs)?.round.cleared ?? false;
  const move = bidMove({
    status: row.status,
    nowMs: state.nowMs,
    openMs: row.openMs,
    cleared,
    revealMs: state.constants.revealMs,
    clearMs: state.constants.clearMs,
  });
  const label = (text: string) => (send.phase && send.pending ? say(AUCTION_PHASE[send.phase]) : text);
  if (move === "reveal" && row.saved) {
    const bid = row.saved;
    return (
      <button type="button" className="btn" disabled={send.pending} onClick={() => send.send({ kind: "reveal", bid })}>
        {label(t("reveal"))}
      </button>
    );
  }
  if (move === "forfeit" && row.saved) {
    const bid = row.saved;
    return (
      <>
        <button type="button" className="btn btn-quiet" disabled={send.pending} onClick={() => send.send({ kind: "forfeit", bid })}>
          {label(t("forfeit"))}
        </button>
        <small>{t("forfeitNote")}</small>
      </>
    );
  }
  if (move === "claim" && row.index !== undefined) {
    const index = row.index;
    return (
      <button
        type="button"
        className="btn"
        disabled={send.pending}
        onClick={() => send.send({ kind: "claim", asset: row.asset, openMs: row.openMs, index })}
      >
        {label(t("claim"))}
      </button>
    );
  }
  if (move === "waitReveal" || move === "waitClear") return <span className="text-[13px] text-muted">{t(move)}</span>;
  return <span aria-hidden="true">—</span>;
}

function Clear({ state, rounds }: { state: AuctionState; rounds: RoundView[] }) {
  const t = useTranslations("auction.clear");
  const f = useFormat();
  const say = useMsg();
  const send = useAuctionSend();
  const wallet = useWallet();
  const open = rounds.filter((r) => !r.round.cleared && r.round.floor > 0n);
  if (state.now.stage !== "clearing" || open.length === 0) return null;
  return (
    <section aria-labelledby="clear-title" className="mt-10 pb-16">
      <h2 id="clear-title" className="label">
        {t("title")}
      </h2>
      <p className="mt-4 max-w-[62ch] text-[15px] text-body">{t("body")}</p>
      <div className="mt-5 grid gap-3 text-[14px]">
        {open.map((view) => {
          const price = proposedPrice(view);
          return (
            <div key={view.asset} className="flex flex-wrap items-center gap-3">
              {price === undefined ? (
                <span className="text-body">{t("none", { asset: view.asset })}</span>
              ) : (
                <button
                  type="button"
                  className="btn btn-quiet"
                  disabled={send.pending || wallet.state !== "connected"}
                  onClick={() => send.send({ kind: "clear", asset: view.asset, openMs: view.openMs, price })}
                >
                  {send.phase && send.pending ? say(AUCTION_PHASE[send.phase]) : t("button", { asset: view.asset, price: f.price(price) })}
                </button>
              )}
            </div>
          );
        })}
      </div>
      <Outcome send={send} />
    </section>
  );
}

function Outcome({ send }: { send: Send }) {
  const t = useTranslations("auction.page");
  const explain = useExplain();
  if (send.error)
    return (
      <p role="alert" className="mt-3 max-w-[60ch] text-[14px] text-[var(--error)]">
        {explain(send.error)}
      </p>
    );
  if (send.succeeded && !send.pending)
    return (
      <p role="status" className="mt-3 text-[14px] text-[var(--gain)]">
        {t("done")}
      </p>
    );
  return null;
}
