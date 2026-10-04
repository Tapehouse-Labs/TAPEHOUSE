// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { stockLending } from "@tapehouse/sdk";
import { useQueries } from "@tanstack/react-query";
import { useState } from "react";
import type { Hex } from "viem";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import type { AccountState, Market, Position, Stock } from "@/lib/accountRead";
import { useSend } from "@/lib/accountHook";
import { parseAmount, tokens, usd, usdg } from "@/lib/amounts";
import { explain } from "@/lib/errors";
import { price } from "@/lib/format";
import { useMinute } from "@/lib/hooks";

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
  if (state.positions.length === 0) return null;
  return (
    <section aria-labelledby="positions-title" className="mt-14">
      <h2 id="positions-title" className="label">
        Positions
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
                {p.name !== "Cross" && <span className="text-[13px] text-muted"> isolated</span>}
                <span className="mt-0.5 block text-[12.5px] text-muted num">
                  {p.stocks.map((s) => `${tokens(s.held + s.lent + s.inBaskets)} ${s.asset}`).join(" · ") || "No Stock Tokens"}
                  {p.usdg > 0n && ` · ${usdg(p.usdg)} USDG`}
                  {p.weth > 0n && ` · ${tokens(p.weth)} WETH`}
                </span>
              </span>
              <span className="text-right num">
                <b className="font-semibold text-strong">{usd(p.equity)}</b> <span className="text-[12.5px] text-muted">equity</span>
                <span className="block text-[12.5px] text-muted">
                  {usd(p.requirement)} required
                  {p.debt + p.premium > 0n && ` · owes ${usdg(p.debt + p.premium)}`}
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
  return (
    <div className="grid gap-5 pb-6">
      <dl className="grid grid-cols-2 gap-x-6 gap-y-2 text-[13px] sm:grid-cols-4 num">
        <Fact term="Debt" value={`${usdg(position.debt)} USDG`} />
        <Fact term="Premium" value={`${usdg(position.premium)} USDG`} />
        <Fact term="Leverage" value={position.equity > 0n ? `${(Number(position.leverage) / 10_000).toFixed(2)}×` : "—"} />
        <Fact term="Requirement" value={usd(position.requirement)} />
      </dl>
      {position.premium > 0n && position.debt === 0n && (
        <p className="text-[13px] text-[var(--warning)]">
          Owing only premium still counts as debt: the guardian&apos;s pause stops this position&apos;s withdrawals.
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
const LEND_VERB: Record<LendAction, string> = { lend: "Lend", unlend: "Take back", recall: "Recall" };

function StockRow({ stock, market, position }: { stock: Stock; market: Market | undefined; position: Position }) {
  const send = useSend();
  const [action, setAction] = useState<LendAction>("lend");
  const [typed, setTyped] = useState("");
  const amount = parseAmount(typed, 18);
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
          liquidation{" "}
          <span className="text-[var(--liquidation)]">
            {stock.liquidationPrice === 0n ? "none" : `$${price(stock.liquidationPrice)}`}
          </span>{" "}
          · band low ${price(low)}
        </span>
      </div>
      <dl className="mt-2 grid grid-cols-2 gap-x-6 gap-y-2 text-[13px] sm:grid-cols-4 num">
        <Fact term="Held" value={tokens(stock.held)} />
        <Fact term="Lent" value={tokens(stock.lent)} />
        <Fact term="Sellable now" value={tokens(stock.sellable)} />
        <Fact term={position.name === "Cross" ? "In baskets" : "Recalled"} value={tokens(position.name === "Cross" ? stock.inBaskets : stock.claim)} />
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
          <div className="choice" role="group" aria-label={`${stock.asset} lending`}>
            {(Object.keys(LEND_VERB) as LendAction[]).map((a) => (
              <button key={a} type="button" aria-pressed={a === action} onClick={() => setAction(a)}>
                {LEND_VERB[a]}
              </button>
            ))}
          </div>
          <input
            className="amount !w-[6ch] !text-[19px]"
            inputMode="decimal"
            placeholder="0"
            aria-label={`${stock.asset} to ${LEND_VERB[action].toLowerCase()}`}
            value={typed}
            onChange={(e) => setTyped(e.target.value)}
          />
          <button type="submit" className="btn btn-quiet" disabled={send.pending || amount === undefined || queuesDust}>
            {send.pending ? "Waiting…" : LEND_VERB[action]}
          </button>
          {stock.recalls.length > 0 && (
            <button
              type="button"
              className="btn btn-quiet"
              disabled={send.pending}
              onClick={() => send.send({ kind: "settle", position: position.id, token: stock.token })}
            >
              Settle what came back
            </button>
          )}
        </form>
      )}
      <ul className="notes mt-3 max-w-[62ch]">
        {action === "lend" && amount !== undefined && (
          <li>
            Your equity would be {usd(position.equity - haircut)}: a lent token counts 5% less in it, while the engine
            still counts it whole.{position.debt > 0n && " A position in debt lends only while it meets its requirement."}
          </li>
        )}
        {action === "unlend" && <li>Takes back only what the vault can return now: its queue&apos;s turn and its free tokens.</li>}
        {action === "recall" && (
          <li>
            What the vault holds free comes back at once; the rest waits at the end of its queue, returned within a day or
            bought in after it. You may hold {OWNER_RECALLS} recalls open at once; you have {owned}.
          </li>
        )}
        {queuesDust && <li data-tone="warning">A recall must queue at least a hundredth of a token.</li>}
        {owned >= OWNER_RECALLS && action === "recall" && (
          <li data-tone="warning">You have {OWNER_RECALLS} recalls of {stock.asset} open, the most an owner may.</li>
        )}
        {stock.lent > 0n && (
          <li>If this position falls short, anyone may recall all it lent through the liquidator.</li>
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
  const deployments = useDeployments();
  const client = usePublicClient();
  const now = useMinute();
  const tickets = useQueries({
    queries: ids.map((id) => ({
      queryKey: ["ticket", deployments.chainId, asset, id.toString()],
      queryFn: () => stockLending.ticket(client!, deployments, asset, id),
      enabled: client !== undefined,
    })),
  });
  return (
    <ol aria-label={`${asset} recalls`} className="mt-3 grid gap-1.5 text-[13px] num">
      {tickets.map((t, i) => {
        if (!t.data) return <li key={i} className="skeleton h-4 w-56" />;
        const { start, end, taken, dueAt, claimable, head } = t.data;
        const id = ids[i]!;
        const turn = id - head;
        const hours = Number(dueAt) * 1000 > now ? Math.ceil((Number(dueAt) * 1000 - now) / 3_600_000) : 0;
        return (
          <li key={i} className="flex flex-wrap justify-between gap-2">
            <span className="text-strong">
              {tokens(end - start - taken)} {asset} recalled
            </span>
            <span className="text-muted">
              {turn === 0n ? "next in turn" : `${turn} ahead of it`} · {tokens(claimable)} held for it ·{" "}
              {hours > 0 ? `notice ends in ${hours}H` : "past its notice, anyone may buy it in"}
            </span>
          </li>
        );
      })}
    </ol>
  );
}
