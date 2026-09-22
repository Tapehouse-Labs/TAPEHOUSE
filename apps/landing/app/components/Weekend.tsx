"use client";

import Link from "next/link";
import { useEffect, useRef } from "react";
import { clamp01, onScrollFrame, prefersStill } from "./scroll";

/**
 * Chapter three: one measured weekend, drawn to scale.
 * Axis runs Fri 00:00 UTC to Tue 00:00 UTC (96 h). The reference feed's last
 * update before the weekend was Fri 16:18 UTC and its next Mon 00:00 UTC.
 */
const W = 1200;
const PAD = 40;
const SPAN = 96;
const x = (hours: number) => PAD + (hours / SPAN) * (W - PAD * 2);
const LAST = 16.3;
const NEXT = 72;
const DAYS = ["Fri", "Sat", "Sun", "Mon"];
const TICKS = Array.from({ length: 288 }, (_, i) => 24 + (i * 48) / 288);

export function Weekend() {
  const root = useRef<HTMLDivElement>(null);
  const hours = useRef<HTMLSpanElement>(null);
  const sweep = useRef<SVGRectElement>(null);

  useEffect(() => {
    const el = root.current;
    if (!el) return;
    const paint = (p: number) => {
      sweep.current?.setAttribute("width", (PAD + p * (W - PAD * 2)).toFixed(1));
      const t = p * SPAN;
      const silent = Math.min(NEXT - LAST, Math.max(0, t - LAST));
      if (hours.current) hours.current.textContent = silent.toFixed(1);
    };
    if (prefersStill()) return paint(1);
    return onScrollFrame(() => {
      const r = el.getBoundingClientRect();
      const vh = window.innerHeight;
      paint(clamp01((vh * 0.9 - r.top) / (vh * 0.75 + r.height * 0.35)));
    });
  }, []);

  return (
    <section id="measured" className="field field-navy" aria-labelledby="ch-weekend">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <h2 id="ch-weekend" className="chapter-title max-w-[17ch]">
          The market closes on Friday. Your loan doesn’t.
        </h2>
        <p className="lede mt-8 max-w-[60ch]">
          We measured it before we built it. Every weekend the reference price feed on this chain goes
          silent for <span className="num font-semibold text-strong">55.7 to 80.8 hours</span>, while an
          around-the-clock source kept answering{" "}
          <span className="num font-semibold text-strong">288 times out of 288</span>. People don’t stop
          either: <span className="num font-semibold text-strong">29.9%</span> of SPY transfers happen
          on a Saturday or a Sunday.
        </p>

        <div ref={root} className="mt-[clamp(3.5rem,7vw,6rem)]">
          <p className="flex items-baseline gap-3 text-strong" aria-live="off">
            <span className="num font-sans text-[clamp(3rem,8vw,6.5rem)] font-semibold leading-none tracking-[-0.04em]">
              <span ref={hours}>55.7</span>
            </span>
            <span className="text-[clamp(1rem,1.4vw,1.25rem)] text-muted">
              hours since the reference feed last spoke
            </span>
          </p>

          <div className="-mx-5 mt-8 overflow-x-auto px-5 sm:mx-0 sm:px-0">
            <svg
              viewBox={`0 0 ${W} 300`}
              className="block h-auto w-full min-w-[760px]"
              role="img"
              aria-label="Timeline of one weekend. The reference equity feed updates until Friday 16:18 UTC and not again until Monday 00:00 UTC, 55.7 hours later. An around-the-clock source reports every ten minutes throughout, 288 points out of 288."
            >
              <defs>
                <clipPath id="sweep">
                  <rect ref={sweep} x="0" y="0" height="300" width={W} />
                </clipPath>
              </defs>

              {DAYS.map((d, i) => (
                <g key={d}>
                  <line x1={x(i * 24)} x2={x(i * 24)} y1="28" y2="262" stroke="currentColor" strokeOpacity="0.14" />
                  <text x={x(i * 24) + 10} y="22" fill="currentColor" fillOpacity="0.62" fontSize="15" className="font-mono" letterSpacing="1.4">
                    {d.toUpperCase()}
                  </text>
                </g>
              ))}
              <line x1={x(96)} x2={x(96)} y1="28" y2="262" stroke="currentColor" strokeOpacity="0.14" />

              <text x={PAD} y="86" fill="currentColor" fillOpacity="0.7" fontSize="16">Reference equity feed</text>
              <text x={PAD} y="186" fill="currentColor" fillOpacity="0.7" fontSize="16">Around-the-clock source</text>

              {/* What the weekend looks like before anyone measures it */}
              <line x1={x(LAST)} x2={x(NEXT)} y1="112" y2="112" stroke="currentColor" strokeOpacity="0.22" strokeDasharray="2 7" />

              <g clipPath="url(#sweep)">
                <line x1={x(0)} x2={x(LAST)} y1="112" y2="112" stroke="currentColor" strokeWidth="3" />
                <line x1={x(NEXT)} x2={x(96)} y1="112" y2="112" stroke="currentColor" strokeWidth="3" />
                <line x1={x(LAST)} x2={x(LAST)} y1="100" y2="124" stroke="currentColor" strokeWidth="3" />
                <line x1={x(NEXT)} x2={x(NEXT)} y1="100" y2="124" stroke="currentColor" strokeWidth="3" />
                <text x={x(LAST)} y="146" textAnchor="middle" fill="currentColor" fillOpacity="0.62" fontSize="13" className="font-mono">FRI 16:18 UTC</text>
                <text x={x(NEXT)} y="146" textAnchor="middle" fill="currentColor" fillOpacity="0.62" fontSize="13" className="font-mono">MON 00:00 UTC</text>

                <line x1={x(0)} x2={x(96)} y1="212" y2="212" stroke="var(--accent)" strokeOpacity="0.45" strokeWidth="1" />
                {TICKS.map((t) => (
                  <line key={t} x1={x(t)} x2={x(t)} y1="203" y2="221" stroke="var(--accent)" strokeWidth="1.3" />
                ))}
                <text x={x(48)} y="248" textAnchor="middle" fill="var(--accent)" fontSize="13" className="font-mono">
                  288 OF 288 TEN-MINUTE POINTS
                </text>
              </g>
            </svg>
          </div>
          <p className="mt-5 max-w-[70ch] text-[13.5px] leading-relaxed text-muted">
            One measured weekend, 28–31 August 2026, drawn to scale. Over the Labor Day weekend the same
            silence lasted 80.8 hours.{" "}
            <Link href="/methodology" className="text-strong underline">How we measured it</Link>.
          </p>
        </div>
      </div>
    </section>
  );
}
