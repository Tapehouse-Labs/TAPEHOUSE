// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Address, type Client, erc4626Abi } from 'viem'
import { readContract } from 'viem/actions'

/**
 * What `owner` holds of the ERC-4626 vault `address` at `blockNumber`: its shares, the assets they are worth, the
 * assets and shares it may take out now, and the assets it may put in.
 */
export async function holder(client: Client, address: Address, owner: Address, blockNumber: bigint) {
  const read = { address, abi: erc4626Abi, blockNumber } as const
  const [shares, maxWithdraw, maxRedeem, maxDeposit] = await Promise.all([
    readContract(client, { ...read, functionName: 'balanceOf', args: [owner] }),
    readContract(client, { ...read, functionName: 'maxWithdraw', args: [owner] }),
    readContract(client, { ...read, functionName: 'maxRedeem', args: [owner] }),
    readContract(client, { ...read, functionName: 'maxDeposit', args: [owner] }),
  ])
  const assets = await readContract(client, { ...read, functionName: 'convertToAssets', args: [shares] })
  return { shares, assets, maxWithdraw, maxRedeem, maxDeposit }
}
