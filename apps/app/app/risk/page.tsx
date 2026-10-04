// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { Header } from "../components/Header";
import { Risk } from "../components/Risk";

export const metadata: Metadata = {
  title: "Risk",
  description: "Every band, the gap backstop's exposure and the latest reopening rounds, read from the chain.",
};

export default function RiskPage() {
  return (
    <>
      <Header />
      <Risk />
    </>
  );
}
