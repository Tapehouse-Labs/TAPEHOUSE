// SPDX-License-Identifier: MIT OR Apache-2.0
import { type Chain, defineChain } from "viem";
import { robinhood, robinhoodTestnet } from "viem/chains";

/** The Nitro dev node the repository's `make devnode` runs. */
const devnode = defineChain({
  id: 412346,
  name: "Nitro dev node",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: ["http://127.0.0.1:8547"] } },
  testnet: true,
});

const CHAINS: readonly Chain[] = [robinhood, robinhoodTestnet, devnode];

/** The chain `id`, reached through `rpcUrl` when one is given. */
export function chainOf(id: number, rpcUrl?: string): Chain {
  const chain = CHAINS.find((c) => c.id === id);
  if (!chain) throw new Error(`Tapehouse runs on no chain ${id}.`);
  return rpcUrl ? { ...chain, rpcUrls: { default: { http: [rpcUrl] } } } : chain;
}
