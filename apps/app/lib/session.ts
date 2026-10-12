// SPDX-License-Identifier: MIT OR Apache-2.0

/** The band's states, indexed as `quote` reports them. */
const BAND_STATES = ["HALTED", "DEGRADED", "CLOSED", "OPEN"] as const;

export type BandState = (typeof BAND_STATES)[number];

/** The name of a band state; a state the app does not know reads as halted. */
export function bandState(state: number): BandState {
  return BAND_STATES[state] ?? "HALTED";
}

export type Seal = {
  state: "OPEN" | "CLOSED" | "UNKNOWN";
  hours: number | undefined;
  eta: string | undefined;
  label: string;
};

/** The session seal: the market's state and the time to the session's next boundary, in whole hours, or in minutes under an hour. */
export function seal(session: { state: number; boundaryMs: bigint }, nowMs: number): Seal {
  if (session.state !== 1 && session.state !== 2) {
    return { state: "UNKNOWN", hours: undefined, eta: undefined, label: "Market session not known, every band degraded" };
  }
  const ahead = session.boundaryMs - BigInt(nowMs);
  const minutes = ahead > 0n ? Number(ahead / 60_000n) : 0;
  const hours = Math.floor(minutes / 60);
  const [eta, span] =
    hours >= 1
      ? [`${hours}H`, `${hours} hours`]
      : minutes >= 1
        ? [`${minutes}M`, minutes === 1 ? "1 minute" : `${minutes} minutes`]
        : ["<1M", "under a minute"];
  return session.state === 2
    ? { state: "OPEN", hours, eta, label: `Market open, closes in ${span}` }
    : { state: "CLOSED", hours, eta, label: `Market closed, reopens in ${span}` };
}
