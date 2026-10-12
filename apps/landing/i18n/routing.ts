// SPDX-License-Identifier: MIT OR Apache-2.0
import { defineRouting } from "next-intl/routing";

/** The languages the site is written in, English first: the one served without a prefix. */
export const LOCALES = ["en", "es", "pt"] as const;
export type Locale = (typeof LOCALES)[number];

export const routing = defineRouting({ locales: LOCALES, defaultLocale: "en", localePrefix: "as-needed" });

/** Each language's name in itself, as its speaker looks for it. */
export const LANGUAGE_NAMES: Readonly<Record<Locale, string>> = { en: "English", es: "Español", pt: "Português" };

/** The tag each language declares to search engines and to the Open Graph protocol. */
export const OPEN_GRAPH_LOCALES: Readonly<Record<Locale, string>> = { en: "en_US", es: "es_ES", pt: "pt_BR" };

export function isLocale(value: string): value is Locale {
  return LOCALES.some((locale) => locale === value);
}

/** The path of a page in a language: English has no prefix. */
export function localePath(locale: Locale, path: string): string {
  if (locale === routing.defaultLocale) return path;
  return path === "/" ? `/${locale}` : `/${locale}${path}`;
}

/** The canonical address of a page and its twins in the other languages, for `<link rel="alternate" hreflang>`. */
export function alternates(locale: Locale, path: string) {
  return {
    canonical: localePath(locale, path),
    languages: {
      ...Object.fromEntries(LOCALES.map((l) => [l, localePath(l, path)])),
      "x-default": localePath(routing.defaultLocale, path),
    },
  };
}
