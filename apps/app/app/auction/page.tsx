// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { Auction } from "../components/Auction";
import { Header } from "../components/Header";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("pages.auction");
  return { title: t("title"), description: t("description") };
}

export default function AuctionPage() {
  return (
    <>
      <Header />
      <Auction />
    </>
  );
}
