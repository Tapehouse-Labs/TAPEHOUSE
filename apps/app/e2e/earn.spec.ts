// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import AxeBuilder from "@axe-core/playwright";
import { expect, type Page, test } from "@playwright/test";
import { backstop, gapCover, parseDeployments, sponsorship, supply } from "@tapehouse/sdk";
import { createPublicClient, createWalletClient, erc20Abi, type Hex, http, parseAbi } from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
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
async function holder(page: Page, amounts: { usdg?: bigint; nvda?: bigint }, path = "/earn") {
  const key = generatePrivateKey();
  const owner = privateKeyToAccount(key).address;
  for (const [token, amount] of [
    [usdg, amounts.usdg],
    [nvda, amounts.nvda],
  ] as const)
    if (amount) await fund(() => funder.writeContract({ address: token, abi: mintable, functionName: "mint", args: [owner, amount] }));
  await installWallet(page, { account: owner, chainId: deployments.chainId, rpcUrl, key });
  await page.goto(path);
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  return { owner, account: await sponsorship.accountAddress(client, deployments, owner) };
}

const balance = (token: Hex, of: Hex) => client.readContract({ address: token, abi: erc20Abi, functionName: "balanceOf", args: [of] });

test("a visitor reads every place's terms before connecting a wallet", async ({ page }) => {
  await page.goto("/earn");
  const ledger = page.getByRole("table", { name: /Each place your money can work/ });
  await expect(ledger.getByRole("rowheader")).toContainText(["Supply", "Backstop", "Lend NVDA", "Write cover"]);
  const market = await supply.market(client, deployments);
  await expect(ledger.getByRole("button", { name: "Supply" })).toHaveAttribute("aria-expanded", "true");
  await expect(page.getByRole("list", { name: "Before you deposit" })).toContainText("Lenders are last in the loss cascade");
  await expect(page.getByText(/You may withdraw at most the USDG the vault holds/)).toBeVisible();
  expect(market.supplyRate).toBeGreaterThanOrEqual(0n);

  await ledger.getByRole("button", { name: "Backstop" }).click();
  await expect(page.getByRole("table", { name: /Exposure this closure/ }).getByRole("rowheader")).toContainText(["Cross positions"]);
  await expect(page.getByText(/The team's seed is an ordinary deposit/)).toBeVisible();

  await ledger.getByRole("button", { name: "Lend NVDA" }).click();
  await expect(page.getByText(/then a buy-in anyone may call/)).toBeVisible();
  await ledger.getByRole("button", { name: "Write cover" }).click();
  await expect(page.getByText(/Every cover reserves its whole payout/)).toBeVisible();

  const scan = await new AxeBuilder({ page }).analyze();
  expect(scan.violations).toEqual([]);
});

test("a holder supplies USDG for free and withdraws it to the wallet", async ({ page }) => {
  test.setTimeout(180_000);
  const { owner, account } = await holder(page, { usdg: 500_000_000n });
  const panel = page.locator("#panel-supply");
  await panel.getByLabel("Amount to deposit").fill("300");
  await panel.getByRole("button", { name: "Deposit, free" }).click();
  await expect(panel.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  const lent = await supply.lender(client, deployments, account);
  expect(lent.assets).toBeGreaterThanOrEqual(299_999_999n);
  expect(await balance(usdg, owner)).toBe(200_000_000n);
  await expect(page.locator(".figure")).toHaveText(/^\$(299\.99|300\.00)$/);

  await panel.getByRole("button", { name: "Withdraw", exact: true }).click();
  await panel.getByRole("button", { name: /^Max / }).click();
  await panel.getByRole("button", { name: "Withdraw, free" }).click();
  await expect(panel.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  expect((await supply.lender(client, deployments, account)).shares).toBe(0n);
  expect(await balance(usdg, owner)).toBeGreaterThanOrEqual(499_999_999n);
  expect(await client.getBalance({ address: owner })).toBe(0n);
});

test("a depositor starts the backstop's cooldown, and deposits only while they are open", async ({ page }) => {
  test.setTimeout(180_000);
  const key = generatePrivateKey();
  const owner = privateKeyToAccount(key).address;
  const account = await sponsorship.accountAddress(client, deployments, owner);
  const vault = deployments.tapehouse.GapBackstop!;
  await fund(() => funder.writeContract({ address: vault, abi: erc20Abi, functionName: "transfer", args: [account, 50n * 10n ** 12n] }));
  await fund(() => funder.writeContract({ address: usdg, abi: mintable, functionName: "mint", args: [owner, 100_000_000n] }));
  await installWallet(page, { account: owner, chainId: deployments.chainId, rpcUrl, key });
  await page.goto("/earn");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();

  await page.getByRole("button", { name: "Backstop" }).click();
  const panel = page.locator("#panel-backstop");
  const open = (await backstop.depositor(client, deployments, account, [])).maxDeposit > 0n;
  if (open) await expect(panel.getByLabel("Amount to deposit")).toBeEditable();
  else await expect(panel.getByText(/Deposits are closed now/)).toBeVisible();

  await panel.getByRole("button", { name: "Leave" }).click();
  await expect(panel.getByText(/Start a cooldown and your 50.00 USDG of shares may leave in a week/)).toBeVisible();
  await panel.getByRole("button", { name: "Start cooldown, free" }).click();
  await expect(panel.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  const { cooldown } = await backstop.depositor(client, deployments, account, []);
  expect(cooldown?.shares).toBe(50n * 10n ** 12n);
  await expect(panel.getByText(/Your cooldown runs: your shares may leave in (7 days|6 days, \d+ hours?), for 6 days/)).toBeVisible();
});

test("a position lends NVDA from the Earn ledger, with the recall order stated first", async ({ page }) => {
  test.setTimeout(240_000);
  const { account } = await holder(page, { nvda: 6n * 10n ** 18n }, "/");
  const act = page.getByRole("region", { name: "Act on your account" });
  await act.getByLabel("Amount to deposit").fill("6");
  await act.getByLabel("Token").selectOption({ label: "NVDA" });
  await act.getByRole("button", { name: "Deposit, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });

  await page.getByRole("navigation", { name: "Sections" }).getByRole("link", { name: "Earn" }).filter({ visible: true }).click();
  await page.getByRole("button", { name: "Lend NVDA" }).click();
  const panel = page.locator("#panel-lend-NVDA");
  await expect(panel.getByText(/No one takes idle tokens first come, first served/)).toBeVisible();
  await panel.getByLabel("Amount to lend").fill("4");
  await panel.getByRole("button", { name: "Lend, free" }).click();
  await expect(panel.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  await expect(page.getByRole("row", { name: /Lend NVDA/ }).first()).toContainText("4 NVDA");
  const shares = await client.readContract({
    address: deployments.stockLending.NVDA!,
    abi: erc20Abi,
    functionName: "balanceOf",
    args: [deployments.tapehouse.MarginAccounts!],
  });
  expect(shares).toBeGreaterThan(0n);
  expect(await sponsorship.freeOperationsLeft(client, deployments, account)).toBe(1n);
});

test("a writer deposits with the gap cover only while its sales are open", async ({ page }) => {
  test.setTimeout(180_000);
  const { account } = await holder(page, { usdg: 100_000_000n });
  await page.getByRole("button", { name: "Write cover" }).click();
  const panel = page.locator("#panel-cover");
  const open = (await gapCover.writer(client, deployments, account)).maxDeposit > 0n;
  if (!open) {
    await expect(panel.getByText(/Deposits are closed now: they open only while the session is open/)).toBeVisible();
    await panel.getByLabel("Amount to deposit").fill("50");
    await expect(panel.getByRole("button", { name: "Deposit, free" })).toBeDisabled();
    return;
  }
  await panel.getByLabel("Amount to deposit").fill("50");
  await panel.getByRole("button", { name: "Deposit, free" }).click();
  await expect(panel.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  expect((await gapCover.writer(client, deployments, account)).assets).toBe(50_000_000n);
});
