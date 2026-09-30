// SPDX-License-Identifier: MIT OR Apache-2.0
//! OpenZeppelin's `Ownable2Step` for Stylus programs, with the same functions, events, errors and selectors.
//! The first owner is an explicit argument: a constructor reached through StylusDeployer sees StylusDeployer
//! as the sender.

#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

use alloc::vec::Vec;

use stylus_sdk::alloy_primitives::Address;
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::StorageAddress;

sol! {
    event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    #[derive(Debug, PartialEq, Eq)]
    error OwnableUnauthorizedAccount(address account);
    #[derive(Debug, PartialEq, Eq)]
    error OwnableInvalidOwner(address owner);
}

/// The owner and the pending owner of a program.
#[storage]
pub struct Ownable2Step {
    owner: StorageAddress,
    pending_owner: StorageAddress,
}

/// The ownership interface of OpenZeppelin's `Ownable2Step`.
#[public]
pub trait IOwnable2Step {
    type Error: Into<Vec<u8>>;

    /// The current owner.
    fn owner(&self) -> Address;

    /// The address that may accept ownership; zero when no transfer is pending.
    fn pending_owner(&self) -> Address;

    /// Starts transferring ownership to `new_owner`, who must accept it. Zero cancels a pending transfer.
    /// Owner only.
    fn transfer_ownership(&mut self, new_owner: Address) -> Result<(), Self::Error>;

    /// Accepts a pending ownership transfer. Pending owner only.
    fn accept_ownership(&mut self) -> Result<(), Self::Error>;

    /// Leaves the program without an owner, for good. Owner only.
    fn renounce_ownership(&mut self) -> Result<(), Self::Error>;
}

impl Ownable2Step {
    /// Sets the first owner, which may not be zero.
    pub fn initialize<E: From<OwnableInvalidOwner>>(
        &mut self,
        initial_owner: Address,
    ) -> Result<(), E> {
        if initial_owner == Address::ZERO {
            return Err(OwnableInvalidOwner {
                owner: Address::ZERO,
            }
            .into());
        }
        self.transfer_to(initial_owner);
        Ok(())
    }

    /// The current owner.
    pub fn owner(&self) -> Address {
        self.owner.get()
    }

    /// The address that may accept ownership; zero when no transfer is pending.
    pub fn pending_owner(&self) -> Address {
        self.pending_owner.get()
    }

    /// Fails unless the caller is the owner.
    pub fn only_owner<E: From<OwnableUnauthorizedAccount>>(&self) -> Result<(), E> {
        let sender = self.vm().msg_sender();
        if sender != self.owner.get() {
            return Err(OwnableUnauthorizedAccount { account: sender }.into());
        }
        Ok(())
    }

    /// Starts transferring ownership to `new_owner`; zero cancels a pending transfer. Owner only.
    pub fn transfer_ownership<E: From<OwnableUnauthorizedAccount>>(
        &mut self,
        new_owner: Address,
    ) -> Result<(), E> {
        self.only_owner()?;
        self.pending_owner.set(new_owner);
        self.vm().log(OwnershipTransferStarted {
            previousOwner: self.owner.get(),
            newOwner: new_owner,
        });
        Ok(())
    }

    /// Makes the pending owner the owner. Pending owner only.
    pub fn accept_ownership<E: From<OwnableUnauthorizedAccount>>(&mut self) -> Result<(), E> {
        let sender = self.vm().msg_sender();
        if sender != self.pending_owner.get() {
            return Err(OwnableUnauthorizedAccount { account: sender }.into());
        }
        self.transfer_to(sender);
        Ok(())
    }

    /// Leaves the program without an owner. Owner only.
    pub fn renounce_ownership<E: From<OwnableUnauthorizedAccount>>(&mut self) -> Result<(), E> {
        self.only_owner()?;
        self.transfer_to(Address::ZERO);
        Ok(())
    }

    fn transfer_to(&mut self, new_owner: Address) {
        let previous = self.owner.get();
        self.owner.set(new_owner);
        self.pending_owner.set(Address::ZERO);
        self.vm().log(OwnershipTransferred {
            previousOwner: previous,
            newOwner: new_owner,
        });
    }
}
