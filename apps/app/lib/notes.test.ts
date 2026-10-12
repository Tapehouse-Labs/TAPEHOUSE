// SPDX-License-Identifier: MIT OR Apache-2.0
import { CROSS } from "@tapehouse/sdk";
import { describe, expect, test } from "vitest";
import type { AccountState, Position } from "./accountRead";
import { formatters } from "./format";
import { borrowNotes, depositNotes, issuerNotes, type Note, positionNotes, weekendAhead } from "./notes";
import { say } from "./say";

const f = formatters("en");
const words = (notes: Note[]) => notes.map((n) => say(n.message));

const D = 10n ** 18n;
const NVDA = "0x525c2aBA45F66987217323E8a05EA400C65D06DC";
const account = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
const owner = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC";
const accounts = "0xE85035F1145aC49333d105632a0d254E479a75bE";

function state(over: Partial<AccountState> = {}): AccountState {
  return {
    blockNumber: 1n,
    lendable: 1_000_000_000n,
    markets: [
      {
        asset: "NVDA",
        token: NVDA,
        low: 17_120_000_000n,
        state: 3,
        lending: true,
        holding: { units: 61n * D, scale: D, quantity: 61n * D, cap: 100n * D, balance: 61n * D },
        issuer: { paused: false, blocked: { [accounts]: false, [account]: false, [owner]: false } },
      },
    ],
    positions: [],
    limits: {
      debt: 412_000_000_000n,
      debtCap: 1_000_000_000_000n,
      weekendDebtCap: 500_000_000_000n,
      weekendLeverage: 50_000,
      borrowingPaused: false,
      premiumRate: 500,
    },
    session: { state: 2, nyse: 1, nyseNext: 1, changeMs: 0n, boundaryMs: 0n },
    ...over,
  };
}

const position: Position = {
  id: CROSS,
  name: "Cross",
  equity: 10_000n * D,
  requirement: 4_000n * D,
  missing: 0,
  regime: 2,
  leverage: 60_000n,
  gross: 60_000n * D,
  debt: 2_000_000_000n,
  premium: 0n,
  usdg: 0n,
  weth: 0n,
  stocks: [
    { asset: "NVDA", token: NVDA, held: 12n * D, lent: 0n, sellable: 12n * D, claim: 0n, recalls: [], inBaskets: 0n, liquidationPrice: 9_610_000_000n },
  ],
  open: 3_200n * D,
  closed: 5_000n * D,
};

describe("weekendAhead", () => {
  test("a closure lies ahead from the last trading day's open, and while the session is unknown", () => {
    expect(weekendAhead(state(), position)).toBe(false);
    expect(weekendAhead(state({ session: { ...state().session, boundaryMs: 1n } }), position)).toBe(true);
    expect(weekendAhead(state(), { ...position, regime: 0 })).toBe(true);
    expect(weekendAhead(state({ session: { ...state().session, boundaryMs: 1n } }), { ...position, regime: 1 })).toBe(false);
  });
});

describe("positionNotes", () => {
  test("before a close, says what keeps the position within 5×", () => {
    const notes = positionNotes(f, state({ session: { ...state().session, boundaryMs: 1n } }), position);
    expect(words(notes)).toContain("To stay within 5× at the close, add $2,000.00 of equity or reduce the exposure by $10,000.00.");
    expect(notes.find((n) => n.message.key === "account.notes.weekendPlan")?.tone).toBe("warning");
  });

  test("names a halted asset, a missing pool and an unknown session", () => {
    const s = state();
    s.markets[0]!.state = 0;
    const texts = words(positionNotes(f, s, { ...position, regime: 0, missing: 1 }));
    expect(texts.slice(0, 3)).toEqual([
      "The band cannot tell whether the market is open, so the requirement counts as closed and borrowing stops at the price above.",
      "NVDA is halted and counts for nothing.",
      "NVDA's pool is absent or unread, so its liquidation is charged at the governed depth.",
    ]);
    expect(texts[3]).toMatch(/^From the close before a weekend or holiday to the regular open after it/);
  });

  test("on an ordinary open day the weekend rules wait", () => {
    expect(words(positionNotes(f, state(), position)).some((t) => t.startsWith("From the close"))).toBe(false);
  });
});

describe("borrowNotes", () => {
  test("states the premium over a 48-hour closure and the caps", () => {
    expect(words(borrowNotes(f, state(), position, 2_000_000_000n))).toEqual([
      "While the market is closed, 5% a year is added on top of interest: this loan would owe 1.10 USDG over a 48-hour closure. A repayment pays the debt first, then the premium.",
      "All loans together owe 412,000.00 of their 1,000,000.00 USDG cap.",
    ]);
  });

  test("before a closure names the weekend cap, and the guardian's pause first", () => {
    const s = state({ session: { ...state().session, boundaryMs: 1n }, limits: { ...state().limits, borrowingPaused: true } });
    const texts = words(borrowNotes(f, s, position, 0n));
    expect(texts[0]).toMatch(/^The guardian has paused new loans/);
    expect(texts.at(-1)).toBe(
      "All loans together owe 412,000.00 of their 1,000,000.00 USDG cap; until the reopen, new loans stop at 500,000.00.",
    );
  });
});

describe("depositNotes", () => {
  test("a Stock Token deposit names its cap, a pending write-down and the issuer", () => {
    const s = state();
    s.markets[0]!.holding = { units: 61n * D, scale: D, quantity: 61n * D, cap: 100n * D, balance: 60n * D };
    s.markets[0]!.issuer = { paused: true, blocked: { [accounts]: false, [account]: false, [owner]: true } };
    expect(words(depositNotes(f, s, NVDA, account, owner, true))).toEqual([
      "Your wallet signs a permit for free and Tapehouse pays the gas.",
      "NVDA deposits: 61 of their 100 cap.",
      "A write-down by the issuer is pending: the accounts hold 60 of the 61 NVDA they count.",
      "The issuer has paused NVDA's transfers.",
      "The issuer has blocked your address from moving NVDA.",
    ]);
  });

  test("ether is wrapped by the wallet first", () => {
    expect(words(depositNotes(f, state(), "ETH", account, owner, true))[0]).toBe(
      "Your wallet wraps the ether into WETH for your account and pays that gas; the deposit follows.",
    );
  });
});

describe("issuerNotes", () => {
  test("a block on the accounts is told apart from one on the user", () => {
    const notes = issuerNotes("NVDA", { paused: false, blocked: { [accounts]: true } }, account, owner);
    expect(words(notes)).toEqual(["The issuer has blocked the accounts from moving NVDA."]);
    expect(notes[0]?.tone).toBe("warning");
  });
});

describe("in Spanish and Portuguese", () => {
  test("a note says the same with the language's own words and separators", () => {
    const es = formatters("es");
    const pt = formatters("pt");
    const notes = positionNotes(es, state({ session: { ...state().session, boundaryMs: 1n } }), position);
    expect(say(notes.find((n) => n.message.key === "account.notes.weekendPlan")!.message, "es")).toBe(
      "Para mantenerte dentro de 5× al cierre, agrega $2000,00 de capital o reduce la exposición en $10.000,00.",
    );
    const loan = borrowNotes(pt, state(), position, 2_000_000_000n);
    expect(say(loan[0]!.message, "pt")).toBe(
      "Enquanto o mercado está fechado, soma-se 5% ao ano aos juros: este empréstimo passaria a dever 1,10 USDG em um fechamento de 48 horas. Um pagamento quita primeiro a dívida e depois o prêmio.",
    );
  });
});
