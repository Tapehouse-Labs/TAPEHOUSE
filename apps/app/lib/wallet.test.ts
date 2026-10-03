// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { walletState } from "./wallet";

const expected = 46630;

describe("walletState", () => {
  test("no wallet in the browser and none connected", () => {
    expect(walletState({ status: "disconnected", chainId: undefined, expected, wallets: 0 })).toBe("no-wallet");
  });

  test("a wallet in the browser, not connected", () => {
    expect(walletState({ status: "disconnected", chainId: undefined, expected, wallets: 1 })).toBe("disconnected");
  });

  test("connecting and reconnecting read as connecting", () => {
    expect(walletState({ status: "connecting", chainId: undefined, expected, wallets: 1 })).toBe("connecting");
    expect(walletState({ status: "reconnecting", chainId: undefined, expected, wallets: 1 })).toBe("connecting");
  });

  test("connected on another chain", () => {
    expect(walletState({ status: "connected", chainId: 1, expected, wallets: 1 })).toBe("wrong-network");
  });

  test("connected on the expected chain", () => {
    expect(walletState({ status: "connected", chainId: expected, expected, wallets: 1 })).toBe("connected");
  });
});
