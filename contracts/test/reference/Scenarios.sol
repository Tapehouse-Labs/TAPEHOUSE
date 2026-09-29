// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

/// @notice The margin engine's parameters, in the stored order of the assets.
struct Parameters {
    bytes32[] symbols;
    uint32[] volatilities;
    uint16[] correlations;
    uint32[] gaps;
    /// @dev The market asset's position plus one, or zero for the equal-weighted portfolio.
    uint256 market;
}

/// @notice Everything a row needs, computed once per set.
struct Set {
    Parameters parameters;
    uint256 size;
    uint256[] order;
    uint256[] rank;
    int256[] factor;
    int256 root;
    int256[] moves;
    int256[] means;
    uint256[] generators;
    uint256[] shifts;
}

/// @notice The engine's scenario set, as the Stylus program draws it: joint, comonotone and independent draws
/// on a shifted Korobov lattice, the market moves, and every asset's weekend gap at once and alone.
library Scenarios {
    int256 internal constant S = 1e12;
    int256 internal constant PPM = 1e6;
    int256 internal constant ONE = 10_000;
    int256 internal constant PIVOT_FLOOR = S / PPM;
    int256 internal constant Z_MARKET = 2_665_214;

    bytes internal constant MEANS_32 = hex"ffdda24dffe650b5ffea530bffed3642ffef8b82fff18b50fff35229fff4f0a4fff67160fff7dba9fff934cdfffa80d5"
        hex"fffbc2f9fffcfde0fffe33d8ffff66ee";
    bytes internal constant MEANS_64 = hex"ffd9b297ffe19203ffe51491ffe78cd8ffe98108ffeb250effec92eaffedd99bffef02a0fff01464fff11374fff2032c"
        hex"fff2e619fff3be3afff48d24fff55424fff61449fff6ce77fff78372fff833e0fff8e052fff98948fffa2f32fffad278"
        hex"fffb7375fffc127dfffcafdffffd4be2fffde6cefffe80e2ffff1a5dffffb37e";
    bytes internal constant MEANS_128 = hex"ffd61275ffdd52b9ffe07a23ffe2a9e3ffe45f1cffe5ca07ffe702b4ffe816fcffe90f9bffe9f275ffeac3b4ffeb8668"
        hex"ffec3ce3ffece8f1ffed8bffffee2736ffeebb86ffef49b9ffefd277fff05650fff0d5bdfff1512afff1c8f2fff23d65"
        hex"fff2aeccfff31d65fff38969fff3f30afff45a75fff4bfd4fff5234bfff584fdfff5e508fff64389fff6a09afff6fc54"
        hex"fff756ccfff7b018fff8084afff85f75fff8b5aafff90af9fff95f70fff9b31ffffa0611fffa5854fffaa9f4fffafafc"
        hex"fffb4b78fffb9b72fffbeaf3fffc3a07fffc88b5fffcd708fffd2508fffd72bdfffdc031fffe0d6bfffe5a72fffea751"
        hex"fffef40dffff40aeffff8d3dffffd9c0";
    bytes internal constant MEANS_256 = hex"ffd2b1f2ffd972f7ffdc5564ffde500fffdfd87cffe11bcbffe23071ffe32354ffe3fc93ffe4c1a5ffe57665ffe61da9"
        hex"ffe6b995ffe74bd3ffe7d5b3ffe85844ffe8d466ffe94ad0ffe9bc1cffea28ceffea9155ffeaf613ffeb575bffebb576"
        hex"ffec10a5ffec6921ffecbf1cffed12c5ffed6443ffedb3bcffee014fffee4d1dffee973fffeedfceffef26e2ffef6c90"
        hex"ffefb0ebffeff404fff035ecfff076b3fff0b667fff0f514fff132c7fff16f8cfff1ab6efff1e676fff220adfff25a1e"
        hex"fff292cffff2cac9fff30214fff338b7fff36eb7fff3a41bfff3d8eafff40d29fff440defff4740dfff4a6bbfff4d8ed"
        hex"fff50aa7fff53beffff56cc6fff59d33fff5cd38fff5fcd8fff62c18fff65afafff68982fff6b7b3fff6e58ffff71319"
        hex"fff74055fff76d44fff799e9fff7c647fff7f25ffff81e35fff849cafff87521fff8a03bfff8cb1afff8f5c1fff92031"
        hex"fff94a6cfff97474fff99e4bfff9c7f2fff9f16bfffa1ab7fffa43d8fffa6cd0fffa959ffffabe48fffae6ccfffb0f2c"
        hex"fffb376afffb5f86fffb8783fffbaf61fffbd721fffbfec5fffc264ffffc4dbefffc7515fffc9c55fffcc37efffcea92"
        hex"fffd1191fffd387efffd5f59fffd8622fffdacdbfffdd386fffdfa23fffe20b2fffe4736fffe6daffffe941efffeba83"
        hex"fffee0e1ffff0738ffff2d88ffff53d4ffff7a1bffffa05fffffc6a0ffffece0";

    /// @notice Scenarios in a set of lattice size `size` over `n` assets.
    function count(uint256 size, uint256 n) internal pure returns (uint256) {
        return 3 * size + 4 + 2 * n;
    }

    /// @notice Index of the pair `(i, j)`, `i < j`, in the upper triangle of an `n × n` matrix read row by row.
    function pair(uint256 n, uint256 i, uint256 j) internal pure returns (uint256) {
        return i * n - i * (i + 1) / 2 + (j - i - 1);
    }

    /// @notice Integer square root, rounded down.
    function isqrt(uint256 n) internal pure returns (uint256 x) {
        if (n == 0) return 0;
        x = n;
        uint256 y = (n + 1) / 2;
        while (y < x) {
            x = y;
            y = (x + n / x) / 2;
        }
    }

    /// @notice Prepares the set for lattice size `size` over a horizon of `horizon` seconds.
    function create(Parameters memory p, uint256 size, uint64 horizon) internal pure returns (Set memory set) {
        uint256 n = p.symbols.length;
        set.parameters = p;
        set.size = size;
        set.order = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            set.order[i] = i;
        }
        for (uint256 i = 1; i < n; ++i) {
            for (uint256 k = i; k > 0 && p.symbols[set.order[k - 1]] > p.symbols[set.order[k]]; --k) {
                (set.order[k - 1], set.order[k]) = (set.order[k], set.order[k - 1]);
            }
        }
        set.rank = new uint256[](n);
        for (uint256 r; r < n; ++r) {
            set.rank[set.order[r]] = r;
        }
        set.factor = _factor(p, set.order);
        set.root = int256(isqrt(uint256(horizon) * 1e12 / 86_400));
        set.moves = _moves(p, set.root);
        set.means = _means(size);
        set.generators = new uint256[](n);
        set.shifts = new uint256[](n);
        uint256 generator = size == 32 ? 5 : size == 64 ? 11 : 13;
        uint256 g = 1;
        for (uint256 d; d < n; ++d) {
            set.generators[d] = g;
            set.shifts[d] = size * ((d * 618_033_988_749_894_848) % 1e18) / 1e18;
            g = g * generator % size;
        }
    }

    /// @notice Scenario `index` as returns in millionths, in the stored order of the assets, written to `row`;
    /// `draws` is scratch space of one word per asset.
    function fill(Set memory set, uint256 index, int256[] memory row, int256[] memory draws) internal pure {
        unchecked {
            Parameters memory p = set.parameters;
            uint256 n = p.symbols.length;
            uint256 size = set.size;
            if (index < size) {
                for (uint256 d; d < n; ++d) {
                    draws[d] = _draw(set, index, d);
                }
                for (uint256 c; c < n; ++c) {
                    int256 x;
                    for (uint256 k; k <= c; ++k) {
                        x += set.factor[c * n + k] * draws[k];
                    }
                    uint256 i = set.order[c];
                    row[i] = _scaled(set, x / S, i);
                }
            } else if (index < 2 * size) {
                int256 x = _draw(set, index - size, 0);
                for (uint256 i; i < n; ++i) {
                    row[i] = _scaled(set, x, i);
                }
            } else if (index < 3 * size) {
                for (uint256 i; i < n; ++i) {
                    row[i] = _scaled(set, _draw(set, index - 2 * size, set.rank[i]), i);
                }
            } else if (index < 3 * size + 4) {
                uint256 shock = index - 3 * size;
                for (uint256 i; i < n; ++i) {
                    int256 magnitude = shock < 2 ? set.moves[i] : int256(uint256(p.gaps[i]));
                    row[i] = _clamp(shock % 2 == 0 ? -magnitude : magnitude);
                }
            } else {
                uint256 alone = index - 3 * size - 4;
                uint256 i = set.order[alone / 2];
                for (uint256 k; k < n; ++k) {
                    row[k] = 0;
                }
                int256 gap = int256(uint256(p.gaps[i]));
                row[i] = _clamp(alone % 2 == 0 ? -gap : gap);
            }
        }
    }

    /// @notice Every scenario in the ascending order of the symbols, as big-endian int32.
    function encoded(Set memory set) internal pure returns (bytes memory out) {
        uint256 n = set.parameters.symbols.length;
        uint256 rows = count(set.size, n);
        uint256 length = rows * n * 4;
        out = new bytes(length + 32);
        assembly ("memory-safe") {
            mstore(out, length)
        }
        int256[] memory row = new int256[](n);
        int256[] memory draws = new int256[](n);
        uint256 offset;
        for (uint256 index; index < rows; ++index) {
            fill(set, index, row, draws);
            for (uint256 c; c < n; ++c) {
                uint256 word = uint256(uint32(int32(row[set.order[c]]))) << 224;
                assembly ("memory-safe") {
                    mstore(add(add(out, 32), offset), word)
                }
                offset += 4;
            }
        }
    }

    function _rho(Parameters memory p, uint256 a, uint256 b) private pure returns (int256) {
        if (a == b) return ONE;
        (uint256 i, uint256 j) = a < b ? (a, b) : (b, a);
        return int256(uint256(p.correlations[pair(p.symbols.length, i, j)]));
    }

    function _factor(Parameters memory p, uint256[] memory order) private pure returns (int256[] memory factor) {
        uint256 n = order.length;
        factor = new int256[](n * n);
        for (uint256 j; j < n; ++j) {
            int256 acc = S * S;
            for (uint256 k; k < j; ++k) {
                acc -= factor[j * n + k] * factor[j * n + k];
            }
            int256 pivot = int256(isqrt(acc > 0 ? uint256(acc) : 0));
            factor[j * n + j] = pivot >= PIVOT_FLOOR ? pivot : int256(0);
            for (uint256 i = j + 1; i < n; ++i) {
                acc = _rho(p, order[i], order[j]) * (S / ONE) * S;
                for (uint256 k; k < j; ++k) {
                    acc -= factor[i * n + k] * factor[j * n + k];
                }
                if (factor[j * n + j] > 0) {
                    int256 value = acc / factor[j * n + j];
                    factor[i * n + j] = value < -S ? -S : value > S ? S : value;
                }
            }
        }
    }

    function _moves(Parameters memory p, int256 root) private pure returns (int256[] memory moves) {
        uint256 n = p.symbols.length;
        int256 variance;
        for (uint256 a; a < n; ++a) {
            for (uint256 b; b < n; ++b) {
                variance += _weight(p, a) * _weight(p, b) * _vol(p, a) * _vol(p, b) * _rho(p, a, b);
            }
        }
        int256 sigma = int256(isqrt(uint256(variance * ONE)));
        moves = new int256[](n);
        for (uint256 i; i < n; ++i) {
            int256 covariance;
            for (uint256 k; k < n; ++k) {
                covariance += _weight(p, k) * _vol(p, i) * _vol(p, k) * _rho(p, i, k);
            }
            moves[i] = Z_MARKET * root * covariance / (1e12 * sigma);
        }
    }

    function _weight(Parameters memory p, uint256 k) private pure returns (int256) {
        return p.market == 0 || p.market == k + 1 ? int256(1) : int256(0);
    }

    function _vol(Parameters memory p, uint256 k) private pure returns (int256) {
        return int256(uint256(p.volatilities[k]));
    }

    function _means(uint256 size) private pure returns (int256[] memory means) {
        bytes memory table = size == 32 ? MEANS_32 : size == 64 ? MEANS_64 : size == 128 ? MEANS_128 : MEANS_256;
        uint256 half = size / 2;
        means = new int256[](size);
        for (uint256 j; j < half; ++j) {
            uint256 word;
            assembly ("memory-safe") {
                word := mload(add(add(table, 32), mul(j, 4)))
            }
            int256 value = int256(int32(uint32(word >> 224)));
            means[j] = value;
            means[size - 1 - j] = -value;
        }
    }

    function _draw(Set memory set, uint256 j, uint256 d) private pure returns (int256) {
        unchecked {
            return set.means[(j * set.generators[d] + set.shifts[d]) % set.size];
        }
    }

    function _scaled(Set memory set, int256 x, uint256 i) private pure returns (int256) {
        unchecked {
            return _clamp(x * int256(uint256(set.parameters.volatilities[i])) * set.root / (PPM * PPM));
        }
    }

    function _clamp(int256 r) private pure returns (int256) {
        return r < -PPM ? -PPM : r > int256(type(int32).max) ? int256(type(int32).max) : r;
    }
}
