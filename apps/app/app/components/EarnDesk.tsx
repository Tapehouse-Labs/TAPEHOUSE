// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { useState } from "react";
import type { Address, Hex } from "viem";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import type { AccountState } from "@/lib/accountRead";
import { PHASE, useAccountState, useMarginAccount, useSend, useWalletBalances } from "@/lib/accountHook";
import { cooldown, span } from "@/lib/earn";
import { claimable, type EarnState, type Holding } from "@/lib/earnRead";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import type { Intent, Vault } from "@/lib/intent";
import { msg } from "@/lib/messages";
import { EarnView, type OpenRow } from "./Earn";
import { SmartAccount } from "./SmartAccount";

const EMPTY: readonly Address[] = [];

type PanelContext = { sponsored: boolean; usdg: bigint; account: Address | undefined };

/** The Earn surface for a connected account: the ledger with its smart account, positions and composers. */
export function EarnDesk({ row }: { row: OpenRow }) {
  const deployments = useDeployments();
  const { account, owner, smart } = useMarginAccount();
  const positions = useAccountState(Object.keys(deployments.stockLending).length > 0 ? account : undefined, owner);
  const balances = useWalletBalances(owner, EMPTY);
  const context: PanelContext = {
    sponsored: smart.enabled && (smart.standing?.left ?? 0n) > 0n,
    usdg: balances.data?.tokens[deployments.tokens.USDG as Address] ?? 0n,
    account,
  };
  return (
    <EarnView
      row={row}
      holder={account}
      positions={positions.data}
      account={<SmartAccount />}
      actions={{
        vault: (vault, holding) => <VaultComposer vault={vault} holding={holding} context={context} />,
        backstop: (backstop, timestamp) => <BackstopActions backstop={backstop} timestamp={timestamp} context={context} />,
        lend: (lending) => <LendActions lending={lending} positions={positions.data} context={context} />,
      }}
    />
  );
}

type BackstopView = "deposit" | "leave" | "gains";

function BackstopActions({
  backstop,
  timestamp,
  context,
}: {
  backstop: NonNullable<EarnState["backstop"]>;
  timestamp: bigint;
  context: PanelContext;
}) {
  const t = useTranslations("earn.desk");
  const f = useFormat();
  const deployments = useDeployments();
  const client = usePublicClient();
  const [view, setView] = useState<BackstopView>("deposit");
  const { state: b, depositor, tokens: stocks } = backstop;
  const { account } = context;
  if (!depositor || !client || !account) return null;
  return (
    <div>
      <Choice
        options={{ deposit: t("deposit"), leave: t("leave"), gains: t("gains") }}
        value={view}
        onChange={setView}
        label={t("backstopAction")}
      />
      {view === "deposit" && (
        <AmountForm
          key="deposit"
          verb="deposit"
          decimals={6}
          max={least(context.usdg, depositor.maxDeposit)}
          sponsored={context.sponsored}
          intent={(amount) => ({ kind: "vaultDeposit", vault: "backstop", amount })}
        />
      )}
      {view === "leave" && <Leave depositor={depositor} timestamp={timestamp} window={b.terms.window} context={context} />}
      {view === "gains" && (
        <div className="mt-5 grid gap-3 text-[14px]">
          {stocks.map(({ asset, token }) => {
            const gained = depositor.gains[token] ?? 0n;
            return (
              <div key={token} className="flex flex-wrap items-center gap-3">
                <span className="min-w-[9rem] text-strong">
                  {f.tokens(gained)} {asset}
                </span>
                <SendButton
                  label={t("claim", { asset })}
                  quiet
                  disabled={gained === 0n}
                  sponsored={context.sponsored}
                  intent={async () => ({
                    kind: "claimGains",
                    token,
                    amount: await claimable(client, deployments, account, token),
                  })}
                />
              </div>
            );
          })}
          <div className="flex flex-wrap items-center gap-3 border-t border-rule pt-3">
            <span className="min-w-[9rem] text-strong">{t("premium", { amount: f.usdg(b.premium.waiting) })}</span>
            <SendButton
              label={t("claimPremium")}
              quiet
              disabled={b.premium.waiting === 0n}
              sponsored={context.sponsored}
              intent={{ kind: "claimPremium" }}
            />
          </div>
        </div>
      )}
    </div>
  );
}

function Leave({
  depositor,
  timestamp,
  window,
  context,
}: {
  depositor: NonNullable<NonNullable<EarnState["backstop"]>["depositor"]>;
  timestamp: bigint;
  window: bigint;
  context: PanelContext;
}) {
  const t = useTranslations("earn.desk");
  const say = useMsg();
  const f = useFormat();
  const phase = cooldown(depositor.cooldown, timestamp);
  if (depositor.shares === 0n) return <p className="mt-5 text-[14px] text-body">{t("noShares")}</p>;
  if (phase.kind === "window" && depositor.maxRedeem > 0n)
    return (
      <>
        <p className="mt-5 text-[14px] text-body">{say(msg("earn.cooldown.windowOpen", { closes: span(phase.closesIn) }))}</p>
        <AmountForm
          key="redeem"
          verb="redeem"
          decimals={6}
          max={depositor.maxWithdraw}
          sponsored={context.sponsored}
          intent={(assets, all) =>
            all
              ? { kind: "vaultWithdraw", vault: "backstop", assets, shares: depositor.maxRedeem }
              : { kind: "vaultWithdraw", vault: "backstop", assets }
          }
        />
      </>
    );
  return (
    <div className="mt-5 grid gap-3 text-[14px] text-body">
      <p className="max-w-[60ch]">
        {phase.kind === "cooling"
          ? say(msg("earn.cooldown.cooling", { opens: span(phase.opensIn), window: span(window) }))
          : phase.kind === "window"
            ? say(
                msg(depositor.cooldown?.shares === 0n ? "earn.cooldown.windowSpent" : "earn.cooldown.windowShut", {
                  closes: span(phase.closesIn),
                }),
              )
            : say(msg("earn.cooldown.start", { assets: f.usdg(depositor.assets), window: span(window) }))}
      </p>
      <div>
        <SendButton
          label={phase.kind === "none" || phase.kind === "lapsed" ? t("startCooldown") : t("startAgain")}
          quiet={phase.kind === "cooling" || phase.kind === "window"}
          sponsored={context.sponsored}
          intent={{ kind: "cooldown" }}
        />
      </div>
    </div>
  );
}

type LendAction = "lend" | "unlend";

function LendActions({
  lending,
  positions,
  context,
}: {
  lending: EarnState["lending"][number];
  positions: AccountState | undefined;
  context: PanelContext;
}) {
  const t = useTranslations("earn.desk");
  const { asset, token } = lending;
  const [action, setAction] = useState<LendAction>("lend");
  const holders = (positions?.positions ?? []).filter((p) =>
    p.stocks.some((s) => s.asset === asset && (action === "lend" ? s.held > 0n : s.lent > 0n)),
  );
  const [chosen, setChosen] = useState<Hex>();
  const position = holders.find((p) => p.id === chosen) ?? holders[0];
  const stock = position?.stocks.find((s) => s.asset === asset);
  return (
    <div>
      <Choice
        options={{ lend: t("verb.lend"), unlend: t("verb.unlend") }}
        value={action}
        onChange={setAction}
        label={t("lendingAction")}
      />
      {position && stock ? (
        <AmountForm
          key={`${action}-${position.id}`}
          verb={action}
          unit={asset}
          decimals={18}
          max={action === "lend" ? stock.held : stock.lent}
          sponsored={context.sponsored}
          intent={(amount) => ({ kind: action, position: position.id, token, amount })}
          from={
            holders.length > 1 ? (
              <select
                className="picker"
                aria-label={t("position")}
                value={position.id}
                onChange={(e) => setChosen(e.target.value as Hex)}
              >
                {holders.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name === "Cross" ? t("cross") : t("isolated", { name: p.name })}
                  </option>
                ))}
              </select>
            ) : (
              <span className="text-body">{position.name === "Cross" ? t("fromCross") : t("fromIsolated", { name: position.name })}</span>
            )
          }
        />
      ) : (
        <p className="mt-5 text-[14px] text-body">
          {action === "lend" ? t("lendEmpty", { asset }) : t("unlendEmpty", { asset })}{" "}
          <Link href="/" prefetch={false} className="text-strong underline underline-offset-4">
            {t("openAccount")}
          </Link>
        </p>
      )}
      <p className="mt-3 text-[13px] text-muted">
        {t.rich("recalls", {
          link: (positions) => (
            <Link href="/" prefetch={false} className="underline underline-offset-4">
              {positions}
            </Link>
          ),
        })}
      </p>
    </div>
  );
}

type VaultAction = "deposit" | "withdraw";

function VaultComposer({ vault, holding, context }: { vault: Vault; holding: Holding | undefined; context: PanelContext }) {
  const t = useTranslations("earn.desk");
  const [action, setAction] = useState<VaultAction>("deposit");
  if (!holding) return null;
  return (
    <div>
      <Choice
        options={{ deposit: t("verb.deposit"), withdraw: t("verb.withdraw") }}
        value={action}
        onChange={setAction}
        label={t("action")}
      />
      <AmountForm
        key={action}
        verb={action}
        decimals={6}
        max={action === "deposit" ? least(context.usdg, holding.maxDeposit) : holding.maxWithdraw}
        sponsored={context.sponsored}
        intent={(amount, all) =>
          action === "deposit"
            ? { kind: "vaultDeposit", vault, amount }
            : all
              ? { kind: "vaultWithdraw", vault, assets: amount, shares: holding.maxRedeem }
              : { kind: "vaultWithdraw", vault, assets: amount }
        }
      />
    </div>
  );
}

function AmountForm({
  verb,
  unit = "USDG",
  decimals,
  max,
  sponsored,
  intent,
  from,
}: {
  verb: "deposit" | "withdraw" | "redeem" | "lend" | "unlend";
  unit?: string;
  decimals: number;
  max: bigint;
  sponsored: boolean;
  intent: (amount: bigint, all: boolean) => Intent;
  from?: React.ReactNode;
}) {
  const t = useTranslations("earn.desk");
  const f = useFormat();
  const say = useMsg();
  const send = useSend();
  const [typed, setTyped] = useState("");
  const amount = f.parseAmount(typed, decimals);
  const over = amount !== undefined && amount > max;
  return (
    <>
      <form
        className="mt-5 flex flex-wrap items-center gap-x-3 gap-y-3 text-[19px] text-strong"
        onSubmit={(e) => {
          e.preventDefault();
          if (amount !== undefined && !over) send.send(intent(amount, amount === max), { onSuccess: () => setTyped("") });
        }}
      >
        <span>{t(`verb.${verb}`)}</span>
        <input
          className="amount"
          inputMode="decimal"
          autoComplete="off"
          placeholder="0"
          aria-label={t(`amountTo.${verb}`)}
          aria-invalid={typed !== "" && (amount === undefined || over)}
          value={typed}
          onChange={(e) => {
            setTyped(f.normalize(e.target.value));
            if (!send.pending) send.reset();
          }}
        />
        <span>{unit}</span>
        {from}
        <button type="submit" className="btn ml-auto" disabled={send.pending || amount === undefined || over}>
          {send.phase ? say(PHASE[send.phase]) : sponsored ? t("free", { label: t(`verb.${verb}`) }) : t(`verb.${verb}`)}
        </button>
      </form>
      <p className="mt-2 text-[13px] text-muted num">
        <button type="button" className="underline-offset-4 hover:underline" onClick={() => setTyped(f.unparse(max, decimals))}>
          {t("max", { amount: decimals === 6 ? f.usdg(max) : f.tokens(max) })}
        </button>
      </p>
      <Outcome send={send} />
    </>
  );
}

function SendButton({
  label,
  intent,
  sponsored,
  quiet,
  disabled,
}: {
  label: string;
  intent: Intent | (() => Promise<Intent>);
  sponsored: boolean;
  quiet?: boolean;
  disabled?: boolean;
}) {
  const t = useTranslations("earn.desk");
  const say = useMsg();
  const send = useSend();
  return (
    <>
      <button
        type="button"
        className={quiet ? "btn btn-quiet" : "btn"}
        disabled={disabled || send.pending}
        onClick={() => send.send(intent)}
      >
        {send.phase ? say(PHASE[send.phase]) : sponsored ? t("free", { label }) : label}
      </button>
      <Outcome send={send} />
    </>
  );
}

function Outcome({ send }: { send: ReturnType<typeof useSend> }) {
  const t = useTranslations("earn.desk");
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

function Choice<value extends string>({
  options,
  value,
  onChange,
  label,
}: {
  options: Record<value, string>;
  value: value;
  onChange: (value: value) => void;
  label: string;
}) {
  return (
    <div className="choice mt-6" role="group" aria-label={label}>
      {(Object.keys(options) as value[]).map((option) => (
        <button key={option} type="button" aria-pressed={option === value} onClick={() => onChange(option)}>
          {options[option]}
        </button>
      ))}
    </div>
  );
}

function least(a: bigint, b: bigint): bigint {
  return a < b ? a : b;
}
