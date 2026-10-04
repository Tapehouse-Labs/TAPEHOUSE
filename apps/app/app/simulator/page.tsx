// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { Header } from "../components/Header";
import { Simulator } from "../components/Simulator";

export const metadata: Metadata = {
  title: "Simulator",
  description: "Margin any portfolio of Stock Tokens with Tapehouse's engine, live from the chain.",
};

export default function SimulatorPage() {
  return (
    <>
      <Header />
      <Simulator />
    </>
  );
}
