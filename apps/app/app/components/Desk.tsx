// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { CROSS, toBytes32 } from "@tapehouse/sdk";
import { useTranslations } from "next-intl";
import { useMemo, useState } from "react";
import { type Address, type Hex, zeroAddress } from "viem";
import { useDeployments } from "@/app/providers";
import { type AccountState, leadAsset, type Position, positionName } from "@/lib/accountRead";
import { PHASE, useAccountState, useLiquidationPreview, useMarginAccount, useSend, useWalletBalances } from "@/lib/accountHook";
import type { Formatters } from "@/lib/format";
import { useMinute } from "@/lib/hooks";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import type { Intent } from "@/lib/intent";
import { msg } from "@/lib/messages";
import { borrowNotes, depositNotes, issuerNotes, type Note, positionNotes, weekendAhead } from "@/lib/notes";
import { capacity, regimeWord, week } from "@/lib/week";
import { Positions } from "./Positions";
import { SmartAccount } from "./SmartAccount";
import { Week } from "./Week";

type Action = "deposit" | "borrow" | "repay" | "withdraw";
const ACTIONS: readonly Action[] = ["deposit", "borrow", "repay", "withdraw"];

/** The connected account: its liquidation price at the crossing, the composer, its week down the stem, its positions. */
export function Desk() {
  const t = useTranslations("account.desk");
  const f = useFormat();
  const say = useMsg();
  const explain = useExplain();
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
  const amount = f.parseAmount(typed, decimals);
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
    ...(action === "borrow" && position ? borrowNotes(f, data, position, borrowing) : []),
    ...(action === "deposit" && chosen ? depositNotes(f, data, chosen.value, account, owner, sponsored) : []),
    ...(action === "repay" ? [{ message: msg("account.desk.notes.repayOrder") }] : []),
    ...(action === "withdraw" && chosen && chosen.value !== "ETH" ? withdrawNotes(data, position, chosen.value, account, owner) : []),
    ...(position ? positionNotes(f, data, position) : []),
  ];

  const figure = preview.data ?? lead?.liquidationPrice;
  const regime = position?.regime ?? 2;

  return (
    <>
      <div className="relative pt-14">
        <span className="bar" aria-hidden="true" />
        <h1 className="label">
          {regime === 0 ? t("borrowingStops") : t("liquidationPrice")}
          {borrowing > 0n && ` ${t("withLoan")}`}
          {lead && ` · ${t(position?.name === "Cross" ? "leadCross" : "leadIsolated", { asset: lead.asset })}`}
        </h1>
        {lead && figure !== undefined ? (
          <p key={String(figure)} className="figure num swap mt-3" data-live="" aria-live="polite">
            {figure === 0n ? t("none") : `$${f.price(figure)}`}
          </p>
        ) : (
          <p className="figure mt-3" data-empty>
            {t("noneYet")}
          </p>
        )}
        <p className="mt-3 text-[14px] text-muted num">
          {lead ? (
            <>
              {borrowing > 0n && (
                <>
                  {t("without", { value: lead.liquidationPrice === 0n ? t("none") : `$${f.price(lead.liquidationPrice)}` })} ·{" "}
                </>
              )}
              {t("bandLow", { asset: lead.asset, price: `$${f.price(data.markets.find((m) => m.asset === lead.asset)?.low ?? 0n)}` })}
              {figure !== undefined && figure > 0n && (
                <> · {t("below", { drop: dropBelow(f, figure, data.markets.find((m) => m.asset === lead.asset)?.low ?? 0n) })}</>
              )}
              <span className="text-faint"> · {say(regimeWord(regime))}</span>
            </>
          ) : (
            t("empty")
          )}
        </p>
        <p className="label mt-2 num">{t("block", { block: Number(data.blockNumber) })}</p>
      </div>

      <SmartAccount />

      <section aria-label={t("actAria")} className="mt-8">
        <div className="choice" role="group" aria-label={t("actionAria")}>
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
              {t(`verb.${a}`)}
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
          <span>{t(`verb.${action}`)}</span>
          <input
            className="amount"
            inputMode="decimal"
            autoComplete="off"
            placeholder="0"
            aria-label={t(`amountAria.${action}`)}
            aria-invalid={typed !== "" && (amount === undefined || (max !== undefined && amount > max))}
            value={typed}
            onChange={(e) => {
              setTyped(f.normalize(e.target.value));
              if (!send.pending) send.reset();
            }}
          />
          {tokenChoices.length > 1 ? (
            <select
              className="picker"
              aria-label={t("tokenAria")}
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
          <span className="text-body">{action === "deposit" || action === "repay" ? t("into") : t("from")}</span>
          <select
            className="picker"
            aria-label={t("positionAria")}
            value={selected}
            onChange={(e) => setSelected(e.target.value as Hex)}
          >
            {positionChoices(data, action).map((id) => (
              <option key={id} value={id}>
                {id === CROSS ? positionName(id, data.markets) : t("isolatedOption", { name: positionName(id, data.markets) })}
              </option>
            ))}
          </select>
          <button
            type="submit"
            className="btn ml-auto"
            disabled={send.pending || !intent || (max !== undefined && amount !== undefined && amount > max)}
          >
            {send.phase
              ? say(PHASE[send.phase])
              : sponsored && !(action === "deposit" && chosen?.value === "ETH")
                ? t("free", { verb: t(`verb.${action}`) })
                : t(`verb.${action}`)}
          </button>
        </form>
        <p className="mt-2 text-[13px] text-muted num">
          {max !== undefined && (
            <button type="button" className="underline-offset-4 hover:underline" onClick={() => setTyped(f.unparse(max, decimals))}>
              {action === "borrow"
                ? t("upTo", { amount: f.usdg(max) })
                : t("max", { amount: decimals === 6 ? f.usdg(max) : f.tokens(max) })}
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
            {t("done")}
          </p>
        )}
      </section>

      {position && (position.open > 0n || position.requirement > 0n) && (
        <Week
          label={t("weekAria", { name: position.name })}
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
        <ul aria-label={t("notesAria")} className="notes mt-8 max-w-[62ch]">
          {notes.map((n) => {
            const text = say(n.message);
            return (
              <li key={text} data-tone={n.tone}>
                {text}
              </li>
            );
          })}
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
    lendable: state.lendable,
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
    notes.push({ message: msg("account.desk.notes.withdrawDebt") });
    if (state.limits.borrowingPaused)
      notes.push({ message: msg("account.desk.notes.withdrawPaused"), tone: "warning" });
  }
  const market = state.markets.find((m) => m.token === token);
  if (market?.issuer) notes.push(...issuerNotes(market.asset, market.issuer, account, owner));
  return notes;
}

function dropBelow(f: Formatters, price: bigint, low: bigint): string {
  if (low === 0n || price >= low) return `0${f.decimal}0%`;
  const tenths = ((low - price) * 1_000n) / low;
  return `${tenths / 10n}${f.decimal}${tenths % 10n}%`;
}
