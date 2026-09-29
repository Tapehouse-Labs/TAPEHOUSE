//! Validity of a correlation matrix given in basis points, with 10 000 on the diagonal.

use stylus_sdk::alloy_primitives::I256;

/// Correlation of an asset with itself, in basis points.
pub const ONE: u16 = 10_000;

/// Index of the pair `(i, j)`, `i < j`, in the upper triangle of an `n × n` matrix read row by row.
pub fn pair(n: usize, i: usize, j: usize) -> usize {
    i * n - i * (i + 1) / 2 + (j - i - 1)
}

/// Whether the symmetric matrix with a unit diagonal and the upper triangle `upper` is positive definite.
///
/// Sylvester's criterion: every leading principal minor is positive. Bareiss's fraction-free elimination
/// produces them exactly, the k-th pivot being the k × k minor, so no rounding decides the answer. By
/// Hadamard's inequality every intermediate stays below 2^255 for up to 9 assets.
pub fn positive_definite(n: usize, upper: &[u16]) -> bool {
    let mut a = vec![I256::ZERO; n * n];
    for i in 0..n {
        a[i * n + i] = I256::try_from(ONE).unwrap();
        for j in i + 1..n {
            let value = I256::try_from(upper[pair(n, i, j)]).unwrap();
            a[i * n + j] = value;
            a[j * n + i] = value;
        }
    }
    let mut previous = I256::ONE;
    for k in 0..n {
        let pivot = a[k * n + k];
        if pivot <= I256::ZERO {
            return false;
        }
        for i in k + 1..n {
            for j in k + 1..n {
                a[i * n + j] = (a[i * n + j] * pivot - a[i * n + k] * a[k * n + j]) / previous;
            }
        }
        previous = pivot;
    }
    true
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;

    #[test]
    fn pairs_are_numbered_row_by_row() {
        let n = 6;
        let mut k = 0;
        for i in 0..n {
            for j in i + 1..n {
                assert_eq!(pair(n, i, j), k);
                k += 1;
            }
        }
        assert_eq!(k, 15);
    }

    #[test]
    fn matrices_are_judged_as_the_reference_judges_them() {
        let vectors: serde_json::Value =
            serde_json::from_str(include_str!("../testdata/matrix-vectors.json")).unwrap();
        let cases = vectors["cases"].as_array().unwrap();
        assert_eq!(cases.len(), 42);
        for case in cases {
            let n = case["n"].as_u64().unwrap() as usize;
            let upper: Vec<u16> = case["upper"]
                .as_array()
                .unwrap()
                .iter()
                .map(|v| v.as_u64().unwrap() as u16)
                .collect();
            assert_eq!(
                positive_definite(n, &upper),
                case["positiveDefinite"].as_bool().unwrap(),
                "{}",
                case["name"]
            );
        }
    }

    fn matrix() -> impl Strategy<Value = (usize, Vec<u16>)> {
        (1usize..=8)
            .prop_flat_map(|n| (Just(n), prop::collection::vec(1u16..=ONE, n * (n - 1) / 2)))
    }

    fn permuted(n: usize, upper: &[u16], order: &[usize]) -> Vec<u16> {
        let at = |a: usize, b: usize| upper[pair(n, a.min(b), a.max(b))];
        let mut out = vec![];
        for i in 0..order.len() {
            for j in i + 1..order.len() {
                out.push(at(order[i], order[j]));
            }
        }
        out
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, ..ProptestConfig::default() })]

        #[test]
        fn reordering_the_assets_keeps_the_answer(
            (n, upper) in matrix(),
            order in Just((0..8).collect::<Vec<usize>>()).prop_shuffle(),
        ) {
            let order: Vec<usize> = order.into_iter().filter(|&i| i < n).collect();
            prop_assert_eq!(positive_definite(n, &upper), positive_definite(n, &permuted(n, &upper, &order)));
        }

        #[test]
        fn every_leading_block_of_a_positive_definite_matrix_is_positive_definite((n, upper) in matrix()) {
            if positive_definite(n, &upper) {
                for m in 1..n {
                    prop_assert!(positive_definite(m, &permuted(n, &upper, &(0..m).collect::<Vec<_>>())));
                }
            }
        }

        #[test]
        fn a_constant_correlation_below_one_is_positive_definite(n in 1usize..=8, c in 0u16..ONE) {
            prop_assert!(positive_definite(n, &vec![c; n * (n - 1) / 2]));
        }
    }
}
