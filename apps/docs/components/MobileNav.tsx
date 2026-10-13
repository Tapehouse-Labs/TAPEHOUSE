// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { usePathname } from "next/navigation";
import { useEffect, useState, type ReactNode } from "react";

/** The sidebar on a phone: a button in the header and a drawer that closes when the page changes. */
export function MobileNav({ children }: { children: ReactNode }) {
  const path = usePathname();
  const [openOn, setOpenOn] = useState<string | null>(null);
  const open = openOn === path;
  const setOpen = (value: boolean) => setOpenOn(value ? path : null);
  useEffect(() => {
    document.body.style.overflow = open ? "hidden" : "";
    return () => {
      document.body.style.overflow = "";
    };
  }, [open]);
  return (
    <>
      <button type="button" className="pill !p-2.5 lg:hidden" aria-label="Open the menu" aria-expanded={open} onClick={() => setOpen(true)}>
        <svg viewBox="0 0 16 16" aria-hidden="true" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.5">
          <path d="M2.5 4.5h11M2.5 8h11M2.5 11.5h11" />
        </svg>
      </button>
      {open ? (
        <div className="fixed inset-0 z-50 lg:hidden" role="dialog" aria-modal="true" aria-label="Menu">
          <button type="button" className="absolute inset-0 bg-ink/40" aria-label="Close the menu" onClick={() => setOpen(false)} />
          <div className="absolute inset-y-0 left-0 w-[19rem] max-w-[85vw] overflow-y-auto bg-paper p-4 shadow-xl">{children}</div>
        </div>
      ) : null}
    </>
  );
}
