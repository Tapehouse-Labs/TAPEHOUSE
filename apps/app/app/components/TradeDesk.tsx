// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { type ReactNode, useEffect, useMemo } from "react";
import type { Address } from "viem";
import { PHASE, useMarginAccount, useSend, useWalletBalances } from "@/lib/accountHook";
import { explain } from "@/lib/errors";
import type { Intent } from "@/lib/intent";
import { useTradeState } from "@/lib/tradeHook";
import { SmartAccount } from "./SmartAccount";
import type { Desk } from "./Trade";

/**
 * The connected account's part of the Trade surface, handed up to the ticket once this chunk loads: who trades, what
 * the wallet holds of USDG and the baskets' Stock Tokens, the button that sends, and the smart account's block.
 */
export function TradeDesk({ onDesk }: { onDesk: (desk: { desk: Desk; account: ReactNode } | undefined) => void }) {
  const { account, owner, smart } = useMarginAccount();
  const state = useTradeState(account);
  const tokens = useMemo(() => [...new Set((state.data?.baskets ?? []).flatMap((b) => b.tokens))], [state.data]);
  const balances = useWalletBalances(owner, tokens);
  const sponsored = smart.enabled && (smart.standing?.left ?? 0n) > 0n;
  const held = balances.data?.tokens;
  useEffect(() => {
    onDesk(
      account
        ? {
            desk: {
              account,
              owner: owner ?? account,
              balances: held ?? ({} as Record<Address, bigint>),
              send: (intent, label) => <SendButton intent={intent} label={label} sponsored={sponsored} />,
            },
            account: <SmartAccount />,
          }
        : undefined,
    );
  }, [account, owner, held, sponsored, onDesk]);
  useEffect(() => () => onDesk(undefined), [onDesk]);
  return null;
}

function SendButton({ intent, label, sponsored }: { intent: Intent | undefined; label: string; sponsored: boolean }) {
  const send = useSend();
  return (
    <span className="inline-flex flex-col items-start gap-2">
      <button type="button" className="btn" disabled={!intent || send.pending} onClick={() => intent && send.send(intent)}>
        {send.phase ? PHASE[send.phase] : `${label}${sponsored ? ", free" : ""}`}
      </button>
      {send.error && (
        <span role="alert" className="max-w-[60ch] text-[14px] text-[var(--error)]">
          {explain(send.error)}
        </span>
      )}
      {send.succeeded && !send.pending && (
        <span role="status" className="text-[14px] text-[var(--gain)]">
          Done. Every line is read again from the chain.
        </span>
      )}
    </span>
  );
}
