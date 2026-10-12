// SPDX-License-Identifier: MIT OR Apache-2.0
import { useFormatter, useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import record from "../data/record.json";
import { Artwork } from "./Artwork";

const MONTHS = ["jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"] as const;

function useFigures() {
  const t = useTranslations("record");
  const format = useFormatter();
  const percent = (value: number) =>
    t("percent", { value: format.number(value, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) });
  const signed = (value: number) => `${value > 0 ? "+" : value < 0 ? "−" : ""}${percent(Math.abs(value))}`;
  return {
    day: (iso: string) => {
      const [year, month, date] = iso.split("-");
      return t("date", { day: date, month: t(`months.${MONTHS[Number(month) - 1]}`), year });
    },
    percent,
    pct: (bps: number) => signed(bps / 100),
    move: signed,
  };
}

type Row = { reopen: string; token: string; weekendMovePct: number; errorBps: number; halfBps: number; inside: boolean };

function Result({ inside }: { inside: boolean }) {
  const t = useTranslations("record.result");
  return (
    <span className="inline-flex items-center gap-2">
      {inside ? (
        <span className="h-2 w-2 rounded-full shadow-[inset_0_0_0_1.5px_var(--text-muted)]" aria-hidden="true" />
      ) : (
        <span className="h-2 w-2 bg-strong" aria-hidden="true" />
      )}
      <span className={inside ? "text-muted" : "font-semibold text-strong"}>
        <span className="sm:hidden">{inside ? t("insideShort") : t("wrongShort")}</span>
        <span className="hidden sm:inline">{inside ? t("inside") : t("wrong")}</span>
      </span>
    </span>
  );
}

function Table({ rows, caption }: { rows: Row[]; caption: string }) {
  const t = useTranslations("record.table");
  const { day, pct, move, percent } = useFigures();
  return (
    <div className="overflow-x-auto rounded-[4px] bg-panel shadow-[inset_0_0_0_1px_var(--rule)]">
      <table className="w-full border-collapse text-[13px] sm:text-[14px]">
        <caption className="sr-only">{caption}</caption>
        <thead>
          <tr className="text-right font-mono text-[10.5px] uppercase tracking-[0.1em] text-faint">
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">{t("reopening")}</th>
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">{t("token")}</th>
            <th scope="col" className="hidden px-4 py-3 font-medium md:table-cell">{t("weekendMove")}</th>
            <th scope="col" className="px-2.5 py-3 font-medium sm:px-4">{t("offBy")}</th>
            <th scope="col" className="px-2.5 py-3 font-medium sm:px-4">{t("band")}</th>
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">{t("result")}</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr
              key={`${r.reopen}-${r.token}`}
              className={`h-10 border-t border-rule text-right ${r.inside ? "" : "bg-[color-mix(in_srgb,var(--ink)_5%,transparent)]"}`}
            >
              <th scope="row" className="num whitespace-nowrap px-2.5 text-left font-medium text-strong sm:px-4">
                {day(r.reopen)}
              </th>
              <td className="px-2.5 text-left font-semibold text-strong sm:px-4">{r.token}</td>
              <td className="num hidden px-4 text-muted md:table-cell">{move(r.weekendMovePct)}</td>
              <td className="num px-2.5 font-semibold text-strong sm:px-4">{pct(r.errorBps)}</td>
              <td className="num px-2.5 text-accent sm:px-4">±{percent(r.halfBps / 100)}</td>
              <td className="px-2.5 text-left sm:px-4">
                <Result inside={r.inside} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/** Chapter six: the record. Every number here is computed from the data file. */
export function Record() {
  const t = useTranslations("record");
  const { percent } = useFigures();
  const live: Row[] = record.live.entries;
  const stress: Row[] = record.stress.entries;
  const misses = stress.filter((r) => !r.inside);
  const liveInside = live.filter((r) => r.inside).length;
  const rich = {
    num: (chunks: React.ReactNode) => <span className="num">{chunks}</span>,
    strong: (chunks: React.ReactNode) => <span className="num font-semibold text-strong">{chunks}</span>,
    link: (chunks: React.ReactNode) => (
      <Link href="/methodology" className="text-strong underline">
        {chunks}
      </Link>
    ),
  };

  return (
    <section className="field field-bone" aria-labelledby="ch-record">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <div className="grid gap-x-16 gap-y-10 lg:grid-cols-[1fr_17rem]">
          <div>
            <h2 id="ch-record" className="chapter-title max-w-[17ch]">
              {t.rich("title", { ...rich, misses: misses.length, total: stress.length })}
            </h2>
            <p className="lede mt-8 max-w-[58ch]">
              {t.rich("lede", { ...rich, weekends: record.stress.weekends, misses: misses.length })}
            </p>
          </div>
          <Artwork src="/art/i3-cards.jpg" alt={t("artAlt")} className="order-first lg:order-none" />
        </div>

        <h3 className="mt-[clamp(3.5rem,7vw,6rem)] font-display text-[1.35rem] font-semibold tracking-[-0.012em] text-strong">
          {t("live.title")}
        </h3>
        <p className="mt-3 max-w-[64ch] text-body">
          {t.rich("live.body", { ...rich, inside: liveInside, total: live.length })}
        </p>
        <div className="mt-6">
          <Table rows={live} caption={t("live.caption")} />
        </div>
        {!record.live.sealed && (
          <p className="mt-4 flex max-w-[72ch] items-start gap-3 text-[14px] leading-snug text-muted">
            <span className="mt-[5px] h-2 w-2 shrink-0 rotate-45 bg-warning" aria-hidden="true" />
            <span>{t("live.note")}</span>
          </p>
        )}

        <h3 className="mt-[clamp(3.5rem,7vw,6rem)] font-display text-[1.35rem] font-semibold tracking-[-0.012em] text-strong">
          {t("stress.title")}
        </h3>
        <p className="mt-3 max-w-[64ch] text-body">{t.rich("stress.body", { ...rich, misses: misses.length })}</p>
        <div className="mt-6">
          <Table rows={misses} caption={t("stress.captionMisses")} />
        </div>

        <details className="fold mt-6">
          <summary>
            <span>{t.rich("stress.showAll", { ...rich, count: stress.length })}</span>
            <span className="plus" aria-hidden="true" />
          </summary>
          <div className="!max-w-none !pr-0">
            <Table rows={stress} caption={t("stress.captionAll")} />
          </div>
        </details>

        <p className="mt-8 max-w-[64ch] text-body">
          {t.rich("baseline", {
            ...rich,
            median: percent(record.baseline.medianGapBps / 100),
            max: percent(record.baseline.maxGapBps / 100),
            reopenings: record.baseline.reopenings,
          })}
        </p>
      </div>
    </section>
  );
}
