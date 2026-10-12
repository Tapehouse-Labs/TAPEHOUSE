// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import { createServer } from "node:http";
import type { AddressInfo } from "node:net";
import { expect, test } from "@playwright/test";
import { auction, parseDeployments, toBytes32 } from "@tapehouse/sdk";
import { type Address, createPublicClient, http, keccak256 } from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { emptyRound, Overlay } from "./overlay";
import { installWallet } from "./wallet";

const registry = process.env.TAPEHOUSE_REGISTRY;
if (!registry) throw new Error("Set TAPEHOUSE_REGISTRY, as make test-app-devnode does.");
const deployments = parseDeployments(JSON.parse(readFileSync(registry, "utf8")));
const rpcUrl = process.env.TAPEHOUSE_RPC_URL ?? "http://127.0.0.1:8547";
const chain = { id: deployments.chainId, name: "dev node", nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 }, rpcUrls: { default: { http: [rpcUrl] } } } as const;
const client = createPublicClient({ chain, transport: http(rpcUrl) });
const asset = Object.keys(deployments.bandFeeds)[0]!;
const auctionAddress = deployments.tapehouse.ReopeningAuction as Address;

const money = (value: bigint, decimals: number) => {
  const unit = 10n ** BigInt(decimals);
  return `${value / unit}.${(value % unit).toString().padStart(decimals, "0")}`;
};

test("a bidder commits, reveals and claims in the reopening auction", async ({ page }) => {
  test.setTimeout(120_000);
  const overlay = new Overlay(rpcUrl, auctionAddress, deployments.tokens.USDG as Address);
  const server = createServer((request, response) => {
    let body = "";
    request.on("data", (chunk) => (body += chunk));
    request.on("end", async () => {
      response.setHeader("content-type", "application/json");
      response.end(JSON.stringify(await overlay.handle(JSON.parse(body))));
    });
  });
  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const overlayUrl = `http://127.0.0.1:${(server.address() as AddressInfo).port}`;
  const cors = { "access-control-allow-origin": "*", "access-control-allow-headers": "*", "access-control-allow-methods": "POST, OPTIONS" };
  await page.route(
    (url) => url.href.startsWith(rpcUrl),
    async (route) => {
      if (route.request().method() === "OPTIONS") return route.fulfill({ status: 204, headers: cors });
      const answer = await overlay.handle(JSON.parse(route.request().postData() ?? "[]"));
      return route.fulfill({ json: answer, headers: cors });
    },
  );

  const key = generatePrivateKey();
  const signer = privateKeyToAccount(key);
  await installWallet(page, { account: signer.address, chainId: deployments.chainId, rpcUrl: overlayUrl, key });

  const head = await client.getBlock();
  const openMs = (head.timestamp + 10n * 3_600n) * 1000n;
  const bond = 100_000_000n;
  const floor = await auction.previewFloor(client, deployments, { asset, openMs });
  const price = floor * 2n;
  const quantity = (300_000_000n * 10n ** 20n) / price;
  const deposit = 300_000_000n;
  overlay.world.phase = { openMs, revealing: false };
  overlay.world.usdg = { balance: 1_000_000_000n, allowance: 0n };

  await page.goto("/auction");
  await expect(page.getByText(/Bids are being committed/)).toBeVisible();
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await expect(page.getByText(/You hold 1,000.00 USDG/)).toBeVisible();

  await page.getByLabel("Asset to bid on").selectOption(asset);
  await page.getByLabel("Quantity to buy").fill(money(quantity, 18));
  await page.getByLabel("Highest price you will pay, in USD per token").fill(money(price, 8));
  await expect(page.getByLabel("Deposit, in USDG")).toHaveValue(money(deposit, 6).replace(/0+$/, "").replace(/\.$/, ""));
  expect(deposit).toBeGreaterThanOrEqual(bond);
  await page.getByRole("button", { name: "Commit bid" }).click();
  await expect.poll(() => overlay.sent.map((s) => s.name)).toEqual(["approve", "commit"]);

  const salt = keccak256(
    await signer.signMessage({ message: `Tapehouse reopening auction\n${asset} · regular open ${openMs} · bid 0\nSigning this derives the secret of your sealed bid. It sends nothing and costs nothing.` }),
  );
  const commitment = auction.commitment({ bidder: signer.address, asset, openMs, quantity, price, salt });
  const [, committed] = overlay.sent;
  expect(committed?.args).toEqual([toBytes32(asset), commitment, deposit]);
  await expect(page.getByRole("table", { name: /Your bids by round/ })).toContainText("Committed");

  overlay.world.shift = openMs / 1000n - 600n - head.timestamp;
  overlay.world.phase = { openMs, revealing: true };
  overlay.world.rounds.set(`${toBytes32(asset)}:${openMs}`, { round: { ...emptyRound(), floor, sealMs: openMs - 13n * 3_600_000n, deposits: deposit }, bids: [] });
  await page.reload();
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await page.getByRole("button", { name: "Reveal" }).click();
  await expect.poll(() => overlay.sent.map((s) => s.name)).toEqual(["approve", "commit", "reveal"]);
  expect(overlay.sent[2]?.args).toEqual([toBytes32(asset), openMs, quantity, price, salt]);

  overlay.world.phase = { openMs: 0n, revealing: false };
  overlay.world.commitments.clear();
  const bid = { bidder: signer.address, quantity, price, escrow: deposit, claimed: false };
  const key32 = `${toBytes32(asset)}:${openMs}`;
  overlay.world.rounds.set(key32, { round: { ...emptyRound(), floor, sealMs: 1n, supply: quantity }, bids: [bid] });
  overlay.world.lastOpenMs = openMs;
  overlay.world.shift = openMs / 1000n + 300n - head.timestamp;
  await page.reload();
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await page.getByRole("button", { name: new RegExp(`^Clear ${asset} at`) }).click();
  await expect.poll(() => overlay.sent.map((s) => s.name)).toEqual(["approve", "commit", "reveal", "clear"]);
  expect(overlay.sent[3]?.args).toEqual([toBytes32(asset), openMs, price]);

  overlay.world.rounds.set(key32, { round: { ...emptyRound(), floor, sealMs: 1n, supply: quantity, cleared: true, price, sold: quantity, paid: deposit }, bids: [bid] });
  await page.reload();
  await page.getByRole("button", { name: "Connect wallet" }).first().click();
  await page.getByRole("button", { name: "Claim" }).click();
  await expect.poll(() => overlay.sent.map((s) => s.name)).toEqual(["approve", "commit", "reveal", "clear", "claim"]);
  expect(overlay.sent[4]?.args).toEqual([toBytes32(asset), openMs, 0n]);
  server.close();
});
