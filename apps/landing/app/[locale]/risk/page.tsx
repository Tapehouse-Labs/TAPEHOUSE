// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { hasLocale } from "next-intl";
import { getTranslations, setRequestLocale } from "next-intl/server";
import { alternates, routing } from "@/i18n/routing";
import { Entry, PlainPage } from "../../components/PlainPage";

export async function generateMetadata({ params }: PageProps<"/[locale]/risk">): Promise<Metadata> {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) return {};
  const t = await getTranslations({ locale, namespace: "risk.meta" });
  return { title: t("title"), description: t("description"), alternates: alternates(locale, "/risk") };
}

export default async function Risk({ params }: PageProps<"/[locale]/risk">) {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) notFound();
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "risk" });
  return (
    <PlainPage title={t("title")} lede={t("lede")}>
      <Entry heading={t("entries.what.heading")}>
        <p>{t("entries.what.p1")}</p>
      </Entry>
      <Entry heading={t("entries.stock.heading")}>
        <p>{t("entries.stock.p1")}</p>
        <p>{t("entries.stock.p2")}</p>
      </Entry>
      <Entry heading={t("entries.borrow.heading")}>
        <p>{t("entries.borrow.p1")}</p>
      </Entry>
      <Entry heading={t("entries.fail.heading")}>
        <p>{t("entries.fail.p1")}</p>
      </Entry>
      <Entry heading={t("entries.where.heading")}>
        <p>{t("entries.where.p1")}</p>
      </Entry>
    </PlainPage>
  );
}
