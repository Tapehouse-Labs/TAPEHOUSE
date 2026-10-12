// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { type Deployments, parseDeployments } from "@tapehouse/sdk";

export type Registry = {
  deployments: Deployments;
  rpcUrl: string | undefined;
  bundlerUrl: string | undefined;
  sponsorUrl: string | undefined;
  indexerUrl: string | undefined;
};

const BUNDLED: Record<string, () => Promise<unknown>> = {
  "46630": () => import("../../../deployments/46630.json").then((m) => m.default),
  "4663": () => import("../../../deployments/4663.json").then((m) => m.default),
  "42161": () => import("../../../deployments/42161.json").then((m) => m.default),
};

/**
 * The address registry the app runs on, bundled with the app: `deployments/<TAPEHOUSE_CHAIN_ID>.json`, Robinhood Chain
 * testnet by default, or the registry file `TAPEHOUSE_REGISTRY` names, read from disk. `TAPEHOUSE_RPC_URL`
 * replaces the chain's public RPC; `TAPEHOUSE_BUNDLER_URL` names the ERC-4337 bundler smart accounts send through, and
 * `TAPEHOUSE_SPONSOR_URL` Tapehouse's sponsor service, which signs the operations the sponsor paymaster pays for, and
 * `TAPEHOUSE_INDEXER_URL` the indexer whose public API serves the history the chain's logs hold.
 */
export async function loadRegistry(env: Record<string, string | undefined>): Promise<Registry> {
  const chain = env.TAPEHOUSE_CHAIN_ID ?? "46630";
  const bundled = BUNDLED[chain];
  if (!env.TAPEHOUSE_REGISTRY && !bundled) throw new Error(`There is no registry for chain ${chain}.`);
  const json: unknown = env.TAPEHOUSE_REGISTRY
    ? JSON.parse(await readFile(resolve(/*turbopackIgnore: true*/ process.cwd(), env.TAPEHOUSE_REGISTRY), "utf8"))
    : await bundled?.();
  return {
    deployments: parseDeployments(json),
    rpcUrl: env.TAPEHOUSE_RPC_URL || undefined,
    bundlerUrl: env.TAPEHOUSE_BUNDLER_URL || undefined,
    sponsorUrl: env.TAPEHOUSE_SPONSOR_URL || undefined,
    indexerUrl: env.TAPEHOUSE_INDEXER_URL || undefined,
  };
}
