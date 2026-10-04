// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useQuery } from "@tanstack/react-query";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
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
