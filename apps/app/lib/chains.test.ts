// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { chainOf, faucetUrl } from "./chains";

describe("chainOf", () => {
  test("Robinhood Chain and its testnet keep viem's definitions", () => {
    expect(chainOf(4663).name).toBe("Robinhood Chain");
    expect(chainOf(46630)).toMatchObject({ name: "Robinhood Chain Testnet", testnet: true });
  });

  test("the Nitro dev node", () => {
    expect(chainOf(412346)).toMatchObject({ id: 412346, rpcUrls: { default: { http: ["http://127.0.0.1:8547"] } } });
  });

  test("an RPC URL replaces the chain's default", () => {
    expect(chainOf(46630, "https://rpc.example").rpcUrls.default.http).toEqual(["https://rpc.example"]);
  });

  test("refuses a chain Tapehouse is not deployed on", () => {
    expect(() => chainOf(1)).toThrow("Tapehouse runs on no chain 1.");
  });
});

describe("faucetUrl", () => {
  test("only Robinhood Chain's testnet has a USDG faucet", () => {
    expect(faucetUrl(46630)).toBe("https://faucet.paxos.com/?network=robinhood");
    expect(faucetUrl(4663)).toBeUndefined();
    expect(faucetUrl(412346)).toBeUndefined();
  });
});
