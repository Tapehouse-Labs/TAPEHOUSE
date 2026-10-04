// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { Header } from "../components/Header";
import { Trade } from "../components/Trade";

export const metadata: Metadata = {
  title: "Trade",
  description: "Short a Stock Token, buy gap cover over the weekend and hold baskets, every line of the ticket read from the chain first.",
};

export default function TradePage() {
  return (
    <>
      <Header />
      <Trade />
    </>
  );
}
