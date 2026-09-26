#![cfg_attr(not(any(test, feature = "export-abi")), no_main)]
#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

pub mod error;
pub mod redstone;

use alloc::vec::Vec;

use stylus_sdk::abi::Bytes;
use stylus_sdk::alloy_primitives::{Address, B256, U64, U256, address};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::call::RawCall;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::{StorageMap, StorageU64, StorageU256};

use crate::error::{BandError, PackageNotNewer};

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
#[entrypoint]
pub struct Band {
    prices: StorageMap<B256, Price>,
}

#[public]
impl Band {
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
    use stylus_sdk::alloy_primitives::b256;
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
