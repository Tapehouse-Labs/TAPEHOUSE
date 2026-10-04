// SPDX-License-Identifier: MIT OR Apache-2.0
import { accounts, backstop, type Deployments, gapCover, sponsorship, supply } from "@tapehouse/sdk";
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
  | { kind: "claimGains"; token: Address; amount: bigint };

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
 * user operation. Loans and withdrawals go to the wallet, and so do the Stock Tokens the backstop's gains pay out. A
 * repayment ends by taking back any allowance it left, as it pays only what is owed.
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
    default:
      return { wrap: 0n, pulls: [], calls: [asCall(direct(deployments, intent, account, owner))] };
  }
}

/** The one transaction `intent` is when it takes nothing from the wallet, sent for `account`. */
export function direct(
  deployments: Deployments,
  intent: Exclude<Intent, { kind: "deposit" | "repay" | "vaultDeposit" }>,
  account: Address,
  receiver: Address,
) {
  switch (intent.kind) {
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
    default:
      return accounts[intent.kind](deployments, {
        position: intent.position,
        token: intent.token,
        amount: intent.amount,
        account,
      });
  }
}
