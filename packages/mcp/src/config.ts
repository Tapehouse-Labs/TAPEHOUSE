// SPDX-License-Identifier: MIT OR Apache-2.0

/** The server's configuration, all of it from the environment: it takes no key. */
export type Config = {
  /** The chain's JSON-RPC endpoint, `TAPEHOUSE_RPC_URL`, which may carry an API key and is never shown to a client. */
  rpcUrl: string
  /** The path of the chain's registry, `TAPEHOUSE_DEPLOYMENTS`: `deployments/<chainId>.json`. */
  deployments: string
  /** The most slippage a sale or buy-back may state, `TAPEHOUSE_MAX_SLIPPAGE_BPS`: 100 basis points unless set. */
  maxSlippageBps: number
  /** The most tool calls the server answers in a minute, `TAPEHOUSE_MAX_CALLS_PER_MINUTE`: 120 unless set. */
  maxCallsPerMinute: number
}

/** Reads the configuration from `env`, refusing what the server cannot use. */
export function configFromEnv(env: Record<string, string | undefined>): Config {
  const rpcUrl = env.TAPEHOUSE_RPC_URL
  if (!rpcUrl) throw new Error("Set TAPEHOUSE_RPC_URL to the chain's JSON-RPC endpoint.")
  if (!/^https?:$/.test(URL.parse(rpcUrl)?.protocol ?? ''))
    throw new Error('TAPEHOUSE_RPC_URL is not an http or https URL.')
  const deployments = env.TAPEHOUSE_DEPLOYMENTS
  if (!deployments) throw new Error("Set TAPEHOUSE_DEPLOYMENTS to the path of the chain's deployments/<chainId>.json.")
  const cap = env.TAPEHOUSE_MAX_SLIPPAGE_BPS ?? '100'
  if (!/^\d{1,5}$/.test(cap) || Number(cap) > 10_000)
    throw new Error('TAPEHOUSE_MAX_SLIPPAGE_BPS is not a whole number of basis points from 0 to 10000.')
  const calls = env.TAPEHOUSE_MAX_CALLS_PER_MINUTE ?? '120'
  if (!/^[1-9]\d{0,5}$/.test(calls) || Number(calls) > 100_000)
    throw new Error('TAPEHOUSE_MAX_CALLS_PER_MINUTE is not a whole number of calls from 1 to 100000.')
  return { rpcUrl, deployments, maxSlippageBps: Number(cap), maxCallsPerMinute: Number(calls) }
}
