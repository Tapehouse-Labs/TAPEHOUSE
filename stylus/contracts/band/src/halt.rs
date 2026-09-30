// SPDX-License-Identifier: MIT OR Apache-2.0
//! Signed trading halts. Tapehouse's halt signer signs an EIP-712 `HaltState` per asset from
//! Robinhood's reported trading halt; anyone may write it to the band, which holds the asset halted
//! until the message expires or a newer one lifts it.

use stylus_sdk::alloy_primitives::{Address, B256, b256};
use stylus_sdk::crypto::keccak;

use crate::error::{BandError, HaltNotNewer, HaltOutsideWindow, InvalidSignature};

/// Longest a signed message stays valid, in seconds. A halt the signer stops refreshing lapses
/// within this time.
pub const MAX_HALT_S: u64 = 3_600;

const DOMAIN_TYPEHASH: B256 =
    b256!("0x8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f");
const NAME_HASH: B256 = b256!("0x3a9c8c66b6f492cbee4758bf28448f26f766ca26a2e140e5661eff1c6a3fa547");
const VERSION_HASH: B256 =
    b256!("0xc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6");
const HALT_TYPEHASH: B256 =
    b256!("0x056addf8508be82b0feb7d43591333e6da6d0bee53cb2b9ea21326cac3917222");
const SIGNATURE_BS: usize = 65;

/// A written halt: when its message was issued and until when the asset is halted, in seconds.
/// All zero before the first message.
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq)]
pub struct Halt {
    /// The `issuedAt` of the last message written.
    pub issued_at: u64,
    /// The asset is halted while the block timestamp is below this; zero once lifted.
    pub until: u64,
}

/// The EIP-712 domain separator of the band at `band` on chain `chain_id`: name "Tapehouse Band",
/// version "1".
pub fn domain_separator(chain_id: u64, band: Address) -> B256 {
    let mut encoded = [0u8; 160];
    encoded[..32].copy_from_slice(DOMAIN_TYPEHASH.as_slice());
    encoded[32..64].copy_from_slice(NAME_HASH.as_slice());
    encoded[64..96].copy_from_slice(VERSION_HASH.as_slice());
    encoded[120..128].copy_from_slice(&chain_id.to_be_bytes());
    encoded[140..].copy_from_slice(band.as_slice());
    keccak(encoded)
}

/// The EIP-712 hash of `HaltState(bytes32 symbol,bool halted,uint64 issuedAt,uint64 expiresAt)`
/// under `domain`.
pub fn signing_hash(
    domain: B256,
    symbol: B256,
    halted: bool,
    issued_at: u64,
    expires_at: u64,
) -> B256 {
    let mut encoded = [0u8; 160];
    encoded[..32].copy_from_slice(HALT_TYPEHASH.as_slice());
    encoded[32..64].copy_from_slice(symbol.as_slice());
    encoded[95] = halted.into();
    encoded[120..128].copy_from_slice(&issued_at.to_be_bytes());
    encoded[152..].copy_from_slice(&expires_at.to_be_bytes());
    let mut message = [0u8; 66];
    message[..2].copy_from_slice(&[0x19, 0x01]);
    message[2..34].copy_from_slice(domain.as_slice());
    message[34..].copy_from_slice(keccak(encoded).as_slice());
    keccak(message)
}

/// The signer of `hash` from a 65-byte `r ‖ s ‖ v` signature, with RedStone's checks: `v` of 27 or
/// 28, low `s`, and a non-zero signer.
pub fn signer(
    hash: B256,
    signature: &[u8],
    recover: impl Fn(B256, u8, B256, B256) -> Option<Address>,
) -> Result<Address, BandError> {
    if signature.len() != SIGNATURE_BS {
        return Err(BandError::InvalidSignature(InvalidSignature {
            signedHash: hash,
        }));
    }
    crate::redstone::recover_signer(hash, signature, &recover)
}

/// The halt to store after a message for `symbol` signed by the halt signer, at block timestamp
/// `now`. The message must be issued no later than `now`, expire after it, span at most
/// [`MAX_HALT_S`], and be issued after the stored one.
pub fn accept(
    symbol: B256,
    stored: Halt,
    halted: bool,
    issued_at: u64,
    expires_at: u64,
    now: u64,
) -> Result<Halt, BandError> {
    if issued_at > now || expires_at <= now || expires_at - issued_at > MAX_HALT_S {
        return Err(BandError::HaltOutsideWindow(HaltOutsideWindow {
            issuedAt: issued_at,
            expiresAt: expires_at,
            blockTimestamp: now,
        }));
    }
    if issued_at <= stored.issued_at {
        return Err(BandError::HaltNotNewer(HaltNotNewer {
            symbol,
            storedIssuedAt: stored.issued_at,
            issuedAt: issued_at,
        }));
    }
    Ok(Halt {
        issued_at,
        until: if halted { expires_at } else { 0 },
    })
}

/// Whether `halt` holds the asset halted at `now`.
pub fn active(halt: Halt, now: u64) -> bool {
    now < halt.until
}

/// Whether `halt` expired at or before `now` without a message lifting it.
pub fn lapsed(halt: Halt, now: u64) -> bool {
    halt.until != 0 && now >= halt.until
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;
    use serde_json::Value;
    use stylus_sdk::alloy_primitives::{Signature, hex};

    const VECTORS: &str = include_str!("../testdata/halt-vectors.json");

    fn int(value: &Value) -> u64 {
        value.as_str().unwrap().parse().unwrap()
    }

    fn b256(value: &Value) -> B256 {
        value.as_str().unwrap().parse().unwrap()
    }

    fn halt(value: &Value) -> Halt {
        Halt {
            issued_at: int(&value["issued_at"]),
            until: int(&value["until"]),
        }
    }

    fn cases(section: &str) -> Vec<Value> {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        vectors[section].as_array().unwrap().clone()
    }

    fn recover(hash: B256, v: u8, r: B256, s: B256) -> Option<Address> {
        let mut bytes = [0u8; 65];
        bytes[..32].copy_from_slice(r.as_slice());
        bytes[32..64].copy_from_slice(s.as_slice());
        bytes[64] = v;
        Signature::from_raw(&bytes)
            .ok()?
            .recover_address_from_prehash(&hash)
            .ok()
    }

    #[test]
    fn type_hashes_are_the_hashes_of_their_strings() {
        let domain =
            "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)";
        assert_eq!(DOMAIN_TYPEHASH, keccak(domain));
        assert_eq!(NAME_HASH, keccak("Tapehouse Band"));
        assert_eq!(VERSION_HASH, keccak("1"));
        let halt = "HaltState(bytes32 symbol,bool halted,uint64 issuedAt,uint64 expiresAt)";
        assert_eq!(HALT_TYPEHASH, keccak(halt));
    }

    #[test]
    fn domain_separators_match_the_reference() {
        for case in cases("domain") {
            let band: Address = case["band"].as_str().unwrap().parse().unwrap();
            let got = domain_separator(int(&case["chain_id"]), band);
            assert_eq!(got, b256(&case["expected"]), "{}", case["name"]);
        }
    }

    #[test]
    fn signed_messages_recover_to_the_halt_signer() {
        for case in cases("signed") {
            let band: Address = case["band"].as_str().unwrap().parse().unwrap();
            let hash = signing_hash(
                domain_separator(int(&case["chain_id"]), band),
                b256(&case["symbol"]),
                case["halted"].as_bool().unwrap(),
                int(&case["issued_at"]),
                int(&case["expires_at"]),
            );
            assert_eq!(hash, b256(&case["hash"]), "{}", case["name"]);
            let signature = hex::decode(case["signature"].as_str().unwrap()).unwrap();
            let expected: Address = case["signer"].as_str().unwrap().parse().unwrap();
            assert_eq!(
                signer(hash, &signature, recover),
                Ok(expected),
                "{}",
                case["name"]
            );
        }
    }

    #[test]
    fn a_signature_of_the_wrong_length_is_invalid() {
        let case = &cases("signed")[0];
        let hash = b256(&case["hash"]);
        let signature = hex::decode(case["signature"].as_str().unwrap()).unwrap();
        for bytes in [&signature[..64], &[signature.as_slice(), &[0]].concat()[..]] {
            assert_eq!(
                signer(hash, bytes, recover),
                Err(BandError::InvalidSignature(InvalidSignature {
                    signedHash: hash
                }))
            );
        }
    }

    #[test]
    fn acceptance_matches_the_reference() {
        let symbol = B256::repeat_byte(7);
        for case in cases("accept") {
            let got = accept(
                symbol,
                halt(&case["stored"]),
                case["halted"].as_bool().unwrap(),
                int(&case["issued_at"]),
                int(&case["expires_at"]),
                int(&case["now"]),
            );
            match case["expected"].as_str() {
                Some("HaltNotNewer") => {
                    assert!(
                        matches!(got, Err(BandError::HaltNotNewer(_))),
                        "{}",
                        case["name"]
                    )
                }
                Some("HaltOutsideWindow") => {
                    assert!(
                        matches!(got, Err(BandError::HaltOutsideWindow(_))),
                        "{}",
                        case["name"]
                    )
                }
                _ => assert_eq!(got, Ok(halt(&case["expected"])), "{}", case["name"]),
            }
        }
    }

    #[test]
    fn lapses_match_the_reference() {
        for case in cases("lapsed") {
            let got = lapsed(halt(&case["halt"]), int(&case["now"]));
            assert_eq!(got, case["expected"].as_bool().unwrap(), "{}", case["name"]);
        }
    }

    #[test]
    fn activity_matches_the_reference() {
        for case in cases("active") {
            let got = active(halt(&case["halt"]), int(&case["now"]));
            assert_eq!(got, case["expected"].as_bool().unwrap(), "{}", case["name"]);
        }
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, ..ProptestConfig::default() })]

        #[test]
        fn an_accepted_message_is_newer_and_lapses_within_the_longest_window(
            stored_issued in any::<u64>(), stored_until in any::<u64>(), halted in any::<bool>(),
            issued_at in any::<u64>(), expires_at in any::<u64>(), now in any::<u64>(),
        ) {
            let stored = Halt { issued_at: stored_issued, until: stored_until };
            if let Ok(next) = accept(B256::ZERO, stored, halted, issued_at, expires_at, now) {
                prop_assert!(next.issued_at > stored.issued_at);
                prop_assert!(next.issued_at <= now);
                prop_assert!(next.until == 0 || (now < next.until && next.until - now <= MAX_HALT_S));
            }
        }
    }
}
