// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import { expect, type Page, test } from "@playwright/test";
import { accounts, band, CROSS, engine, parseDeployments, sponsorship } from "@tapehouse/sdk";
import { createPublicClient, createWalletClient, erc20Abi, type Hex, http, parseAbi } from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { usd } from "../lib/amounts";
import { price } from "../lib/format";
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
const nvda = deployments.tokens.NVDA!;
const usdg = deployments.tokens.USDG!;
const ONE = 10n ** 18n;

test.describe.configure({ mode: "serial" });

/** Mints `amount` of `token` to `to` from the funder, again if a worker in another project took its nonce first. */
async function mint(token: Hex, to: Hex, amount: bigint) {
  for (let attempt = 0; ; attempt++) {
    try {
      const hash = await funder.writeContract({ address: token, abi: mintable, functionName: "mint", args: [to, amount] });
      return await client.waitForTransactionReceipt({ hash });
    } catch (error) {
      if (attempt === 4 || !String(error).includes("nonce")) throw error;
    }
  }
}

/** A wallet with `tokens` NVDA and no ether, connected to the app. */
async function holder(page: Page, tokens: bigint) {
  const key = generatePrivateKey();
  const owner = privateKeyToAccount(key).address;
  await mint(nvda, owner, tokens * ONE);
  await installWallet(page, { account: owner, chainId: deployments.chainId, rpcUrl, key });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  return owner;
}

test("a holder with Stock Tokens and no ether deposits and borrows for free, knowing the liquidation price first", async ({ page }) => {
  test.setTimeout(180_000);
  const owner = await holder(page, 20n);
  await expect(page.getByText("None yet")).toBeVisible();

  const act = page.getByRole("region", { name: "Act on your account" });
  await act.getByLabel("Amount to deposit").fill("20");
  await act.getByLabel("Token").selectOption({ label: "NVDA" });
  await act.getByRole("button", { name: "Deposit, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });

  const figure = page.locator(".figure");
  await expect(figure).toHaveText("None");

  await act.getByRole("button", { name: "Borrow" }).first().click();
  await act.getByLabel("Amount to borrow").fill("1000");
  const account = (await page.getByRole("region", { name: "Smart account" }).locator(".addr").getAttribute("title")) as Hex;
  const expected = await accounts.liquidationPrice(client, deployments, account, CROSS, "NVDA", 1_000_000_000n);
  await expect(figure).toHaveText(`$${price(expected)}`);
  await expect(page.getByRole("list", { name: /requirement against its equity/ })).toBeVisible();
  await act.getByRole("button", { name: "Borrow, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });

  expect(await client.readContract({ address: usdg, abi: erc20Abi, functionName: "balanceOf", args: [owner] })).toBe(1_000_000_000n);
  expect(await client.getBalance({ address: owner })).toBe(0n);
  expect(await sponsorship.freeOperationsLeft(client, deployments, account)).toBe(1n);
  const { debt } = await accounts.repayment(client, deployments, account, CROSS);
  expect(debt).toBeGreaterThanOrEqual(1_000_000_000n);
});

test("the third free operation repays from the wallet's USDG with a permit", async ({ page }) => {
  test.setTimeout(180_000);
  const owner = await holder(page, 10n);
  const act = page.getByRole("region", { name: "Act on your account" });
  await act.getByLabel("Amount to deposit").fill("10");
  await act.getByLabel("Token").selectOption({ label: "NVDA" });
  await act.getByRole("button", { name: "Deposit, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  await act.getByRole("button", { name: "Borrow" }).first().click();
  await act.getByLabel("Amount to borrow").fill("300");
  await act.getByRole("button", { name: "Borrow, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });

  await act.getByRole("button", { name: "Repay" }).first().click();
  await act.getByLabel("Amount to repay").fill("120");
  await act.getByRole("button", { name: "Repay, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });
  const account = (await page.getByRole("region", { name: "Smart account" }).locator(".addr").getAttribute("title")) as Hex;
  expect(await client.readContract({ address: usdg, abi: erc20Abi, functionName: "balanceOf", args: [owner] })).toBe(180_000_000n);
  const { debt } = await accounts.repayment(client, deployments, account, CROSS);
  expect(debt).toBeGreaterThanOrEqual(180_000_000n);
  expect(debt).toBeLessThan(180_010_000n);
  expect(await sponsorship.freeOperationsLeft(client, deployments, account)).toBe(0n);
  await expect(page.getByRole("region", { name: "Smart account" })).toContainText("Free actions used");
});

test("a position lends its NVDA and takes it back, its equity shown with the recall haircut first", async ({ page }) => {
  test.setTimeout(180_000);
  await holder(page, 8n);
  const act = page.getByRole("region", { name: "Act on your account" });
  await act.getByLabel("Amount to deposit").fill("8");
  await act.getByLabel("Token").selectOption({ label: "NVDA" });
  await act.getByRole("button", { name: "Deposit, free" }).click();
  await expect(act.getByRole("status")).toContainText("Done", { timeout: 60_000 });

  const position = page.locator("details.position").first();
  await position.locator("summary").click();
  await position.getByLabel("NVDA to lend").fill("5");
  await expect(position).toContainText("a lent token counts 5% less in it");
  await position.getByRole("button", { name: "Lend", exact: true }).last().click();
  await expect(position.getByRole("definition").nth(5)).toHaveText("5", { timeout: 60_000 });

  await position.getByRole("button", { name: "Take back" }).click();
  await position.getByLabel("NVDA to take back").fill("5");
  await position.getByRole("button", { name: "Take back" }).last().click();
  await expect(position.getByRole("definition").nth(5)).toHaveText("0", { timeout: 60_000 });
  await expect(position.getByRole("definition").nth(4)).toHaveText("8");
});

test("the simulator margins any portfolio with the engine, no wallet needed", async ({ page }) => {
  await page.goto("/simulator");
  await expect(page.getByText("Any portfolio")).toBeVisible();
  await page.getByLabel("NVDA tokens").fill("20");
  const assets = await accounts.stocks(client, deployments);
  const [{ low }] = await Promise.all([band.quote(client, deployments, "NVDA")]);
  const quantities = assets.map((a) => (a.asset === "NVDA" ? 20n * ONE : 0n));
  const prices = await Promise.all(assets.map(async (a) => (await band.quote(client, deployments, a.asset)).low));
  const { margin } = await engine.currentRequirement(client, deployments, { quantities, prices });
  const equity = 20n * low * 10n ** 10n;
  await expect(page.getByText(`${usd(margin)} required against ${usd(equity)} of equity`)).toBeVisible();
  await expect(page.getByRole("list", { name: /requirement against its equity through the week/ })).toBeVisible();
});
