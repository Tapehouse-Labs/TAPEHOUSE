// SPDX-License-Identifier: MIT OR Apache-2.0
import { accounts, band, type Deployments, engine } from "@tapehouse/sdk";
import type { Client } from "viem";
import { getBlockNumber } from "viem/actions";

/** The engine's assets and their bands' low edges, the accounts' limits and the session, at one block. */
export async function readMarket(client: Client, deployments: Deployments) {
  const blockNumber = await getBlockNumber(client, { cacheTime: 0 });
  const at = { blockNumber };
  const assets = await accounts.stocks(client, deployments, at);
  const [quotes, limits, session] = await Promise.all([
    Promise.all(assets.map(({ asset }) => band.quote(client, deployments, asset, at))),
    accounts.limits(client, deployments, at),
    band.session(client, deployments, at),
  ]);
  return { blockNumber, assets: assets.map(({ asset }, i) => ({ asset, low: quotes[i]!.low, state: quotes[i]!.state })), limits, session };
}

/**
 * What the engine requires of `quantities`, in the order of its assets, at `prices`: now, over two days without and
 * with the weekend-gap scenarios, and the regime, as `currentRequirement` and `requirement` read them at one block.
 */
export async function simulate(client: Client, deployments: Deployments, quantities: bigint[], prices: bigint[], blockNumber: bigint) {
  const portfolio = { quantities, prices };
  const at = { blockNumber };
  const [current, open, closed] = await Promise.all([
    engine.currentRequirement(client, deployments, portfolio, at),
    engine.requirement(client, deployments, portfolio, { spansClosure: false }, at),
    engine.requirement(client, deployments, portfolio, { spansClosure: true }, at),
  ]);
  return { current: current.margin, regime: current.regime, missing: current.missing, open: open.margin, closed: closed.margin };
}
