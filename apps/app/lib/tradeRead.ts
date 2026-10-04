// SPDX-License-Identifier: MIT OR Apache-2.0
import {
  accounts,
  band,
  baskets,
  CROSS,
  type Deployments,
  gapCover,
  shortPositionsAbi,
  shorts,
  stockLending,
  toBytes32,
} from "@tapehouse/sdk";
import { type Abi, type Address, type Client, encodeFunctionData, type Hex, hexToString, zeroAddress } from "viem";
import { getBlock, simulateCalls } from "viem/actions";

export type Market = {
  asset: string;
  token: Address;
  quote: Awaited<ReturnType<typeof band.quote>>;
  fee: number;
  restricted: boolean;
  epoch: bigint;
  borrowRate: bigint | undefined;
  borrowable: bigint | undefined;
  short?: {
    position: Awaited<ReturnType<typeof shorts.position>>;
    health: Awaited<ReturnType<typeof shorts.health>>;
  };
};

export type BasketState = {
  key: string;
  address: Address;
  assets: string[];
  tokens: readonly Address[];
  target: readonly bigint[];
  pending: { units: readonly bigint[]; effectiveAt: bigint } | undefined;
  held?: bigint;
};

export type TradeState = {
  blockNumber: bigint;
  timestamp: bigint;
  session: Awaited<ReturnType<typeof band.session>>;
  markets: Market[];
  cover?: {
    sales: Awaited<ReturnType<typeof gapCover.sales>>;
    capacity: bigint;
    payouts: bigint | undefined;
    /** The account's covers, undefined where the node would not return the cover's logs. */
    holdings:
      | (Awaited<ReturnType<typeof gapCover.holdings>>[number] & {
          asset: string;
          series: Awaited<ReturnType<typeof gapCover.series>>;
        })[]
      | undefined;
  };
  baskets: BasketState[];
  cross?: { health: Awaited<ReturnType<typeof accounts.health>>; inBaskets: Record<string, bigint> };
};

/**
 * Everything the Trade surface shows, read at one block: each Stock Token's band, whether it may be shorted and at
 * what fee tier, the restriction after a 10% fall, its book and lending rate, and `account`'s short of it; the gap
 * cover's sales, capacity and `account`'s covers with their series and payouts; and each basket's components, target
 * and pending target with what `account`'s cross position holds of it and through it.
 */
export async function readTrade(client: Client, deployments: Deployments, account: Address | undefined): Promise<TradeState> {
  const block = await getBlock(client, { blockTag: "latest" });
  const at = { blockNumber: block.number };
  const shortable = deployments.tapehouse.ShortPositions !== undefined;
  const stocks = (await accounts.stocks(client, deployments, at)).filter((s) => s.token !== zeroAddress);
  const [session, markets, cover, basketStates, cross] = await Promise.all([
    band.session(client, deployments, at),
    Promise.all(
      stocks.map(async ({ asset, token }): Promise<Market> => {
        const [quote, fee, restriction, epoch, terms] = await Promise.all([
          band.quote(client, deployments, asset, at),
          shortable ? shorts.fee(client, deployments, asset, at).catch(() => 0) : 0,
          shortable ? shorts.restriction(client, deployments, asset, at) : { restricted: false },
          shortable ? shorts.epoch(client, deployments, asset, at) : 0n,
          deployments.stockLending[asset] ? stockLending.terms(client, deployments, asset, at) : undefined,
        ]);
        const market: Market = {
          asset,
          token,
          quote,
          fee,
          restricted: restriction.restricted,
          epoch,
          borrowRate: terms?.borrowRate,
          borrowable: terms?.borrowable,
        };
        if (account && fee > 0) {
          const [position, health] = await Promise.all([
            shorts.position(client, deployments, account, asset, at),
            shorts.health(client, deployments, account, asset, at),
          ]);
          market.short = { position, health };
        }
        return market;
      }),
    ),
    deployments.tapehouse.GapCover ? readCover(client, deployments, account, at.blockNumber) : undefined,
    Promise.all(
      Object.entries(deployments.baskets).map(async ([key, address]): Promise<BasketState> => {
        const [{ assets, tokens }, target, pending, held] = await Promise.all([
          baskets.components(client, deployments, key, at),
          baskets.target(client, deployments, key, at),
          baskets.pendingTarget(client, deployments, key, at),
          account ? accounts.collateral(client, deployments, account, CROSS, address, at) : undefined,
        ]);
        return {
          key,
          address,
          assets,
          tokens,
          target,
          pending: pending.effectiveAt === 0n ? undefined : pending,
          held,
        };
      }),
    ),
    account && Object.keys(deployments.baskets).length > 0
      ? Promise.all([
          accounts.health(client, deployments, account, CROSS, at),
          accounts.inBaskets(client, deployments, account, CROSS, at),
        ]).then(([health, inBaskets]) => ({ health, inBaskets }))
      : undefined,
  ]);
  return {
    blockNumber: block.number,
    timestamp: block.timestamp,
    session,
    markets,
    cover,
    baskets: basketStates,
    cross,
  };
}

async function readCover(client: Client, deployments: Deployments, account: Address | undefined, blockNumber: bigint) {
  const at = { blockNumber };
  const [sales, capacity, payouts, held] = await Promise.all([
    gapCover.sales(client, deployments, at),
    gapCover.capacity(client, deployments, at),
    account ? gapCover.payouts(client, deployments, account, at) : undefined,
    account ? gapCover.holdings(client, deployments, account, at).catch(() => undefined) : [],
  ]);
  const holdings =
    held &&
    (await Promise.all(
      held.map(async (h) => {
        const asset = symbolName(h.symbol);
        return { ...h, asset, series: await gapCover.series(client, deployments, asset, h.closesMs, at) };
      }),
    ));
  return { sales, capacity, payouts, holdings };
}

function symbolName(symbol: Hex): string {
  return hexToString(symbol, { size: 32 });
}

/**
 * The short's health after `call` runs from `account`, simulated with `eth_simulateV1` at the latest block: the call,
 * then `ShortPositions.health`; the call's error where it reverts, and undefined where the node cannot simulate calls.
 */
export async function healthAfter(
  client: Client,
  deployments: Deployments,
  account: Address,
  asset: string,
  call: { address: Address; abi: Abi; functionName: string; args: readonly unknown[] },
): Promise<{ equity: bigint; requirement: bigint } | { error: Error } | undefined> {
  try {
    const { results } = await simulateCalls(client, {
      account,
      calls: [
        { to: call.address, data: encodeFunctionData(call as never) },
        {
          to: deployments.tapehouse.ShortPositions as Address,
          abi: shortPositionsAbi,
          functionName: "health",
          args: [account, toBytes32(asset)],
        },
      ],
    });
    const [ran, after] = results;
    if (ran?.status === "failure") return { error: ran.error };
    if (after?.status !== "success") return undefined;
    const [equity, requirement] = after.result;
    return equity === undefined || requirement === undefined ? undefined : { equity, requirement };
  } catch {
    return undefined;
  }
}
