// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, erc20Abi, parseAbi, type PublicClient } from "viem";

/** The wrapper's part of a Stock Token the registry lists as a wrapper: what it wraps, and wrapping into it. */
export const wrapperAbi = parseAbi([
  "function underlying() view returns (address)",
  "function depositFor(address account, uint256 amount) returns (bool)",
]);

/** The Stock Token `token` wraps, or undefined where `token` is a Stock Token itself. */
export async function underlyingOf(client: PublicClient, token: Address): Promise<Address | undefined> {
  try {
    return await client.readContract({ address: token, abi: wrapperAbi, functionName: "underlying" });
  } catch {
    return undefined;
  }
}

/** What the wallet holds of `token` and of the Stock Token it wraps, which wraps into it one to one. */
export async function stockHeld(client: PublicClient, token: Address, owner: Address): Promise<bigint> {
  const held = await client.readContract({ address: token, abi: erc20Abi, functionName: "balanceOf", args: [owner] });
  const stock = await underlyingOf(client, token);
  if (!stock) return held;
  return held + (await client.readContract({ address: stock, abi: erc20Abi, functionName: "balanceOf", args: [owner] }));
}

/**
 * What the wallet wraps before it deposits `amount` of a wrapped Stock Token: the shortfall of what it holds of the
 * wrapper, and whether the Stock Token's allowance to the wrapper must be raised first. Undefined when it holds enough.
 */
export function wrapPlan({
  held,
  allowance,
  amount,
}: {
  held: bigint;
  allowance: bigint;
  amount: bigint;
}): { amount: bigint; approve: boolean } | undefined {
  if (held >= amount) return undefined;
  const missing = amount - held;
  return { amount: missing, approve: allowance < missing };
}
