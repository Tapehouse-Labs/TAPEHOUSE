#![cfg_attr(not(any(test, feature = "export-abi")), no_main)]
#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

pub mod chainlink;
pub mod error;
pub mod redstone;

use alloc::vec::Vec;

use stylus_sdk::abi::Bytes;
use stylus_sdk::alloy_primitives::{Address, B256, U64, U256, address};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::call::RawCall;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::{StorageAddress, StorageB256, StorageMap, StorageU64, StorageU256};

use crate::error::{
    BandError, DuplicateAsset, InvalidFeed, LengthMismatch, NoLegs, PackageNotNewer, ZeroSymbol,
};

const ECRECOVER: Address = address!("0x0000000000000000000000000000000000000001");

sol! {
    event PriceWritten(bytes32 indexed feedId, uint256 value, uint64 packageTimestampMs);
}

#[storage]
pub struct Price {
    value: StorageU256,
    package_timestamp_ms: StorageU64,
    written_at: StorageU64,
}

#[storage]
pub struct Asset {
    chainlink_feed: StorageAddress,
    redstone_feed_id: StorageB256,
}

#[storage]
#[entrypoint]
pub struct Band {
    prices: StorageMap<B256, Price>,
    assets: StorageMap<B256, Asset>,
}

#[public]
impl Band {
    /// Sets the per-asset configuration once. A zero feed or feed ID means the asset has no such leg.
    #[constructor]
    pub fn constructor(
        &mut self,
        symbols: Vec<B256>,
        chainlink_feeds: Vec<Address>,
        redstone_feed_ids: Vec<B256>,
    ) -> Result<(), BandError> {
        if symbols.len() != chainlink_feeds.len() || symbols.len() != redstone_feed_ids.len() {
            return Err(BandError::LengthMismatch(LengthMismatch {}));
        }
        for ((symbol, feed), feed_id) in symbols
            .into_iter()
            .zip(chainlink_feeds)
            .zip(redstone_feed_ids)
        {
            if symbol == B256::ZERO {
                return Err(BandError::ZeroSymbol(ZeroSymbol {}));
            }
            if feed == Address::ZERO && feed_id == B256::ZERO {
                return Err(BandError::NoLegs(NoLegs { symbol }));
            }
            if self.asset(symbol) != (Address::ZERO, B256::ZERO) {
                return Err(BandError::DuplicateAsset(DuplicateAsset { symbol }));
            }
            if feed != Address::ZERO && !chainlink::has_expected_decimals(self.vm(), feed) {
                return Err(BandError::InvalidFeed(InvalidFeed { feed }));
            }
            let mut asset = self.assets.setter(symbol);
            asset.chainlink_feed.set(feed);
            asset.redstone_feed_id.set(feed_id);
        }
        Ok(())
    }

    /// The Chainlink feed and RedStone feed ID of `symbol`. Both zero for an unknown symbol.
    pub fn asset(&self, symbol: B256) -> (Address, B256) {
        let asset = self.assets.getter(symbol);
        (asset.chainlink_feed.get(), asset.redstone_feed_id.get())
    }

    /// Both legs of `symbol`: the Chainlink answer and its `updatedAt` in seconds, then the stored
    /// RedStone 24/7 value and its package timestamp in milliseconds. Zero for a leg that is unset or
    /// unreadable. Never reverts.
    pub fn legs(&self, symbol: B256) -> (U256, u64, U256, u64) {
        let (feed, feed_id) = self.asset(symbol);
        let (chainlink_price, chainlink_updated_at) = chainlink::latest(self.vm(), feed);
        let (price_247, package_timestamp_ms) = if feed_id == B256::ZERO {
            (U256::ZERO, 0)
        } else {
            let (value, package_timestamp_ms, _) = self.price(feed_id);
            (value, package_timestamp_ms)
        };
        (
            chainlink_price,
            chainlink_updated_at,
            price_247,
            package_timestamp_ms,
        )
    }

    /// Verifies a RedStone payload for `feed_ids` and stores each value. Anyone may call it.
    pub fn write_prices(
        &mut self,
        feed_ids: Vec<B256>,
        payload: Bytes,
    ) -> Result<Vec<U256>, BandError> {
        let now = self.vm().block_timestamp();
        let verified = redstone::verify(&payload, &feed_ids, now, |hash, v, r, s| {
            ecrecover(self.vm(), hash, v, r, s)
        })?;
        for (feed_id, value) in feed_ids.iter().zip(&verified.values) {
            let stored = self
                .prices
                .getter(*feed_id)
                .package_timestamp_ms
                .get()
                .to::<u64>();
            if verified.timestamp_ms <= stored {
                return Err(BandError::PackageNotNewer(PackageNotNewer {
                    feedId: *feed_id,
                    storedTimestampMs: stored,
                    packageTimestampMs: verified.timestamp_ms,
                }));
            }
            let mut price = self.prices.setter(*feed_id);
            price.value.set(*value);
            price
                .package_timestamp_ms
                .set(U64::from(verified.timestamp_ms));
            price.written_at.set(U64::from(now));
            self.vm().log(PriceWritten {
                feedId: *feed_id,
                value: *value,
                packageTimestampMs: verified.timestamp_ms,
            });
        }
        Ok(verified.values)
    }

    /// The stored value of `feed_id`, its package timestamp in milliseconds, and the block timestamp
    /// it was written at. All zero for a feed never written.
    pub fn price(&self, feed_id: B256) -> (U256, u64, u64) {
        let price = self.prices.getter(feed_id);
        (
            price.value.get(),
            price.package_timestamp_ms.get().to::<u64>(),
            price.written_at.get().to::<u64>(),
        )
    }
}

fn ecrecover(host: &impl Host, hash: B256, v: u8, r: B256, s: B256) -> Option<Address> {
    let mut input = [0u8; 128];
    input[..32].copy_from_slice(hash.as_slice());
    input[63] = v;
    input[64..96].copy_from_slice(r.as_slice());
    input[96..].copy_from_slice(s.as_slice());
    let output = unsafe { RawCall::new_static(host).call(ECRECOVER, &input) }.ok()?;
    (output.len() == 32).then(|| Address::from_slice(&output[12..]))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::error::CalldataMustHaveValidPayload;
    use stylus_sdk::alloy_primitives::{I256, b256};
    use stylus_sdk::alloy_sol_types::SolValue;
    use stylus_sdk::testing::*;

    const HASH: B256 = b256!("0x1111111111111111111111111111111111111111111111111111111111111111");
    const R: B256 = b256!("0x2222222222222222222222222222222222222222222222222222222222222222");
    const S: B256 = b256!("0x3333333333333333333333333333333333333333333333333333333333333333");

    fn precompile_input(v: u8) -> Vec<u8> {
        [
            HASH.as_slice(),
            B256::with_last_byte(v).as_slice(),
            R.as_slice(),
            S.as_slice(),
        ]
        .concat()
    }

    #[test]
    fn unwritten_feed_reads_as_zero() {
        let vm = TestVM::default();
        let band = Band::from(&vm);
        assert_eq!(band.price(B256::repeat_byte(1)), (U256::ZERO, 0, 0));
    }

    #[test]
    fn ecrecover_returns_the_precompile_output() {
        let vm = TestVM::default();
        let signer = address!("0x8BB8F32Df04c8b654987DAaeD53D6B6091e3B774");
        vm.mock_static_call(
            ECRECOVER,
            precompile_input(28),
            Ok(signer.into_word().to_vec()),
        );
        assert_eq!(ecrecover(&vm, HASH, 28, R, S), Some(signer));
    }

    #[test]
    fn ecrecover_without_output_is_none() {
        let vm = TestVM::default();
        vm.mock_static_call(ECRECOVER, precompile_input(27), Ok(Vec::new()));
        assert_eq!(ecrecover(&vm, HASH, 27, R, S), None);
    }

    #[test]
    fn ecrecover_sends_the_precompile_layout() {
        let vm = TestVM::default();
        let word = address!("0x8BB8F32Df04c8b654987DAaeD53D6B6091e3B774").into_word();
        vm.mock_static_call(ECRECOVER, precompile_input(28), Err(word.to_vec()));
        assert_eq!(ecrecover(&vm, HASH, 28, R, S), None);
    }

    const FEED: Address = address!("0x4444444444444444444444444444444444444444");
    const DECIMALS_CALL: [u8; 4] = [0x31, 0x3c, 0xe5, 0x67];
    const LATEST_ROUND_DATA_CALL: [u8; 4] = [0xfe, 0xaf, 0x96, 0x8c];

    fn symbol(name: &str) -> B256 {
        let mut id = [0u8; 32];
        id[..name.len()].copy_from_slice(name.as_bytes());
        B256::from(id)
    }

    fn round(answer: i64, updated_at: u64) -> Vec<u8> {
        (
            U256::from(7u8),
            I256::try_from(answer).unwrap(),
            U256::from(updated_at - 12),
            U256::from(updated_at),
            U256::from(7u8),
        )
            .abi_encode_params()
    }

    #[test]
    fn constructor_rejects_arrays_of_different_lengths() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(vec![symbol("NVDA")], vec![], vec![symbol("NVDA---24_7")]);
        assert_eq!(result, Err(BandError::LengthMismatch(LengthMismatch {})));
        let result = band.constructor(vec![symbol("NVDA")], vec![Address::ZERO], vec![]);
        assert_eq!(result, Err(BandError::LengthMismatch(LengthMismatch {})));
    }

    #[test]
    fn constructor_rejects_a_zero_symbol() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![B256::ZERO],
            vec![Address::ZERO],
            vec![symbol("NVDA---24_7")],
        );
        assert_eq!(result, Err(BandError::ZeroSymbol(ZeroSymbol {})));
    }

    #[test]
    fn constructor_rejects_an_asset_without_legs() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(vec![symbol("SPY")], vec![Address::ZERO], vec![B256::ZERO]);
        assert_eq!(
            result,
            Err(BandError::NoLegs(NoLegs {
                symbol: symbol("SPY")
            }))
        );
    }

    #[test]
    fn constructor_rejects_a_duplicate_symbol() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![symbol("NVDA"), symbol("NVDA")],
            vec![Address::ZERO, Address::ZERO],
            vec![symbol("NVDA---24_7"), symbol("TSLA---24_7")],
        );
        assert_eq!(
            result,
            Err(BandError::DuplicateAsset(DuplicateAsset {
                symbol: symbol("NVDA")
            }))
        );
    }

    #[test]
    fn constructor_rejects_a_feed_with_other_decimals() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(
            FEED,
            DECIMALS_CALL.to_vec(),
            Ok(U256::from(18u8).to_be_bytes::<32>().to_vec()),
        );
        let result = band.constructor(vec![symbol("NVDA")], vec![FEED], vec![B256::ZERO]);
        assert_eq!(
            result,
            Err(BandError::InvalidFeed(InvalidFeed { feed: FEED }))
        );
    }

    #[test]
    fn constructor_rejects_a_feed_without_code() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(vec![symbol("NVDA")], vec![FEED], vec![B256::ZERO]);
        assert_eq!(
            result,
            Err(BandError::InvalidFeed(InvalidFeed { feed: FEED }))
        );
    }

    #[test]
    fn constructor_stores_each_asset() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(
            FEED,
            DECIMALS_CALL.to_vec(),
            Ok(U256::from(8u8).to_be_bytes::<32>().to_vec()),
        );
        band.constructor(
            vec![symbol("NVDA"), symbol("TSLA")],
            vec![FEED, Address::ZERO],
            vec![symbol("NVDA---24_7"), symbol("TSLA---24_7")],
        )
        .unwrap();
        assert_eq!(band.asset(symbol("NVDA")), (FEED, symbol("NVDA---24_7")));
        assert_eq!(
            band.asset(symbol("TSLA")),
            (Address::ZERO, symbol("TSLA---24_7"))
        );
        assert_eq!(band.asset(symbol("AAPL")), (Address::ZERO, B256::ZERO));
    }

    #[test]
    fn legs_read_the_feed_and_the_stored_price() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        band.assets.setter(symbol("NVDA")).chainlink_feed.set(FEED);
        band.assets
            .setter(symbol("NVDA"))
            .redstone_feed_id
            .set(symbol("NVDA---24_7"));
        band.prices
            .setter(symbol("NVDA---24_7"))
            .value
            .set(U256::from(22_865_213_250u64));
        band.prices
            .setter(symbol("NVDA---24_7"))
            .package_timestamp_ms
            .set(U64::from(1_790_126_530_000u64));
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(22_902_739_517, 1_790_089_676)),
        );
        assert_eq!(
            band.legs(symbol("NVDA")),
            (
                U256::from(22_902_739_517u64),
                1_790_089_676,
                U256::from(22_865_213_250u64),
                1_790_126_530_000
            )
        );
    }

    #[test]
    fn legs_of_an_unknown_symbol_make_no_call() {
        let vm = TestVM::default();
        let band = Band::from(&vm);
        vm.mock_static_call(
            Address::ZERO,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(22_902_739_517, 1_790_089_676)),
        );
        assert_eq!(band.legs(symbol("NVDA")), (U256::ZERO, 0, U256::ZERO, 0));
    }

    #[test]
    fn a_non_positive_answer_is_no_reading() {
        let vm = TestVM::default();
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(0, 1_790_089_676)),
        );
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0));
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(-1, 1_790_089_676)),
        );
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0));
    }

    #[test]
    fn a_reverting_feed_is_no_reading() {
        let vm = TestVM::default();
        vm.mock_static_call(FEED, LATEST_ROUND_DATA_CALL.to_vec(), Err(vec![]));
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0));
    }

    #[test]
    fn a_feed_without_code_is_no_reading() {
        let vm = TestVM::default();
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0));
    }

    #[test]
    fn write_without_a_payload_reverts() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.write_prices(vec![B256::repeat_byte(1)], Bytes::from(vec![0u8; 32]));
        assert_eq!(
            result,
            Err(BandError::CalldataMustHaveValidPayload(
                CalldataMustHaveValidPayload {}
            ))
        );
        assert!(vm.get_emitted_logs().is_empty());
    }
}
