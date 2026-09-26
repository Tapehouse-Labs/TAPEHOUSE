#![cfg_attr(not(any(test, feature = "export-abi")), no_main)]
#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

pub mod chainlink;
pub mod error;
pub mod index;
pub mod multiplier;
pub mod quote;
pub mod redstone;
pub mod session;
pub mod token;

use alloc::vec::Vec;

use stylus_sdk::abi::Bytes;
use stylus_sdk::alloy_primitives::{Address, B256, U64, U128, U256, address};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::call::RawCall;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::{
    StorageAddress, StorageB256, StorageBool, StorageMap, StorageU64, StorageU128, StorageU256,
};

use crate::error::{
    AmbiguousLeg, BandError, DuplicateAsset, IncompleteStatus, IndexInUse, IndexWithoutChainlink,
    InvalidFeed, InvalidToken, LengthMismatch, NoLegs, NoToken, PackageNotNewer, ZeroSymbol,
};
use crate::index::Anchor;
use crate::multiplier::{Change, Status as CorporateAction};
use crate::session::{Session, Status};

const ECRECOVER: Address = address!("0x0000000000000000000000000000000000000001");

sol! {
    event PriceWritten(bytes32 indexed feedId, uint256 value, uint64 packageTimestampMs);
    event Anchored(bytes32 indexed symbol, uint64 chainlinkPrice, uint64 indexPrice, uint64 updatedAt);
    event MultiplierRecorded(bytes32 indexed symbol, uint128 before, uint128 after, uint64 effectiveAt);
    event MultiplierConfirmed(bytes32 indexed symbol, uint64 effectiveAt);
}

#[storage]
pub struct Price {
    value: StorageU256,
    package_timestamp_ms: StorageU64,
    written_at: StorageU64,
    sample_ms: StorageU64,
    tracked: StorageBool,
    index: StorageBool,
}

#[storage]
pub struct Volatility {
    var_cpb2: StorageU128,
    sample_px: StorageU64,
    step_at: StorageU64,
}

#[storage]
pub struct Asset {
    chainlink_feed: StorageAddress,
    redstone_feed_id: StorageB256,
    index_feed_id: StorageB256,
    anchor_cl_px: StorageU64,
    anchor_index_px: StorageU64,
    anchor_at: StorageU64,
    anchor_started_at: StorageU64,
    token: StorageAddress,
    effective_at: StorageU64,
    confirmed: StorageBool,
    multiplier_before: StorageU128,
    multiplier_after: StorageU128,
}

#[storage]
#[entrypoint]
pub struct Band {
    prices: StorageMap<B256, Price>,
    assets: StorageMap<B256, Asset>,
    volatility: StorageMap<B256, Volatility>,
    index_assets: StorageMap<B256, StorageB256>,
    close_ms: StorageU64,
}

#[public]
impl Band {
    /// Sets the per-asset configuration once. A zero feed or feed ID means the asset has no such leg.
    /// An asset takes its 24/7 leg from its own RedStone feed or from an index feed, never both; an
    /// index leg needs a Chainlink feed to anchor it, and backs one asset. A Stock Token turns the
    /// RedStone share price into the token's price through its ERC-8056 multiplier.
    #[constructor]
    pub fn constructor(
        &mut self,
        symbols: Vec<B256>,
        chainlink_feeds: Vec<Address>,
        redstone_feed_ids: Vec<B256>,
        index_feed_ids: Vec<B256>,
        tokens: Vec<Address>,
    ) -> Result<(), BandError> {
        let n = symbols.len();
        if chainlink_feeds.len() != n
            || redstone_feed_ids.len() != n
            || index_feed_ids.len() != n
            || tokens.len() != n
        {
            return Err(BandError::LengthMismatch(LengthMismatch {}));
        }
        for ((((symbol, feed), feed_id), index_id), token) in symbols
            .into_iter()
            .zip(chainlink_feeds)
            .zip(redstone_feed_ids)
            .zip(index_feed_ids)
            .zip(tokens)
        {
            if symbol == B256::ZERO {
                return Err(BandError::ZeroSymbol(ZeroSymbol {}));
            }
            if feed == Address::ZERO && feed_id == B256::ZERO && index_id == B256::ZERO {
                return Err(BandError::NoLegs(NoLegs { symbol }));
            }
            if feed_id != B256::ZERO && index_id != B256::ZERO {
                return Err(BandError::AmbiguousLeg(AmbiguousLeg { symbol }));
            }
            if index_id != B256::ZERO && feed == Address::ZERO {
                return Err(BandError::IndexWithoutChainlink(IndexWithoutChainlink {
                    symbol,
                }));
            }
            if self.asset(symbol) != (Address::ZERO, B256::ZERO, B256::ZERO, Address::ZERO) {
                return Err(BandError::DuplicateAsset(DuplicateAsset { symbol }));
            }
            if feed != Address::ZERO && !chainlink::has_expected_decimals(self.vm(), feed) {
                return Err(BandError::InvalidFeed(InvalidFeed { feed }));
            }
            let mut asset = self.assets.setter(symbol);
            asset.token.set(token);
            asset.chainlink_feed.set(feed);
            asset.redstone_feed_id.set(feed_id);
            asset.index_feed_id.set(index_id);
            if feed_id != B256::ZERO {
                self.prices.setter(feed_id).tracked.set(true);
            }
            if index_id != B256::ZERO {
                if self.index_assets.getter(index_id).get() != B256::ZERO {
                    return Err(BandError::IndexInUse(IndexInUse { feedId: index_id }));
                }
                self.index_assets.setter(index_id).set(symbol);
                let mut price = self.prices.setter(index_id);
                price.tracked.set(true);
                price.index.set(true);
            }
            if token != Address::ZERO {
                let (current, new, effective_at) = self
                    .token_terms(token, Change::default())
                    .filter(|(current, _, _)| *current != 0)
                    .ok_or(BandError::InvalidToken(InvalidToken { token }))?;
                let now = self.vm().block_timestamp();
                if let Some(change) =
                    multiplier::record(Change::default(), current, new, effective_at, now)
                {
                    self.record_change(symbol, change);
                }
            }
        }
        Ok(())
    }

    /// The Chainlink feed, RedStone feed ID, index feed ID and Stock Token of `symbol`. All zero for an
    /// unknown symbol.
    pub fn asset(&self, symbol: B256) -> (Address, B256, B256, Address) {
        let asset = self.assets.getter(symbol);
        (
            asset.chainlink_feed.get(),
            asset.redstone_feed_id.get(),
            asset.index_feed_id.get(),
            asset.token.get(),
        )
    }

    /// The token's multiplier change as it affects the band now: its status (0 none, 1 scheduled,
    /// 2 not yet confirmed by Chainlink), when it takes effect, and the multipliers before and after it.
    /// The multiplier before is zero when the change's size is not known. A token that cannot be read
    /// reports status 2 from now. All zero for an asset without a token.
    pub fn corporate_action(&self, symbol: B256) -> (u8, u64, u128, u128) {
        let (_, change) = self.terms_of(symbol);
        (
            multiplier::status(change, self.vm().block_timestamp()) as u8,
            change.effective_at,
            change.before,
            change.after,
        )
    }

    /// Confirms a material change past its step once a Chainlink round that started at or after the
    /// step falls inside the band of the 24/7 leg alone, then records the token's latest change. A
    /// change is not replaced before it is confirmed. Anyone may call it; the multiplier before a step
    /// can only be recorded before the step. Returns the status of the change in force.
    pub fn sync_multiplier(&mut self, symbol: B256) -> Result<u8, BandError> {
        let now = self.vm().block_timestamp();
        let token = self.assets.getter(symbol).token.get();
        if token == Address::ZERO {
            return Err(BandError::NoToken(NoToken { symbol }));
        }
        let mut change = self.change_of(symbol);
        let (current, new, effective_at) = self
            .token_terms(token, change)
            .ok_or(BandError::InvalidToken(InvalidToken { token }))?;
        loop {
            if multiplier::status(change, now) == CorporateAction::Unconfirmed {
                if !self.confirms(symbol, change) {
                    break;
                }
                change.confirmed = true;
                self.assets.setter(symbol).confirmed.set(true);
                self.vm().log(MultiplierConfirmed {
                    symbol,
                    effectiveAt: change.effective_at,
                });
            }
            match multiplier::record(change, current, new, effective_at, now) {
                Some(next) => {
                    change = next;
                    self.record_change(symbol, change);
                }
                None => break,
            }
        }
        Ok(multiplier::status(change, now) as u8)
    }

    /// The anchor of an asset priced from an index: the Chainlink print, the index price at that
    /// print, and the print's `updatedAt` in seconds. All zero before the first anchor.
    pub fn anchor(&self, symbol: B256) -> (u64, u64, u64) {
        let anchor = self.anchor_of(symbol);
        (anchor.cl_px, anchor.index_px, anchor.at_s)
    }

    /// Both legs of `symbol` as the band prices them: the Chainlink answer and its `updatedAt` in
    /// seconds, then the 24/7 price and its package timestamp in milliseconds. Prices are the token's:
    /// the RedStone share price times the token's multiplier, a Chainlink round from before a
    /// multiplier step scaled to the new terms, and for an asset priced from an index, the anchor's
    /// print moved by the index since. Zero for a leg that is unset, unreadable, above `u64`, or not
    /// known in the current terms. Never reverts.
    pub fn legs(&self, symbol: B256) -> (U256, u64, U256, u64) {
        let (legs, _) = self.priced_legs(symbol, self.vm().block_timestamp());
        (
            U256::from(legs.cl_px),
            legs.cl_at,
            U256::from(legs.px_247),
            legs.ms_247,
        )
    }

    /// The band of `symbol`: its state (0 halted, 1 degraded, 2 closed, 3 open), how many legs are
    /// live, its centre, half-width in basis points and bounds. While the session is not known,
    /// Chainlink sleeps and a band that is not halted is degraded. After a material multiplier step,
    /// the band is halted until Chainlink confirms the new terms, because the share price may or may
    /// not have moved with the step.
    pub fn quote(&self, symbol: B256) -> (u8, u8, u64, u64, u64, u128) {
        let now = self.vm().block_timestamp();
        let (inputs, change, session) = self.inputs(symbol, now);
        if multiplier::status(change, now) == CorporateAction::Unconfirmed {
            return (0, 0, 0, 0, 0, 0);
        }
        let mut quote = quote::compute(&inputs);
        if session.is_none() && quote.state != quote::State::Halted {
            quote.state = quote::State::Degraded;
        }
        (
            quote.state as u8,
            quote.live,
            quote.mid,
            quote.half_bps,
            quote.low,
            quote.high,
        )
    }

    /// Chainlink's 24/5 session from the signed New York market status: 0 not known, 1 closed,
    /// 2 open; NYSE's state and its next state (0 not known, 1 regular hours, 2 short close, 3 weekend
    /// or holiday); when NYSE changes state; and the session's next boundary, when an open session
    /// closes or a closed one reopens. Times in milliseconds, zero when not known.
    pub fn session(&self) -> (u8, u8, u8, u64, u64) {
        match self.session_at(self.vm().block_timestamp()) {
            Some(s) => (
                if s.open { 2 } else { 1 },
                s.nyse as u8,
                s.nyse_next.map_or(0, |n| n as u8),
                s.change_ms,
                s.boundary_ms,
            ),
            None => (0, 0, 0, 0, 0),
        }
    }

    /// Verifies a RedStone payload for `feed_ids` and stores each value. Anyone may call it. The three
    /// market-status feeds are written together or not at all.
    pub fn write_prices(
        &mut self,
        feed_ids: Vec<B256>,
        payload: Bytes,
    ) -> Result<Vec<U256>, BandError> {
        let status = session::FEEDS.map(|feed_id| feed_ids.contains(&feed_id));
        if status.contains(&true) && status.contains(&false) {
            return Err(BandError::IncompleteStatus(IncompleteStatus {}));
        }
        let verified = redstone::verify(
            &payload,
            &feed_ids,
            self.vm().block_timestamp(),
            |hash, v, r, s| ecrecover(self.vm(), hash, v, r, s),
        )?;
        self.store(&feed_ids, &verified)?;
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

    /// The variance of a configured asset's 24/7 feed, in centi-basis-points squared per minute. Zero
    /// before its second sample and for a feed no asset uses.
    pub fn variance(&self, feed_id: B256) -> u128 {
        self.volatility.getter(feed_id).var_cpb2.get().to::<u128>()
    }
}

struct Legs {
    cl_px: u64,
    cl_at: u64,
    px_247: u64,
    ms_247: u64,
}

impl Band {
    fn priced_legs(&self, symbol: B256, now: u64) -> (Legs, Change) {
        let asset = self.assets.getter(symbol);
        let (current, change) = self.terms_of(symbol);
        let (cl, started_at, updated_at) = chainlink::latest(self.vm(), asset.chainlink_feed.get());
        let cl_px =
            multiplier::chainlink_px(u64::try_from(cl).unwrap_or(0), started_at, change, now);
        let feed_id = asset.redstone_feed_id.get();
        let (px_247, ms_247) = if feed_id != B256::ZERO {
            let (value, ms, _) = self.price(feed_id);
            (
                u64::try_from(value).map_or(0, |v| multiplier::token_px(v, current)),
                ms,
            )
        } else {
            let index_id = asset.index_feed_id.get();
            if index_id == B256::ZERO {
                (0, 0)
            } else {
                let (value, ms, _) = self.price(index_id);
                let anchor = self.anchor_of(symbol);
                let anchor = Anchor {
                    cl_px: multiplier::chainlink_px(
                        anchor.cl_px,
                        asset.anchor_started_at.get().to::<u64>(),
                        change,
                        now,
                    ),
                    ..anchor
                };
                (
                    u64::try_from(value).map_or(0, |v| index::leg_px(anchor, v)),
                    ms,
                )
            }
        };
        let legs = Legs {
            cl_px,
            cl_at: if cl_px == 0 { 0 } else { updated_at },
            px_247,
            ms_247: if px_247 == 0 { 0 } else { ms_247 },
        };
        (legs, change)
    }

    fn inputs(&self, symbol: B256, now: u64) -> (quote::Inputs, Change, Option<Session>) {
        let (legs, change) = self.priced_legs(symbol, now);
        let asset = self.assets.getter(symbol);
        let feed_id = asset.redstone_feed_id.get();
        let (var_feed_id, basis_bps) = if feed_id == B256::ZERO {
            (asset.index_feed_id.get(), index::INDEX_BASIS_BPS)
        } else {
            (feed_id, 0)
        };
        let session = self.session_at(now);
        let inputs = quote::Inputs {
            live247_px: legs.px_247,
            live247_age_s: quote::package_age_s(now, legs.ms_247),
            cl_px: legs.cl_px,
            cl_age_s: quote::age_s(now, legs.cl_at),
            cl_session_open: session.is_some_and(|s| s.open),
            var_cpb2: self.variance(var_feed_id),
            basis_bps,
        };
        (inputs, change, session)
    }

    fn change_of(&self, symbol: B256) -> Change {
        let asset = self.assets.getter(symbol);
        Change {
            before: asset.multiplier_before.get().to::<u128>(),
            after: asset.multiplier_after.get().to::<u128>(),
            effective_at: asset.effective_at.get().to::<u64>(),
            confirmed: asset.confirmed.get(),
        }
    }

    fn token_terms(&self, token: Address, recorded: Change) -> Option<(u128, u128, u64)> {
        let (new, effective_at) = token::schedule(self.vm(), token)?;
        let now = self.vm().block_timestamp();
        let current = match multiplier::current(recorded, new, effective_at, now) {
            Some(current) => current,
            None => token::multiplier(self.vm(), token)?,
        };
        Some((current, new, effective_at))
    }

    fn terms_of(&self, symbol: B256) -> (u128, Change) {
        let token = self.assets.getter(symbol).token.get();
        if token == Address::ZERO {
            return (multiplier::SCALE, Change::default());
        }
        let recorded = self.change_of(symbol);
        let now = self.vm().block_timestamp();
        match self.token_terms(token, recorded) {
            Some((current, new, effective_at)) => (
                current,
                multiplier::observe(recorded, current, new, effective_at, now),
            ),
            None => (
                0,
                Change {
                    effective_at: now,
                    ..Change::default()
                },
            ),
        }
    }

    fn record_change(&mut self, symbol: B256, change: Change) {
        let mut asset = self.assets.setter(symbol);
        asset.multiplier_before.set(U128::from(change.before));
        asset.multiplier_after.set(U128::from(change.after));
        asset.effective_at.set(U64::from(change.effective_at));
        asset.confirmed.set(change.confirmed);
        let feed_id = asset.redstone_feed_id.get();
        if feed_id != B256::ZERO {
            self.volatility
                .setter(feed_id)
                .step_at
                .set(U64::from(change.effective_at));
        }
        self.vm().log(MultiplierRecorded {
            symbol,
            before: change.before,
            after: change.after,
            effectiveAt: change.effective_at,
        });
    }

    fn confirms(&self, symbol: B256, change: Change) -> bool {
        let feed = self.assets.getter(symbol).chainlink_feed.get();
        if feed == Address::ZERO {
            return true;
        }
        let (_, started_at, _) = chainlink::latest(self.vm(), feed);
        let (inputs, _, _) = self.inputs(symbol, self.vm().block_timestamp());
        let alone = quote::compute(&quote::Inputs {
            cl_px: 0,
            cl_session_open: false,
            ..inputs
        });
        multiplier::confirms(
            change,
            started_at,
            inputs.cl_px,
            alone.mid,
            alone.low,
            alone.high,
        )
    }

    fn store(&mut self, feed_ids: &[B256], verified: &redstone::Verified) -> Result<(), BandError> {
        let now = self.vm().block_timestamp();
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
            self.record_sample(*feed_id, *value, verified.timestamp_ms);
            if self.prices.getter(*feed_id).index.get() {
                self.record_anchor(*feed_id, *value, verified.timestamp_ms);
            }
            self.vm().log(PriceWritten {
                feedId: *feed_id,
                value: *value,
                packageTimestampMs: verified.timestamp_ms,
            });
        }
        if feed_ids.contains(&session::CURRENT_STATUS) {
            self.record_close();
        }
        Ok(())
    }

    fn status(&self) -> Option<Status> {
        let [current, next, change] = session::FEEDS.map(|feed_id| self.price(feed_id));
        Status::decode([current.0, next.0, change.0], [current.1, next.1, change.1])
    }

    fn record_close(&mut self) {
        let close_ms = self.close_ms.get().to::<u64>();
        let next = session::next_close(self.status(), close_ms);
        if next != close_ms {
            self.close_ms.set(U64::from(next));
        }
    }

    fn session_at(&self, now: u64) -> Option<Session> {
        session::derive(
            self.status(),
            self.close_ms.get().to::<u64>(),
            now.saturating_mul(1000),
        )
    }

    fn anchor_of(&self, symbol: B256) -> Anchor {
        let asset = self.assets.getter(symbol);
        Anchor {
            cl_px: asset.anchor_cl_px.get().to::<u64>(),
            index_px: asset.anchor_index_px.get().to::<u64>(),
            at_s: asset.anchor_at.get().to::<u64>(),
        }
    }

    fn record_anchor(&mut self, index_id: B256, value: U256, package_timestamp_ms: u64) {
        let symbol = self.index_assets.getter(index_id).get();
        let feed = self.assets.getter(symbol).chainlink_feed.get();
        let (cl_px, cl_started_at, cl_at) = chainlink::latest(self.vm(), feed);
        let (Ok(cl_px), Ok(index_px)) = (u64::try_from(cl_px), u64::try_from(value)) else {
            return;
        };
        let anchor = self.anchor_of(symbol);
        if let Some(next) = index::next_anchor(anchor, cl_px, cl_at, index_px, package_timestamp_ms)
        {
            let mut asset = self.assets.setter(symbol);
            asset.anchor_cl_px.set(U64::from(next.cl_px));
            asset.anchor_index_px.set(U64::from(next.index_px));
            asset.anchor_at.set(U64::from(next.at_s));
            asset.anchor_started_at.set(U64::from(cl_started_at));
            self.vm().log(Anchored {
                symbol,
                chainlinkPrice: next.cl_px,
                indexPrice: next.index_px,
                updatedAt: next.at_s,
            });
        }
    }

    fn record_sample(&mut self, feed_id: B256, value: U256, package_timestamp_ms: u64) {
        let price = self.prices.getter(feed_id);
        let sample_ms = price.sample_ms.get().to::<u64>();
        let Ok(px) = u64::try_from(value) else {
            return;
        };
        if !price.tracked.get() || !quote::sample_due(sample_ms, package_timestamp_ms) {
            return;
        }
        let volatility = self.volatility.getter(feed_id);
        let step_ms = volatility.step_at.get().to::<u64>().saturating_mul(1000);
        let restart = sample_ms < step_ms && step_ms <= package_timestamp_ms;
        let last = quote::Sample {
            var_cpb2: volatility.var_cpb2.get().to::<u128>(),
            px: if restart {
                0
            } else {
                volatility.sample_px.get().to::<u64>()
            },
            ms: if restart { 0 } else { sample_ms },
        };
        if let Some(next) = quote::next_sample(last, px, package_timestamp_ms) {
            let mut volatility = self.volatility.setter(feed_id);
            volatility.var_cpb2.set(U128::from(next.var_cpb2));
            volatility.sample_px.set(U64::from(next.px));
            self.prices
                .setter(feed_id)
                .sample_ms
                .set(U64::from(next.ms));
        }
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
    use stylus_sdk::alloy_sol_types::{SolEvent, SolValue};
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
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![],
            vec![symbol("NVDA---24_7")],
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
        assert_eq!(result, Err(BandError::LengthMismatch(LengthMismatch {})));
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![Address::ZERO],
            vec![],
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
        assert_eq!(result, Err(BandError::LengthMismatch(LengthMismatch {})));
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![Address::ZERO],
            vec![symbol("NVDA---24_7")],
            vec![],
            vec![Address::ZERO],
        );
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
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
        assert_eq!(result, Err(BandError::ZeroSymbol(ZeroSymbol {})));
    }

    #[test]
    fn constructor_rejects_an_asset_without_legs() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![symbol("SPY")],
            vec![Address::ZERO],
            vec![B256::ZERO],
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
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
            vec![B256::ZERO, B256::ZERO],
            vec![Address::ZERO; 2],
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
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![FEED],
            vec![B256::ZERO],
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
        assert_eq!(
            result,
            Err(BandError::InvalidFeed(InvalidFeed { feed: FEED }))
        );
    }

    #[test]
    fn constructor_rejects_a_feed_without_code() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![FEED],
            vec![B256::ZERO],
            vec![B256::ZERO],
            vec![Address::ZERO],
        );
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
            vec![symbol("NVDA"), symbol("TSLA"), symbol("SPY")],
            vec![FEED, Address::ZERO, FEED],
            vec![symbol("NVDA---24_7"), symbol("TSLA---24_7"), B256::ZERO],
            vec![B256::ZERO, B256::ZERO, symbol("USA500.Y---24_7")],
            vec![Address::ZERO; 3],
        )
        .unwrap();
        assert_eq!(
            band.asset(symbol("NVDA")),
            (FEED, symbol("NVDA---24_7"), B256::ZERO, Address::ZERO)
        );
        assert_eq!(
            band.asset(symbol("TSLA")),
            (
                Address::ZERO,
                symbol("TSLA---24_7"),
                B256::ZERO,
                Address::ZERO
            )
        );
        assert_eq!(
            band.asset(symbol("SPY")),
            (FEED, B256::ZERO, symbol("USA500.Y---24_7"), Address::ZERO)
        );
        assert_eq!(
            band.asset(symbol("AAPL")),
            (Address::ZERO, B256::ZERO, B256::ZERO, Address::ZERO)
        );
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
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0, 0));
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(-1, 1_790_089_676)),
        );
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0, 0));
    }

    #[test]
    fn a_reverting_feed_is_no_reading() {
        let vm = TestVM::default();
        vm.mock_static_call(FEED, LATEST_ROUND_DATA_CALL.to_vec(), Err(vec![]));
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0, 0));
    }

    #[test]
    fn a_feed_without_code_is_no_reading() {
        let vm = TestVM::default();
        assert_eq!(chainlink::latest(&vm, FEED), (U256::ZERO, 0, 0));
    }

    #[test]
    fn the_constructor_tracks_each_24_7_feed() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(
            FEED,
            DECIMALS_CALL.to_vec(),
            Ok(U256::from(8u8).to_be_bytes::<32>().to_vec()),
        );
        band.constructor(
            vec![symbol("NVDA"), symbol("SPY")],
            vec![Address::ZERO, FEED],
            vec![symbol("NVDA---24_7"), B256::ZERO],
            vec![B256::ZERO, symbol("USA500.Y---24_7")],
            vec![Address::ZERO; 2],
        )
        .unwrap();
        assert!(band.prices.getter(symbol("NVDA---24_7")).tracked.get());
        assert!(!band.prices.getter(symbol("NVDA---24_7")).index.get());
        let index = band.prices.getter(symbol("USA500.Y---24_7"));
        assert!(index.tracked.get() && index.index.get());
        assert_eq!(
            band.index_assets.getter(symbol("USA500.Y---24_7")).get(),
            symbol("SPY")
        );
        assert!(!band.prices.getter(B256::ZERO).tracked.get());
        assert!(!band.prices.getter(symbol("NY_MARKET_STATUS")).tracked.get());
    }

    #[test]
    fn samples_fold_into_the_variance_of_a_tracked_feed() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let feed_id = symbol("NVDA---24_7");
        band.prices.setter(feed_id).tracked.set(true);
        band.record_sample(feed_id, U256::from(22_000_000_000u64), 1_790_000_000_000);
        assert_eq!(band.variance(feed_id), 0);
        band.record_sample(feed_id, U256::from(22_300_000_000u64), 1_790_000_049_999);
        assert_eq!(band.variance(feed_id), 0);
        band.record_sample(feed_id, U256::from(22_300_000_000u64), 1_790_000_050_000);
        assert_eq!(
            band.variance(feed_id),
            quote::ewma_update(0, 22_000_000_000, 22_300_000_000, 50)
        );
        assert_eq!(
            band.prices.getter(feed_id).sample_ms.get().to::<u64>(),
            1_790_000_050_000
        );
        assert_eq!(
            band.volatility.getter(feed_id).sample_px.get().to::<u64>(),
            22_300_000_000
        );
    }

    #[test]
    fn an_untracked_feed_keeps_no_variance() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let feed_id = symbol("NY_MARKET_STATUS");
        band.record_sample(feed_id, U256::from(101_000_000u64), 1_790_000_000_000);
        band.record_sample(feed_id, U256::from(100_000_000u64), 1_790_000_060_000);
        assert_eq!(band.variance(feed_id), 0);
        assert_eq!(band.prices.getter(feed_id).sample_ms.get().to::<u64>(), 0);
    }

    #[test]
    fn a_value_beyond_u64_is_not_a_sample() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let feed_id = symbol("NY_MARKET_NEXT_CHANGE_TIME");
        band.prices.setter(feed_id).tracked.set(true);
        band.record_sample(
            feed_id,
            U256::from(179_017_020_000_000_000_000u128),
            1_790_000_000_000,
        );
        assert_eq!(band.prices.getter(feed_id).sample_ms.get().to::<u64>(), 0);
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
    #[test]
    fn constructor_rejects_an_index_next_to_a_24_7_feed() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![symbol("SPY")],
            vec![FEED],
            vec![symbol("SPY---24_7")],
            vec![symbol("USA500.Y---24_7")],
            vec![Address::ZERO],
        );
        assert_eq!(
            result,
            Err(BandError::AmbiguousLeg(AmbiguousLeg {
                symbol: symbol("SPY")
            }))
        );
    }

    #[test]
    fn constructor_rejects_an_index_without_chainlink() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let result = band.constructor(
            vec![symbol("SPY")],
            vec![Address::ZERO],
            vec![B256::ZERO],
            vec![symbol("USA500.Y---24_7")],
            vec![Address::ZERO],
        );
        assert_eq!(
            result,
            Err(BandError::IndexWithoutChainlink(IndexWithoutChainlink {
                symbol: symbol("SPY")
            }))
        );
    }

    #[test]
    fn constructor_rejects_an_index_that_backs_another_asset() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(
            FEED,
            DECIMALS_CALL.to_vec(),
            Ok(U256::from(8u8).to_be_bytes::<32>().to_vec()),
        );
        let result = band.constructor(
            vec![symbol("SPY"), symbol("VOO")],
            vec![FEED, FEED],
            vec![B256::ZERO, B256::ZERO],
            vec![symbol("USA500.Y---24_7"), symbol("USA500.Y---24_7")],
            vec![Address::ZERO; 2],
        );
        assert_eq!(
            result,
            Err(BandError::IndexInUse(IndexInUse {
                feedId: symbol("USA500.Y---24_7")
            }))
        );
    }

    const FRI_CLOSE_S: u64 = 1_789_761_600;
    const SAT_NOON_S: u64 = 1_789_819_200;
    const MON_OPEN_S: u64 = 1_789_997_400;
    const WED_3AM_S: u64 = 1_790_132_400;
    const WED_OPEN_S: u64 = 1_790_170_200;
    const REGULAR: u64 = 100_000_000;
    const CLOSED_SHORT: u64 = 101_000_000;
    const CLOSED_LONG: u64 = 102_000_000;

    fn store(band: &mut Band, feed_id: B256, value: U256, package_timestamp_ms: u64) {
        let mut price = band.prices.setter(feed_id);
        price.value.set(value);
        price
            .package_timestamp_ms
            .set(U64::from(package_timestamp_ms));
    }

    fn sign_status(band: &mut Band, current: u64, next: u64, change_s: u64, package_s: u64) {
        let ms = package_s * 1000;
        store(band, session::CURRENT_STATUS, U256::from(current), ms);
        store(band, session::NEXT_STATUS, U256::from(next), ms);
        let change = U256::from(change_s * 1000) * U256::from(session::SCALE);
        store(band, session::NEXT_CHANGE_TIME, change, ms);
    }

    fn nvda(vm: &TestVM, band: &mut Band, cl_px: i64, cl_at: u64, px_247: u64, at_247: u64) {
        let mut asset = band.assets.setter(symbol("NVDA"));
        asset.chainlink_feed.set(FEED);
        asset.redstone_feed_id.set(symbol("NVDA---24_7"));
        store(
            band,
            symbol("NVDA---24_7"),
            U256::from(px_247),
            at_247 * 1000,
        );
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(cl_px, cl_at)),
        );
    }

    #[test]
    fn a_regular_package_records_its_close() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(
            &mut band,
            REGULAR,
            CLOSED_LONG,
            FRI_CLOSE_S,
            FRI_CLOSE_S - 60,
        );
        band.record_close();
        assert_eq!(band.close_ms.get().to::<u64>(), FRI_CLOSE_S * 1000);
        sign_status(
            &mut band,
            CLOSED_LONG,
            REGULAR,
            MON_OPEN_S,
            FRI_CLOSE_S + 10,
        );
        band.record_close();
        assert_eq!(band.close_ms.get().to::<u64>(), FRI_CLOSE_S * 1000);
    }

    #[test]
    fn the_session_is_not_known_without_a_status() {
        let vm = TestVM::default();
        let band = Band::from(&vm);
        vm.set_block_timestamp(WED_3AM_S);
        assert_eq!(band.session(), (0, 0, 0, 0, 0));
    }

    #[test]
    fn the_session_reports_the_weekend_and_its_reopen() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(
            &mut band,
            REGULAR,
            CLOSED_LONG,
            FRI_CLOSE_S,
            FRI_CLOSE_S - 60,
        );
        band.record_close();
        sign_status(&mut band, CLOSED_LONG, REGULAR, MON_OPEN_S, SAT_NOON_S - 10);
        vm.set_block_timestamp(SAT_NOON_S);
        assert_eq!(
            band.session(),
            (1, 3, 1, MON_OPEN_S * 1000, (MON_OPEN_S - 48_600) * 1000)
        );
        vm.set_block_timestamp(FRI_CLOSE_S + 14_399);
        sign_status(
            &mut band,
            CLOSED_LONG,
            REGULAR,
            MON_OPEN_S,
            FRI_CLOSE_S + 14_390,
        );
        assert_eq!(
            band.session(),
            (2, 3, 1, MON_OPEN_S * 1000, (FRI_CLOSE_S + 14_400) * 1000)
        );
    }

    #[test]
    fn an_open_session_with_both_legs_is_open() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_SHORT, REGULAR, WED_OPEN_S, WED_3AM_S - 10);
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            WED_3AM_S - 300,
            22_000_000_000,
            WED_3AM_S - 10,
        );
        vm.set_block_timestamp(WED_3AM_S);
        let quote = quote::compute(&quote::Inputs {
            live247_px: 22_000_000_000,
            live247_age_s: 10,
            cl_px: 22_050_000_000,
            cl_age_s: 300,
            cl_session_open: true,
            var_cpb2: 0,
            basis_bps: 0,
        });
        assert_eq!(quote.state, quote::State::Open);
        assert_eq!(
            band.quote(symbol("NVDA")),
            (3, 2, quote.mid, quote.half_bps, quote.low, quote.high)
        );
    }

    #[test]
    fn a_closed_session_puts_chainlink_to_sleep() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_LONG, REGULAR, MON_OPEN_S, SAT_NOON_S - 10);
        nvda(
            &vm,
            &mut band,
            30_000_000_000,
            SAT_NOON_S - 300,
            22_000_000_000,
            SAT_NOON_S - 10,
        );
        vm.set_block_timestamp(SAT_NOON_S);
        assert_eq!(
            band.quote(symbol("NVDA")),
            (2, 1, 22_000_000_000, 55, 21_879_000_000, 22_121_000_000)
        );
    }

    #[test]
    fn a_status_gone_quiet_degrades_the_band() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(
            &mut band,
            CLOSED_SHORT,
            REGULAR,
            WED_OPEN_S,
            WED_3AM_S - 3_600,
        );
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            WED_3AM_S - 300,
            22_000_000_000,
            WED_3AM_S - 10,
        );
        vm.set_block_timestamp(WED_3AM_S);
        assert_eq!(band.quote(symbol("NVDA")).0, 3);
        vm.set_block_timestamp(WED_3AM_S + 1);
        assert_eq!(
            band.quote(symbol("NVDA")),
            (1, 1, 22_000_000_000, 55, 21_879_000_000, 22_121_000_000)
        );
        assert_eq!(band.session(), (0, 0, 0, 0, 0));
    }

    #[test]
    fn a_24_7_leg_gone_quiet_degrades_the_band() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_SHORT, REGULAR, WED_OPEN_S, WED_3AM_S - 10);
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            WED_3AM_S - 300,
            22_000_000_000,
            WED_3AM_S - 121,
        );
        vm.set_block_timestamp(WED_3AM_S);
        let (state, live, mid, half, _, _) = band.quote(symbol("NVDA"));
        assert_eq!((state, live, mid, half), (1, 1, 22_050_000_000, 105));
    }

    #[test]
    fn chainlink_gone_quiet_degrades_the_band() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_SHORT, REGULAR, WED_OPEN_S, WED_3AM_S - 10);
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            WED_3AM_S - 86_461,
            22_000_000_000,
            WED_3AM_S - 10,
        );
        vm.set_block_timestamp(WED_3AM_S);
        let (state, live, mid, half, _, _) = band.quote(symbol("NVDA"));
        assert_eq!((state, live, mid, half), (1, 1, 22_000_000_000, 55));
    }

    #[test]
    fn no_live_leg_halts_the_band() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_LONG, REGULAR, MON_OPEN_S, SAT_NOON_S - 10);
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            SAT_NOON_S - 300,
            22_000_000_000,
            SAT_NOON_S - 121,
        );
        vm.set_block_timestamp(SAT_NOON_S);
        assert_eq!(band.quote(symbol("NVDA")), (0, 0, 0, 0, 0, 0));
        assert_eq!(band.quote(symbol("AAPL")), (0, 0, 0, 0, 0, 0));
    }

    #[test]
    fn a_leg_above_u64_is_absent() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        sign_status(&mut band, CLOSED_SHORT, REGULAR, WED_OPEN_S, WED_3AM_S - 10);
        nvda(
            &vm,
            &mut band,
            22_050_000_000,
            WED_3AM_S - 300,
            0,
            WED_3AM_S - 10,
        );
        store(
            &mut band,
            symbol("NVDA---24_7"),
            U256::from(u64::MAX) + U256::from(1u8),
            (WED_3AM_S - 10) * 1000,
        );
        vm.set_block_timestamp(WED_3AM_S);
        let (state, live, mid, _, _, _) = band.quote(symbol("NVDA"));
        assert_eq!((state, live, mid), (1, 1, 22_050_000_000));
    }

    fn spy(vm: &TestVM, band: &mut Band, cl_px: i64, cl_at: u64) {
        let index = symbol("USA500.Y---24_7");
        let mut asset = band.assets.setter(symbol("SPY"));
        asset.chainlink_feed.set(FEED);
        asset.index_feed_id.set(index);
        band.index_assets.setter(index).set(symbol("SPY"));
        vm.mock_static_call(
            FEED,
            LATEST_ROUND_DATA_CALL.to_vec(),
            Ok(round(cl_px, cl_at)),
        );
    }

    #[test]
    fn an_index_price_near_a_new_print_anchors_it() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let index = symbol("USA500.Y---24_7");
        spy(&vm, &mut band, 77_232_802_713, 1_790_352_180);
        band.record_anchor(index, U256::from(774_263_746_783u64), 1_790_352_300_001);
        assert_eq!(band.anchor(symbol("SPY")), (0, 0, 0));
        band.record_anchor(index, U256::from(774_263_746_783u64), 1_790_352_200_000);
        assert_eq!(
            band.anchor(symbol("SPY")),
            (77_232_802_713, 774_263_746_783, 1_790_352_180)
        );
        let logs = vm.get_emitted_logs();
        assert_eq!(logs.len(), 1);
        assert_eq!(logs[0].0, vec![Anchored::SIGNATURE_HASH, symbol("SPY")]);
        band.record_anchor(index, U256::from(774_135_000_000u64), 1_790_352_210_000);
        assert_eq!(vm.get_emitted_logs().len(), 1);
        assert_eq!(
            band.anchor(symbol("SPY")),
            (77_232_802_713, 774_263_746_783, 1_790_352_180)
        );
    }

    #[test]
    fn writing_the_index_anchors_its_asset() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let index = symbol("USA500.Y---24_7");
        spy(&vm, &mut band, 77_232_802_713, 1_790_352_180);
        band.prices.setter(index).index.set(true);
        vm.set_block_timestamp(1_790_352_205);
        let verified = redstone::Verified {
            values: vec![U256::from(774_263_746_783u64)],
            timestamp_ms: 1_790_352_200_000,
        };
        band.store(&[index], &verified).unwrap();
        assert_eq!(
            band.anchor(symbol("SPY")),
            (77_232_802_713, 774_263_746_783, 1_790_352_180)
        );
        assert_eq!(
            band.assets.getter(symbol("SPY")).anchor_started_at.get(),
            U64::from(1_790_352_168u64)
        );
    }

    #[test]
    fn the_index_leg_prices_spy() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let index = symbol("USA500.Y---24_7");
        spy(&vm, &mut band, 77_232_802_713, 1_790_352_180);
        store(
            &mut band,
            index,
            U256::from(774_263_746_783u64),
            1_790_352_200_000,
        );
        assert_eq!(
            band.legs(symbol("SPY")),
            (U256::from(77_232_802_713u64), 1_790_352_180, U256::ZERO, 0)
        );
        band.record_anchor(index, U256::from(774_263_746_783u64), 1_790_352_200_000);
        store(
            &mut band,
            index,
            U256::from(774_135_000_000u64),
            1_790_370_020_000,
        );
        band.volatility
            .setter(index)
            .var_cpb2
            .set(U128::from(9_000_000u64));
        sign_status(
            &mut band,
            CLOSED_LONG,
            REGULAR,
            1_790_602_200,
            1_790_370_020,
        );
        band.close_ms.set(U64::from(1_790_366_400_000u64));
        vm.set_block_timestamp(1_790_370_033);
        assert_eq!(
            band.legs(symbol("SPY")),
            (
                U256::from(77_232_802_713u64),
                1_790_352_180,
                U256::from(77_219_960_222u64),
                1_790_370_020_000
            )
        );
        let expected = quote::compute(&quote::Inputs {
            live247_px: 77_219_960_222,
            live247_age_s: 13,
            cl_px: 77_232_802_713,
            cl_age_s: 17_853,
            cl_session_open: true,
            var_cpb2: 9_000_000,
            basis_bps: index::INDEX_BASIS_BPS,
        });
        assert_eq!(expected.half_bps, 180);
        assert_eq!(
            band.quote(symbol("SPY")),
            (
                3,
                2,
                expected.mid,
                expected.half_bps,
                expected.low,
                expected.high
            )
        );
    }
    #[test]
    fn a_status_written_in_regular_hours_records_its_close() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let change = U256::from(FRI_CLOSE_S * 1000) * U256::from(session::SCALE);
        let verified = redstone::Verified {
            values: vec![U256::from(REGULAR), U256::from(CLOSED_LONG), change],
            timestamp_ms: (FRI_CLOSE_S - 600) * 1000,
        };
        vm.set_block_timestamp(FRI_CLOSE_S - 590);
        band.store(&session::FEEDS, &verified).unwrap();
        assert_eq!(band.close_ms.get().to::<u64>(), FRI_CLOSE_S * 1000);
        assert_eq!(
            band.session(),
            (2, 1, 3, FRI_CLOSE_S * 1000, (FRI_CLOSE_S + 14_400) * 1000)
        );
        assert_eq!(vm.get_emitted_logs().len(), 3);
    }

    #[test]
    fn a_status_written_in_part_reverts() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        for feed_ids in [
            vec![session::CURRENT_STATUS],
            vec![session::NEXT_STATUS, session::NEXT_CHANGE_TIME],
            vec![symbol("NVDA---24_7"), session::NEXT_CHANGE_TIME],
        ] {
            assert_eq!(
                band.write_prices(feed_ids, Bytes::from(vec![0u8; 32])),
                Err(BandError::IncompleteStatus(IncompleteStatus {}))
            );
        }
        assert_eq!(
            band.write_prices(session::FEEDS.to_vec(), Bytes::from(vec![0u8; 32])),
            Err(BandError::CalldataMustHaveValidPayload(
                CalldataMustHaveValidPayload {}
            ))
        );
    }
    const TOKEN: Address = address!("0x5555555555555555555555555555555555555555");

    #[test]
    fn constructor_rejects_a_token_without_a_multiplier() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(TOKEN, vec![0xa6, 0x0b, 0xf1, 0x3d], Err(vec![]));
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![Address::ZERO],
            vec![symbol("NVDA---24_7")],
            vec![B256::ZERO],
            vec![TOKEN],
        );
        assert_eq!(
            result,
            Err(BandError::InvalidToken(InvalidToken { token: TOKEN }))
        );
    }

    #[test]
    fn constructor_rejects_a_token_whose_multiplier_is_zero() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        vm.mock_static_call(
            TOKEN,
            vec![0xdc, 0x76, 0x70, 0x07],
            Ok(U256::ZERO.to_be_bytes::<32>().to_vec()),
        );
        let result = band.constructor(
            vec![symbol("NVDA")],
            vec![Address::ZERO],
            vec![symbol("NVDA---24_7")],
            vec![B256::ZERO],
            vec![TOKEN],
        );
        assert_eq!(
            result,
            Err(BandError::InvalidToken(InvalidToken { token: TOKEN }))
        );
    }

    #[test]
    fn a_token_that_cannot_be_read_halts_its_asset() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        band.assets.setter(symbol("NVDA")).token.set(TOKEN);
        vm.mock_static_call(TOKEN, vec![0xdc, 0x76, 0x70, 0x07], Err(vec![]));
        vm.set_block_timestamp(WED_3AM_S);
        assert_eq!(band.corporate_action(symbol("NVDA")), (2, WED_3AM_S, 0, 0));
        assert_eq!(band.quote(symbol("NVDA")), (0, 0, 0, 0, 0, 0));
    }

    #[test]
    fn syncing_an_asset_without_a_token_reverts() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        assert_eq!(
            band.sync_multiplier(symbol("NVDA")),
            Err(BandError::NoToken(NoToken {
                symbol: symbol("NVDA")
            }))
        );
    }

    #[test]
    fn a_sample_across_a_multiplier_step_restarts() {
        let vm = TestVM::default();
        let mut band = Band::from(&vm);
        let feed_id = symbol("NVDA---24_7");
        band.prices.setter(feed_id).tracked.set(true);
        band.record_sample(feed_id, U256::from(22_000_000_000u64), 1_790_000_000_000);
        band.record_sample(feed_id, U256::from(22_300_000_000u64), 1_790_000_060_000);
        let variance = band.variance(feed_id);
        assert!(variance > 0);
        band.volatility
            .setter(feed_id)
            .step_at
            .set(U64::from(1_790_000_100u64));
        band.record_sample(feed_id, U256::from(5_575_000_000u64), 1_790_000_120_000);
        assert_eq!(band.variance(feed_id), variance);
        assert_eq!(
            band.volatility.getter(feed_id).sample_px.get().to::<u64>(),
            5_575_000_000
        );
        band.record_sample(feed_id, U256::from(5_575_000_000u64), 1_790_000_180_000);
        assert_eq!(
            band.variance(feed_id),
            quote::ewma_update(variance, 5_575_000_000, 5_575_000_000, 60)
        );
    }
}
