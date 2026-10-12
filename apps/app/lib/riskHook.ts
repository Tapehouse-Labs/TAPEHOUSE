// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useQuery } from "@tanstack/react-query";
import { usePublicClient } from "wagmi";
import { useDeployments, useIndexerUrl } from "@/app/providers";
import { backstopActions, clearedRounds, fetchEvents } from "./history";
import { REFRESH_MS } from "./hooks";
import { readRisk } from "./riskRead";

/** Every band, the backstop's exposure and the latest reopening rounds, read at one block every few seconds. */
export function useRiskState() {
  const deployments = useDeployments();
  const client = usePublicClient();
  return useQuery({
    queryKey: ["risk", deployments.chainId],
    queryFn: () => readRisk(client!, deployments),
    enabled: client !== undefined,
    placeholderData: (previous) => previous,
    refetchInterval: REFRESH_MS,
  });
}

/** What the indexer holds of the auction's cleared rounds and the backstop's covers, once the app has an indexer. */
export function useRecord() {
  const indexerUrl = useIndexerUrl();
  return useQuery({
    queryKey: ["record", indexerUrl],
    queryFn: async () => {
      const [rounds, covers, bought] = await Promise.all([
        fetchEvents(indexerUrl!, { contract: "tapehouse.ReopeningAuction", event: "Cleared" }),
        fetchEvents(indexerUrl!, { contract: "tapehouse.GapBackstop", event: "Covered" }),
        fetchEvents(indexerUrl!, { contract: "tapehouse.GapBackstop", event: "RemainderBought" }),
      ]);
      return { rounds: clearedRounds(rounds), actions: backstopActions([...covers, ...bought]) };
    },
    enabled: indexerUrl !== undefined,
    refetchInterval: REFRESH_MS * 4,
  });
}
