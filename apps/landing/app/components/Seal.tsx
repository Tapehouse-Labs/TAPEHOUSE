"use client";

import { useEffect, useState } from "react";

type Session = { state: "OPEN" | "CLOSED"; detail: string };

const OPEN_MIN = 9 * 60 + 30;
const CLOSE_MIN = 16 * 60;

/** Wall-clock parts in New York, independent of the visitor's time zone. */
function newYorkParts(date: Date) {
  const parts = new Intl.DateTimeFormat("en-US", {
    timeZone: "America/New_York",
    weekday: "short",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).formatToParts(date);
  const get = (type: string) => parts.find((p) => p.type === type)?.value ?? "";
  const day = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].indexOf(get("weekday"));
  const minutes = (Number(get("hour")) % 24) * 60 + Number(get("minute"));
  return { day, minutes };
}

/**
 * Regular-hours session of the underlying US equity market. Exchange holidays
 * are not modelled here; the product reads the real session state on-chain.
 */
function readSession(now: Date): Session {
  const { day, minutes } = newYorkParts(now);
  const weekday = day >= 1 && day <= 5;
  if (weekday && minutes >= OPEN_MIN && minutes < CLOSE_MIN) {
    return { state: "OPEN", detail: "LIVE" };
  }
  // Walk back to the most recent 16:00 close.
  let daysBack = 0;
  let d = day;
  if (!(weekday && minutes >= CLOSE_MIN)) {
    do {
      daysBack += 1;
      d = (d + 6) % 7;
    } while (d === 0 || d === 6);
  }
  const hours = Math.max(1, Math.round((daysBack * 1440 + minutes - CLOSE_MIN) / 60));
  return { state: "CLOSED", detail: `${hours}h` };
}

/** The session seal: the market's state, declared next to every money figure. */
export function Seal({ tone = "night" }: { tone?: "night" | "paper" }) {
  const [session, setSession] = useState<Session | null>(null);

  useEffect(() => {
    const tick = () => setSession(readSession(new Date()));
    tick();
    const id = window.setInterval(tick, 60_000);
    return () => window.clearInterval(id);
  }, []);

  const night = tone === "night";
  return (
    <span
      role="status"
      aria-label={
        session
          ? session.state === "OPEN"
            ? "US equity market is open"
            : `US equity market is closed, ${session.detail} since the close`
          : "Reading the market session"
      }
      className={[
        "inline-flex items-center gap-2 whitespace-nowrap rounded-[2px] px-2.5 py-1",
        "font-mono text-[11px] tracking-[0.09em] num",
        night
          ? "bg-[rgb(245_242_236/0.08)] text-[rgb(245_242_236/0.92)] shadow-[inset_0_0_0_1px_rgb(245_242_236/0.2)]"
          : "bg-[#fff] text-[#1a1713] shadow-[inset_0_0_0_1px_#dcd6ca]",
      ].join(" ")}
    >
      <span aria-hidden="true">[ {session?.state ?? "······"} ]</span>
      <span aria-hidden="true" className="opacity-50">·</span>
      <span aria-hidden="true" className="opacity-75">{session?.detail ?? "··"}</span>
    </span>
  );
}
