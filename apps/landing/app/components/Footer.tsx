// SPDX-FileCopyrightText: Tapehouse contributors
// SPDX-License-Identifier: MIT OR Apache-2.0
import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { APP_URL, CONTACT_EMAIL, DOCS_URL, EXTERNAL, SECURITY_EMAIL } from "../site";
import { LanguageSwitch } from "./LanguageSwitch";
import { Mark } from "./Mark";

type Item = { label: string; href: string; external?: boolean };

export function Footer() {
  const t = useTranslations("site.footer");
  const columns: { title: string; items: Item[] }[] = [
    {
      title: t("product"),
      items: [
        { label: t("openApp"), href: APP_URL },
        { label: t("capacity"), href: "/#ch-capacity" },
        { label: t("weekend"), href: "/#ch-weekend" },
        { label: t("record"), href: "/#ch-record" },
      ],
    },
    {
      title: t("builders"),
      items: [
        { label: t("docs"), href: DOCS_URL, external: true },
        { label: t("feed"), href: "/#ch-method" },
        { label: t("methodology"), href: "/methodology" },
        { label: t("vulnerability"), href: `mailto:${SECURITY_EMAIL}`, external: true },
      ],
    },
    {
      title: "Tapehouse",
      items: [
        { label: t("risk"), href: "/risk" },
        ...EXTERNAL.filter((l): l is { label: string; href: string } => Boolean(l.href)).map((l) => ({
          ...l,
          external: true,
        })),
      ],
    },
  ];
  return (
    <footer className="field field-night">
      <div className="shell pb-10 pt-[clamp(4.5rem,9vw,7.5rem)]">
        <div className="grid gap-x-16 gap-y-14 lg:grid-cols-[1.15fr_1.6fr]">
          <div>
            <Link href="/" className="inline-flex items-center gap-3.5 text-strong" aria-label={t("home")}>
              <Mark className="h-7 w-auto" />
              <span className="font-display text-[1.35rem] font-semibold tracking-[0.17em]">TAPEHOUSE</span>
            </Link>
            <p className="mt-4 max-w-[30ch] text-[1.02rem] leading-snug text-body">{t("tagline")}</p>
            {CONTACT_EMAIL && (
              <a href={`mailto:${CONTACT_EMAIL}`} className="mt-6 inline-block text-[15px] text-strong underline">
                {CONTACT_EMAIL}
              </a>
            )}
          </div>

          <nav aria-label={t("aria")} className="grid grid-cols-2 gap-x-8 gap-y-10 sm:grid-cols-3">
            {columns.map((col) => (
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
            {t("legal")}{" "}
            <Link href="/risk" className="text-muted underline hover:text-strong">
              {t("readRisk")}
            </Link>
            .
          </p>
          <div className="mt-6 flex flex-wrap items-center justify-between gap-x-8 gap-y-3 text-[12.5px] text-faint">
            <span>© 2026 Tapehouse</span>
            <LanguageSwitch />
            <span>{t("network")}</span>
          </div>
        </div>
      </div>
    </footer>
  );
}
