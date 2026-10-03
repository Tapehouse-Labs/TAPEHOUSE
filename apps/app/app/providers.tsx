// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import type { Deployments } from "@tapehouse/sdk";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createContext, use, useState } from "react";
import { createConfig, http, injected, WagmiProvider } from "wagmi";
import { chainOf } from "@/lib/chains";

const DeploymentsContext = createContext<Deployments | null>(null);

/** The registry of the chain the app runs on. */
export function useDeployments(): Deployments {
  const deployments = use(DeploymentsContext);
  if (!deployments) throw new Error("useDeployments is used outside Providers.");
  return deployments;
}

export function Providers({
  deployments,
  rpcUrl,
  children,
}: {
  deployments: Deployments;
  rpcUrl: string | undefined;
  children: React.ReactNode;
}) {
  const [config] = useState(() => {
    const chain = chainOf(deployments.chainId, rpcUrl);
    return createConfig({
      chains: [chain],
      connectors: [injected()],
      transports: { [chain.id]: http() },
      ssr: true,
    });
  });
  const [queryClient] = useState(() => new QueryClient());
  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <DeploymentsContext value={deployments}>{children}</DeploymentsContext>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
