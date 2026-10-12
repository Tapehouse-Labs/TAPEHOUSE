// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Hex, keccak256 } from "viem";
import { type Msg, msg } from "./messages";
import type { Note } from "./notes";

const TOKEN_TO_USDG = 10n ** 20n;

/** Where a round stands: bids are committed, revealed, the clearing price awaited, or the round is over. */
export type Stage = "idle" | "committing" | "revealing" | "clearing" | "cleared" | "lapsed";

/**
 * Where the round the auction shows stands at `nowMs`. `phase` is the auction's own reading, which names the regular open
 * `openMs` while bids are committed or revealed; once it reads zero, the round of `lastOpenMs` is clearing until its
 * window of `clearMs` closes, or is over.
 */
export function stage(args: {
  phase: { openMs: bigint; revealing: boolean };
  lastOpenMs: bigint;
  round: { floor: bigint; cleared: boolean } | undefined;
  nowMs: number;
  clearMs: bigint;
}): { stage: Stage; openMs: bigint; endsMs: bigint | undefined } {
  const { phase, lastOpenMs, round, nowMs, clearMs } = args;
  if (phase.openMs !== 0n) {
    return {
      stage: phase.revealing ? "revealing" : "committing",
      openMs: phase.openMs,
      endsMs: phase.revealing ? phase.openMs : undefined,
    };
  }
  if (lastOpenMs === 0n || !round || round.floor === 0n) return { stage: "idle", openMs: lastOpenMs, endsMs: undefined };
  if (round.cleared) return { stage: "cleared", openMs: lastOpenMs, endsMs: undefined };
  const closes = lastOpenMs + clearMs;
  if (BigInt(nowMs) < closes) return { stage: "clearing", openMs: lastOpenMs, endsMs: closes };
  return { stage: "lapsed", openMs: lastOpenMs, endsMs: undefined };
}

/** What a bid escrows once revealed, in USDG units: what it is worth at its price, rounded up. */
export function escrowOf(quantity: bigint, price: bigint): bigint {
  return (quantity * price + TOKEN_TO_USDG - 1n) / TOKEN_TO_USDG;
}

/** What a bid of `quantity` is worth at `price`, in USDG units, rounded down. */
export function worth(quantity: bigint, price: bigint): bigint {
  return (quantity * price) / TOKEN_TO_USDG;
}

/**
 * Why the auction would refuse a bid, if it would: a price below the round's floor, a bid worth less than the least
 * at the floor, a deposit below the bond, or one that does not cover what the bid escrows. A deposit larger than the
 * bid hides its size.
 */
export function bidProblem(args: {
  quantity: bigint;
  price: bigint;
  deposit: bigint;
  floor: bigint;
  bond: bigint;
  minBid: bigint;
}): Msg | undefined {
  const { quantity, price, deposit, floor, bond, minBid } = args;
  if (price < floor) return msg("auction.problem.belowFloor");
  if (worth(quantity, floor) < minBid) return msg("auction.problem.tooSmall");
  if (deposit < bond) return msg("auction.problem.belowBond");
  if (deposit < escrowOf(quantity, price)) return msg("auction.problem.belowEscrow");
  return undefined;
}

/** A bid the visitor committed, kept in the browser until it is revealed, with what is needed to reveal it. */
export type SavedBid = {
  asset: string;
  openMs: bigint;
  quantity: bigint;
  price: bigint;
  deposit: bigint;
  nonce: number;
  salt: Hex;
  commitment: Hex;
};

type Stored = Omit<SavedBid, "openMs" | "quantity" | "price" | "deposit"> & {
  openMs: string;
  quantity: string;
  price: string;
  deposit: string;
};

/** The saved bids as text, with each amount in decimal. */
export function serialize(bids: readonly SavedBid[]): string {
  return JSON.stringify(
    bids.map(
      (b): Stored => ({
        ...b,
        openMs: b.openMs.toString(),
        quantity: b.quantity.toString(),
        price: b.price.toString(),
        deposit: b.deposit.toString(),
      }),
    ),
  );
}

/** The bids `serialize` wrote; none for text that is not them. */
export function deserialize(text: string | null): SavedBid[] {
  if (!text) return [];
  try {
    const parsed: unknown = JSON.parse(text);
    if (!Array.isArray(parsed)) return [];
    return parsed.flatMap((entry: Partial<Stored>) => {
      try {
        if (
          typeof entry.asset !== "string" ||
          typeof entry.nonce !== "number" ||
          typeof entry.salt !== "string" ||
          typeof entry.commitment !== "string"
        )
          return [];
        return [
          {
            asset: entry.asset,
            nonce: entry.nonce,
            salt: entry.salt as Hex,
            commitment: entry.commitment as Hex,
            openMs: BigInt(entry.openMs ?? ""),
            quantity: BigInt(entry.quantity ?? ""),
            price: BigInt(entry.price ?? ""),
            deposit: BigInt(entry.deposit ?? ""),
          },
        ];
      } catch {
        return [];
      }
    });
  } catch {
    return [];
  }
}

/** The words a wallet signs to derive the salt of a commitment, so the same wallet can derive it again. */
export function saltMessage(args: { asset: string; openMs: bigint; nonce: number }): string {
  return `Tapehouse reopening auction\n${args.asset} · regular open ${args.openMs} · bid ${args.nonce}\nSigning this derives the secret of your sealed bid. It sends nothing and costs nothing.`;
}

/** The salt a signature of `saltMessage` gives. */
export function saltFrom(signature: Hex): Hex {
  return keccak256(signature);
}

/** The next free nonce for a bid of `asset` in the round of `openMs` among `saved`. */
export function nextNonce(saved: readonly SavedBid[], asset: string, openMs: bigint): number {
  const used = saved.filter((b) => b.asset === asset && b.openMs === openMs).map((b) => b.nonce);
  return used.length === 0 ? 0 : Math.max(...used) + 1;
}

/** Where a saved bid stands, from what the auction holds for it. */
export type BidStatus = "committed" | "revealed" | "claimed" | "returned";

/**
 * A saved bid's standing: committed while the auction still holds its deposit; revealed once the round lists a bid of
 * the bidder at its price and quantity; claimed after its payout; returned when the deposit is gone and no such bid is
 * listed, as after a forfeit or an outbid reveal.
 */
export function bidStatus(args: {
  deposit: bigint;
  listed: { quantity: bigint; price: bigint; claimed: boolean } | undefined;
}): BidStatus {
  if (args.deposit > 0n) return "committed";
  if (!args.listed) return "returned";
  return args.listed.claimed ? "claimed" : "revealed";
}

/** The deposit a bid is offered with unless its bidder sets another: the bond, or what the bid escrows when that is more. */
export function defaultDeposit(args: { quantity: bigint | undefined; price: bigint | undefined; bond: bigint }): bigint {
  const { quantity, price, bond } = args;
  if (quantity === undefined || price === undefined) return bond;
  const escrow = escrowOf(quantity, price);
  return escrow > bond ? escrow : bond;
}

/** What a bid of the visitor's can do next. */
export type BidMove = "reveal" | "forfeit" | "claim" | "waitReveal" | "waitClear" | "none";

/**
 * The next step of a bid in the round of `openMs`: reveal in the window before the open, forfeit once it has passed
 * unrevealed, claim once the round is cleared or its clearing window has closed.
 */
export function bidMove(args: {
  status: BidStatus;
  nowMs: number;
  openMs: bigint;
  cleared: boolean;
  revealMs: bigint;
  clearMs: bigint;
}): BidMove {
  const { status, nowMs, openMs, cleared, revealMs, clearMs } = args;
  const now = BigInt(nowMs);
  if (status === "committed") {
    if (now >= openMs) return "forfeit";
    return now + revealMs >= openMs ? "reveal" : "waitReveal";
  }
  if (status === "revealed") return cleared || now >= openMs + clearMs ? "claim" : "waitClear";
  return "none";
}

/** A bid of the visitor's in a round, as the screen lists it: saved in this browser, or only listed by the auction. */
export type BidRow = {
  asset: string;
  openMs: bigint;
  quantity: bigint;
  price: bigint;
  status: BidStatus;
  index: number | undefined;
  saved: SavedBid | undefined;
};

/** The visitor's bids, newest round first: the saved ones, then the listed ones this browser holds no secret of. */
export function bidRows(
  saved: readonly (SavedBid & { status: BidStatus; index: number | undefined })[],
  listed: readonly { asset: string; openMs: bigint; index: number; quantity: bigint; price: bigint; claimed: boolean }[],
): BidRow[] {
  const known = new Set(saved.flatMap((b) => (b.index === undefined ? [] : [`${b.asset}:${b.openMs}:${b.index}`])));
  const rows: BidRow[] = [
    ...saved.map((b) => ({ asset: b.asset, openMs: b.openMs, quantity: b.quantity, price: b.price, status: b.status, index: b.index, saved: b })),
    ...listed
      .filter((b) => !known.has(`${b.asset}:${b.openMs}:${b.index}`))
      .map(
        (b): BidRow => ({
          asset: b.asset,
          openMs: b.openMs,
          quantity: b.quantity,
          price: b.price,
          status: b.claimed ? "claimed" : "revealed",
          index: b.index,
          saved: undefined,
        }),
      ),
  ];
  return rows.sort((a, b) => (a.openMs === b.openMs ? a.asset.localeCompare(b.asset) : a.openMs > b.openMs ? -1 : 1));
}

/** How the auction works, said before a bid: sealed bids, one price for all, the floor and the backstop. */
export function auctionNotes(): Note[] {
  return [
    { term: msg("auction.notes.sealed.term"), message: msg("auction.notes.sealed.body") },
    { term: msg("auction.notes.uniform.term"), message: msg("auction.notes.uniform.body") },
    { term: msg("auction.notes.floor.term"), message: msg("auction.notes.floor.body") },
    { term: msg("auction.notes.backstop.term"), message: msg("auction.notes.backstop.body") },
    { term: msg("auction.notes.reveal.term"), message: msg("auction.notes.reveal.body"), tone: "warning" },
  ];
}
