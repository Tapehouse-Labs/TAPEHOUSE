// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import AxeBuilder from "@axe-core/playwright";
import { expect, type Page, test } from "@playwright/test";
import { accounts, CROSS, gapCover, parseDeployments, shorts, sponsorship } from "@tapehouse/sdk";
import { createPublicClient, createWalletClient, erc20Abi, type Hex, http, parseAbi } from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { usdg as usdgWords } from "../lib/amounts";
import { installWallet } from "./wallet";

const registry = process.env.TAPEHOUSE_REGISTRY;
const funderKey = process.env.PRIVATE_KEY as Hex | undefined;
if (!registry || !funderKey) throw new Error("Set TAPEHOUSE_REGISTRY and PRIVATE_KEY, as make test-app-devnode does.");
const deployments = parseDeployments(JSON.parse(readFileSync(registry, "utf8")));
const rpcUrl = process.env.TAPEHOUSE_RPC_URL ?? "http://127.0.0.1:8547";
const chain = {
  id: deployments.chainId,
  name: "dev node",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: [rpcUrl] } },
} as const;
const client = createPublicClient({ chain, transport: http(rpcUrl) });
const funder = createWalletClient({ account: privateKeyToAccount(funderKey), chain, transport: http(rpcUrl) });
const mintable = parseAbi(["function mint(address to, uint256 value)"]);
const usdg = deployments.tokens.USDG!;
const nvda = deployments.tokens.NVDA!;
const spy = deployments.tokens.SPY!;

test.describe.configure({ mode: "serial" });

/** Waits on the funder's transaction `write` sends, sending it again if a worker in another project took its nonce first. */
async function fund(write: () => Promise<Hex>) {
  for (let attempt = 0; ; attempt++) {
    try {
      return await client.waitForTransactionReceipt({ hash: await write() });
    } catch (error) {
      if (attempt === 4 || !String(error).includes("nonce")) throw error;
    }
  }
}

/** A wallet with `amounts` of tokens and no ether, connected to the app at `path`; its smart account's address. */
async function holder(page: Page, amounts: { usdg?: bigint; nvda?: bigint; spy?: bigint }, path = "/trade") {
  const key = generatePrivateKey();
  const owner = privateKeyToAccount(key).address;
  for (const [token, amount] of [
    [usdg, amounts.usdg],
    [nvda, amounts.nvda],
    [spy, amounts.spy],
  ] as const)
    if (amount) await fund(() => funder.writeContract({ address: token, abi: mintable, functionName: "mint", args: [owner, amount] }));
  await installWallet(page, { account: owner, chainId: deployments.chainId, rpcUrl, key });
  await page.goto(path);
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  return { owner, account: await sponsorship.accountAddress(client, deployments, owner) };
}

test.beforeAll(async () => {
  const lender = funder.account.address;
  const marginAccounts = deployments.tapehouse.MarginAccounts!;
  const lent = await accounts.lending(client, deployments, lender, CROSS, spy);
  if (lent.lent >= 5n * ONE) return;
  await fund(() => funder.writeContract({ address: spy, abi: mintable, functionName: "mint", args: [lender, 5n * ONE] }));
  await fund(() => funder.writeContract({ address: spy, abi: erc20Abi, functionName: "approve", args: [marginAccounts, 5n * ONE] }));
  await fund(() => funder.writeContract({ ...accounts.deposit(deployments, { position: CROSS, token: spy, amount: 5n * ONE, account: lender }) }));
  await fund(() => funder.writeContract({ ...accounts.lend(deployments, { position: CROSS, token: spy, amount: 5n * ONE, account: lender }) }));
});

const balance = (token: Hex, of: Hex) => client.readContract({ address: token, abi: erc20Abi, functionName: "balanceOf", args: [of] });
const ONE = 10n ** 18n;

test("a visitor reads every line of a short's ticket before connecting a wallet", async ({ page }) => {
  await page.goto("/trade");
  const ticket = page.getByRole("definition").first();
  await expect(ticket).toBeVisible();
  await page.getByLabel("Stock Token").selectOption("SPY");
  await page.getByLabel("Amount, sell").fill("1");
  const { proceeds, minProceeds } = await shorts.quoteSale(client, deployments, { asset: "SPY", amount: ONE, slippageBps: 50n });
  await expect(page.getByText(`$${usdgWords(proceeds)}`, { exact: true })).toBeVisible();
  await expect(page.getByText(`$${usdgWords(minProceeds)}`, { exact: true })).toBeVisible();
  await expect(page.getByRole("list", { name: "Before you sign" })).toContainText("only once it falls short at both edges");
  await page.getByRole("button", { name: "Gap cover" }).click();
  await expect(page.getByText(/last round before the market closes/)).toBeVisible();
  await page.getByRole("button", { name: "Baskets" }).click();
  await expect(page.getByText(/The target is .* per share/)).toBeVisible();
  const scan = await new AxeBuilder({ page }).analyze();
  expect(scan.violations).toEqual([]);
});

test("a short of SPY takes its margin from the wallet, sells at the quote and buys back, all for free", async ({ page }) => {
  test.setTimeout(240_000);
  const { owner, account } = await holder(page, { usdg: 1_000_000_000n });
  await page.getByRole("region", { name: "Smart account" }).waitFor();
  await page.getByLabel("Stock Token").selectOption("SPY");
  await page.getByLabel("Short action").selectOption("marginIn");
  await page.getByLabel("Amount, add margin").fill("900");
  await page.getByRole("button", { name: "Add margin, free" }).click();
  await expect(page.getByRole("status").filter({ hasText: "Done" })).toBeVisible({ timeout: 60_000 });
  expect((await shorts.position(client, deployments, account, "SPY")).usdgHeld).toBe(900_000_000n);
  expect(await balance(usdg, owner)).toBe(100_000_000n);

  await page.getByLabel("Short action").selectOption("sell");
  await page.getByLabel("Amount, sell").fill("1");
  const { minProceeds } = await shorts.quoteSale(client, deployments, { asset: "SPY", amount: ONE, slippageBps: 50n });
  await expect(page.getByText(`$${usdgWords(minProceeds)}`, { exact: true })).toBeVisible();
  await expect(page.getByRole("term").filter({ hasText: "Margin after" }).locator("..")).toContainText("against");
  await page.getByRole("button", { name: "Sell short, free" }).click();
  await expect(page.getByRole("status").filter({ hasText: "Done" })).toBeVisible({ timeout: 60_000 });
  const sold = await shorts.position(client, deployments, account, "SPY");
  expect(sold.debt).toBeGreaterThanOrEqual(ONE);
  expect(sold.usdgHeld).toBeGreaterThanOrEqual(900_000_000n + minProceeds);
  await expect(page.getByRole("term").filter({ hasText: "Short SPY" })).toBeVisible();

  await page.getByRole("button", { name: "Buy back", exact: true }).click();
  await expect(page.getByLabel("Short action")).toHaveValue("buyBack");
  await page.getByRole("button", { name: /^All it owes/ }).click();
  await expect(page.getByText("your limit alone, a full buy-back")).toBeVisible();
  await page.getByRole("button", { name: "Buy back, free" }).click();
  await expect(page.getByRole("status").filter({ hasText: "Done" })).toBeVisible({ timeout: 60_000 });
  expect((await shorts.position(client, deployments, account, "SPY")).debt).toBe(0n);
  expect(await sponsorship.freeOperationsLeft(client, deployments, account)).toBe(0n);
  expect(await client.getBalance({ address: owner })).toBe(0n);
});

test("a basket mints from the wallet's Stock Tokens into the cross position, margined as what it holds", async ({ page }) => {
  test.setTimeout(240_000);
  const { owner, account } = await holder(page, { nvda: 2n * ONE, spy: ONE });
  await page.getByRole("region", { name: "Smart account" }).waitFor();
  await page.getByRole("button", { name: "Baskets" }).click();
  await page.getByLabel("Shares to mint").fill("1");
  await expect(page.getByText(/1 NVDA \+ 0\.5 SPY/)).toBeVisible();
  await page.getByRole("button", { name: "Mint, free" }).click();
  await expect(page.getByRole("status").filter({ hasText: "Done" })).toBeVisible({ timeout: 60_000 });
  const basket = deployments.baskets.PAIR!;
  expect(await accounts.collateral(client, deployments, account, CROSS, basket)).toBe(ONE);
  const through = await accounts.inBaskets(client, deployments, account, CROSS);
  expect(through.NVDA).toBe(ONE);
  await expect(page.getByText(/margined as 1 NVDA \+ 0\.5 SPY through baskets/)).toBeVisible();
  await expect(page.getByRole("term").filter({ hasText: "Basket PAIR" })).toBeVisible();

  await page.getByRole("button", { name: "Redeem", exact: true }).click();
  await expect(page.getByLabel("Shares to redeem")).toHaveValue("1");
  await expect(page.getByText(/1 NVDA \+ 0\.5 SPY/).first()).toBeVisible();
  await page.getByRole("button", { name: "Redeem, free" }).click();
  await expect(page.getByRole("status").filter({ hasText: "Done" })).toBeVisible({ timeout: 60_000 });
  expect(await accounts.collateral(client, deployments, account, CROSS, basket)).toBe(0n);
  expect(await balance(nvda, owner)).toBe(2n * ONE);
  expect(await balance(spy, owner)).toBe(ONE);
});

test("gap cover is bought only while its sales are open", async ({ page }) => {
  await page.goto("/trade");
  await page.getByRole("button", { name: "Gap cover" }).click();
  const sales = await gapCover.sales(client, deployments);
  if (!sales) {
    await expect(page.getByText(/closed until the market reopens/)).toBeVisible();
    return;
  }
  await expect(page.getByText(/close in /)).toBeVisible();
});
