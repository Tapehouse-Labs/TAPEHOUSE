// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import { tabs, type Tab } from "@/lib/nav";
import { APP_URL, GITHUB_URL, HOME_URL } from "@/lib/site";
import { LiveHead } from "./LiveHead";
import { Mark } from "./Mark";
import { MobileNav } from "./MobileNav";
import { SearchDialog } from "./SearchDialog";
import { ThemeToggle } from "./ThemeToggle";

/** The site's chrome: the mark, the search, the tabs. `side` is the current tab's sidebar, for phones. */
export function Header({ current, side }: { current: Tab; side: React.ReactNode }) {
  return (
    <header className="chrome">
      <div className="shell flex h-[var(--header-h)] items-center gap-3 sm:gap-5">
        <MobileNav>{side}</MobileNav>
        <Link href="/" className="flex shrink-0 items-center gap-3 text-accent" aria-label="Tapehouse documentation, home">
          <Mark className="h-[20px] w-auto" />
          <span className="font-display text-[14.5px] font-semibold tracking-[0.18em] text-strong">TAPEHOUSE</span>
          <span className="hidden border-l border-rule pl-3 text-[14px] text-muted sm:inline">Docs</span>
        </Link>
        <div className="flex min-w-0 flex-1 justify-center md:justify-start md:pl-6">
          <SearchDialog />
        </div>
        <div className="flex shrink-0 items-center gap-2 sm:gap-3">
          <LiveHead />
          <a href={GITHUB_URL} className="pill hidden sm:inline-flex" rel="noreferrer">
            GitHub
          </a>
          <ThemeToggle />
          <a href={APP_URL} className="pill pill-solid hidden md:inline-flex" rel="noreferrer">
            Open the app
          </a>
        </div>
      </div>
      <nav aria-label="Sections" className="shell flex h-[var(--tabs-h)] overflow-x-auto">
        {tabs.map((t) => (
          <Link key={t.id} href={t.home} className="tab-link" aria-current={t.id === current.id ? "page" : undefined}>
            {t.title}
          </Link>
        ))}
        <a href={HOME_URL} className="tab-link ml-auto hidden text-faint lg:inline-flex" rel="noreferrer">
          tapehouse.xyz
        </a>
      </nav>
    </header>
  );
}
