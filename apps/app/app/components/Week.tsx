// SPDX-License-Identifier: MIT OR Apache-2.0
import { usd } from "@/lib/amounts";
import type { Mark } from "@/lib/week";

const at = (h: number | undefined, otherwise: string) => (h === undefined ? otherwise : h === 0 ? "Now" : `In ${h}H`);

const WHEN: Record<Mark["key"], (hours: number | undefined) => [string, string]> = {
  now: () => ["Now", ""],
  ramp: (h) => [at(h, ""), "ramp starts"],
  close: (h) => [at(h, "Weekend"), "closed"],
  reopen: (h) => [at(h, "After"), "24/5 reopens"],
  regular: () => ["After", "regular open"],
};

/** What each mark changes, in a few words. */
const WHAT: Partial<Record<Mark["key"], string>> = {
  ramp: "Climbs for seven hours, leverage capped at 5×.",
  close: "Liquidations slow down: both band edges, a tenth an hour.",
  reopen: "Buffer back; the auction may enrol a short position.",
  regular: "The auction sells at its clearing price.",
};

/**
 * The account's week down the stem: at each mark, the requirement there against the equity. Past the equity the bar's
 * edge turns liquidation red and says by how much it falls short.
 */
export function Week({ marks, equity, label }: { marks: readonly Mark[]; equity: bigint; label: string }) {
  return (
    <ol aria-label={label} className="mt-10 grid gap-1">
      {marks.map((mark) => {
        const [when, what] = WHEN[mark.key](mark.hours);
        const short = mark.requirement > equity;
        const width = equity <= 0n || short ? 100 : Number((mark.requirement * 10_000n) / equity) / 100;
        return (
          <li key={mark.key}>
            <div className="mark" data-now={mark.key === "now" ? "" : undefined}>
              <span className="when">
                {when}
                {what && (
                  <>
                    <br />
                    {what}
                  </>
                )}
              </span>
              <span className="req num" data-short={short ? "" : undefined}>
                <i style={{ "--fill": `${width}%` } as React.CSSProperties} />
                <span>
                  {short ? (
                    <span className="text-[var(--liquidation)]">
                      {usd(mark.requirement)} · short {usd(mark.requirement - equity)}
                    </span>
                  ) : (
                    <>
                      {usd(mark.requirement)} <span className="text-muted">of {usd(equity)}</span>
                    </>
                  )}
                </span>
              </span>
            </div>
            {WHAT[mark.key] && (
              <p className={`ramp ${mark.key === "ramp" ? "" : "!text-muted"}`}>{WHAT[mark.key]}</p>
            )}
          </li>
        );
      })}
    </ol>
  );
}
