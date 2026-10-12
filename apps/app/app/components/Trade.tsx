// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { gapCover, shorts } from "@tapehouse/sdk";
import { useTranslations } from "next-intl";
import { lazy, type ReactNode, Suspense, useRef, useState } from "react";
import type { Address } from "viem";
import { useChains } from "wagmi";
import { useDeployments } from "@/app/providers";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import { useMinute, useWallet } from "@/lib/hooks";
import type { Intent } from "@/lib/intent";
import { msg } from "@/lib/messages";
import type { Note } from "@/lib/notes";
import {
  basketNotes,
  countdown,
  coverPurchaseNotes,
  DEFAULT_SLIPPAGE,
  feeTier,
  gapPercent,
  parseSlippage,
  saleFloor,
  shortNotes,
  shortStanding,
  withSlippage,
} from "@/lib/trade";
import { useBasketPreview, useCoverQuote, useHealthAfter, usePoolQuote, useTradeState } from "@/lib/tradeHook";
import type { BasketState, Market, TradeState } from "@/lib/tradeRead";
import { Wallets } from "./Wallets";

const TradeDesk = lazy(() => import("./TradeDesk").then((m) => ({ default: m.TradeDesk })));

/** What a connected account brings to a ticket: who trades, the wallet that owns it and what it holds, and the button that sends. */
export type Desk = {
  account: Address;
  owner: Address;
  balances: Record<Address, bigint>;
  send: (intent: Intent | undefined, label: string) => ReactNode;
};

export type Kind = "short" | "cover" | "basket";

/** The ticket shown and its setter. */
export type KindState = [Kind, (kind: Kind) => void];

/** A ticket opened from a line of the book: its kind, action, what it trades and the amount typed. */
type Preset = { kind: Kind; action: string; pick: string; typed: string; n: number };
type Open = (preset: Omit<Preset, "n">) => void;

/** The Trade surface: one ticket for a short, gap cover or a basket, and the account's trades beneath it. */
export function Trade() {
  const wallet = useWallet();
  const kind: KindState = useState<Kind>("short");
  const [desk, setDesk] = useState<{ desk: Desk; account: ReactNode }>();
  const connected = wallet.state === "connected";
  return (
    <>
      {connected && (
        <Suspense fallback={null}>
          <TradeDesk onDesk={setDesk} />
        </Suspense>
      )}
      <TradeView kind={kind} desk={connected ? desk?.desk : undefined} account={connected ? desk?.account : undefined} />
    </>
  );
}

/**
 * The ticket under the bar and the account's trades, read at one block. `kind` is the ticket shown; `desk` the
 * connected account, which the composers' chunk hands up once it loads, and `account` its smart account's block; a
 * visitor reads every line but signs nothing.
 */
export function TradeView({
  kind: [kind, setKind],
  desk,
  account,
}: {
  kind: KindState;
  desk?: Desk;
  account?: ReactNode;
}) {
  const t = useTranslations("trade");
  const f = useFormat();
  const explain = useExplain();
  const wallet = useWallet();
  const [chain] = useChains();
  const state = useTradeState(desk?.account);
  const data = state.data;
  const top = useRef<HTMLDivElement>(null);
  const [preset, setPreset] = useState<Preset>();
  const open: Open = (next) => {
    setKind(next.kind);
    setPreset({ ...next, n: (preset?.n ?? 0) + 1 });
    top.current?.scrollIntoView({ block: "start" });
  };
  const shown = preset?.kind === kind ? preset : undefined;
  const connect =
    wallet.state === "connected" ? null : (
      <>
        <p className="max-w-[44ch] text-[15px] text-body">
          {wallet.state === "wrong-network"
            ? t.rich("connect.wrongNetwork", { chain: chain.name, b: (c) => <span className="text-strong">{c}</span> })
            : t("connect.visitor")}
        </p>
        <div className="mt-4 flex min-h-[44px] flex-wrap items-center gap-3">
          {wallet.state === "wrong-network" ? (
            <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
              {wallet.pending ? t("connect.waiting") : t("connect.switch", { chain: chain.name })}
            </button>
          ) : (
            (wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />
          )}
        </div>
      </>
    );
  const sign = (intent: Intent | undefined, label: string) => (desk ? desk.send(intent, label) : connect);
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div ref={top} className="relative scroll-mt-24 pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">{t("page.title")}</h1>
          <Choice
            options={{ short: t("page.kind.short"), cover: t("page.kind.cover"), basket: t("page.kind.basket") }}
            value={kind}
            onChange={setKind}
            label={t("page.ticket")}
          />
          {!data ? (
            <div aria-busy="true" className="mt-6">
              {state.error ? (
                <p role="alert" className="text-[14px] text-[var(--error)]">
                  {explain(state.error)}
                </p>
              ) : (
                <span className="skeleton block h-[18px] w-[18rem]" />
              )}
            </div>
          ) : kind === "short" ? (
            <ShortTicket key={shown?.n} state={data} desk={desk} sign={sign} preset={shown} />
          ) : kind === "cover" ? (
            <CoverTicket state={data} sign={sign} />
          ) : (
            <BasketTicket key={shown?.n} state={data} desk={desk} sign={sign} preset={shown} />
          )}
          {data && <p className="label mt-3 num">{t("block", { block: new Intl.NumberFormat(f.locale).format(data.blockNumber) })}</p>}
        </div>
        {account}
        {data && desk && <Book state={data} desk={desk} open={open} />}
      </div>
    </main>
  );
}

type ShortAction = "sell" | "buyBack" | "marginIn" | "marginOut";

/** What a ticket's foot holds: the account's send button, or how a visitor connects. */
type Sign = (intent: Intent | undefined, label: string) => ReactNode;

function ShortTicket({ state, desk, sign, preset }: { state: TradeState; desk: Desk | undefined; sign: Sign; preset?: Preset }) {
  const t = useTranslations("trade");
  const f = useFormat();
  const say = useMsg();
  const explain = useExplain();
  const deployments = useDeployments();
  const markets = state.markets.filter((m) => m.fee > 0);
  const [chosen, setChosen] = useState(preset?.pick);
  const [action, setAction] = useState((preset?.action as ShortAction | undefined) ?? "sell");
  const [typed, setTyped] = useState(preset?.typed ?? "");
  const [slippage, setSlippage] = useState(f.normalize(DEFAULT_SLIPPAGE));
  const market = markets.find((m) => m.asset === chosen) ?? markets[0];
  const usdgUnits = action === "marginIn" || action === "marginOut";
  const amount = f.parseAmount(typed, usdgUnits ? 6 : 18);
  const bps = parseSlippage(f, slippage);
  const side = action === "buyBack" ? "buyBack" : "sell";
  const position = market?.short?.position;
  const quoted = action === "buyBack" && position && amount !== undefined && amount > position.debt ? position.debt : amount;
  const pool = usePoolQuote(side, market?.asset ?? "", usdgUnits ? undefined : quoted, bps);
  const standing = position && market ? shortStanding(position, market.epoch) : "none";
  const intent = market && amount !== undefined ? shortIntent(action, market.asset, amount, pool.data?.limit) : undefined;
  const call = intent && desk && intent.kind !== "shortDeposit" ? simulated(deployments, intent, desk) : undefined;
  const after = useHealthAfter(desk?.account, market?.asset ?? "", call);

  if (!deployments.tapehouse.ShortPositions || !market)
    return <p className="mt-6 text-[15px] text-body">{t("short.none")}</p>;
  const floor = saleFloor(market.quote, market.restricted);
  const health = market.short?.health;
  const full = action === "buyBack" && position && amount !== undefined && amount >= position.debt;
  const simulatedAfter = after.data && "equity" in after.data ? after.data : undefined;
  const reverts = after.data && "error" in after.data ? after.data.error : undefined;
  const marginAfter =
    health && amount !== undefined && action === "marginIn"
      ? { equity: health.equity + amount * 10n ** 12n, requirement: health.requirement }
      : simulatedAfter;
  const lines: [string, ReactNode][] = [];
  if (action === "sell" || action === "buyBack") {
    lines.push(
      [
        action === "sell" ? t("short.poolPays") : t("short.poolAsks"),
        pool.error ? t("short.noQuote") : pool.data ? `$${f.usdg(pool.data.quoted)}` : "—",
      ],
      [
        action === "sell" ? t("short.atLeast") : t("short.atMost"),
        <span key="limit" className="inline-flex items-center gap-2">
          <span>{pool.data && !pool.error ? `$${f.usdg(pool.data.limit)}` : "—"}</span>
          <span className="text-muted">{t("short.at")}</span>
          <input
            className="amount amount-small"
            inputMode="decimal"
            aria-label={t("short.slippageAria")}
            aria-invalid={bps === undefined}
            value={slippage}
            onChange={(e) => setSlippage(f.normalize(e.target.value))}
          />
          <span className="text-muted">%</span>
        </span>,
      ],
      action === "sell"
        ? [t("short.floor"), t("short.floorValue", { price: `$${f.price(floor.price)}`, floor: say(floor.words) })]
        : full
          ? [t("short.cost"), t("short.fullLimit")]
          : [t("short.cost"), t("short.highEdge", { price: `$${f.price(market.quote.high)}` })],
      [
        t("short.lendingFee"),
        market.borrowRate === undefined ? "—" : t("short.lendingFeeValue", { rate: f.percent(market.borrowRate), asset: market.asset }),
      ],
      [t("short.lenders"), market.borrowable === undefined ? "—" : `${f.tokens(market.borrowable)} ${market.asset}`],
      [t("short.feeTier"), feeTier(f, market.fee)],
    );
  } else {
    lines.push([action === "marginIn" ? t("short.source") : t("short.to"), t("short.wallet")]);
  }
  if (health)
    lines.push([t("short.marginNow"), t("short.against", { equity: f.usd(health.equity), requirement: f.usd(health.requirement) })]);
  if (desk)
    lines.push([
      t("short.marginAfter"),
      marginAfter
        ? t("short.against", { equity: f.usd(marginAfter.equity), requirement: f.usd(marginAfter.requirement) })
        : reverts
          ? t("short.revert")
          : "—",
    ]);
  const belowFloor =
    action === "sell" && amount !== undefined && pool.data && pool.data.quoted * 10n ** 20n < amount * floor.price;
  const notes: Note[] = [
    ...(standing === "earlier"
      ? [{ term: msg("trade.short.earlierTerm"), message: msg("trade.short.earlier", { held: f.usdg(position!.usdgHeld) }) }]
      : []),
    ...(standing === "deficit"
      ? [{ term: msg("trade.short.deficitTerm"), message: msg("trade.short.deficit"), tone: "warning" as const }]
      : []),
    ...shortNotes(f, { asset: market.asset, borrowRate: market.borrowRate ?? 0n, fee: market.fee, restricted: market.restricted }),
  ];
  const max =
    action === "marginIn"
      ? desk?.balances[deployments.tokens.USDG as Address]
      : action === "marginOut"
        ? position && position.usdgHeld > 0n
          ? position.usdgHeld
          : 0n
        : action === "buyBack" && position
          ? position.debt + position.debt / 1_000n + 1n
          : undefined;
  return (
    <>
      <dl className="ticket num" aria-label={t("page.ticket")}>
        <div className="ticket-head">
          <dt>
            {t("short.you")}{" "}
            <Verb
              options={{
                sell: t("short.verb.sell"),
                buyBack: t("short.verb.buyBack"),
                marginIn: t("short.verb.marginIn"),
                marginOut: t("short.verb.marginOut"),
              }}
              value={action}
              onChange={(a) => (setAction(a), setTyped(""))}
              label={t("short.action.aria")}
            />
          </dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label={t("short.amountAria", { action })}
              aria-invalid={typed !== "" && (amount === undefined || (max !== undefined && amount > max))}
              value={typed}
              onChange={(e) => setTyped(f.normalize(e.target.value))}
            />{" "}
            {action === "marginIn" ? t("short.marginIn") : action === "marginOut" ? t("short.marginOut") : ""}
            <select className="picker" aria-label={t("short.stockToken")} value={market.asset} onChange={(e) => setChosen(e.target.value)}>
              {markets.map((m) => (
                <option key={m.asset} value={m.asset}>
                  {m.asset}
                </option>
              ))}
            </select>
            {usdgUnits ? t("short.suffix") : ""}
          </dd>
        </div>
        {lines.map(([name, value]) => (
          <div key={name}>
            <dt>{name}</dt>
            <dd>{value}</dd>
          </div>
        ))}
      </dl>
      {max !== undefined && (
        <p className="mt-2 text-[13px] text-muted num">
          <button type="button" className="underline-offset-4 hover:underline" onClick={() => setTyped(f.unparse(max, usdgUnits ? 6 : 18))}>
            {action === "buyBack" ? t("short.maxBuyBack") : t("short.max", { amount: usdgUnits ? f.usdg(max) : f.tokens(max) })}
          </button>
        </p>
      )}
      <div className="mt-5">{sign(pool.error ? undefined : intent, action === "sell" ? t("short.submit") : t(`short.action.${action}`))}</div>
      {pool.error && <p className="mt-3 text-[14px] text-[var(--error)]">{explain(pool.error)}</p>}
      {belowFloor && (
        <p className="mt-3 max-w-[60ch] text-[14px] text-[var(--warning)]">
          {t("short.belowFloor", { floor: say(floor.words) })}
        </p>
      )}
      {reverts && !belowFloor && <p className="mt-3 max-w-[60ch] text-[14px] text-[var(--error)]">{explain(reverts)}</p>}
      <Notes notes={notes} />
    </>
  );
}

function shortIntent(action: ShortAction, asset: string, amount: bigint, limit: bigint | undefined): Intent | undefined {
  if (action === "marginIn") return { kind: "shortDeposit", asset, amount };
  if (action === "marginOut") return { kind: "shortWithdraw", asset, amount };
  if (limit === undefined) return undefined;
  return action === "sell"
    ? { kind: "sell", asset, amount, minProceeds: limit }
    : { kind: "buyBack", asset, amount, maxCost: limit };
}

function simulated(deployments: ReturnType<typeof useDeployments>, intent: Intent, { account, owner }: Desk) {
  switch (intent.kind) {
    case "sell":
      return shorts.sell(deployments, { asset: intent.asset, amount: intent.amount, minProceeds: intent.minProceeds, account });
    case "buyBack":
      return shorts.cover(deployments, { asset: intent.asset, amount: intent.amount, maxCost: intent.maxCost, account });
    case "shortWithdraw":
      return shorts.withdraw(deployments, { asset: intent.asset, amount: intent.amount, account, receiver: owner });
    default:
      return undefined;
  }
}

function CoverTicket({ state, sign }: { state: TradeState; sign: Sign }) {
  const t = useTranslations("trade");
  const f = useFormat();
  const say = useMsg();
  const explain = useExplain();
  const now = useMinute();
  const [chosen, setChosen] = useState<string>();
  const [notional, setNotional] = useState("");
  const [deductible, setDeductible] = useState("3");
  const [limit, setLimit] = useState("10");
  const markets = state.markets;
  const market = markets.find((m) => m.asset === chosen) ?? markets.find((m) => m.asset === "SPY") ?? markets[0];
  const amount = f.parseAmount(notional, 6);
  const deductibleBps = parseSlippage(f, deductible);
  const limitBps = parseSlippage(f, limit);
  const layer =
    amount !== undefined && deductibleBps !== undefined && limitBps !== undefined && limitBps > deductibleBps
      ? { notional: amount, deductibleBps, limitBps }
      : undefined;
  const quote = useCoverQuote(market?.asset ?? "", state.cover?.sales && layer ? layer : undefined);
  if (!state.cover) return <p className="mt-6 text-[15px] text-body">{t("cover.none")}</p>;
  const { sales, capacity } = state.cover;
  if (!market) return null;
  const maxPremium = quote.data && withSlippage(quote.data.premium, COVER_SLIPPAGE);
  const intent: Intent | undefined =
    layer && maxPremium !== undefined && !quote.error ? { kind: "buyCover", asset: market.asset, ...layer, maxPremium } : undefined;
  return (
    <>
      <dl className="ticket num" aria-label={t("page.ticket")}>
        <div className="ticket-head">
          <dt>{t("cover.you")}</dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label={t("cover.notionalAria")}
              aria-invalid={notional !== "" && amount === undefined}
              value={notional}
              onChange={(e) => setNotional(f.normalize(e.target.value))}
            />{" "}
            {t("cover.usdgOf")}{" "}
            <select className="picker" aria-label={t("short.stockToken")} value={market.asset} onChange={(e) => setChosen(e.target.value)}>
              {markets.map((m) => (
                <option key={m.asset} value={m.asset}>
                  {m.asset}
                </option>
              ))}
            </select>
          </dd>
        </div>
        <div>
          <dt>{t("cover.paysBetween")}</dt>
          <dd className="inline-flex items-center gap-2">
            <input
              className="amount amount-small"
              inputMode="decimal"
              aria-label={t("cover.deductibleAria")}
              value={deductible}
              onChange={(e) => setDeductible(f.normalize(e.target.value))}
            />
            <span className="text-muted">{t("cover.percentAnd")}</span>
            <input
              className="amount amount-small"
              inputMode="decimal"
              aria-label={t("cover.limitAria")}
              value={limit}
              onChange={(e) => setLimit(f.normalize(e.target.value))}
            />
            <span className="text-muted">{t("cover.percent")}</span>
          </dd>
        </div>
        <div>
          <dt>{t("cover.sales")}</dt>
          <dd data-tone={sales ? undefined : "warning"}>
            {sales ? t("cover.closeIn", { time: say(countdown(sales.endsMs, now)) }) : t("cover.closed")}
          </dd>
        </div>
        {sales && (
          <>
            <div>
              <dt>{t("cover.premium")}</dt>
              <dd>
                {quote.error ? t("short.noQuote") : quote.data ? `$${f.usdg(quote.data.premium)}` : "—"}
                {maxPremium !== undefined && !quote.error && <small>{t("cover.atMost", { max: `$${f.usdg(maxPremium)}` })}</small>}
              </dd>
            </div>
            <div>
              <dt>{t("cover.reserves")}</dt>
              <dd>
                {quote.data
                  ? t("cover.reserveValue", { reserve: f.usdg(quote.data.reserve), capacity: f.usdg(capacity) })
                  : t("cover.free", { capacity: f.usdg(capacity) })}
              </dd>
            </div>
            <div>
              <dt>{t("cover.gap")}</dt>
              <dd>
                {quote.data
                  ? quote.data.weekMove > 0n && quote.data.weekMove < 2n ** 64n - 1n
                    ? t("cover.gapFrom", { gap: gapPercent(f, quote.data.gap), week: gapPercent(f, quote.data.weekMove) })
                    : t("cover.gapValue", { gap: gapPercent(f, quote.data.gap) })
                  : "—"}
              </dd>
            </div>
            <div>
              <dt>{t("cover.minDeductible")}</dt>
              <dd>{quote.data ? f.bpsPercent(quote.data.minDeductible) : "—"}</dd>
            </div>
          </>
        )}
      </dl>
      <div className="mt-5">{sign(sales ? intent : undefined, t("cover.submit"))}</div>
      {quote.error && <p className="mt-3 text-[14px] text-[var(--error)]">{explain(quote.error)}</p>}
      <Notes notes={coverPurchaseNotes(market.asset)} />
    </>
  );
}

type BasketAction = "mint" | "redeem" | "unwrap";
const COVER_SLIPPAGE = 50n;
const MINT_SLIPPAGE = 10n;

function BasketTicket({ state, desk, sign, preset }: { state: TradeState; desk: Desk | undefined; sign: Sign; preset?: Preset }) {
  const t = useTranslations("trade");
  const f = useFormat();
  const [chosen, setChosen] = useState(preset?.pick);
  const [action, setAction] = useState((preset?.action as BasketAction | undefined) ?? "mint");
  const [typed, setTyped] = useState(preset?.typed ?? "");
  const basket = state.baskets.find((b) => b.key === chosen) ?? state.baskets[0];
  const shares = f.parseAmount(typed, 18);
  const preview = useBasketPreview(action === "mint" ? "mint" : "redeem", basket?.key ?? "", basket ? shares : undefined);
  if (!basket) return <p className="mt-6 text-[15px] text-body">{t("basket.none")}</p>;
  const amounts = preview.data;
  const maxAssets = action === "mint" ? amounts?.map((a) => withSlippage(a, MINT_SLIPPAGE)) : undefined;
  const short = maxAssets && desk ? basket.tokens.some((t, i) => (maxAssets[i] ?? 0n) > (desk.balances[t] ?? 0n)) : false;
  const intent: Intent | undefined =
    shares === undefined || short
      ? undefined
      : action === "mint"
        ? maxAssets && { kind: "basketMint", basket: basket.key, shares, tokens: basket.tokens, maxAssets }
        : { kind: action === "redeem" ? "basketRedeem" : "unwrap", basket: basket.key, shares };
  const words = (units: readonly bigint[] | undefined) =>
    units ? basket.assets.map((a, i) => `${f.tokens(units[i] ?? 0n)} ${a}`).join(" + ") : "—";
  const inBaskets = Object.entries(state.cross?.inBaskets ?? {})
    .filter(([, amount]) => amount > 0n)
    .map(([asset, amount]) => `${f.tokens(amount)} ${asset}`)
    .join(" + ");
  return (
    <>
      <dl className="ticket num" aria-label={t("page.ticket")}>
        <div className="ticket-head">
          <dt>
            {t("basket.you")}{" "}
            <Verb
              options={{ mint: t("basket.verb.mint"), redeem: t("basket.verb.redeem"), unwrap: t("basket.verb.unwrap") }}
              value={action}
              onChange={(a) => (setAction(a), setTyped(""))}
              label={t("basket.action.aria")}
            />
          </dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label={t("basket.sharesAria", { action })}
              aria-invalid={typed !== "" && (shares === undefined || (action !== "mint" && (basket.held ?? 0n) < shares))}
              value={typed}
              onChange={(e) => setTyped(f.normalize(e.target.value))}
            />{" "}
            {t("basket.sharesOf", { one: String(shares === 10n ** 18n) })}{" "}
            <select className="picker" aria-label={t("basket.aria")} value={basket.key} onChange={(e) => setChosen(e.target.value)}>
              {state.baskets.map((b) => (
                <option key={b.key} value={b.key}>
                  {b.key}
                </option>
              ))}
            </select>
          </dd>
        </div>
        <div>
          <dt>{action === "mint" ? t("basket.takes") : action === "redeem" ? t("basket.receives") : t("basket.unwrapHolds")}</dt>
          <dd data-tone={short ? "warning" : undefined}>
            {words(amounts)}
            {maxAssets && <small>{t("basket.atMost", { assets: words(maxAssets) })}</small>}
          </dd>
        </div>
        <div>
          <dt>{action === "mint" ? t("basket.into") : t("basket.source")}</dt>
          <dd>{t("basket.cross")}</dd>
        </div>
        {basket.held !== undefined && (
          <div>
            <dt>{t("basket.held")}</dt>
            <dd>{t("basket.shares", { shares: f.tokens(basket.held) })}</dd>
          </div>
        )}
        {state.cross && (
          <div>
            <dt>{t("basket.marginNow")}</dt>
            <dd>
              {t("basket.against", { equity: f.usd(state.cross.health.equity), requirement: f.usd(state.cross.health.requirement) })}
              <small>
                {inBaskets ? t("basket.marginedAs", { assets: inBaskets }) : t("basket.nothing")}
              </small>
            </dd>
          </div>
        )}
      </dl>
      <div className="mt-5">{sign(intent, t(`basket.action.${action}`))}</div>
      <Notes
        notes={basketNotes(f, {
          basket: basket.key,
          assets: basket.assets,
          target: basket.target,
          pending: basket.pending,
          now: state.timestamp,
        })}
      />
    </>
  );
}

function Book({ state, desk, open }: { state: TradeState; desk: Desk; open: Open }) {
  const t = useTranslations("trade");
  const f = useFormat();
  const shortsHeld = state.markets.filter((m) => m.short && shortStanding(m.short.position, m.epoch) !== "none");
  const covers = state.cover?.holdings ?? [];
  const unread = state.cover !== undefined && state.cover.holdings === undefined;
  const basketsHeld = state.baskets.filter((b) => (b.held ?? 0n) > 0n);
  const payouts = state.cover?.payouts ?? 0n;
  if (shortsHeld.length + covers.length + basketsHeld.length === 0 && payouts === 0n && !unread) return null;
  return (
    <section aria-labelledby="book-title" className="mt-14">
      <h2 id="book-title" className="label">
        {t("book.title")}
      </h2>
      <dl className="ticket ticket-book num">
        {shortsHeld.map((m) => (
          <ShortLine key={m.asset} market={m} open={open} />
        ))}
        {covers.map((c) => (
          <div key={String(c.id)}>
            <dt>
              {t("book.cover", { asset: c.asset, id: c.id.toString() })}
              <small>
                {c.series.settlement
                  ? t("book.coverSettled", {
                      notional: f.usdg(c.notional),
                      deductible: f.bpsPercent(c.deductibleBps),
                      limit: f.bpsPercent(c.limitBps),
                      status: c.series.status,
                      price: `$${f.price(c.series.settlement.price)}`,
                      reference: `$${f.price(c.series.settlement.referencePrice)}`,
                      payout: f.usdg(gapCover.payout(c, c.series.settlement)),
                    })
                  : t("book.coverDetail", {
                      notional: f.usdg(c.notional),
                      deductible: f.bpsPercent(c.deductibleBps),
                      limit: f.bpsPercent(c.limitBps),
                      status: c.series.status,
                    })}
              </small>
            </dt>
            <dd>{c.series.status === "open" ? t("book.waits") : desk.send({ kind: "release", id: c.id }, t("book.release"))}</dd>
          </div>
        ))}
        {unread && (
          <div>
            <dt>{t("book.unreadTerm")}</dt>
            <dd data-tone="warning">{t("book.unread")}</dd>
          </div>
        )}
        {payouts > 0n && (
          <div>
            <dt>{t("book.credited")}</dt>
            <dd>
              {t("book.amountUsdg", { amount: f.usdg(payouts) })} {desk.send({ kind: "claimCover" }, t("book.claim"))}
            </dd>
          </div>
        )}
        {basketsHeld.map((b) => (
          <BasketLine key={b.key} basket={b} open={open} />
        ))}
      </dl>
    </section>
  );
}

function ShortLine({ market, open }: { market: Market; open: Open }) {
  const t = useTranslations("trade");
  const f = useFormat();
  const { position, health } = market.short!;
  const standing = shortStanding(position, market.epoch);
  return (
    <div>
      <dt>
        {t("book.short", { asset: market.asset })}
        <small>
          {standing === "earlier"
            ? t("book.shortEarlier")
            : t("book.shortOwes", {
                debt: f.tokens(position.debt),
                asset: market.asset,
                equity: f.usd(health.equity),
                requirement: f.usd(health.requirement),
              })}
        </small>
      </dt>
      <dd data-tone={standing === "deficit" ? "warning" : undefined}>
        {t("book.amountUsdg", { amount: position.usdgHeld < 0n ? `−${f.usdg(-position.usdgHeld)}` : f.usdg(position.usdgHeld) })}
        {standing === "open" && (
          <button
            type="button"
            className="pill"
            onClick={() =>
              open({ kind: "short", action: "buyBack", pick: market.asset, typed: f.unparse(position.debt + position.debt / 1_000n + 1n, 18) })
            }
          >
            {t("book.buyBack")}
          </button>
        )}
        {standing === "earlier" && position.usdgHeld > 0n && (
          <button
            type="button"
            className="pill"
            onClick={() => open({ kind: "short", action: "marginOut", pick: market.asset, typed: f.unparse(position.usdgHeld, 6) })}
          >
            {t("book.takeOut")}
          </button>
        )}
      </dd>
    </div>
  );
}

function BasketLine({ basket, open }: { basket: BasketState; open: Open }) {
  const t = useTranslations("trade");
  const f = useFormat();
  return (
    <div>
      <dt>
        {t("book.basket", { key: basket.key })}
        <small>{t("book.inCross")}</small>
      </dt>
      <dd>
        {t("book.shares", { shares: f.tokens(basket.held ?? 0n) })}
        <button
          type="button"
          className="pill"
          onClick={() => open({ kind: "basket", action: "redeem", pick: basket.key, typed: f.unparse(basket.held ?? 0n, 18) })}
        >
          {t("book.redeem")}
        </button>
      </dd>
    </div>
  );
}

/** A segmented choice of `options`, the current one pressed. */
export function Choice<value extends string>({
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
    <div className="choice mt-4" role="group" aria-label={label}>
      {(Object.keys(options) as value[]).map((option) => (
        <button key={option} type="button" aria-pressed={option === value} onClick={() => onChange(option)}>
          {options[option]}
        </button>
      ))}
    </div>
  );
}

/** The ticket's verb, picked in its head line. */
function Verb<value extends string>({
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
    <select className="picker" aria-label={label} value={value} onChange={(e) => onChange(e.target.value as value)}>
      {(Object.keys(options) as value[]).map((option) => (
        <option key={option} value={option}>
          {options[option]}
        </option>
      ))}
    </select>
  );
}

function Notes({ notes }: { notes: Note[] }) {
  const t = useTranslations("trade");
  const say = useMsg();
  return (
    <ul aria-label={t("notesAria")} className="notes mt-6 max-w-[62ch]">
      {notes.map((n) => (
        <li key={say(n.message)} data-tone={n.tone}>
          {n.term && <b className="font-semibold text-strong">{say(n.term)}. </b>}
          {say(n.message)}
        </li>
      ))}
    </ul>
  );
}

