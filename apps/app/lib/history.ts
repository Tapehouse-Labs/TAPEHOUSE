// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Hex, hexToString } from "viem";

/** An event as the indexer serves it. */
export type Indexed = {
  block: number;
  logIndex: number;
  time: number;
  tx: string;
  contract: string;
  event: string;
  args: Record<string, string>;
  finalized: boolean;
};

/** A reopening round the auction cleared. */
export type ClearedRound = {
  asset: string;
  openMs: bigint;
  price: bigint;
  sold: bigint;
  paid: bigint;
  taken: bigint;
  tx: string;
  finalized: boolean;
};

/** What the backstop covered in a closure: a position's shortfall paid, or what the auction's bids left bought. */
export type BackstopAction = {
  kind: "covered" | "bought";
  closesMs: bigint;
  position: Hex;
  amount: bigint;
  tx: string;
  finalized: boolean;
};

const bytes32 = (value: string | undefined): string => hexToString((value ?? "0x") as Hex, { size: 32 });

export function clearedRounds(events: readonly Indexed[]): ClearedRound[] {
  return events
    .filter((e) => e.event === "Cleared")
    .map((e) => ({
      asset: bytes32(e.args.symbol),
      openMs: BigInt(e.args.openMs ?? 0),
      price: BigInt(e.args.price ?? 0),
      sold: BigInt(e.args.sold ?? 0),
      paid: BigInt(e.args.paid ?? 0),
      taken: BigInt(e.args.taken ?? 0),
      tx: e.tx,
      finalized: e.finalized,
    }))
    .sort((a, b) => (a.openMs === b.openMs ? 0 : a.openMs > b.openMs ? -1 : 1));
}

export function backstopActions(events: readonly Indexed[]): BackstopAction[] {
  return events
    .flatMap((e): BackstopAction[] => {
      if (e.event === "Covered")
        return BigInt(e.args.paid ?? 0) === 0n && BigInt(e.args.written ?? 0) === 0n ? [] : [{ kind: "covered", closesMs: BigInt(e.args.closesMs ?? 0), position: (e.args.position ?? "0x") as Hex, amount: BigInt(e.args.paid ?? 0), tx: e.tx, finalized: e.finalized }];
      if (e.event === "RemainderBought")
        return [{ kind: "bought", closesMs: BigInt(e.args.closesMs ?? 0), position: (e.args.position ?? "0x") as Hex, amount: BigInt(e.args.cost ?? 0), tx: e.tx, finalized: e.finalized }];
      return [];
    })
    .sort((a, b) => (a.closesMs === b.closesMs ? 0 : a.closesMs > b.closesMs ? -1 : 1));
}

/** Every page of `event` of `contract` the indexer holds, oldest first, at most `limit` pages of 1,000. */
export async function fetchEvents(
  indexerUrl: string,
  query: { contract: string; event?: string },
  fetcher: typeof fetch = fetch,
  limit = 5,
): Promise<Indexed[]> {
  const events: Indexed[] = [];
  let after: string | null = null;
  for (let page = 0; page < limit; page++) {
    const url = new URL("/v1/events", indexerUrl);
    url.searchParams.set("contract", query.contract);
    if (query.event) url.searchParams.set("event", query.event);
    url.searchParams.set("limit", "1000");
    if (after) url.searchParams.set("after", after);
    const response = await fetcher(url);
    if (!response.ok) throw new Error(`The indexer answered ${response.status}.`);
    const body = (await response.json()) as { events: Indexed[]; next: string | null };
    events.push(...body.events);
    if (!body.next) break;
    after = body.next;
  }
  return events;
}
