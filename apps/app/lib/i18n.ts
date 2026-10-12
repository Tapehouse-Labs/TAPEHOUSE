// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useLocale, useTranslations } from "next-intl";
import { useCallback, useMemo } from "react";
import type { Locale } from "@/i18n/config";
import { explain } from "./errors";
import { type Formatters, formatters } from "./format";
import { type Msg, render } from "./messages";

/** The function that writes a message in the visitor's language. */
export function useMsg(): (message: Msg) => string {
  const t = useTranslations();
  return useCallback((message) => render(t as never, message), [t]);
}

/** How amounts are written in the visitor's language. */
export function useFormat(): Formatters {
  const locale = useLocale() as Locale;
  return useMemo(() => formatters(locale), [locale]);
}

/** The function that says, in the visitor's language, what went wrong with an error. */
export function useExplain(): (error: unknown) => string {
  const f = useFormat();
  const say = useMsg();
  return useCallback((error) => say(explain(f, error)), [f, say]);
}
