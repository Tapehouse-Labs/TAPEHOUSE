// SPDX-License-Identifier: MIT OR Apache-2.0
import { backstop, band, type Deployments, reopeningAuctionAbi, toBytes32 } from "@tapehouse/sdk";
import type { Client } from "viem";
import { getBlock, readContract } from "viem/actions";

export type Round = {
  asset: string;
  openMs: bigint;
  round: Awaited<ReturnType<typeof readRound>>;
};

export type RiskState = {
  blockNumber: bigint;
  bands: { asset: string; band: Awaited<ReturnType<typeof band.latestBand>> }[];
  backstop?: Awaited<ReturnType<typeof backstop.state>>;
  auction?: { phase: { openMs: bigint; revealing: boolean }; rounds: Round[] };
};

async function readRound(client: Client, address: `0x${string}`, asset: string, openMs: bigint, blockNumber: bigint) {
  const [round, lots] = await Promise.all([
    readContract(client, { address, abi: reopeningAuctionAbi, functionName: "round", args: [toBytes32(asset), openMs], blockNumber }),
    readContract(client, { address, abi: reopeningAuctionAbi, functionName: "lots", args: [toBytes32(asset), openMs], blockNumber }),
  ]);
  return { ...round, lots: lots.length };
}

/**
 * The risk surface at one block: each asset's band with its TWAP and premium, the backstop's exposure limits, and
 * the latest reopening round of each asset, each only where the registry deploys it.
 */
export async function readRisk(client: Client, deployments: Deployments): Promise<RiskState> {
  const block = await getBlock(client, { blockTag: "latest" });
  const at = { blockNumber: block.number };
  const { GapBackstop, ReopeningAuction: auction } = deployments.tapehouse;
  const assets = Object.keys(deployments.bandFeeds);
  const [bands, backstopState, auctionState] = await Promise.all([
    Promise.all(assets.map(async (asset) => ({ asset, band: await band.latestBand(client, deployments, asset, at) }))),
    GapBackstop ? backstop.state(client, deployments, at) : undefined,
    auction
      ? Promise.all([
          readContract(client, { address: auction, abi: reopeningAuctionAbi, functionName: "phase", blockNumber: block.number }),
          readContract(client, { address: auction, abi: reopeningAuctionAbi, functionName: "lastOpenMs", blockNumber: block.number }),
        ]).then(async ([[openMs, revealing], last]) => ({
          phase: { openMs, revealing },
          rounds: await Promise.all(
            assets.map(async (asset) => ({ asset, openMs: last, round: await readRound(client, auction, asset, last, block.number) })),
          ),
        }))
      : undefined,
  ]);
  return { blockNumber: block.number, bands, backstop: backstopState, auction: auctionState };
}
