// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import dynamic from "next/dynamic";
import { useTranslations } from "next-intl";
import { useChains } from "wagmi";
import { useWallet } from "@/lib/hooks";
import { useExplain } from "@/lib/i18n";
import type { WalletState } from "@/lib/wallet";
import { Bands } from "./Bands";
import { Wallets } from "./Wallets";

const Desk = dynamic(() => import("./Desk").then((m) => m.Desk), { ssr: false, loading: () => <Figure /> });

const COPY: Record<WalletState, "noWallet" | "disconnected" | "connecting" | undefined> = {
  "no-wallet": "noWallet",
  disconnected: "disconnected",
  connecting: "connecting",
  "wrong-network": undefined,
  connected: undefined,
};

/** The liquidation price's place under the bar, before there is one. */
export function Figure() {
  const t = useTranslations("account.page");
  return (
    <div className="relative pt-14">
      <span className="bar" aria-hidden="true" />
      <h1 className="label">{t("figureLabel")}</h1>
      <p className="figure mt-3" data-empty>
        {t("none")}
      </p>
    </div>
  );
}

/** The account, drawn as the mark: the liquidation price sits where the bar crosses the stem. */
export function Account() {
  const t = useTranslations("account.page");
  const explain = useExplain();
  const wallet = useWallet();
  const [chain] = useChains();
  const copy = COPY[wallet.state];
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />

      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        {wallet.state === "connected" ? (
          <Desk />
        ) : (
          <>
            <Figure />
            <div key={wallet.state} className="swap" aria-live="polite">
              <p className="mt-6 max-w-[34ch] text-[19px] leading-snug text-body">
                {wallet.state === "wrong-network" ? (
                  t.rich("wrongNetwork", {
                    chain: chain.name,
                    strong: (c) => <span className="text-strong">{c}</span>,
                  })
                ) : (
                  copy && t(copy)
                )}
              </p>
              <div className="mt-8 flex min-h-[44px] flex-wrap items-center gap-3">
                {wallet.state === "wrong-network" && (
                  <button type="button" onClick={wallet.switchChain} disabled={wallet.pending} className="btn">
                    {wallet.pending ? t("waiting") : t("switchTo", { chain: chain.name })}
                  </button>
                )}
                {(wallet.state === "disconnected" || wallet.state === "no-wallet") && <Wallets />}
              </div>
              {wallet.error && (
                <p role="alert" className="mt-4 max-w-[52ch] text-[14px] text-[var(--error)]">
                  {explain(wallet.error)}
                </p>
              )}
            </div>
          </>
        )}

        <Bands />
      </div>
    </main>
  );
}
