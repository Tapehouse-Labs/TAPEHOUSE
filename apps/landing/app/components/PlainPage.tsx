import Link from "next/link";
import { Footer } from "./Footer";
import { Mark } from "./Mark";

/** The plain register: for pages that are read, not sold. */
export function PlainPage({ title, lede, children }: { title: string; lede: string; children: React.ReactNode }) {
  return (
    <>
      <div className="field field-bone min-h-[70svh]">
        <header className="shell flex items-center justify-between gap-4 pt-6">
          <Link href="/" className="flex items-center gap-3 text-accent" aria-label="Tapehouse, home">
            <Mark className="h-[22px] w-auto" />
            <span className="font-display text-[15px] font-semibold tracking-[0.18em] text-strong">TAPEHOUSE</span>
          </Link>
          <Link href="/" className="text-[15px] text-muted underline hover:text-strong">
            Back to the overview
          </Link>
        </header>
        <main className="shell pb-[clamp(5rem,10vw,9rem)] pt-[clamp(4rem,9vw,8rem)]">
          <h1 className="chapter-title max-w-[18ch]">{title}</h1>
          <p className="lede mt-7 max-w-[58ch]">{lede}</p>
          <div className="mt-[clamp(3rem,6vw,5rem)] max-w-[68ch] space-y-10 text-[1.02rem] leading-[1.65]">
            {children}
          </div>
        </main>
      </div>
      <Footer />
    </>
  );
}

export function Entry({ heading, children }: { heading: string; children: React.ReactNode }) {
  return (
    <section>
      <h2 className="font-display text-[1.2rem] font-semibold tracking-[-0.01em] text-strong">{heading}</h2>
      <div className="mt-3 space-y-4 text-body">{children}</div>
    </section>
  );
}
