// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import type { Deployments } from "@tapehouse/sdk";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createContext, use, useState } from "react";
import { createConfig, http, injected, WagmiProvider } from "wagmi";
import { chainOf } from "@/lib/chains";

const DeploymentsContext = createContext<Deployments | null>(null);
const SmartAccountsContext = createContext<{ bundlerUrl?: string; sponsorUrl?: string }>({});

/** The registry of the chain the app runs on. */
export function useDeployments(): Deployments {
  const deployments = use(DeploymentsContext);
  if (!deployments) throw new Error("useDeployments is used outside Providers.");
  return deployments;
}

/** The ERC-4337 bundler smart accounts send their operations through, and the sponsor service, when the app has them. */
export function useSmartAccountUrls(): { bundlerUrl?: string; sponsorUrl?: string } {
  return use(SmartAccountsContext);
}

export function Providers({
  deployments,
  rpcUrl,
  bundlerUrl,
  sponsorUrl,
  children,
}: {
  deployments: Deployments;
  rpcUrl: string | undefined;
  bundlerUrl: string | undefined;
  sponsorUrl: string | undefined;
  children: React.ReactNode;
}) {
  const [config] = useState(() => {
    const chain = chainOf(deployments.chainId, rpcUrl);
    return createConfig({
      chains: [chain],
      connectors: [injected()],
      transports: { [chain.id]: http(undefined, { batch: true }) },
      ssr: true,
    });
  });
  const [queryClient] = useState(() => new QueryClient());
  const [urls] = useState(() => ({ bundlerUrl, sponsorUrl }));
  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <DeploymentsContext value={deployments}>
          <SmartAccountsContext value={urls}>{children}</SmartAccountsContext>
        </DeploymentsContext>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
