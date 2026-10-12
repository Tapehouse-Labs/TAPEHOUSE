// SPDX-License-Identifier: MIT OR Apache-2.0
import { cookies, headers } from "next/headers";
import { getRequestConfig } from "next-intl/server";
import { DEFAULT_LOCALE, isLocale, LOCALE_COOKIE, type Locale, negotiate } from "./config";
import { catalogs } from "@/messages";

/** The visitor's language: the one they chose, else the one their browser prefers, else English. */
async function locale(): Promise<Locale> {
  const chosen = (await cookies()).get(LOCALE_COOKIE)?.value;
  if (isLocale(chosen)) return chosen;
  return negotiate((await headers()).get("accept-language")) ?? DEFAULT_LOCALE;
}

export default getRequestConfig(async () => {
  const resolved = await locale();
  return { locale: resolved, messages: catalogs[resolved] };
});
