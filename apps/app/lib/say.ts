// SPDX-License-Identifier: MIT OR Apache-2.0
import { createTranslator } from "next-intl";
import type { Locale } from "@/i18n/config";
import { catalogs } from "@/messages";
import { type Msg, render } from "./messages";

/** A message written in a language, outside a component: what a test reads, and what a log carries. */
export function say(message: Msg, locale: Locale = "en"): string {
  const t = createTranslator({ locale, messages: catalogs[locale] });
  return render(t as never, message);
}
