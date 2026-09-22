import { Mark } from "./Mark";
import { Seal } from "./Seal";

/**
 * The lit screen the camera pushes into. Its surface never inverts with the
 * theme: it is a screen glowing in a dark room. All figures are illustrative
 * and labelled as such on the screen itself.
 */
const HOLDINGS = [
  { token: "NVDA", units: "120.00", price: "182.50", band: "±1.1%", value: "21,900.00", move: "+1.24%", up: true },
  { token: "TSLA", units: "45.00", price: "396.12", band: "±1.6%", value: "17,825.40", move: "−0.87%", up: false },
  { token: "AAPL", units: "60.00", price: "241.08", band: "±0.9%", value: "14,464.80", move: "+0.31%", up: true },
  { token: "SPY", units: "18.00", price: "766.14", band: "±0.6%", value: "13,790.52", move: "−0.12%", up: false },
];

export function PortfolioScreen() {
  return (
    <div
      className="overflow-hidden rounded-[10px] bg-[#fbfaf7] text-[#1a1713] shadow-[0_40px_120px_-20px_rgb(0_0_0/0.65),0_0_0_1px_rgb(245_242_236/0.14)]"
      aria-label="Example Tapehouse account"
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
          <p className="text-[13px] text-[#1a1713]/60">Your liquidation price, stated before you sign</p>
          <p className="num mt-2 font-sans text-[clamp(2.75rem,7vw,4.5rem)] font-semibold leading-none tracking-[-0.03em] text-[#9e1409]">
            141.08
            <span className="ml-2 align-baseline text-[0.3em] font-medium tracking-normal text-[#1a1713]/60">
              NVDA
            </span>
          </p>
          <p className="mt-5 max-w-[34ch] text-[17px] leading-snug text-[#1a1713]/85">
            If NVDA falls to <span className="num font-semibold">141.08</span> we sell{" "}
            <span className="num font-semibold">12%</span> of your SPY, not your whole position.
          </p>
          <p className="mt-7 border-t border-[#dcd6ca] pt-5 text-[14px] leading-snug text-[#1a1713]/70">
            <span className="num text-[19px] font-semibold text-[#1a1713]">10,000.00</span>{" "}
            <span className="text-[12px] font-medium">USDG</span> borrowed, backed by all four holdings
            as one position.
          </p>
        </div>

        <div className="-mx-1 overflow-x-auto">
          <table className="w-full border-collapse text-[13px] sm:text-[14px]">
            <caption className="sr-only">Holdings backing the loan</caption>
            <thead>
              <tr className="text-right font-mono text-[10.5px] uppercase tracking-[0.1em] text-[#1a1713]/55">
                <th scope="col" className="px-1 pb-2.5 text-left font-medium">Token</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">Price</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">Band</th>
                <th scope="col" className="hidden px-1 pb-2.5 font-medium sm:table-cell">24h</th>
                <th scope="col" className="px-1 pb-2.5 font-medium">Value</th>
              </tr>
            </thead>
            <tbody>
              {HOLDINGS.map((h) => (
                <tr key={h.token} className="h-10 border-t border-[#dcd6ca] text-right">
                  <th scope="row" className="px-1 text-left font-semibold">{h.token}</th>
                  <td className="num px-1 font-semibold">{h.price}</td>
                  <td className="num px-1 text-[#16326b]">{h.band}</td>
                  <td className={`num hidden px-1 sm:table-cell ${h.up ? "text-[#1c6b4a]" : "text-[#b0503f]"}`}>
                    <svg
                      viewBox="0 0 8 8"
                      aria-hidden="true"
                      className={`mr-1 inline-block h-[7px] w-[7px] ${h.up ? "" : "rotate-180"}`}
                      fill="currentColor"
                    >
                      <path d="M4 1 L7.5 7 H0.5 Z" />
                    </svg>
                    {h.move}
                  </td>
                  <td className="num px-1 font-semibold">{h.value}</td>
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr className="h-11 border-t border-[#1a1713]/35 text-right">
                <th scope="row" colSpan={2} className="px-1 text-left font-medium text-[#1a1713]/60">
                  Whole portfolio
                </th>
                <td className="px-1" />
                <td className="hidden px-1 sm:table-cell" />
                <td className="num px-1 text-[15px] font-semibold">67,980.72</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <p className="border-t border-[#dcd6ca] bg-[#f5f2ec] px-5 py-2.5 font-mono text-[10.5px] uppercase tracking-[0.1em] text-[#1a1713]/60 sm:px-7">
        Example portfolio · figures are illustrative
      </p>
    </div>
  );
}
