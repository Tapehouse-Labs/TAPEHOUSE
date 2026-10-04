// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import {
  accounts,
  backstop,
  baskets,
  CROSS,
  decodeRevertData,
  type Deployments,
  gapCover,
  shorts,
  sponsorship,
  supply,
} from "@tapehouse/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  type Account,
  type Address,
  type Chain,
  encodeFunctionData,
  erc20Abi,
  type Hex,
  type PublicClient,
  type Transport,
  type WalletClient,
  zeroAddress,
} from "viem";
import { useConnection, usePublicClient, useWalletClient } from "wagmi";
import { useDeployments } from "@/app/providers";
import { readAccount } from "./accountRead";
import { describeRevert } from "./errors";
import { REFRESH_MS, useSettled } from "./hooks";
import { asCall, type Call, direct, type Intent, shortfall, smartAccountPlan, vaultAddress } from "./intent";
import { useSmartAccount } from "./smartAccountHook";

/**
 * The call gas a user operation carries beyond its estimate, on top of a fifth or a quarter more: another operation in
 * between can turn a write the estimate found warm into a cold one.
 */
const CALL_GAS_MARGIN = 20_000n;

/** How long a permit the wallet signs stays good, in seconds. */
const PERMIT_SECONDS = 3_600n;

/**
 * The margin account the visitor acts on: their smart account when the app sponsors one, else their wallet, with the
 * wallet that owns it.
 */
export function useMarginAccount() {
  const smart = useSmartAccount();
  const { address: owner } = useConnection();
  if (smart.enabled) return { account: smart.address, owner, smart } as const;
  return { account: owner, owner, smart } as const;
}

/** Everything the account surface shows of `account`, which `owner` owns, read at one block and every few seconds. */
export function useAccountState(account: Address | undefined, owner: Address | undefined) {
  const deployments = useDeployments();
  const client = usePublicClient();
  return useQuery({
    queryKey: ["account", deployments.chainId, account, owner],
    queryFn: () => readAccount(client!, deployments, account!, owner!),
    enabled: client !== undefined && account !== undefined && owner !== undefined,
    refetchInterval: REFRESH_MS,
  });
}

/** What the wallet holds of ether, USDG, WETH and each Stock Token, for the deposit and repay amounts. */
export function useWalletBalances(owner: Address | undefined, stockTokens: readonly Address[]) {
  const deployments = useDeployments();
  const client = usePublicClient();
  return useQuery({
    queryKey: ["walletBalances", deployments.chainId, owner, stockTokens],
    queryFn: async () => {
      const tokens = [deployments.tokens.USDG, deployments.tokens.WETH, ...stockTokens].filter(
        (t): t is Address => t !== undefined && t !== zeroAddress,
      );
      const [eth, ...amounts] = await Promise.all([
        client!.getBalance({ address: owner! }),
        ...tokens.map((address) =>
          client!.readContract({ address, abi: erc20Abi, functionName: "balanceOf", args: [owner!] }),
        ),
      ]);
      return { eth, tokens: Object.fromEntries(tokens.map((t, i) => [t, amounts[i]!])) as Record<Address, bigint> };
    },
    enabled: client !== undefined && owner !== undefined,
    refetchInterval: REFRESH_MS,
  });
}

/**
 * The liquidation price of `asset` for `account`'s `position` once it owes `borrowing` more USDG, in USD with 8
 * decimals, as the accounts compute it.
 */
export function useLiquidationPreview(
  account: Address | undefined,
  position: Hex | undefined,
  asset: string | undefined,
  borrowing: bigint,
) {
  const deployments = useDeployments();
  const client = usePublicClient();
  const settled = useSettled(borrowing);
  return useQuery({
    queryKey: ["liquidationPreview", deployments.chainId, account, position, asset, settled.toString()],
    queryFn: () => accounts.liquidationPrice(client!, deployments, account!, position!, asset!, settled),
    enabled: client !== undefined && account !== undefined && position !== undefined && asset !== undefined,
    placeholderData: (previous) => previous,
  });
}

export type SendPhase = "wrapping" | "permit" | "checking" | "signing" | "sending";

/** What the primary button says while a send is in each phase. */
export const PHASE: Record<SendPhase, string> = {
  wrapping: "Wrapping your ether…",
  permit: "Sign the permit in your wallet…",
  checking: "Checking it against the chain…",
  signing: "Waiting for your wallet…",
  sending: "Sending…",
};

type Sender = {
  deployments: Deployments;
  client: PublicClient;
  wallet: WalletClient<Transport, Chain, Account>;
  account: Address;
  owner: Address;
  bundler: ReturnType<typeof useSmartAccount>["bundler"];
  phase: (phase: SendPhase) => void;
};

/**
 * Carries out an intent for the margin account. From a smart account: the wallet wraps any ether into WETH for it,
 * signs a free permit for each token it takes from the wallet, and signs one sponsored user operation, which the
 * bundler estimates, and so simulates, before the wallet is asked. A token the wallet already allows the account to take
 * needs no permit, so a permit someone else sent first does not fail the operation. From the wallet: each transaction is simulated,
 * then sent. Loans and withdrawals go to the wallet.
 */
export function useSend() {
  const deployments = useDeployments();
  const client = usePublicClient();
  const { data: wallet } = useWalletClient();
  const { account, owner, smart } = useMarginAccount();
  const queryClient = useQueryClient();
  const [phase, setPhase] = useState<SendPhase>();
  const mutation = useMutation({
    mutationKey: ["send"],
    mutationFn: async (request: Intent | (() => Promise<Intent>)) => {
      const intent = typeof request === "function" ? await request() : request;
      if (!client || !wallet || !account || !owner) throw new Error("Connect your wallet first.");
      const sender = { deployments, client, wallet, account, owner, bundler: smart.bundler, phase: setPhase };
      if (smart.enabled) await sendFromSmartAccount(sender, intent);
      else await sendFromWallet(sender, intent);
    },
    onSettled: async () => {
      setPhase(undefined);
      await queryClient.invalidateQueries();
    },
  });
  return {
    send: (intent: Intent | (() => Promise<Intent>), options?: { onSuccess?: () => void }) => mutation.mutate(intent, options),
    phase,
    pending: mutation.isPending,
    error: mutation.error ?? undefined,
    reset: mutation.reset,
    succeeded: mutation.isSuccess,
  } as const;
}

function permitDeadline(): bigint {
  return BigInt(Math.floor(Date.now() / 1000)) + PERMIT_SECONDS;
}

async function sendFromSmartAccount(s: Sender, intent: Intent) {
  const plan = smartAccountPlan(s.deployments, intent, s.account, s.owner);
  if (plan.wrap > 0n) {
    s.phase("wrapping");
    const { request } = await s.client.simulateContract({
      ...accounts.wrap(s.deployments, { account: s.account, value: plan.wrap }),
      account: s.owner,
    });
    await confirm(s, await s.wallet.writeContract(request));
  }
  const deadline = permitDeadline();
  const held = await Promise.all(
    plan.pulls.map(({ token }) =>
      s.client.readContract({ address: token, abi: erc20Abi, functionName: "balanceOf", args: [s.account] }),
    ),
  );
  const pulls = [];
  for (const { token, amount } of shortfall(plan.pulls, held)) {
    const allowed = await s.client.readContract({
      address: token,
      abi: erc20Abi,
      functionName: "allowance",
      args: [s.owner, s.account],
    });
    if (allowed >= amount) {
      pulls.push(asCall({ address: token, abi: erc20Abi, functionName: "transferFrom", args: [s.owner, s.account, amount] }));
      continue;
    }
    s.phase("permit");
    const permit = await sponsorship.permit(s.client, { token, owner: s.owner, account: s.account, value: amount, deadline });
    const signature = await s.wallet.signTypedData({ account: s.owner, ...permit });
    pulls.push(...sponsorship.pull({ token, owner: s.owner, account: s.account, value: amount, deadline, signature }));
  }
  if (plan.pulls.length === 0 && plan.calls.length === 1) {
    s.phase("checking");
    const [call] = plan.calls;
    await s.client.simulateContract({
      address: call!.to,
      abi: call!.abi,
      functionName: call!.functionName,
      args: call!.args,
      account: s.account,
    } as never);
  }
  s.phase("signing");
  const bundler = s.bundler!();
  const calls = [...pulls, ...plan.calls] as never;
  const hash = await sendWithFreshGas(s, bundler, calls);
  s.phase("sending");
  const { success, reason } = await bundler.waitForUserOperationReceipt({ hash });
  if (!success) {
    const revert = reason?.startsWith("0x") ? decodeRevertData(reason as Hex) : undefined;
    throw new Error(revert ? describeRevert(revert) : "The operation reverted on chain.");
  }
}

/**
 * Sends `calls` from the smart account. Once the account exists, its call gas is half again the calls' own cold L2 gas:
 * a bundler that searches inside one simulation runs a Stylus program, such as the margin engine, warm after its first
 * pass, and underestimates it; and the cold figure is net of storage refunds, which hide up to a quarter of what a call
 * that clears storage, such as a short's full buy-back, must hold at its peak. Before, it is a quarter over the bundler's estimate. The sponsor service bounds the
 * pre-verification gas by the L1 gas it estimates as it signs, which falls as the L2 base fee rises; an estimate that
 * went stale between the two is estimated again, once, before the wallet is asked.
 */
async function sendWithFreshGas(s: Sender, bundler: ReturnType<NonNullable<Sender["bundler"]>>, calls: never) {
  for (let attempt = 0; ; attempt++) {
    try {
      let callGasLimit: bigint;
      if (await bundler.account.isDeployed()) {
        const callData = await bundler.account.encodeCalls(
          (calls as readonly Call[]).map((c) => ({ to: c.to, data: encodeFunctionData(c as never), value: 0n })),
        );
        callGasLimit = ((await sponsorship.callGas(s.client, s.deployments, s.account, callData)) * 3n) / 2n;
      } else {
        callGasLimit = ((await bundler.estimateUserOperationGas({ calls })).callGasLimit * 5n) / 4n;
      }
      return await bundler.sendUserOperation({ calls, callGasLimit: callGasLimit + CALL_GAS_MARGIN });
    } catch (error) {
      if (attempt > 0 || !String(error).includes("pre-verification gas is above")) throw error;
    }
  }
}

async function sendFromWallet(s: Sender, intent: Intent) {
  const transactions = await walletTransactions(s.deployments, intent, s.owner, async (token, amount, spender) => {
    s.phase("permit");
    const deadline = permitDeadline();
    const permit = await sponsorship.permit(s.client, { token, owner: s.owner, account: spender, value: amount, deadline });
    return { deadline, signature: await s.wallet.signTypedData({ account: s.owner, ...permit }) };
  });
  for (const transaction of transactions) {
    s.phase("checking");
    const { request } = await s.client.simulateContract({ ...transaction, account: s.owner } as never);
    s.phase("signing");
    await confirm(s, await s.wallet.writeContract(request as never));
  }
}

async function confirm(s: Sender, hash: Hex) {
  s.phase("sending");
  const receipt = await s.client.waitForTransactionReceipt({ hash });
  if (receipt.status !== "success") throw new Error("The transaction reverted.");
}

/**
 * The transactions the wallet sends for `intent` when it is the margin account: a deposit with a permit `sign`
 * returns for `spender`, ether wrapped first, a supply with a USDG permit, a repayment, a short's margin or a deposit
 * with the backstop or the cover's writers after its approval, a cover's premium approved exactly and taken back, a
 * basket minted from approved tokens into the cross position or withdrawn and redeemed, and anything else as it is.
 */
async function walletTransactions(
  deployments: Deployments,
  intent: Intent,
  owner: Address,
  sign: (token: Address, amount: bigint, spender: Address) => Promise<{ deadline: bigint; signature: Hex }>,
) {
  const spender = deployments.tapehouse.MarginAccounts as Address;
  const usdg = deployments.tokens.USDG as Address;
  if (intent.kind === "vaultDeposit") {
    const vault = vaultAddress(deployments, intent.vault);
    const assets = intent.amount;
    if (intent.vault === "supply") {
      const { deadline, signature } = await sign(usdg, assets, vault);
      return [supply.depositWithPermit(deployments, { assets, receiver: owner, deadline, signature })];
    }
    const deposit = (intent.vault === "backstop" ? backstop : gapCover).deposit(deployments, { assets, receiver: owner });
    return [{ address: usdg, abi: erc20Abi, functionName: "approve", args: [vault, assets] } as const, deposit];
  }
  if (intent.kind === "deposit") {
    const eth = intent.token === "ETH";
    const token = eth ? (deployments.tokens.WETH as Address) : (intent.token as Address);
    const { deadline, signature } = await sign(token, intent.amount, spender);
    const deposit = accounts.depositWithPermit(deployments, {
      position: intent.position,
      token,
      amount: intent.amount,
      account: owner,
      deadline,
      signature,
    });
    return eth ? [accounts.wrap(deployments, { account: owner, value: intent.amount }), deposit] : [deposit];
  }
  const approve = (token: Address, to: Address, amount: bigint) =>
    ({ address: token, abi: erc20Abi, functionName: "approve", args: [to, amount] }) as const;
  switch (intent.kind) {
    case "repay":
      return [
        approve(usdg, spender, intent.amount),
        accounts.repay(deployments, { position: intent.position, assets: intent.amount, account: owner }),
      ];
    case "shortDeposit":
      return [
        approve(usdg, deployments.tapehouse.ShortPositions as Address, intent.amount),
        shorts.deposit(deployments, { asset: intent.asset, amount: intent.amount, account: owner }),
      ];
    case "buyCover": {
      const cover = deployments.tapehouse.GapCover as Address;
      const { asset, notional, deductibleBps, limitBps, maxPremium } = intent;
      return [
        approve(usdg, cover, intent.maxPremium),
        gapCover.buy(deployments, { asset, notional, deductibleBps, limitBps, maxPremium, holder: owner }),
        approve(usdg, cover, 0n),
      ];
    }
    case "basketMint": {
      const basket = deployments.baskets[intent.basket] as Address;
      return [
        ...intent.tokens.map((token, i) => approve(token, basket, intent.maxAssets[i] ?? 0n)),
        baskets.mint(deployments, { basket: intent.basket, shares: intent.shares, receiver: owner, maxAssets: intent.maxAssets }),
        ...intent.tokens.map((token) => approve(token, basket, 0n)),
        approve(basket, spender, intent.shares),
        accounts.deposit(deployments, { position: CROSS, token: basket, amount: intent.shares, account: owner }),
      ];
    }
    case "basketRedeem": {
      const basket = deployments.baskets[intent.basket] as Address;
      return [
        accounts.withdraw(deployments, { position: CROSS, token: basket, amount: intent.shares, account: owner, receiver: owner }),
        baskets.redeem(deployments, { basket: intent.basket, shares: intent.shares, receiver: owner, owner }),
      ];
    }
    default:
      return [direct(deployments, intent, owner, owner)];
  }
}

