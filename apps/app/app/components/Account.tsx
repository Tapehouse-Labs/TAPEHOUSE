// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import dynamic from "next/dynamic";
import { useChains } from "wagmi";
import { explain } from "@/lib/errors";
import { useWallet } from "@/lib/hooks";
import type { WalletState } from "@/lib/wallet";
import { Bands } from "./Bands";
import { Wallets } from "./Wallets";

const SmartAccount = dynamic(() => import("./SmartAccount").then((m) => m.SmartAccount), { ssr: false });

const COPY: Record<WalletState, string> = {
  "no-wallet": "Add a browser wallet and your liquidation price appears here, before you sign anything.",
  disconnected: "Connect a wallet and your liquidation price appears here, before you sign anything.",
  connecting: "Waiting for your wallet.",
  "wrong-network": "",
  connected: "Your liquidation price appears here the moment you deposit.",
};

/** The account, drawn as the mark: the liquidation price sits where the bar crosses the stem. */
export function Account() {
  const wallet = useWallet();
  const [chain] = useChains();
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />

      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">Liquidation price</h1>
          <p className="figure mt-3" data-empty>
            None yet
          </p>
        </div>

        <div key={wallet.state} className="swap" aria-live="polite">
          <p className="mt-6 max-w-[34ch] text-[19px] leading-snug text-body">
            {wallet.state === "wrong-network" ? (
              <>
                Your wallet is on another network. Tapehouse runs on <span className="text-strong">{chain.name}</span>.
              </>
            ) : (
              COPY[wallet.state]
            )}
          </p>
          {wallet.state === "connected" ? (
            <SmartAccount />
          ) : (
            <div className="mt-8 flex min-h-[44px] flex-wrap items-center gap-3">
              {wallet.state === "wrong-network" && (
                <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
                  {wallet.pending ? "Waiting for your wallet…" : `Switch to ${chain.name}`}
                </button>
              )}
              {(wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />}
            </div>
          )}
          {wallet.error && (
            <p role="alert" className="mt-4 max-w-[52ch] text-[14px] text-[var(--error)]">
              {explain(wallet.error)}
            </p>
          )}
        </div>

        <Bands />
      </div>
    </main>
  );
}
