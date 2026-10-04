// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { type Deployments, parseDeployments } from "@tapehouse/sdk";

export type Registry = {
  deployments: Deployments;
  rpcUrl: string | undefined;
  bundlerUrl: string | undefined;
  sponsorUrl: string | undefined;
};

/**
 * The address registry the app runs on, read when the app is built: `TAPEHOUSE_REGISTRY` names a registry file,
 * otherwise `deployments/<TAPEHOUSE_CHAIN_ID>.json`, Robinhood Chain testnet by default. `TAPEHOUSE_RPC_URL`
 * replaces the chain's public RPC; `TAPEHOUSE_BUNDLER_URL` names the ERC-4337 bundler smart accounts send through, and
 * `TAPEHOUSE_SPONSOR_URL` Tapehouse's sponsor service, which signs the operations the sponsor paymaster pays for.
 */
export async function loadRegistry(env: Record<string, string | undefined>): Promise<Registry> {
  const file = env.TAPEHOUSE_REGISTRY ?? `../../deployments/${env.TAPEHOUSE_CHAIN_ID ?? "46630"}.json`;
  const json: unknown = JSON.parse(await readFile(resolve(/*turbopackIgnore: true*/ process.cwd(), file), "utf8"));
  return {
    deployments: parseDeployments(json),
    rpcUrl: env.TAPEHOUSE_RPC_URL || undefined,
    bundlerUrl: env.TAPEHOUSE_BUNDLER_URL || undefined,
    sponsorUrl: env.TAPEHOUSE_SPONSOR_URL || undefined,
  };
}
