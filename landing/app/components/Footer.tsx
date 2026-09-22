import Link from "next/link";
import { APP_URL, CONTACT_EMAIL, EXTERNAL, SECURITY_EMAIL } from "../site";
import { Mark } from "./Mark";

type Item = { label: string; href: string; external?: boolean };

const COLUMNS: { title: string; items: Item[] }[] = [
  {
    title: "Product",
    items: [
      { label: "Open the app", href: APP_URL },
      { label: "How capacity works", href: "/#ch-capacity" },
      { label: "The weekend, measured", href: "/#ch-weekend" },
      { label: "The record", href: "/#ch-record" },
    ],
  },
  {
    title: "Builders",
    items: [
      { label: "Integrate the feed", href: "/#ch-method" },
      { label: "Methodology", href: "/methodology" },
      { label: "Report a vulnerability", href: `mailto:${SECURITY_EMAIL}`, external: true },
    ],
  },
  {
    title: "Tapehouse",
    items: [
      { label: "Risk disclosure", href: "/risk" },
      ...EXTERNAL.filter((l): l is { label: string; href: string } => Boolean(l.href)).map((l) => ({
        ...l,
        external: true,
      })),
    ],
  },
];

export function Footer() {
  return (
    <footer className="field field-night">
      <div className="shell pb-10 pt-[clamp(4.5rem,9vw,7.5rem)]">
        <div className="grid gap-x-16 gap-y-14 lg:grid-cols-[1.15fr_1.6fr]">
          <div>
            <Link href="/" className="inline-flex items-center gap-3.5 text-strong" aria-label="Tapehouse, home">
              <Mark className="h-7 w-auto" />
              <span className="font-display text-[1.35rem] font-semibold tracking-[0.17em]">TAPEHOUSE</span>
            </Link>
            <p className="mt-4 max-w-[30ch] text-[1.02rem] leading-snug text-body">
              Portfolio margin for Stock Tokens.
            </p>
            {CONTACT_EMAIL && (
              <a href={`mailto:${CONTACT_EMAIL}`} className="mt-6 inline-block text-[15px] text-strong underline">
                {CONTACT_EMAIL}
              </a>
            )}
          </div>

          <nav aria-label="Footer" className="grid grid-cols-2 gap-x-8 gap-y-10 sm:grid-cols-3">
            {COLUMNS.map((col) => (
              <div key={col.title}>
                <h2 className="text-[13px] text-faint">{col.title}</h2>
                <ul className="mt-4 space-y-2.5 text-[15px]">
                  {col.items.map((item) => (
                    <li key={item.label}>
                      {item.external ? (
                        <a
                          href={item.href}
                          {...(item.href.startsWith("mailto:") ? {} : { target: "_blank", rel: "noreferrer" })}
                          className="text-body transition-colors duration-200 hover:text-strong"
                        >
                          {item.label}
                        </a>
                      ) : (
                        <Link href={item.href} className="text-body transition-colors duration-200 hover:text-strong">
                          {item.label}
                        </Link>
                      )}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </nav>
        </div>

        <div className="mt-[clamp(3.5rem,7vw,6rem)] border-t border-rule pt-7">
          <p className="max-w-[108ch] text-[12.5px] leading-[1.6] text-faint">
            Tapehouse is software, not a broker or a bank, and it holds no licence to provide financial
            services. Nothing here is investment advice or an offer of credit, and today it runs on a test
            network, where assets have no value. A Stock Token is a debt security of its issuer, not a share:
            it follows the price and gives you no shareholder rights. Borrowing against a portfolio can lose
            you part or all of it, including while the underlying market is closed.{" "}
            <Link href="/risk" className="text-muted underline hover:text-strong">
              Read the full risk disclosure
            </Link>
            .
          </p>
          <p className="mt-6 flex flex-wrap justify-between gap-x-8 gap-y-2 text-[12.5px] text-faint">
            <span>© 2026 Tapehouse</span>
            <span>Runs on Robinhood Chain testnet</span>
          </p>
        </div>
      </div>
    </footer>
  );
}
