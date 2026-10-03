// SPDX-License-Identifier: MIT OR Apache-2.0
import { errorsAbi } from "@tapehouse/sdk";
import {
  BaseError,
  ContractFunctionExecutionError,
  ContractFunctionRevertedError,
  encodeErrorResult,
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

  test("a request the wallet's owner declined", () => {
    expect(explain(new UserRejectedRequestError(new Error("denied")))).toBe(
      "You declined the request in your wallet.",
    );
  });

  test("any other error keeps viem's short message, or the error's own", () => {
    expect(explain(new BaseError("The request took too long to respond.", { details: "timeout" }))).toBe(
      "The request took too long to respond.",
    );
    expect(explain(new Error("socket hang up"))).toBe("socket hang up");
  });
});
