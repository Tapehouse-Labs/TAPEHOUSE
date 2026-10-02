// SPDX-License-Identifier: MIT OR Apache-2.0
import type { McpServer } from '@modelcontextprotocol/server'
import { accounts, baskets, gapCover, gapCoverAbi, shorts } from '@tapehouse/sdk'
import {
  type Abi,
  type Address,
  encodeFunctionData,
  type EncodeFunctionDataParameters,
  erc20Abi,
  hexToString,
  maxUint256,
} from 'viem'
import { getBlockNumber, readContract, simulateContract } from 'viem/actions'
import * as z from 'zod'
import {
  address,
  amounts,
  asset,
  basket,
  basketAddress,
  bps,
  byAsset,
  checksummed,
  collateralAddress,
  type Context,
  explainRevert,
  gapCoverAddress,
  position,
  positionId,
  shares,
  token,
  tokenAddress,
  tool,
  uint,
  units,
} from './tool.js'

type ContractCall = { address: Address; abi: Abi; functionName: string; args: readonly unknown[] }

const prepared = z.object({
  chainId: z.number().int(),
  signer: z.string().describe('Who signs and sends the calls, in order.'),
  calls: z.array(
    z.object({
      to: z.string(),
      data: z.string().describe('The calldata, as the wallet sends it.'),
      value: z.literal('0'),
      functionName: z.string(),
      args: z.record(z.string(), z.unknown()).describe("The call's arguments by name, for the signer to check."),
    }),
  ),
})

function call({ address: to, abi, functionName, args }: ContractCall) {
  const item = abi.find(
    (entry) => entry.type === 'function' && entry.name === functionName && entry.inputs.length === args.length,
  )
  if (item?.type !== 'function') throw new Error(`No ${functionName} with ${args.length} arguments.`)
  return {
    to,
    data: encodeFunctionData({ abi, functionName, args } as EncodeFunctionDataParameters),
    value: '0',
    functionName,
    args: Object.fromEntries(item.inputs.map((input, i) => [input.name ?? `${i}`, args[i]])),
  }
}

/**
 * The calls `from` sends: for each token and amount `then` takes, an approval of `spender` for exactly that amount
 * first where `from`'s allowance of the token is short of it, read at `blockNumber`, or at the latest block.
 */
async function approved(
  { client }: Context,
  from: Address,
  spender: Address,
  takes: readonly (readonly [Address, bigint])[],
  then: ContractCall,
  blockNumber?: bigint,
) {
  const at = blockNumber ?? (await getBlockNumber(client, { cacheTime: 0 }))
  const allowances = await Promise.all(
    takes.map(([tokenAt]) =>
      readContract(client, {
        address: tokenAt,
        abi: erc20Abi,
        functionName: 'allowance',
        args: [from, spender],
        blockNumber: at,
      }),
    ),
  )
  const approvals = takes
    .filter(([, amount], i) => (allowances[i] ?? 0n) < amount)
    .map(([tokenAt, amount]) =>
      call({ address: tokenAt, abi: erc20Abi, functionName: 'approve', args: [spender, amount] }),
    )
  return [...approvals, call(then)]
}

const account = address('The account the call acts for.')
const grants = (operator: string) =>
  `${operator} can withdraw every collateral token of every margin position and the USDG of every short to any address, borrow USDG to any address, and sell or buy back the account's shorts, until the account signs setAuthorization(${operator}, false).`
const operatorOf = (account: string) => `${checksummed(account)}, or an address it authorized`

/**
 * Registers the tools that prepare transactions. Each returns unsigned calls, never signs one: the account, or an
 * address it authorized, checks and signs them in its own wallet.
 */
export function registerTransactions(server: McpServer, context: Context) {
  const { client, deployments: d, maxSlippageBps } = context
  const chainId = d.chainId
  const marginAccounts = () => {
    const found = d.tapehouse.MarginAccounts
    if (found === undefined) throw new Error('The registry has no .tapehouse.MarginAccounts.')
    return found
  }
  const shortPositions = () => {
    const found = d.tapehouse.ShortPositions
    if (found === undefined) throw new Error('The registry has no .tapehouse.ShortPositions.')
    return found
  }
  const slippageBps = z
    .number()
    .int()
    .min(0)
    .max(maxSlippageBps)
    .describe(
      `The most the price may move against the quote, in basis points, from 0 to ${maxSlippageBps}, this server's cap. The server chooses none.`,
    )

  tool(
    server,
    context,
    'accounts_set_authorization',
    {
      title: 'Authorize an operator',
      description: `Prepares the margin accounts' setAuthorization for the account itself to sign. With allowed true, ${grants('operator')} With allowed false, it stops the operator. An agent that acts for an account needs this first.`,
      input: z.strictObject({
        operator: address("The address that may act for the account, such as an agent's own wallet."),
        allowed: z.boolean(),
      }),
      output: prepared.extend({ grants: z.string().describe('What the signed call lets the operator do.') }),
      openWorld: false,
    },
    async (args) => {
      const operator = checksummed(args.operator)
      return {
        chainId,
        signer: 'the account, which lets the operator act for it',
        grants: args.allowed ? grants(operator) : `Nothing: ${operator} can no longer act for the account.`,
        calls: [call(accounts.setAuthorization(d, operator, args.allowed))],
      }
    },
  )

  tool(
    server,
    context,
    'accounts_deposit',
    {
      title: 'Deposit collateral',
      description:
        "Prepares a deposit of a token from `from`, the account or an address it authorized, into the account's margin position, preceded by an approval where from's allowance is short. A basket's key deposits its shares, into the cross position alone, where the engine margins them as the Stock Tokens they redeem for.",
      input: z.strictObject({
        from: address('The address that sends the tokens and signs.'),
        account,
        position,
        token,
        amount: units("The amount, in the token's base units: 6 decimals for USDG, 18 for WETH and Stock Tokens."),
      }),
      output: prepared,
      openWorld: true,
    },
    async (args) => {
      const from = checksummed(args.from)
      const held = collateralAddress(d, args.token, args.position)
      const amount = BigInt(args.amount)
      const deposit = accounts.deposit(d, {
        position: positionId(args.position),
        token: held,
        amount,
        account: checksummed(args.account),
      })
      return {
        chainId,
        signer: from,
        calls: await approved(context, from, marginAccounts(), [[held, amount]], deposit),
      }
    },
  )

  tool(
    server,
    context,
    'accounts_withdraw',
    {
      title: 'Withdraw collateral',
      description:
        "Prepares a withdrawal of a token, or of a basket's shares from the cross position, from the account's margin position to receiver.",
      input: z.strictObject({
        account,
        position,
        token,
        amount: units("The amount, in the token's base units."),
        receiver: address('Who receives the tokens.'),
      }),
      output: prepared,
      openWorld: false,
    },
    async (args) => ({
      chainId,
      signer: operatorOf(args.account),
      calls: [
        call(
          accounts.withdraw(d, {
            position: positionId(args.position),
            token: collateralAddress(d, args.token, args.position),
            amount: BigInt(args.amount),
            account: checksummed(args.account),
            receiver: checksummed(args.receiver),
          }),
        ),
      ],
    }),
  )

  tool(
    server,
    context,
    'accounts_borrow',
    {
      title: 'Borrow USDG',
      description: "Prepares a borrow of USDG against the account's margin position, to receiver.",
      input: z.strictObject({
        account,
        position,
        assets: units('The USDG to borrow, in base units (6 decimals).'),
        receiver: address('Who receives the USDG.'),
      }),
      output: prepared,
      openWorld: false,
    },
    async (args) => ({
      chainId,
      signer: operatorOf(args.account),
      calls: [
        call(
          accounts.borrow(d, {
            position: positionId(args.position),
            assets: BigInt(args.assets),
            account: checksummed(args.account),
            receiver: checksummed(args.receiver),
          }),
        ),
      ],
    }),
  )

  tool(
    server,
    context,
    'accounts_repay',
    {
      title: 'Repay USDG',
      description:
        "Prepares a repayment of the account's margin position with from's USDG, preceded by an approval where its allowance is short: the debt first, then the premium. Anyone may repay; accounts_repayment reads what clears it.",
      input: z.strictObject({
        from: address('The address that pays and signs.'),
        account,
        position,
        assets: units('The USDG to repay, in base units (6 decimals).'),
      }),
      output: prepared,
      openWorld: true,
    },
    async (args) => {
      const from = checksummed(args.from)
      const assets = BigInt(args.assets)
      const repay = accounts.repay(d, {
        position: positionId(args.position),
        assets,
        account: checksummed(args.account),
      })
      return {
        chainId,
        signer: from,
        calls: await approved(context, from, marginAccounts(), [[tokenAddress(d, 'USDG'), assets]], repay),
      }
    },
  )

  tool(
    server,
    context,
    'accounts_unwrap',
    {
      title: "Unwrap a position's basket",
      description:
        "Prepares the margin accounts' unwrap of shares of a basket in the account's cross position: the basket redeems them, and the position holds their Stock Tokens in their place, at the same equity and requirement. The account or an address it authorized signs it, or anyone once the position falls short.",
      input: z.strictObject({ account, basket, shares }),
      output: prepared,
      openWorld: false,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const owner = checksummed(args.account)
      return {
        chainId,
        signer: `${owner}, an address it authorized, or anyone once the position falls short`,
        calls: [call(accounts.unwrap(d, { account: owner, basket: args.basket, shares: BigInt(args.shares) }))],
      }
    },
  )

  tool(
    server,
    context,
    'baskets_mint',
    {
      title: 'Mint basket shares',
      description:
        "Prepares a mint of a basket's shares to receiver, paid in kind by from: each Stock Token's part of what the basket holds, rounded up, as baskets_preview_mint reads it at this block. That is the mint's maxAssets, so it reverts with AboveMaximum rather than take more. Each token is approved for exactly its part where from's allowance is short.",
      input: z.strictObject({
        from: address('The address that pays the Stock Tokens and signs.'),
        basket,
        shares,
        receiver: address('Who receives the shares.'),
      }),
      output: prepared.extend({ assets: amounts.describe('What the mint takes at most of each Stock Token.') }),
      openWorld: true,
    },
    async (args) => {
      const spender = basketAddress(d, args.basket)
      const from = checksummed(args.from)
      const shares = BigInt(args.shares)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const [held, maxAssets] = await Promise.all([
        baskets.components(client, d, args.basket, { blockNumber }),
        baskets.previewMint(client, d, args.basket, shares, { blockNumber }),
      ])
      const mint = baskets.mint(d, { basket: args.basket, shares, receiver: checksummed(args.receiver), maxAssets })
      const takes = held.tokens.map((tokenAt, i) => [tokenAt, maxAssets[i] ?? 0n] as const)
      return {
        chainId,
        signer: from,
        assets: byAsset(held, maxAssets, 'amount'),
        calls: await approved(context, from, spender, takes, mint, blockNumber),
      }
    },
  )

  tool(
    server,
    context,
    'baskets_redeem',
    {
      title: 'Redeem basket shares',
      description:
        "Prepares a redemption of owner's shares of a basket for its Stock Tokens, to receiver: each token's part of what the basket holds, rounded down, as baskets_preview_redeem reads it. Shares redeem whole, in kind: a paused Stock Token or an address the issuer blocklisted stops the redemption. Signed by owner, or by an address it allowed to spend the shares.",
      input: z.strictObject({
        basket,
        shares,
        receiver: address('Who receives the Stock Tokens.'),
        owner: address('Whose shares are redeemed.'),
      }),
      output: prepared.extend({ assets: amounts.describe('What the redemption gives of each Stock Token now.') }),
      openWorld: true,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const owner = checksummed(args.owner)
      const shares = BigInt(args.shares)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const [held, assets] = await Promise.all([
        baskets.components(client, d, args.basket, { blockNumber }),
        baskets.previewRedeem(client, d, args.basket, shares, { blockNumber }),
      ])
      return {
        chainId,
        signer: `${owner}, or an address it allowed to spend the shares`,
        assets: byAsset(held, assets, 'amount'),
        calls: [call(baskets.redeem(d, { basket: args.basket, shares, receiver: checksummed(args.receiver), owner }))],
      }
    },
  )

  tool(
    server,
    context,
    'shorts_deposit',
    {
      title: 'Add USDG to a short',
      description:
        "Prepares a deposit of from's USDG into the account's short of an asset, preceded by an approval where its allowance is short. from is the account or an address it authorized.",
      input: z.strictObject({
        from: address('The address that sends the USDG and signs.'),
        account,
        asset,
        amount: units('The USDG, in base units (6 decimals).'),
      }),
      output: prepared,
      openWorld: true,
    },
    async (args) => {
      const from = checksummed(args.from)
      const amount = BigInt(args.amount)
      const deposit = shorts.deposit(d, { asset: args.asset, amount, account: checksummed(args.account) })
      return {
        chainId,
        signer: from,
        calls: await approved(context, from, shortPositions(), [[tokenAddress(d, 'USDG'), amount]], deposit),
      }
    },
  )

  tool(
    server,
    context,
    'shorts_withdraw',
    {
      title: 'Take USDG from a short',
      description: "Prepares a withdrawal of USDG from the account's short of an asset, to receiver.",
      input: z.strictObject({
        account,
        asset,
        amount: units('The USDG, in base units (6 decimals).'),
        receiver: address('Who receives the USDG.'),
      }),
      output: prepared,
      openWorld: false,
    },
    async (args) => ({
      chainId,
      signer: operatorOf(args.account),
      calls: [
        call(
          shorts.withdraw(d, {
            asset: args.asset,
            amount: BigInt(args.amount),
            account: checksummed(args.account),
            receiver: checksummed(args.receiver),
          }),
        ),
      ],
    }),
  )

  tool(
    server,
    context,
    'shorts_sell',
    {
      title: 'Sell short',
      description:
        "Prepares a short sale: the shorts borrow amount of the asset from its lending vault and sell it through its Uniswap pool for at least Uniswap QuoterV2's quote now less the slippage you state.",
      input: z.strictObject({
        account,
        asset,
        amount: units("The amount to sell, in the token's base units (18 decimals)."),
        slippageBps,
      }),
      output: prepared.extend({
        quote: z.object({ proceeds: uint, minProceeds: uint, slippageBps: z.number().int() }),
      }),
      openWorld: true,
    },
    async (args) => {
      const amount = BigInt(args.amount)
      const owner = checksummed(args.account)
      const sale = await shorts.quoteSale(client, d, {
        asset: args.asset,
        amount,
        slippageBps: BigInt(args.slippageBps),
      })
      return {
        chainId,
        signer: operatorOf(args.account),
        calls: [call(shorts.sell(d, { asset: args.asset, amount, minProceeds: sale.minProceeds, account: owner }))],
        quote: { ...sale, slippageBps: args.slippageBps },
      }
    },
  )

  tool(
    server,
    context,
    'shorts_cover',
    {
      title: 'Buy back a short',
      description:
        "Prepares a buy-back of the account's short, paid from the short's own USDG, for at most Uniswap QuoterV2's quote now plus the slippage you state. amount max buys back all it owes.",
      input: z.strictObject({
        account,
        asset,
        amount: z
          .union([z.literal('max'), units('')])
          .describe("The amount to buy back, in the token's base units (18 decimals), or max for all the short owes."),
        slippageBps,
      }),
      output: prepared.extend({
        quote: z.object({ amount: uint, cost: uint, maxCost: uint, slippageBps: z.number().int() }),
      }),
      openWorld: true,
    },
    async (args) => {
      const owner = checksummed(args.account)
      let amount = args.amount === 'max' ? 0n : BigInt(args.amount)
      if (args.amount === 'max') {
        amount = (await shorts.position(client, d, owner, args.asset)).debt
        if (amount === 0n) throw new Error(`${owner} owes no ${args.asset}.`)
      }
      const buyBack = await shorts.quoteCover(client, d, {
        asset: args.asset,
        amount,
        slippageBps: BigInt(args.slippageBps),
      })
      const cover = shorts.cover(d, {
        asset: args.asset,
        amount: args.amount === 'max' ? maxUint256 : amount,
        maxCost: buyBack.maxCost,
        account: owner,
      })
      return {
        chainId,
        signer: operatorOf(args.account),
        calls: [call(cover)],
        quote: { amount, ...buyBack, slippageBps: args.slippageBps },
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_buy',
    {
      title: 'Buy gap cover',
      description:
        "Prepares a purchase of gap cover on an asset for holder over the coming closure, paying notional times the asset's fall from its last Chainlink round before the close to its reopening price beyond deductibleBps, up to limitBps. Its maxPremium is the cover's quote at this block, so it reverts with PremiumAboveLimit rather than pay more, and from's USDG is approved for exactly that premium where its allowance is short. Refused while no cover is on sale, as the purchase would revert with SalesClosed. The cover is held by holder and cannot be transferred.",
      input: z.strictObject({
        from: address('The address that pays the premium and signs.'),
        asset,
        notional: units('The notional, in USDG base units (6 decimals).'),
        deductibleBps: bps("The fall the cover pays beyond, in basis points: at least the quote's minDeductibleBps."),
        limitBps: bps('The fall the cover pays up to, in basis points: above the deductible, at most 10000.'),
        holder: address('Who holds the cover and is credited its payout.'),
      }),
      output: prepared.extend({
        quote: z.object({ premium: uint, reserve: uint, closesMs: uint, endsMs: uint }),
      }),
      openWorld: true,
    },
    async (args) => {
      const cover = gapCoverAddress(d)
      const from = checksummed(args.from)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const sales = await gapCover.sales(client, d, { blockNumber })
      if (sales === undefined)
        throw new Error(
          `No cover is on sale at block ${blockNumber}: the purchase would revert with ${explainRevert({ errorName: 'SalesClosed', args: [] })}. gap_cover_sales reads whether it is.`,
        )
      const layer = {
        notional: BigInt(args.notional),
        deductibleBps: BigInt(args.deductibleBps),
        limitBps: BigInt(args.limitBps),
      }
      const offer = await gapCover.quote(client, d, args.asset, layer, { blockNumber })
      const purchase = gapCover.buy(d, {
        asset: args.asset,
        ...layer,
        maxPremium: offer.premium,
        holder: checksummed(args.holder),
      })
      return {
        chainId,
        signer: from,
        quote: { ...offer, ...sales },
        calls: await approved(context, from, cover, [[tokenAddress(d, 'USDG'), offer.premium]], purchase, blockNumber),
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_release',
    {
      title: 'Release a gap cover',
      description:
        "Prepares the release of cover id once its series has settled or is void, which anyone may send: it credits the holder the payout, or the premium of a void series, frees the writers' reserve and adds the premium less the payout to their USDG. The holder then claims with gap_cover_claim. Refused before the series settles, as the release would revert with NotSettled.",
      input: z.strictObject({ id: units('The cover, as Bought numbered it.') }),
      output: prepared.extend({ holder: z.string(), payout: uint, refund: uint }),
      openWorld: true,
    },
    async (args) => {
      gapCoverAddress(d)
      const id = BigInt(args.id)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const held = await gapCover.cover(client, d, id, { blockNumber })
      if (held === undefined)
        throw new Error(
          `Cover ${id} does not exist or has been released: the release would revert with ${explainRevert({ errorName: 'NoCover', args: [id] })}.`,
        )
      const asset = hexToString(held.symbol, { size: 32 })
      const series = await gapCover.series(client, d, asset, held.closesMs, { blockNumber })
      if (series.status === 'open')
        throw new Error(
          `Cover ${id}'s series has not settled: the release would revert with ${explainRevert({ errorName: 'NotSettled', args: [id] })}.`,
        )
      return {
        chainId,
        signer: 'anyone',
        holder: held.holder,
        payout: series.settlement === undefined ? 0n : gapCover.payout(held, series.settlement),
        refund: series.status === 'void' ? held.premium : 0n,
        calls: [call(gapCover.release(d, id))],
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_claim',
    {
      title: 'Claim gap cover credits',
      description:
        'Prepares the claim of the USDG credited to holder by its released covers, sent to receiver. The holder signs it. Refused while nothing is credited.',
      input: z.strictObject({
        holder: address('The holder the credits are owed to, who signs.'),
        receiver: address('Who receives the USDG.'),
      }),
      output: prepared.extend({ amount: uint.describe('The USDG the claim sends, in base units.') }),
      openWorld: true,
    },
    async (args) => {
      gapCoverAddress(d)
      const holder = checksummed(args.holder)
      const amount = await gapCover.payouts(client, d, holder)
      if (amount === 0n) throw new Error(`Nothing is credited to ${holder}.`)
      return { chainId, signer: holder, amount, calls: [call(gapCover.claim(d, checksummed(args.receiver)))] }
    },
  )

  tool(
    server,
    context,
    'gap_cover_deposit',
    {
      title: 'Write gap cover',
      description:
        "Prepares a deposit of from's USDG into the gap cover's writers' vault, for shares to receiver, preceded by an approval for exactly that amount where from's allowance is short. Every cover reserves its whole payout from the writers' USDG, and their premiums join it at release. Deposits are open only before the sales end, while every outstanding cover is over the coming closure: refused beyond what the vault takes now.",
      input: z.strictObject({
        from: address('The address that sends the USDG and signs.'),
        assets: units('The USDG, in base units (6 decimals).'),
        receiver: address('Who receives the shares.'),
      }),
      output: prepared.extend({ shares: uint.describe('The shares the deposit mints at this block, 12 decimals.') }),
      openWorld: true,
    },
    async (args) => {
      const cover = gapCoverAddress(d)
      const from = checksummed(args.from)
      const receiver = checksummed(args.receiver)
      const assets = BigInt(args.assets)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const vault = { address: cover, abi: gapCoverAbi, blockNumber } as const
      const [most, shares] = await Promise.all([
        readContract(client, { ...vault, functionName: 'maxDeposit', args: [receiver] }),
        readContract(client, { ...vault, functionName: 'previewDeposit', args: [assets] }),
      ])
      if (assets > most)
        throw new Error(
          `The vault takes at most ${most} now: the deposit would revert with ${explainRevert({ errorName: 'ERC4626ExceededMaxDeposit', args: [receiver, assets, most] })}.`,
        )
      const deposit = gapCover.deposit(d, { assets, receiver })
      return {
        chainId,
        signer: from,
        shares,
        calls: await approved(context, from, cover, [[tokenAddress(d, 'USDG'), assets]], deposit, blockNumber),
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_redeem',
    {
      title: "Redeem a writer's shares",
      description:
        "Prepares a redemption of owner's shares of the gap cover's writers' vault for USDG to receiver, signed by owner or an address it allowed to spend them. Redemptions wait until every cover is released: refused beyond what owner may redeem now.",
      input: z.strictObject({
        shares: units('The shares, in base units (12 decimals).'),
        receiver: address('Who receives the USDG.'),
        owner: address('Whose shares are redeemed.'),
      }),
      output: prepared.extend({ assets: uint.describe('The USDG the redemption gives at this block, in base units.') }),
      openWorld: true,
    },
    async (args) => {
      const cover = gapCoverAddress(d)
      const owner = checksummed(args.owner)
      const shares = BigInt(args.shares)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const vault = { address: cover, abi: gapCoverAbi, blockNumber } as const
      const [most, assets] = await Promise.all([
        readContract(client, { ...vault, functionName: 'maxRedeem', args: [owner] }),
        readContract(client, { ...vault, functionName: 'previewRedeem', args: [shares] }),
      ])
      if (shares > most)
        throw new Error(
          `${owner} may redeem at most ${most} shares now: the redemption would revert with ${explainRevert({ errorName: 'ERC4626ExceededMaxRedeem', args: [owner, shares, most] })}.`,
        )
      return {
        chainId,
        signer: `${owner}, or an address it allowed to spend the shares`,
        assets,
        calls: [call(gapCover.redeem(d, { shares, receiver: checksummed(args.receiver), owner }))],
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_measure',
    {
      title: "Measure an asset's week",
      description:
        "Prepares the measurement of an asset's week for the closure on sale, which anyone may send: the cover keeps the realised move of the trading week before the close, which its series' covers are priced at, so that the series' first buyer does not pay for reading it from the feed. Gives the move, in millionths. Refused outside the sales, as the measurement would revert with SalesClosed, and before the week's last day, a day before the close, has passed (TooEarlyToMeasure).",
      input: z.strictObject({ asset }),
      output: prepared.extend({
        closesMs: uint.describe('The close that keys the series, in milliseconds.'),
        weekMove: uint.describe(
          'The realised move of the week before the close, in millionths; 18446744073709551615 for a week the feed cannot show.',
        ),
      }),
      openWorld: true,
    },
    async (args) => {
      gapCoverAddress(d)
      const blockNumber = await getBlockNumber(client, { cacheTime: 0 })
      const sales = await gapCover.sales(client, d, { blockNumber })
      if (sales === undefined)
        throw new Error(
          `No cover is on sale at block ${blockNumber}: the measurement would revert with ${explainRevert({ errorName: 'SalesClosed', args: [] })}. gap_cover_sales reads whether it is.`,
        )
      const measure = gapCover.measure(d, args.asset)
      const { result } = await simulateContract(client, { ...measure, blockNumber })
      return { chainId, signer: 'anyone', closesMs: sales.closesMs, weekMove: result, calls: [call(measure)] }
    },
  )
}
