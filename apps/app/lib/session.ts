// SPDX-License-Identifier: MIT OR Apache-2.0

/** The band's states, indexed as `quote` reports them. */
const BAND_STATES = ["HALTED", "DEGRADED", "CLOSED", "OPEN"] as const;

export type BandState = (typeof BAND_STATES)[number];

/** The name of a band state; a state the app does not know reads as halted. */
export function bandState(state: number): BandState {
  return BAND_STATES[state] ?? "HALTED";
}

export type Seal = { state: "OPEN" | "CLOSED" | "UNKNOWN"; hours: number | undefined; label: string };

/** The session seal: the market's state and the whole hours to the session's next boundary. */
export function seal(session: { state: number; boundaryMs: bigint }, nowMs: number): Seal {
  if (session.state !== 1 && session.state !== 2) {
    return { state: "UNKNOWN", hours: undefined, label: "Market session not known, every band degraded" };
  }
  const ahead = session.boundaryMs - BigInt(nowMs);
  const hours = ahead > 0n ? Number(ahead / 3_600_000n) : 0;
  return session.state === 2
    ? { state: "OPEN", hours, label: `Market open, closes in ${hours} hours` }
    : { state: "CLOSED", hours, label: `Market closed, reopens in ${hours} hours` };
}
