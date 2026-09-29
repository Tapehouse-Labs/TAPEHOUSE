# Test data

`matrix-vectors.json` holds 42 correlation matrices in basis points, each given by its size `n` and its upper triangle row by row, with whether it is positive definite. They include the identity from 1 to 8 assets, the launch assets' floors and first values, their first values without SPY, pairs that move as one (singular), a matrix whose pairs are each plausible but which is not positive definite, a constant correlation of 0.9999 over 8 assets, 12 matrices from random two-factor models and 12 of random entries.

Every case was classified twice, independently of this crate: by exact integer elimination, and by the smallest eigenvalue in floating point. They agree on every case; where the smallest eigenvalue is zero, the matrix is singular and the exact elimination finds a leading minor of exactly zero.
