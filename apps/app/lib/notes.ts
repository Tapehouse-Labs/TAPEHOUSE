// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, zeroAddress } from "viem";
import type { AccountState, Position } from "./accountRead";
import type { Formatters } from "./format";
import { type Msg, msg } from "./messages";
import { premiumOver, weekendPlan } from "./week";

/** A line of what the visitor should know before signing: what it says, how loud, and an optional lead word. */
export type Note = { message: Msg; tone?: "warning" | "loss"; term?: Msg };

/** Whether a closure lies ahead of new loans: from the regular open on its last trading day, or while unknown. */
export function weekendAhead(state: AccountState, position: Position): boolean {
  return position.regime === 0 || ((position.regime === 2 || position.regime === 3) && state.session.boundaryMs !== 0n);
}

/** What the position's market says before any signature: halts, missing pools, the regime and the 5× rule. */
export function positionNotes(f: Formatters, state: AccountState, position: Position): Note[] {
  const notes: Note[] = [];
  if (position.regime === 0) notes.push({ message: msg("account.notes.sessionUnknown"), tone: "warning" });
  state.markets.forEach((market, i) => {
    const held = position.stocks.find((s) => s.asset === market.asset);
    if (!held) return;
    if (market.state === 0) notes.push({ message: msg("account.notes.halted", { asset: market.asset }), tone: "warning" });
    if ((position.missing >> i) & 1)
      notes.push({ message: msg("account.notes.poolUnread", { asset: market.asset }), tone: "warning" });
    if (held.inBaskets > 0n)
      notes.push({ message: msg("account.notes.viaBaskets", { tokens: f.tokens(held.inBaskets), asset: market.asset }) });
  });
  if (position.regime === 0 || position.regime === 1 || state.session.boundaryMs !== 0n)
    notes.push({ message: msg("account.notes.closureRules") });
  if (weekendAhead(state, position) && position.regime !== 0) {
    const plan = weekendPlan({
      gross: position.gross,
      equity: position.equity,
      weekendLeverage: state.limits.weekendLeverage,
    });
    if (plan)
      notes.push({ message: msg("account.notes.weekendPlan", { add: f.usd(plan.add), reduce: f.usd(plan.reduce) }), tone: "warning" });
  }
  return notes;
}

/** What a loan of `amount` USDG says before its signature: the guardian, the premium and the caps. */
export function borrowNotes(f: Formatters, state: AccountState, position: Position, amount: bigint): Note[] {
  const { limits } = state;
  const notes: Note[] = [];
  if (limits.borrowingPaused) notes.push({ message: msg("account.notes.borrowPaused"), tone: "warning" });
  if (limits.premiumRate > 0)
    notes.push({
      message: msg("account.notes.premium", {
        rate: limits.premiumRate / 100,
        owed: f.usdg(premiumOver(amount + position.debt, limits.premiumRate, 48)),
      }),
    });
  notes.push({
    message: weekendAhead(state, position)
      ? msg("account.notes.capsUntilReopen", {
          debt: f.usdg(limits.debt),
          cap: f.usdg(limits.debtCap),
          weekendCap: f.usdg(limits.weekendDebtCap),
        })
      : msg("account.notes.caps", { debt: f.usdg(limits.debt), cap: f.usdg(limits.debtCap) }),
  });
  return notes;
}

/** What a deposit of `token` says before its signature: the asset's cap, a pending write-down and the issuer. */
export function depositNotes(
  f: Formatters,
  state: AccountState,
  token: Address | "ETH",
  account: Address,
  owner: Address,
  sponsored: boolean,
): Note[] {
  const notes: Note[] = [];
  if (token === "ETH") notes.push({ message: msg("account.notes.wrapsEther") });
  else if (sponsored) notes.push({ message: msg("account.notes.signsPermit") });
  const market = state.markets.find((m) => m.token === token && m.token !== zeroAddress);
  if (market?.holding) {
    const { quantity, cap, balance } = market.holding;
    notes.push({
      message: msg("account.notes.depositCap", { asset: market.asset, quantity: f.tokens(quantity), cap: f.tokens(cap) }),
    });
    if (balance < quantity)
      notes.push({
        message: msg("account.notes.writeDown", { balance: f.tokens(balance), quantity: f.tokens(quantity), asset: market.asset }),
        tone: "warning",
      });
  }
  if (market?.issuer) notes.push(...issuerNotes(market.asset, market.issuer, account, owner));
  return notes;
}

/** The issuer's pause and blocks, for `asset`. */
export function issuerNotes(
  asset: string,
  issuer: { paused: boolean; blocked: Record<string, boolean> },
  account: Address,
  owner: Address,
): Note[] {
  const notes: Note[] = [];
  if (issuer.paused) notes.push({ message: msg("account.notes.issuerPaused", { asset }), tone: "warning" });
  const blocked = Object.entries(issuer.blocked)
    .filter(([, is]) => is)
    .map(([address]) => address);
  if (blocked.some((a) => a === account || a === owner))
    notes.push({ message: msg("account.notes.issuerBlockedYou", { asset }), tone: "warning" });
  else if (blocked.length > 0) notes.push({ message: msg("account.notes.issuerBlockedAccounts", { asset }), tone: "warning" });
  return notes;
}
