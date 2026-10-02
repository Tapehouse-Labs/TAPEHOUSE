// SPDX-License-Identifier: MIT OR Apache-2.0
use std::error::Error as _;

use alloy::primitives::{Address, B256, Bytes, U64, U256, address, aliases::U24};
use alloy::providers::{Provider, ProviderBuilder};
use alloy::rpc::json_rpc::ErrorPayload;
use alloy::sol_types::{Revert as SolidityRevert, SolCall, SolError, SolEvent, SolValue};
use alloy::transports::mock::Asserter;
use tapehouse_sdk::bindings::{
    Band, BandFeed, Basket, GapCover, MarginAccounts, MorphoBandOracle, ShortPositions, StockToken,
};
use tapehouse_sdk::{
    CROSS, Deployments, Error, Layer, NoPrice, OraclePrice, PackageSource, PricingGap, Sales,
    SeriesStatus, Tapehouse, decode_revert_data, payout, to_bytes32,
};

const ALICE: Address = address!("0x70997970C51812dc3A010C7d01b50e0d17dc79C8");
const BOB: Address = address!("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC");

fn registry(chain_id: u64) -> Deployments {
    Deployments::load(format!(
        "{}/../../deployments/{chain_id}.json",
        env!("CARGO_MANIFEST_DIR")
    ))
    .unwrap()
}

fn b32(name: &str) -> B256 {
    to_bytes32(name).unwrap()
}

fn offline(deployments: Deployments) -> Tapehouse<impl Provider + Clone> {
    Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(Asserter::new()),
        deployments,
    )
}

#[test]
fn chainlink_feeds_price_the_token_on_robinhood_chain_and_the_share_on_arbitrum_one() {
    let (robinhood, arbitrum) = (registry(4663), registry(42161));
    assert_eq!(
        robinhood.token_price_feed("NVDA_USD").unwrap().0,
        address!("0x379EC4f7C378F34a1B47E4F3cbeBCbAC3E8E9F15")
    );
    assert_eq!(
        arbitrum.share_price_feed("NVDA_USD").unwrap().0,
        address!("0x4881A4418b5F2460B21d6F08CD5aA0678a7f262F")
    );
    assert_eq!(
        robinhood
            .share_price_feed("NVDA_USD")
            .unwrap_err()
            .to_string(),
        ".chainlink.NVDA_USD prices the Stock Token on chain 4663"
    );
    assert_eq!(
        arbitrum
            .token_price_feed("NVDA_USD")
            .unwrap_err()
            .to_string(),
        ".chainlink.NVDA_USD prices the share on chain 42161"
    );
}

#[test]
fn every_group_is_read_and_a_missing_one_is_empty() {
    let (testnet, robinhood) = (registry(46630), registry(4663));
    assert_eq!(testnet.chain_id, 46630);
    let raw: serde_json::Value =
        serde_json::from_str(include_str!("../../../deployments/46630.json")).unwrap();
    for group in ["tapehouse", "bandFeeds"] {
        for (name, value) in raw[group].as_object().unwrap() {
            if let Some(text) = value.as_str() {
                assert!(
                    Address::parse_checksummed(text, None).is_ok(),
                    ".{group}.{name} is not checksummed: {text}"
                );
            }
        }
    }
    assert_eq!(
        robinhood.uniswap_v3["QuoterV2"],
        address!("0x33e885eD0Ec9bF04EcfB19341582aADCb4c8A9E7")
    );
}

#[test]
fn a_mixed_case_address_must_carry_its_checksum() {
    let lower = Deployments::parse(&format!(
        r#"{{"chainId":412346,"tapehouse":{{"Band":"0xa70118d3324d90532e7d2854627b13cace305641","StockLending":{{"SPY":"{BOB}"}}}}}}"#
    ))
    .unwrap();
    assert_eq!(
        lower.tapehouse["Band"],
        address!("0xa70118d3324D90532E7D2854627b13CacE305641")
    );
    assert_eq!(lower.stock_lending["SPY"], BOB);
    assert_eq!(lower.tapehouse.len(), 1);
    for (json, message) in [
        (
            r#"{"chainId":1,"tokens":{"USDG":"0xa70118D3324D90532E7D2854627b13CacE305641"}}"#,
            ".tokens.USDG is not an address",
        ),
        (r#"{"tokens":{}}"#, "the registry has no chainId"),
        (
            r#"{"chainId":1,"tapehouse":{"Band":1}}"#,
            ".tapehouse.Band is not an address",
        ),
    ] {
        assert_eq!(
            Deployments::parse(json).unwrap_err().to_string(),
            message,
            "{json}"
        );
    }
}

#[test]
fn the_baskets_are_read_from_tapehouse_baskets() {
    let deployments = Deployments::parse(&format!(
        r#"{{"chainId":412346,"tapehouse":{{"MarginAccounts":"{ALICE}","Baskets":{{"PAIR":"{}"}}}}}}"#,
        BOB.to_string().to_lowercase()
    ))
    .unwrap();
    assert_eq!(deployments.baskets["PAIR"], BOB);
    assert_eq!(deployments.tapehouse.len(), 1);
    assert!(registry(46630).baskets.is_empty());
    assert_eq!(
        Deployments::parse(r#"{"chainId":1,"tapehouse":{"Baskets":{"PAIR":1}}}"#)
            .unwrap_err()
            .to_string(),
        ".tapehouse.Baskets.PAIR is not an address"
    );
}

#[test]
fn a_basket_mints_and_redeems_in_kind_and_the_accounts_unwrap_one() {
    let mut deployments = registry(46630);
    deployments.baskets.insert("PAIR".into(), BOB);
    let tapehouse = offline(deployments.clone());
    let pair = tapehouse.basket("PAIR").unwrap();
    assert_eq!(*pair.address(), BOB);
    let mint = pair.mint(
        U256::from(10u64.pow(18)),
        ALICE,
        vec![U256::from(3), U256::from(4)],
    );
    let decoded = Basket::mintCall::abi_decode(mint.calldata()).unwrap();
    assert_eq!(
        (decoded.shares, decoded.receiver, decoded.maxAssets),
        (
            U256::from(10u64.pow(18)),
            ALICE,
            vec![U256::from(3), U256::from(4)]
        )
    );
    let redeem = pair.redeem(U256::from(5), BOB, ALICE);
    let decoded = Basket::redeemCall::abi_decode(redeem.calldata()).unwrap();
    assert_eq!((decoded.receiver, decoded.owner), (BOB, ALICE));
    let unwrap = tapehouse.unwrap(ALICE, "PAIR", U256::from(5)).unwrap();
    assert_eq!(unwrap.calldata()[..4], MarginAccounts::unwrapCall::SELECTOR);
    let decoded = MarginAccounts::unwrapCall::abi_decode(unwrap.calldata()).unwrap();
    assert_eq!(
        (decoded.account, decoded.basket, decoded.shares),
        (ALICE, BOB, U256::from(5))
    );
    assert_eq!(
        tapehouse.basket("NONE").unwrap_err().to_string(),
        "the registry has no .tapehouse.Baskets.NONE"
    );
}

#[test]
fn a_sale_and_a_cover_carry_the_limits_a_quote_gives() {
    let mut deployments = registry(46630);
    deployments.tapehouse.insert("ShortPositions".into(), BOB);
    let shorts = offline(deployments).shorts().unwrap();
    let sale = shorts.sell(
        b32("SPY"),
        U256::from(10u64.pow(18)),
        U256::from(770),
        ALICE,
    );
    assert_eq!(sale.calldata()[..4], ShortPositions::sellCall::SELECTOR);
    let decoded = ShortPositions::sellCall::abi_decode(sale.calldata()).unwrap();
    assert_eq!(
        (
            decoded.symbol,
            decoded.amount,
            decoded.minProceeds,
            decoded.account
        ),
        (
            b32("SPY"),
            U256::from(10u64.pow(18)),
            U256::from(770),
            ALICE
        )
    );
    assert_eq!(*shorts.address(), BOB);
}

#[test]
fn an_authorization_lets_one_address_act_for_the_sender() {
    let deployments = registry(46630);
    let accounts = offline(deployments.clone()).accounts().unwrap();
    let call = accounts.setAuthorization(BOB, true);
    let decoded = MarginAccounts::setAuthorizationCall::abi_decode(call.calldata()).unwrap();
    assert_eq!((decoded.authorized, decoded.allowed), (BOB, true));
    assert_eq!(*accounts.address(), deployments.tapehouse["MarginAccounts"]);
    let repay = accounts.repay(CROSS, U256::from(5), ALICE);
    let decoded = MarginAccounts::repayCall::abi_decode(repay.calldata()).unwrap();
    assert_eq!(
        (decoded.position, decoded.assets, decoded.account),
        (CROSS, U256::from(5), ALICE)
    );
}

#[test]
fn a_contract_missing_from_the_registry_is_named() {
    let empty = Deployments::parse(r#"{"chainId":1}"#).unwrap();
    assert_eq!(
        offline(empty.clone()).shorts().unwrap_err().to_string(),
        "the registry has no .tapehouse.ShortPositions"
    );
    assert_eq!(
        offline(empty).band_feed("NVDA").unwrap_err().to_string(),
        "the registry has no .bandFeeds.NVDA"
    );
}

struct FileSource(std::sync::Mutex<Vec<Vec<String>>>);

impl PackageSource for FileSource {
    async fn payload(
        &self,
        feed_ids: &[&str],
    ) -> Result<Bytes, Box<dyn std::error::Error + Send + Sync>> {
        self.0
            .lock()
            .unwrap()
            .push(feed_ids.iter().map(|id| id.to_string()).collect());
        Ok(packages())
    }
}

fn packages() -> Bytes {
    std::fs::read_to_string(concat!(
        env!("CARGO_MANIFEST_DIR"),
        "/../../stylus/contracts/band/testdata/nvda-24_7.hex"
    ))
    .unwrap()
    .trim()
    .parse()
    .unwrap()
}

#[tokio::test]
async fn signed_redstone_packages_come_from_the_source_the_integrator_plugs_in() {
    let deployments = registry(46630);
    let source = FileSource(Default::default());
    let call = offline(deployments.clone())
        .write_prices(&source, &["NVDA---24_7"])
        .await
        .unwrap();
    assert_eq!(*source.0.lock().unwrap(), vec![vec!["NVDA---24_7"]]);
    let decoded = Band::writePricesCall::abi_decode(call.calldata()).unwrap();
    assert_eq!(decoded.feedIds, vec![b32("NVDA---24_7")]);
    assert_eq!(decoded.payload, packages());
}

#[test]
fn the_reverts_of_tapehouse_the_band_the_issuer_and_the_router_are_decoded() {
    let cases = [
        (
            ShortPositions::Unauthorized {
                caller: BOB,
                account: ALICE,
            }
            .abi_encode(),
            format!("Unauthorized({BOB}, {ALICE})"),
        ),
        (
            BandFeed::NotSealWindow {
                session: 2,
                reopenMs: 1790200000000,
            }
            .abi_encode(),
            "NotSealWindow(2, 1790200000000)".into(),
        ),
        (StockToken::IsPaused {}.abi_encode(), "IsPaused()".into()),
        (
            Band::TimestampIsTooOld {
                receivedTimestampSeconds: U256::from(1790119750u64),
                blockTimestamp: U256::from(1790300000u64),
            }
            .abi_encode(),
            "TimestampIsTooOld(1790119750, 1790300000)".into(),
        ),
        (
            MarginAccounts::DebtCapExceeded {
                debt: U256::from(2),
                cap: U256::from(1),
            }
            .abi_encode(),
            "DebtCapExceeded(2, 1)".into(),
        ),
        (
            MarginAccounts::BasketFrozen { basket: BOB }.abi_encode(),
            format!("BasketFrozen({BOB})"),
        ),
        (
            Basket::PastTarget { symbol: b32("SPY") }.abi_encode(),
            format!("PastTarget({})", b32("SPY")),
        ),
        (
            SolidityRevert::from("Too little received").abi_encode(),
            "Error(Too little received)".into(),
        ),
    ];
    for (data, expected) in cases {
        assert_eq!(decode_revert_data(&data).unwrap().to_string(), expected);
    }
    assert_eq!(decode_revert_data(&[0xde, 0xad, 0xbe, 0xef]), None);
}

#[tokio::test]
async fn the_revert_behind_a_failed_call_is_decoded() {
    let asserter = Asserter::new();
    let data = MarginAccounts::Unauthorized {
        caller: BOB,
        account: ALICE,
    }
    .abi_encode();
    asserter.push_failure(ErrorPayload {
        code: 3,
        message: "execution reverted".into(),
        data: Some(serde_json::value::to_raw_value(&Bytes::from(data)).unwrap()),
    });
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter),
        registry(46630),
    );
    let accounts = tapehouse.accounts().unwrap();
    let error = Error::from(
        accounts
            .withdraw(CROSS, Address::ZERO, U256::from(1), ALICE, BOB)
            .call()
            .await
            .unwrap_err(),
    );
    let Error::Revert(revert) = error else {
        panic!("not a revert: {error}");
    };
    assert_eq!(revert.to_string(), format!("Unauthorized({BOB}, {ALICE})"));
}

#[tokio::test]
async fn a_sale_takes_its_quote_less_the_slippage_and_a_cover_pays_its_quote_plus_it() {
    let asserter = Asserter::new();
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("ShortPositions".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    let quote = (U256::from(77_232_802u64), U256::ZERO, 0u32, U256::ZERO).abi_encode_params();
    for _ in 0..2 {
        asserter.push_success(&Bytes::from(U24::from(500).abi_encode()));
        asserter.push_success(&Bytes::from(quote.clone()));
    }
    let sale = tapehouse
        .quote_sale("SPY", U256::from(10u64.pow(17)), 50)
        .await
        .unwrap();
    assert_eq!(
        (sale.proceeds, sale.min_proceeds),
        (U256::from(77_232_802u64), U256::from(76_846_637u64))
    );
    let cover = tapehouse
        .quote_cover("SPY", U256::from(10u64.pow(17)), 50)
        .await
        .unwrap();
    assert_eq!(
        (cover.cost, cover.max_cost),
        (U256::from(77_232_802u64), U256::from(77_618_967u64))
    );
}

#[tokio::test]
async fn a_slippage_outside_0_to_10000_basis_points_is_refused() {
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("ShortPositions".into(), BOB);
    let tapehouse = offline(deployments);
    let error = tapehouse
        .quote_sale("SPY", U256::from(1), 10_001)
        .await
        .unwrap_err();
    assert!(matches!(error, Error::Argument(_)));
    assert_eq!(
        error.to_string(),
        "slippage_bps 10001 is outside 0 to 10000"
    );
    assert!(
        tapehouse
            .quote_cover("SPY", U256::from(1), 10_001)
            .await
            .is_err()
    );
}

#[tokio::test]
async fn a_quote_that_fails_is_an_error_never_a_zero_limit() {
    let asserter = Asserter::new();
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("ShortPositions".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    asserter.push_success(&Bytes::from(U24::ZERO.abi_encode()));
    asserter.push_failure(ErrorPayload {
        code: 3,
        message: "execution reverted".into(),
        data: Some(serde_json::value::to_raw_value(&Bytes::new()).unwrap()),
    });
    let error = tapehouse
        .quote_sale("SPY", U256::from(1), 50)
        .await
        .unwrap_err();
    assert!(matches!(error, Error::Contract(_)), "{error}");
    assert!(error.source().is_some());
}

#[test]
fn a_name_longer_than_32_bytes_is_refused() {
    let long = "A".repeat(33);
    let error = to_bytes32(&long).unwrap_err();
    assert!(matches!(error, Error::Argument(_)));
    assert_eq!(
        error.to_string(),
        format!("\"{long}\" is longer than 32 bytes")
    );
    assert!(error.source().is_none());
    assert_eq!(b32("SPY"), B256::right_padding_from(b"SPY"));
}

#[tokio::test]
async fn a_repayment_reads_the_debt_and_the_premium_at_one_block() {
    let asserter = Asserter::new();
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        registry(46630),
    );
    asserter.push_success(&U64::from(7));
    asserter.push_success(&Bytes::from(U256::from(100).abi_encode()));
    asserter.push_success(&Bytes::from(U256::from(3).abi_encode()));
    let repayment = tapehouse.repayment(ALICE, CROSS).await.unwrap();
    assert_eq!(
        (repayment.debt, repayment.premium, repayment.assets),
        (U256::from(100), U256::from(3), U256::from(103))
    );
    assert!(asserter.read_q().is_empty());
}

#[tokio::test]
async fn a_positions_collateral_leverage_and_liquidation_price_are_read() {
    let asserter = Asserter::new();
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        registry(46630),
    );
    for answer in [2 * 10u64.pow(18), 25_000, 69_412_000_000] {
        asserter.push_success(&Bytes::from(U256::from(answer).abi_encode()));
    }
    assert_eq!(
        tapehouse.collateral(ALICE, CROSS, BOB).await.unwrap(),
        U256::from(2 * 10u64.pow(18))
    );
    assert_eq!(
        tapehouse.leverage(ALICE, CROSS).await.unwrap(),
        U256::from(25_000)
    );
    assert_eq!(
        tapehouse
            .liquidation_price(ALICE, CROSS, "SPY", U256::from(5))
            .await
            .unwrap(),
        U256::from(69_412_000_000u64)
    );
}

const BAND: Address = address!("0xa70118d3324D90532E7D2854627b13CacE305641");

#[test]
fn morpho_blue_and_the_bands_morpho_oracles_are_read() {
    let robinhood = registry(4663);
    assert_eq!(
        robinhood.morpho["Blue"],
        address!("0x9D53d5E3bd5E8d4Cbfa6DB1ca238AEA02E651010")
    );
    assert_eq!(
        robinhood.morpho["AdaptiveCurveIrm"],
        address!("0x2BD3d5965B26B51814AC95127B2b80dD6CcC0fa1")
    );
    let deployments = Deployments::parse(&format!(
        r#"{{"chainId":412346,"morphoOracles":{{"NVDA":"{BOB}"}}}}"#
    ))
    .unwrap();
    assert_eq!(deployments.morpho_oracles["NVDA"], BOB);
    assert_eq!(
        offline(registry(46630))
            .morpho_oracle("NVDA")
            .unwrap_err()
            .to_string(),
        "the registry has no .morphoOracles.NVDA"
    );
}

#[test]
fn morpho_blues_markets_are_read_from_markets_as_32_byte_ids() {
    let id = "0x3a85e619751152991742810df6ec69ce473daef99e28a64ab2340d7b7ccfee49";
    let deployments = Deployments::parse(&format!(
        r#"{{"chainId":412346,"morpho":{{"Blue":"{ALICE}","Markets":{{"NVDA_USDG":"0x{}"}}}}}}"#,
        id[2..].to_uppercase()
    ))
    .unwrap();
    assert_eq!(deployments.morpho.len(), 1);
    assert_eq!(deployments.morpho["Blue"], ALICE);
    assert_eq!(
        deployments.morpho_markets["NVDA_USDG"],
        id.parse::<B256>().unwrap()
    );
    assert!(registry(4663).morpho_markets.is_empty());
    for wrong in [
        format!(r#""{ALICE}""#),
        format!(r#""{id}00""#),
        format!(r#""{}""#, &id[..65]),
        format!(r#""{}""#, &id[2..]),
        "7".into(),
    ] {
        assert_eq!(
            Deployments::parse(&format!(
                r#"{{"chainId":1,"morpho":{{"Markets":{{"NVDA_USDG":{wrong}}}}}}}"#
            ))
            .unwrap_err()
            .to_string(),
            ".morpho.Markets.NVDA_USDG is not a 32-byte id",
            "{wrong}"
        );
    }
}

fn oracles(asserter: &Asserter) -> Tapehouse<impl Provider + Clone> {
    let mut deployments = registry(46630);
    deployments.morpho_oracles.insert("NVDA".into(), BOB);
    Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    )
}

fn revert(asserter: &Asserter, data: Vec<u8>) {
    asserter.push_failure(ErrorPayload {
        code: 3,
        message: "execution reverted".into(),
        data: Some(serde_json::value::to_raw_value(&Bytes::from(data)).unwrap()),
    });
}

fn no_answer(asserter: &Asserter, halt: (bool, u64, u64, bool), step: u8) {
    asserter.push_success(&U64::from(7));
    revert(
        asserter,
        MorphoBandOracle::NoAnswer {
            symbol: b32("NVDA"),
        }
        .abi_encode(),
    );
    asserter.push_success(&Bytes::from(halt.abi_encode_params()));
    asserter.push_success(&Bytes::from(BAND.abi_encode()));
    asserter.push_success(&Bytes::from(b32("NVDA").abi_encode()));
    asserter.push_success(&Bytes::from(
        (U256::from(step), U256::ZERO, U256::ZERO, U256::ZERO).abi_encode_params(),
    ));
}

#[tokio::test]
async fn an_oracle_answers_its_price() {
    let asserter = Asserter::new();
    let price = U256::from(765_517_754_420u64) * U256::from(10u64.pow(15));
    asserter.push_success(&U64::from(7));
    asserter.push_success(&Bytes::from(price.abi_encode()));
    assert_eq!(
        oracles(&asserter).oracle_price("NVDA").await.unwrap(),
        OraclePrice::Price(price)
    );
    assert!(asserter.read_q().is_empty());
}

#[tokio::test]
async fn no_answer_with_no_halt_pause_or_step_is_a_stale_band_never_a_price_of_zero() {
    let asserter = Asserter::new();
    no_answer(&asserter, (false, 0, 0, false), 0);
    let OraclePrice::NoPrice { reason, revert } =
        oracles(&asserter).oracle_price("NVDA").await.unwrap()
    else {
        panic!("a price without an answer");
    };
    assert_eq!(reason, NoPrice::Stale);
    assert_eq!(revert.to_string(), format!("NoAnswer({})", b32("NVDA")));
    assert!(asserter.read_q().is_empty());
}

#[tokio::test]
async fn no_answer_under_a_signed_halt_the_pause_or_an_unconfirmed_step_is_a_halt() {
    for (halt, step) in [
        ((true, 1_790_200_000, 1_790_196_400, false), 0),
        ((false, 0, 0, true), 0),
        ((false, 0, 0, false), 2),
    ] {
        let asserter = Asserter::new();
        no_answer(&asserter, halt, step);
        assert!(matches!(
            oracles(&asserter).oracle_price("NVDA").await.unwrap(),
            OraclePrice::NoPrice {
                reason: NoPrice::Halted,
                ..
            }
        ));
    }
}

#[tokio::test]
async fn sequencer_not_settled_is_no_price_and_any_other_revert_an_error() {
    let asserter = Asserter::new();
    asserter.push_success(&U64::from(7));
    revert(
        &asserter,
        MorphoBandOracle::SequencerNotSettled {}.abi_encode(),
    );
    assert!(matches!(
        oracles(&asserter).oracle_price("NVDA").await.unwrap(),
        OraclePrice::NoPrice {
            reason: NoPrice::SequencerNotSettled,
            ..
        }
    ));
    asserter.push_success(&U64::from(7));
    revert(&asserter, Vec::new());
    assert!(oracles(&asserter).oracle_price("NVDA").await.is_err());
}

#[tokio::test]
async fn an_oracles_band_symbol_tokens_scale_owner_and_halt_are_read() {
    let asserter = Asserter::new();
    let oracle = oracles(&asserter).morpho_oracle("NVDA").unwrap();
    assert_eq!(*oracle.address(), BOB);
    asserter.push_success(&Bytes::from(
        (true, 1_790_200_000u64, 1_790_196_400u64, false).abi_encode_params(),
    ));
    let halt = oracle.halt().call().await.unwrap();
    assert_eq!(
        (
            halt.signedHalt,
            halt.until,
            halt.issuedAt,
            halt.oraclePaused
        ),
        (true, 1_790_200_000, 1_790_196_400, false)
    );
    asserter.push_success(&Bytes::from(U256::from(10u64.pow(16)).abi_encode()));
    assert_eq!(
        oracle.scaleFactor().call().await.unwrap(),
        U256::from(10u64.pow(16))
    );
    asserter.push_success(&Bytes::from(b32("NVDA").abi_encode()));
    assert_eq!(oracle.symbol().call().await.unwrap(), b32("NVDA"));
    asserter.push_success(&Bytes::from(BAND.abi_encode()));
    assert_eq!(oracle.band().call().await.unwrap(), BAND);
    asserter.push_success(&Bytes::from(ALICE.abi_encode()));
    assert_eq!(oracle.collateralToken().call().await.unwrap(), ALICE);
    asserter.push_success(&Bytes::from(BOB.abi_encode()));
    assert_eq!(oracle.loanToken().call().await.unwrap(), BOB);
    asserter.push_success(&Bytes::from(ALICE.abi_encode()));
    assert_eq!(oracle.owner().call().await.unwrap(), ALICE);
    assert!(asserter.read_q().is_empty());
}

#[test]
fn a_repoints_asset_mismatch_and_band_set_are_decoded() {
    let data = MorphoBandOracle::AssetMismatch {
        band: BAND,
        token: BOB,
    }
    .abi_encode();
    assert_eq!(
        decode_revert_data(&data).unwrap().to_string(),
        format!("AssetMismatch({BAND}, {BOB})")
    );
    let set = MorphoBandOracle::BandSet::decode_raw_log(
        [
            MorphoBandOracle::BandSet::SIGNATURE_HASH,
            ALICE.into_word(),
            BAND.into_word(),
        ],
        &[],
    )
    .unwrap();
    assert_eq!((set.previousBand, set.newBand), (ALICE, BAND));
}

#[tokio::test]
async fn a_baskets_components_previews_and_targets_are_read() {
    let asserter = Asserter::new();
    let mut deployments = registry(46630);
    deployments.baskets.insert("PAIR".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    let symbols = vec![b32("NVDA"), b32("SPY")];
    asserter.push_success(&Bytes::from(
        (symbols, vec![ALICE, BOB]).abi_encode_params(),
    ));
    asserter.push_success(&Bytes::from(
        vec![U256::from(4), U256::from(2)].abi_encode(),
    ));
    asserter.push_success(&Bytes::from(
        vec![U256::from(3), U256::from(1)].abi_encode(),
    ));
    asserter.push_success(&Bytes::from(
        vec![U256::from(7), U256::from(8)].abi_encode(),
    ));
    asserter.push_success(&Bytes::from(
        (vec![U256::from(1), U256::from(2)], 1_790_604_800u64).abi_encode_params(),
    ));
    let components = tapehouse.basket_components("PAIR").await.unwrap();
    assert_eq!(components.assets, vec!["NVDA", "SPY"]);
    assert_eq!(components.tokens, vec![ALICE, BOB]);
    assert_eq!(
        tapehouse.preview_mint("PAIR", U256::from(3)).await.unwrap(),
        vec![U256::from(4), U256::from(2)]
    );
    assert_eq!(
        tapehouse
            .preview_redeem("PAIR", U256::from(3))
            .await
            .unwrap(),
        vec![U256::from(3), U256::from(1)]
    );
    assert_eq!(
        tapehouse.basket_target("PAIR").await.unwrap(),
        vec![U256::from(7), U256::from(8)]
    );
    let pending = tapehouse.pending_target("PAIR").await.unwrap();
    assert_eq!(
        (pending.units, pending.effective_at),
        (vec![U256::from(1), U256::from(2)], 1_790_604_800)
    );
    assert!(asserter.read_q().is_empty());
}

#[tokio::test]
async fn what_a_position_holds_through_its_baskets_is_read_by_asset_at_one_block() {
    let asserter = Asserter::new();
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        registry(46630),
    );
    asserter.push_success(&U64::from(9));
    asserter.push_success(&Bytes::from(
        (vec![b32("NVDA"), b32("SPY")], vec![ALICE, BOB]).abi_encode_params(),
    ));
    asserter.push_success(&Bytes::from(
        vec![U256::from(10), U256::from(5)].abi_encode(),
    ));
    let held = tapehouse.in_baskets(ALICE, CROSS).await.unwrap();
    assert_eq!(
        held.into_iter().collect::<Vec<_>>(),
        vec![
            ("NVDA".into(), U256::from(10)),
            ("SPY".into(), U256::from(5))
        ]
    );
    assert!(asserter.read_q().is_empty());
}

fn layer() -> Layer {
    Layer {
        notional: U256::from(10_000_000_000u64),
        deductible_bps: U256::from(218),
        limit_bps: U256::from(1_218),
    }
}

#[test]
fn a_gap_cover_purchase_carries_its_layer_its_premium_limit_and_its_holder() {
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("GapCover".into(), BOB);
    let cover = offline(deployments).gap_cover().unwrap();
    let l = layer();
    let buy = cover.buy(
        b32("NVDA"),
        l.notional,
        l.deductible_bps,
        l.limit_bps,
        U256::from(13_908_962),
        ALICE,
    );
    let decoded = GapCover::buyCall::abi_decode(buy.calldata()).unwrap();
    assert_eq!(
        (
            decoded.symbol,
            decoded.notional,
            decoded.deductibleBps,
            decoded.limitBps,
            decoded.maxPremium,
            decoded.holder
        ),
        (
            b32("NVDA"),
            l.notional,
            l.deductible_bps,
            l.limit_bps,
            U256::from(13_908_962),
            ALICE
        )
    );
    let reference = alloy::primitives::aliases::U80::from((1u128 << 64) | 7);
    let last = alloy::primitives::aliases::U80::from((1u128 << 64) | 8);
    let settle = cover.settle(b32("SPY"), 1_790_985_600_000, reference, last);
    let decoded = GapCover::settleCall::abi_decode(settle.calldata()).unwrap();
    assert_eq!(
        (
            decoded.symbol,
            decoded.closesMs,
            decoded.referenceRound,
            decoded.lastRound
        ),
        (b32("SPY"), 1_790_985_600_000, reference, last)
    );
    let measure = cover.measure(b32("SPY"));
    assert_eq!(
        GapCover::measureCall::abi_decode(measure.calldata())
            .unwrap()
            .symbol,
        b32("SPY")
    );
    assert_eq!(*cover.address(), BOB);
    assert_eq!(
        offline(Deployments::parse(r#"{"chainId":1}"#).unwrap())
            .gap_cover()
            .unwrap_err()
            .to_string(),
        "the registry has no .tapehouse.GapCover"
    );
}

#[tokio::test]
async fn a_gap_cover_quote_reads_the_premium_and_gives_the_writers_usdg_it_reserves() {
    let asserter = Asserter::new();
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("GapCover".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    asserter.push_success(&Bytes::from(U256::from(13_908_962).abi_encode()));
    let quote = tapehouse.quote_gap_cover("NVDA", layer()).await.unwrap();
    assert_eq!(
        (quote.premium, quote.reserve),
        (U256::from(13_908_962), U256::from(1_000_000_000u64))
    );
    asserter.push_success(&Bytes::from(
        (U256::from(279_862), U256::from(111_945)).abi_encode_params(),
    ));
    assert_eq!(
        tapehouse.gap_cover_pricing_gap("NVDA").await.unwrap(),
        PricingGap {
            gap: U256::from(279_862),
            week_move: U256::from(111_945)
        }
    );
}

#[tokio::test]
async fn a_series_reads_where_it_stands_and_its_settlement_at_one_block() {
    let asserter = Asserter::new();
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("GapCover".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    for status in [1u8, 2, 3] {
        asserter.push_success(&U64::from(7));
        asserter.push_success(&Bytes::from(GapCover::seriesCall::abi_encode_returns(
            &GapCover::seriesReturn {
                notional: 10_000_000_000,
                referencePrice: 22_244_729_849,
                price: 22_280_257_368,
                shift: 0,
                status,
                flagged: false,
            },
        )));
        if status < 3 {
            asserter.push_success(&Bytes::from(U256::from(1_789_948_800_000u64).abi_encode()));
        }
    }
    let series = tapehouse
        .gap_cover_series("NVDA", 1_789_776_000_000)
        .await
        .unwrap();
    assert_eq!(series.status, SeriesStatus::Settled);
    assert_eq!(
        (series.reference_price, series.price, series.flagged),
        (22_244_729_849, 22_280_257_368, false)
    );
    assert_eq!(
        (series.notional, series.reopen_ms),
        (
            U256::from(10_000_000_000u64),
            U256::from(1_789_948_800_000u64)
        )
    );
    let void = tapehouse.gap_cover_series("NVDA", 1).await.unwrap();
    assert_eq!(void.status, SeriesStatus::Void);
    assert_eq!(
        tapehouse
            .gap_cover_series("NVDA", 1)
            .await
            .unwrap_err()
            .to_string(),
        "unknown series status 3"
    );
    assert!(asserter.read_q().is_empty());
}

#[tokio::test]
async fn the_gap_cover_sales_read_the_closure_on_sale_and_when_they_end() {
    let asserter = Asserter::new();
    let mut deployments = registry(4663);
    deployments.tapehouse.insert("GapCover".into(), BOB);
    let tapehouse = Tapehouse::new(
        ProviderBuilder::new().connect_mocked_client(asserter.clone()),
        deployments,
    );
    for (closes, ends) in [(1_790_985_600_000u64, 1_790_985_600_000u64), (0, 0)] {
        asserter.push_success(&Bytes::from(GapCover::salesCall::abi_encode_returns(
            &GapCover::salesReturn {
                closesMs: closes,
                endsMs: ends,
            },
        )));
    }
    assert_eq!(
        tapehouse.gap_cover_sales().await.unwrap(),
        Some(Sales {
            closes_ms: 1_790_985_600_000,
            ends_ms: 1_790_985_600_000
        })
    );
    assert_eq!(tapehouse.gap_cover_sales().await.unwrap(), None);
}

#[test]
fn a_gap_cover_pays_the_fall_beyond_its_deductible_up_to_its_limit() {
    let layer = Layer {
        notional: U256::from(10_000_000_000u64),
        deductible_bps: U256::from(500),
        limit_bps: U256::from(1_500),
    };
    for (price, paid) in [
        (210e8 as u64, 0u64),
        (200e8 as u64, 0),
        (190e8 as u64, 0),
        (185e8 as u64, 250_000_000),
        (100e8 as u64, 1_000_000_000),
    ] {
        assert_eq!(
            payout(layer, 200e8 as u64, price),
            U256::from(paid),
            "{price}"
        );
    }
    let full = GapCover::NoCapacity {
        need: U256::from(2),
        free: U256::from(1),
    }
    .abi_encode();
    assert_eq!(
        decode_revert_data(&full).unwrap().to_string(),
        "NoCapacity(2, 1)"
    );
    assert_eq!(
        decode_revert_data(&GapCover::SalesClosed {}.abi_encode())
            .unwrap()
            .name,
        "SalesClosed"
    );
    let stale = GapCover::StaleReference {
        deductibleBps: U256::from(267),
        minimum: U256::from(268),
    }
    .abi_encode();
    assert_eq!(
        decode_revert_data(&stale).unwrap().to_string(),
        "StaleReference(267, 268)"
    );
    let early = GapCover::TooEarlyToMeasure {
        fromMs: U256::from(1_790_899_200_000u64),
    }
    .abi_encode();
    assert_eq!(
        decode_revert_data(&early).unwrap().to_string(),
        "TooEarlyToMeasure(1790899200000)"
    );
}
