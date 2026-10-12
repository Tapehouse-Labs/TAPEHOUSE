// SPDX-License-Identifier: MIT OR Apache-2.0
import { auction, type Deployments } from "@tapehouse/sdk";
import { type Address, type Client, erc20Abi } from "viem";
import { getBlock, readContract } from "viem/actions";
import { type BidStatus, bidStatus, type SavedBid, stage } from "./auction";

export type RoundView = { asset: string; openMs: bigint; round: Awaited<ReturnType<typeof auction.round>> };

/** A bid of the visitor's in a round: what is saved of it, what the auction holds, and where it stands. */
export type MyBid = SavedBid & { status: BidStatus; index: number | undefined };

/** A revealed bid of the visitor's that the auction lists, saved in this browser or not. */
export type Listed = { asset: string; openMs: bigint; index: number; quantity: bigint; price: bigint; claimed: boolean };

export type AuctionState = {
  blockNumber: bigint;
  nowMs: number;
  constants: Awaited<ReturnType<typeof auction.constants>>;
  phase: { openMs: bigint; revealing: boolean };
  lastOpenMs: bigint;
  now: ReturnType<typeof stage>;
  assets: string[];
  rounds: RoundView[];
  /** The floor of each asset's round before anything has opened it, from the band the round will take it from. */
  previews: Record<string, bigint>;
  wallet?: { usdg: bigint; allowance: bigint; saved: MyBid[]; listed: Listed[] };
};

const key = (asset: string, openMs: bigint) => `${asset}:${openMs}`;

/**
 * The reopening auction at one block: its constants, the round that is live or the last one, each asset's state in it,
 * and, once a wallet is connected, what it holds of the bond's token, what it allowed the auction to take, its saved
 * bids against what the auction keeps of them, and every bid of its that the rounds list.
 */
export async function readAuction(
  client: Client,
  deployments: Deployments,
  wallet?: { account: Address; saved: readonly SavedBid[] },
): Promise<AuctionState> {
  const block = await getBlock(client, { blockTag: "latest" });
  const at = { blockNumber: block.number };
  const nowMs = Number(block.timestamp) * 1000;
  const [constants, phase, lastOpenMs] = await Promise.all([
    auction.constants(client, deployments, at),
    auction.phase(client, deployments, at),
    auction.lastOpenMs(client, deployments, at),
  ]);
  const assets = Object.keys(deployments.bandFeeds);
  const shown = phase.openMs !== 0n ? phase.openMs : lastOpenMs;
  const wanted = new Map<string, { asset: string; openMs: bigint }>();
  if (shown !== 0n) for (const asset of assets) wanted.set(key(asset, shown), { asset, openMs: shown });
  for (const bid of wallet?.saved ?? []) wanted.set(key(bid.asset, bid.openMs), { asset: bid.asset, openMs: bid.openMs });
  const rounds: RoundView[] = await Promise.all(
    [...wanted.values()].map(async (w) => ({ ...w, round: await auction.round(client, deployments, w, at) })),
  );
  const previews: Record<string, bigint> = {};
  if (phase.openMs !== 0n) {
    await Promise.all(
      assets
        .filter((asset) => rounds.find((r) => r.asset === asset && r.openMs === shown)?.round.floor === 0n)
        .map(async (asset) => {
          previews[asset] = await auction.previewFloor(client, deployments, { asset, openMs: shown }, at).catch(() => 0n);
        }),
    );
  }
  const current = rounds.find((r) => r.openMs === shown);
  const state: AuctionState = {
    blockNumber: block.number,
    nowMs,
    constants,
    phase,
    lastOpenMs,
    now: stage({ phase, lastOpenMs, round: current?.round, nowMs, clearMs: constants.clearMs }),
    assets,
    rounds,
    previews,
  };
  if (!wallet) return state;
  const usdg = deployments.tokens.USDG as Address;
  const address = deployments.tapehouse.ReopeningAuction as Address;
  const [balance, allowance, saved] = await Promise.all([
    readContract(client, { address: usdg, abi: erc20Abi, functionName: "balanceOf", args: [wallet.account], blockNumber: block.number }),
    readContract(client, { address: usdg, abi: erc20Abi, functionName: "allowance", args: [wallet.account, address], blockNumber: block.number }),
    Promise.all(
      wallet.saved.map(async (bid): Promise<MyBid> => {
        const { deposit } = await auction.commitmentOf(client, deployments, bid, at);
        const view = rounds.find((r) => r.asset === bid.asset && r.openMs === bid.openMs);
        const index = view?.round.bids.findIndex(
          (b) => b.bidder.toLowerCase() === wallet.account.toLowerCase() && b.quantity === bid.quantity && b.price === bid.price,
        );
        const found = index === undefined || index < 0 ? undefined : index;
        const listed = found === undefined ? undefined : view?.round.bids[found];
        return { ...bid, status: bidStatus({ deposit, listed }), index: found };
      }),
    ),
  ]);
  const listed = rounds.flatMap((r) =>
    r.round.bids.flatMap((b, index) =>
      b.bidder.toLowerCase() === wallet.account.toLowerCase()
        ? [{ asset: r.asset, openMs: r.openMs, index, quantity: b.quantity, price: b.price, claimed: b.claimed }]
        : [],
    ),
  );
  return { ...state, wallet: { usdg: balance, allowance, saved, listed } };
}
