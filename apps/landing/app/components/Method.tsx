// SPDX-License-Identifier: MIT OR Apache-2.0
import { useTranslations } from "next-intl";
import { Artwork } from "./Artwork";

const FOLDS = ["band", "failure", "integrate", "changes"] as const;

const CODE: ["c" | "", string][] = [
  ["c", "read"],
  ["", "ITapehouseFeed feed = ITapehouseFeed(NVDA_FEED);"],
  ["", ""],
  ["c", "interface"],
  ["", "(, int256 price, , uint256 updatedAt, ) = feed.latestRoundData();"],
  ["", ""],
  ["c", "band"],
  ["", "(uint256 lower, uint256 upper, Session session, uint8 live)"],
  ["", "    = feed.latestBand();"],
  ["", ""],
  ["", "require(session != Session.Degraded, \"source degraded\");"],
  ["c", "lend"],
  ["", "uint256 collateralPrice = lower;"],
];

/** Chapter five: the method, in an engineer's register. */
export function Method() {
  const t = useTranslations("method");
  return (
    <section className="field field-night" aria-labelledby="ch-method">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <div className="grid gap-x-16 gap-y-10 lg:grid-cols-[1fr_17rem]">
          <div>
            <h2 id="ch-method" className="chapter-title max-w-[15ch]">
              {t("title")}
            </h2>
            <p className="lede mt-8 max-w-[56ch]">
              {t("lede")}
            </p>
          </div>
          <Artwork src="/art/i2-socket.jpg" alt={t("art")} className="order-first lg:order-none" />
        </div>

        <div className="mt-[clamp(3rem,6vw,5rem)] grid gap-x-16 gap-y-12 lg:grid-cols-[1fr_1.1fr]">
          <div>
            {FOLDS.map((f, i) => (
              <details key={f} className="fold" open={i === 0} name="method">
                <summary>
                  {t(`folds.${f}.q`)}
                  <span className="plus" aria-hidden="true" />
                </summary>
                <div>{t(`folds.${f}.a`)}</div>
              </details>
            ))}
          </div>

          <figure className="overflow-hidden rounded-[8px] bg-panel shadow-[inset_0_0_0_1px_var(--rule)]">
            <figcaption className="flex items-center justify-between gap-4 border-b border-rule px-5 py-3 font-mono text-[11px] uppercase tracking-[0.12em] text-muted">
              <span>Consumer.sol</span>
              <span>{t("preview")}</span>
            </figcaption>
            <pre className="overflow-x-auto px-5 py-5 font-mono text-[13px] leading-[1.75] text-body" tabIndex={0}>
              <code>
                {CODE.map(([kind, line], i) => (
                  <span key={i} className={`block ${kind === "c" ? "text-faint" : ""}`}>
                    {kind === "c" ? t(`code.${line as "read" | "interface" | "band" | "lend"}`) : line || " "}
                  </span>
                ))}
              </code>
            </pre>
          </figure>
        </div>
        <p className="mt-6 max-w-[70ch] text-[13px] leading-relaxed text-muted">
          {t("footnote")}
        </p>
      </div>
    </section>
  );
}
