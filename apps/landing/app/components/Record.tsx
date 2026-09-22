import Link from "next/link";
import record from "../data/record.json";
import { Artwork } from "./Artwork";

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
const day = (iso: string) => {
  const [y, m, d] = iso.split("-");
  return `${d} ${MONTHS[Number(m) - 1]} ${y}`;
};
const pct = (bps: number) => `${bps > 0 ? "+" : bps < 0 ? "−" : ""}${(Math.abs(bps) / 100).toFixed(2)}%`;
const move = (p: number) => `${p > 0 ? "+" : p < 0 ? "−" : ""}${Math.abs(p).toFixed(2)}%`;

type Row = { reopen: string; token: string; weekendMovePct: number; errorBps: number; halfBps: number; inside: boolean };

function Result({ inside }: { inside: boolean }) {
  return (
    <span className="inline-flex items-center gap-2">
      {inside ? (
        <span className="h-2 w-2 rounded-full shadow-[inset_0_0_0_1.5px_var(--text-muted)]" aria-hidden="true" />
      ) : (
        <span className="h-2 w-2 bg-strong" aria-hidden="true" />
      )}
      <span className={inside ? "text-muted" : "font-semibold text-strong"}>
        <span className="sm:hidden">{inside ? "Inside" : "Wrong"}</span>
        <span className="hidden sm:inline">{inside ? "Inside the band" : "We were wrong"}</span>
      </span>
    </span>
  );
}

function Table({ rows, caption, offLabel }: { rows: Row[]; caption: string; offLabel: string }) {
  return (
    <div className="overflow-x-auto rounded-[4px] bg-panel shadow-[inset_0_0_0_1px_var(--rule)]">
      <table className="w-full border-collapse text-[13px] sm:text-[14px]">
        <caption className="sr-only">{caption}</caption>
        <thead>
          <tr className="text-right font-mono text-[10.5px] uppercase tracking-[0.1em] text-faint">
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">Reopening</th>
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">Token</th>
            <th scope="col" className="hidden px-4 py-3 font-medium md:table-cell">Weekend move</th>
            <th scope="col" className="px-2.5 py-3 font-medium sm:px-4">{offLabel}</th>
            <th scope="col" className="px-2.5 py-3 font-medium sm:px-4">Band</th>
            <th scope="col" className="px-2.5 py-3 text-left font-medium sm:px-4">Result</th>
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
              <td className="num px-2.5 text-accent sm:px-4">±{(r.halfBps / 100).toFixed(2)}%</td>
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
  const live: Row[] = record.live.entries;
  const stress: Row[] = record.stress.entries;
  const misses = stress.filter((r) => !r.inside);
  const liveInside = live.filter((r) => r.inside).length;

  return (
    <section className="field field-bone" aria-labelledby="ch-record">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <div className="grid gap-x-16 gap-y-10 lg:grid-cols-[1fr_17rem]">
          <div>
            <h2 id="ch-record" className="chapter-title max-w-[17ch]">
              We were wrong <span className="num">{misses.length}</span> of{" "}
              <span className="num">{stress.length}</span> times. Here they are.
            </h2>
            <p className="lede mt-8 max-w-[58ch]">
              We gave our band the worst input we could find, the on-chain exchange price on its own, across{" "}
              <span className="num">{record.stress.weekends}</span> weekends, and graded it against the first
              fresh price when the market came back. It missed{" "}
              <span className="num">{misses.length}</span> times. Anyone can rebuild every row from the chain.
            </p>
          </div>
          <Artwork
            src="/art/i3-cards.jpg"
            alt="A tall stack of sealed cards with one pulled out of line."
            className="order-first lg:order-none"
          />
        </div>

        <h3 className="mt-[clamp(3.5rem,7vw,6rem)] font-display text-[1.35rem] font-semibold tracking-[-0.012em] text-strong">
          The live record
        </h3>
        <p className="mt-3 max-w-[64ch] text-body">
          The band as the product runs it, fed by an around-the-clock source that answers every ten seconds.
          It is one weekend old: <span className="num">{liveInside}</span> of{" "}
          <span className="num">{live.length}</span> reopenings inside.
        </p>
        <div className="mt-6">
          <Table rows={live} caption="Live record of weekend reopenings" offLabel="Off by" />
        </div>
        {!record.live.sealed && (
          <p className="mt-4 flex max-w-[72ch] items-start gap-3 text-[14px] leading-snug text-muted">
            <span className="mt-[5px] h-2 w-2 shrink-0 rotate-45 bg-warning" aria-hidden="true" />
            <span>
              This first weekend was graded after the fact, from our own minute-by-minute capture: the band at
              23:59 UTC against the price published at 00:00. Later entries are sealed and timestamped before
              the market reopens, and marked as sealed.
            </span>
          </p>
        )}

        <h3 className="mt-[clamp(3.5rem,7vw,6rem)] font-display text-[1.35rem] font-semibold tracking-[-0.012em] text-strong">
          The stress test
        </h3>
        <p className="mt-3 max-w-[64ch] text-body">
          The same band, fed only a thirty-minute average of on-chain trades. At a weekend that price drifts
          away from the real one, because few people are trading. These are the{" "}
          <span className="num">{misses.length}</span> times it left the band.
        </p>
        <div className="mt-6">
          <Table rows={misses} caption="Reopenings where the band was wrong" offLabel="Off by" />
        </div>

        <details className="fold mt-6">
          <summary>
            <span>
              Show all <span className="num">{stress.length}</span> reopenings
            </span>
            <span className="plus" aria-hidden="true" />
          </summary>
          <div className="!max-w-none !pr-0">
            <Table rows={stress} caption="Every reopening in the stress test" offLabel="Off by" />
          </div>
        </details>

        <p className="mt-8 max-w-[64ch] text-body">
          For scale: doing nothing, and carrying Friday’s price through the weekend, was off by{" "}
          <span className="num font-semibold text-strong">{(record.baseline.medianGapBps / 100).toFixed(2)}%</span> at
          the median and as much as{" "}
          <span className="num font-semibold text-strong">{(record.baseline.maxGapBps / 100).toFixed(2)}%</span>, over{" "}
          <span className="num">{record.baseline.reopenings}</span> reopenings.{" "}
          <Link href="/methodology" className="text-strong underline">
            How to rebuild these tables
          </Link>
          .
        </p>
      </div>
    </section>
  );
}
