---
name: tapehouse-x402
description: Read Tapehouse's public REST API without an API key by paying per request in USDG over x402. Use when a call to the Tapehouse indexer API answers 402 Payment Required or 429, when asked how an agent can pay for API reads, or when writing a client that pays with an x402 library.
---

# Paying for Tapehouse's API over x402

The indexer's API (`/v1/status`, `/v1/events`, `/v1/contracts/{contract}/views/{function}`, the composite reads) is free for a client address up to 10 requests a second with bursts of 20. Beyond that the operator may sell each request instead of refusing it: the API answers **402 Payment Required** with the price, and an x402 client pays in USDG and repeats the request. No key, no account.

The MCP server needs none of this: it reads the chain through the RPC. Use x402 to read the indexer's history, streams of events and composite reads.

## What a 402 says

The `PAYMENT-REQUIRED` header is the base64 of a JSON document of x402 version 2:

- `accepts[0].scheme` is `exact`, `network` is `eip155:46630` for Robinhood Chain's testnet, `asset` is USDG's address, `amount` is its atomic units (6 decimals: `1000` is 0.001 USDG), `payTo` is who receives it, and `extra` is the token's EIP-712 domain (`name`, `version`) the payment is signed under.
- The client signs an EIP-3009 `transferWithAuthorization` of that amount, to `payTo`, and repeats the request with `PAYMENT-SIGNATURE`, the base64 of the payment payload.
- The server verifies it, answers the read, and settles the transfer through its facilitator. The response is released only after the payment settles, and `PAYMENT-RESPONSE` carries the settlement, including the transaction hash. A payment that does not settle gets 402 again and the read is not served.

The stream (`/v1/stream`) and the specification (`/v1/openapi.yaml`) are not sold; over the allowance they answer 429 with `Retry-After`.

## Paying

Use the official client, which signs the authorization and retries, and **allow USDG explicitly with a cap**: x402 clients refuse any asset they were not told to accept, which is the protection against a server naming another token or a large amount.

TypeScript, with `@x402/fetch`, `@x402/evm` and `viem`:

```ts
import { wrapFetchWithPayment, x402Client } from '@x402/fetch'
import { ExactEvmScheme } from '@x402/evm/exact/client'
import { privateKeyToAccount } from 'viem/accounts'

const usdg = '0x7E955252E15c84f5768B83c41a71F9eba181802F'
const client = x402Client.fromConfig({
  schemes: [{ network: 'eip155:46630', client: new ExactEvmScheme(privateKeyToAccount(process.env.PAYER_KEY as `0x${string}`)) }],
  spendControls: { allowedAssets: [{ network: 'eip155:46630', asset: usdg, maxAmountPerPayment: '10000' }] },
})
const fetchWithPayment = wrapFetchWithPayment(fetch, client)
const response = await fetchWithPayment('https://indexer-production-f5dc.up.railway.app/v1/contracts')
```

Go, with `github.com/x402-foundation/x402/go/v2`: `x402.Newx402Client(x402.WithSpendControls(x402.SpendControls{AllowedAssets: []x402.SpendControlAsset{{Network: "eip155:46630", Asset: usdg, MaxAmountPerPayment: "10000"}}}))`, registered with `exactclient.NewExactEvmScheme(signer, nil)` and wrapped with `x402http.WrapHTTPClientWithPayment`.

## Rules for an agent

- **The payer's key is the user's.** Never ask for it in chat. A key in an environment variable the user set is the way; a smart account's signature works too.
- **Pay from a wallet that holds only what it may spend.** Fund it with a few USDG; each request costs a fraction of one.
- **Read the price from the 402 before paying and refuse a change.** Cap `MaxAmountPerPayment` (or the library's equivalent) at a small multiple of the advertised price.
- **Do not loop.** If a read is needed repeatedly, page `/v1/events` once and keep the result; or ask the operator for a key.
- Testnet USDG comes from Paxos' faucet at https://faucet.paxos.com/?network=robinhood.
