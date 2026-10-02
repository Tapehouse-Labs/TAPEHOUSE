// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, type Hex, hexToString } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry, toBytes32 } from './deployments.js'
import { marginAccountsAbi } from './generated.js'

/**
 * Lets `authorized` borrow and withdraw for the caller's account, deposit Stock Tokens and USDG into it, and act for
 * its shorts, or stops it. Until then the margin accounts and the shorts revert with `Unauthorized(caller, account)`.
 */
export function setAuthorization(deployments: Deployments, authorized: Address, allowed: boolean) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'setAuthorization',
    args: [authorized, allowed],
  } as const
}

/**
 * Deposits `amount` of `token` from the caller into `account`'s `position`: `CROSS` or an asset's symbol. A basket's
 * address deposits its shares, into `CROSS` alone.
 */
export function deposit(
  deployments: Deployments,
  { position, token, amount, account }: { position: Hex; token: Address; amount: bigint; account: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'deposit',
    args: [position, token, amount, account],
  } as const
}

/** Withdraws `amount` of `token`, a basket's shares included, from `account`'s `position` to `receiver`. */
export function withdraw(
  deployments: Deployments,
  {
    position,
    token,
    amount,
    account,
    receiver,
  }: { position: Hex; token: Address; amount: bigint; account: Address; receiver: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'withdraw',
    args: [position, token, amount, account, receiver],
  } as const
}

/** Borrows `assets` of USDG against `account`'s `position`, to `receiver`. */
export function borrow(
  deployments: Deployments,
  { position, assets, account, receiver }: { position: Hex; assets: bigint; account: Address; receiver: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'borrow',
    args: [position, assets, account, receiver],
  } as const
}

/**
 * Repays `assets` of USDG for `account`'s `position` with the caller's USDG: its debt first, then its premium with what
 * is left. Anyone may. `repayment` reads what a full repayment takes.
 */
export function repay(
  deployments: Deployments,
  { position, assets, account }: { position: Hex; assets: bigint; account: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'repay',
    args: [position, assets, account],
  } as const
}

/**
 * Redeems `shares` of the registry's basket `basket` in `account`'s cross position for its Stock Tokens, which the
 * position then holds. The account or an address it authorized may, and anyone once the position falls short.
 */
export function unwrap(
  deployments: Deployments,
  { account, basket, shares }: { account: Address; basket: string; shares: bigint },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'unwrap',
    args: [account, entry(deployments.baskets, basket, '.tapehouse.Baskets'), shares],
  } as const
}

/** Whether `authorized` may act for `account`. */
export function isAuthorized(
  client: Client,
  deployments: Deployments,
  account: Address,
  authorized: Address,
  at: At = {},
) {
  return readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'isAuthorized',
    args: [account, authorized],
    ...at,
  })
}

/** What repays all `account`'s `position` owes now: its debt and its premium, in USDG, both read at one block. */
export async function repayment(
  client: Client,
  deployments: Deployments,
  account: Address,
  position: Hex,
  at: At = {},
) {
  const address = accounts(deployments)
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: marginAccountsAbi, args: [account, position], blockNumber } as const
  const [debt, premium] = await Promise.all([
    readContract(client, { ...read, functionName: 'debt' }),
    readContract(client, { ...read, functionName: 'premium' }),
  ])
  return { debt, premium, assets: debt + premium }
}

/**
 * What `account`'s `position` holds of each Stock Token through its baskets, by asset, as its shares would redeem now:
 * the engine margins them as those Stock Tokens. Both reads are at one block.
 */
export async function inBaskets(
  client: Client,
  deployments: Deployments,
  account: Address,
  position: Hex,
  at: At = {},
): Promise<Record<string, bigint>> {
  const address = accounts(deployments)
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const [[symbols], amounts] = await Promise.all([
    readContract(client, { address, abi: marginAccountsAbi, functionName: 'stocks', blockNumber }),
    readContract(client, {
      address,
      abi: marginAccountsAbi,
      functionName: 'inBaskets',
      args: [account, position],
      blockNumber,
    }),
  ])
  return Object.fromEntries(symbols.map((symbol, i) => [hexToString(symbol, { size: 32 }), amounts[i] ?? 0n]))
}

/**
 * `account`'s `position` as the engine sees it: equity and requirement in USD with 18 decimals, the engine's
 * missing bits and the regime.
 */
export async function health(client: Client, deployments: Deployments, account: Address, position: Hex, at: At = {}) {
  const [equity, requirement, missing, regime] = await readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'health',
    args: [account, position],
    ...at,
  })
  return { equity, requirement, missing, regime }
}

/** What `account`'s `position` holds of `token`: USDG, WETH or a Stock Token. */
export function collateral(
  client: Client,
  deployments: Deployments,
  account: Address,
  position: Hex,
  token: Address,
  at: At = {},
) {
  return readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'collateral',
    args: [account, position, token],
    ...at,
  })
}

/**
 * The gross exposure of `account`'s `position` over its equity, in basis points; `maxUint256` for a position without
 * equity.
 */
export function leverage(client: Client, deployments: Deployments, account: Address, position: Hex, at: At = {}) {
  return readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'leverage',
    args: [account, position],
    ...at,
  })
}

/**
 * The price of `asset`, in USD with 8 decimals, at or above which `account`'s `position` meets its requirement once it
 * owes `borrowing` more USDG: its band's low edge where it falls short already, zero where no price leaves it short.
 */
export function liquidationPrice(
  client: Client,
  deployments: Deployments,
  account: Address,
  position: Hex,
  asset: string,
  borrowing: bigint,
  at: At = {},
) {
  return readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'liquidationPrice',
    args: [account, position, toBytes32(asset), borrowing],
    ...at,
  })
}

function accounts(deployments: Deployments) {
  return entry(deployments.tapehouse, 'MarginAccounts', '.tapehouse')
}
