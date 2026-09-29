# Test data

`matrix-vectors.json` holds 42 correlation matrices in basis points, each given by its size `n` and its upper triangle row by row, with whether it is positive definite. They include the identity from 1 to 8 assets, the launch assets' floors and first values, their first values without SPY, pairs that move as one (singular), a matrix whose pairs are each plausible but which is not positive definite, a constant correlation of 0.9999 over 8 assets, 12 matrices from random two-factor models and 12 of random entries.

Every case was classified twice, independently of this crate: by exact integer elimination, and by the smallest eigenvalue in floating point. They agree on every case; where the smallest eigenvalue is zero, the matrix is singular and the exact elimination finds a leading minor of exactly zero.

`scenario-vectors.json` holds the scenario sets of nine configurations:
- the launch assets at their first values, and the same assets in another order;
- the testnet's five, without SPY, so with the equal-weighted market;
- the dev node's three;
- SPY alone;
- eight synthetic assets;
- three assets whose correlation matrix is at the edge of positive definiteness;
- two matrices outside it, which the program refuses but the set is still defined for: one whose factor clamps an entry to ±1, and a singular one whose last pivot falls below the floor.

For each configuration it gives the set's keccak-256 at horizons of 0, 1, 2 and 3.5 days and for 32, 64, 128 and 256 lattice points: 144 digests. It also gives fourteen scenarios of the launch set, one from each family and its edges. Every value was computed in unbounded integer arithmetic, independently of this crate. The launch assets and their reordering have the same digest.

`requirement-vectors.json` holds what the engine reads from Robinhood Chain for the six launch assets at block 75,093,578 (28 September 2026): each asset's pool from `parameters.json`, its tokens, fee and decimals, its 30-minute observation (`observe([1800, 0])`), and the ETH/USD round. From that state it gives 120 requirements:
- ten portfolios: one token, long and short positions, a hedge, all six assets, and positions past their pool's depth;
- at horizons of 2 and 3.5 days, with and without a closure;
- with the pools read; with them unread, where each asset held is still charged its pool's fee and governed depth and sets its bit; and without them, as on Arbitrum One, where every asset held sets its bit.

It also gives the requirement the dev-node suite expects from its three assets, with NVDA's stub pool repeating the real pool's mean tick and liquidity. Every value was computed in unbounded integer arithmetic, independently of this crate.

`session-timeline.json` holds the session the band's `session()` returns every 15 minutes from 27 August to 28 September 2026: 3,071 instants, each `[now_ms, open, nyse, nyse_next, change_ms, boundary_ms]` with `open` 0 unknown, 1 closed and 2 open. It was derived from the New York market status RedStone signed over those days, with the same reference that checks the band's `session-vectors.json`. It holds five Friday post-markets and Labor Day.
