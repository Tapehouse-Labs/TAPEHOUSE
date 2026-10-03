// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { type Deployments, parseDeployments } from "@tapehouse/sdk";

export type Registry = { deployments: Deployments; rpcUrl: string | undefined };

/**
 * The address registry the app runs on, read when the app is built: `TAPEHOUSE_REGISTRY` names a registry file,
 * otherwise `deployments/<TAPEHOUSE_CHAIN_ID>.json`, Robinhood Chain testnet by default. `TAPEHOUSE_RPC_URL`
 * replaces the chain's public RPC.
 */
export async function loadRegistry(env: Record<string, string | undefined>): Promise<Registry> {
  const file = env.TAPEHOUSE_REGISTRY ?? `../../deployments/${env.TAPEHOUSE_CHAIN_ID ?? "46630"}.json`;
  const json: unknown = JSON.parse(await readFile(resolve(/*turbopackIgnore: true*/ process.cwd(), file), "utf8"));
  return { deployments: parseDeployments(json), rpcUrl: env.TAPEHOUSE_RPC_URL || undefined };
}
