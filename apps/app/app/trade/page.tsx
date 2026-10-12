// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { Header } from "../components/Header";
import { Trade } from "../components/Trade";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("pages.trade");
  return { title: t("title"), description: t("description") };
}

export default function TradePage() {
  return (
    <>
      <Header />
      <Trade />
    </>
  );
}
