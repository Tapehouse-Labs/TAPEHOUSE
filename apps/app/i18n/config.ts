// SPDX-License-Identifier: MIT OR Apache-2.0

/** The languages the app is written in, English first: the one it falls back to. */
export const LOCALES = ["en", "es", "pt"] as const;
export type Locale = (typeof LOCALES)[number];
export const DEFAULT_LOCALE: Locale = "en";

/** The cookie that holds the language the visitor chose. */
export const LOCALE_COOKIE = "locale";

/** Each language's name in itself, as its speaker looks for it. */
export const LANGUAGE_NAMES: Readonly<Record<Locale, string>> = { en: "English", es: "Español", pt: "Português" };

export function isLocale(value: string | undefined | null): value is Locale {
  return LOCALES.some((locale) => locale === value);
}

/**
 * The language an `Accept-Language` header prefers among LOCALES, by quality and then by order, matching a language
 * by its primary subtag (`pt-BR` is Portuguese); English where none matches or the header is absent.
 */
export function negotiate(header: string | null | undefined): Locale {
  const ranked = (header ?? "")
    .split(",")
    .map((part, order) => {
      const [tag = "", ...parameters] = part.trim().split(";");
      const q = parameters.map((p) => p.trim()).find((p) => p.startsWith("q="));
      const quality = q === undefined ? 1 : Number(q.slice(2));
      return { primary: tag.trim().toLowerCase().split("-")[0], quality: Number.isNaN(quality) ? 0 : quality, order };
    })
    .filter((entry) => entry.quality > 0)
    .sort((a, b) => b.quality - a.quality || a.order - b.order);
  for (const { primary } of ranked) if (isLocale(primary)) return primary;
  return DEFAULT_LOCALE;
}
