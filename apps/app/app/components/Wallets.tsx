// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useId } from "react";
import type { Connector } from "wagmi";
import { useWallet } from "@/lib/hooks";

const FIND_A_WALLET = "https://ethereum.org/en/wallets/find-wallet/";

/** Connects one of the browser's wallets: directly when there is one, from a list when there are several. */
export function Wallets({ compact = false }: { compact?: boolean }) {
  const wallet = useWallet();
  const list = useId();
  const size = compact ? "pill" : "btn";

  if (wallet.wallets.length === 0) {
    return (
      <a href={FIND_A_WALLET} target="_blank" rel="noreferrer" className={size}>
        Get a wallet
        <svg viewBox="0 0 12 12" aria-hidden="true" className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth="1.5">
          <path d="M4 2.5h5.5V8M9.5 2.5 2.5 9.5" />
        </svg>
        <span className="sr-only">(opens in a new tab)</span>
      </a>
    );
  }
  if (wallet.wallets.length === 1) {
    const only = wallet.wallets[0];
    return (
      <button type="button" onClick={() => wallet.connect(only)} disabled={wallet.pending} className={size}>
        {wallet.pending ? "Waiting for your wallet…" : "Connect wallet"}
      </button>
    );
  }
  return (
    <div className="relative">
      <button type="button" popoverTarget={list} disabled={wallet.pending} className={size}>
        {wallet.pending ? "Waiting for your wallet…" : "Connect wallet"}
      </button>
      <ul
        id={list}
        popover="auto"
        aria-label="Wallets in this browser"
        className="fixed right-4 top-16 w-64 rounded-[10px] bg-panel p-2 shadow-[0_18px_50px_-12px_rgb(0_0_0/0.35),0_0_0_1px_var(--rule)]"
      >
        {wallet.wallets.map((connector) => (
          <li key={connector.uid}>
            <WalletOption connector={connector} onPick={() => wallet.connect(connector)} target={list} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function WalletOption({ connector, onPick, target }: { connector: Connector; onPick: () => void; target: string }) {
  return (
    <button
      type="button"
      popoverTarget={target}
      popoverTargetAction="hide"
      onClick={onPick}
      className="flex w-full items-center gap-3 rounded-[6px] px-3 py-2.5 text-left text-[14px] text-strong hover:bg-sunk"
    >
      {connector.icon ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={connector.icon} alt="" width={22} height={22} className="rounded-[5px]" />
      ) : (
        <span aria-hidden="true" className="h-[22px] w-[22px] rounded-[5px] bg-sunk" />
      )}
      {connector.name}
    </button>
  );
}
