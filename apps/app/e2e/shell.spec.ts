// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import AxeBuilder from "@axe-core/playwright";
import { expect, test } from "@playwright/test";
import { band, parseDeployments } from "@tapehouse/sdk";
import { createPublicClient, http } from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { bandPercent, price, shortAddress } from "../lib/format";
import { bandState } from "../lib/session";
import { installWallet, moveWallet } from "./wallet";

const registry = process.env.TAPEHOUSE_REGISTRY;
if (!registry) throw new Error("Set TAPEHOUSE_REGISTRY to the dev node's registry, as make test-app-devnode does.");
const deployments = parseDeployments(JSON.parse(readFileSync(registry, "utf8")));
const rpcUrl = process.env.TAPEHOUSE_RPC_URL ?? "http://127.0.0.1:8547";
const client = createPublicClient({ transport: http(rpcUrl) });
const account = "0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E";

test("with no wallet, the app says so and still shows every band as the chain reports it", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("link", { name: "Get a wallet" }).first()).toBeVisible();
  for (const asset of Object.keys(deployments.bandFeeds)) {
    const q = await band.quote(client, deployments, asset);
    const row = page.getByRole("row", { name: new RegExp(`^${asset}\\b`) });
    if (q.mid > 0n) {
      await expect(row).toContainText(price(q.mid));
      await expect(row).toContainText(bandPercent(q.halfBps));
    } else {
      await expect(row).toContainText(bandState(q.state));
    }
  }
});

test("the seal reads the band's session", async ({ page }) => {
  const session = await band.session(client, deployments);
  await page.goto("/");
  const state = session.state === 2 ? "OPEN" : session.state === 1 ? "CLOSED" : "UNKNOWN";
  await expect(page.getByRole("status").filter({ hasText: state }).first()).toBeVisible();
});

test("a wallet connects on the app's chain", async ({ page }) => {
  await installWallet(page, { account, chainId: deployments.chainId, rpcUrl });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await expect(page.getByRole("button", { name: shortAddress(account) })).toBeVisible();
  await expect(page.getByText("appears here the moment you deposit")).toBeVisible();
});

test("a wallet with no ether creates its smart account for free", async ({ page }) => {
  const key = generatePrivateKey();
  const owner = privateKeyToAccount(key).address;
  await installWallet(page, { account: owner, chainId: deployments.chainId, rpcUrl, key });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  const smart = page.getByRole("region", { name: "Smart account" });
  await expect(smart).toContainText("NOT CREATED");
  await expect(smart).toContainText("Creating it costs you nothing, and so do your first 3 actions.");
  await smart.getByRole("button", { name: "Create account" }).click();
  await expect(smart).toContainText("LIVE", { timeout: 30_000 });
  await expect(smart).toContainText("3 free actions left");
  await expect(smart.getByRole("button")).toHaveCount(0);
  expect(await client.getBalance({ address: owner })).toBe(0n);
});

test("a wallet on another chain is moved to the app's chain as it connects", async ({ page }) => {
  await installWallet(page, { account, chainId: 1, rpcUrl });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await expect(page.getByText("appears here the moment you deposit")).toBeVisible();
});

test("a wallet that moves to another chain is asked to switch back, and switches", async ({ page }) => {
  await installWallet(page, { account, chainId: deployments.chainId, rpcUrl });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await expect(page.getByText("appears here the moment you deposit")).toBeVisible();
  await moveWallet(page, 1);
  const switchBack = page.getByRole("button", { name: /^Switch to / });
  await expect(page.getByText("Your wallet is on another network.")).toBeVisible();
  await switchBack.click();
  await expect(page.getByText("appears here the moment you deposit")).toBeVisible();
});

test("a declined connection is explained", async ({ page }) => {
  await installWallet(page, { account, chainId: deployments.chainId, rpcUrl, decline: true });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await expect(page.getByRole("alert").filter({ hasText: "declined" })).toHaveText(
    "You declined the request in your wallet.",
  );
});

test("a disconnect returns to the connect state", async ({ page }) => {
  await installWallet(page, { account, chainId: deployments.chainId, rpcUrl });
  await page.goto("/");
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await page.getByRole("button", { name: shortAddress(account) }).click();
  await page.getByRole("button", { name: "Disconnect" }).click();
  await expect(page.getByRole("button", { name: "Connect wallet" }).first()).toBeVisible();
});

for (const colorScheme of ["light", "dark"] as const) {
  test(`meets WCAG 2.2 AA in the ${colorScheme} theme, connected, with its smart account, and not`, async ({ page }) => {
    await page.emulateMedia({ colorScheme, reducedMotion: "reduce" });
    await installWallet(page, { account, chainId: deployments.chainId, rpcUrl });
    await page.goto("/");
    await expect(page.getByRole("button", { name: "Connect wallet" }).first()).toBeVisible();
    const tags = ["wcag2a", "wcag2aa", "wcag21a", "wcag21aa", "wcag22aa"];
    expect((await new AxeBuilder({ page }).withTags(tags).analyze()).violations).toEqual([]);
    await page.getByRole("button", { name: "Connect wallet" }).first().click();
    await expect(page.getByRole("region", { name: "Smart account" })).toContainText("NOT CREATED");
    expect((await new AxeBuilder({ page }).withTags(tags).analyze()).violations).toEqual([]);
    await moveWallet(page, 1);
    await expect(page.getByRole("button", { name: /^Switch to / })).toBeVisible();
    expect((await new AxeBuilder({ page }).withTags(tags).analyze()).violations).toEqual([]);
  });
}
