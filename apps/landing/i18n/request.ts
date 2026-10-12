// SPDX-License-Identifier: MIT OR Apache-2.0
import { hasLocale } from "next-intl";
import { getRequestConfig } from "next-intl/server";
import { catalogs } from "@/messages";
import { routing } from "./routing";

export default getRequestConfig(async ({ requestLocale }) => {
  const requested = await requestLocale;
  const locale = hasLocale(routing.locales, requested) ? requested : routing.defaultLocale;
  return { locale, messages: catalogs[locale] };
});
