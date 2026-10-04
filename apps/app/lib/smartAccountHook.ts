// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { sponsorship } from "@tapehouse/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toSimpleSmartAccount } from "permissionless/accounts";
import { useState } from "react";
import { createPublicClient, createWalletClient, custom, type EIP1193Provider, http, isAddressEqual } from "viem";
import { createBundlerClient, createPaymasterClient } from "viem/account-abstraction";
import { useConnection, usePublicClient } from "wagmi";
import { useDeployments, useSmartAccountUrls } from "@/app/providers";
import { REFRESH_MS } from "./hooks";

/**
 * The visitor's smart account: a SimpleAccount v0.7 their wallet owns, whether it exists yet, and how many of its
 * operations the sponsor paymaster still pays for. `create` deploys it in a sponsored user operation that makes no call;
 * `bundler` builds the bundler client that sends its operations, asking the sponsor service for one stub an operation;
 * the service signs them over ERC-7677.
 * Its address is read through the app's node and the wallet's alike, as the wallet will sign permits to it. Off when
 * the app has no bundler or sponsor service, or the registry no sponsor paymaster.
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
      const account = await toSimpleSmartAccount({
        client: client!,
        owner: createWalletClient({ account: owner!, transport: custom(provider) }),
        factoryAddress: sponsorship.accountFactory(deployments),
        entryPoint: sponsorship.entryPoint(deployments),
      });
      const viaWallet = await sponsorship.accountAddress(
        createPublicClient({ transport: custom(provider) }),
        deployments,
        owner!,
      );
      if (!isAddressEqual(viaWallet, account.address))
        throw new Error("Your wallet and the app's node disagree on your smart account's address.");
      return account;
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
  const bundler = () => {
    const sponsor = createPaymasterClient({ transport: http(sponsorUrl) });
    const stubs = new Map<string, ReturnType<typeof sponsor.getPaymasterStubData>>();
    return createBundlerClient({
      account: account.data!,
      client: client!,
      paymaster: {
        getPaymasterStubData: (request) => {
          const key = `${request.sender}:${request.nonce}:${request.callData}`;
          if (!stubs.has(key)) stubs.set(key, sponsor.getPaymasterStubData(request));
          return stubs.get(key)!;
        },
        getPaymasterData: (request) => sponsor.getPaymasterData(request),
      },
      transport: http(bundlerUrl),
    });
  };
  const create = useMutation({
    mutationKey: ["createSmartAccount"],
    mutationFn: async () => {
      const client = bundler();
      if (await account.data!.isDeployed()) return;
      const hash = await client.sendUserOperation({ callData: "0x" });
      setSubmitted(true);
      const { success } = await client.waitForUserOperationReceipt({ hash });
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
    owner,
    standing: standing.data,
    bundler: account.data ? bundler : undefined,
    create: () => create.mutate(),
    phase: create.isPending ? (submitted ? "creating" : "signing") : undefined,
    error: create.error ?? account.error ?? standing.error ?? undefined,
  } as const;
}
