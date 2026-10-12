// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import { REFRESH_MS, useMinute, useSettled } from "@/lib/hooks";
import { useExplain, useFormat, useMsg } from "@/lib/i18n";
import { readMarket, simulate } from "@/lib/simulation";
import { capacity, regimeWord, week } from "@/lib/week";
import { Week } from "./Week";

const PRICE_UNIT = 10n ** 8n;

/**
 * The public margin simulator: any portfolio of the engine's assets, priced at each band's low edge and margined by
 * the engine itself, with what it could borrow now and its week down the stem. No wallet needed.
 */
export function Simulator() {
  const t = useTranslations("simulator");
  const f = useFormat();
  const say = useMsg();
  const explain = useExplain();
  const deployments = useDeployments();
  const client = usePublicClient();
  const now = useMinute();
  const [typed, setTyped] = useState<Record<string, string>>({});
  const [typedDebt, setTypedDebt] = useState("");
  const market = useQuery({
    queryKey: ["market", deployments.chainId],
    queryFn: () => readMarket(client!, deployments),
    enabled: client !== undefined,
    refetchInterval: REFRESH_MS,
  });
  const assets = market.data?.assets ?? [];
  const quantities = assets.map((a) => f.parseAmount(typed[a.asset] ?? "", 18) ?? 0n);
  const prices = assets.map((a) => a.low);
  const debt = f.parseAmount(typedDebt, 6) ?? 0n;
  const settled = useSettled(quantities.join(","));
  const result = useQuery({
    queryKey: ["simulation", deployments.chainId, market.data?.blockNumber.toString(), settled],
    queryFn: () => simulate(client!, deployments, settled.split(",").map(BigInt), prices, market.data!.blockNumber),
    enabled: client !== undefined && market.data !== undefined && quantities.some((q) => q > 0n),
    placeholderData: (previous) => previous,
  });

  const gross = quantities.reduce((sum, q, i) => sum + (q * prices[i]!) / PRICE_UNIT, 0n);
  const equity = gross - debt * 10n ** 12n;
  const data = market.data;
  const sim = quantities.some((q) => q > 0n) ? result.data : undefined;
  const ahead = sim ? sim.regime === 0 || ((sim.regime === 2 || sim.regime === 3) && data!.session.boundaryMs !== 0n) : false;
  const room =
    sim && data
      ? capacity({
          equity,
          requirement: sim.current,
          gross,
          weekendAhead: ahead,
          weekendLeverage: data.limits.weekendLeverage,
          debt: data.limits.debt,
          debtCap: data.limits.debtCap,
          weekendDebtCap: data.limits.weekendDebtCap,
          lendable: data.lendable,
        })
      : undefined;

  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <div className="relative pt-14">
          <span className="bar" aria-hidden="true" />
          <h1 className="label">{t("title")}</h1>
          {room !== undefined ? (
            <p className="figure num mt-3 !text-strong" aria-live="polite">
              {f.usdg(room)} <span className="text-[0.32em] font-medium tracking-normal text-muted">USDG</span>
            </p>
          ) : (
            <p className="figure mt-3" data-empty>
              {t("any")}
            </p>
          )}
          <p className="mt-3 max-w-[52ch] text-[14px] text-muted num">
            {sim ? (
              t("required", { required: f.usd(sim.current), equity: f.usd(equity), regime: say(regimeWord(sim.regime)) })
            ) : (
              t("hint")
            )}
          </p>
          {data && (
            <p className="mt-2 max-w-[52ch] text-[14px] text-muted num">
              {t("lendable", { amount: f.usdg(data.lendable) })}
            </p>
          )}
        </div>

        <section aria-label={t("portfolio")} className="mt-10">
          {market.isError && (
            <p role="alert" className="text-[14px] text-[var(--error)]">
              {explain(market.error)}
            </p>
          )}
          <div className="grid gap-1">
            {assets.map((a) => (
              <label key={a.asset} className="row flex min-h-[52px] items-center justify-between gap-4 border-t border-rule">
                <span className="font-semibold text-strong">{a.asset}</span>
                <span className="flex items-center gap-3 text-[13px] text-muted num">
                  {a.low > 0n ? t("lowEdge", { price: f.price(a.low) }) : t("halted")}
                  <input
                    className="amount !text-[19px]"
                    inputMode="decimal"
                    placeholder="0"
                    aria-label={t("tokens", { asset: a.asset })}
                    value={typed[a.asset] ?? ""}
                    onChange={(e) => setTyped({ ...typed, [a.asset]: f.normalize(e.target.value) })}
                  />
                </span>
              </label>
            ))}
            {data && (
              <label className="row flex min-h-[52px] items-center justify-between gap-4 border-t border-rule">
                <span className="font-semibold text-strong">{t("owed")}</span>
                <span className="flex items-center gap-3 text-[13px] text-muted num">
                  USDG
                  <input
                    className="amount !text-[19px]"
                    inputMode="decimal"
                    placeholder="0"
                    aria-label={t("owedAria")}
                    value={typedDebt}
                    onChange={(e) => setTypedDebt(f.normalize(e.target.value))}
                  />
                </span>
              </label>
            )}
          </div>
        </section>

        {sim && data && (
          <Week
            label={t("week")}
            equity={equity}
            marks={week({
              regime: sim.regime,
              current: sim.current,
              open: sim.open,
              closed: sim.closed,
              gross,
              weekendLeverage: data.limits.weekendLeverage,
              boundaryMs: data.session.boundaryMs,
              nowMs: now,
            })}
          />
        )}
      </div>
    </main>
  );
}
