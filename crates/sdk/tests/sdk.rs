// SPDX-License-Identifier: MIT OR Apache-2.0
use std::error::Error as _;

use alloy::primitives::{Address, B256, Bytes, U64, U256, address, aliases::U24};
use alloy::providers::{Provider, ProviderBuilder};
use alloy::rpc::json_rpc::ErrorPayload;
use alloy::sol_types::{Revert as SolidityRevert, SolCall, SolError, SolValue};
use alloy::transports::mock::Asserter;
use tapehouse_sdk::bindings::{Band, BandFeed, MarginAccounts, ShortPositions, StockToken};
use tapehouse_sdk::{
    CROSS, Deployments, Error, PackageSource, Tapehouse, decode_revert_data, to_bytes32,
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
