// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { stockLending } from "@tapehouse/sdk";
import { useQueries } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useState } from "react";
import type { Hex } from "viem";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import type { AccountState, Market, Position, Stock } from "@/lib/accountRead";
import { useSend } from "@/lib/accountHook";
import { useMinute } from "@/lib/hooks";
import { useExplain, useFormat } from "@/lib/i18n";

/** What a lent Stock Token counts for less in a position's equity, in basis points: the accounts' `RECALL_HAIRCUT`. */
const RECALL_HAIRCUT = 500n;
/** The least a recall that is not a position's first and whole one may queue: the accounts' `MIN_RECALL`. */
const MIN_RECALL = 10n ** 16n;
/** The recalls of one Stock Token an owner may hold open: one fewer than the accounts' `MAX_RECALLS`. */
const OWNER_RECALLS = 3;

/** The account's positions, hanging off the stem: each opens in place on what it holds, lends and recalls. */
export function Positions({
  state,
  selected,
  onSelect,
}: {
  state: AccountState;
  selected: Hex;
  onSelect: (id: Hex) => void;
}) {
  const t = useTranslations("account.positions");
  const f = useFormat();
  if (state.positions.length === 0) return null;
  return (
    <section aria-labelledby="positions-title" className="mt-14">
      <h2 id="positions-title" className="label">
        {t("title")}
      </h2>
      <div className="mt-3">
        {state.positions.map((p) => (
          <details
            key={p.id}
            className="position"
            data-selected={p.id === selected ? "" : undefined}
            onToggle={(e) => (e.currentTarget.open ? onSelect(p.id) : undefined)}
          >
            <summary>
              <span>
                <b className="font-semibold text-strong">{p.name}</b>
                {p.name !== "Cross" && <span className="text-[13px] text-muted"> {t("isolated")}</span>}
                <span className="mt-0.5 block text-[12.5px] text-muted num">
                  {p.stocks.map((s) => `${f.tokens(s.held + s.lent + s.inBaskets)} ${s.asset}`).join(" · ") || t("noStockTokens")}
                  {p.usdg > 0n && ` · ${f.usdg(p.usdg)} USDG`}
                  {p.weth > 0n && ` · ${f.tokens(p.weth)} WETH`}
                </span>
              </span>
              <span className="text-right num">
                <b className="font-semibold text-strong">{f.usd(p.equity)}</b> <span className="text-[12.5px] text-muted">{t("equity")}</span>
                <span className="block text-[12.5px] text-muted">
                  {t("required", { amount: f.usd(p.requirement) })}
                  {p.debt + p.premium > 0n && ` · ${t("owes", { amount: f.usdg(p.debt + p.premium) })}`}
                </span>
              </span>
            </summary>
            <PositionDetail state={state} position={p} />
          </details>
        ))}
      </div>
    </section>
  );
}

function PositionDetail({ state, position }: { state: AccountState; position: Position }) {
  const t = useTranslations("account.positions");
  const f = useFormat();
  const leverage = new Intl.NumberFormat(f.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  return (
    <div className="grid gap-5 pb-6">
      <dl className="grid grid-cols-2 gap-x-6 gap-y-2 text-[13px] sm:grid-cols-4 num">
        <Fact term={t("debt")} value={`${f.usdg(position.debt)} USDG`} />
        <Fact term={t("premium")} value={`${f.usdg(position.premium)} USDG`} />
        <Fact
          term={t("leverage")}
          value={position.equity > 0n ? `${leverage.format(Number(position.leverage) / 10_000)}×` : "—"}
        />
        <Fact term={t("requirement")} value={f.usd(position.requirement)} />
      </dl>
      {position.premium > 0n && position.debt === 0n && (
        <p className="text-[13px] text-[var(--warning)]">
          {t("premiumOnly")}
        </p>
      )}
      {position.stocks.map((s) => (
        <StockRow key={s.asset} stock={s} market={state.markets.find((m) => m.asset === s.asset)} position={position} />
      ))}
    </div>
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

type LendAction = "lend" | "unlend" | "recall";
const LEND_ACTIONS: readonly LendAction[] = ["lend", "unlend", "recall"];

function StockRow({ stock, market, position }: { stock: Stock; market: Market | undefined; position: Position }) {
  const t = useTranslations("account.positions");
  const f = useFormat();
  const explain = useExplain();
  const send = useSend();
  const [action, setAction] = useState<LendAction>("lend");
  const [typed, setTyped] = useState("");
  const amount = f.parseAmount(typed, 18);
  const low = market?.low ?? 0n;
  const haircut = amount && action === "lend" ? (amount * low * RECALL_HAIRCUT) / (10n ** 8n * 10_000n) : 0n;
  const owned = stock.recalls.length;
  const queuesDust =
    action === "recall" && amount !== undefined && amount < MIN_RECALL && !(owned === 0 && amount === stock.lent);
  return (
    <div className="border-t border-rule pt-4">
      <div className="flex flex-wrap items-baseline justify-between gap-2">
        <b className="font-semibold text-strong">{stock.asset}</b>
        <span className="text-[13px] text-muted num">
          {t.rich("liquidation", {
            value: stock.liquidationPrice === 0n ? t("noLiquidation") : `$${f.price(stock.liquidationPrice)}`,
            low: `$${f.price(low)}`,
            loss: (c) => <span className="text-[var(--liquidation)]">{c}</span>,
          })}
        </span>
      </div>
      <dl className="mt-2 grid grid-cols-2 gap-x-6 gap-y-2 text-[13px] sm:grid-cols-4 num">
        <Fact term={t("held")} value={f.tokens(stock.held)} />
        <Fact term={t("lent")} value={f.tokens(stock.lent)} />
        <Fact term={t("sellable")} value={f.tokens(stock.sellable)} />
        <Fact
          term={position.name === "Cross" ? t("inBaskets") : t("recalled")}
          value={f.tokens(position.name === "Cross" ? stock.inBaskets : stock.claim)}
        />
      </dl>
      {market?.lending && (
        <form
          className="mt-4 flex flex-wrap items-center gap-2"
          onSubmit={(e) => {
            e.preventDefault();
            if (amount !== undefined)
              send.send({ kind: action, position: position.id, token: stock.token, amount });
          }}
        >
          <div className="choice" role="group" aria-label={t("lendingAria", { asset: stock.asset })}>
            {LEND_ACTIONS.map((a) => (
              <button key={a} type="button" aria-pressed={a === action} onClick={() => setAction(a)}>
                {t(`verb.${a}`)}
              </button>
            ))}
          </div>
          <input
            className="amount !w-[6ch] !text-[19px]"
            inputMode="decimal"
            placeholder="0"
            aria-label={t(`amountAria.${action}`, { asset: stock.asset })}
            value={typed}
            onChange={(e) => setTyped(f.normalize(e.target.value))}
          />
          <button type="submit" className="btn btn-quiet" disabled={send.pending || amount === undefined || queuesDust}>
            {send.pending ? t("waiting") : t(`verb.${action}`)}
          </button>
          {stock.recalls.length > 0 && (
            <button
              type="button"
              className="btn btn-quiet"
              disabled={send.pending}
              onClick={() => send.send({ kind: "settle", position: position.id, token: stock.token })}
            >
              {t("settle")}
            </button>
          )}
        </form>
      )}
      <ul className="notes mt-3 max-w-[62ch]">
        {action === "lend" && amount !== undefined && (
          <li>
            {t("lendEquity", { equity: f.usd(position.equity - haircut), indebt: position.debt > 0n ? "yes" : "no" })}
          </li>
        )}
        {action === "unlend" && <li>{t("unlendNote")}</li>}
        {action === "recall" && (
          <li>
            {t("recallNote", { max: OWNER_RECALLS, owned })}
          </li>
        )}
        {queuesDust && <li data-tone="warning">{t("dust")}</li>}
        {owned >= OWNER_RECALLS && action === "recall" && (
          <li data-tone="warning">{t("recallsFull", { max: OWNER_RECALLS, asset: stock.asset })}</li>
        )}
        {stock.lent > 0n && (
          <li>{t("lentWarning")}</li>
        )}
      </ul>
      {stock.recalls.length > 0 && <Tickets asset={stock.asset} ids={stock.recalls} />}
      {send.error && (
        <p role="alert" className="mt-3 max-w-[60ch] text-[13px] text-[var(--error)]">
          {explain(send.error)}
        </p>
      )}
    </div>
  );
}

function Tickets({ asset, ids }: { asset: string; ids: readonly bigint[] }) {
  const t = useTranslations("account.positions");
  const f = useFormat();
  const deployments = useDeployments();
  const client = usePublicClient();
  const now = useMinute();
  const queries = useQueries({
    queries: ids.map((id) => ({
      queryKey: ["ticket", deployments.chainId, asset, id.toString()],
      queryFn: () => stockLending.ticket(client!, deployments, asset, id),
      enabled: client !== undefined,
    })),
  });
  return (
    <ol aria-label={t("recallsAria", { asset })} className="mt-3 grid gap-1.5 text-[13px] num">
      {queries.map((q, i) => {
        if (!q.data) return <li key={i} className="skeleton h-4 w-56" />;
        const { start, end, taken, dueAt, claimable, head } = q.data;
        const id = ids[i]!;
        const turn = id - head;
        const hours = Number(dueAt) * 1000 > now ? Math.ceil((Number(dueAt) * 1000 - now) / 3_600_000) : 0;
        return (
          <li key={i} className="flex flex-wrap justify-between gap-2">
            <span className="text-strong">
              {t("recalledAmount", { amount: f.tokens(end - start - taken), asset })}
            </span>
            <span className="text-muted">
              {turn === 0n ? t("nextInTurn") : t("ahead", { turn: String(turn) })} ·{" "}
              {t("heldFor", { amount: f.tokens(claimable) })} · {hours > 0 ? t("noticeEnds", { hours }) : t("pastNotice")}
            </span>
          </li>
        );
      })}
    </ol>
  );
}
