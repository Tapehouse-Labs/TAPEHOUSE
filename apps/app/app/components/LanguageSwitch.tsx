// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useLocale, useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { isLocale, LANGUAGE_NAMES, LOCALE_COOKIE, LOCALES } from "@/i18n/config";

const A_YEAR = 60 * 60 * 24 * 365;

/** The visitor's choice of language, kept in a cookie, and the page drawn again in it. */
export function LanguageSwitch() {
  const t = useTranslations("header");
  const locale = useLocale();
  const router = useRouter();
  const choose = (value: string) => {
    if (!isLocale(value)) return;
    document.cookie = `${LOCALE_COOKIE}=${value}; path=/; max-age=${A_YEAR}; samesite=lax`;
    router.refresh();
  };
  return (
    <label className="pill !px-2.5">
      <span className="sr-only">{t("language")}</span>
      <select
        value={locale}
        onChange={(event) => choose(event.target.value)}
        className="cursor-pointer bg-transparent text-[13px] text-strong focus:outline-none"
      >
        {LOCALES.map((code) => (
          <option key={code} value={code} lang={code}>
            {LANGUAGE_NAMES[code]}
          </option>
        ))}
      </select>
    </label>
  );
}
