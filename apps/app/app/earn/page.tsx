// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { Earn } from "../components/Earn";
import { Header } from "../components/Header";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("pages.earn");
  return { title: t("title"), description: t("description") };
}

export default function EarnPage() {
  return (
    <>
      <Header />
      <Earn />
    </>
  );
}
