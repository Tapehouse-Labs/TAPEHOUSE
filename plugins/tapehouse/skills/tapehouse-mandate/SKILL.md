---
name: tapehouse-mandate
description: Act for a Tapehouse margin account as an agent under a bounded on-chain mandate, or help an account owner give one. Use when asked to let an AI agent or another wallet borrow, withdraw, deposit or repay for an account within limits, to check what a mandate still allows, or when a call reverts with MandateExpired, MandateDebtLimit, MandateLeverageLimit, MandateReceiver, InvalidMandate or Unauthorized.
---

# Tapehouse mandates

A mandate lets an account give an address, such as an agent's own wallet, a bounded authority until a time. It is the alternative to `accounts_set_authorization`, which lets an address do everything: withdraw every collateral token to anywhere and borrow to anywhere, until the account signs the reverse. **Recommend a mandate, not an authorization, for an agent.**

## What a mandate allows

Until `expiresAt`, the agent may:

- **Borrow and withdraw only to the account itself**, never to another receiver, while the position's debt stays at most `maxDebt` (USDG base units, 6 decimals) and its leverage at most `maxLeverage` (basis points, at least 10000, where 10000 is 1x);
- deposit, repay and lend.

An address authorized in full is not bound by a mandate. An `expiresAt` of 0 ends one.

## Giving one (the owner's side)

1. Ask the owner for the agent's address, how long it may act, the debt limit and the leverage limit. Do not choose them for the owner.
2. Call `accounts_set_mandate` with `agent`, `expiresAt` (Unix seconds, in the future), `maxDebt` and `maxLeverage`. Its result names the signer, which is the account itself, and says in `limits` what the signed call lets the agent do. Show the owner that sentence.
3. The owner signs and sends the call. You do not.
4. To end it early: `accounts_set_mandate` with `expiresAt: 0`.

## Acting under one (the agent's side)

Before preparing a borrow or a withdrawal for the account, call `accounts_mandate` with the account, your address as `agent` and the position (`CROSS` by default or an asset's symbol). Its `reason` says where you stand:

| `reason` | Meaning | What to do |
| --- | --- | --- |
| `ok` | Within the mandate | `debtRoom` is how much more USDG you may borrow under the debt limit. Borrow no more than that, and only with `receiver` set to the account |
| `none` | No mandate | Stop. Ask the owner to give one |
| `expired` | It ended | Stop. Ask the owner to renew it |
| `debtLimit` | The position's debt is at its limit | Do not borrow; repaying frees room |
| `leverageLimit` | Leverage is at its limit | Do not borrow or withdraw; repaying or depositing lowers it |

The chain enforces the same rules and reverts with `MandateExpired`, `MandateDebtLimit`, `MandateLeverageLimit`, `MandateReceiver` (a receiver other than the account) or `InvalidMandate` (terms the account cannot set: an expiry in the past, or leverage below 10000). `Unauthorized(caller, account)` means there is no mandate and no authorization at all.

Within the mandate the position can still be unsafe: it is limited by the band, and it can be liquidated. Before a borrow, call `accounts_liquidation_price` with the amount as `borrowing` and tell the owner the price. Prepare the call and say that the **agent's wallet** signs it; the server and you do not.
