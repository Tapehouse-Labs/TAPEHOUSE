// SPDX-License-Identifier: MIT OR Apache-2.0
import { tabOf, type Tab } from "@/lib/nav";
import { tree } from "@/lib/tree";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";

/** The page's frame: the chrome, the current tab's sidebar and a column for the page. */
export function Shell({ slug, tab, aside, wide, children }: { slug: string; tab?: Tab; aside?: React.ReactNode; wide?: boolean; children: React.ReactNode }) {
  const current = tab ?? tabOf(slug);
  const sections = tree(current);
  const side = <Sidebar sections={sections} current={slug} />;
  return (
    <>
      <Header current={current} side={side} />
      <div className="shell lg:grid lg:grid-cols-[16.5rem_minmax(0,1fr)] lg:gap-10">
        <div className="hidden lg:block">
          <div className="sticky top-[var(--chrome-h)] max-h-[calc(100vh-var(--chrome-h))] overflow-y-auto py-8 pr-2">{side}</div>
        </div>
        <div className={wide ? "min-w-0 py-8 lg:py-10" : "min-w-0 py-8 lg:grid lg:grid-cols-[minmax(0,46rem)_14rem] lg:justify-between lg:gap-12 lg:py-10"}>
          <main id="main" className="min-w-0">
            {children}
          </main>
          {wide ? null : (
            <div className="hidden lg:block">
              <div className="sticky top-[calc(var(--chrome-h)+2rem)] max-h-[calc(100vh-var(--chrome-h)-3rem)] overflow-y-auto">{aside}</div>
            </div>
          )}
        </div>
      </div>
      <footer className="border-t border-rule">
        <div className="shell flex flex-wrap items-center justify-between gap-3 py-6 text-[13px] text-muted">
          <span>Testnet only. Not audited. Stock Tokens are not offered to United States persons.</span>
          <span>MIT OR Apache-2.0</span>
        </div>
      </footer>
    </>
  );
}
