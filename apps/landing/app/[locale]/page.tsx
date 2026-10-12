// SPDX-License-Identifier: MIT OR Apache-2.0
import { notFound } from "next/navigation";
import { hasLocale } from "next-intl";
import { setRequestLocale } from "next-intl/server";
import { routing } from "@/i18n/routing";
import { Capacity } from "../components/Capacity";
import { Closing } from "../components/Closing";
import { Footer } from "../components/Footer";
import { Frustrations } from "../components/Frustrations";
import { Hero } from "../components/Hero";
import { Method } from "../components/Method";
import { Record } from "../components/Record";
import { Thread } from "../components/Thread";
import { Weekend } from "../components/Weekend";

export default async function Home({ params }: PageProps<"/[locale]">) {
  const { locale } = await params;
  if (!hasLocale(routing.locales, locale)) notFound();
  setRequestLocale(locale);
  return (
    <>
      <main>
        <Hero />
        <div className="relative">
          <Thread />
          <Frustrations />
          <Weekend />
          <Capacity />
          <Method />
          <Record />
        </div>
        <Closing />
      </main>
      <Footer />
    </>
  );
}
