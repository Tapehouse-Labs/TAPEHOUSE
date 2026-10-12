// SPDX-License-Identifier: MIT OR Apache-2.0
import type { MetadataRoute } from "next";
import { LOCALES, localePath } from "@/i18n/routing";
import { SITE_URL } from "./site";

const PAGES = [
  { path: "/", changeFrequency: "weekly", priority: 1 },
  { path: "/methodology", changeFrequency: "weekly", priority: 0.7 },
  { path: "/risk", changeFrequency: "monthly", priority: 0.5 },
] as const;

export default function sitemap(): MetadataRoute.Sitemap {
  return PAGES.flatMap(({ path, changeFrequency, priority }) =>
    LOCALES.map((locale) => ({
      url: `${SITE_URL}${localePath(locale, path) === "/" ? "" : localePath(locale, path)}`,
      changeFrequency,
      priority,
      alternates: {
        languages: Object.fromEntries(LOCALES.map((l) => [l, `${SITE_URL}${localePath(l, path) === "/" ? "" : localePath(l, path)}`])),
      },
    })),
  );
}
