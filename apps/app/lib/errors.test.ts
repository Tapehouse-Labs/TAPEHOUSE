// SPDX-License-Identifier: MIT OR Apache-2.0
import { errorsAbi } from "@tapehouse/sdk";
import {
  BaseError,
  ContractFunctionExecutionError,
  ContractFunctionRevertedError,
  encodeErrorResult,
  InvalidParamsRpcError,
  RpcRequestError,
  UserRejectedRequestError,
} from "viem";
import { describe, expect, test } from "vitest";
import { explain } from "./errors";

function reverted(data: `0x${string}`) {
  const cause = new ContractFunctionRevertedError({ abi: errorsAbi, data, functionName: "quote" });
  return new ContractFunctionExecutionError(cause, { abi: errorsAbi, functionName: "quote" });
}

describe("explain", () => {
  test("names a contract's custom error and its arguments", () => {
    const data = encodeErrorResult({
      abi: errorsAbi,
      errorName: "UnknownAsset",
      args: ["0x4d53465400000000000000000000000000000000000000000000000000000000"],
    });
    expect(explain(reverted(data))).toBe(
      "UnknownAsset(0x4d53465400000000000000000000000000000000000000000000000000000000)",
    );
  });

  test("says what the accounts' refusals mean, in words, with their figures", () => {
    const say = (errorName: string, args: readonly unknown[]) =>
      explain(reverted(encodeErrorResult({ abi: errorsAbi, errorName, args } as never)));
    expect(say("InsufficientMargin", [9_000n * 10n ** 18n, 9_500n * 10n ** 18n])).toBe(
      "This leaves the position short of its requirement: $9,000.00 of equity against $9,500.00 required.",
    );
    expect(say("WeekendLeverageExceeded", [60_000n * 10n ** 18n, 10_000n * 10n ** 18n])).toBe(
      "Before the close the position may hold at most 5× its equity: $60,000.00 against $10,000.00 of equity.",
    );
    expect(say("BorrowingIsPaused", [])).toBe(
      "The guardian has paused new loans and withdrawals from positions in debt. Deposits and repayments stay open.",
    );
    expect(say("DebtCapExceeded", [1_200_000_000_000n, 1_000_000_000_000n])).toBe(
      "All loans together would owe 1,200,000.00 USDG, past the cap of 1,000,000.00.",
    );
    expect(say("OutOfReach", ["0x4e56444100000000000000000000000000000000000000000000000000000000", 5n * 10n ** 18n, 2n * 10n ** 18n])).toBe(
      "The NVDA lending vault can return 2 of the 5 asked now; recall the rest.",
    );
    expect(say("TooManyRecalls", ["0x4e56444100000000000000000000000000000000000000000000000000000000"])).toBe(
      "The position already has 3 NVDA recalls open, the most an owner may.",
    );
    expect(say("BasketFrozen", ["0x152c42eda481651B6053f7f59bf617046c960EAB"])).toBe(
      "A basket the position holds is frozen: it backs no new loan or withdrawal until it moves again.",
    );
  });

  test("a request the wallet's owner declined", () => {
    expect(explain(new UserRejectedRequestError(new Error("denied")))).toBe(
      "You declined the request in your wallet.",
    );
  });

  test("a service's refusal over JSON-RPC is told in its own words", () => {
    const refusal = new InvalidParamsRpcError(new RpcRequestError({ body: {}, url: "http://127.0.0.1:4338", error: { code: -32000, message: "the call gas limit is above 1000000" } }));
    expect(explain(refusal)).toBe("Refused: the call gas limit is above 1000000.");
  });

  test("any other error keeps viem's short message, or the error's own", () => {
    expect(explain(new BaseError("The request took too long to respond.", { details: "timeout" }))).toBe(
      "The request took too long to respond.",
    );
    expect(explain(new Error("socket hang up"))).toBe("socket hang up");
  });
});
