// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useId } from "react";
import { shortAddress } from "@/lib/format";
import { useWallet } from "@/lib/hooks";
import { Wallets } from "./Wallets";

/** The header's wallet control: connect, or the connected address and a way out. */
export function WalletControl() {
  const wallet = useWallet();
  const menu = useId();
  if (wallet.state === "connecting") return <span className="skeleton inline-block h-[34px] w-[132px] rounded-full" />;
  if (wallet.state === "no-wallet") return null;
  if (wallet.state === "disconnected") return <Wallets compact />;
  return (
    <div className="relative">
      <button type="button" popoverTarget={menu} className="pill">
        <span
          aria-hidden="true"
          className={`h-2 w-2 rounded-full ${wallet.state === "connected" ? "bg-gain" : "bg-warning"}`}
        />
        <span className="addr num">{wallet.address ? shortAddress(wallet.address) : ""}</span>
      </button>
      <div
        id={menu}
        popover="auto"
        className="fixed right-4 top-16 w-64 rounded-[10px] bg-panel p-4 text-[14px] shadow-[0_18px_50px_-12px_rgb(0_0_0/0.35),0_0_0_1px_var(--rule)]"
      >
        <p className="label">Connected</p>
        <p className="addr num mt-1.5 break-all text-[12.5px] text-strong">{wallet.address}</p>
        <button type="button" onClick={wallet.disconnect} className="btn btn-quiet mt-4 w-full">
          Disconnect
        </button>
      </div>
    </div>
  );
}
