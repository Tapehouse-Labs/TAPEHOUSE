# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

**Primary — the holder who will not sell.** Mid-thirties, technically fluent, holds a diversified portfolio of Stock Tokens (for example NVDA, TSLA, AAPL and SPY) on Robinhood Chain and arrived from a retail brokerage app. Uses a browser wallet, understands loan-to-value, does not read contracts. Needs liquidity without selling: selling means a tax event and losing the position. Her blocking fear is a surprise liquidation, ahead of the interest rate.

**Primary — the collateral optimiser.** Late twenties, crypto-native, compares risk parameters across protocols and reads the methodology. Wants capital efficiency, public parameters, a measurable track record, and the option to isolate a position.

**Secondary — the integrator.** A protocol, wallet or aggregator that needs a price and a risk number for Stock Tokens and does not want to build them. Reads the contract before the website and decides in about twenty minutes.

## Product Purpose

Tapehouse is a portfolio margin account for Stock Tokens. A user deposits their whole portfolio and borrows against it as one position. Success is a user who borrows on a Sunday, knew their liquidation price before signing, and was never surprised by it.

## Positioning

Portfolio margin for Stock Tokens. The whole portfolio backs one loan, so a diversified portfolio earns more borrowing capacity than the same assets priced one by one, and a concentrated one earns less. Every limit is published before it matters: the liquidation price is stated as a number before the user signs, the price keeps updating while the underlying market is closed, and the record of when the price band was wrong is public.

## Operating Context

- Underlying equity markets are closed for most of the week, while Stock Tokens trade on-chain around the clock. A meaningful share of on-chain activity happens at weekends.
- Reference price feeds for equities stop updating while the market is closed. The product's price band widens when live sources disagree, when the price moves fast, and when fewer sources answer; a source that has gone quiet is set aside, never averaged in. It declares its own state: open, closed, degraded or halted.
- The session state is always visible next to any money figure, as a compact status seal stating the state, the hours elapsed and the number of live price sources.
- The launch surface is Robinhood Chain testnet. The landing page's primary action opens the app on testnet.

## Capabilities and Constraints

- Deterministic on-chain confidence band around each Stock Token price.
- Scenario-based portfolio margin: capacity is derived from the expected loss of the whole portfolio in its worst one percent of scenarios, not from a per-asset loan-to-value sum.
- Session-aware margin accounts with deferred liquidation and partial liquidation.
- An adapter exposing the price and the band through a standard aggregator interface, with a runnable example.
- Terminology is fixed: "Stock Tokens", "the band", "the session", "capacity", "liquidation price", "the seal".
- Undecided: the headline capacity figure. It comes from the backtest and must not be stated until that result exists.
- Tapehouse is not a broker. A Stock Token is a debt security of its issuer, not a share: no shareholder rights, cash redemption only, and not offered to United States persons.

## Brand Commitments

- Name: Tapehouse. Category line, always in the same words: "Portfolio margin for Stock Tokens."
- Organising idea: "Published before it matters."
- Mark: the continuity mark — a support with a line that crosses it and keeps rising. "The market's line ends. Ours keeps going." Always shown with the name.
- Colour: warm bone paper, warm ink, and a single deep registry navy accent. Gain and loss are their own semantic colours, distinct from success and error.
- Type: Archivo for display, Inter Tight for interface and figures, Spline Sans Mono for addresses and labels. Every figure that decides money is set in tabular numerals.
- Voice: second person, concrete verbs, a human in every sentence. Every public figure carries a footnote with its source and method. Banned words: "institutional-grade", "seamless", "revolutionary", "the future of".
- The landing is a narrative that follows one example portfolio from top to bottom. The product interface is dense and exact. Documentation is plain.
- Light and dark themes both ship from day one.

## Evidence on Hand

- Measured: the reference equity feed on this chain falls silent for 55.7 to 80.8 hours every weekend.
- Measured: an around-the-clock price source served 288 of 288 expected data points across a full weekend.
- Measured: 29.9% of SPY token transfers on this chain occur at weekends.
- A live record of the band graded against the first fresh market print after each weekend (one weekend so far, graded after the fact and marked as such).
- A reproducible on-chain stress test: the band fed only the DEX price, 30 reopenings over 9 weekends, wrong 8 times.
- Absent, and not to be fabricated: customers, testimonials, press, audits, total value locked, and the headline capacity figure.

## Product Principles

1. Publish the number before it matters, and grade it afterwards.
2. A person is always in front of the mechanism.
3. Warm in voice, exact in figures; the footnote is the hinge between the two.
4. State the failure mode out loud: a degraded source is shown, never hidden.
5. The data decides the design, never the reverse.

## Accessibility & Inclusion

WCAG 2.2 AA. Gain and loss must never rely on colour alone. Motion respects reduced-motion preferences. Product copy is English.
