// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useLocale, useTranslations } from "next-intl";
import { Link, usePathname } from "@/i18n/navigation";
import { LANGUAGE_NAMES, LOCALES } from "@/i18n/routing";

/** The site's languages, each a link to the same page in that language. */
export function LanguageSwitch({ className = "", onDark = false }: { className?: string; onDark?: boolean }) {
  const t = useTranslations("site.language");
  const current = useLocale();
  const pathname = usePathname();
  return (
    <nav aria-label={t("aria")} className={`flex items-center gap-1 text-[13px] ${className}`}>
      {LOCALES.map((locale) => (
        <Link
          key={locale}
          href={pathname}
          locale={locale}
          lang={locale}
          hrefLang={locale}
          title={LANGUAGE_NAMES[locale]}
          aria-current={locale === current ? "true" : undefined}
          className={`inline-flex min-h-6 min-w-6 items-center justify-center px-1.5 ${
            onDark
              ? locale === current
                ? "text-[rgb(245_242_236)] underline underline-offset-4"
                : "text-[rgb(245_242_236/0.7)] hover:text-[rgb(245_242_236)]"
              : locale === current
                ? "text-strong underline underline-offset-4"
                : "text-muted hover:text-strong"
          }`}
        >
          {locale.toUpperCase()}
        </Link>
      ))}
    </nav>
  );
}
