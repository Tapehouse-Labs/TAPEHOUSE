// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Msg, msg } from "./messages";

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
  label: Msg;
};

/** The session seal: the market's state and the time to the session's next boundary, in whole hours or, under an hour, minutes. A boundary of zero is not known, and has no countdown. */
export function seal(session: { state: number; boundaryMs: bigint }, nowMs: number): Seal {
  if (session.state !== 1 && session.state !== 2) {
    return { state: "UNKNOWN", hours: undefined, eta: undefined, label: msg("account.seal.unknown") };
  }
  const open = session.state === 2;
  if (session.boundaryMs === 0n) {
    return { state: open ? "OPEN" : "CLOSED", hours: undefined, eta: undefined, label: msg(open ? "account.seal.open" : "account.seal.closed") };
  }
  const ahead = session.boundaryMs - BigInt(nowMs);
  const minutes = ahead > 0n ? Number(ahead / 60_000n) : 0;
  const hours = Math.floor(minutes / 60);
  const [eta, label] =
    hours >= 1
      ? [`${hours}H`, msg(open ? "account.seal.closesInHours" : "account.seal.reopensInHours", { hours })]
      : minutes >= 1
        ? [`${minutes}M`, msg(open ? "account.seal.closesInMinutes" : "account.seal.reopensInMinutes", { minutes })]
        : ["<1M", msg(open ? "account.seal.closesSoon" : "account.seal.reopensSoon")];
  return { state: open ? "OPEN" : "CLOSED", hours, eta, label };
}
