import type { Metadata } from "next";
import { Entry, PlainPage } from "../components/PlainPage";

export const metadata: Metadata = {
  title: "Methodology",
  alternates: { canonical: "/methodology" },
  description: "How every figure on the Tapehouse site was measured, and what the illustrative model assumes.",
};

export default function Methodology() {
  return (
    <PlainPage
      title="How we measured it"
      lede="Every figure on this site comes from somewhere you can check. This is where, and how."
    >
      <Entry heading="$109.2 billion against $31.2 billion">
        <p>
          Morgan Stanley’s annual report on Form 10-K for fiscal year 2025, Wealth Management loans:
          “Securities-based lending and Other” of $109,201 million, and margin loans of $31,214 million.
          The first category includes lending other than securities-based loans, so read it as an upper
          bound. That most of it is borrowed to avoid selling is our reading, not the filing’s.
        </p>
      </Entry>
      <Entry heading="55.7 to 80.8 hours of silence">
        <p>
          The gaps between consecutive updates of the reference equity price feed on Robinhood Chain, read
          on-chain across weekends in August and September 2026. The longest, 80.8 hours, spans the Labor
          Day weekend. The weekend drawn on the overview is 28–31 August 2026: last update Friday 16:18 UTC,
          next update Monday 00:00 UTC.
        </p>
      </Entry>
      <Entry heading="288 of 288">
        <p>
          An around-the-clock price source, sampled at ten-minute resolution across one full 48-hour
          weekend. All 288 expected points were present, and the longest gap between two of them was 600
          seconds.
        </p>
      </Entry>
      <Entry heading="29.9% of SPY transfers">
        <p>
          The share of SPY Stock Token transfers on Robinhood Chain that occurred on a Saturday or a Sunday,
          over 69 consecutive days. A perfectly even week would give 28.6%.
        </p>
      </Entry>
      <Entry heading="The 62.5% flat rule">
        <p>
          The per-token limit that lending markets on this chain applied to large-cap Stock Tokens such as
          SPY, NVDA and TSLA, as read from their public interfaces in September 2026.
        </p>
      </Entry>
      <Entry heading="The illustrative portfolio model">
        <p>
          The interactive model takes the worst one percent of outcomes over a three-day closed market,
          from assumed volatilities and correlations, doubles it as a stress, and holds back a further ten
          percent. The assumptions exist to show the mechanism. They are not the parameters Tapehouse
          lends on, which are published with the backtest.
        </p>
      </Entry>
      <Entry heading="The price band">
        <p>
          The band answers one question: what is this token worth right now, and how sure are we. Its
          centre is the live around-the-clock price. Three things widen it: live sources that disagree, a
          price that is moving fast, and fewer sources answering. A reference feed that has gone quiet for
          the weekend is set aside, never averaged in, and never used to widen the band on its own.
        </p>
        <p>
          How far the price might move before anyone can act is a different question. It belongs to the
          margin engine, not to the band.
        </p>
        <p>
          With one source live, as at a weekend, the band is ±0.55%. That calibration is provisional: it
          rests on a small number of reopenings in a calm month.
        </p>
      </Entry>
      <Entry heading="The live record">
        <p>
          For each reopening we take the band from the last sample before 00:00 UTC on Monday, when the
          reference equity feed on Robinhood Chain resumes, and compare its centre with the first price
          that feed publishes. Nothing later than the band’s own timestamp is used.
        </p>
        <p>
          The first weekend, 21 September 2026, was graded after the fact from our own minute-by-minute
          capture and is marked as such. Later entries are sealed and timestamped before the reopening.
        </p>
      </Entry>
      <Entry heading="The stress test">
        <p>
          The same band fed a single, deliberately weak source: the thirty-minute time-weighted price of
          the Uniswap v3 pool for each token against USDG, built from swap events read from the chain and
          ending at the reopening. It is graded the same way, against the first reference print after a
          silence of thirty hours or more. It covers the nine weekends for which those pools existed, from
          27 July to 21 September 2026: thirty reopenings across NVDA, TSLA, AAPL and SPY.
        </p>
        <p>
          At a weekend the on-chain price carries a premium or a discount of up to about 1.2%, because
          liquidity is thin. That gap is not noise to a lender: it is the price a liquidation would
          actually get. Tapehouse accounts for it in the margin engine.
        </p>
      </Entry>
      <Entry heading="Friday’s price, carried through the weekend">
        <p>
          The change from the last reference print before each weekend to the first one after it, over all
          forty-eight reopenings of NVDA, TSLA, AAPL and SPY since 1 July 2026, read from the chain. The
          median is 0.34% and the largest 1.33%.
        </p>
      </Entry>
    </PlainPage>
  );
}
