// SPDX-License-Identifier: MIT OR Apache-2.0
import { accounts, backstop, type Deployments, gapCover, stockLending, supply } from "@tapehouse/sdk";
import { type Address, type Client, erc20Abi, zeroAddress } from "viem";
import { getBlock, getBlockNumber, readContract } from "viem/actions";

export type Holding = Awaited<ReturnType<typeof supply.lender>>;

export type EarnState = {
  blockNumber: bigint;
  timestamp: bigint;
  supply?: { market: Awaited<ReturnType<typeof supply.market>>; lender?: Holding };
  backstop?: {
    state: Awaited<ReturnType<typeof backstop.state>>;
    depositor?: Awaited<ReturnType<typeof backstop.depositor>>;
    tokens: { asset: string; token: Address }[];
  };
  lending: { asset: string; token: Address; terms: Awaited<ReturnType<typeof stockLending.terms>> }[];
  cover?: { writers: Awaited<ReturnType<typeof gapCover.writers>>; writer?: Holding };
};

/**
 * Everything the Earn surface shows, read at one block: the supply vault, the gap backstop with its exposure and
 * premium, each Stock Token's lending vault and the gap cover's writers, each only where the registry deploys it, and
 * what `account` holds in each once a wallet is connected.
 */
export async function readEarn(client: Client, deployments: Deployments, account: Address | undefined): Promise<EarnState> {
  const block = await getBlock(client, { blockTag: "latest" });
  const at = { blockNumber: block.number };
  const { tapehouse } = deployments;
  const stocks = tapehouse.MarginAccounts ? await accounts.stocks(client, deployments, at) : [];
  const tokens = stocks.filter((s) => s.token !== zeroAddress);
  const [supplyState, backstopState, lending, cover] = await Promise.all([
    tapehouse.SupplyVault
      ? Promise.all([
          supply.market(client, deployments, at),
          account ? supply.lender(client, deployments, account, at) : undefined,
        ]).then(([market, lender]) => ({ market, lender }))
      : undefined,
    tapehouse.GapBackstop
      ? Promise.all([
          backstop.state(client, deployments, at),
          account
            ? backstop.depositor(
                client,
                deployments,
                account,
                tokens.map((s) => s.token),
                at,
              )
            : undefined,
        ]).then(([state, depositor]) => ({ state, depositor, tokens }))
      : undefined,
    Promise.all(
      Object.keys(deployments.stockLending).map(async (asset) => ({
        asset,
        token: tokens.find((s) => s.asset === asset)?.token ?? zeroAddress,
        terms: await stockLending.terms(client, deployments, asset, at),
      })),
    ),
    tapehouse.GapCover
      ? Promise.all([
          gapCover.writers(client, deployments, at),
          account ? gapCover.writer(client, deployments, account, at) : undefined,
        ]).then(([writers, writer]) => ({ writers, writer }))
      : undefined,
  ]);
  return {
    blockNumber: block.number,
    timestamp: block.timestamp,
    supply: supplyState,
    backstop: backstopState,
    lending,
    cover,
  };
}

/**
 * What the smart account `account` sends its owner when it claims its gains of `token`, read at one block: the gains
 * the claim pays, with any of `token` it already holds, so nothing it claimed is left behind.
 */
export async function claimable(client: Client, deployments: Deployments, account: Address, token: Address): Promise<bigint> {
  const blockNumber = await getBlockNumber(client, { cacheTime: 0 });
  const [{ gains }, held] = await Promise.all([
    backstop.depositor(client, deployments, account, [token], { blockNumber }),
    readContract(client, { address: token, abi: erc20Abi, functionName: "balanceOf", args: [account], blockNumber }),
  ]);
  return (gains[token] ?? 0n) + held;
}
