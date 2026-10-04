// SPDX-License-Identifier: MIT OR Apache-2.0
// Usage: PRIVATE_KEY=0x… node examples/sponsored.ts RPC_URL BUNDLER_URL SPONSOR_URL DEPLOYMENTS_JSON
// Sends the same two calls, an approval of USDG to the margin accounts and an authorization in them, first from a fresh
// address the funded account pays the gas of, then as one user operation of a smart account whose owner holds no
// ether: Tapehouse's sponsor service signs it over ERC-7677 and the sponsor paymaster pays for it and for creating the
// account. Checks that both reach the margin accounts, that the service refuses a call outside Tapehouse, and prints
// the L2 gas of each way: Arbitrum's receipts count the L1 calldata's cost apart.
import { readFileSync } from 'node:fs'
import {
  type Address,
  createPublicClient,
  createWalletClient,
  defineChain,
  type Hex,
  http,
  maxUint256,
  parseEther,
} from 'viem'
import { createBundlerClient, createPaymasterClient } from 'viem/account-abstraction'
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts'
import { toSimpleSmartAccount } from 'permissionless/accounts'
import { accounts, marginAccountsAbi, parseDeployments, sponsorship, usdgAbi } from '@tapehouse/sdk'

const [rpc, bundlerUrl, sponsorUrl, registry] = process.argv.slice(2) as [string, string, string, string]
const key = process.env.PRIVATE_KEY as Hex | undefined
if (key === undefined) throw new Error("Set PRIVATE_KEY to a funded account's key.")
const deployments = parseDeployments(JSON.parse(readFileSync(registry, 'utf8')))
const chain = defineChain({
  id: deployments.chainId,
  name: 'dev node',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: [rpc] } },
})
const client = createPublicClient({ chain, transport: http(rpc) })
const funded = createWalletClient({ account: privateKeyToAccount(key), chain, transport: http(rpc) })
const fresh = createWalletClient({ account: privateKeyToAccount(generatePrivateKey()), chain, transport: http(rpc) })
const owner = privateKeyToAccount(generatePrivateKey())
const operator = privateKeyToAccount(generatePrivateKey()).address
const marginAccounts = deployments.tapehouse.MarginAccounts as Address
const usdg = deployments.tokens.USDG as Address

async function l2Gas(hash: Hex) {
  const receipt = (await client.request({ method: 'eth_getTransactionReceipt', params: [hash] })) as {
    gasUsed: Hex
    gasUsedForL1: Hex
  }
  return BigInt(receipt.gasUsed) - BigInt(receipt.gasUsedForL1)
}

function check(condition: boolean, message: string): asserts condition {
  if (!condition) throw new Error(`FAIL: ${message}`)
}

const approval = {
  address: usdg,
  abi: usdgAbi,
  functionName: 'approve',
  args: [marginAccounts, maxUint256],
} as const
const authorization = accounts.setAuthorization(deployments, operator, true)

await client.waitForTransactionReceipt({
  hash: await funded.sendTransaction({ to: fresh.account.address, value: parseEther('0.01') }),
})
let eoaGas = 0n
for (const hash of [await fresh.writeContract(approval), await fresh.writeContract(authorization)]) {
  const receipt = await client.waitForTransactionReceipt({ hash })
  check(receipt.status === 'success', `the transaction ${hash} reverted`)
  eoaGas += await l2Gas(hash)
}

const account = await toSimpleSmartAccount({
  client,
  owner,
  factoryAddress: sponsorship.accountFactory(deployments),
  entryPoint: sponsorship.entryPoint(deployments),
})
const bundler = createBundlerClient({
  account,
  client,
  paymaster: createPaymasterClient({ transport: http(sponsorUrl) }),
  transport: http(bundlerUrl),
})
check((await client.getCode({ address: account.address })) === undefined, 'the smart account exists already')
const hash = await bundler.sendUserOperation({
  calls: [sponsorship.asCall(approval), sponsorship.asCall(authorization)],
})
const { success, actualGasCost, receipt, paymaster } = await bundler.waitForUserOperationReceipt({ hash })
check(success, 'the user operation reverted')
check(paymaster === sponsorship.paymaster(deployments), 'the sponsor paymaster did not pay')
check((await client.getCode({ address: account.address })) !== undefined, 'the smart account was not created')
check((await client.getBalance({ address: account.address })) === 0n, 'the smart account holds ether')
check((await client.getBalance({ address: owner.address })) === 0n, 'the owner holds ether')
check(
  await client.readContract({
    address: marginAccounts,
    abi: marginAccountsAbi,
    functionName: 'isAuthorized',
    args: [account.address, operator],
  }),
  'the margin accounts did not authorize the operator for the smart account',
)
check(
  (await client.readContract({ ...approval, functionName: 'allowance', args: [account.address, marginAccounts] })) ===
    maxUint256,
  'the smart account did not approve the margin accounts',
)
const left = await sponsorship.freeOperationsLeft(client, deployments, account.address)
check(left === 2n, `the smart account has ${left} free operations left`)

const refusal = await bundler
  .sendUserOperation({ calls: [{ to: funded.account.address, value: 0n }] })
  .then(
    () => undefined,
    (error: unknown) => String(error),
  )
check(refusal?.includes('is not to a Tapehouse contract') === true, 'the service sponsored a call outside Tapehouse')

console.log(`smart account ${account.address}, created and authorized in ${receipt.transactionHash}`)
console.log(`from an address holding ether: ${eoaGas} L2 gas in two transactions`)
console.log(`as a sponsored user operation creating the account: ${await l2Gas(receipt.transactionHash)} L2 gas; the paymaster paid ${actualGasCost} wei`)
console.log(`free operations left: ${left}; a call outside Tapehouse is refused`)
