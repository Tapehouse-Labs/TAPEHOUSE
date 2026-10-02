// SPDX-License-Identifier: MIT OR Apache-2.0
use std::collections::HashMap;
use std::fmt;
use std::sync::LazyLock;

use alloy::dyn_abi::{DynSolValue, JsonAbiExt};
use alloy::json_abi::{Error as AbiError, JsonAbi};
use alloy::primitives::{Bytes, Selector};
use alloy::sol_types::{Panic, Revert as SolidityRevert, SolError};

/// A call's revert, decoded: a custom error of Tapehouse's contracts, its baskets, the band, the Stock Tokens or USDG, or
/// Solidity's `Error(string)` and `Panic(uint256)`.
#[derive(Clone, Debug, PartialEq)]
pub struct Revert {
    /// The error's name: `Error` and `Panic` for Solidity's own.
    pub name: String,
    /// Its arguments, decoded.
    pub args: Vec<DynSolValue>,
    /// The raw revert data.
    pub data: Bytes,
}

impl fmt::Display for Revert {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let args: Vec<String> = self.args.iter().map(format_value).collect();
        write!(f, "{}({})", self.name, args.join(", "))
    }
}

static ERRORS: LazyLock<HashMap<Selector, AbiError>> = LazyLock::new(|| {
    let mut errors = HashMap::new();
    for json in [
        include_str!("../abi/ShortPositions.json"),
        include_str!("../abi/MarginAccounts.json"),
        include_str!("../abi/Basket.json"),
        include_str!("../abi/BandFeed.json"),
        include_str!("../abi/MorphoBandOracle.json"),
        include_str!("../abi/Band.json"),
        include_str!("../abi/StockToken.json"),
        include_str!("../abi/Usdg.json"),
    ] {
        let abi: JsonAbi = serde_json::from_str(json).expect("the bindings' ABIs parse");
        for error in abi.errors() {
            errors
                .entry(error.selector())
                .or_insert_with(|| error.clone());
        }
    }
    errors
});

/// Decodes the revert behind a failed call or transaction of an alloy contract instance.
pub fn decode_revert(error: &alloy::contract::Error) -> Option<Revert> {
    decode_revert_data(&error.as_revert_data()?)
}

/// Decodes raw revert data.
pub fn decode_revert_data(data: &[u8]) -> Option<Revert> {
    let selector = Selector::try_from(data.get(..4)?).ok()?;
    let (name, args) = if selector == SolidityRevert::SELECTOR {
        let reason = SolidityRevert::abi_decode(data).ok()?.reason;
        ("Error".to_string(), vec![DynSolValue::String(reason)])
    } else if selector == Panic::SELECTOR {
        let code = Panic::abi_decode(data).ok()?.code;
        ("Panic".to_string(), vec![DynSolValue::Uint(code, 256)])
    } else {
        let error = ERRORS.get(&selector)?;
        (error.name.clone(), error.abi_decode_input(&data[4..]).ok()?)
    };
    Some(Revert {
        name,
        args,
        data: Bytes::copy_from_slice(data),
    })
}

fn format_value(value: &DynSolValue) -> String {
    match value {
        DynSolValue::Address(address) => address.to_checksum(None),
        DynSolValue::Bool(value) => value.to_string(),
        DynSolValue::Int(value, _) => value.to_string(),
        DynSolValue::Uint(value, _) => value.to_string(),
        DynSolValue::FixedBytes(word, size) => Bytes::copy_from_slice(&word[..*size]).to_string(),
        DynSolValue::Bytes(bytes) => Bytes::copy_from_slice(bytes).to_string(),
        DynSolValue::String(text) => text.clone(),
        other => format!("{other:?}"),
    }
}
