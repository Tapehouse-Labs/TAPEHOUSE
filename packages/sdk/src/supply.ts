// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, type Hex, parseSignature } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
import type { At } from './band.js'
import { type Deployments, entry } from './deployments.js'
import { supplyVaultAbi } from './generated.js'
import { holder } from './vault.js'

/** Deposits `assets` of the caller's USDG with the lenders, for shares to `receiver`. */
export function deposit(deployments: Deployments, { assets, receiver }: { assets: bigint; receiver: Address }) {
  return { address: vault(deployments), abi: supplyVaultAbi, functionName: 'deposit', args: [assets, receiver] } as const
}

/**
 * Deposits as `deposit` does, with the caller's EIP-2612 `signature` of a USDG permit for the vault to spend `assets`
 * until `deadline`. A permit already used, as by a front-runner, does not stop the deposit if the allowance stands.
 */
export function depositWithPermit(
  deployments: Deployments,
  { assets, receiver, deadline, signature }: { assets: bigint; receiver: Address; deadline: bigint; signature: Hex },
) {
  const { v, yParity, r, s } = parseSignature(signature)
  return {
    address: vault(deployments),
    abi: supplyVaultAbi,
    functionName: 'depositWithPermit',
    args: [assets, receiver, deadline, Number(v ?? BigInt(yParity) + 27n), r, s],
  } as const
}

/** Sends `assets` of `owner`'s USDG to `receiver`, at most what the vault holds and is not lent. */
export function withdraw(
  deployments: Deployments,
  { assets, receiver, owner }: { assets: bigint; receiver: Address; owner: Address },
) {
  return {
    address: vault(deployments),
    abi: supplyVaultAbi,
    functionName: 'withdraw',
    args: [assets, receiver, owner],
  } as const
}

/** Redeems `owner`'s `shares` for USDG to `receiver`, at most what the vault holds and is not lent. */
export function redeem(
  deployments: Deployments,
  { shares, receiver, owner }: { shares: bigint; receiver: Address; owner: Address },
) {
  return {
    address: vault(deployments),
    abi: supplyVaultAbi,
    functionName: 'redeem',
    args: [shares, receiver, owner],
  } as const
}

/**
 * The supply vault at one block: what lenders earn and the borrower pays, a year with 18 decimals; the share lent,
 * with 18 decimals; the USDG it holds and what the borrower owes; and whether it takes deposits, which stop until the
 * borrower is set and while USDG is paused, the vault frozen or its USDG wiped.
 */
export async function market(client: Client, deployments: Deployments, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const read = { address: vault(deployments), abi: supplyVaultAbi, blockNumber } as const
  const [supplyRate, borrowRate, utilization, idle, debt, totalAssets, maxDeposit] = await Promise.all([
    readContract(client, { ...read, functionName: 'supplyRate' }),
    readContract(client, { ...read, functionName: 'borrowRate' }),
    readContract(client, { ...read, functionName: 'utilization' }),
    readContract(client, { ...read, functionName: 'idle' }),
    readContract(client, { ...read, functionName: 'debt' }),
    readContract(client, { ...read, functionName: 'totalAssets' }),
    readContract(client, { ...read, functionName: 'maxDeposit', args: [read.address] }),
  ])
  return { supplyRate, borrowRate, utilization, idle, debt, totalAssets, open: maxDeposit > 0n }
}

/**
 * What `owner` lent, at one block: its shares, with 12 decimals, and the USDG they are worth; what it may withdraw
 * now, at most what the vault holds; and what it may deposit.
 */
export async function lender(client: Client, deployments: Deployments, owner: Address, at: At = {}) {
  const blockNumber = at.blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  return holder(client, vault(deployments), owner, blockNumber)
}

function vault(deployments: Deployments) {
  return entry(deployments.tapehouse, 'SupplyVault', '.tapehouse')
}
