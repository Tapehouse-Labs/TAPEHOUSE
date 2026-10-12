// SPDX-License-Identifier: MIT OR Apache-2.0
import { getTranslations } from "next-intl/server";
import Link from "next/link";
import { Header } from "./components/Header";

export default async function NotFound() {
  const t = await getTranslations("pages.notFound");
  return (
    <>
      <Header />
      <main className="column">
        <span className="stem" aria-hidden="true" />
        <div className="pt-[2.5rem] sm:pt-[3.5rem]">
          <div className="relative pt-14">
            <span className="bar" aria-hidden="true" />
            <h1 className="label">{t("title")}</h1>
            <Link href="/" className="btn mt-8 inline-flex">
              {t("back")}
            </Link>
          </div>
        </div>
      </main>
    </>
  );
}
