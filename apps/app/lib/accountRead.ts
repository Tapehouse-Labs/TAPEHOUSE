// SPDX-License-Identifier: MIT OR Apache-2.0
import { accounts, band, CROSS, type Deployments, engine, toBytes32 } from "@tapehouse/sdk";
import { type Address, type Client, type Hex, maxUint256, zeroAddress } from "viem";
import { getBlockNumber } from "viem/actions";

/** A Stock Token in a position: what it holds, what it lent and what the engine counts. */
export type Stock = {
  asset: string;
  token: Address;
  held: bigint;
  lent: bigint;
  sellable: bigint;
  claim: bigint;
  recalls: readonly bigint[];
  inBaskets: bigint;
  liquidationPrice: bigint;
};

export type Position = {
  id: Hex;
  name: string;
  equity: bigint;
  requirement: bigint;
  missing: number;
  regime: number;
  leverage: bigint;
  gross: bigint;
  debt: bigint;
  premium: bigint;
  usdg: bigint;
  weth: bigint;
  stocks: Stock[];
  open: bigint;
  closed: bigint;
};

export type Market = {
  asset: string;
  token: Address;
  low: bigint;
  state: number;
  lending: boolean;
  holding: Awaited<ReturnType<typeof accounts.holding>> | undefined;
  issuer: Awaited<ReturnType<typeof accounts.issuer>> | undefined;
};

export type AccountState = {
  blockNumber: bigint;
  markets: Market[];
  positions: Position[];
  limits: Awaited<ReturnType<typeof accounts.limits>>;
  session: Awaited<ReturnType<typeof band.session>>;
};

const PRICE_UNIT = 10n ** 8n;

/** The position `id` of the engine's assets, in words. */
export function positionName(id: Hex, markets: readonly Market[]): string {
  return id === CROSS ? "Cross" : (markets.find((m) => toBytes32(m.asset) === id)?.asset ?? id);
}

/**
 * Everything the account surface shows of `account`, which `owner` owns, read at one block: the market of each asset,
 * the issuer's pause and blocks on the accounts, the account and its owner among it, every position it holds anything
 * in, the accounts' limits and the session. A position is margined on the effective quantities the
 * engine counts, held, lent and in baskets, at each band's low edge.
 */
export async function readAccount(
  client: Client,
  deployments: Deployments,
  account: Address,
  owner: Address,
): Promise<AccountState> {
  const blockNumber = await getBlockNumber(client, { cacheTime: 0 });
  const at = { blockNumber };
  const assets = await accounts.stocks(client, deployments, at);
  const holders = [deployments.tapehouse.MarginAccounts as Address, account, owner];
  const [quotes, limits, session, holdings, issuers] = await Promise.all([
    Promise.all(assets.map(({ asset }) => band.quote(client, deployments, asset, at))),
    accounts.limits(client, deployments, at),
    band.session(client, deployments, at),
    Promise.all(
      assets.map(({ asset, token }) =>
        token === zeroAddress ? undefined : accounts.holding(client, deployments, asset, token, at),
      ),
    ),
    Promise.all(
      assets.map(({ token }) => (token === zeroAddress ? undefined : accounts.issuer(client, token, holders, at))),
    ),
  ]);
  const markets: Market[] = assets.map(({ asset, token }, i) => ({
    asset,
    token,
    low: quotes[i]!.low,
    state: quotes[i]!.state,
    lending: Object.hasOwn(deployments.stockLending, asset),
    holding: holdings[i],
    issuer: issuers[i],
  }));
  const traded = markets.filter((m) => m.token !== zeroAddress);
  const ids = [CROSS, ...traded.map((m) => toBytes32(m.asset))];
  const positions = await Promise.all(ids.map((id) => readPosition(client, deployments, account, id, markets, at)));
  return { blockNumber, markets, positions: positions.filter((p) => p !== undefined), limits, session };
}

async function readPosition(
  client: Client,
  deployments: Deployments,
  account: Address,
  id: Hex,
  markets: readonly Market[],
  at: { blockNumber: bigint },
): Promise<Position | undefined> {
  const usdgToken = deployments.tokens.USDG as Address;
  const wethToken = (deployments.tokens.WETH ?? zeroAddress) as Address;
  const traded = markets.filter((m) => m.token !== zeroAddress && (id === CROSS || toBytes32(m.asset) === id));
  const [health, leverage, { debt, premium }, usdg, weth, baskets, held] = await Promise.all([
    accounts.health(client, deployments, account, id, at),
    accounts.leverage(client, deployments, account, id, at),
    accounts.repayment(client, deployments, account, id, at),
    accounts.collateral(client, deployments, account, id, usdgToken, at),
    wethToken === zeroAddress ? 0n : accounts.collateral(client, deployments, account, id, wethToken, at),
    id === CROSS ? accounts.inBaskets(client, deployments, account, id, at) : ({} as Record<string, bigint>),
    Promise.all(
      traded.map(async (m) => {
        const [amount, lending] = await Promise.all([
          accounts.collateral(client, deployments, account, id, m.token, at),
          m.lending
            ? accounts.lending(client, deployments, account, id, m.token, at)
            : { lent: 0n, sellable: 0n, claim: 0n, recalls: [] as readonly bigint[] },
        ]);
        return { m, amount, lending };
      }),
    ),
  ]);
  const stocks = held
    .map(({ m, amount, lending }) => ({
      asset: m.asset,
      token: m.token,
      held: amount,
      ...lending,
      inBaskets: baskets[m.asset] ?? 0n,
      liquidationPrice: 0n,
    }))
    .filter((s) => s.held + s.lent + s.inBaskets + s.claim > 0n);
  if (health.equity === 0n && debt === 0n && premium === 0n && usdg === 0n && weth === 0n && stocks.length === 0)
    return undefined;
  const quantities = markets.map((m) => {
    const s = stocks.find((x) => x.asset === m.asset);
    return s ? s.held + s.lent + s.inBaskets : 0n;
  });
  const prices = markets.map((m) => m.low);
  const exposed = quantities.some((q) => q > 0n);
  const [open, closed, liquidationPrices] = await Promise.all([
    exposed ? engine.requirement(client, deployments, { quantities, prices }, { spansClosure: false }, at) : undefined,
    exposed ? engine.requirement(client, deployments, { quantities, prices }, { spansClosure: true }, at) : undefined,
    Promise.all(stocks.map((s) => accounts.liquidationPrice(client, deployments, account, id, s.asset, 0n, at))),
  ]);
  stocks.forEach((s, i) => (s.liquidationPrice = liquidationPrices[i]!));
  const gross =
    leverage === maxUint256 || health.equity <= 0n
      ? quantities.reduce((sum, q, i) => sum + (q * prices[i]!) / PRICE_UNIT, 0n)
      : (leverage * health.equity) / 10_000n;
  return {
    id,
    name: positionName(id, markets),
    equity: health.equity,
    requirement: health.requirement,
    missing: health.missing,
    regime: health.regime,
    leverage,
    gross,
    debt,
    premium,
    usdg,
    weth,
    stocks,
    open: open?.margin ?? 0n,
    closed: closed?.margin ?? 0n,
  };
}

/** The asset whose price governs a position: its own when isolated, else the one it holds the most value of. */
export function leadAsset(position: Position, markets: readonly Market[]): Stock | undefined {
  const value = (s: Stock) => ((s.held + s.lent + s.inBaskets) * (markets.find((m) => m.asset === s.asset)?.low ?? 0n));
  return position.stocks.reduce<Stock | undefined>((best, s) => (best && value(best) >= value(s) ? best : s), undefined);
}
