// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Abi, Address, Client, ContractFunctionArgs, ContractFunctionName } from 'viem'
import { readContract } from 'viem/actions'
import { type Deployments, entry } from './deployments.js'
import { sponsorPaymasterAbi } from './generated.js'

/** The ERC-4337 EntryPoint smart accounts send their operations through, as viem's smart accounts take it. */
export function entryPoint(deployments: Deployments) {
  return { address: entry(deployments.erc4337, 'EntryPoint', '.erc4337'), version: '0.7' } as const
}

/** The factory of the SimpleAccount v0.7 smart accounts the sponsor paymaster understands. */
export function accountFactory(deployments: Deployments): Address {
  return entry(deployments.erc4337, 'SimpleAccountFactory', '.erc4337')
}

/** How many sponsored operations `account` has left. */
export function freeOperationsLeft(client: Client, deployments: Deployments, account: Address): Promise<bigint> {
  return readContract(client, {
    address: paymaster(deployments),
    abi: sponsorPaymasterAbi,
    functionName: 'freeOperationsLeft',
    args: [account],
  })
}

/** A transaction of this SDK as a call of a user operation. */
export function asCall<
  const abi extends Abi,
  functionName extends ContractFunctionName<abi, 'nonpayable' | 'payable'>,
  args extends ContractFunctionArgs<abi, 'nonpayable' | 'payable', functionName>,
>({ address, ...call }: { address: Address; abi: abi; functionName: functionName; args: args }) {
  return { to: address, ...call }
}

/**
 * The sponsor paymaster. Tapehouse's sponsor service signs the operations it pays for over ERC-7677: viem's
 * `createPaymasterClient` on the service's URL is a bundler client's `paymaster`.
 */
export function paymaster(deployments: Deployments): Address {
  return entry(deployments.tapehouse, 'SponsorPaymaster', '.tapehouse')
}
