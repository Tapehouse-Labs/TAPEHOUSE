import { Artwork } from "./Artwork";

const FOLDS = [
  {
    q: "How the band is built",
    a: "Three things widen it: live sources that disagree, a price that is moving fast, and fewer sources answering. A source that has gone quiet for the weekend is set aside, never averaged in. The parameters are public, so anyone can recompute last Sunday’s band and check it against ours.",
  },
  {
    q: "What happens when a source fails",
    a: "The band widens and the seal says DEGRADED, with the count of sources still alive. It never fills the hole with a number it doesn’t have.",
  },
  {
    q: "How you integrate it",
    a: "Each feed answers the aggregator interface your contracts already read, and adds one call that returns the band, the session and the live-source count.",
  },
  {
    q: "What changes, and when you hear about it",
    a: "Parameters change only with notice, on-chain, ahead of taking effect. You cannot build on something that moves without warning, so it doesn’t.",
  },
];

const CODE: [string, string][] = [
  ["c", "// Read a Stock Token price together with its band."],
  ["", "ITapehouseFeed feed = ITapehouseFeed(NVDA_FEED);"],
  ["", ""],
  ["c", "// The interface your contracts already read."],
  ["", "(, int256 price, , uint256 updatedAt, ) = feed.latestRoundData();"],
  ["", ""],
  ["c", "// One more call: the band, the session, the live sources."],
  ["", "(uint256 lower, uint256 upper, Session session, uint8 live)"],
  ["", "    = feed.latestBand();"],
  ["", ""],
  ["", "require(session != Session.Degraded, \"source degraded\");"],
  ["c", "// Lend against the conservative edge, not the midpoint."],
  ["", "uint256 collateralPrice = lower;"],
];

/** Chapter five: the method, in an engineer's register. */
export function Method() {
  return (
    <section className="field field-night" aria-labelledby="ch-method">
      <div className="shell py-[clamp(6rem,13vw,11rem)]">
        <div className="grid gap-x-16 gap-y-10 lg:grid-cols-[1fr_17rem]">
          <div>
            <h2 id="ch-method" className="chapter-title max-w-[15ch]">
              A number you can check. Not a promise.
            </h2>
            <p className="lede mt-8 max-w-[56ch]">
              If you lend, build a wallet, or run a protocol on this chain, you need a price for Stock
              Tokens on a Sunday. Here is ours, with everything required to prove us wrong.
            </p>
          </div>
          <Artwork src="/art/i2-socket.jpg" alt="A plug with three pegs hovering just above the block it fits into." className="order-first lg:order-none" />
        </div>

        <div className="mt-[clamp(3rem,6vw,5rem)] grid gap-x-16 gap-y-12 lg:grid-cols-[1fr_1.1fr]">
          <div>
            {FOLDS.map((f, i) => (
              <details key={f.q} className="fold" open={i === 0} name="method">
                <summary>
                  {f.q}
                  <span className="plus" aria-hidden="true" />
                </summary>
                <div>{f.a}</div>
              </details>
            ))}
          </div>

          <figure className="overflow-hidden rounded-[8px] bg-panel shadow-[inset_0_0_0_1px_var(--rule)]">
            <figcaption className="flex items-center justify-between gap-4 border-b border-rule px-5 py-3 font-mono text-[11px] uppercase tracking-[0.12em] text-muted">
              <span>Consumer.sol</span>
              <span>Interface preview</span>
            </figcaption>
            <pre className="overflow-x-auto px-5 py-5 font-mono text-[13px] leading-[1.75] text-body" tabIndex={0}>
              <code>
                {CODE.map(([kind, line], i) => (
                  <span key={i} className={`block ${kind === "c" ? "text-faint" : ""}`}>
                    {line || " "}
                  </span>
                ))}
              </code>
            </pre>
          </figure>
        </div>
        <p className="mt-6 max-w-[70ch] text-[13px] leading-relaxed text-muted">
          The interface is a preview and may change before mainnet.
        </p>
      </div>
    </section>
  );
}
