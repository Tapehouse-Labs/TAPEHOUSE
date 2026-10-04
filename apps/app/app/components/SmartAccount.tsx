// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { explain } from "@/lib/errors";
import { shortAddress } from "@/lib/format";
import { useSmartAccount } from "@/lib/smartAccountHook";
import { freeActions, notches } from "@/lib/smartAccount";

/**
 * The visitor's smart account under the figure: its address in a seal, and its free actions drawn as notches across
 * the stem at the bar's angle, each going out as it is spent.
 */
export function SmartAccount() {
  const smart = useSmartAccount();
  if (!smart.enabled) return null;
  const { address, standing } = smart;
  const left = standing?.left ?? 0n;
  return (
    <section aria-label="Smart account" className="relative mt-6">
      <span className="notches" data-pending={standing ? undefined : ""} aria-hidden="true">
        {notches(left).map((lit, i) => (
          <i key={i} data-spent={lit ? undefined : ""} />
        ))}
      </span>
      {address && standing ? (
        <>
          <p className="seal num">
            <span aria-hidden="true">[</span>
            <span className="addr" title={address}>
              {shortAddress(address)}
            </span>
            <span aria-hidden="true">]</span>
            <span className={standing.deployed ? "text-[var(--gain)]" : "text-muted"}>
              · {standing.deployed ? "LIVE" : "NOT CREATED"}
            </span>
          </p>
          <p className="mt-2 max-w-[44ch] text-[14px] text-body">{freeActions(standing)}</p>
        </>
      ) : (
        <>
          <span className="skeleton block h-[22px] w-[15rem]" />
          <span className="skeleton mt-3 block h-[14px] w-[19rem] max-w-full" />
        </>
      )}
      {!standing?.deployed && (
        <div className="mt-6 flex min-h-[44px] items-center">
          {standing && (
            <button type="button" onClick={smart.create} disabled={smart.phase !== undefined} className="btn">
              {smart.phase === "signing"
                ? "Waiting for your wallet…"
                : smart.phase === "creating"
                  ? "Creating your account…"
                  : "Create account"}
            </button>
          )}
        </div>
      )}
      {smart.error && (
        <p role="alert" className="mt-4 max-w-[52ch] text-[14px] text-[var(--error)]">
          {explain(smart.error)}
        </p>
      )}
    </section>
  );
}
