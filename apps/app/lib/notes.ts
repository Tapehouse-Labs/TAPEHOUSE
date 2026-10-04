// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, zeroAddress } from "viem";
import type { AccountState, Position } from "./accountRead";
import { tokens, usd, usdg } from "./amounts";
import { premiumOver, weekendPlan } from "./week";

export type Note = { text: string; tone?: "warning" | "loss"; term?: string };

/** Whether a closure lies ahead of new loans: from the regular open on its last trading day, or while unknown. */
export function weekendAhead(state: AccountState, position: Position): boolean {
  return position.regime === 0 || ((position.regime === 2 || position.regime === 3) && state.session.boundaryMs !== 0n);
}

/** What the position's market says before any signature: halts, missing pools, the regime and the 5× rule. */
export function positionNotes(state: AccountState, position: Position): Note[] {
  const notes: Note[] = [];
  if (position.regime === 0)
    notes.push({
      text: "The band cannot tell whether the market is open, so the requirement counts as closed and borrowing stops at the price above.",
      tone: "warning",
    });
  state.markets.forEach((market, i) => {
    const held = position.stocks.find((s) => s.asset === market.asset);
    if (!held) return;
    if (market.state === 0) notes.push({ text: `${market.asset} is halted and counts for nothing.`, tone: "warning" });
    if ((position.missing >> i) & 1)
      notes.push({
        text: `${market.asset}'s pool is absent or unread, so its liquidation is charged at the governed depth.`,
        tone: "warning",
      });
    if (held.inBaskets > 0n)
      notes.push({ text: `Holds ${tokens(held.inBaskets)} ${market.asset} through baskets, counted as the tokens they redeem for.` });
  });
  if (position.regime === 0 || position.regime === 1 || state.session.boundaryMs !== 0n)
    notes.push({
      text: "From the close before a weekend or holiday to the regular open after it, and while the session is unknown, a position is liquidated only if it falls short at both edges of its bands, at most a tenth of each holding an hour, and a purchase repays at most half of what it owes above 2,000 USDG. From the 24/5 reopen, a position short at its low edges may be enrolled in the reopening auction and sold at its clearing price at the regular open, with no such half, up to what repays it.",
    });
  if (weekendAhead(state, position) && position.regime !== 0) {
    const plan = weekendPlan({
      gross: position.gross,
      equity: position.equity,
      weekendLeverage: state.limits.weekendLeverage,
    });
    if (plan)
      notes.push({
        text: `To stay within 5× at the close, add ${usd(plan.add)} of equity or reduce the exposure by ${usd(plan.reduce)}.`,
        tone: "warning",
      });
  }
  return notes;
}

/** What a loan of `amount` USDG says before its signature: the guardian, the premium and the caps. */
export function borrowNotes(state: AccountState, position: Position, amount: bigint): Note[] {
  const { limits } = state;
  const notes: Note[] = [];
  if (limits.borrowingPaused)
    notes.push({
      text: "The guardian has paused new loans and withdrawals from positions in debt. Deposits and repayments stay open.",
      tone: "warning",
    });
  if (limits.premiumRate > 0)
    notes.push({
      text: `While the market is closed, ${limits.premiumRate / 100}% a year is added on top of interest: this loan would owe ${usdg(premiumOver(amount + position.debt, limits.premiumRate, 48))} USDG over a 48-hour closure. A repayment pays the debt first, then the premium.`,
    });
  notes.push({
    text: `All loans together owe ${usdg(limits.debt)} of their ${usdg(limits.debtCap)} USDG cap${weekendAhead(state, position) ? `; until the reopen, new loans stop at ${usdg(limits.weekendDebtCap)}` : ""}.`,
  });
  return notes;
}

/** What a deposit of `token` says before its signature: the asset's cap, a pending write-down and the issuer. */
export function depositNotes(state: AccountState, token: Address | "ETH", account: Address, owner: Address, sponsored: boolean): Note[] {
  const notes: Note[] = [];
  if (token === "ETH")
    notes.push({ text: "Your wallet wraps the ether into WETH for your account and pays that gas; the deposit follows." });
  else if (sponsored) notes.push({ text: "Your wallet signs a permit for free and Tapehouse pays the gas." });
  const market = state.markets.find((m) => m.token === token && m.token !== zeroAddress);
  if (market?.holding) {
    const { quantity, cap, balance } = market.holding;
    notes.push({ text: `${market.asset} deposits: ${tokens(quantity)} of their ${tokens(cap)} cap.` });
    if (balance < quantity)
      notes.push({
        text: `A write-down by the issuer is pending: the accounts hold ${tokens(balance)} of the ${tokens(quantity)} ${market.asset} they count.`,
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
  if (issuer.paused) notes.push({ text: `The issuer has paused ${asset}'s transfers.`, tone: "warning" });
  const blocked = Object.entries(issuer.blocked)
    .filter(([, is]) => is)
    .map(([address]) => address);
  if (blocked.some((a) => a === account || a === owner))
    notes.push({ text: `The issuer has blocked your address from moving ${asset}.`, tone: "warning" });
  else if (blocked.length > 0) notes.push({ text: `The issuer has blocked the accounts from moving ${asset}.`, tone: "warning" });
  return notes;
}
