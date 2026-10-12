// SPDX-License-Identifier: MIT OR Apache-2.0
import { decodeRevert, type Revert } from "@tapehouse/sdk";
import { BaseError, type Hex, hexToString, RpcRequestError, UserRejectedRequestError } from "viem";
import type { Formatters } from "./format";
import { MessageError, type Msg, msg } from "./messages";

const asset = (symbol: unknown) => hexToString(symbol as Hex, { size: 32 });

/** What the accounts' and the issuer's refusals mean, by error name. */
const WORDS: Record<string, (f: Formatters, args: readonly unknown[]) => Msg> = {
  InsufficientMargin: (f, [equity, requirement]) =>
    msg("errors.revert.InsufficientMargin", { equity: f.usd(equity as bigint), requirement: f.usd(requirement as bigint) }),
  WeekendLeverageExceeded: (f, [gross, equity]) =>
    msg("errors.revert.WeekendLeverageExceeded", { gross: f.usd(gross as bigint), equity: f.usd(equity as bigint) }),
  BorrowingIsPaused: () => msg("errors.revert.BorrowingIsPaused"),
  DebtCapExceeded: (f, [debt, cap]) =>
    msg("errors.revert.DebtCapExceeded", { debt: f.usdg(debt as bigint), cap: f.usdg(cap as bigint) }),
  WeekendDebtCapExceeded: (f, [debt, cap]) =>
    msg("errors.revert.WeekendDebtCapExceeded", { debt: f.usdg(debt as bigint), cap: f.usdg(cap as bigint) }),
  AssetCapExceeded: (f, [symbol, cap]) =>
    msg("errors.revert.AssetCapExceeded", { cap: f.tokens(cap as bigint), asset: asset(symbol) }),
  BasketFrozen: () => msg("errors.revert.BasketFrozen"),
  AssetFrozen: (_, [symbol]) => msg("errors.revert.AssetFrozen", { asset: asset(symbol) }),
  AssetHalted: (_, [symbol]) => msg("errors.revert.AssetHalted", { asset: asset(symbol) }),
  OutOfReach: (f, args) =>
    args.length === 3
      ? msg("errors.revert.OutOfReachOfAsset", {
          asset: asset(args[0]),
          returned: f.tokens(args[2] as bigint),
          asked: f.tokens(args[1] as bigint),
        })
      : msg("errors.revert.OutOfReach", { returned: f.tokens(args[1] as bigint), asked: f.tokens(args[0] as bigint) }),
  RecallTooSmall: (f, [amount]) => msg("errors.revert.RecallTooSmall", { amount: f.tokens(amount as bigint) }),
  TooManyRecalls: (_, [symbol]) => msg("errors.revert.TooManyRecalls", { asset: asset(symbol) }),
  NothingToSettle: () => msg("errors.revert.NothingToSettle"),
  InsufficientCollateral: () => msg("errors.revert.InsufficientCollateral"),
  BorrowingAgainstUsdg: () => msg("errors.revert.BorrowingAgainstUsdg"),
  AssetInOtherPosition: (_, [symbol]) => msg("errors.revert.AssetInOtherPosition", { asset: asset(symbol) }),
  SessionUnknown: () => msg("errors.revert.SessionUnknown"),
  SequencerNotSettled: () => msg("errors.revert.SequencerNotSettled"),
  CorporateActionPending: (_, [symbol]) => msg("errors.revert.CorporateActionPending", { asset: asset(symbol) }),
  LiquidityUnknown: (_, [symbol]) => msg("errors.revert.LiquidityUnknown", { asset: asset(symbol) }),
  AssetWrittenOff: (_, [symbol]) => msg("errors.revert.AssetWrittenOff", { asset: asset(symbol) }),
  Blocked: () => msg("errors.revert.Blocked"),
  IsPaused: () => msg("errors.revert.IsPaused"),
  AddressFrozen: () => msg("errors.revert.AddressFrozen"),
  ZeroAmount: () => msg("errors.revert.ZeroAmount"),
};

/** A revert in words where the app knows its meaning, else by its error's name and arguments. */
export function describeRevert(f: Formatters, { errorName, args }: Revert): Msg {
  return WORDS[errorName]?.(f, args) ?? msg("errors.unknownRevert", { name: errorName, args: args.map(String).join(", ") });
}

/**
 * What went wrong, for the visitor: a revert in words, one of the app's own messages, a declined request, a service's
 * refusal in its own words, or viem's short message.
 */
export function explain(f: Formatters, error: unknown): Msg {
  const revert = decodeRevert(error);
  if (revert) return describeRevert(f, revert);
  if (error instanceof MessageError) return error.msg;
  if (error instanceof BaseError) {
    if (error.walk((cause) => cause instanceof UserRejectedRequestError)) return msg("errors.declined");
    const refused = error.walk((cause) => cause instanceof RpcRequestError);
    if (refused instanceof RpcRequestError && refused.details) return msg("errors.refused", { details: refused.details });
    return msg("errors.text", { text: error.shortMessage });
  }
  return msg("errors.text", { text: error instanceof Error ? error.message : String(error) });
}
