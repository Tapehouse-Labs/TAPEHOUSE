// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useTranslations } from "next-intl";
import { useMinute, useSession } from "@/lib/hooks";
import { useMsg } from "@/lib/i18n";
import { seal } from "@/lib/session";

/** The session seal: the market's state as the band reads it, and the hours to its next boundary. */
export function Seal() {
  const t = useTranslations("account.seal");
  const say = useMsg();
  const session = useSession();
  const now = useMinute();
  if (!session.data) {
    return (
      <span role="status" aria-label={t("reading")} className="seal">
        <span aria-hidden="true">[ ······ ]</span>
      </span>
    );
  }
  const s = seal(session.data, now);
  return (
    <span role="status" aria-label={say(s.label)} data-state={s.state} className="seal num">
      <span aria-hidden="true">[ {t(`state.${s.state}`)} ]</span>
      {s.eta !== undefined && (
        <>
          <span aria-hidden="true" className="text-muted">·</span>
          <span aria-hidden="true">
            {s.state === "OPEN" ? t("closes") : t("opens")} {s.eta}
          </span>
        </>
      )}
    </span>
  );
}
