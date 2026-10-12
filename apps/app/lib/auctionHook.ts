// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { auction } from "@tapehouse/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState, useSyncExternalStore } from "react";
import { type Address, erc20Abi, type Hex } from "viem";
import { useConnection, usePublicClient, useWalletClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import { deserialize, nextNonce, type SavedBid, saltFrom, saltMessage, serialize } from "./auction";
import { readAuction } from "./auctionRead";
import { isSanctioned } from "./compliance";
import { REFRESH_MS } from "./hooks";
import { MessageError, type Msg, msg } from "./messages";

const listeners = new Set<() => void>();
const storageKey = (chainId: number, account: Address) => `tapehouse.auction.${chainId}.${account.toLowerCase()}`;

function read(key: string): string | null {
  try {
    return window.localStorage.getItem(key);
  } catch {
    return null;
  }
}

function write(key: string, bids: readonly SavedBid[]) {
  try {
    window.localStorage.setItem(key, serialize(bids));
  } catch {}
  for (const listener of listeners) listener();
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  window.addEventListener("storage", listener);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", listener);
  };
}

const NONE: readonly SavedBid[] = [];
const cache = new Map<string, { text: string | null; bids: readonly SavedBid[] }>();

/** The bids this browser saved for `account`, which the wallet's signature can derive the secret of again. */
export function useSavedBids(account: Address | undefined): readonly SavedBid[] {
  const { chainId } = useDeployments();
  const key = account ? storageKey(chainId, account) : undefined;
  return useSyncExternalStore(
    subscribe,
    () => {
      if (!key) return NONE;
      const text = read(key);
      const cached = cache.get(key);
      if (cached && cached.text === text) return cached.bids;
      const bids = deserialize(text);
      cache.set(key, { text, bids });
      return bids;
    },
    () => NONE,
  );
}

/** The reopening auction at one block, with the wallet's bids once it is connected, read every few seconds. */
export function useAuctionState() {
  const deployments = useDeployments();
  const client = usePublicClient();
  const { address } = useConnection();
  const saved = useSavedBids(address);
  return useQuery({
    queryKey: ["auction", deployments.chainId, address, saved.map((b) => b.commitment)],
    queryFn: () => readAuction(client!, deployments, address ? { account: address, saved } : undefined),
    enabled: client !== undefined && deployments.tapehouse.ReopeningAuction !== undefined,
    placeholderData: (previous) => previous,
    refetchInterval: REFRESH_MS,
  });
}

export type AuctionPhase = "salt" | "approve" | "checking" | "signing" | "sending";

/** What the primary button says while an auction action is in each phase. */
export const AUCTION_PHASE: Record<AuctionPhase, Msg> = {
  salt: msg("auction.phase.salt"),
  approve: msg("auction.phase.approve"),
  checking: msg("common.phase.checking"),
  signing: msg("common.phase.signing"),
  sending: msg("common.phase.sending"),
};

type Action =
  | { kind: "commit"; asset: string; openMs: bigint; quantity: bigint; price: bigint; deposit: bigint }
  | { kind: "reveal"; bid: SavedBid }
  | { kind: "claim"; asset: string; openMs: bigint; index: number }
  | { kind: "clear"; asset: string; openMs: bigint; price: bigint }
  | { kind: "forfeit"; bid: SavedBid };

/** Sends the auction's transactions from the wallet itself: the bidder is whoever signs, and the tokens come back to it. */
export function useAuctionSend() {
  const deployments = useDeployments();
  const client = usePublicClient();
  const { data: wallet } = useWalletClient();
  const { address } = useConnection();
  const queryClient = useQueryClient();
  const saved = useSavedBids(address);
  const [phase, setPhase] = useState<AuctionPhase>();
  const mutation = useMutation({
    mutationKey: ["auction-send"],
    mutationFn: async (action: Action) => {
      if (!client || !wallet || !address) throw new MessageError(msg("errors.connectFirst"));
      if (isSanctioned(address)) throw new MessageError(msg("errors.sanctioned"));
      const key = storageKey(deployments.chainId, address);
      const send = async (transaction: { address: Address; abi: never; functionName: never; args?: never }) => {
        setPhase("checking");
        const { request } = await client.simulateContract({ ...transaction, account: address } as never);
        setPhase("signing");
        const hash = await wallet.writeContract(request as never);
        setPhase("sending");
        const receipt = await client.waitForTransactionReceipt({ hash });
        if (receipt.status !== "success") throw new MessageError(msg("errors.transactionReverted"));
      };
      switch (action.kind) {
        case "commit": {
          const nonce = nextNonce(saved, action.asset, action.openMs);
          setPhase("salt");
          const signature = await wallet.signMessage({
            account: address,
            message: saltMessage({ asset: action.asset, openMs: action.openMs, nonce }),
          });
          const salt: Hex = saltFrom(signature);
          const commitment = auction.commitment({ bidder: address, ...action, salt });
          const bid: SavedBid = { ...action, nonce, salt, commitment };
          const usdg = deployments.tokens.USDG as Address;
          const spender = deployments.tapehouse.ReopeningAuction as Address;
          const allowance = await client.readContract({ address: usdg, abi: erc20Abi, functionName: "allowance", args: [address, spender] });
          if (allowance < action.deposit) {
            setPhase("approve");
            await send({ address: usdg, abi: erc20Abi as never, functionName: "approve" as never, args: [spender, action.deposit] as never });
          }
          write(key, [...saved, bid]);
          try {
            await send(auction.commit(deployments, { asset: action.asset, commitment, deposit: action.deposit }) as never);
          } catch (error) {
            write(key, saved);
            throw error;
          }
          return;
        }
        case "reveal":
          return send(auction.reveal(deployments, action.bid) as never);
        case "claim":
          return send(auction.claim(deployments, { ...action, index: BigInt(action.index) }) as never);
        case "clear":
          return send(auction.clear(deployments, action) as never);
        case "forfeit":
          return send(auction.forfeit(deployments, { asset: action.bid.asset, openMs: action.bid.openMs, commitment: action.bid.commitment }) as never);
      }
    },
    onSettled: async () => {
      setPhase(undefined);
      await queryClient.invalidateQueries();
    },
  });
  const run = useCallback((action: Action) => mutation.mutate(action), [mutation]);
  return { send: run, phase, pending: mutation.isPending, succeeded: mutation.isSuccess, error: mutation.error ?? undefined, reset: mutation.reset } as const;
}
