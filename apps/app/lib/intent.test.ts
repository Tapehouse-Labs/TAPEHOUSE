// SPDX-License-Identifier: MIT OR Apache-2.0
import { CROSS, marginAccountsAbi, parseDeployments, toBytes32 } from "@tapehouse/sdk";
import { decodeFunctionData, encodeFunctionData, erc20Abi } from "viem";
import { describe, expect, test } from "vitest";
import { smartAccountPlan } from "./intent";

const deployments = parseDeployments({
  chainId: 412346,
  tokens: {
    USDG: "0x4A2bA922052bA54e29c5417bC979Daaf7D5Fe4f4",
    WETH: "0x4Af567288e68caD4aA93A272fe6139Ca53859C70",
    NVDA: "0x525c2aBA45F66987217323E8a05EA400C65D06DC",
  },
  tapehouse: { MarginAccounts: "0xE85035F1145aC49333d105632a0d254E479a75bE" },
});
const account = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
const owner = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC";
const { USDG, WETH, NVDA } = deployments.tokens as Record<string, `0x${string}`>;
const accounts = deployments.tapehouse.MarginAccounts!;
const decoded = (call: { abi: unknown; functionName: string; args: readonly unknown[]; to: string }) => {
  const data = encodeFunctionData(call as never);
  return { to: call.to, ...decodeFunctionData({ abi: [...marginAccountsAbi, ...erc20Abi], data }) };
};

describe("smartAccountPlan", () => {
  test("a Stock Token deposit takes the tokens from the wallet by permit, approves and deposits in one operation", () => {
    const plan = smartAccountPlan(deployments, { kind: "deposit", position: CROSS, token: NVDA!, amount: 5n }, account, owner);
    expect(plan.wrap).toBe(0n);
    expect(plan.pulls).toEqual([{ token: NVDA, amount: 5n }]);
    expect(plan.calls.map(decoded)).toEqual([
      { to: NVDA, functionName: "approve", args: [accounts, 5n] },
      { to: accounts, functionName: "deposit", args: [CROSS, NVDA, 5n, account] },
    ]);
  });

  test("an ether deposit is wrapped into WETH for the account by the wallet, then deposited", () => {
    const plan = smartAccountPlan(deployments, { kind: "deposit", position: CROSS, token: "ETH", amount: 7n }, account, owner);
    expect(plan.wrap).toBe(7n);
    expect(plan.pulls).toEqual([]);
    expect(plan.calls.map(decoded).map((c) => [c.to, c.functionName])).toEqual([
      [WETH, "approve"],
      [accounts, "deposit"],
    ]);
  });

  test("a repayment takes the wallet's USDG by permit, and leaves no allowance", () => {
    const plan = smartAccountPlan(deployments, { kind: "repay", position: CROSS, amount: 9n }, account, owner);
    expect(plan.pulls).toEqual([{ token: USDG, amount: 9n }]);
    expect(plan.calls.map(decoded).slice(1)).toEqual([
      { to: accounts, functionName: "repay", args: [CROSS, 9n, account] },
      { to: USDG, functionName: "approve", args: [accounts, 0n] },
    ]);
  });

  test("a loan and a withdrawal go to the wallet", () => {
    const borrow = smartAccountPlan(deployments, { kind: "borrow", position: CROSS, amount: 3n }, account, owner);
    expect(borrow.calls.map(decoded)).toEqual([
      { to: accounts, functionName: "borrow", args: [CROSS, 3n, account, owner] },
    ]);
    const nvda = toBytes32("NVDA");
    const withdraw = smartAccountPlan(deployments, { kind: "withdraw", position: nvda, token: NVDA!, amount: 2n }, account, owner);
    expect(withdraw.calls.map(decoded)[0]?.args).toEqual([nvda, NVDA, 2n, account, owner]);
  });

  test("lending, taking back, recalling and settling act for the account", () => {
    for (const kind of ["lend", "unlend", "recall"] as const) {
      const plan = smartAccountPlan(deployments, { kind, position: CROSS, token: NVDA!, amount: 1n }, account, owner);
      expect(plan.calls.map(decoded)[0]).toEqual({ to: accounts, functionName: kind, args: [CROSS, NVDA, 1n, account] });
    }
    const settle = smartAccountPlan(deployments, { kind: "settle", position: CROSS, token: NVDA! }, account, owner);
    expect(settle.calls.map(decoded)[0]?.args).toEqual([account, CROSS, NVDA]);
  });
});
