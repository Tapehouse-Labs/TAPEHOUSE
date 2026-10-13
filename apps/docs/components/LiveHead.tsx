// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useEffect, useState } from "react";
import { INDEXER_URL } from "@/lib/site";

/** The testnet's index head, read from the public API: the docs' own proof that what they describe is running. */
export function LiveHead() {
  const [head, setHead] = useState<number | null>(null);
  useEffect(() => {
    let alive = true;
    const read = () =>
      fetch(`${INDEXER_URL}/v1/status`)
        .then((r) => (r.ok ? r.json() : Promise.reject(new Error(String(r.status)))))
        .then((s: { head: { number: number } }) => alive && setHead(s.head.number))
        .catch(() => alive && setHead(null));
    read();
    const t = setInterval(read, 30_000);
    return () => {
      alive = false;
      clearInterval(t);
    };
  }, []);
  return (
    <span className="hidden items-center gap-2 text-[12.5px] text-muted lg:inline-flex" title="Robinhood Chain testnet, chain 46630">
      <span aria-hidden="true" className="h-1.5 w-1.5 rounded-full" style={{ background: head ? "var(--gain)" : "var(--text-faint)" }} />
      <span>Testnet</span>
      {head ? <span className="num addr text-faint">{head.toLocaleString("en-US")}</span> : null}
    </span>
  );
}
