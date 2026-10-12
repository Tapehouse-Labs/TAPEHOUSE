// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { hasLocale } from "next-intl";
import { getTranslations, setRequestLocale } from "next-intl/server";
import { alternates, routing } from "@/i18n/routing";
import { Entry, PlainPage } from "../../components/PlainPage";

export async function generateMetadata({ params }: PageProps<"/[locale]/methodology">): Promise<Metadata> {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) return {};
  const t = await getTranslations({ locale, namespace: "methodology.meta" });
  return { title: t("title"), description: t("description"), alternates: alternates(locale, "/methodology") };
}

export default async function Methodology({ params }: PageProps<"/[locale]/methodology">) {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) notFound();
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "methodology" });
  return (
    <PlainPage title={t("title")} lede={t("lede")}>
      <Entry heading={t("entries.lending.heading")}>
        <p>{t("entries.lending.p1")}</p>
      </Entry>
      <Entry heading={t("entries.silence.heading")}>
        <p>{t("entries.silence.p1")}</p>
      </Entry>
      <Entry heading={t("entries.coverage.heading")}>
        <p>{t("entries.coverage.p1")}</p>
      </Entry>
      <Entry heading={t("entries.weekendShare.heading")}>
        <p>{t("entries.weekendShare.p1")}</p>
      </Entry>
      <Entry heading={t("entries.flatRule.heading")}>
        <p>{t("entries.flatRule.p1")}</p>
      </Entry>
      <Entry heading={t("entries.model.heading")}>
        <p>{t("entries.model.p1")}</p>
      </Entry>
      <Entry heading={t("entries.band.heading")}>
        <p>{t("entries.band.p1")}</p>
        <p>{t("entries.band.p2")}</p>
        <p>{t("entries.band.p3")}</p>
      </Entry>
      <Entry heading={t("entries.liveRecord.heading")}>
        <p>{t("entries.liveRecord.p1")}</p>
        <p>{t("entries.liveRecord.p2")}</p>
      </Entry>
      <Entry heading={t("entries.stress.heading")}>
        <p>{t("entries.stress.p1")}</p>
        <p>{t("entries.stress.p2")}</p>
      </Entry>
      <Entry heading={t("entries.friday.heading")}>
        <p>{t("entries.friday.p1")}</p>
      </Entry>
    </PlainPage>
  );
}
