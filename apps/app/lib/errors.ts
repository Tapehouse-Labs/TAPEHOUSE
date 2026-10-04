// SPDX-License-Identifier: MIT OR Apache-2.0
import { decodeRevert, type Revert } from "@tapehouse/sdk";
import { BaseError, type Hex, hexToString, RpcRequestError, UserRejectedRequestError } from "viem";
import { tokens, usd, usdg } from "./amounts";

const asset = (symbol: unknown) => hexToString(symbol as Hex, { size: 32 });

/** What the accounts' and the issuer's refusals mean, in words, by error name. */
const WORDS: Record<string, (args: readonly unknown[]) => string> = {
  InsufficientMargin: ([equity, requirement]) =>
    `This leaves the position short of its requirement: ${usd(equity as bigint)} of equity against ${usd(requirement as bigint)} required.`,
  WeekendLeverageExceeded: ([gross, equity]) =>
    `Before the close the position may hold at most 5× its equity: ${usd(gross as bigint)} against ${usd(equity as bigint)} of equity.`,
  BorrowingIsPaused: () =>
    "The guardian has paused new loans and withdrawals from positions in debt. Deposits and repayments stay open.",
  DebtCapExceeded: ([debt, cap]) =>
    `All loans together would owe ${usdg(debt as bigint)} USDG, past the cap of ${usdg(cap as bigint)}.`,
  WeekendDebtCapExceeded: ([debt, cap]) =>
    `All loans together would owe ${usdg(debt as bigint)} USDG, past the cap of ${usdg(cap as bigint)} that holds until the reopen.`,
  AssetCapExceeded: ([symbol, cap]) =>
    `The accounts may hold at most ${tokens(cap as bigint)} ${asset(symbol)}, through baskets included.`,
  BasketFrozen: () => "A basket the position holds is frozen: it backs no new loan or withdrawal until it moves again.",
  AssetFrozen: ([symbol]) =>
    `${asset(symbol)} is frozen by its lending vault or its issuer: the position borrows nothing and withdraws nothing while in debt.`,
  AssetHalted: ([symbol]) => `${asset(symbol)} is halted: it backs nothing until trading resumes.`,
  OutOfReach: (args) =>
    args.length === 3
      ? `The ${asset(args[0])} lending vault can return ${tokens(args[2] as bigint)} of the ${tokens(args[1] as bigint)} asked now; recall the rest.`
      : `The lending vault can return ${tokens(args[1] as bigint)} of the ${tokens(args[0] as bigint)} asked now.`,
  RecallTooSmall: ([amount]) =>
    `A recall must queue at least a hundredth of a token; this one would queue ${tokens(amount as bigint)}.`,
  TooManyRecalls: ([symbol]) => `The position already has 3 ${asset(symbol)} recalls open, the most an owner may.`,
  NothingToSettle: () => "Nothing has come back for this recall yet.",
  InsufficientCollateral: () => "The position holds less than that.",
  BorrowingAgainstUsdg: () => "A position holding USDG repays with it before it borrows.",
  AssetInOtherPosition: ([symbol]) => `Your ${asset(symbol)} sits in another of your positions.`,
  SessionUnknown: () => "The band cannot tell whether the market is open, so new loans wait.",
  SequencerNotSettled: () => "The chain's sequencer has just restarted; new loans wait a moment.",
  CorporateActionPending: ([symbol]) => `${asset(symbol)} has a corporate action pending; new loans against it wait.`,
  LiquidityUnknown: ([symbol]) => `${asset(symbol)}'s pool cannot be read, so it backs no new loan.`,
  AssetWrittenOff: ([symbol]) => `${asset(symbol)} has been written off by its issuer.`,
  Blocked: () => "The issuer has blocked this address from moving the token.",
  IsPaused: () => "The issuer has paused the token's transfers.",
  AddressFrozen: () => "USDG's issuer has frozen this address.",
  ZeroAmount: () => "Enter an amount above zero.",
};

/** A revert in words where the app knows its meaning, else by its error's name and arguments. */
export function describeRevert({ errorName, args }: Revert): string {
  return WORDS[errorName]?.(args) ?? `${errorName}(${args.map(String).join(", ")})`;
}

/**
 * What went wrong, for the visitor: a revert in words, a declined request, a service's refusal in its own words, or
 * viem's short message.
 */
export function explain(error: unknown): string {
  const revert = decodeRevert(error);
  if (revert) return describeRevert(revert);
  if (error instanceof BaseError) {
    if (error.walk((cause) => cause instanceof UserRejectedRequestError))
      return "You declined the request in your wallet.";
    const refused = error.walk((cause) => cause instanceof RpcRequestError);
    if (refused instanceof RpcRequestError && refused.details) return `Refused: ${refused.details}.`;
    return error.shortMessage;
  }
  return error instanceof Error ? error.message : String(error);
}
