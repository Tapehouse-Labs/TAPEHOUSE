// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useQuery } from "@tanstack/react-query";
import type { Address } from "viem";
import { usePublicClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import { readEarn } from "./earnRead";
import { REFRESH_MS } from "./hooks";

/** Everything the Earn surface shows, with what `account` holds once a wallet connects, read at one block every few seconds. */
export function useEarnState(account: Address | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  return useQuery({
    queryKey: ["earn", deployments.chainId, account],
    queryFn: () => readEarn(client!, deployments, account),
    enabled: client !== undefined,
    placeholderData: (previous) => previous,
    refetchInterval: REFRESH_MS,
  });
}
