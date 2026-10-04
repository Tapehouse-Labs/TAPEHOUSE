// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { Earn } from "../components/Earn";
import { Header } from "../components/Header";

export const metadata: Metadata = {
  title: "Earn",
  description: "Supply USDG, back the margin accounts, lend Stock Tokens or write gap cover, with every term read from the chain first.",
};

export default function EarnPage() {
  return (
    <>
      <Header />
      <Earn />
    </>
  );
}
