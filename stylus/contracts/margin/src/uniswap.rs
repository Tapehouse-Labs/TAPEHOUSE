// SPDX-FileCopyrightText: Tapehouse contributors
// SPDX-License-Identifier: MIT OR Apache-2.0
//! Reads of Uniswap v3 pools and Chainlink feeds, and the pool arithmetic the liquidity add-on needs.

use stylus_sdk::alloy_primitives::{Address, U256, aliases::U160};
use stylus_sdk::prelude::*;

/// Decimals the ETH/USD feed reports.
pub const FEED_DECIMALS: u8 = 8;
/// The window of the pool's time-weighted price and liquidity, in seconds, as the band's feeds use.
pub const WINDOW: u32 = 1_800;
/// Maximum age of a Chainlink answer: its heartbeat plus 60 s.
pub const MAX_FEED_AGE: u64 = 86_400 + 60;
/// Value a pool pays out when the token's price falls 10%, per unit of liquidity and of √price:
/// 1 − √0.9, in billionths, rounded down.
pub const SELL_FACTOR: u64 = 51_316_701;
/// Value a pool takes in when the token's price rises 10%, per unit of liquidity and of √price:
/// √1.1 − 1, in billionths, rounded down.
pub const BUY_FACTOR: u64 = 48_808_848;
/// Lowest tick of a Uniswap v3 pool.
pub const MIN_TICK: i32 = -887_272;
/// Highest tick of a Uniswap v3 pool.
pub const MAX_TICK: i32 = 887_272;

sol_interface! {
    interface IUniswapV3Pool {
        function token0() external view returns (address);
        function token1() external view returns (address);
        function fee() external view returns (uint24);
        function slot0() external view returns (uint160, int24, uint16, uint16, uint16, uint8, bool);
        function observe(uint32[]) external view returns (int56[], uint160[]);
    }

    interface IERC20Metadata {
        function decimals() external view returns (uint8);
    }

    interface AggregatorV3Interface {
        function decimals() external view returns (uint8);
        function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80);
    }
}

/// `pool`'s two tokens and its fee in millionths. `None` when a call fails.
pub fn tokens(host: &impl Host, pool: Address) -> Option<(Address, Address, u32)> {
    let pool = IUniswapV3Pool::new(pool);
    Some((
        pool.token_0(host, Call::new()).ok()?,
        pool.token_1(host, Call::new()).ok()?,
        pool.fee(host, Call::new()).ok()?.to::<u32>(),
    ))
}

/// How many observations `pool` keeps. `None` when the call fails.
pub fn cardinality(host: &impl Host, pool: Address) -> Option<u16> {
    let (_, _, _, cardinality, _, _, _) =
        IUniswapV3Pool::new(pool).slot_0(host, Call::new()).ok()?;
    Some(cardinality)
}

/// `token`'s decimals. `None` when the call fails.
pub fn decimals(host: &impl Host, token: Address) -> Option<u8> {
    IERC20Metadata::new(token).decimals(host, Call::new()).ok()
}

/// Whether `feed` answers `decimals()` with [`FEED_DECIMALS`].
pub fn has_feed_decimals(host: &impl Host, feed: Address) -> bool {
    AggregatorV3Interface::new(feed).decimals(host, Call::new()) == Ok(FEED_DECIMALS)
}

/// Whether the asset's `stock` token is the pool's token0, and whether the other token is WETH:
/// `(stock_is_token0, quote_is_weth)`. `None` unless the pool trades `stock` against `usdg` or `weth`.
pub fn classify(
    token0: Address,
    token1: Address,
    stock: Address,
    usdg: Address,
    weth: Address,
) -> Option<(bool, bool)> {
    let quote = |token: Address| token != Address::ZERO && (token == usdg || token == weth);
    match (token0 == stock, token1 == stock) {
        (true, false) if quote(token1) => Some((true, token1 == weth)),
        (false, true) if quote(token0) => Some((false, token0 == weth)),
        _ => None,
    }
}

/// The latest answer of `feed`, or `None` when the call fails, the answer is not positive or it is older
/// than [`MAX_FEED_AGE`] at `now`.
pub fn answer(host: &impl Host, feed: Address, now: u64) -> Option<U256> {
    let (_, answer, _, updated_at, _) = AggregatorV3Interface::new(feed)
        .latest_round_data(host, Call::new())
        .ok()?;
    let updated_at = u64::try_from(updated_at).ok()?;
    (answer.is_positive() && updated_at <= now && now - updated_at <= MAX_FEED_AGE)
        .then(|| answer.into_raw())
}

/// The arithmetic-mean tick and harmonic-mean liquidity of `pool` over the last [`WINDOW`] seconds, as
/// Uniswap's `OracleLibrary.consult` computes them. `None` when the call fails, as it does without
/// [`WINDOW`] seconds of observations.
pub fn consult(host: &impl Host, pool: Address) -> Option<(i32, u128)> {
    let (ticks, seconds) = IUniswapV3Pool::new(pool)
        .observe(host, Call::new(), vec![WINDOW, 0])
        .ok()?;
    mean(
        i64::try_from(*ticks.first()?).ok()?,
        i64::try_from(*ticks.get(1)?).ok()?,
        *seconds.first()?,
        *seconds.get(1)?,
    )
}

/// `OracleLibrary.consult`'s arithmetic on two observations [`WINDOW`] seconds apart, with its cumulatives
/// wrapping as the pool's `int56` and `uint160` do. `None` if the mean tick lies outside [`MIN_TICK`,
/// `MAX_TICK`], which no Uniswap v3 pool reports.
pub fn mean(
    tick_then: i64,
    tick_now: i64,
    seconds_then: U160,
    seconds_now: U160,
) -> Option<(i32, u128)> {
    let window = i64::from(WINDOW);
    let delta = tick_now.wrapping_sub(tick_then) << 8 >> 8;
    let mut tick = delta / window;
    if delta < 0 && delta % window != 0 {
        tick -= 1;
    }
    if !(i64::from(MIN_TICK)..=i64::from(MAX_TICK)).contains(&tick) {
        return None;
    }
    let per_liquidity: U256 = U256::from(seconds_now.wrapping_sub(seconds_then)) << 32;
    if per_liquidity.is_zero() {
        return None;
    }
    let liquidity = U256::from(WINDOW) * U256::from(U160::MAX) / per_liquidity;
    Some((i32::try_from(tick).ok()?, u128::try_from(liquidity).ok()?))
}

// SPDX-SnippetBegin
// SPDX-SnippetCopyrightText: 2023 Universal Navigation Inc.
// SPDX-License-Identifier: MIT
/// √(1.0001^`tick`) as a Q64.96, exactly as Uniswap v4-core's `TickMath.getSqrtPriceAtTick`, whose factors
/// this uses. `tick` lies within [`MIN_TICK`, `MAX_TICK`].
pub fn sqrt_ratio_at_tick(tick: i32) -> U256 {
    const FACTORS: [u128; 19] = [
        0xfff97272373d413259a46990580e213a,
        0xfff2e50f5f656932ef12357cf3c7fdcc,
        0xffe5caca7e10e4e61c3624eaa0941cd0,
        0xffcb9843d60f6159c9db58835c926644,
        0xff973b41fa98c081472e6896dfb254c0,
        0xff2ea16466c96a3843ec78b326b52861,
        0xfe5dee046a99a2a811c461f1969c3053,
        0xfcbe86c7900a88aedcffc83b479aa3a4,
        0xf987a7253ac413176f2b074cf7815e54,
        0xf3392b0822b70005940c7a398e4b70f3,
        0xe7159475a2c29b7443b29c7fa6e889d9,
        0xd097f3bdfd2022b8845ad8f792aa5825,
        0xa9f746462d870fdf8a65dc1f90e061e5,
        0x70d869a156d2a1b890bb3df62baf32f7,
        0x31be135f97d08fd981231505542fcfa6,
        0x9aa508b5b7a84e1c677de54f3e99bc9,
        0x5d6af8dedb81196699c329225ee604,
        0x2216e584f5fa1ea926041bedfe98,
        0x48a170391f7dc42444e8fa2,
    ];
    let abs = tick.unsigned_abs();
    let mut ratio = if abs & 1 != 0 {
        U256::from(0xfffcb933bd6fad37aa2d162d1a594001_u128)
    } else {
        U256::from(1) << 128
    };
    for (bit, &factor) in FACTORS.iter().enumerate() {
        if abs & (2 << bit) != 0 {
            ratio = (ratio * U256::from(factor)) >> 128;
        }
    }
    if tick > 0 {
        ratio = U256::MAX / ratio;
    }
    let rounded = if (ratio & U256::from(u32::MAX)).is_zero() {
        0
    } else {
        1
    };
    (ratio >> 32) + U256::from(rounded)
}
// SPDX-SnippetEnd

/// From a pool's mean `tick` and `liquidity`: the asset's price in USD with 8 decimals, and the USD with 18
/// decimals the pool pays out for a 10% fall and takes in for a 10% rise at that liquidity. `usd` is the
/// quote token's price in USD with 8 decimals. `None` if a value does not fit.
pub fn terms(
    stock_is_token0: bool,
    stock_decimals: u8,
    quote_decimals: u8,
    tick: i32,
    liquidity: u128,
    usd: U256,
) -> Option<(U256, U256, U256)> {
    let sqrt = sqrt_ratio_at_tick(if stock_is_token0 { tick } else { -tick });
    let stock_unit = U256::from(10).pow(U256::from(stock_decimals));
    let quote_unit = U256::from(10).pow(U256::from(quote_decimals));
    let per_stock = mul_shr(mul_shr(sqrt, sqrt, 64), stock_unit, 128);
    let price = per_stock.checked_mul(usd)? / quote_unit;
    let moved = mul_shr(U256::from(liquidity), sqrt, 96);
    let wad = U256::from(10).pow(U256::from(18));
    let value = |factor: u64| -> Option<U256> {
        let raw = moved.checked_mul(U256::from(factor))? / U256::from(1_000_000_000);
        Some(raw.checked_mul(usd)?.checked_mul(wad / quote_unit)? / U256::from(100_000_000))
    };
    Some((price, value(SELL_FACTOR)?, value(BUY_FACTOR)?))
}

/// `a` · `b` / 2^`shift`, rounded down, without overflowing on the way.
pub fn mul_shr(a: U256, b: U256, shift: usize) -> U256 {
    U256::from(a.widening_mul::<256, 4, 512, 8>(b) >> shift)
}

#[cfg(test)]
mod tests {
    use super::*;
    use stylus_sdk::alloy_primitives::{I256, address, hex};
    use stylus_sdk::alloy_sol_types::SolValue;
    use stylus_sdk::testing::*;

    const USDG: Address = address!("0x0000000000000000000000000000000000000d01");
    const WETH: Address = address!("0x0000000000000000000000000000000000000d02");
    const STOCK: Address = address!("0x0000000000000000000000000000000000000d04");

    #[test]
    fn a_pool_is_the_asset_against_usdg_or_weth() {
        let other = address!("0x0000000000000000000000000000000000000d05");
        assert_eq!(
            classify(USDG, STOCK, STOCK, USDG, WETH),
            Some((false, false))
        );
        assert_eq!(classify(STOCK, WETH, STOCK, USDG, WETH), Some((true, true)));
        assert_eq!(
            classify(WETH, STOCK, STOCK, USDG, WETH),
            Some((false, true))
        );
        assert_eq!(classify(USDG, other, STOCK, USDG, WETH), None);
        assert_eq!(classify(USDG, WETH, STOCK, USDG, WETH), None);
        assert_eq!(classify(STOCK, other, STOCK, USDG, WETH), None);
        assert_eq!(
            classify(STOCK, Address::ZERO, STOCK, USDG, Address::ZERO),
            None
        );
    }

    const POOL: Address = address!("0xd4EB21209C4D6093f80B5b84f5C45cc093EA14a3");
    const FEED: Address = address!("0x0000000000000000000000000000000000000d03");

    #[test]
    fn the_square_root_ratio_matches_uniswap_at_its_ends_and_at_zero() {
        assert_eq!(sqrt_ratio_at_tick(0), U256::from(1) << 96);
        assert_eq!(sqrt_ratio_at_tick(MIN_TICK), U256::from(4295128739_u64));
        assert_eq!(
            sqrt_ratio_at_tick(MAX_TICK),
            U256::from_str_radix("1461446703485210103287273052203988822378723970342", 10).unwrap()
        );
    }

    #[test]
    fn the_square_root_ratio_matches_uniswap_on_its_test_vectors() {
        for (tick, expected) in [
            (-887_271, "4295343490"),
            (887_271, "1461373636630004318706518188784493106690254656249"),
            (-50, "79030349367926598376800521322"),
            (50, "79426470787362580746886972461"),
            (-100, "78833030112140176575862854579"),
            (100, "79625275426524748796330556128"),
            (-250, "78244023372248365697264290337"),
            (250, "80224679980005306637834519095"),
            (-500, "77272108795590369356373805297"),
            (500, "81233731461783161732293370115"),
            (-1_000, "75364347830767020784054125655"),
            (1_000, "83290069058676223003182343270"),
            (-2_500, "69919044979842180277688105136"),
            (2_500, "89776708723587163891445672585"),
            (-3_000, "68192822843687888778582228483"),
            (3_000, "92049301871182272007977902845"),
            (-4_000, "64867181785621769311890333195"),
            (4_000, "96768528593268422080558758223"),
            (-5_000, "61703726247759831737814779831"),
            (5_000, "101729702841318637793976746270"),
            (-50_000, "6504256538020985011912221507"),
            (50_000, "965075977353221155028623082916"),
            (-150_000, "43836292794701720435367485"),
            (150_000, "143194173941309278083010301478497"),
            (-250_000, "295440463448801648376846"),
            (250_000, "21246587762933397357449903968194344"),
            (-500_000, "1101692437043807371"),
            (500_000, "5697689776495288729098254600827762987878"),
            (-738_203, "7409801140451"),
            (738_203, "847134979253254120489401328389043031315994541"),
        ] {
            assert_eq!(
                sqrt_ratio_at_tick(tick),
                U256::from_str_radix(expected, 10).unwrap(),
                "{tick}"
            );
        }
    }

    #[test]
    fn the_mean_tick_rounds_toward_minus_infinity() {
        let tick = |delta: i64| mean(0, delta, U160::ZERO, U160::MAX).unwrap().0;
        assert_eq!(tick(1_800 * 7), 7);
        assert_eq!(tick(-1_800 * 7), -7);
        assert_eq!(tick(-1_800 * 7 - 1), -8);
        assert_eq!(tick(1_800 * 7 + 1), 7);
    }

    #[test]
    fn the_harmonic_mean_liquidity_inverts_seconds_per_liquidity() {
        let liquidity: u128 = 11_245_526_858_841_909_681;
        let per_second = (U256::from(1) << 128) / U256::from(liquidity);
        let delta = U160::from(per_second * U256::from(WINDOW));
        let (_, mean_liquidity) = mean(0, 0, U160::ZERO, delta).unwrap();
        assert!(mean_liquidity.abs_diff(liquidity) * 1_000_000_000 < liquidity);
        assert_eq!(mean(0, 0, delta, delta), None);
    }

    #[test]
    fn cumulatives_wrap_as_the_pool_s_do() {
        let per_second = U160::from(1) << 100;
        let span = per_second * U160::from(WINDOW);
        let then = U160::MAX - U160::from(5);
        let expected = mean(0, 1800 * 42, U160::ZERO, span);
        assert_eq!(mean(0, 1800 * 42, then, then.wrapping_add(span)), expected);
        let top = (1_i64 << 55) - 1;
        let wrapped = top + 1800 * 42 - (1_i64 << 56);
        assert_eq!(mean(top, wrapped, U160::ZERO, span), expected);
    }

    #[test]
    fn a_mean_tick_outside_uniswap_s_range_is_no_answer() {
        let span = (U160::from(1) << 100) * U160::from(WINDOW);
        let window = i64::from(WINDOW);
        assert!(mean(0, i64::from(MAX_TICK) * window, U160::ZERO, span).is_some());
        assert!(mean(0, i64::from(MIN_TICK) * window, U160::ZERO, span).is_some());
        assert_eq!(
            mean(0, (i64::from(MAX_TICK) + 1) * window, U160::ZERO, span),
            None
        );
        assert_eq!(
            mean(0, (i64::from(MIN_TICK) - 1) * window, U160::ZERO, span),
            None
        );
    }

    #[test]
    fn consult_reads_the_last_thirty_minutes_and_nothing_without_them() {
        let vm = TestVM::default();
        let call = [
            hex!("883bdbfd").to_vec(),
            (vec![WINDOW, 0u32],).abi_encode_params(),
        ]
        .concat();
        let observed = (
            vec![
                I256::try_from(399_581_400_000_i64).unwrap(),
                I256::try_from(399_581_400_000_i64 + 1_800 * 221_989).unwrap(),
            ],
            vec![
                U256::from(1_000_000_u64),
                U256::from(1_000_000_u64)
                    + (U256::from(WINDOW) << 128) / U256::from(11_245_526_858_841_909_681_u128),
            ],
        )
            .abi_encode_params();
        vm.mock_static_call(POOL, call.clone(), Ok(observed));
        let (tick, liquidity) = consult(&vm, POOL).unwrap();
        assert_eq!(tick, 221_989);
        assert!(liquidity.abs_diff(11_245_526_858_841_909_681) * 1_000_000_000 < liquidity);
        vm.mock_static_call(POOL, call, Err(b"OLD".to_vec()));
        assert_eq!(consult(&vm, POOL), None);
    }

    #[test]
    fn an_ether_price_counts_while_it_is_fresh_and_positive() {
        let vm = TestVM::default();
        let now = 1_790_000_000_u64;
        let round = |answer: i64, at: u64| {
            Ok((
                U256::from(1),
                I256::try_from(answer).unwrap(),
                U256::from(at),
                U256::from(at),
                U256::from(1),
            )
                .abi_encode_params())
        };
        let call = hex!("feaf968c").to_vec();
        vm.mock_static_call(
            FEED,
            call.clone(),
            round(268_330_550_000, now - MAX_FEED_AGE),
        );
        assert_eq!(
            answer(&vm, FEED, now),
            Some(U256::from(268_330_550_000_u64))
        );
        vm.mock_static_call(
            FEED,
            call.clone(),
            round(268_330_550_000, now - MAX_FEED_AGE - 1),
        );
        assert_eq!(answer(&vm, FEED, now), None);
        vm.mock_static_call(FEED, call.clone(), round(0, now));
        assert_eq!(answer(&vm, FEED, now), None);
        vm.mock_static_call(FEED, call, Err(vec![]));
        assert_eq!(answer(&vm, FEED, now), None);
        vm.mock_static_call(
            FEED,
            hex!("313ce567").to_vec(),
            Ok(U256::from(8).abi_encode()),
        );
        assert!(has_feed_decimals(&vm, FEED));
    }

    #[test]
    fn a_pool_prices_its_asset_from_its_mean_tick() {
        let (price, selling, buying) = terms(
            false,
            18,
            6,
            221_989,
            11_245_526_858_841_909_681,
            U256::from(100_000_000),
        )
        .unwrap();
        assert_eq!(price, U256::from(22_888_758_400_u64));
        assert!(selling > buying);
        let (price, _, _) = terms(
            false,
            18,
            18,
            12_513,
            16_029_297_629_534_329_325_587,
            U256::from(268_330_550_000_u64),
        )
        .unwrap();
        assert_eq!(price, U256::from(76_782_916_718_u64));
    }
}
