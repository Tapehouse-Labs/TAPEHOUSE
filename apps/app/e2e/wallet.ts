// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Page } from "@playwright/test";
import type { Hex } from "viem";
import { privateKeyToAccount } from "viem/accounts";

export type TestWallet = {
  /** The account the wallet holds. */
  account: string;
  /** The chain the wallet starts on. */
  chainId: number;
  /** The node the wallet forwards reads to. */
  rpcUrl: string;
  /** Whether the wallet's owner declines the connection. */
  decline?: boolean;
  /** The key the wallet signs messages with, if it signs at all; `account` must be its address. */
  key?: Hex;
};

/** Moves the page's test wallet to `chainId`. */
export async function moveWallet(page: Page, chainId: number) {
  await page.evaluate((id) => (window as unknown as { testWallet: { moveTo(id: number): void } }).testWallet.moveTo(id), chainId);
}

/**
 * Installs a browser wallet in the page before it loads: an EIP-1193 provider announced through EIP-6963, which
 * holds `account`, starts on `chainId`, switches to the chains it is asked for and forwards every other request to
 * `rpcUrl`. `window.testWallet.moveTo(chainId)` moves it to another chain, as its owner would. A test double: it
 * signs messages (`personal_sign`) with `key` when it is given one, and nothing else.
 */
export async function installWallet(page: Page, wallet: TestWallet) {
  if (wallet.key) {
    const signer = privateKeyToAccount(wallet.key);
    if (signer.address.toLowerCase() !== wallet.account.toLowerCase()) throw new Error("The key does not hold account.");
    await page.exposeFunction("testWalletSign", (message: Hex) => signer.signMessage({ message: { raw: message } }));
  }
  await page.addInitScript((w: TestWallet) => {
    let chainId = w.chainId;
    let connected = false;
    const listeners = new Map<string, Set<(value: unknown) => void>>();
    const emit = (event: string, value: unknown) => listeners.get(event)?.forEach((fn) => fn(value));
    const hex = (n: number) => `0x${n.toString(16)}`;
    const provider = {
      async request({ method, params }: { method: string; params?: unknown[] }) {
        switch (method) {
          case "eth_requestAccounts":
            if (w.decline) throw Object.assign(new Error("User rejected the request."), { code: 4001 });
            connected = true;
            emit("connect", { chainId: hex(chainId) });
            emit("accountsChanged", [w.account]);
            return [w.account];
          case "eth_accounts":
            return connected ? [w.account] : [];
          case "eth_chainId":
            return hex(chainId);
          case "personal_sign": {
            const sign = (window as unknown as { testWalletSign?: (message: string) => Promise<string> }).testWalletSign;
            if (!sign) throw Object.assign(new Error("The test wallet signs nothing."), { code: 4200 });
            return sign((params as [string, string])[0]);
          }
          case "wallet_requestPermissions":
            return [{ parentCapability: "eth_accounts" }];
          case "wallet_switchEthereumChain": {
            const [{ chainId: next }] = params as [{ chainId: string }];
            chainId = Number(next);
            emit("chainChanged", next);
            return null;
          }
          default: {
            const response = await fetch(w.rpcUrl, {
              method: "POST",
              headers: { "content-type": "application/json" },
              body: JSON.stringify({ jsonrpc: "2.0", id: 1, method, params: params ?? [] }),
            });
            const { result, error } = await response.json();
            if (error) throw Object.assign(new Error(error.message), { code: error.code, data: error.data });
            return result;
          }
        }
      },
      on(event: string, fn: (value: unknown) => void) {
        if (!listeners.has(event)) listeners.set(event, new Set());
        listeners.get(event)!.add(fn);
      },
      removeListener(event: string, fn: (value: unknown) => void) {
        listeners.get(event)?.delete(fn);
      },
    };
    const info = { uuid: "2f3b1c9a-0000-4000-8000-000000000001", name: "Test Wallet", icon: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg'/>", rdns: "xyz.tapehouse.test" };
    const announce = () =>
      window.dispatchEvent(new CustomEvent("eip6963:announceProvider", { detail: Object.freeze({ info, provider }) }));
    Object.assign(window, {
      testWallet: {
        moveTo(next: number) {
          chainId = next;
          emit("chainChanged", hex(next));
        },
      },
    });
    window.addEventListener("eip6963:requestProvider", announce);
    announce();
  }, wallet);
}
