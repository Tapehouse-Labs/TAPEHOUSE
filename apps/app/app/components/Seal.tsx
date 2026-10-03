// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useMinute, useSession } from "@/lib/hooks";
import { seal } from "@/lib/session";

/** The session seal: the market's state as the band reads it, and the hours to its next boundary. */
export function Seal() {
  const session = useSession();
  const now = useMinute();
  if (!session.data) {
    return (
      <span role="status" aria-label="Reading the market session" className="seal">
        <span aria-hidden="true">[ ······ ]</span>
      </span>
    );
  }
  const s = seal(session.data, now);
  return (
    <span role="status" aria-label={s.label} data-state={s.state} className="seal num">
      <span aria-hidden="true">[ {s.state} ]</span>
      {s.hours !== undefined && (
        <>
          <span aria-hidden="true" className="text-muted">·</span>
          <span aria-hidden="true">
            {s.state === "OPEN" ? "CLOSES" : "OPENS"} {s.hours}H
          </span>
        </>
      )}
    </span>
  );
}
