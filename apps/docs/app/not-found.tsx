// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import { Shell } from "@/components/Shell";

export default function NotFound() {
  return (
    <Shell slug="index" wide>
      <h1 className="font-display text-[2.4rem] font-normal tracking-[-0.03em] text-strong">This page does not exist</h1>
      <p className="lede max-w-[52ch]">
        It may have moved. Search for it with <kbd className="kbd">⌘</kbd> <kbd className="kbd">K</kbd>, or start from the{" "}
        <Link href="/" className="text-accent underline">introduction</Link>.
      </p>
    </Shell>
  );
}
