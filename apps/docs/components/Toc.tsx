// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useEffect, useState } from "react";
import type { Heading } from "@/lib/mdx";

/** The headings of the page, with the one being read marked. */
export function Toc({ headings }: { headings: Heading[] }) {
  const [active, setActive] = useState<string | null>(headings[0]?.id ?? null);
  useEffect(() => {
    const els = headings.map((h) => document.getElementById(h.id)).filter((e): e is HTMLElement => e !== null);
    if (els.length === 0) return;
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries.filter((e) => e.isIntersecting).sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
        if (visible[0]) setActive(visible[0].target.id);
      },
      { rootMargin: "-120px 0px -65% 0px" },
    );
    els.forEach((e) => observer.observe(e));
    return () => observer.disconnect();
  }, [headings]);
  if (headings.length < 2) return null;
  return (
    <nav aria-label="On this page">
      <div className="side-group-title !px-0">On this page</div>
      <ul className="border-l border-transparent">
        {headings.map((h) => (
          <li key={h.id}>
            <a href={`#${h.id}`} data-depth={h.depth} className="toc-link" aria-current={h.id === active ? "location" : undefined}>
              {h.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
