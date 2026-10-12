// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { Mark } from "../components/Mark";

export const metadata: Metadata = { title: "Not available", robots: { index: false, follow: false } };

export default function Restricted() {
  return (
    <main className="column">
      <span className="stem" aria-hidden="true" />
      <div className="pt-[2.5rem] sm:pt-[3.5rem]">
        <Mark className="h-[22px] w-auto text-accent" />
        <h1 className="label mt-8">Not available</h1>
        <p className="mt-6 max-w-[34ch] text-[19px] leading-snug text-body">
          Tapehouse is not offered where you are, because the law of the United States restricts dealing with it.
        </p>
        <p className="mt-4 max-w-[50ch] text-[15px] leading-snug text-muted">
          The block follows the location your connection reports. If it is wrong, connect from where you are.
        </p>
      </div>
    </main>
  );
}
