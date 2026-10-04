// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useChains } from "wagmi";
import { faucetUrl } from "@/lib/chains";
import { Mark } from "./Mark";
import { Seal } from "./Seal";
import { ThemeToggle } from "./ThemeToggle";
import { WalletControl } from "./WalletControl";

/** The app's sections. Each surface adds itself here as it ships. */
export const SECTIONS = [
  { href: "/", label: "Account" },
  { href: "/earn", label: "Earn" },
  { href: "/simulator", label: "Simulator" },
] as const;

export function Header() {
  const pathname = usePathname();
  const [chain] = useChains();
  const faucet = faucetUrl(chain.id);
  return (
    <>
      <header className="sticky top-0 z-40 border-b border-rule bg-paper/90 backdrop-blur-md">
        <div className="mx-auto flex h-16 max-w-[1320px] items-center justify-between gap-4 px-4 sm:px-8">
          <Link href="/" className="flex items-center gap-2.5 text-strong" aria-label="Tapehouse, account">
            <Mark className="h-[18px] w-auto text-accent" />
            <span className="font-display text-[13px] font-semibold tracking-[0.18em]">TAPEHOUSE</span>
          </Link>
          <nav aria-label="Sections" className="hidden md:block">
            <ul className="flex gap-1">
              {SECTIONS.map((s) => (
                <li key={s.href}>
                  <Link
                    href={s.href}
                    prefetch={false}
                    aria-current={pathname === s.href ? "page" : undefined}
                    className="rounded-full px-3.5 py-2 text-[14px] text-muted hover:text-strong aria-[current=page]:bg-panel aria-[current=page]:text-strong aria-[current=page]:shadow-[inset_0_0_0_1px_var(--rule)]"
                  >
                    {s.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
          <div className="flex items-center gap-2 sm:gap-2.5">
            <span className="hidden sm:inline-flex">
              <Seal />
            </span>
            <span className="pill hidden md:inline-flex">
              <span aria-hidden="true" className="h-1.5 w-1.5 rounded-full bg-gain" />
              {chain.name}
            </span>
            {faucet && (
              <a href={faucet} target="_blank" rel="noopener noreferrer" className="pill hidden md:inline-flex">
                Get test USDG
              </a>
            )}
            <WalletControl />
            <ThemeToggle />
          </div>
        </div>
      </header>
      <nav
        aria-label="Sections"
        className="fixed inset-x-0 bottom-0 z-40 border-t border-rule bg-paper/95 pb-[env(safe-area-inset-bottom)] backdrop-blur-md md:hidden"
      >
        <ul className="mx-auto flex max-w-[480px] justify-center gap-1 p-2">
          {SECTIONS.map((s) => (
            <li key={s.href} className="flex-1">
              <Link
                href={s.href}
                prefetch={false}
                aria-current={pathname === s.href ? "page" : undefined}
                className="flex min-h-[44px] items-center justify-center rounded-full text-[14px] text-muted aria-[current=page]:bg-panel aria-[current=page]:text-strong aria-[current=page]:shadow-[inset_0_0_0_1px_var(--rule)]"
              >
                {s.label}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </>
  );
}
