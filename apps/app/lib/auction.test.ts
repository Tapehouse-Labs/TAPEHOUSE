// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import {
  auctionNotes,
  bidMove,
  bidProblem,
  bidRows,
  defaultDeposit,
  bidStatus,
  deserialize,
  escrowOf,
  nextNonce,
  type SavedBid,
  saltFrom,
  saltMessage,
  serialize,
  stage,
  worth,
} from "./auction";
import { say } from "./say";

const E18 = 10n ** 18n;
const open = 1_791_000_000_000n;
const clearMs = 3_600_000n;
const live = { floor: 90n, cleared: false };

describe("stage", () => {
  test("bids are committed, then revealed, as the auction's phase says", () => {
    const base = { lastOpenMs: 0n, round: undefined, nowMs: 0, clearMs };
    expect(stage({ ...base, phase: { openMs: open, revealing: false } })).toEqual({ stage: "committing", openMs: open, endsMs: undefined });
    expect(stage({ ...base, phase: { openMs: open, revealing: true } })).toEqual({ stage: "revealing", openMs: open, endsMs: open });
  });

  test("after the open the round awaits its clearing price, then lapses or is cleared", () => {
    const base = { phase: { openMs: 0n, revealing: false }, lastOpenMs: open, clearMs };
    expect(stage({ ...base, round: live, nowMs: Number(open) + 1000 })).toEqual({ stage: "clearing", openMs: open, endsMs: open + clearMs });
    expect(stage({ ...base, round: live, nowMs: Number(open + clearMs) }).stage).toBe("lapsed");
    expect(stage({ ...base, round: { floor: 90n, cleared: true }, nowMs: Number(open) + 1000 }).stage).toBe("cleared");
  });

  test("with no round opened there is nothing to show", () => {
    const base = { phase: { openMs: 0n, revealing: false }, clearMs, nowMs: 5 };
    expect(stage({ ...base, lastOpenMs: 0n, round: undefined }).stage).toBe("idle");
    expect(stage({ ...base, lastOpenMs: open, round: { floor: 0n, cleared: false } }).stage).toBe("idle");
  });
});

describe("bids", () => {
  const terms = { floor: 90n * 10n ** 8n, bond: 100_000_000n, minBid: 100_000_000n };

  test("a bid escrows what it is worth at its price, rounded up", () => {
    expect(escrowOf(E18, 60_000_000_000n)).toBe(600_000_000n);
    expect(escrowOf(1n, 1n)).toBe(1n);
    expect(worth(E18 / 3n, 60_000_000_000n)).toBe(199_999_999n);
  });

  test("the auction refuses a bid under the floor, too small, under the bond or not covered", () => {
    const ok = { quantity: 2n * E18, price: 100n * 10n ** 8n, deposit: 300_000_000n, ...terms };
    expect(bidProblem(ok)).toBeUndefined();
    expect(say(bidProblem({ ...ok, price: 89n * 10n ** 8n })!)).toMatch(/floor/);
    expect(say(bidProblem({ ...ok, quantity: E18 / 2n })!)).toMatch(/too small|least/i);
    expect(say(bidProblem({ ...ok, deposit: 99_000_000n })!)).toMatch(/bond|least/i);
    expect(say(bidProblem({ ...ok, deposit: 150_000_000n })!)).toMatch(/escrow|cover/i);
  });
});

describe("saved bids", () => {
  const bid: SavedBid = {
    asset: "SPY",
    openMs: open,
    quantity: 5n * 10n ** 17n,
    price: 60_000_000_000n,
    deposit: 300_000_000n,
    nonce: 1,
    salt: "0x01",
    commitment: "0x02",
  };

  test("they survive a round trip through text, with their amounts exact", () => {
    expect(deserialize(serialize([bid, { ...bid, nonce: 2 }]))).toEqual([bid, { ...bid, nonce: 2 }]);
  });

  test("text that is not them gives none, and a broken entry is dropped alone", () => {
    expect(deserialize(null)).toEqual([]);
    expect(deserialize("{")).toEqual([]);
    expect(deserialize('{"a":1}')).toEqual([]);
    const text = JSON.stringify([JSON.parse(serialize([bid]))[0], { asset: "SPY", nonce: 0, salt: "0x", commitment: "0x", openMs: "x" }]);
    expect(deserialize(text)).toEqual([bid]);
  });

  test("the next nonce is free among the bids of the round", () => {
    expect(nextNonce([], "SPY", open)).toBe(0);
    expect(nextNonce([bid, { ...bid, nonce: 4 }, { ...bid, asset: "TSLA", nonce: 9 }], "SPY", open)).toBe(5);
    expect(nextNonce([bid], "SPY", open + 1n)).toBe(0);
  });

  test("a salt is the hash of the signature of the words that name the bid", () => {
    expect(saltMessage({ asset: "SPY", openMs: open, nonce: 1 })).toContain("SPY · regular open 1791000000000 · bid 1");
    expect(saltFrom("0x1234")).toBe("0x56570de287d73cd1cb6092bb8fdee6173974955fdef345ae579ee9f475ea7432");
  });

  test("a bid is committed, revealed, claimed or returned", () => {
    const listed = { quantity: bid.quantity, price: bid.price, claimed: false };
    expect(bidStatus({ deposit: 1n, listed: undefined })).toBe("committed");
    expect(bidStatus({ deposit: 0n, listed })).toBe("revealed");
    expect(bidStatus({ deposit: 0n, listed: { ...listed, claimed: true } })).toBe("claimed");
    expect(bidStatus({ deposit: 0n, listed: undefined })).toBe("returned");
  });
});

describe("defaultDeposit", () => {
  const bond = 100_000_000n;
  test("is the bond until the bid escrows more", () => {
    expect(defaultDeposit({ quantity: undefined, price: undefined, bond })).toBe(bond);
    expect(defaultDeposit({ quantity: E18, price: 50n * 10n ** 8n, bond })).toBe(bond);
    expect(defaultDeposit({ quantity: 4n * E18, price: 50n * 10n ** 8n, bond })).toBe(200_000_000n);
  });
});

describe("bidMove", () => {
  const base = { openMs: open, cleared: false, revealMs: 1_800_000n, clearMs };
  const at = (ms: bigint) => Number(open + ms);
  test("a commitment waits, is revealed in the half hour before the open, and is forfeited after it", () => {
    expect(bidMove({ ...base, status: "committed", nowMs: at(-3_600_000n) })).toBe("waitReveal");
    expect(bidMove({ ...base, status: "committed", nowMs: at(-1_800_000n) })).toBe("reveal");
    expect(bidMove({ ...base, status: "committed", nowMs: at(0n) })).toBe("forfeit");
  });
  test("a revealed bid is claimed once the round is cleared or its window has closed", () => {
    expect(bidMove({ ...base, status: "revealed", nowMs: at(1000n) })).toBe("waitClear");
    expect(bidMove({ ...base, status: "revealed", cleared: true, nowMs: at(1000n) })).toBe("claim");
    expect(bidMove({ ...base, status: "revealed", nowMs: at(clearMs) })).toBe("claim");
  });
  test("a claimed or returned bid has nothing left to do", () => {
    expect(bidMove({ ...base, status: "claimed", nowMs: at(0n) })).toBe("none");
    expect(bidMove({ ...base, status: "returned", nowMs: at(0n) })).toBe("none");
  });
});

describe("bidRows", () => {
  const saved = (asset: string, index: number | undefined, openMs = open) => ({
    asset,
    openMs,
    quantity: E18,
    price: 5n,
    deposit: 1n,
    nonce: 0,
    salt: "0x01" as const,
    commitment: "0x02" as const,
    status: "revealed" as const,
    index,
  });
  const listed = (asset: string, index: number, claimed = false, openMs = open) => ({ asset, openMs, index, quantity: E18, price: 5n, claimed });
  test("lists a revealed bid once, saved or not, newest round first", () => {
    const rows = bidRows([saved("TSLA", 0), saved("AAPL", undefined, open - 1n)], [listed("TSLA", 0), listed("NVDA", 1, true)]);
    expect(rows.map((r) => [r.asset, r.status, r.saved !== undefined])).toEqual([
      ["NVDA", "claimed", false],
      ["TSLA", "revealed", true],
      ["AAPL", "revealed", true],
    ]);
  });
});

describe("auctionNotes", () => {
  test("say how the auction works in every language", () => {
    for (const locale of ["en", "es", "pt"] as const)
      for (const note of auctionNotes()) {
        expect(say(note.message, locale).length).toBeGreaterThan(20);
        expect(note.term && say(note.term, locale).length).toBeGreaterThan(2);
      }
  });
});
