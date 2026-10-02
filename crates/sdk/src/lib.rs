// SPDX-License-Identifier: MIT OR Apache-2.0
//! Typed reads and transactions for Tapehouse's price band, its feeds, the margin accounts and the short positions,
//! over alloy bindings generated from the contracts' ABIs. Addresses come from a chain's registry,
//! `deployments/<chainId>.json`, read at runtime.
#![warn(missing_docs)]

pub mod bindings;
mod client;
mod deployments;
mod revert;

pub use client::{BuyBack, PackageSource, Repayment, Sale, Tapehouse};
pub use deployments::{
    CROSS, Deployments, SHARE_PRICE_CHAINS, SharePriceFeed, TokenPriceFeed, to_bytes32,
};
pub use revert::{Revert, decode_revert, decode_revert_data};

/// What the SDK fails with.
#[derive(Debug)]
#[non_exhaustive]
pub enum Error {
    /// The registry is malformed or lacks an entry.
    Registry(String),
    /// An argument is out of range: a slippage above 10000 basis points, or a name longer than 32 bytes.
    Argument(String),
    /// A call reverted with an error the SDK decodes.
    Revert(Revert),
    /// A call failed otherwise.
    Contract(alloy::contract::Error),
    /// The package source failed.
    Source(Box<dyn std::error::Error + Send + Sync>),
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::Registry(message) | Self::Argument(message) => f.write_str(message),
            Self::Revert(revert) => write!(f, "reverted with {revert}"),
            Self::Contract(error) => error.fmt(f),
            Self::Source(error) => write!(f, "the package source failed: {error}"),
        }
    }
}

impl std::error::Error for Error {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        match self {
            Self::Contract(error) => Some(error),
            Self::Source(error) => Some(error.as_ref()),
            Self::Registry(_) | Self::Argument(_) | Self::Revert(_) => None,
        }
    }
}

impl From<alloy::contract::Error> for Error {
    fn from(error: alloy::contract::Error) -> Self {
        decode_revert(&error).map_or(Self::Contract(error), Self::Revert)
    }
}

/// The SDK's result.
pub type Result<T> = std::result::Result<T, Error>;
