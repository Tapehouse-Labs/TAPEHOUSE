//! Verification of RedStone signed data packages, matching
//! `PrimaryProdDataServiceConsumerBase` in `@redstone-finance/evm-connector` 1.0.0.

use alloc::vec;
use alloc::vec::Vec;

use stylus_sdk::alloy_primitives::{Address, B256, U256, address, b256};
use stylus_sdk::crypto::keccak;

use crate::error::*;

/// Signers of the `redstone-primary-prod` data service, in the reference contract's index order.
pub const AUTHORISED_SIGNERS: [Address; 5] = [
    address!("0x8BB8F32Df04c8b654987DAaeD53D6B6091e3B774"),
    address!("0xdEB22f54738d54976C4c0fe5ce6d408E40d88499"),
    address!("0x51Ce04Be4b3E32572C4Ec9135221d0691Ba7d202"),
    address!("0xDD682daEC5A90dD295d14DA4b0bec9281017b5bE"),
    address!("0x9c5AE89C4Af6aA32cE58588DBaF90d18a855B6de"),
];

/// Unique signers required per data feed.
pub const UNIQUE_SIGNERS_THRESHOLD: usize = 3;

/// Maximum age of a package, in seconds.
pub const MAX_DATA_AGE_S: u64 = 180;

/// Maximum distance of a package into the future, in seconds.
pub const MAX_DATA_AHEAD_S: u64 = 60;

const REDSTONE_MARKER: [u8; 9] = [0x00, 0x00, 0x02, 0xed, 0x57, 0x01, 0x1e, 0x00, 0x00];
const SIGNATURE_BS: usize = 65;
const DATA_POINTS_COUNT_BS: usize = 3;
const DATA_POINT_VALUE_BYTE_SIZE_BS: usize = 4;
const TIMESTAMP_BS: usize = 6;
const DATA_PACKAGES_COUNT_BS: usize = 2;
const UNSIGNED_METADATA_BYTE_SIZE_BS: usize = 3;
const SECP256K1_HALF_ORDER: B256 =
    b256!("0x7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0");

/// Values for the requested feeds, in request order, and the packages' common timestamp in milliseconds.
#[derive(Debug, PartialEq, Eq)]
pub struct Verified {
    pub values: Vec<U256>,
    pub timestamp_ms: u64,
}

/// Verifies `payload` for `feed_ids` at `block_timestamp` (seconds).
///
/// `recover` returns the signer of a prehashed message from `(hash, v, r, s)`, or `None` when
/// recovery fails. Packages are read from the end of the payload, as the reference contract does,
/// and each requested feed keeps the values of its first three unique authorised signers.
pub fn verify(
    payload: &[u8],
    feed_ids: &[B256],
    block_timestamp: u64,
    recover: impl Fn(B256, u8, B256, B256) -> Option<Address>,
) -> Result<Verified, BandError> {
    let mut end = trailer_start(payload)?;
    let packages_count = read_size(payload, &mut end, DATA_PACKAGES_COUNT_BS)?;

    let mut counts = vec![0usize; feed_ids.len()];
    let mut bitmaps = vec![0u8; feed_ids.len()];
    let mut values = vec![[U256::ZERO; UNIQUE_SIGNERS_THRESHOLD]; feed_ids.len()];
    let mut timestamp_ms = 0u64;

    for _ in 0..packages_count {
        let signature_start = end.checked_sub(SIGNATURE_BS).ok_or(overflow())?;
        let signature = &payload[signature_start..end];
        end = signature_start;
        let points_count = read_size(payload, &mut end, DATA_POINTS_COUNT_BS)?;
        let value_size = read_size(payload, &mut end, DATA_POINT_VALUE_BYTE_SIZE_BS)?;
        let package_timestamp = read_uint(payload, &mut end, TIMESTAMP_BS)?;
        let points_size = points_count
            .checked_mul(32 + value_size)
            .ok_or(overflow())?;
        let points_start = end.checked_sub(points_size).ok_or(overflow())?;

        let signed_hash = keccak(&payload[points_start..signature_start]);
        let signer = recover_signer(signed_hash, signature, &recover)?;
        let signer_index = AUTHORISED_SIGNERS.iter().position(|s| *s == signer).ok_or(
            BandError::SignerNotAuthorised(SignerNotAuthorised {
                receivedSigner: signer,
            }),
        )?;

        for point in (0..points_count).rev() {
            if value_size > 32 {
                return Err(BandError::TooLargeValueByteSize(TooLargeValueByteSize {
                    valueByteSize: U256::from(value_size),
                }));
            }
            let offset = points_start + point * (32 + value_size);
            let feed_id = B256::from_slice(&payload[offset..offset + 32]);
            let Some(feed) = feed_ids.iter().position(|id| *id == feed_id) else {
                continue;
            };
            let bit = 1u8 << signer_index;
            if bitmaps[feed] & bit == 0 && counts[feed] < UNIQUE_SIGNERS_THRESHOLD {
                values[feed][counts[feed]] =
                    U256::from_be_slice(&payload[offset + 32..offset + 32 + value_size]);
                counts[feed] += 1;
                bitmaps[feed] |= bit;
            }
        }
        if package_timestamp == 0 {
            return Err(BandError::DataTimestampCannotBeZero(
                DataTimestampCannotBeZero {},
            ));
        }
        if timestamp_ms == 0 {
            timestamp_ms = package_timestamp;
        } else if package_timestamp != timestamp_ms {
            return Err(BandError::TimestampsMustBeEqual(TimestampsMustBeEqual {}));
        }
        end = points_start;
    }

    let mut medians = Vec::with_capacity(feed_ids.len());
    for (feed, collected) in values.iter_mut().enumerate() {
        if counts[feed] < UNIQUE_SIGNERS_THRESHOLD {
            return Err(BandError::InsufficientNumberOfUniqueSigners(
                InsufficientNumberOfUniqueSigners {
                    receivedSignersCount: U256::from(counts[feed]),
                    requiredSignersCount: U256::from(UNIQUE_SIGNERS_THRESHOLD),
                },
            ));
        }
        collected.sort_unstable();
        medians.push(collected[UNIQUE_SIGNERS_THRESHOLD / 2]);
    }

    validate_timestamp(timestamp_ms / 1000, block_timestamp)?;
    Ok(Verified {
        values: medians,
        timestamp_ms,
    })
}

fn trailer_start(payload: &[u8]) -> Result<usize, BandError> {
    if payload.len() < REDSTONE_MARKER.len()
        || payload[payload.len() - REDSTONE_MARKER.len()..] != REDSTONE_MARKER
    {
        return Err(BandError::CalldataMustHaveValidPayload(
            CalldataMustHaveValidPayload {},
        ));
    }
    let mut end = payload.len() - REDSTONE_MARKER.len();
    let metadata_size = read_size(payload, &mut end, UNSIGNED_METADATA_BYTE_SIZE_BS)?;
    end.checked_sub(metadata_size).ok_or(overflow())
}

fn read_uint(payload: &[u8], end: &mut usize, size: usize) -> Result<u64, BandError> {
    let start = end.checked_sub(size).ok_or(overflow())?;
    let value = payload[start..*end]
        .iter()
        .fold(0u64, |acc, byte| (acc << 8) | u64::from(*byte));
    *end = start;
    Ok(value)
}

fn read_size(payload: &[u8], end: &mut usize, size: usize) -> Result<usize, BandError> {
    usize::try_from(read_uint(payload, end, size)?).map_err(|_| overflow())
}

fn recover_signer(
    signed_hash: B256,
    signature: &[u8],
    recover: &impl Fn(B256, u8, B256, B256) -> Option<Address>,
) -> Result<Address, BandError> {
    let invalid = || {
        BandError::InvalidSignature(InvalidSignature {
            signedHash: signed_hash,
        })
    };
    let r = B256::from_slice(&signature[..32]);
    let s = B256::from_slice(&signature[32..64]);
    let v = signature[64];
    if (v != 27 && v != 28) || s > SECP256K1_HALF_ORDER {
        return Err(invalid());
    }
    match recover(signed_hash, v, r, s) {
        Some(signer) if signer != Address::ZERO => Ok(signer),
        _ => Err(invalid()),
    }
}

fn validate_timestamp(timestamp_s: u64, block_timestamp: u64) -> Result<(), BandError> {
    if block_timestamp < timestamp_s {
        if timestamp_s - block_timestamp > MAX_DATA_AHEAD_S {
            return Err(BandError::TimestampFromTooLongFuture(
                TimestampFromTooLongFuture {
                    receivedTimestampSeconds: U256::from(timestamp_s),
                    blockTimestamp: U256::from(block_timestamp),
                },
            ));
        }
    } else if block_timestamp - timestamp_s > MAX_DATA_AGE_S {
        return Err(BandError::TimestampIsTooOld(TimestampIsTooOld {
            receivedTimestampSeconds: U256::from(timestamp_s),
            blockTimestamp: U256::from(block_timestamp),
        }));
    }
    Ok(())
}

fn overflow() -> BandError {
    BandError::CalldataOverOrUnderFlow(CalldataOverOrUnderFlow {})
}

#[cfg(test)]
mod tests {
    use super::*;
    use alloy_primitives::{Signature, hex};

    const NVDA: &str = include_str!("../testdata/nvda-24_7.hex");
    const STATUS: &str = include_str!("../testdata/status.hex");
    const NY_MARKET_STATUS: &str = include_str!("../testdata/ny-market-status.hex");
    const PACKAGE_BS: usize = 142;
    const TIMESTAMP_S: u64 = 1_790_119_750;

    fn feed(name: &str) -> B256 {
        let mut id = [0u8; 32];
        id[..name.len()].copy_from_slice(name.as_bytes());
        B256::from(id)
    }

    fn recover(hash: B256, v: u8, r: B256, s: B256) -> Option<Address> {
        let signature = Signature::new(r.into(), s.into(), v == 28);
        signature.recover_address_from_prehash(&hash).ok()
    }

    fn nvda() -> Vec<u8> {
        hex::decode(NVDA).unwrap()
    }

    fn verify_nvda(payload: &[u8], block_timestamp: u64) -> Result<Verified, BandError> {
        verify(payload, &[feed("NVDA---24_7")], block_timestamp, recover)
    }

    fn repack(indices: &[usize]) -> Vec<u8> {
        let payload = nvda();
        let mut out = Vec::new();
        for i in indices {
            out.extend_from_slice(&payload[i * PACKAGE_BS..(i + 1) * PACKAGE_BS]);
        }
        out.extend_from_slice(&(indices.len() as u16).to_be_bytes());
        out.extend_from_slice(&[0, 0, 0]);
        out.extend_from_slice(&REDSTONE_MARKER);
        out
    }

    fn flipped(offset: usize) -> Vec<u8> {
        let mut payload = nvda();
        payload[offset] ^= 0x01;
        payload
    }

    fn insufficient(received: u64) -> BandError {
        BandError::InsufficientNumberOfUniqueSigners(InsufficientNumberOfUniqueSigners {
            receivedSignersCount: U256::from(received),
            requiredSignersCount: U256::from(3),
        })
    }

    #[test]
    fn nvda_median_matches_reference() {
        let verified = verify_nvda(&nvda(), TIMESTAMP_S).unwrap();
        assert_eq!(verified.values, vec![U256::from(22_846_110_842u64)]);
        assert_eq!(verified.timestamp_ms, 1_790_119_750_000);
    }

    #[test]
    fn market_status_feeds_match_reference() {
        let feeds = [
            feed("NY_MARKET_CURRENT_STATUS"),
            feed("NY_MARKET_NEXT_STATUS"),
            feed("NY_MARKET_NEXT_CHANGE_TIME"),
        ];
        let verified = verify(&hex::decode(STATUS).unwrap(), &feeds, TIMESTAMP_S, recover).unwrap();
        let change_time = U256::from(179_017_020_000_000_000_000u128);
        assert_eq!(
            verified.values,
            vec![
                U256::from(101_000_000u64),
                U256::from(100_000_000u64),
                change_time
            ]
        );
    }

    #[test]
    fn multi_point_packages_match_reference() {
        let feeds = [
            feed("NY_MARKET_CURRENT_STATUS"),
            feed("NY_MARKET_NEXT_STATUS"),
            feed("NY_MARKET_NEXT_CHANGE_TIME"),
        ];
        let payload = hex::decode(NY_MARKET_STATUS).unwrap();
        let verified = verify(&payload, &feeds, 1_790_122_370, recover).unwrap();
        let change_time = U256::from(179_017_020_000_000_000_000u128);
        assert_eq!(
            verified.values,
            vec![
                U256::from(101_000_000u64),
                U256::from(100_000_000u64),
                change_time
            ]
        );
    }

    #[test]
    fn duplicate_feed_id_fills_only_its_first_slot() {
        let feeds = [feed("NY_MARKET_NEXT_STATUS"), feed("NY_MARKET_NEXT_STATUS")];
        let payload = hex::decode(NY_MARKET_STATUS).unwrap();
        assert_eq!(
            verify(&payload, &feeds, 1_790_122_370, recover),
            Err(insufficient(0))
        );
    }

    #[test]
    fn flipped_r_recovers_an_unknown_signer() {
        let signer = address!("0xcfFc4C8940d9E1A2B279AEd3bF793AaA843070c4");
        assert_eq!(
            verify_nvda(&flipped(4 * PACKAGE_BS + 77), TIMESTAMP_S),
            Err(BandError::SignerNotAuthorised(SignerNotAuthorised {
                receivedSigner: signer
            }))
        );
    }

    #[test]
    fn flipped_s_recovers_an_unknown_signer() {
        let signer = address!("0xFd5341950D4FE1A88D64046Ea9F1386950c9d458");
        assert_eq!(
            verify_nvda(&flipped(4 * PACKAGE_BS + 109), TIMESTAMP_S),
            Err(BandError::SignerNotAuthorised(SignerNotAuthorised {
                receivedSigner: signer
            }))
        );
    }

    #[test]
    fn invalid_v_is_an_invalid_signature() {
        let hash = b256!("0xf8a8a216a21ce20c4c14f35e003345bf1213f1f38a3ed611458fee26d80f9e22");
        assert_eq!(
            verify_nvda(&flipped(4 * PACKAGE_BS + 141), TIMESTAMP_S),
            Err(BandError::InvalidSignature(InvalidSignature {
                signedHash: hash
            }))
        );
    }

    #[test]
    fn flipped_value_recovers_an_unknown_signer() {
        let signer = address!("0x2BA9ca41588443baa99386a8D1797d9648c61f1E");
        assert_eq!(
            verify_nvda(&flipped(4 * PACKAGE_BS + 63), TIMESTAMP_S),
            Err(BandError::SignerNotAuthorised(SignerNotAuthorised {
                receivedSigner: signer
            }))
        );
    }

    #[test]
    fn high_s_is_an_invalid_signature() {
        let mut payload = nvda();
        let s = 4 * PACKAGE_BS + 109;
        payload[s..s + 32].copy_from_slice(&[0xff; 32]);
        assert!(matches!(
            verify_nvda(&payload, TIMESTAMP_S),
            Err(BandError::InvalidSignature(_))
        ));
    }

    #[test]
    fn malleable_twin_is_an_invalid_signature() {
        let n = U256::from_be_bytes(
            b256!("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141").0,
        );
        let mut payload = nvda();
        let s = 4 * PACKAGE_BS + 109;
        let twin = n - U256::from_be_slice(&payload[s..s + 32]);
        payload[s..s + 32].copy_from_slice(&twin.to_be_bytes::<32>());
        payload[s + 32] = 55 - payload[s + 32];
        let hash = b256!("0xf8a8a216a21ce20c4c14f35e003345bf1213f1f38a3ed611458fee26d80f9e22");
        assert_eq!(
            verify_nvda(&payload, TIMESTAMP_S),
            Err(BandError::InvalidSignature(InvalidSignature {
                signedHash: hash
            }))
        );
    }

    #[test]
    fn oversized_value_recovers_an_unknown_signer() {
        let mut payload = nvda();
        payload[4 * PACKAGE_BS + 73] = 33;
        let signer = address!("0x2f47599de0cc267e880666db678d0996cddd9d26");
        assert_eq!(
            verify_nvda(&payload, TIMESTAMP_S),
            Err(BandError::SignerNotAuthorised(SignerNotAuthorised {
                receivedSigner: signer
            }))
        );
    }

    #[test]
    fn packages_with_different_timestamps_are_rejected() {
        let status = hex::decode(NY_MARKET_STATUS).unwrap();
        let mut payload = nvda()[..3 * PACKAGE_BS].to_vec();
        payload.extend_from_slice(&status[..270]);
        payload.extend_from_slice(&4u16.to_be_bytes());
        payload.extend_from_slice(&[0, 0, 0]);
        payload.extend_from_slice(&REDSTONE_MARKER);
        assert_eq!(
            verify_nvda(&payload, TIMESTAMP_S),
            Err(BandError::TimestampsMustBeEqual(TimestampsMustBeEqual {}))
        );
    }

    #[test]
    fn two_signers_are_insufficient() {
        assert_eq!(
            verify_nvda(&repack(&[3, 4]), TIMESTAMP_S),
            Err(insufficient(2))
        );
    }

    #[test]
    fn duplicate_signer_is_ignored() {
        let verified = verify_nvda(&repack(&[0, 1, 2, 2]), TIMESTAMP_S).unwrap();
        assert_eq!(verified.values, vec![U256::from(22_846_115_210u64)]);
    }

    #[test]
    fn duplicate_signer_does_not_count() {
        assert_eq!(
            verify_nvda(&repack(&[0, 1, 1]), TIMESTAMP_S),
            Err(insufficient(2))
        );
    }

    #[test]
    fn unrequested_feed_is_insufficient() {
        assert_eq!(
            verify(&nvda(), &[feed("TSLA---24_7")], TIMESTAMP_S, recover),
            Err(insufficient(0))
        );
    }

    #[test]
    fn age_of_180_seconds_is_accepted() {
        assert!(verify_nvda(&nvda(), TIMESTAMP_S + 180).is_ok());
    }

    #[test]
    fn age_of_181_seconds_is_too_old() {
        assert_eq!(
            verify_nvda(&nvda(), TIMESTAMP_S + 181),
            Err(BandError::TimestampIsTooOld(TimestampIsTooOld {
                receivedTimestampSeconds: U256::from(TIMESTAMP_S),
                blockTimestamp: U256::from(TIMESTAMP_S + 181),
            }))
        );
    }

    #[test]
    fn sixty_seconds_ahead_is_accepted() {
        assert!(verify_nvda(&nvda(), TIMESTAMP_S - 60).is_ok());
    }

    #[test]
    fn sixty_one_seconds_ahead_is_too_far() {
        assert_eq!(
            verify_nvda(&nvda(), TIMESTAMP_S - 61),
            Err(BandError::TimestampFromTooLongFuture(
                TimestampFromTooLongFuture {
                    receivedTimestampSeconds: U256::from(TIMESTAMP_S),
                    blockTimestamp: U256::from(TIMESTAMP_S - 61),
                }
            ))
        );
    }

    #[test]
    fn missing_marker_is_not_a_payload() {
        let mut payload = nvda();
        *payload.last_mut().unwrap() = 0x01;
        assert_eq!(
            verify_nvda(&payload, TIMESTAMP_S),
            Err(BandError::CalldataMustHaveValidPayload(
                CalldataMustHaveValidPayload {}
            ))
        );
    }

    #[test]
    fn truncated_payload_overflows() {
        let payload = nvda();
        let truncated = [
            &payload[PACKAGE_BS * 4 + 10..PACKAGE_BS * 5],
            &payload[PACKAGE_BS * 5..],
        ]
        .concat();
        assert_eq!(verify_nvda(&truncated, TIMESTAMP_S), Err(overflow()));
    }

    #[test]
    fn empty_payload_is_not_a_payload() {
        assert_eq!(
            verify_nvda(&[], TIMESTAMP_S),
            Err(BandError::CalldataMustHaveValidPayload(
                CalldataMustHaveValidPayload {}
            ))
        );
    }
}
