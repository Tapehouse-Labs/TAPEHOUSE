// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import { Header } from "./components/Header";

export default function NotFound() {
  return (
    <>
      <Header />
      <main className="column">
        <span className="stem" aria-hidden="true" />
        <div className="pt-[2.5rem] sm:pt-[3.5rem]">
          <div className="relative pt-14">
            <span className="bar" aria-hidden="true" />
            <h1 className="label">404 · This page does not exist</h1>
            <Link href="/" className="btn mt-8 inline-flex">
              Back to your account
            </Link>
          </div>
        </div>
      </main>
    </>
  );
}
