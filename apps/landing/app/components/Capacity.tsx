// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useLocale, useTranslations } from "next-intl";
import { useMemo, useState } from "react";
import { Link } from "@/i18n/navigation";
import { Artwork } from "./Artwork";

/**
 * Chapter four: the same portfolio, priced two ways.
 * The model is illustrative and says so on the page. It shows the mechanism:
 * a flat per-asset rule cannot see diversification, a portfolio model can.
 */
const TOKENS = ["NVDA", "TSLA", "AAPL", "SPY"] as const;
type Token = (typeof TOKENS)[number];
type Weights = Record<Token, number>;

const PORTFOLIO_USD = 67980.72;
const FLAT_RULE = 0.625;
const VOL: Record<Token, number> = { NVDA: 0.48, TSLA: 0.58, AAPL: 0.27, SPY: 0.16 };
const CORR: Record<string, number> = {
  "NVDA-TSLA": 0.45, "NVDA-AAPL": 0.55, "NVDA-SPY": 0.7,
  "TSLA-AAPL": 0.4, "TSLA-SPY": 0.55, "AAPL-SPY": 0.75,
};
const rho = (a: Token, b: Token) => (a === b ? 1 : CORR[`${a}-${b}`] ?? CORR[`${b}-${a}`]);

const PRESETS: { id: "hers" | "single" | "index"; weights: Weights }[] = [
  { id: "hers", weights: { NVDA: 32, TSLA: 26, AAPL: 22, SPY: 20 } },
  { id: "single", weights: { NVDA: 0, TSLA: 100, AAPL: 0, SPY: 0 } },
  { id: "index", weights: { NVDA: 8, TSLA: 4, AAPL: 8, SPY: 80 } },
];

/** Worst one-percent loss over a closed weekend, stressed, plus a fixed buffer. */
function portfolioCapacity(w: Weights) {
  let variance = 0;
  for (const a of TOKENS) for (const b of TOKENS) variance += (w[a] / 100) * (w[b] / 100) * VOL[a] * VOL[b] * rho(a, b);
  const loss = 2.6 * Math.sqrt(variance) * Math.sqrt(3 / 252) * 2;
  return Math.min(0.9, Math.max(0.35, 1 - loss - 0.1));
}

const SCALE = { from: 0.35, to: 0.9 };
const at = (v: number) => `${((v - SCALE.from) / (SCALE.to - SCALE.from)) * 100}%`;

export function Capacity() {
  const t = useTranslations("capacity");
  const locale = useLocale();
  const usd = useMemo(() => new Intl.NumberFormat(locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }), [locale]);
  const pct = useMemo(() => new Intl.NumberFormat(locale, { minimumFractionDigits: 1, maximumFractionDigits: 1 }), [locale]);
  const [weights, setWeights] = useState<Weights>(PRESETS[0].weights);
  const capacity = useMemo(() => portfolioCapacity(weights), [weights]);
  const delta = capacity - FLAT_RULE;
  const active = PRESETS.findIndex((p) => TOKENS.every((t) => p.weights[t] === weights[t]));

  /** Move one weight and rescale the others so the book always sums to 100. */
  const move = (token: Token, value: number) => {
    setWeights((prev) => {
      const rest = TOKENS.filter((t) => t !== token);
      const restSum = rest.reduce((sum, t) => sum + prev[t], 0);
      const room = 100 - value;
      const next = { ...prev, [token]: value } as Weights;
      let used = 0;
      rest.forEach((t, i) => {
        const share = restSum > 0 ? prev[t] / restSum : 1 / rest.length;
        next[t] = i === rest.length - 1 ? room - used : Math.round(share * room);
        used += next[t];
      });
      return next;
    });
  };

  const verdict = delta > 0.015 ? t("verdict.more") : delta < -0.015 ? t("verdict.less") : t("verdict.same");

  return (
    <section className="field field-white" aria-labelledby="ch-capacity">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <div className="grid gap-x-16 gap-y-10 lg:grid-cols-[1fr_17rem]">
          <div>
            <h2 id="ch-capacity" className="chapter-title max-w-[16ch]">
              {t("title")}
            </h2>
            <p className="lede mt-8 max-w-[56ch]">
              {t("lede")}
            </p>
          </div>
          <Artwork src="/art/i1-beam.jpg" alt={t("art")} className="order-first lg:order-none" />
        </div>

        <div className="mt-[clamp(3rem,6vw,5rem)] grid gap-x-16 gap-y-12 border-t border-rule pt-10 lg:grid-cols-[1fr_1.15fr]">
          <div>
            <div className="flex flex-wrap gap-2" role="group" aria-label={t("presetsLabel")}>
              {PRESETS.map((p, i) => (
                <button key={p.id} type="button" className="chip" aria-pressed={i === active} onClick={() => setWeights(p.weights)}>
                  {t(`presets.${p.id}`)}
                </button>
              ))}
            </div>
            <div className="mt-8 space-y-5">
              {TOKENS.map((token) => (
                <div key={token} className="grid grid-cols-[3.5rem_1fr_3.25rem] items-center gap-4">
                  <label htmlFor={`w-${token}`} className="font-semibold text-strong">{token}</label>
                  <input
                    id={`w-${token}`}
                    className="weight"
                    type="range"
                    min={0}
                    max={100}
                    step={1}
                    value={weights[token]}
                    onChange={(e) => move(token, Number(e.target.value))}
                    aria-valuetext={t("valueText", { value: weights[token] })}
                  />
                  <output htmlFor={`w-${token}`} className="num text-right font-semibold text-strong">{weights[token]}%</output>
                </div>
              ))}
            </div>
          </div>

          <div aria-live="polite">
            <div className="grid grid-cols-2 gap-x-8">
              <div>
                <p className="text-[14px] text-muted">{t("oneByOne")}</p>
                <p className="num mt-1.5 text-[clamp(1.6rem,3vw,2.4rem)] font-semibold leading-none tracking-[-0.02em] text-muted">
                  {usd.format(PORTFOLIO_USD * FLAT_RULE)}
                </p>
                <p className="num mt-2 text-[14px] text-muted">{pct.format(FLAT_RULE * 100)}% · USDG</p>
              </div>
              <div>
                <p className="text-[14px] text-body">{t("asPortfolio")}</p>
                <p className="num mt-1.5 text-[clamp(1.6rem,3vw,2.4rem)] font-semibold leading-none tracking-[-0.02em] text-strong">
                  {usd.format(PORTFOLIO_USD * capacity)}
                </p>
                <p className="num mt-2 text-[14px] text-accent">{pct.format(capacity * 100)}% · USDG</p>
              </div>
            </div>

            <div className="relative mt-10 h-14" aria-hidden="true">
              <div className="absolute inset-x-0 top-6 h-px bg-rule" />
              <div className="absolute top-3 h-7 w-px bg-[var(--text-faint)]" style={{ left: at(FLAT_RULE) }} />
              <div
                className="absolute top-1 h-11 w-[3px] rounded-[1px] bg-accent transition-[left] duration-700 ease-[cubic-bezier(0.16,1,0.3,1)]"
                style={{ left: at(capacity) }}
              />
              <span className="num absolute top-[3.1rem] -translate-x-1/2 text-[12px] text-faint" style={{ left: at(FLAT_RULE) }}>
                {t("flatRule")}
              </span>
            </div>

            <p className="mt-8 max-w-[44ch] text-[1.12rem] leading-snug text-strong">{verdict}</p>
            <p className="mt-6 max-w-[58ch] text-[13px] leading-relaxed text-muted">
              {t.rich("note", {
                amount: usd.format(PORTFOLIO_USD),
                link: (c) => <Link href="/methodology" className="text-strong underline">{c}</Link>,
              })}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
