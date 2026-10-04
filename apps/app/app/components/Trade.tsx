// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { gapCover, shorts } from "@tapehouse/sdk";
import { lazy, type ReactNode, Suspense, useRef, useState } from "react";
import type { Address } from "viem";
import { useChains } from "wagmi";
import { useDeployments } from "@/app/providers";
import { parseAmount, tokens, unparse, usd, usdg } from "@/lib/amounts";
import { percent } from "@/lib/earn";
import { explain } from "@/lib/errors";
import { price } from "@/lib/format";
import { useMinute, useWallet } from "@/lib/hooks";
import type { Intent } from "@/lib/intent";
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
  slippageWords,
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
const KIND: Record<Kind, string> = { short: "Short", cover: "Gap cover", basket: "Baskets" };

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
          {wallet.state === "wrong-network" ? (
            <>
              Your wallet is on another network. Tapehouse runs on <span className="text-strong">{chain.name}</span>.
            </>
          ) : (
            "Every line of the ticket is read from the chain. Connect a wallet to sign it."
          )}
        </p>
        <div className="mt-4 flex min-h-[44px] flex-wrap items-center gap-3">
          {wallet.state === "wrong-network" ? (
            <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
              {wallet.pending ? "Waiting for your wallet…" : `Switch to ${chain.name}`}
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
          <h1 className="label">Trade · shorts, gap cover and baskets</h1>
          <Choice options={KIND} value={kind} onChange={setKind} label="Ticket" />
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
          {data && <p className="label mt-3 num">Read at block {data.blockNumber.toLocaleString("en-US")}</p>}
        </div>
        {account}
        {data && desk && <Book state={data} desk={desk} open={open} />}
      </div>
    </main>
  );
}

type ShortAction = "sell" | "buyBack" | "marginIn" | "marginOut";
const SHORT_ACTION: Record<ShortAction, string> = {
  sell: "Sell",
  buyBack: "Buy back",
  marginIn: "Add margin",
  marginOut: "Take out",
};
const SHORT_VERB: Record<ShortAction, string> = { sell: "sell", buyBack: "buy back", marginIn: "add", marginOut: "take out" };

/** What a ticket's foot holds: the account's send button, or how a visitor connects. */
type Sign = (intent: Intent | undefined, label: string) => ReactNode;

function ShortTicket({ state, desk, sign, preset }: { state: TradeState; desk: Desk | undefined; sign: Sign; preset?: Preset }) {
  const deployments = useDeployments();
  const markets = state.markets.filter((m) => m.fee > 0);
  const [chosen, setChosen] = useState(preset?.pick);
  const [action, setAction] = useState((preset?.action as ShortAction | undefined) ?? "sell");
  const [typed, setTyped] = useState(preset?.typed ?? "");
  const [slippage, setSlippage] = useState(DEFAULT_SLIPPAGE);
  const market = markets.find((m) => m.asset === chosen) ?? markets[0];
  const usdgUnits = action === "marginIn" || action === "marginOut";
  const amount = parseAmount(typed, usdgUnits ? 6 : 18);
  const bps = parseSlippage(slippage);
  const side = action === "buyBack" ? "buyBack" : "sell";
  const position = market?.short?.position;
  const quoted = action === "buyBack" && position && amount !== undefined && amount > position.debt ? position.debt : amount;
  const pool = usePoolQuote(side, market?.asset ?? "", usdgUnits ? undefined : quoted, bps);
  const standing = position && market ? shortStanding(position, market.epoch) : "none";
  const intent = market && amount !== undefined ? shortIntent(action, market.asset, amount, pool.data?.limit) : undefined;
  const call = intent && desk && intent.kind !== "shortDeposit" ? simulated(deployments, intent, desk) : undefined;
  const after = useHealthAfter(desk?.account, market?.asset ?? "", call);

  if (!deployments.tapehouse.ShortPositions || !market)
    return <p className="mt-6 text-[15px] text-body">No Stock Token can be shorted on this chain yet.</p>;
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
      [action === "sell" ? "The pool pays" : "The pool asks", pool.error ? "No quote" : pool.data ? `$${usdg(pool.data.quoted)}` : "—"],
      [
        action === "sell" ? "At least" : "At most",
        <span key="limit" className="inline-flex items-center gap-2">
          <span>{pool.data && !pool.error ? `$${usdg(pool.data.limit)}` : "—"}</span>
          <span className="text-muted">at</span>
          <input
            className="amount amount-small"
            inputMode="decimal"
            aria-label="Slippage, in percent"
            aria-invalid={bps === undefined}
            value={slippage}
            onChange={(e) => setSlippage(e.target.value)}
          />
          <span className="text-muted">%</span>
        </span>,
      ],
      action === "sell"
        ? ["A token fetches no less than", `$${price(floor.price)}, ${floor.words}`]
        : full
          ? ["A token costs no more than", "your limit alone, a full buy-back"]
          : ["A token costs no more than", `$${price(market.quote.high)}, the band's high edge`],
      ["Lending fee", market.borrowRate === undefined ? "—" : `${percent(market.borrowRate)} a year, in ${market.asset}`],
      ["Lenders can lend now", market.borrowable === undefined ? "—" : `${tokens(market.borrowable)} ${market.asset}`],
      ["Pool fee tier", feeTier(market.fee)],
    );
  } else {
    lines.push([action === "marginIn" ? "From" : "To", "your wallet"]);
  }
  if (health)
    lines.push([
      "Margin now",
      `${usd(health.equity)} against ${usd(health.requirement)}`,
    ]);
  if (desk)
    lines.push([
      "Margin after",
      marginAfter ? `${usd(marginAfter.equity)} against ${usd(marginAfter.requirement)}` : reverts ? "It would revert" : "—",
    ]);
  const belowFloor =
    action === "sell" && amount !== undefined && pool.data && pool.data.quoted * 10n ** 20n < amount * floor.price;
  const notes: Note[] = [
    ...(standing === "earlier"
      ? [{ term: "Earlier book", text: `This short is of an earlier book: it owes nothing and holds only its ${usdg(position!.usdgHeld)} USDG to take out.` }]
      : []),
    ...(standing === "deficit"
      ? [{ term: "Deficit", text: "This short's USDG cannot buy back what it owes: only a liquidation closes it (Insolvent).", tone: "warning" as const }]
      : []),
    ...shortNotes({ asset: market.asset, borrowRate: market.borrowRate ?? 0n, fee: market.fee, restricted: market.restricted }),
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
      <dl className="ticket num" aria-label="Ticket">
        <div className="ticket-head">
          <dt>
            You{" "}
            <Verb options={SHORT_VERB} value={action} onChange={(a) => (setAction(a), setTyped(""))} label="Short action" />
          </dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label={`Amount, ${SHORT_ACTION[action].toLowerCase()}`}
              aria-invalid={typed !== "" && (amount === undefined || (max !== undefined && amount > max))}
              value={typed}
              onChange={(e) => setTyped(e.target.value)}
            />{" "}
            {action === "marginIn" ? "USDG to your" : action === "marginOut" ? "USDG from your" : ""}
            <select className="picker" aria-label="Stock Token" value={market.asset} onChange={(e) => setChosen(e.target.value)}>
              {markets.map((m) => (
                <option key={m.asset} value={m.asset}>
                  {m.asset}
                </option>
              ))}
            </select>
            {usdgUnits ? "short" : ""}
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
          <button type="button" className="underline-offset-4 hover:underline" onClick={() => setTyped(unparse(max, usdgUnits ? 6 : 18))}>
            {action === "buyBack" ? "All it owes, with room for the interest until it lands" : `Max ${usdgUnits ? usdg(max) : tokens(max)}`}
          </button>
        </p>
      )}
      <div className="mt-5">{sign(pool.error ? undefined : intent, action === "sell" ? "Sell short" : SHORT_ACTION[action])}</div>
      {pool.error && <p className="mt-3 text-[14px] text-[var(--error)]">{explain(pool.error)}</p>}
      {belowFloor && (
        <p className="mt-3 max-w-[60ch] text-[14px] text-[var(--warning)]">
          The pool pays less than {floor.words} for each token, so the sale would revert: a short never sells below it.
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
  const now = useMinute();
  const [chosen, setChosen] = useState<string>();
  const [notional, setNotional] = useState("");
  const [deductible, setDeductible] = useState("3");
  const [limit, setLimit] = useState("10");
  const markets = state.markets;
  const market = markets.find((m) => m.asset === chosen) ?? markets.find((m) => m.asset === "SPY") ?? markets[0];
  const amount = parseAmount(notional, 6);
  const deductibleBps = parseSlippage(deductible);
  const limitBps = parseSlippage(limit);
  const layer =
    amount !== undefined && deductibleBps !== undefined && limitBps !== undefined && limitBps > deductibleBps
      ? { notional: amount, deductibleBps, limitBps }
      : undefined;
  const quote = useCoverQuote(market?.asset ?? "", state.cover?.sales && layer ? layer : undefined);
  if (!state.cover) return <p className="mt-6 text-[15px] text-body">Gap cover is not deployed on this chain yet.</p>;
  const { sales, capacity } = state.cover;
  if (!market) return null;
  const maxPremium = quote.data && withSlippage(quote.data.premium, COVER_SLIPPAGE);
  const intent: Intent | undefined =
    layer && maxPremium !== undefined && !quote.error ? { kind: "buyCover", asset: market.asset, ...layer, maxPremium } : undefined;
  return (
    <>
      <dl className="ticket num" aria-label="Ticket">
        <div className="ticket-head">
          <dt>You cover</dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label="Notional, in USDG"
              aria-invalid={notional !== "" && amount === undefined}
              value={notional}
              onChange={(e) => setNotional(e.target.value)}
            />{" "}
            USDG of{" "}
            <select className="picker" aria-label="Stock Token" value={market.asset} onChange={(e) => setChosen(e.target.value)}>
              {markets.map((m) => (
                <option key={m.asset} value={m.asset}>
                  {m.asset}
                </option>
              ))}
            </select>
          </dd>
        </div>
        <div>
          <dt>Pays the fall between</dt>
          <dd className="inline-flex items-center gap-2">
            <input className="amount amount-small" inputMode="decimal" aria-label="Deductible, in percent" value={deductible} onChange={(e) => setDeductible(e.target.value)} />
            <span className="text-muted">% and</span>
            <input className="amount amount-small" inputMode="decimal" aria-label="Limit, in percent" value={limit} onChange={(e) => setLimit(e.target.value)} />
            <span className="text-muted">%</span>
          </dd>
        </div>
        <div>
          <dt>Sales</dt>
          <dd data-tone={sales ? undefined : "warning"}>
            {sales ? `close in ${countdown(sales.endsMs, now)}` : "closed until the market reopens with a close ahead"}
          </dd>
        </div>
        {sales && (
          <>
            <div>
              <dt>Premium</dt>
              <dd>
                {quote.error ? "No quote" : quote.data ? `$${usdg(quote.data.premium)}` : "—"}
                {maxPremium !== undefined && !quote.error && <small>at most ${usdg(maxPremium)}, as the feed moves it</small>}
              </dd>
            </div>
            <div>
              <dt>Writers&apos; USDG it reserves</dt>
              <dd>{quote.data ? `${usdg(quote.data.reserve)} of ${usdg(capacity)} free` : `${usdg(capacity)} free`}</dd>
            </div>
            <div>
              <dt>Priced at a weekend gap of</dt>
              <dd>
                {quote.data
                  ? `${gapPercent(quote.data.gap)}${quote.data.weekMove > 0n && quote.data.weekMove < 2n ** 64n - 1n ? `, from last week's move of ${gapPercent(quote.data.weekMove)}` : ""}`
                  : "—"}
              </dd>
            </div>
            <div>
              <dt>Smallest deductible now</dt>
              <dd>{quote.data ? `${slippageWords(quote.data.minDeductible)}` : "—"}</dd>
            </div>
          </>
        )}
      </dl>
      <div className="mt-5">{sign(sales ? intent : undefined, "Buy cover")}</div>
      {quote.error && <p className="mt-3 text-[14px] text-[var(--error)]">{explain(quote.error)}</p>}
      <Notes notes={coverPurchaseNotes(market.asset)} />
    </>
  );
}

type BasketAction = "mint" | "redeem" | "unwrap";
const COVER_SLIPPAGE = 50n;
const MINT_SLIPPAGE = 10n;
const BASKET_ACTION: Record<BasketAction, string> = { mint: "Mint", redeem: "Redeem", unwrap: "Unwrap" };
const BASKET_VERB: Record<BasketAction, string> = { mint: "mint", redeem: "redeem", unwrap: "unwrap" };

function BasketTicket({ state, desk, sign, preset }: { state: TradeState; desk: Desk | undefined; sign: Sign; preset?: Preset }) {
  const [chosen, setChosen] = useState(preset?.pick);
  const [action, setAction] = useState((preset?.action as BasketAction | undefined) ?? "mint");
  const [typed, setTyped] = useState(preset?.typed ?? "");
  const basket = state.baskets.find((b) => b.key === chosen) ?? state.baskets[0];
  const shares = parseAmount(typed, 18);
  const preview = useBasketPreview(action === "mint" ? "mint" : "redeem", basket?.key ?? "", basket ? shares : undefined);
  if (!basket) return <p className="mt-6 text-[15px] text-body">No basket is deployed on this chain yet.</p>;
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
    units ? basket.assets.map((a, i) => `${tokens(units[i] ?? 0n)} ${a}`).join(" + ") : "—";
  return (
    <>
      <dl className="ticket num" aria-label="Ticket">
        <div className="ticket-head">
          <dt>
            You{" "}
            <Verb options={BASKET_VERB} value={action} onChange={(a) => (setAction(a), setTyped(""))} label="Basket action" />
          </dt>
          <dd>
            <input
              className="amount"
              inputMode="decimal"
              autoComplete="off"
              placeholder="0"
              aria-label={`Shares to ${action}`}
              aria-invalid={typed !== "" && (shares === undefined || (action !== "mint" && (basket.held ?? 0n) < shares))}
              value={typed}
              onChange={(e) => setTyped(e.target.value)}
            />{" "}
            {shares === 10n ** 18n ? "share" : "shares"} of{" "}
            <select className="picker" aria-label="Basket" value={basket.key} onChange={(e) => setChosen(e.target.value)}>
              {state.baskets.map((b) => (
                <option key={b.key} value={b.key}>
                  {b.key}
                </option>
              ))}
            </select>
          </dd>
        </div>
        <div>
          <dt>{action === "mint" ? "It takes from your wallet" : action === "redeem" ? "Your wallet receives" : "Your cross position holds"}</dt>
          <dd data-tone={short ? "warning" : undefined}>
            {words(amounts)}
            {maxAssets && <small>at most {words(maxAssets)}, as other mints and redemptions move it</small>}
          </dd>
        </div>
        <div>
          <dt>{action === "mint" ? "Into" : "From"}</dt>
          <dd>your cross position</dd>
        </div>
        {basket.held !== undefined && (
          <div>
            <dt>Your cross position holds</dt>
            <dd>{tokens(basket.held)} shares</dd>
          </div>
        )}
        {state.cross && (
          <div>
            <dt>Cross margin now</dt>
            <dd>
              {usd(state.cross.health.equity)} against {usd(state.cross.health.requirement)}
              <small>
                {Object.entries(state.cross.inBaskets)
                  .filter(([, amount]) => amount > 0n)
                  .map(([asset, amount]) => `${tokens(amount)} ${asset}`)
                  .join(" + ")
                  .replace(/^(.+)$/, "margined as $1 through baskets") || "nothing held through baskets yet"}
              </small>
            </dd>
          </div>
        )}
      </dl>
      <div className="mt-5">{sign(intent, BASKET_ACTION[action])}</div>
      <Notes
        notes={basketNotes({
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
  const shortsHeld = state.markets.filter((m) => m.short && shortStanding(m.short.position, m.epoch) !== "none");
  const covers = state.cover?.holdings ?? [];
  const unread = state.cover !== undefined && state.cover.holdings === undefined;
  const basketsHeld = state.baskets.filter((b) => (b.held ?? 0n) > 0n);
  const payouts = state.cover?.payouts ?? 0n;
  if (shortsHeld.length + covers.length + basketsHeld.length === 0 && payouts === 0n && !unread) return null;
  return (
    <section aria-labelledby="book-title" className="mt-14">
      <h2 id="book-title" className="label">
        Your trades
      </h2>
      <dl className="ticket ticket-book num">
        {shortsHeld.map((m) => (
          <ShortLine key={m.asset} market={m} open={open} />
        ))}
        {covers.map((c) => (
          <div key={String(c.id)}>
            <dt>
              Cover {c.asset} #{c.id.toString()}
              <small>
                {usdg(c.notional)} USDG, {slippageWords(c.deductibleBps)} to {slippageWords(c.limitBps)} · series {c.series.status}
                {c.series.settlement &&
                  ` at $${price(c.series.settlement.price)} from $${price(c.series.settlement.referencePrice)}, pays ${usdg(gapCover.payout(c, c.series.settlement))}`}
              </small>
            </dt>
            <dd>{c.series.status === "open" ? "Waits on its reopen" : desk.send({ kind: "release", id: c.id }, "Release")}</dd>
          </div>
        ))}
        {unread && (
          <div>
            <dt>Your covers</dt>
            <dd data-tone="warning">The node would not return the cover&apos;s logs; reload to read them again</dd>
          </div>
        )}
        {payouts > 0n && (
          <div>
            <dt>Credited to your account</dt>
            <dd>
              {usdg(payouts)} USDG {desk.send({ kind: "claimCover" }, "Claim to your wallet")}
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
  const { position, health } = market.short!;
  const standing = shortStanding(position, market.epoch);
  return (
    <div>
      <dt>
        Short {market.asset}
        <small>
          {standing === "earlier"
            ? "of an earlier book: owes nothing"
            : `owes ${tokens(position.debt)} ${market.asset} · margin ${usd(health.equity)} against ${usd(health.requirement)}`}
        </small>
      </dt>
      <dd data-tone={standing === "deficit" ? "warning" : undefined}>
        {position.usdgHeld < 0n ? `−${usdg(-position.usdgHeld)}` : usdg(position.usdgHeld)} USDG
        {standing === "open" && (
          <button
            type="button"
            className="pill"
            onClick={() =>
              open({ kind: "short", action: "buyBack", pick: market.asset, typed: unparse(position.debt + position.debt / 1_000n + 1n, 18) })
            }
          >
            Buy back
          </button>
        )}
        {standing === "earlier" && position.usdgHeld > 0n && (
          <button
            type="button"
            className="pill"
            onClick={() => open({ kind: "short", action: "marginOut", pick: market.asset, typed: unparse(position.usdgHeld, 6) })}
          >
            Take out
          </button>
        )}
      </dd>
    </div>
  );
}

function BasketLine({ basket, open }: { basket: BasketState; open: Open }) {
  return (
    <div>
      <dt>
        Basket {basket.key}
        <small>in your cross position</small>
      </dt>
      <dd>
        {tokens(basket.held ?? 0n)} shares
        <button
          type="button"
          className="pill"
          onClick={() => open({ kind: "basket", action: "redeem", pick: basket.key, typed: unparse(basket.held ?? 0n, 18) })}
        >
          Redeem
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
  return (
    <ul aria-label="Before you sign" className="notes mt-6 max-w-[62ch]">
      {notes.map((n) => (
        <li key={n.text} data-tone={n.tone}>
          {n.term && <b className="font-semibold text-strong">{n.term}. </b>}
          {n.text}
        </li>
      ))}
    </ul>
  );
}

