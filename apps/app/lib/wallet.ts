// SPDX-License-Identifier: MIT OR Apache-2.0

export type WalletState = "no-wallet" | "disconnected" | "connecting" | "wrong-network" | "connected";

/** Where the visitor's wallet stands against the chain the app runs on. */
export function walletState(input: {
  status: "connected" | "connecting" | "reconnecting" | "disconnected";
  chainId: number | undefined;
  expected: number;
  wallets: number;
}): WalletState {
  switch (input.status) {
    case "connected":
      return input.chainId === input.expected ? "connected" : "wrong-network";
    case "connecting":
    case "reconnecting":
      return "connecting";
    case "disconnected":
      return input.wallets > 0 ? "disconnected" : "no-wallet";
  }
}
