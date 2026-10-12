// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { Mark } from "../components/Mark";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("compliance.restricted");
  return { title: t("title"), robots: { index: false, follow: false } };
}

export default async function Restricted() {
  const t = await getTranslations("compliance.restricted");
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <Mark className="h-[22px] w-auto text-accent" />
        <h1 className="label mt-8">{t("title")}</h1>
        <p className="mt-6 max-w-[34ch] text-[19px] leading-snug text-body">{t("body")}</p>
        <p className="mt-4 max-w-[50ch] text-[15px] leading-snug text-muted">{t("hint")}</p>
      </div>
    </main>
  );
}
