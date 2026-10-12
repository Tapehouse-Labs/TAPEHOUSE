// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Locale } from "@/i18n/config";
import type { Messages } from "@/messages";

declare module "next-intl" {
  interface AppConfig {
    Locale: Locale;
    Messages: Messages;
  }
}
