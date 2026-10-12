// SPDX-License-Identifier: MIT OR Apache-2.0
import { useFormatter, useTranslations } from "next-intl";
import { Mark } from "./Mark";
import { Seal } from "./Seal";

/**
 * The lit screen the camera pushes into. Its surface never inverts with the
 * theme: it is a screen glowing in a dark room. All figures are illustrative
 * and labelled as such on the screen itself.
 */
const HOLDINGS = [
  { token: "NVDA", price: 182.5, band: 1.1, value: 21900, move: 1.24 },
  { token: "TSLA", price: 396.12, band: 1.6, value: 17825.4, move: -0.87 },
  { token: "AAPL", price: 241.08, band: 0.9, value: 14464.8, move: 0.31 },
  { token: "SPY", price: 766.14, band: 0.6, value: 13790.52, move: -0.12 },
];

const FALLS_TO = 141.08;
const SOLD_SHARE = 12;
const BORROWED = 10000;
const TOTAL = 67980.72;

export function PortfolioScreen() {
  const t = useTranslations("screen");
  const format = useFormatter();
  const cents = (value: number) => format.number(value, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  const signed = (value: number) => `${value > 0 ? "+" : "−"}${cents(Math.abs(value))}%`;
  return (
    <div
      className="overflow-hidden rounded-[10px] bg-[#fbfaf7] text-[#1a1713] shadow-[0_40px_120px_-20px_rgb(0_0_0/0.65),0_0_0_1px_rgb(245_242_236/0.14)]"
      aria-label={t("aria")}
    >
      <div className="flex items-center justify-between gap-4 border-b border-[#dcd6ca] px-5 py-3.5 sm:px-7">
        <div className="flex items-center gap-2.5 text-[#16326b]">
          <Mark className="h-[18px] w-auto" />
          <span className="font-display text-[13px] font-semibold tracking-[0.16em] text-[#1a1713]">
            TAPEHOUSE
          </span>
        </div>
        <Seal tone="paper" />
      </div>

      <div className="grid gap-x-10 gap-y-7 px-5 py-7 sm:px-7 lg:grid-cols-[1.05fr_1fr] lg:py-9">
        <div>
          <p className="text-[13px] text-[#1a1713]/60">{t("liquidationLabel")}</p>
          <p className="num mt-2 font-sans text-[clamp(2.75rem,7vw,4.5rem)] font-semibold leading-none tracking-[-0.03em] text-[#9e1409]">
            {cents(FALLS_TO)}
            <span className="ml-2 align-baseline text-[0.3em] font-medium tracking-normal text-[#1a1713]/60">
              NVDA
            </span>
          </p>
          <p className="mt-5 max-w-[34ch] text-[17px] leading-snug text-[#1a1713]/85">
            {t.rich("ifFalls", {
              price: cents(FALLS_TO),
              share: `${SOLD_SHARE}%`,
              n: (chunks) => <span className="num font-semibold">{chunks}</span>,
            })}
          </p>
          <p className="mt-7 border-t border-[#dcd6ca] pt-5 text-[14px] leading-snug text-[#1a1713]/70">
            {t.rich("borrowed", {
              amount: cents(BORROWED),
              big: (chunks) => <span className="num text-[19px] font-semibold text-[#1a1713]">{chunks}</span>,
              unit: (chunks) => <span className="text-[12px] font-medium">{chunks}</span>,
            })}
          </p>
        </div>

        <div className="-mx-1 overflow-x-auto">
          <table className="w-full border-collapse text-[13px] sm:text-[14px]">
            <caption className="sr-only">{t("caption")}</caption>
            <thead>
              <tr className="text-right font-mono text-[10.5px] uppercase tracking-[0.1em] text-[#1a1713]/55">
                <th scope="col" className="px-1 pb-2.5 text-left font-medium">{t("token")}</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">{t("price")}</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">{t("band")}</th>
                <th scope="col" className="hidden px-1 pb-2.5 font-medium sm:table-cell">{t("change")}</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">{t("value")}</th>
              </tr>
            </thead>
            <tbody>
              {HOLDINGS.map((h) => {
                const up = h.move > 0;
                return (
                <tr key={h.token} className="h-10 border-t border-[#dcd6ca] text-right">
                  <th scope="row" className="px-1 text-left font-semibold">{h.token}</th>
                  <td className="num px-1 font-semibold">{cents(h.price)}</td>
                  <td className="num px-1 text-[#16326b]">±{format.number(h.band, { minimumFractionDigits: 1, maximumFractionDigits: 1 })}%</td>
                  <td className={`num hidden px-1 sm:table-cell ${up ? "text-[#1c6b4a]" : "text-[#b0503f]"}`}>
                    <svg
                      viewBox="0 0 8 8"
                      aria-hidden="true"
                      className={`mr-1 inline-block h-[7px] w-[7px] ${up ? "" : "rotate-180"}`}
                      fill="currentColor"
                    >
                      <path d="M4 1 L7.5 7 H0.5 Z" />
                    </svg>
                    {signed(h.move)}
                  </td>
                  <td className="num px-1 font-semibold">{cents(h.value)}</td>
                </tr>
                );
              })}
            </tbody>
            <tfoot>
              <tr className="h-11 border-t border-[#1a1713]/35 text-right">
                <th scope="row" colSpan={2} className="px-1 text-left font-medium text-[#1a1713]/60">
                  {t("total")}
                </th>
                <td className="px-1" />
                <td className="hidden px-1 sm:table-cell" />
                <td className="num px-1 text-[15px] font-semibold">{cents(TOTAL)}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <p className="border-t border-[#dcd6ca] bg-[#f5f2ec] px-5 py-2.5 font-mono text-[10.5px] uppercase tracking-[0.1em] text-[#1a1713]/60 sm:px-7">
        {t("note")}
      </p>
    </div>
  );
}
