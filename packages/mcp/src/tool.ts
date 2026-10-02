// SPDX-License-Identifier: MIT OR Apache-2.0
import type { CallToolResult, McpServer, ToolCallback } from '@modelcontextprotocol/server'
import { CROSS, decodeRevert, type Deployments, type Revert, toBytes32 } from '@tapehouse/sdk'
import { type Address, BaseError, type Client, getAddress, type Hex, isAddress, maxUint256 } from 'viem'
import * as z from 'zod'

/** What every tool reads from: the chain's client, its registry, the slippage cap, and the limit on calls. */
export type Context = { client: Client; deployments: Deployments; maxSlippageBps: number; rateLimit: RateLimit }

/**
 * A sliding window of a minute over every tool call the process answers. It lives in the context, not in a server,
 * because the stateless HTTP handler builds a server for each request.
 */
export class RateLimit {
  readonly perMinute: number
  readonly #taken: number[] = []

  constructor(perMinute: number) {
    this.perMinute = perMinute
  }

  /** Takes one call at `now`, in milliseconds, or returns false where the last minute already holds `perMinute`. */
  take(now = Date.now()): boolean {
    const since = now - 60_000
    while ((this.#taken[0] ?? Infinity) <= since) this.#taken.shift()
    if (this.#taken.length >= this.perMinute) return false
    this.#taken.push(now)
    return true
  }
}

const name = /^[A-Za-z0-9._-]{1,32}$/

/** An address, which must carry its checksum where it is in mixed case. */
export const address = (description: string) =>
  z
    .string()
    .regex(/^0x[0-9a-fA-F]{40}$/, { error: 'not a 20-byte hex address', abort: true })
    .refine((value) => isAddress(value), 'its checksum is wrong')
    .describe(description)

/** A whole number of a token's base units, or of a contract's own units, as a decimal string. */
export const units = (description: string) =>
  z
    .string()
    .regex(/^(0|[1-9]\d{0,77})$/, { error: 'not a whole number in decimal', abort: true })
    .refine((value) => BigInt(value) <= maxUint256, 'above 2^256 - 1')
    .describe(description)

/** An asset's symbol as the band and the registry name it. */
export const asset = z
  .string()
  .regex(name, 'not a symbol of up to 32 letters, digits, dots, dashes or underscores')
  .describe("The asset's symbol as the band names it, such as NVDA or SPY.")

/** A margin account's position: the cross position or an asset's isolated one. */
export const position = z
  .string()
  .regex(name, 'not CROSS or a symbol')
  .describe("CROSS for the account's cross position, or an asset's symbol, such as NVDA, for its isolated position.")

/** A token of the registry, by its key in `.tokens`, or a basket's shares, by its key in `.tapehouse.Baskets`. */
export const token = z
  .string()
  .regex(name, 'not a symbol')
  .describe(
    "The token's key in the registry's .tokens: USDG, WETH or a Stock Token's symbol; or a basket's key in .tapehouse.Baskets, such as PAIR, for its shares, which the cross position alone holds.",
  )

/** A basket of the registry, by its key in `.tapehouse.Baskets`. */
export const basket = z
  .string()
  .regex(name, 'not a basket key')
  .describe("The basket's key in the registry's .tapehouse.Baskets, such as PAIR.")

/** A share of a price, in basis points from 1 to 10,000. */
export const bps = (description: string) => z.number().int().min(1).max(10_000).describe(description)

/** A decimal string of a whole number. */
export const uint = z.string().regex(/^\d+$/)

/** A decimal string of a whole number that may be negative. */
export const int = z.string().regex(/^-?\d+$/)

/** Shares of a basket. */
export const shares = units('The shares, in base units (18 decimals): 10^18 is one share.')

/** A basket's amounts of each Stock Token, by asset, in the token's base units. */
export const amounts = z.array(z.object({ asset: z.string(), token: z.string(), amount: uint }))

/** The `position` argument of the margin accounts. */
export function positionId(value: string): Hex {
  return value === 'CROSS' ? CROSS : toBytes32(value)
}

/** The address of `symbol` in the registry's `.tokens`. */
export function tokenAddress(deployments: Deployments, symbol: string): Address {
  const found = Object.hasOwn(deployments.tokens, symbol) ? deployments.tokens[symbol] : undefined
  if (found === undefined) throw new Error(`The registry has no .tokens.${symbol}.`)
  return found
}

/** The address of the basket `key` in the registry's `.tapehouse.Baskets`. */
export function basketAddress(deployments: Deployments, key: string): Address {
  const found = Object.hasOwn(deployments.baskets, key) ? deployments.baskets[key] : undefined
  if (found === undefined) throw new Error(`The registry has no .tapehouse.Baskets.${key}.`)
  return found
}

/** The registry's gap cover, `.tapehouse.GapCover`. */
export function gapCoverAddress(deployments: Deployments): Address {
  const found = Object.hasOwn(deployments.tapehouse, 'GapCover') ? deployments.tapehouse.GapCover : undefined
  if (found === undefined) throw new Error('The registry has no .tapehouse.GapCover.')
  return found
}

/**
 * The address of the collateral `key` of the margin accounts' `position`: a token of `.tokens`, or the shares of a
 * basket of `.tapehouse.Baskets`, which the cross position alone holds.
 */
export function collateralAddress(deployments: Deployments, key: string, position: string): Address {
  if (Object.hasOwn(deployments.tokens, key) || !Object.hasOwn(deployments.baskets, key))
    return tokenAddress(deployments, key)
  if (position !== 'CROSS')
    throw new Error("A basket's shares are held in the cross position alone: position must be CROSS.")
  return basketAddress(deployments, key)
}

/** A basket's amounts of each Stock Token, in its order, each with its asset and token. */
export function byAsset<Key extends string>(
  components: { assets: readonly string[]; tokens: readonly Address[] },
  amounts: readonly bigint[],
  key: Key,
) {
  return components.assets.flatMap((asset, i) => {
    const amount = amounts[i]
    return amount === undefined ? [] : [{ asset, token: components.tokens[i], [key]: amount }]
  })
}

/** A checksummed address from a validated input. */
export function checksummed(value: string): Address {
  return getAddress(value)
}

/** A revert as a client reads it: its name and arguments, with any text from the chain quoted and cut short. */
export function explainRevert({ errorName, args }: Revert): string {
  const shown = args.map((arg) => {
    if (typeof arg === 'string' && !/^0x[0-9a-fA-F]*$/.test(arg))
      return JSON.stringify(arg.replace(/\p{Cc}/gu, '').slice(0, 200))
    return String(arg)
  })
  return `${errorName}(${shown.join(', ')})`
}

const clean = (text: string) => text.replace(/\p{Cc}/gu, ' ').slice(0, 300)

function explain(error: unknown): string {
  const revert = decodeRevert(error)
  if (revert) return `The call reverted with ${explainRevert(revert)}.`
  if (error instanceof BaseError) return clean(error.shortMessage)
  if (error instanceof Error) return clean(error.message)
  return 'The call failed.'
}

function json(value: unknown): Record<string, unknown> {
  return JSON.parse(JSON.stringify(value, (_, inner) => (typeof inner === 'bigint' ? inner.toString() : inner)))
}

/**
 * Registers a tool whose result is structured content, its bigints as decimal strings, and the same JSON as text. A
 * call beyond the context's rate limit is refused before it reads anything. A failure is a tool error: a revert by name
 * and arguments, any other error by its short message, stripped of control characters and cut to 300 characters,
 * never by the endpoint it came from.
 */
export function tool<Input extends z.ZodObject>(
  server: McpServer,
  { rateLimit }: Context,
  name: string,
  config: {
    title: string
    description: string
    input: Input
    output: z.ZodObject
    openWorld: boolean
  },
  run: (args: z.output<Input>) => Promise<object>,
) {
  const handler = async (args: z.output<Input>): Promise<CallToolResult> => {
    if (!rateLimit.take())
      return {
        content: [{ type: 'text', text: `At most ${rateLimit.perMinute} tool calls a minute: try again later.` }],
        isError: true,
      }
    try {
      const result = json(await run(args))
      return { content: [{ type: 'text', text: JSON.stringify(result) }], structuredContent: result }
    } catch (error) {
      return { content: [{ type: 'text', text: explain(error) }], isError: true }
    }
  }
  server.registerTool(
    name,
    {
      title: config.title,
      description: config.description,
      inputSchema: config.input,
      outputSchema: config.output,
      annotations: { readOnlyHint: true, openWorldHint: config.openWorld },
    },
    handler as ToolCallback<Input>,
  )
}
