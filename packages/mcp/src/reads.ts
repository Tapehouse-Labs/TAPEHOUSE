// SPDX-License-Identifier: MIT OR Apache-2.0
import type { McpServer } from '@modelcontextprotocol/server'
import {
  accounts,
  band,
  baskets,
  gapCover,
  gapCoverAbi,
  morpho,
  sharePriceFeed,
  shorts,
  toBytes32,
  tokenPriceFeed,
} from '@tapehouse/sdk'
import { hexToString } from 'viem'
import { getBlockNumber, readContract } from 'viem/actions'
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
  int,
  position,
  positionId,
  shares,
  token,
  tool,
  uint,
  units,
} from './tool.js'

const states = ['halted', 'degraded', 'closed', 'open'] as const
const sessions = ['unknown', 'closed', 'open'] as const
const block = { blockNumber: uint.describe('The block every value was read at.') }
const quote = {
  state: z.enum(states),
  live: z.number().int().describe('How many legs are live.'),
  mid: uint.describe("The band's centre, in USD with 8 decimals."),
  halfBps: uint.describe('Its half-width, in basis points.'),
  low: uint.describe('Its low edge, in USD with 8 decimals: what lending values the token at.'),
  high: uint.describe('Its high edge, in USD with 8 decimals: what a short is measured at.'),
}
const halt = {
  signedHalt: z.boolean().describe("Whether a halt signed by Tapehouse's halt signer holds."),
  until: uint.describe('Until when, in seconds.'),
  issuedAt: uint.describe('When its last message was issued, in seconds.'),
  oraclePaused: z.boolean().describe("Whether the issuer has paused the Stock Token's oracle."),
}
const health = {
  equity: int.describe('Equity net of debt and premium, in USD with 18 decimals.'),
  requirement: uint.describe('The requirement, in USD with 18 decimals.'),
  surplus: int.describe('Equity less the requirement, in USD with 18 decimals: negative falls short.'),
  missing: z.number().int().describe("The engine's bits for assets whose liquidity is missing."),
  regime: z.number().int().describe("The engine's regime."),
}
const account = address('The account.')
const target = z.array(z.object({ asset: z.string(), token: z.string(), units: uint }))
const statuses = ['open', 'settled', 'void'] as const
const series = {
  notional: uint.describe('The notional the series covers, in USDG base units (6 decimals).'),
  status: z.enum(statuses).describe('open until it settles; void refunds every premium.'),
  referencePrice: uint
    .nullable()
    .describe("Once settled, the reference, the feed's last round before the sales ended."),
  price: uint
    .nullable()
    .describe('Once settled, the reopening price: the first round after the reopen, or the median.'),
  flagged: z.boolean().describe("Whether it settled on the median of the band's centres rather than the first round."),
  reopenMs: uint.describe('When the closure reopens, in milliseconds; 0 while not recorded.'),
}

async function at({ client }: Context) {
  return { blockNumber: await getBlockNumber(client, { cacheTime: 0 }) }
}

async function seriesOf(context: Context, asset: string, closesMs: bigint, read: { blockNumber: bigint }) {
  const s = await gapCover.series(context.client, context.deployments, asset, closesMs, read)
  return {
    notional: s.notional,
    status: s.status,
    referencePrice: s.settlement?.referencePrice ?? null,
    price: s.settlement?.price ?? null,
    flagged: s.flagged,
    reopenMs: s.reopenMs,
  }
}

/**
 * Registers the tools that read the band, its feeds, Chainlink, the Morpho oracles, the baskets, the accounts, the
 * shorts and the gap cover.
 */
export function registerReads(server: McpServer, context: Context) {
  const { client, deployments: d } = context

  tool(
    server,
    context,
    'registry',
    {
      title: 'Registry',
      description: "The chain's registry: every contract, token, basket, feed and oracle address the other tools use.",
      input: z.strictObject({}),
      output: z.object({
        chainId: z.number().int(),
        tokens: z.record(z.string(), z.string()),
        chainlink: z.record(z.string(), z.object({ address: z.string(), kind: z.enum(['token', 'share']) })),
        bandFeeds: z.record(z.string(), z.string()),
        tapehouse: z.record(z.string(), z.string()),
        stockLending: z.record(z.string(), z.string()),
        baskets: z.record(z.string(), z.string()),
        uniswapV3: z.record(z.string(), z.string()),
        morpho: z.record(z.string(), z.string()),
        morphoMarkets: z.record(z.string(), z.string()),
        morphoOracles: z.record(z.string(), z.string()),
      }),
      openWorld: false,
    },
    async () => d,
  )

  tool(
    server,
    context,
    'band_quote',
    {
      title: "An asset's band",
      description:
        "The asset's price band now: its state, live legs, centre, half-width and edges. All zero while it is halted or has no live leg.",
      input: z.strictObject({ asset }),
      output: z.object({ ...block, asset: z.string(), ...quote }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const { state, ...rest } = await band.quote(client, d, args.asset, read)
      return { ...read, asset: args.asset, state: states[state], ...rest }
    },
  )

  tool(
    server,
    context,
    'band_session',
    {
      title: 'Market session',
      description:
        "Chainlink's 24/5 session the band follows, NYSE's state and next state, when NYSE changes state, and the session's next boundary, in milliseconds.",
      input: z.strictObject({}),
      output: z.object({
        ...block,
        state: z.enum(sessions),
        nyse: z.number().int(),
        nyseNext: z.number().int(),
        changeMs: uint,
        boundaryMs: uint.describe('The next session boundary, in milliseconds: a reopen while closed.'),
      }),
      openWorld: true,
    },
    async () => {
      const read = await at(context)
      const { state, ...rest } = await band.session(client, d, read)
      return { ...read, state: sessions[state], ...rest }
    },
  )

  tool(
    server,
    context,
    'band_halt',
    {
      title: "An asset's trading halt",
      description:
        "Whether a halt signed by Tapehouse's halt signer holds for the asset, the one input the band takes on Tapehouse's own signature, and whether the issuer has paused the Stock Token's oracle.",
      input: z.strictObject({ asset }),
      output: z.object({ ...block, asset: z.string(), ...halt }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return { ...read, asset: args.asset, ...(await band.halt(client, d, args.asset, read)) }
    },
  )

  tool(
    server,
    context,
    'band_latest_band',
    {
      title: "An asset's band feed",
      description:
        "The whole band from the asset's BandFeed, the Chainlink-compatible feed of its low side: the band, its variance, the session, both halt sources, the sequencer, and the Uniswap pool's TWAP and premium.",
      input: z.strictObject({ asset }),
      output: z.object({
        ...block,
        asset: z.string(),
        ...quote,
        variance: uint,
        session: z.number().int(),
        nyse: z.number().int(),
        nyseNext: z.number().int(),
        nyseChangeMs: uint,
        sessionBoundaryMs: uint,
        signedHalt: z.boolean(),
        haltUntil: uint,
        haltIssuedAt: uint,
        oraclePaused: z.boolean(),
        sequencerSettled: z.boolean(),
        twapValid: z.boolean(),
        twap: uint,
        premiumBps: int,
      }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const { state, ...rest } = await band.latestBand(client, d, args.asset, read)
      return { ...read, asset: args.asset, state: states[state], ...rest }
    },
  )

  tool(
    server,
    context,
    'band_sealed',
    {
      title: "An asset's sealed band",
      description:
        "The band the asset's BandFeed sealed in the ten minutes before the reopen at reopenMs, and when; all zero where none was sealed.",
      input: z.strictObject({ asset, reopenMs: units('The reopen, in milliseconds, as band_session reports it.') }),
      output: z.object({ ...block, asset: z.string(), ...quote, sealedAt: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const { state, ...rest } = await band.sealed(client, d, args.asset, BigInt(args.reopenMs), read)
      return { ...read, asset: args.asset, state: states[state], ...rest }
    },
  )

  tool(
    server,
    context,
    'chainlink_round',
    {
      title: 'A Chainlink round',
      description:
        "The latest round of a Chainlink feed of the registry's .chainlink. kind is what the caller expects it to price: the Stock Token, its multiplier included, on Robinhood Chain, the share on Arbitrum One; the other is refused.",
      input: z.strictObject({
        feed: z
          .string()
          .regex(/^[A-Za-z0-9_]{1,64}$/, 'not a feed name')
          .describe("The feed's key in .chainlink, such as NVDA_USD."),
        kind: z.enum(['token', 'share']),
      }),
      output: z.object({
        ...block,
        feed: z.string(),
        kind: z.enum(['token', 'share']),
        roundId: uint,
        answer: int,
        startedAt: uint,
        updatedAt: uint,
        answeredInRound: uint,
      }),
      openWorld: true,
    },
    async (args) => {
      const feed = args.kind === 'token' ? tokenPriceFeed(d, args.feed) : sharePriceFeed(d, args.feed)
      const read = await at(context)
      return { ...read, feed: args.feed, ...(await band.chainlinkRound(client, feed, read)) }
    },
  )

  tool(
    server,
    context,
    'morpho_oracle',
    {
      title: "An asset's Morpho oracle",
      description:
        "The price the asset's Morpho Blue oracle answers, the band's low edge in Morpho's scale, or why it has none: a stale band (no live leg, which a fresh writePrices ends), a halt, or an unsettled sequencer. Never a price of zero. With the band it reads, its tokens, its scale, its owner and the asset's halt.",
      input: z.strictObject({ asset }),
      output: z.object({
        ...block,
        asset: z.string(),
        address: z.string(),
        band: z.string(),
        symbol: z.string(),
        collateralToken: z.string(),
        loanToken: z.string(),
        scaleFactor: uint,
        owner: z.string().describe('Who may re-point the oracle to another band.'),
        halt: z.object(halt),
        price: uint
          .nullable()
          .describe(
            'One whole collateral token in loan-token units with 36 decimals more than the collateral has, or null.',
          ),
        noPrice: z.enum(['stale', 'halted', 'sequencerNotSettled']).nullable(),
        revert: z.string().nullable(),
      }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const [settings, haltNow, answer] = await Promise.all([
        morpho.oracle(client, d, args.asset, read),
        morpho.halt(client, d, args.asset, read),
        morpho.price(client, d, args.asset, read),
      ])
      return {
        ...read,
        asset: args.asset,
        ...settings,
        symbol: settings.symbol === toBytes32(args.asset) ? args.asset : settings.symbol,
        halt: haltNow,
        price: answer.price ?? null,
        noPrice: answer.price === undefined ? answer.noPrice : null,
        revert: answer.price === undefined ? explainRevert(answer.revert) : null,
      }
    },
  )

  tool(
    server,
    context,
    'baskets_components',
    {
      title: "A basket's components",
      description:
        "The Stock Tokens a basket holds, by asset, in the order of every list of its amounts, and the basket's address, the ERC-20 of its shares. A basket is minted and redeemed in kind, never at a price.",
      input: z.strictObject({ basket }),
      output: z.object({
        ...block,
        basket: z.string(),
        address: z.string(),
        components: z.array(z.object({ asset: z.string(), token: z.string() })),
      }),
      openWorld: true,
    },
    async (args) => {
      const contract = basketAddress(d, args.basket)
      const read = await at(context)
      const { assets, tokens } = await baskets.components(client, d, args.basket, read)
      return {
        ...read,
        basket: args.basket,
        address: contract,
        components: assets.map((asset, i) => ({ asset, token: tokens[i] })),
      }
    },
  )

  tool(
    server,
    context,
    'baskets_preview_mint',
    {
      title: 'What a mint takes',
      description:
        "What minting shares of a basket takes of each Stock Token now, in the token's base units: its part of what the basket holds, rounded up, or of the target while the basket has no shares. baskets_mint takes exactly this as its limit.",
      input: z.strictObject({ basket, shares }),
      output: z.object({ ...block, basket: z.string(), shares: uint, assets: amounts }),
      openWorld: true,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const read = await at(context)
      const [held, assets] = await Promise.all([
        baskets.components(client, d, args.basket, read),
        baskets.previewMint(client, d, args.basket, BigInt(args.shares), read),
      ])
      return { ...read, basket: args.basket, shares: args.shares, assets: byAsset(held, assets, 'amount') }
    },
  )

  tool(
    server,
    context,
    'baskets_preview_redeem',
    {
      title: 'What a redemption gives',
      description:
        "What redeeming shares of a basket gives of each Stock Token now, in the token's base units: its part of what the basket holds, rounded down.",
      input: z.strictObject({ basket, shares }),
      output: z.object({ ...block, basket: z.string(), shares: uint, assets: amounts }),
      openWorld: true,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const read = await at(context)
      const [held, assets] = await Promise.all([
        baskets.components(client, d, args.basket, read),
        baskets.previewRedeem(client, d, args.basket, BigInt(args.shares), read),
      ])
      return { ...read, basket: args.basket, shares: args.shares, assets: byAsset(held, assets, 'amount') }
    },
  )

  tool(
    server,
    context,
    'baskets_target',
    {
      title: "A basket's target",
      description:
        "The target in effect of a basket: base units of each Stock Token per share (10^18 of the share's base units). Anyone may rebalance the basket toward it, only when what comes in at the bands' low edges is worth at least what goes out at their high edges.",
      input: z.strictObject({ basket }),
      output: z.object({ ...block, basket: z.string(), target }),
      openWorld: true,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const read = await at(context)
      const [held, units] = await Promise.all([
        baskets.components(client, d, args.basket, read),
        baskets.target(client, d, args.basket, read),
      ])
      return { ...read, basket: args.basket, target: byAsset(held, units, 'units') }
    },
  )

  tool(
    server,
    context,
    'baskets_pending_target',
    {
      title: "A basket's pending target",
      description:
        "A target the basket's owner proposed and that is not yet in effect, as baskets_target gives it, and when it takes effect, in seconds: seven days after it was proposed, so a holder who disagrees can redeem in kind first. Empty and 0 where none is pending.",
      input: z.strictObject({ basket }),
      output: z.object({ ...block, basket: z.string(), target, effectiveAt: uint }),
      openWorld: true,
    },
    async (args) => {
      basketAddress(d, args.basket)
      const read = await at(context)
      const [held, pending] = await Promise.all([
        baskets.components(client, d, args.basket, read),
        baskets.pendingTarget(client, d, args.basket, read),
      ])
      return {
        ...read,
        basket: args.basket,
        target: byAsset(held, pending.units, 'units'),
        effectiveAt: pending.effectiveAt,
      }
    },
  )

  tool(
    server,
    context,
    'accounts_health',
    {
      title: "A margin position's health",
      description:
        "An account's margin position as the risk engine sees it now: equity, requirement and their difference. A halted asset counts for nothing. The surplus is room over the requirement in USD, not an amount of USDG that can be borrowed; see accounts_liquidation_price's `borrowing`.",
      input: z.strictObject({ account, position }),
      output: z.object({ ...block, ...health }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const h = await accounts.health(client, d, checksummed(args.account), positionId(args.position), read)
      return {
        ...read,
        equity: h.equity,
        requirement: h.requirement,
        surplus: h.equity - h.requirement,
        missing: h.missing,
        regime: h.regime,
      }
    },
  )

  tool(
    server,
    context,
    'accounts_in_baskets',
    {
      title: 'What a margin position holds through its baskets',
      description:
        "What an account's margin position holds of each Stock Token through its baskets' shares, by asset, in the token's base units, as the shares would redeem now, rounded down. The engine margins the shares as these Stock Tokens, beside what the position holds directly. Only the cross position holds baskets.",
      input: z.strictObject({ account, position }),
      output: z.object({ ...block, assets: z.record(z.string(), uint) }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return {
        ...read,
        assets: await accounts.inBaskets(client, d, checksummed(args.account), positionId(args.position), read),
      }
    },
  )

  tool(
    server,
    context,
    'accounts_collateral',
    {
      title: "A margin position's collateral",
      description:
        "What an account's margin position holds of a token, in the token's base units, or of a basket's shares in its cross position.",
      input: z.strictObject({ account, position, token }),
      output: z.object({ ...block, token: z.string(), amount: uint }),
      openWorld: true,
    },
    async (args) => {
      const held = collateralAddress(d, args.token, args.position)
      const read = await at(context)
      const amount = await accounts.collateral(
        client,
        d,
        checksummed(args.account),
        positionId(args.position),
        held,
        read,
      )
      return { ...read, token: args.token, amount }
    },
  )

  tool(
    server,
    context,
    'accounts_leverage',
    {
      title: "A margin position's leverage",
      description:
        "The gross exposure of an account's margin position over its equity, in basis points; 2^256 - 1 for a position without equity.",
      input: z.strictObject({ account, position }),
      output: z.object({ ...block, leverageBps: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return {
        ...read,
        leverageBps: await accounts.leverage(client, d, checksummed(args.account), positionId(args.position), read),
      }
    },
  )

  tool(
    server,
    context,
    'accounts_liquidation_price',
    {
      title: "A margin position's liquidation price",
      description:
        "The price of an asset, in USD with 8 decimals, at or above which the position meets its requirement once it owes borrowing more USDG: the band's low edge where it falls short already, zero where no price leaves it short or the asset is halted.",
      input: z.strictObject({
        account,
        position,
        asset,
        borrowing: units('More USDG the position would owe, in base units (6 decimals); 0 for none.'),
      }),
      output: z.object({ ...block, asset: z.string(), price: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const price = await accounts.liquidationPrice(
        client,
        d,
        checksummed(args.account),
        positionId(args.position),
        args.asset,
        BigInt(args.borrowing),
        read,
      )
      return { ...read, asset: args.asset, price }
    },
  )

  tool(
    server,
    context,
    'accounts_repayment',
    {
      title: "A margin position's repayment",
      description:
        'What repays all a margin position owes at this block: its debt and its premium, in USDG base units (6 decimals). A repayment pays the debt first; premium keeps accruing until it is sent.',
      input: z.strictObject({ account, position }),
      output: z.object({ ...block, debt: uint, premium: uint, assets: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return {
        ...read,
        ...(await accounts.repayment(client, d, checksummed(args.account), positionId(args.position), read)),
      }
    },
  )

  tool(
    server,
    context,
    'accounts_is_authorized',
    {
      title: 'An authorization',
      description:
        "Whether operator may act for account: borrow, withdraw and deposit in the margin accounts, and trade the account's shorts.",
      input: z.strictObject({ account, operator: address('The address that would act for the account.') }),
      output: z.object({ ...block, authorized: z.boolean() }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return {
        ...read,
        authorized: await accounts.isAuthorized(client, d, checksummed(args.account), checksummed(args.operator), read),
      }
    },
  )

  tool(
    server,
    context,
    'shorts_position',
    {
      title: 'A short position',
      description:
        "An account's short of an asset: its USDG once it pays its part of every buy-in, negative if it cannot, what it owes in the token's base units, its borrow shares and its book.",
      input: z.strictObject({ account, asset }),
      output: z.object({ ...block, asset: z.string(), usdgHeld: int, debt: uint, shares: uint, epoch: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return {
        ...read,
        asset: args.asset,
        ...(await shorts.position(client, d, checksummed(args.account), args.asset, read)),
      }
    },
  )

  tool(
    server,
    context,
    'shorts_health',
    {
      title: "A short's health",
      description: "An account's short of an asset at the band's high edge: equity, requirement and their difference.",
      input: z.strictObject({ account, asset }),
      output: z.object({ ...block, asset: z.string(), ...health }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const h = await shorts.health(client, d, checksummed(args.account), args.asset, read)
      return {
        ...read,
        asset: args.asset,
        equity: h.equity,
        requirement: h.requirement,
        surplus: h.equity - h.requirement,
        missing: h.missing,
        regime: h.regime,
      }
    },
  )

  tool(
    server,
    context,
    'shorts_restriction',
    {
      title: "An asset's short-sale restriction",
      description:
        "Whether a short of the asset must sell at no less than the band's centre now, after a 10% fall from its last close, and that close in USD with 8 decimals.",
      input: z.strictObject({ asset }),
      output: z.object({ ...block, asset: z.string(), restricted: z.boolean(), close: uint }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      return { ...read, asset: args.asset, ...(await shorts.restriction(client, d, args.asset, read)) }
    },
  )

  tool(
    server,
    context,
    'gap_cover_sales',
    {
      title: 'Gap cover sales',
      description:
        "Whether gap cover is on sale now: on the last trading day before a weekend or holiday, from the band's session showing the close until the close itself, the 24/5 session's 20:00 ET or, where the feeds follow NYSE's regular hours, the regular close. With the close that keys the series on sale and when the sales end, in milliseconds. A purchase outside them reverts with SalesClosed.",
      input: z.strictObject({}),
      output: z.object({
        ...block,
        onSale: z.boolean(),
        closesMs: uint.describe('The close that keys the series on sale; 0 while none is.'),
        endsMs: uint.describe('When the sales end; 0 while none is on sale.'),
      }),
      openWorld: true,
    },
    async () => {
      const read = await at(context)
      const sales = await gapCover.sales(client, d, read)
      return { ...read, onSale: sales !== undefined, closesMs: sales?.closesMs ?? 0n, endsMs: sales?.endsMs ?? 0n }
    },
  )

  tool(
    server,
    context,
    'gap_cover_quote',
    {
      title: 'A gap cover quote',
      description:
        "The premium of cover on an asset paying notional times its fall over the coming closure beyond deductibleBps, up to limitBps, as the cover prices it now, and the writers' USDG it reserves. It is priced at the weekend gap the cover reads: twice the week's realised move on the asset's Chainlink feed, at least 55% of the engine's gap and at most eight times it; outside the sales at the engine's alone. The smallest deductible starts the fitted tail, raised where the feed's last answer stands above the band's centre.",
      input: z.strictObject({
        asset,
        notional: units('The notional, in USDG base units (6 decimals).'),
        deductibleBps: bps('The fall the cover pays beyond, in basis points: at least minDeductibleBps.'),
        limitBps: bps('The fall the cover pays up to, in basis points: above the deductible, at most 10000.'),
      }),
      output: z.object({
        ...block,
        asset: z.string(),
        premium: uint.describe('The premium, in USDG base units: what gap_cover_buy pays at most.'),
        reserve: uint.describe("What the cover reserves of the writers' USDG: notional times the layer, rounded up."),
        minDeductibleBps: uint,
        gap: uint.describe('The weekend gap the cover prices at, in millionths.'),
        weekMove: uint.describe(
          "The realised move of the trading week before the close on sale, in millionths; 0 outside the sales, 18446744073709551615 for a week the feed cannot show, which prices at eight times the engine's gap.",
        ),
        capacity: uint.describe("The writers' USDG free to reserve."),
        onSale: z.boolean(),
      }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const layer = {
        notional: BigInt(args.notional),
        deductibleBps: BigInt(args.deductibleBps),
        limitBps: BigInt(args.limitBps),
      }
      const [offer, minimum, priced, capacity, sales] = await Promise.all([
        gapCover.quote(client, d, args.asset, layer, read),
        gapCover.minDeductible(client, d, args.asset, read),
        gapCover.pricingGap(client, d, args.asset, read),
        gapCover.capacity(client, d, read),
        gapCover.sales(client, d, read),
      ])
      return {
        ...read,
        asset: args.asset,
        ...offer,
        minDeductibleBps: minimum,
        ...priced,
        capacity,
        onSale: sales !== undefined,
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_series',
    {
      title: 'A gap cover series',
      description:
        "An asset's gap cover series over the closure keyed by closesMs: the notional it covers, whether it is open, settled at its reference and reopening prices, or void and refunding, and when the closure reopens.",
      input: z.strictObject({ asset, closesMs: units('The close that keys the series, as gap_cover_sales names it.') }),
      output: z.object({ ...block, asset: z.string(), closesMs: uint, ...series }),
      openWorld: true,
    },
    async (args) => {
      gapCoverAddress(d)
      const read = await at(context)
      const closesMs = BigInt(args.closesMs)
      return { ...read, asset: args.asset, closesMs, ...(await seriesOf(context, args.asset, closesMs, read)) }
    },
  )

  tool(
    server,
    context,
    'gap_cover_position',
    {
      title: 'A gap cover',
      description:
        'Cover id: its holder, asset, closure, notional, layer and premium, its series, and what releasing it credits the holder once the series has settled or is void: the payout, or the premium back. With the USDG credited to the holder and not yet claimed. A released cover no longer exists.',
      input: z.strictObject({ id: units('The cover, as Bought numbered it.') }),
      output: z.object({
        ...block,
        id: uint,
        holder: z.string(),
        asset: z.string(),
        closesMs: uint,
        notional: uint,
        deductibleBps: uint,
        limitBps: uint,
        premium: uint,
        series: z.object(series),
        payout: uint.nullable().describe('What its release pays once its series has settled.'),
        refund: uint.nullable().describe('What its release refunds once its series is void: its premium.'),
        credited: uint.describe('The USDG credited to the holder and not yet claimed, in base units.'),
      }),
      openWorld: true,
    },
    async (args) => {
      const read = await at(context)
      const id = BigInt(args.id)
      const held = await gapCover.cover(client, d, id, read)
      if (held === undefined) throw new Error(`Cover ${id} does not exist or has been released.`)
      const asset = hexToString(held.symbol, { size: 32 })
      const [state, credited] = await Promise.all([
        seriesOf(context, asset, held.closesMs, read),
        gapCover.payouts(client, d, held.holder, read),
      ])
      return {
        ...read,
        id,
        holder: held.holder,
        asset,
        closesMs: held.closesMs,
        notional: held.notional,
        deductibleBps: held.deductibleBps,
        limitBps: held.limitBps,
        premium: held.premium,
        series: state,
        payout:
          state.referencePrice !== null && state.price !== null
            ? gapCover.payout(held, { referencePrice: state.referencePrice, price: state.price })
            : null,
        refund: state.status === 'void' ? held.premium : null,
        credited,
      }
    },
  )

  tool(
    server,
    context,
    'gap_cover_vault',
    {
      title: "The gap cover writers' vault",
      description:
        "The ERC-4626 vault of the gap cover's writers: the USDG their shares stand for, what is reserved for outstanding covers, the premiums held until release, the credits not yet claimed, and how many covers are outstanding. With owner, its shares, what they stand for, and how much it may deposit and redeem now: deposits only before the sales end while every outstanding cover is over the coming closure, redemptions only once every cover is released.",
      input: z.strictObject({ owner: address('A writer whose shares and limits to read.').optional() }),
      output: z.object({
        ...block,
        held: uint.describe("The writers' USDG, reserved or not, in base units."),
        reserved: uint,
        capacity: uint,
        premiums: uint,
        owed: uint.describe('Payouts and refunds credited to holders and not yet claimed.'),
        outstanding: uint,
        totalShares: uint.describe('Shares have 12 decimals.'),
        writer: z
          .object({ owner: z.string(), shares: uint, assets: uint, maxDeposit: uint, maxRedeem: uint })
          .nullable(),
      }),
      openWorld: true,
    },
    async (args) => {
      const cover = gapCoverAddress(d)
      const read = await at(context)
      const vault = { address: cover, abi: gapCoverAbi, ...read } as const
      const [held, reserved, capacity, premiums, owed, outstanding, totalShares] = await Promise.all([
        readContract(client, { ...vault, functionName: 'held' }),
        readContract(client, { ...vault, functionName: 'reserved' }),
        readContract(client, { ...vault, functionName: 'capacity' }),
        readContract(client, { ...vault, functionName: 'premiums' }),
        readContract(client, { ...vault, functionName: 'owed' }),
        readContract(client, { ...vault, functionName: 'outstanding' }),
        readContract(client, { ...vault, functionName: 'totalSupply' }),
      ])
      let writer = null
      if (args.owner !== undefined) {
        const owner = checksummed(args.owner)
        const [shares, maxDeposit, maxRedeem] = await Promise.all([
          readContract(client, { ...vault, functionName: 'balanceOf', args: [owner] }),
          readContract(client, { ...vault, functionName: 'maxDeposit', args: [owner] }),
          readContract(client, { ...vault, functionName: 'maxRedeem', args: [owner] }),
        ])
        const assets = await readContract(client, { ...vault, functionName: 'convertToAssets', args: [shares] })
        writer = { owner, shares, assets, maxDeposit, maxRedeem }
      }
      return { ...read, held, reserved, capacity, premiums, owed, outstanding, totalShares, writer }
    },
  )
}
