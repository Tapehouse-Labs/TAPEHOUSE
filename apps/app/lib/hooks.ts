// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { band } from "@tapehouse/sdk";
import { type MutationState, useMutationState, useQueries, useQuery } from "@tanstack/react-query";
import { useEffect, useState, useSyncExternalStore } from "react";
import { type Connector, useClient, useConnect, useConnection, useConnectors, useDisconnect, useSwitchChain } from "wagmi";
import { useDeployments } from "@/app/providers";
import { walletState } from "./wallet";

/** How often the app reads the chain again, in milliseconds. */
export const REFRESH_MS = 15_000;

const subscribeNever = () => () => {};

/** False while the server renders and the page hydrates, true once the browser runs it. */
export function useMounted(): boolean {
  return useSyncExternalStore(
    subscribeNever,
    () => true,
    () => false,
  );
}

/** The wallets the browser offers: those announced through EIP-6963, else the one at `window.ethereum`. */
function useWallets(): readonly Connector[] {
  const connectors = useConnectors();
  const mounted = useMounted();
  if (!mounted) return [];
  const announced = connectors.filter((c) => c.id !== "injected");
  if (announced.length > 0) return announced;
  return (window as { ethereum?: unknown }).ethereum ? connectors.filter((c) => c.id === "injected") : [];
}

/** The latest of wagmi's mutations under `key`, whichever component started it. */
function useLatest(key: string): MutationState<unknown, Error> | undefined {
  return useMutationState({
    filters: { mutationKey: [key] },
    select: (mutation) => mutation.state as MutationState<unknown, Error>,
  }).at(-1);
}

/** The visitor's wallet, where it stands against the app's chain, and what can be done about it. */
export function useWallet() {
  const deployments = useDeployments();
  const connection = useConnection();
  const wallets = useWallets();
  const mounted = useMounted();
  const connect = useConnect();
  const switchChain = useSwitchChain();
  const disconnect = useDisconnect();
  const connecting = useLatest("connect");
  const switching = useLatest("switchChain");
  const failed = [connecting, switching]
    .filter((m) => m?.status === "error")
    .sort((a, b) => b!.submittedAt - a!.submittedAt)[0];
  const state = mounted
    ? walletState({
        status: connection.status,
        chainId: connection.chainId,
        expected: deployments.chainId,
        wallets: wallets.length,
      })
    : "connecting";
  return {
    state,
    address: connection.address,
    chainId: connection.chainId,
    wallets,
    connect: (connector: Connector) => connect.mutate({ connector, chainId: deployments.chainId }),
    switchChain: () => switchChain.mutate({ chainId: deployments.chainId }),
    disconnect: () => disconnect.mutate(),
    pending: connecting?.status === "pending" || switching?.status === "pending",
    error: state === "connected" ? undefined : (failed?.error ?? undefined),
  };
}

/** The clock, to the minute. */
export function useMinute(): number {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), 60_000);
    return () => window.clearInterval(id);
  }, []);
  return now;
}

/** Chainlink's 24/5 session as the band reads it. */
export function useSession() {
  const deployments = useDeployments();
  const client = useClient();
  return useQuery({
    queryKey: ["session", deployments.chainId],
    queryFn: () => band.session(client!, deployments),
    enabled: client !== undefined,
    refetchInterval: REFRESH_MS,
  });
}

/** The band of every asset with a band feed on the app's chain. */
export function useBands() {
  const deployments = useDeployments();
  const client = useClient();
  const assets = Object.keys(deployments.bandFeeds);
  const queries = useQueries({
    queries: assets.map((asset) => ({
      queryKey: ["quote", deployments.chainId, asset],
      queryFn: () => band.quote(client!, deployments, asset),
      enabled: client !== undefined,
      refetchInterval: REFRESH_MS,
    })),
  });
  return assets.map((asset, i) => ({ asset, query: queries[i] }));
}
