// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { CROSS, toBytes32 } from "@tapehouse/sdk";
import { useMemo, useState } from "react";
import { type Address, type Hex, zeroAddress } from "viem";
import { useDeployments } from "@/app/providers";
import { type AccountState, leadAsset, type Position, positionName } from "@/lib/accountRead";
import { PHASE, useAccountState, useLiquidationPreview, useMarginAccount, useSend, useWalletBalances } from "@/lib/accountHook";
import { parseAmount, tokens, unparse, usdg } from "@/lib/amounts";
import { explain } from "@/lib/errors";
import { price } from "@/lib/format";
import { useMinute } from "@/lib/hooks";
import type { Intent } from "@/lib/intent";
import { borrowNotes, depositNotes, issuerNotes, type Note, positionNotes, weekendAhead } from "@/lib/notes";
import { capacity, regimeWord, week } from "@/lib/week";
import { Positions } from "./Positions";
import { SmartAccount } from "./SmartAccount";
import { Week } from "./Week";

type Action = "deposit" | "borrow" | "repay" | "withdraw";
const ACTIONS: readonly Action[] = ["deposit", "borrow", "repay", "withdraw"];
const VERB: Record<Action, string> = { deposit: "Deposit", borrow: "Borrow", repay: "Repay", withdraw: "Withdraw" };

/** The connected account: its liquidation price at the crossing, the composer, its week down the stem, its positions. */
export function Desk() {
  const deployments = useDeployments();
  const { account, owner, smart } = useMarginAccount();
  const state = useAccountState(account, owner);
  const [selected, setSelected] = useState<Hex>(CROSS);
  const [action, setAction] = useState<Action>("deposit");
  const [typed, setTyped] = useState("");
  const [token, setToken] = useState<Address | "ETH">();
  const now = useMinute();
  const send = useSend();
  const stockTokens = useMemo(
    () => (state.data?.markets ?? []).map((m) => m.token).filter((t) => t !== zeroAddress),
    [state.data],
  );
  const balances = useWalletBalances(owner, stockTokens);

  const data = state.data;
  const position = data?.positions.find((p) => p.id === selected);
  const lead = position && data ? leadAsset(position, data.markets) : undefined;
  const tokenChoices = data ? choices(deployments, data, action, selected, position) : [];
  const chosen = tokenChoices.find((c) => c.value === token) ?? tokenChoices[0];
  const decimals = action === "borrow" || action === "repay" || chosen?.value === deployments.tokens.USDG ? 6 : 18;
  const amount = parseAmount(typed, decimals);
  const borrowing = action === "borrow" ? (amount ?? 0n) : 0n;
  const preview = useLiquidationPreview(account, position?.id, lead?.asset, borrowing);
  const sponsored = smart.enabled && (smart.standing?.left ?? 0n) > 0n;

  if (!account || !owner) return null;
  if (!data) {
    return (
      <div aria-busy="true" className="mt-8">
        {state.error ? (
          <p role="alert" className="text-[14px] text-[var(--error)]">
            {explain(state.error)}
          </p>
        ) : (
          <span className="skeleton block h-[18px] w-[18rem]" />
        )}
      </div>
    );
  }

  const free = position ? capacity({ ...capacityInput(data, position) }) : 0n;
  const max = maximum(action, chosen?.value, position, free, balances.data, deployments.tokens.USDG as Address);
  const equityAfter = position
    ? position.equity - borrowing * 10n ** 12n + (action === "repay" ? (amount ?? 0n) * 10n ** 12n : 0n)
    : 0n;
  const intent = buildIntent(action, selected, chosen?.value, amount);
  const notes: Note[] = [
    ...(action === "borrow" && position ? borrowNotes(data, position, borrowing) : []),
    ...(action === "deposit" && chosen ? depositNotes(data, chosen.value, account, owner, sponsored) : []),
    ...(action === "repay" ? [{ text: "A repayment pays the debt first, then the premium." }] : []),
    ...(action === "withdraw" && chosen && chosen.value !== "ETH" ? withdrawNotes(data, position, chosen.value, account, owner) : []),
    ...(position ? positionNotes(data, position) : []),
  ];

  const figure = preview.data ?? lead?.liquidationPrice;
  const regime = position?.regime ?? 2;

  return (
    <>
      <div className="relative pt-14">
        <span className="bar" aria-hidden="true" />
        <h1 className="label">
          {regime === 0 ? "Borrowing stops at" : "Liquidation price"}
          {borrowing > 0n && " with this loan"}
          {lead && ` · ${lead.asset}, ${position?.name === "Cross" ? "cross" : "isolated"}`}
        </h1>
        {lead && figure !== undefined ? (
          <p key={String(figure)} className="figure num swap mt-3" data-live="" aria-live="polite">
            {figure === 0n ? "None" : `$${price(figure)}`}
          </p>
        ) : (
          <p className="figure mt-3" data-empty>
            None yet
          </p>
        )}
        <p className="mt-3 text-[14px] text-muted num">
          {lead ? (
            <>
              {borrowing > 0n && (
                <>{lead.liquidationPrice === 0n ? "None" : `$${price(lead.liquidationPrice)}`} without it · </>
              )}
              {lead.asset}&apos;s band low edge ${price(data.markets.find((m) => m.asset === lead.asset)?.low ?? 0n)}
              {figure !== undefined && figure > 0n && (
                <> · {dropBelow(figure, data.markets.find((m) => m.asset === lead.asset)?.low ?? 0n)} below it</>
              )}
              <span className="text-faint"> · {regimeWord(regime)}</span>
            </>
          ) : (
            "Deposit Stock Tokens and your liquidation price appears here before you borrow."
          )}
        </p>
        <p className="label mt-2 num">Read at block {data.blockNumber.toLocaleString("en-US")}</p>
      </div>

      <SmartAccount />

      <section aria-label="Act on your account" className="mt-8">
        <div className="choice" role="group" aria-label="Action">
          {ACTIONS.map((a) => (
            <button
              key={a}
              type="button"
              aria-pressed={a === action}
              onClick={() => {
                setAction(a);
                setTyped("");
                setToken(undefined);
                send.reset();
              }}
            >
              {VERB[a]}
            </button>
          ))}
        </div>
        <form
          className="mt-5 flex flex-wrap items-center gap-x-3 gap-y-3 text-[19px] text-strong"
          onSubmit={(e) => {
            e.preventDefault();
            if (intent) send.send(intent);
          }}
        >
          <span>{VERB[action]}</span>
          <input
            className="amount"
            inputMode="decimal"
            autoComplete="off"
            placeholder="0"
            aria-label={`Amount to ${action}`}
            aria-invalid={typed !== "" && (amount === undefined || (max !== undefined && amount > max))}
            value={typed}
            onChange={(e) => {
              setTyped(e.target.value);
              if (!send.pending) send.reset();
            }}
          />
          {tokenChoices.length > 1 ? (
            <select
              className="picker"
              aria-label="Token"
              value={String(chosen?.value)}
              onChange={(e) => setToken(e.target.value as Address | "ETH")}
            >
              {tokenChoices.map((c) => (
                <option key={String(c.value)} value={String(c.value)}>
                  {c.label}
                </option>
              ))}
            </select>
          ) : (
            <span>{chosen?.label ?? "USDG"}</span>
          )}
          <span className="text-body">{action === "deposit" || action === "repay" ? "into" : "from"}</span>
          <select
            className="picker"
            aria-label="Position"
            value={selected}
            onChange={(e) => setSelected(e.target.value as Hex)}
          >
            {positionChoices(data, action).map((id) => (
              <option key={id} value={id}>
                {positionName(id, data.markets)}
                {id === CROSS ? "" : " isolated"}
              </option>
            ))}
          </select>
          <button
            type="submit"
            className="btn ml-auto"
            disabled={send.pending || !intent || (max !== undefined && amount !== undefined && amount > max)}
          >
            {send.phase ? PHASE[send.phase] : `${VERB[action]}${sponsored && !(action === "deposit" && chosen?.value === "ETH") ? ", free" : ""}`}
          </button>
        </form>
        <p className="mt-2 text-[13px] text-muted num">
          {max !== undefined && (
            <button type="button" className="underline-offset-4 hover:underline" onClick={() => setTyped(unparse(max, decimals))}>
              {action === "borrow" ? "Up to" : "Max"} {decimals === 6 ? usdg(max) : tokens(max)}
              {action === "borrow" && " USDG you can borrow now"}
            </button>
          )}
        </p>
        {send.error && (
          <p role="alert" className="mt-3 max-w-[60ch] text-[14px] text-[var(--error)]">
            {explain(send.error)}
          </p>
        )}
        {send.succeeded && !send.pending && (
          <p role="status" className="mt-3 text-[14px] text-[var(--gain)]">
            Done. Every figure below is read again from the chain.
          </p>
        )}
      </section>

      {position && (position.open > 0n || position.requirement > 0n) && (
        <Week
          label={`${position.name} position's requirement against its equity through the week`}
          equity={equityAfter}
          marks={week({
            regime: position.regime,
            current: position.requirement,
            open: position.open,
            closed: position.closed,
            gross: position.gross,
            weekendLeverage: data.limits.weekendLeverage,
            boundaryMs: data.session.boundaryMs,
            nowMs: now,
          })}
        />
      )}

      {notes.length > 0 && (
        <ul aria-label="Before you sign" className="notes mt-8 max-w-[62ch]">
          {notes.map((n) => (
            <li key={n.text} data-tone={n.tone}>
              {n.text}
            </li>
          ))}
        </ul>
      )}

      <Positions state={data} selected={selected} onSelect={setSelected} />
    </>
  );
}

function capacityInput(state: AccountState, position: Position) {
  return {
    equity: position.equity,
    requirement: position.requirement,
    gross: position.gross,
    weekendAhead: weekendAhead(state, position),
    weekendLeverage: state.limits.weekendLeverage,
    debt: state.limits.debt,
    debtCap: state.limits.debtCap,
    weekendDebtCap: state.limits.weekendDebtCap,
  };
}

type Choice = { value: Address | "ETH"; label: string };

function choices(
  deployments: ReturnType<typeof useDeployments>,
  state: AccountState,
  action: Action,
  selected: Hex,
  position: Position | undefined,
): Choice[] {
  if (action === "borrow" || action === "repay") return [{ value: deployments.tokens.USDG as Address, label: "USDG" }];
  const stocks = state.markets.filter(
    (m) => m.token !== zeroAddress && (selected === CROSS || toBytes32(m.asset) === selected),
  );
  const all: Choice[] = [
    ...stocks.map((m) => ({ value: m.token, label: m.asset })),
    { value: deployments.tokens.USDG as Address, label: "USDG" },
    ...(deployments.tokens.WETH ? [{ value: deployments.tokens.WETH as Address, label: "WETH" }] : []),
  ];
  if (action === "deposit") return deployments.tokens.WETH ? [...all, { value: "ETH", label: "ETH" }] : all;
  return all.filter((c) => {
    if (!position) return false;
    if (c.value === deployments.tokens.USDG) return position.usdg > 0n;
    if (c.value === deployments.tokens.WETH) return position.weth > 0n;
    return position.stocks.some((s) => s.token === c.value && s.held > 0n);
  });
}

function positionChoices(state: AccountState, action: Action): Hex[] {
  const held = state.positions.map((p) => p.id);
  if (action !== "deposit") return held.length > 0 ? held : [CROSS];
  const isolated = state.markets.filter((m) => m.token !== zeroAddress).map((m) => toBytes32(m.asset));
  return [CROSS, ...isolated];
}

function maximum(
  action: Action,
  token: Address | "ETH" | undefined,
  position: Position | undefined,
  free: bigint,
  wallet: { eth: bigint; tokens: Record<Address, bigint> } | undefined,
  usdgToken: Address,
): bigint | undefined {
  if (action === "borrow") return free;
  if (action === "repay") {
    if (!position) return 0n;
    const owed = position.debt + position.premium;
    const held = wallet?.tokens[usdgToken] ?? 0n;
    return owed < held ? owed : held;
  }
  if (!token) return undefined;
  if (action === "deposit") return token === "ETH" ? wallet?.eth : wallet?.tokens[token];
  if (!position) return 0n;
  if (token === usdgToken) return position.usdg;
  const stock = position.stocks.find((s) => s.token === token);
  return stock ? stock.held : position.weth;
}

function buildIntent(action: Action, position: Hex, token: Address | "ETH" | undefined, amount: bigint | undefined): Intent | undefined {
  if (amount === undefined) return undefined;
  if (action === "borrow" || action === "repay") return { kind: action, position, amount };
  if (!token) return undefined;
  if (action === "deposit") return { kind: "deposit", position, token, amount };
  return token === "ETH" ? undefined : { kind: "withdraw", position, token, amount };
}

function withdrawNotes(
  state: AccountState,
  position: Position | undefined,
  token: Address,
  account: Address,
  owner: Address,
): Note[] {
  const notes: Note[] = [];
  if (position && position.debt + position.premium > 0n) {
    notes.push({ text: "The position is in debt, so what it keeps must still meet its requirement." });
    if (state.limits.borrowingPaused)
      notes.push({ text: "The guardian has paused withdrawals from positions in debt.", tone: "warning" });
  }
  const market = state.markets.find((m) => m.token === token);
  if (market?.issuer) notes.push(...issuerNotes(market.asset, market.issuer, account, owner));
  return notes;
}

function dropBelow(price: bigint, low: bigint): string {
  if (low === 0n || price >= low) return "0.0%";
  const tenths = ((low - price) * 1_000n) / low;
  return `${tenths / 10n}.${tenths % 10n}%`;
}
