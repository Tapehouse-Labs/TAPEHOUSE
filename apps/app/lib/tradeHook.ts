// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { baskets, gapCover, type Layer, shorts } from "@tapehouse/sdk";
import { useQuery } from "@tanstack/react-query";
import type { Abi, Address } from "viem";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import { REFRESH_MS, useSettled } from "./hooks";
import { healthAfter, readTrade } from "./tradeRead";

/** Everything the Trade surface shows, with what `account` holds once a wallet connects, read at one block every few seconds. */
export function useTradeState(account: Address | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  return useQuery({
    queryKey: ["trade", deployments.chainId, account],
    queryFn: () => readTrade(client!, deployments, account),
    enabled: client !== undefined,
    placeholderData: (previous) => previous,
    refetchInterval: REFRESH_MS,
  });
}

/** What selling `amount` of `asset` short pays, or buying it back costs, through its pool, and the limit at `slippageBps`. */
export function usePoolQuote(side: "sell" | "buyBack", asset: string, amount: bigint | undefined, slippageBps: bigint | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  const settled = useSettled(amount);
  return useQuery({
    queryKey: ["poolQuote", deployments.chainId, side, asset, settled?.toString(), slippageBps?.toString()],
    queryFn: async () => {
      const input = { asset, amount: settled!, slippageBps: slippageBps! };
      if (side === "sell") {
        const { proceeds, minProceeds } = await shorts.quoteSale(client!, deployments, input);
        return { quoted: proceeds, limit: minProceeds };
      }
      const { cost, maxCost } = await shorts.quoteCover(client!, deployments, input);
      return { quoted: cost, limit: maxCost };
    },
    enabled: client !== undefined && settled !== undefined && slippageBps !== undefined,
    placeholderData: (previous) => previous,
    retry: false,
  });
}

/** The short's health once `call` runs from `account`, simulated at the latest block. */
export function useHealthAfter(
  account: Address | undefined,
  asset: string,
  call: { address: Address; abi: Abi; functionName: string; args: readonly unknown[] } | undefined,
) {
  const deployments = useDeployments();
  const client = usePublicClient();
  const settled = useSettled(call);
  return useQuery({
    queryKey: [
      "healthAfter",
      deployments.chainId,
      account,
      asset,
      settled ? JSON.stringify([settled.address, settled.functionName, settled.args], big) : "",
    ],
    queryFn: () => healthAfter(client!, deployments, account!, asset, settled!),
    enabled: client !== undefined && account !== undefined && settled !== undefined,
    placeholderData: (previous) => previous,
  });
}

/** Gap cover's premium on `asset` for `layer`, the reserve it takes, the smallest deductible and the gap it prices at. */
export function useCoverQuote(asset: string, layer: Layer | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  const settled = useSettled(layer);
  return useQuery({
    queryKey: ["coverQuote", deployments.chainId, asset, settled ? JSON.stringify(settled, big) : ""],
    queryFn: async () => {
      const [quote, minDeductible, pricing] = await Promise.all([
        gapCover.quote(client!, deployments, asset, settled!),
        gapCover.minDeductible(client!, deployments, asset),
        gapCover.pricingGap(client!, deployments, asset),
      ]);
      return { ...quote, minDeductible, ...pricing };
    },
    enabled: client !== undefined && settled !== undefined,
    placeholderData: (previous) => previous,
    refetchInterval: REFRESH_MS,
    retry: false,
  });
}

/** What minting `shares` of `basket` takes of each Stock Token, rounded up, or redeeming them returns, rounded down. */
export function useBasketPreview(side: "mint" | "redeem", basket: string, shares: bigint | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  const settled = useSettled(shares);
  return useQuery({
    queryKey: ["basketPreview", deployments.chainId, side, basket, settled?.toString()],
    refetchInterval: REFRESH_MS,
    queryFn: () =>
      side === "mint"
        ? baskets.previewMint(client!, deployments, basket, settled!)
        : baskets.previewRedeem(client!, deployments, basket, settled!),
    enabled: client !== undefined && settled !== undefined,
    placeholderData: (previous) => previous,
  });
}

function big(_: string, value: unknown) {
  return typeof value === "bigint" ? value.toString() : value;
}
