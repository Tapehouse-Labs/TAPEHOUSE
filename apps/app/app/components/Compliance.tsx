// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useAccountState, useMarginAccount } from "@/lib/accountHook";
import { notices } from "@/lib/compliance";

/**
 * What the visitor must know on connecting, under the header: that their address is sanctioned, and which Stock Tokens
 * the issuer has blocked it from moving. Nothing shows for a wallet with neither.
 */
export function Compliance() {
  const { account, owner } = useMarginAccount();
  const state = useAccountState(account, owner);
  const blocks = (state.data?.markets ?? []).map((m) => ({ asset: m.asset, blocked: m.issuer?.blocked ?? {} }));
  const list = notices([owner, account], blocks);
  if (list.length === 0) return null;
  return (
    <div role="alert" className="border-b border-rule bg-panel">
      <ul className="notes mx-auto max-w-[1320px] px-4 py-2.5 sm:px-8">
        {list.map((n) => (
          <li key={n.text} data-tone={n.tone}>
            {n.text}
          </li>
        ))}
      </ul>
    </div>
  );
}
