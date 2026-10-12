// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { Header } from "../components/Header";
import { Simulator } from "../components/Simulator";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("pages.simulator");
  return { title: t("title"), description: t("description") };
}

export default function SimulatorPage() {
  return (
    <>
      <Header />
      <Simulator />
    </>
  );
}
