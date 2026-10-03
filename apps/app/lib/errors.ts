// SPDX-License-Identifier: MIT OR Apache-2.0
import { decodeRevert } from "@tapehouse/sdk";
import { BaseError, UserRejectedRequestError } from "viem";

/** What went wrong, for the visitor: a revert by its error's name, a declined request, or viem's short message. */
export function explain(error: unknown): string {
  const revert = decodeRevert(error);
  if (revert) return `${revert.errorName}(${revert.args.map(String).join(", ")})`;
  if (error instanceof BaseError) {
    if (error.walk((cause) => cause instanceof UserRejectedRequestError))
      return "You declined the request in your wallet.";
    return error.shortMessage;
  }
  return error instanceof Error ? error.message : String(error);
}
