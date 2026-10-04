// SPDX-License-Identifier: MIT OR Apache-2.0
import { accounts, backstop, baskets, CROSS, type Deployments, gapCover, shorts, sponsorship, supply } from "@tapehouse/sdk";
import { type Abi, type Address, erc20Abi, type Hex } from "viem";

/** A vault that takes USDG for shares: the supply vault, the gap backstop and the gap cover's writers. */
export type Vault = "supply" | "backstop" | "cover";

/** What the visitor asks the account to do. */
export type Intent =
  | { kind: "deposit"; position: Hex; token: Address | "ETH"; amount: bigint }
  | { kind: "withdraw"; position: Hex; token: Address; amount: bigint }
  | { kind: "borrow"; position: Hex; amount: bigint }
  | { kind: "repay"; position: Hex; amount: bigint }
  | { kind: "lend" | "unlend" | "recall"; position: Hex; token: Address; amount: bigint }
  | { kind: "settle"; position: Hex; token: Address }
  | { kind: "vaultDeposit"; vault: Vault; amount: bigint }
  | { kind: "vaultWithdraw"; vault: Vault; assets: bigint; shares?: bigint }
  | { kind: "cooldown" }
  | { kind: "claimPremium" }
  | { kind: "claimGains"; token: Address; amount: bigint }
  | { kind: "shortDeposit"; asset: string; amount: bigint }
  | { kind: "shortWithdraw"; asset: string; amount: bigint }
  | { kind: "sell"; asset: string; amount: bigint; minProceeds: bigint }
  | { kind: "buyBack"; asset: string; amount: bigint; maxCost: bigint }
  | { kind: "buyCover"; asset: string; notional: bigint; deductibleBps: bigint; limitBps: bigint; maxPremium: bigint }
  | { kind: "release"; id: bigint }
  | { kind: "claimCover" }
  | { kind: "basketMint"; basket: string; shares: bigint; tokens: readonly Address[]; maxAssets: readonly bigint[] }
  | { kind: "basketRedeem"; basket: string; shares: bigint }
  | { kind: "unwrap"; basket: string; shares: bigint };

/** The intents that take tokens from the wallet first, or move them through the account; the rest are one call. */
export type DirectIntent = Exclude<
  Intent,
  { kind: "deposit" | "repay" | "vaultDeposit" | "shortDeposit" | "buyCover" | "basketMint" | "basketRedeem" }
>;

/** The vault's address in the registry. */
export function vaultAddress(deployments: Deployments, vault: Vault): Address {
  const name = { supply: "SupplyVault", backstop: "GapBackstop", cover: "GapCover" }[vault];
  const address = deployments.tapehouse[name];
  if (!address) throw new Error(`The registry has no .tapehouse.${name}.`);
  return address;
}

const vaults = { supply, backstop, cover: gapCover } as const;

/** A call of a user operation. */
export type Call = { to: Address; abi: Abi; functionName: string; args: readonly unknown[] };

/** A transaction of the SDK as a call of a user operation. */
export function asCall(transaction: { address: Address; abi: Abi; functionName: string; args?: readonly unknown[] }): Call {
  return sponsorship.asCall({ args: [], ...transaction } as never) as Call;
}

/** A token the account must take from its owner first, with a permit the owner signs for free. */
export type Pull = { token: Address; amount: bigint };

/**
 * The plan that carries out `intent` for the smart account `account` its wallet `owner` owns: the ether the wallet
 * wraps into WETH for the account first, the tokens the account takes from the wallet by permit, and the calls of one
 * user operation. Loans and withdrawals go to the wallet, and so do the Stock Tokens the backstop's gains pay out, the
 * cover's payouts and a basket's redeemed tokens. A repayment, a cover's premium and a mint end by taking back any
 * allowance they left, as each takes only what it costs.
 */
export function smartAccountPlan(
  deployments: Deployments,
  intent: Intent,
  account: Address,
  owner: Address,
): { wrap: bigint; pulls: Pull[]; calls: Call[] } {
  const spender = deployments.tapehouse.MarginAccounts as Address;
  const usdg = deployments.tokens.USDG as Address;
  const approve = (token: Address, amount: bigint, to: Address = spender) =>
    asCall({ address: token, abi: erc20Abi, functionName: "approve", args: [to, amount] });
  switch (intent.kind) {
    case "deposit": {
      const eth = intent.token === "ETH";
      const token = eth ? (deployments.tokens.WETH as Address) : (intent.token as Address);
      const { position, amount } = intent;
      return {
        wrap: eth ? amount : 0n,
        pulls: eth ? [] : [{ token, amount }],
        calls: [
          approve(token, amount),
          asCall(accounts.deposit(deployments, { position, token, amount, account })),
        ],
      };
    }
    case "repay":
      return {
        wrap: 0n,
        pulls: [{ token: usdg, amount: intent.amount }],
        calls: [
          approve(usdg, intent.amount),
          asCall(accounts.repay(deployments, { position: intent.position, assets: intent.amount, account })),
          approve(usdg, 0n),
        ],
      };
    case "vaultDeposit":
      return {
        wrap: 0n,
        pulls: [{ token: usdg, amount: intent.amount }],
        calls: [
          approve(usdg, intent.amount, vaultAddress(deployments, intent.vault)),
          asCall(vaults[intent.vault].deposit(deployments, { assets: intent.amount, receiver: account })),
        ],
      };
    case "claimGains":
      return {
        wrap: 0n,
        pulls: [],
        calls: [
          asCall(backstop.claimGains(deployments, intent.token)),
          asCall({ address: intent.token, abi: erc20Abi, functionName: "transfer", args: [owner, intent.amount] }),
        ],
      };
    case "shortDeposit": {
      const to = deployments.tapehouse.ShortPositions as Address;
      return {
        wrap: 0n,
        pulls: [{ token: usdg, amount: intent.amount }],
        calls: [
          approve(usdg, intent.amount, to),
          asCall(shorts.deposit(deployments, { asset: intent.asset, amount: intent.amount, account })),
        ],
      };
    }
    case "buyCover": {
      const to = deployments.tapehouse.GapCover as Address;
      const { asset, notional, deductibleBps, limitBps, maxPremium } = intent;
      return {
        wrap: 0n,
        pulls: [{ token: usdg, amount: intent.maxPremium }],
        calls: [
          approve(usdg, intent.maxPremium, to),
          asCall(gapCover.buy(deployments, { asset, notional, deductibleBps, limitBps, maxPremium, holder: account })),
          approve(usdg, 0n, to),
        ],
      };
    }
    case "basketMint": {
      const basket = deployments.baskets[intent.basket] as Address;
      return {
        wrap: 0n,
        pulls: intent.tokens.map((token, i) => ({ token, amount: intent.maxAssets[i] ?? 0n })),
        calls: [
          ...intent.tokens.map((token, i) => approve(token, intent.maxAssets[i] ?? 0n, basket)),
          asCall(baskets.mint(deployments, { basket: intent.basket, shares: intent.shares, receiver: account, maxAssets: intent.maxAssets })),
          ...intent.tokens.map((token) => approve(token, 0n, basket)),
          approve(basket, intent.shares),
          asCall(accounts.deposit(deployments, { position: CROSS, token: basket, amount: intent.shares, account })),
        ],
      };
    }
    case "basketRedeem": {
      const basket = deployments.baskets[intent.basket] as Address;
      return {
        wrap: 0n,
        pulls: [],
        calls: [
          asCall(accounts.withdraw(deployments, { position: CROSS, token: basket, amount: intent.shares, account, receiver: account })),
          asCall(baskets.redeem(deployments, { basket: intent.basket, shares: intent.shares, receiver: owner, owner: account })),
        ],
      };
    }
    default:
      return { wrap: 0n, pulls: [], calls: [asCall(direct(deployments, intent, account, owner))] };
  }
}

/** What `pulls` still take from the wallet once the account spends what it `held` of each token first. */
export function shortfall(pulls: readonly Pull[], held: readonly bigint[]): Pull[] {
  return pulls.flatMap((pull, i) => {
    const amount = pull.amount - (held[i] ?? 0n);
    return amount > 0n ? [{ token: pull.token, amount }] : [];
  });
}

/** The one transaction `intent` is when it takes nothing from the wallet, sent for `account`. */
export function direct(
  deployments: Deployments,
  intent: DirectIntent,
  account: Address,
  receiver: Address,
) {
  switch (intent.kind) {
    case "sell":
      return shorts.sell(deployments, { asset: intent.asset, amount: intent.amount, minProceeds: intent.minProceeds, account });
    case "buyBack":
      return shorts.cover(deployments, { asset: intent.asset, amount: intent.amount, maxCost: intent.maxCost, account });
    case "shortWithdraw":
      return shorts.withdraw(deployments, { asset: intent.asset, amount: intent.amount, account, receiver });
    case "release":
      return gapCover.release(deployments, intent.id);
    case "claimCover":
      return gapCover.claim(deployments, receiver);
    case "unwrap":
      return accounts.unwrap(deployments, { account, basket: intent.basket, shares: intent.shares });
    case "vaultWithdraw": {
      const { vault, assets, shares } = intent;
      return shares === undefined
        ? vaults[vault].withdraw(deployments, { assets, receiver, owner: account })
        : vaults[vault].redeem(deployments, { shares, receiver, owner: account });
    }
    case "cooldown":
      return backstop.startCooldown(deployments);
    case "claimPremium":
      return backstop.claim(deployments);
    case "claimGains":
      return backstop.claimGains(deployments, intent.token);
    case "withdraw":
      return accounts.withdraw(deployments, { ...intent, account, receiver });
    case "borrow":
      return accounts.borrow(deployments, { position: intent.position, assets: intent.amount, account, receiver });
    case "settle":
      return accounts.settle(deployments, { account, position: intent.position, token: intent.token });
    case "lend":
    case "unlend":
    case "recall":
      return accounts[intent.kind](deployments, {
        position: intent.position,
        token: intent.token,
        amount: intent.amount,
        account,
      });
  }
}
