# @tapehouse/sdk

Typed reads and transactions for [Tapehouse](https://github.com/Tapehouse-Labs/TAPEHOUSE)'s price band, its feeds, the margin accounts, the baskets, the short positions, the Morpho oracles and the gap cover, on [viem](https://viem.sh).

```bash
pnpm add @tapehouse/sdk viem
```

Every call takes a chain's registry, `deployments/<chainId>.json` from the Tapehouse repository, parsed with `parseDeployments`: the SDK compiles in no address. Reads take a viem client. Transactions come back as viem contract calls, ready for `simulateContract` and `writeContract`.

```ts
import { readFileSync } from 'node:fs'
import { createWalletClient, http, publicActions } from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { band, parseDeployments, shorts } from '@tapehouse/sdk'

const deployments = parseDeployments(JSON.parse(readFileSync('deployments/4663.json', 'utf8')))
const client = createWalletClient({ account: privateKeyToAccount(key), transport: http(rpc) }).extend(publicActions)
const account = client.account.address

const { state, low, high } = await band.quote(client, deployments, 'NVDA')
const sale = await shorts.quoteSale(client, deployments, { asset: 'NVDA', amount: 10n ** 18n, slippageBps: 50n })
const { request } = await client.simulateContract({
  ...shorts.sell(deployments, { asset: 'NVDA', amount: 10n ** 18n, minProceeds: sale.minProceeds, account }),
  account: client.account,
})
await client.writeContract(request)
```

- `band`: the band's `quote`, `session`, `halt` and stored `price`; each asset's `BandFeed` (`latestRound`, `latestBand`, `sealed`, `seal`); `chainlinkRound` of a `.chainlink` feed; and `writePrices` with the signed RedStone packages a `PackageSource` of your own supplies. The SDK holds no API key.
- `accounts`: `setAuthorization`, `deposit`, `withdraw`, `borrow`, `repay`, `unwrap`, `isAuthorized`, `health`, `collateral`, `leverage`, `liquidationPrice`, `repayment`, the debt and premium a full repayment takes, and `inBaskets`, what a position holds of each Stock Token through its baskets, each pair read at one block. A basket's address deposits and withdraws its shares, in the cross position alone.
- `baskets`: `mint` and `redeem` of a basket of `.tapehouse.Baskets`, in kind, and the views `components`, `previewMint`, `previewRedeem`, `target` and `pendingTarget`.
- `shorts`: `deposit`, `withdraw`, `sell`, `cover`, `liquidate`, `mark`, the views `position`, `health`, `epoch`, `book`, `restriction` and `fee`, and `quoteSale` and `quoteCover`, the limits of a sale and a buy-back from Uniswap's QuoterV2, with a slippage of 0 to 10,000 basis points.
- `morpho`: an asset's Morpho oracle from `.morphoOracles`, its `price`, which is no price, never zero, where it reverts with `NoAnswer` or `SequencerNotSettled`, and says why (`stale`, `halted` or `sequencerNotSettled`), its `halt`, and the `oracle` itself: band, symbol, tokens, scale and owner.
- `gapCover`: `buy`, `deposit`, `redeem`, `record`, `measure`, `observe`, `settle`, `release`, `claim`, `quote` with the writers' USDG a cover reserves, `pricingGap`, the weekend gap it prices at, which follows the week before the close, `minDeductible`, `capacity`, `lastCloseMs`, `sales`, `series` (open, settled with its reference and reopening prices, or void), `cover`, `payouts`, and `payout`, what a settled cover pays.
- `decodeRevert`: the name and arguments of any revert of the contracts, the baskets, the gap cover, the band, the Stock Tokens, USDG, or Solidity's `Error(string)` and `Panic`.
- `tokenPriceFeed` and `sharePriceFeed`: a `.chainlink` feed typed by what it prices, the Stock Token on Robinhood Chain and the share on Arbitrum One.
- The contracts' ABIs, `as const`, generated from the Solidity build.

`examples/devnode.ts` runs all of it against a local Nitro dev node; it reads the account owner's key from `PRIVATE_KEY` in its environment, never from its command line.

## License

Licensed under either of [Apache License, Version 2.0](LICENSE-APACHE) or [MIT license](LICENSE-MIT) at your option.
