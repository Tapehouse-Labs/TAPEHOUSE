// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useState } from "react";

/** Copies `text` to the clipboard and says so for two seconds. */
export function CopyButton({ text, label = "Copy" }: { text: string; label?: string }) {
  const [done, setDone] = useState(false);
  return (
    <button
      type="button"
      className="copy"
      onClick={async () => {
        try {
          await navigator.clipboard.writeText(text);
          setDone(true);
          setTimeout(() => setDone(false), 2000);
        } catch {}
      }}
    >
      <svg viewBox="0 0 16 16" aria-hidden="true" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth="1.4">
        {done ? <path d="M3 8.5l3.2 3.2L13 4.8" /> : <><rect x="5.5" y="5.5" width="7.5" height="8" rx="1.5" /><path d="M10.5 3.5v-.5a1 1 0 0 0-1-1h-5.5a1 1 0 0 0-1 1v7.5a1 1 0 0 0 1 1h.5" /></>}
      </svg>
      <span aria-live="polite">{done ? "Copied" : label}</span>
    </button>
  );
}
