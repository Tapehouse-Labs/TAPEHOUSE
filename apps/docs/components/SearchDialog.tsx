// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import MiniSearch from "minisearch";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

type Record_ = { id: string; url: string; page: string; heading: string; area: string; text: string };

/** The search box: a button in the header, ⌘K or / from anywhere, and a dialog over the whole index. */
export function SearchDialog() {
  const router = useRouter();
  const [open, setOpenState] = useState(false);
  const [query, setQuery] = useState("");
  const [cursor, setCursor] = useState(0);
  const [records, setRecords] = useState<Record_[] | null>(null);
  const input = useRef<HTMLInputElement>(null);
  const setOpen = useCallback((value: boolean) => {
    setOpenState(value);
    if (!value) {
      setQuery("");
      setCursor(0);
    }
  }, []);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const typing = (e.target as HTMLElement | null)?.closest("input, textarea, [contenteditable]");
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || (e.key === "/" && !typing)) {
        e.preventDefault();
        setOpen(true);
      } else if (e.key === "Escape") setOpen(false);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [setOpen]);

  useEffect(() => {
    if (!open || records) return;
    fetch("/search.json")
      .then((r) => r.json())
      .then((d: Record_[]) => setRecords(d))
      .catch(() => setRecords([]));
  }, [open, records]);

  useEffect(() => {
    if (open) input.current?.focus();
  }, [open]);

  const index = useMemo(() => {
    if (!records) return null;
    const mini = new MiniSearch<Record_>({
      fields: ["page", "heading", "text"],
      storeFields: ["url", "page", "heading", "area", "text"],
      searchOptions: { boost: { page: 3, heading: 2 }, prefix: true, fuzzy: 0.15, combineWith: "AND" },
    });
    mini.addAll(records);
    return mini;
  }, [records]);

  const hits = useMemo(() => {
    if (!index || query.trim().length < 2) return [];
    return index.search(query).slice(0, 12) as unknown as Record_[];
  }, [index, query]);

  const go = useCallback(
    (url: string) => {
      setOpen(false);
      router.push(url);
    },
    [router, setOpen],
  );

  return (
    <>
      <button type="button" className="pill min-w-0 flex-1 justify-between md:max-w-[22rem]" onClick={() => setOpen(true)} aria-label="Search the documentation">
        <span className="flex items-center gap-2 text-muted">
          <svg viewBox="0 0 16 16" aria-hidden="true" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth="1.5">
            <circle cx="7" cy="7" r="4.5" />
            <path d="M10.5 10.5 14 14" />
          </svg>
          <span>Search</span>
        </span>
        <span className="hidden gap-1 sm:flex">
          <kbd className="kbd">⌘</kbd>
          <kbd className="kbd">K</kbd>
        </span>
      </button>
      {open ? (
        <div className="search-overlay" role="presentation" onMouseDown={(e) => e.target === e.currentTarget && setOpen(false)}>
          <div className="search-panel" role="dialog" aria-modal="true" aria-label="Search">
            <input
              ref={input}
              value={query}
              onChange={(e) => {
                setQuery(e.target.value);
                setCursor(0);
              }}
              onKeyDown={(e) => {
                if (e.key === "ArrowDown") {
                  e.preventDefault();
                  setCursor((c) => Math.min(c + 1, hits.length - 1));
                } else if (e.key === "ArrowUp") {
                  e.preventDefault();
                  setCursor((c) => Math.max(c - 1, 0));
                } else if (e.key === "Enter" && hits[cursor]) go(hits[cursor].url);
              }}
              placeholder="Search the documentation"
              className="w-full border-b border-rule bg-transparent px-4 py-3.5 text-[15.5px] text-strong outline-none placeholder:text-faint"
              aria-label="Search the documentation"
              spellCheck={false}
              autoComplete="off"
            />
            <div className="overflow-y-auto p-2" role="listbox">
              {query.trim().length < 2 ? (
                <p className="px-3 py-6 text-center text-[14px] text-muted">Type two letters or more.</p>
              ) : !records ? (
                <p className="px-3 py-6 text-center text-[14px] text-muted">Loading the index.</p>
              ) : hits.length === 0 ? (
                <p className="px-3 py-6 text-center text-[14px] text-muted">Nothing matches “{query}”.</p>
              ) : (
                hits.map((h, i) => (
                  <a
                    key={h.id}
                    href={h.url}
                    role="option"
                    aria-selected={i === cursor}
                    className="search-hit"
                    onMouseEnter={() => setCursor(i)}
                    onClick={(e) => {
                      e.preventDefault();
                      go(h.url);
                    }}
                  >
                    <div className="flex items-baseline justify-between gap-3">
                      <span className="text-[14.5px] font-medium text-strong">{h.heading || h.page}</span>
                      <span className="shrink-0 font-mono text-[11.5px] uppercase tracking-wider text-faint">{h.area}</span>
                    </div>
                    {h.heading ? <div className="text-[12.5px] text-muted">{h.page}</div> : null}
                    <div className="mt-0.5 line-clamp-2 text-[13px] text-muted">{h.text.slice(0, 160)}</div>
                  </a>
                ))
              )}
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
