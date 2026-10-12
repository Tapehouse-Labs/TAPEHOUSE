// SPDX-License-Identifier: MIT OR Apache-2.0
import { stringToHex } from "viem";
import { describe, expect, test } from "vitest";
import { backstopActions, clearedRounds, fetchEvents, type Indexed } from "./history";

const event = (over: Partial<Indexed>): Indexed => ({
  block: 1,
  logIndex: 0,
  time: 1,
  tx: "0xabc",
  contract: "tapehouse.ReopeningAuction",
  event: "Cleared",
  args: {},
  finalized: true,
  ...over,
});

describe("clearedRounds", () => {
  test("reads the rounds the auction cleared, the latest open first", () => {
    const rounds = clearedRounds([
      event({ args: { symbol: stringToHex("SPY", { size: 32 }), openMs: "100", price: "60000000000", sold: "5", paid: "300", taken: "1" } }),
      event({ args: { symbol: stringToHex("TSLA", { size: 32 }), openMs: "200", price: "1", sold: "2", paid: "3", taken: "4" }, finalized: false }),
      event({ event: "Claimed" }),
    ]);
    expect(rounds.map((r) => [r.asset, r.openMs, r.price, r.finalized])).toEqual([
      ["TSLA", 200n, 1n, false],
      ["SPY", 100n, 60_000_000_000n, true],
    ]);
  });
});

describe("backstopActions", () => {
  test("reads what the backstop paid on a position and what it bought of the auction's remainder", () => {
    const actions = backstopActions([
      event({ contract: "tapehouse.GapBackstop", event: "Covered", args: { position: "0x00", closesMs: "10", paid: "7", written: "0" } }),
      event({ contract: "tapehouse.GapBackstop", event: "RemainderBought", args: { position: "0x01", closesMs: "20", cost: "9", amount: "1" } }),
      event({ event: "Sync" }),
    ]);
    expect(actions.map((a) => [a.kind, a.closesMs, a.amount])).toEqual([
      ["bought", 20n, 9n],
      ["covered", 10n, 7n],
    ]);
  });
});

describe("fetchEvents", () => {
  test("follows the pages the indexer names until the last", async () => {
    const calls: string[] = [];
    const fetcher = (async (url: URL) => {
      calls.push(url.search);
      const first = !url.searchParams.has("after");
      return Response.json({ events: [event({ block: first ? 1 : 2 })], next: first ? "1:0" : null });
    }) as unknown as typeof fetch;
    const events = await fetchEvents("https://indexer.example", { contract: "tapehouse.ReopeningAuction", event: "Cleared" }, fetcher);
    expect(events.map((e) => e.block)).toEqual([1, 2]);
    expect(calls[0]).toBe("?contract=tapehouse.ReopeningAuction&event=Cleared&limit=1000");
    expect(calls[1]).toContain("after=1%3A0");
  });

  test("refuses an answer that is not ok", async () => {
    const fetcher = (async () => new Response("no", { status: 429 })) as unknown as typeof fetch;
    await expect(fetchEvents("https://indexer.example", { contract: "x" }, fetcher)).rejects.toThrow("429");
  });
});
