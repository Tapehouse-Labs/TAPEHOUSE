// SPDX-License-Identifier: MIT OR Apache-2.0
import {
  type Abi,
  type Address,
  type Client,
  type ContractFunctionArgs,
  ContractFunctionExecutionError,
  type ContractFunctionName,
  erc20Abi,
  hashDomain,
  type Hex,
  isAddressEqual,
  parseAbi,
  parseSignature,
  type TypedDataDomain,
} from 'viem'
import { getChainId, readContract, simulateContract } from 'viem/actions'
import { type Deployments, entry } from './deployments.js'
import { sponsorPaymasterAbi, stockTokenAbi, usdgAbi } from './generated.js'

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

const permitTypes = {
  Permit: [
    { name: 'owner', type: 'address' },
    { name: 'spender', type: 'address' },
    { name: 'value', type: 'uint256' },
    { name: 'nonce', type: 'uint256' },
    { name: 'deadline', type: 'uint256' },
  ],
} as const

/**
 * The EIP-2612 permit `owner` signs, for free, to let `account` take `value` of `token` until `deadline`, ready for
 * viem's `signTypedData`. The domain is read from the chain each time: a Stock Token's from its `eip712Domain()`, as
 * its issuer may rename it; a token without one, such as USDG, from its `name()` and version 1. Either must match the
 * token's `DOMAIN_SEPARATOR()` on the client's chain, so a signature never goes out in a domain the token rejects.
 */
export async function permit(
  client: Client,
  {
    token,
    owner,
    account,
    value,
    deadline,
  }: { token: Address; owner: Address; account: Address; value: bigint; deadline: bigint },
) {
  const [domain, nonce] = await Promise.all([
    permitDomain(client, token),
    readContract(client, { address: token, abi: usdgAbi, functionName: 'nonces', args: [owner] }),
  ])
  return {
    domain,
    types: permitTypes,
    primaryType: 'Permit',
    message: { owner, spender: account, value, nonce, deadline },
  } as const
}

/**
 * The calls by which `account` takes `value` of `token` from `owner` with `owner`'s permit `signature`: the permit,
 * then the transfer, as calls of a user operation. The sponsor paymaster pays for both.
 */
export function pull({
  token,
  owner,
  account,
  value,
  deadline,
  signature,
}: {
  token: Address
  owner: Address
  account: Address
  value: bigint
  deadline: bigint
  signature: Hex
}) {
  const { v, yParity, r, s } = parseSignature(signature)
  return [
    asCall({
      address: token,
      abi: usdgAbi,
      functionName: 'permit',
      args: [owner, account, value, deadline, Number(v ?? BigInt(yParity) + 27n), r, s],
    }),
    asCall({ address: token, abi: erc20Abi, functionName: 'transferFrom', args: [owner, account, value] }),
  ] as const
}

async function permitDomain(client: Client, token: Address): Promise<TypedDataDomain> {
  const [separator, chainId] = await Promise.all([
    readContract(client, { address: token, abi: usdgAbi, functionName: 'DOMAIN_SEPARATOR' }),
    getChainId(client),
  ])
  let domain: { name: string; version: string; chainId: number; verifyingContract: Address }
  try {
    const [fields, name, version, domainChainId, verifyingContract] = await readContract(client, {
      address: token,
      abi: stockTokenAbi,
      functionName: 'eip712Domain',
    })
    if (fields !== '0x0f') throw new Error(`${token}'s permit domain carries fields ${fields}, not 0x0f.`)
    domain = { name, version, chainId: Number(domainChainId), verifyingContract }
  } catch (error) {
    if (!(error instanceof ContractFunctionExecutionError)) throw error
    const name = await readContract(client, { address: token, abi: erc20Abi, functionName: 'name' })
    domain = { name, version: '1', chainId, verifyingContract: token }
  }
  if (
    domain.chainId !== chainId ||
    !isAddressEqual(domain.verifyingContract, token) ||
    hashDomain({
      domain: { ...domain, chainId: BigInt(domain.chainId) },
      types: { EIP712Domain: eip712DomainType },
    }) !== separator
  )
    throw new Error(`${token}'s permit domain does not match its DOMAIN_SEPARATOR on chain ${chainId}.`)
  return domain
}

const eip712DomainType = [
  { name: 'name', type: 'string' },
  { name: 'version', type: 'string' },
  { name: 'chainId', type: 'uint256' },
  { name: 'verifyingContract', type: 'address' },
] as const

/** Arbitrum's NodeInterface, which answers only `eth_call`. */
const nodeInterface = {
  address: '0x00000000000000000000000000000000000000C8',
  abi: parseAbi([
    'function gasEstimateComponents(address to, bool contractCreation, bytes data) payable returns (uint64 gasEstimate, uint64 gasEstimateForL1, uint256 baseFee, uint256 l1BaseFeeEstimate)',
  ]),
} as const

/**
 * The L2 gas of `account`'s calls `callData`, sent by the EntryPoint as a transaction of their own, as Arbitrum's
 * NodeInterface splits it from the L1 gas. A bundler that estimates call gas by searching inside one simulation runs the
 * calls warm after its first pass, and a Stylus program, such as the margin engine, costs more cold: this is the cold
 * figure.
 */
export async function callGas(client: Client, deployments: Deployments, account: Address, callData: Hex) {
  const { result } = await simulateContract(client, {
    ...nodeInterface,
    functionName: 'gasEstimateComponents',
    args: [account, false, callData],
    account: entryPoint(deployments).address,
  })
  return result[0] - result[1]
}

/**
 * The SimpleAccount the registry's factory makes for `owner`, salt 0: the address its permits and wrapped ether go to,
 * read from the factory through `client`. Read it through the wallet's own provider as well before trusting it.
 */
export function accountAddress(client: Client, deployments: Deployments, owner: Address) {
  return readContract(client, {
    address: accountFactory(deployments),
    abi: parseAbi(['function getAddress(address owner, uint256 salt) view returns (address)']),
    functionName: 'getAddress',
    args: [owner, 0n],
  })
}
