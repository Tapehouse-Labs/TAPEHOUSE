// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { sponsorship } from "@tapehouse/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toSimpleSmartAccount } from "permissionless/accounts";
import { useState } from "react";
import { createWalletClient, custom, type EIP1193Provider, http } from "viem";
import { createBundlerClient, createPaymasterClient } from "viem/account-abstraction";
import { useConnection, usePublicClient } from "wagmi";
import { useDeployments, useSmartAccountUrls } from "@/app/providers";
import { REFRESH_MS } from "./hooks";

/**
 * The visitor's smart account: a SimpleAccount v0.7 their wallet owns, whether it exists yet, and how many of its
 * operations the sponsor paymaster still pays for. `create` deploys it in a sponsored user operation that makes no call.
 * Tapehouse's sponsor service signs it over ERC-7677. Off when the app has no bundler or sponsor service, or the
 * registry no sponsor paymaster.
 */
export function useSmartAccount() {
  const deployments = useDeployments();
  const { bundlerUrl, sponsorUrl } = useSmartAccountUrls();
  const client = usePublicClient();
  const { connector, address: owner } = useConnection();
  const queryClient = useQueryClient();
  const [submitted, setSubmitted] = useState(false);
  const enabled =
    bundlerUrl !== undefined &&
    sponsorUrl !== undefined &&
    Object.hasOwn(deployments.tapehouse, "SponsorPaymaster") &&
    Object.hasOwn(deployments.erc4337, "EntryPoint") &&
    Object.hasOwn(deployments.erc4337, "SimpleAccountFactory");
  const account = useQuery({
    queryKey: ["smartAccount", deployments.chainId, connector?.uid, owner],
    queryFn: async () => {
      const provider = (await connector!.getProvider()) as EIP1193Provider;
      return toSimpleSmartAccount({
        client: client!,
        owner: createWalletClient({ account: owner!, transport: custom(provider) }),
        factoryAddress: sponsorship.accountFactory(deployments),
        entryPoint: sponsorship.entryPoint(deployments),
      });
    },
    enabled: enabled && client !== undefined && connector !== undefined && owner !== undefined,
    staleTime: Infinity,
    structuralSharing: false,
  });
  const address = account.data?.address;
  const standing = useQuery({
    queryKey: ["smartAccountStanding", deployments.chainId, address],
    queryFn: async () => {
      const [code, left] = await Promise.all([
        client!.getCode({ address: address! }),
        sponsorship.freeOperationsLeft(client!, deployments, address!),
      ]);
      return { deployed: code !== undefined, left };
    },
    enabled: address !== undefined,
    refetchInterval: REFRESH_MS,
  });
  const create = useMutation({
    mutationKey: ["createSmartAccount"],
    mutationFn: async () => {
      const bundler = createBundlerClient({
        account: account.data!,
        client: client!,
        paymaster: createPaymasterClient({ transport: http(sponsorUrl) }),
        transport: http(bundlerUrl),
      });
      if (await account.data!.isDeployed()) return;
      const hash = await bundler.sendUserOperation({ callData: "0x" });
      setSubmitted(true);
      const { success } = await bundler.waitForUserOperationReceipt({ hash });
      if (!success) throw new Error("The account's creation reverted.");
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({ queryKey: ["smartAccountStanding", deployments.chainId, address] });
      setSubmitted(false);
    },
  });
  return {
    enabled,
    address,
    standing: standing.data,
    create: () => create.mutate(),
    phase: create.isPending ? (submitted ? "creating" : "signing") : undefined,
    error: create.error ?? account.error ?? standing.error ?? undefined,
  } as const;
}
