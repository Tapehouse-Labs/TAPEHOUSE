//! A 24/7 leg made from an index: the asset's last Chainlink print, moved by the index since.
//!
//! SPY has no 24/7 feed, but the S&P 500 index has one. The anchor pairs a Chainlink print with the
//! index price signed within [`ANCHOR_MAX_LAG_MS`] of it, so the leg never assumes a fixed ratio
//! between the fund and the index.

/// Extra half-width of a leg made from an index, for the gap between the fund and its index, in basis
/// points.
pub const INDEX_BASIS_BPS: u64 = 25;

/// Maximum distance between a Chainlink print and the index package that anchors it, in milliseconds.
pub const ANCHOR_MAX_LAG_MS: u64 = 120_000;

/// A Chainlink print and the index price at that print. All zero before the first anchor.
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq)]
pub struct Anchor {
    /// The Chainlink print.
    pub cl_px: u64,
    /// The index price signed within [`ANCHOR_MAX_LAG_MS`] of the print.
    pub index_px: u64,
    /// The print's `updatedAt`, in seconds.
    pub at_s: u64,
}

/// The anchor after an index price `index_px` with package timestamp `index_ms` is written, when
/// Chainlink's latest print is `cl_px` at `cl_at_s`. `None` keeps the anchor: the print is not newer
/// or repeats the anchored price, a price is zero, or the package is too far from the print. An older
/// anchor still pairs a print with the index at that print.
pub fn next_anchor(
    anchor: Anchor,
    cl_px: u64,
    cl_at_s: u64,
    index_px: u64,
    index_ms: u64,
) -> Option<Anchor> {
    let near =
        u128::from(index_ms).abs_diff(u128::from(cl_at_s) * 1000) <= u128::from(ANCHOR_MAX_LAG_MS);
    let new = cl_at_s > anchor.at_s && cl_px != anchor.cl_px;
    (cl_px > 0 && index_px > 0 && new && near).then_some(Anchor {
        cl_px,
        index_px,
        at_s: cl_at_s,
    })
}

/// The leg's price for the index at `index_px`, rounded down. Zero without an anchor, or when it does
/// not fit in `u64`.
pub fn leg_px(anchor: Anchor, index_px: u64) -> u64 {
    if anchor.index_px == 0 {
        return 0;
    }
    let px = u128::from(anchor.cl_px) * u128::from(index_px) / u128::from(anchor.index_px);
    u64::try_from(px).unwrap_or(0)
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::Value;

    const VECTORS: &str = include_str!("../testdata/quote-vectors.json");

    fn int(value: &Value) -> u64 {
        value.as_str().unwrap().parse().unwrap()
    }

    #[test]
    fn legs_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["index_legs"].as_array().unwrap() {
            let anchor = Anchor {
                cl_px: int(&case["anchor_cl_px"]),
                index_px: int(&case["anchor_index_px"]),
                at_s: 0,
            };
            assert_eq!(
                leg_px(anchor, int(&case["index_px"])),
                int(&case["expected"]),
                "{}",
                case["name"]
            );
        }
    }

    #[test]
    fn a_print_beyond_u64_milliseconds_is_far_from_any_package() {
        let cl_at_s = u64::MAX / 1000 + 1000;
        let anchor = Anchor {
            cl_px: 1,
            index_px: 1,
            at_s: cl_at_s - 1,
        };
        assert_eq!(next_anchor(anchor, 2, cl_at_s, 3, u64::MAX), None);
    }

    #[test]
    fn anchors_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["anchors"].as_array().unwrap() {
            let anchor = |v: &Value| Anchor {
                cl_px: int(&v["cl_px"]),
                index_px: int(&v["index_px"]),
                at_s: int(&v["at_s"]),
            };
            let before = anchor(&case["anchor"]);
            let after = next_anchor(
                before,
                int(&case["cl_px"]),
                int(&case["cl_at_s"]),
                int(&case["index_px"]),
                int(&case["index_ms"]),
            )
            .unwrap_or(before);
            assert_eq!(after, anchor(&case["expected"]), "{}", case["name"]);
        }
    }
}
