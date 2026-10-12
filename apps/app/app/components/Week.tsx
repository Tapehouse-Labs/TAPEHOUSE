// SPDX-License-Identifier: MIT OR Apache-2.0
import { useTranslations } from "next-intl";
import { useFormat } from "@/lib/i18n";
import type { Mark } from "@/lib/week";

type Translate = ReturnType<typeof useTranslations<"account.week">>;

const WHAT: Partial<Record<Mark["key"], "whatRamp" | "whatClose" | "whatReopen" | "whatRegular">> = {
  ramp: "whatRamp",
  close: "whatClose",
  reopen: "whatReopen",
  regular: "whatRegular",
};

function when(t: Translate, mark: Mark): [string, string] {
  const at = (otherwise: string) =>
    mark.hours === undefined ? otherwise : mark.hours === 0 ? t("now") : t("inHours", { hours: mark.hours });
  switch (mark.key) {
    case "now":
      return [t("now"), ""];
    case "ramp":
      return [at(""), t("rampStarts")];
    case "close":
      return [at(t("weekend")), t("closedWord")];
    case "reopen":
      return [at(t("after")), t("reopens")];
    case "regular":
      return [t("after"), t("regularOpen")];
  }
}

/**
 * The account's week down the stem: at each mark, the requirement there against the equity. Past the equity the bar's
 * edge turns liquidation red and says by how much it falls short.
 */
export function Week({ marks, equity, label }: { marks: readonly Mark[]; equity: bigint; label: string }) {
  const t = useTranslations("account.week");
  const f = useFormat();
  return (
    <ol aria-label={label} className="mt-10 grid gap-1">
      {marks.map((mark) => {
        const [at, what] = when(t, mark);
        const note = WHAT[mark.key];
        const short = mark.requirement > equity;
        const width = equity <= 0n || short ? 100 : Number((mark.requirement * 10_000n) / equity) / 100;
        return (
          <li key={mark.key}>
            <div className="mark" data-now={mark.key === "now" ? "" : undefined}>
              <span className="when">
                {at}
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
                      {t("short", { requirement: f.usd(mark.requirement), shortfall: f.usd(mark.requirement - equity) })}
                    </span>
                  ) : (
                    t.rich("of", {
                      requirement: f.usd(mark.requirement),
                      equity: f.usd(equity),
                      muted: (c) => <span className="text-muted">{c}</span>,
                    })
                  )}
                </span>
              </span>
            </div>
            {note && <p className={`ramp ${mark.key === "ramp" ? "" : "!text-muted"}`}>{t(note)}</p>}
          </li>
        );
      })}
    </ol>
  );
}
