// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import type { Item } from "@/lib/tree";

export function PrevNext({ prev, next }: { prev?: Item; next?: Item }) {
  if (!prev && !next) return null;
  return (
    <nav aria-label="Previous and next page" className="mt-16 grid gap-3 border-t border-rule pt-6 sm:grid-cols-2">
      {prev ? (
        <Link href={prev.href} className="card plain">
          <div className="font-mono text-[11.5px] uppercase tracking-wider text-faint">Previous</div>
          <div className="card-title mt-1">{prev.title}</div>
        </Link>
      ) : (
        <span />
      )}
      {next ? (
        <Link href={next.href} className="card plain sm:text-right">
          <div className="font-mono text-[11.5px] uppercase tracking-wider text-faint">Next</div>
          <div className="card-title mt-1">{next.title}</div>
        </Link>
      ) : null}
    </nav>
  );
}
