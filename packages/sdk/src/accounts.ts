// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, erc20Abi, type Hex, hexToString, parseSignature, zeroAddress } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { CROSS, type Deployments, entry, toBytes32 } from './deployments.js'
import { marginAccountsAbi, stockTokenAbi, stockTokenRegistryAbi, supplyVaultAbi, wethAbi } from './generated.js'

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

/**
 * Deposits as `deposit` does, with the caller's EIP-2612 `signature` of a permit for the accounts to spend `amount`
 * of `token` until `deadline`. A permit already used, as by a front-runner, does not stop the deposit if the
 * allowance stands.
 */
export function depositWithPermit(
  deployments: Deployments,
  {
    position,
    token,
    amount,
    account,
    deadline,
    signature,
  }: { position: Hex; token: Address; amount: bigint; account: Address; deadline: bigint; signature: Hex },
) {
  const { v, yParity, r, s } = parseSignature(signature)
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'depositWithPermit',
    args: [position, token, amount, account, deadline, Number(v ?? BigInt(yParity) + 27n), r, s],
  } as const
}

/**
 * Wraps `value` of ether into the registry's WETH for `account`, which then deposits it: the accounts take ether as
 * WETH.
 */
export function wrap(deployments: Deployments, { account, value }: { account: Address; value: bigint }) {
  return {
    address: entry(deployments.tokens, 'WETH', '.tokens'),
    abi: wethAbi,
    functionName: 'depositTo',
    args: [account],
    value,
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

/**
 * Lends `amount` of `account`'s `position`'s holding of the Stock Token `token` through the asset's lending vault. The
 * position keeps the vault's shares, which count for the asset at `RECALL_HAIRCUT` less in its equity; a position in
 * debt must still meet its requirement.
 */
export function lend(
  deployments: Deployments,
  { position, token, amount, account }: { position: Hex; token: Address; amount: bigint; account: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'lend',
    args: [position, token, amount, account],
  } as const
}

/**
 * Takes `amount` of what `account`'s `position` lent of `token` back into its holding, as far as the vault can return
 * it now; past that it reverts with `OutOfReach(symbol, amount, reachable)`.
 */
export function unlend(
  deployments: Deployments,
  { position, token, amount, account }: { position: Hex; token: Address; amount: bigint; account: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'unlend',
    args: [position, token, amount, account],
  } as const
}

/**
 * Recalls `amount` of what `account`'s `position` lent of `token`: what the vault holds free comes back at once, and
 * the rest is a ticket at the end of the vault's queue, which the borrower returns within the vault's `NOTICE`.
 */
export function recall(
  deployments: Deployments,
  { position, token, amount, account }: { position: Hex; token: Address; amount: bigint; account: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'recall',
    args: [position, token, amount, account],
  } as const
}

/** Takes what the vault holds for `account`'s `position`'s recall of `token`, next in turn, into its holding. Anyone may. */
export function settle(
  deployments: Deployments,
  { account, position, token }: { account: Address; position: Hex; token: Address },
) {
  return {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'settle',
    args: [account, position, token],
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

/** The assets, in the engine's order, and their Stock Tokens; an asset without a token takes no deposit. */
export async function stocks(client: Client, deployments: Deployments, at: At = {}) {
  const [symbols, tokens] = await readContract(client, {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    functionName: 'stocks',
    ...at,
  })
  return symbols.map((symbol, i) => ({ asset: hexToString(symbol, { size: 32 }), token: tokens[i] ?? zeroAddress }))
}

/**
 * Every position an account may hold: `CROSS`, then the isolated position of each asset with a Stock Token, in the
 * engine's order.
 */
export async function positions(client: Client, deployments: Deployments, at: At = {}): Promise<Hex[]> {
  const assets = await stocks(client, deployments, at)
  return [CROSS, ...assets.filter(({ token }) => token !== zeroAddress).map(({ asset }) => toBytes32(asset))]
}

/**
 * What `account`'s `position` has of the Stock Token `token` beyond what it holds, at one block: what it lent, its fee
 * included; what the liquidator may take now, held and reachable together; what it recalled and has not taken back;
 * and the vault's tickets of its open recalls, oldest first.
 */
export async function lending(
  client: Client,
  deployments: Deployments,
  account: Address,
  position: Hex,
  token: Address,
  at: At = {},
) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = {
    address: accounts(deployments),
    abi: marginAccountsAbi,
    args: [account, position, token],
    blockNumber,
  } as const
  const [lent, sellable, claim, recalls] = await Promise.all([
    readContract(client, { ...read, functionName: 'lent' }),
    readContract(client, { ...read, functionName: 'sellable' }),
    readContract(client, { ...read, functionName: 'claim' }),
    readContract(client, { ...read, functionName: 'recalls' }),
  ])
  return { lent, sellable, claim, recalls }
}

/**
 * The accounts' holding of `asset`, at one block: the units its positions hold, each worth `scale` / 10^18 of a token,
 * so `quantity` in raw units; its cap in raw units; and the Stock Tokens the accounts hold, below `quantity` while a
 * write-down by the issuer is pending.
 */
export async function holding(client: Client, deployments: Deployments, asset: string, token: Address, at: At = {}) {
  const address = accounts(deployments)
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const [[units, scale, cap], balance] = await Promise.all([
    readContract(client, {
      address,
      abi: marginAccountsAbi,
      functionName: 'holding',
      args: [toBytes32(asset)],
      blockNumber,
    }),
    readContract(client, { address: token, abi: erc20Abi, functionName: 'balanceOf', args: [address], blockNumber }),
  ])
  return { units, scale, quantity: (units * scale) / 10n ** 18n, cap, balance }
}

/**
 * The accounts' limits on borrowing, at one block: what the positions owe the vault together, the debt cap, the debt
 * cap from the regular open on the last trading day before a closure, the weekend leverage cap in basis points,
 * whether the guardian paused new borrowing, and the premium charged while the session is closed, in basis points a
 * year.
 */
export async function limits(client: Client, deployments: Deployments, at: At = {}) {
  const address = accounts(deployments)
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address, abi: marginAccountsAbi, blockNumber } as const
  const [debt, debtCap, weekendDebtCap, weekendLeverage, borrowingPaused, premiumRate] = await Promise.all([
    readContract(client, {
      address: entry(deployments.tapehouse, 'SupplyVault', '.tapehouse'),
      abi: supplyVaultAbi,
      functionName: 'debt',
      blockNumber,
    }),
    readContract(client, { ...read, functionName: 'debtCap' }),
    readContract(client, { ...read, functionName: 'weekendDebtCap' }),
    readContract(client, { ...read, functionName: 'weekendLeverage' }),
    readContract(client, { ...read, functionName: 'borrowingPaused' }),
    readContract(client, { ...read, functionName: 'premiumRate' }),
  ])
  return { debt, debtCap, weekendDebtCap, weekendLeverage, borrowingPaused, premiumRate }
}

/**
 * What stops the Stock Token `token` moving, at one block: the issuer's pause, and whether its registry blocks each of
 * `holders`, by address. A pause or a block on the accounts or on the user stops its deposits and withdrawals.
 */
export async function issuer(client: Client, token: Address, holders: readonly Address[], at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: token, abi: stockTokenAbi, blockNumber } as const
  const [paused, registry] = await Promise.all([
    readContract(client, { ...read, functionName: 'paused' }),
    readContract(client, { ...read, functionName: 'ACCESS_CONTROLLED_REGISTRY' }),
  ])
  const blocked = await Promise.all(
    holders.map((holder) =>
      readContract(client, {
        address: registry,
        abi: stockTokenRegistryAbi,
        functionName: 'isBlocked',
        args: [holder],
        blockNumber,
      }),
    ),
  )
  return { paused, blocked: Object.fromEntries(holders.map((holder, i) => [holder, blocked[i] ?? false])) }
}

function accounts(deployments: Deployments) {
  return entry(deployments.tapehouse, 'MarginAccounts', '.tapehouse')
}
