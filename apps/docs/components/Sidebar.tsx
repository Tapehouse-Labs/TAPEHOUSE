// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import type { Section } from "@/lib/tree";

/** One tab's pages, grouped; the current page is marked for assistive technology as well as for the eye. */
export function Sidebar({ sections, current }: { sections: Section[]; current: string }) {
  return (
    <nav aria-label="Pages" className="space-y-7">
      {sections.map((s) => (
        <div key={s.title}>
          <div className="side-group-title">{s.title}</div>
          <ul className="space-y-0.5">
            {s.items.map((i) => (
              <li key={i.href}>
                <Link href={i.href} className="side-link" aria-current={i.slug === current ? "page" : undefined}>
                  {i.method ? <span className="side-method">{i.method}</span> : null}
                  {i.title}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </nav>
  );
}
