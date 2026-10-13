// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useState } from "react";
import { GITHUB_URL, SITE_URL } from "@/lib/site";

/** Copy the page as Markdown, ask an assistant about it, or fix it at its source. */
export function PageActions({ slug, source }: { slug: string; source?: string }) {
  const [state, setState] = useState<"idle" | "done" | "failed">("idle");
  const md = `${SITE_URL}/md/${slug}`;
  const prompt = encodeURIComponent(`Read ${md} so I can ask questions about it.`);
  return (
    <div className="mt-5 flex flex-wrap items-center gap-2 text-[13px]">
      <button
        type="button"
        className="pill"
        onClick={async () => {
          try {
            const text = await fetch(`/md/${slug}`).then((r) => r.text());
            await navigator.clipboard.writeText(text);
            setState("done");
          } catch {
            setState("failed");
          }
          setTimeout(() => setState("idle"), 2000);
        }}
      >
        {state === "done" ? "Copied" : state === "failed" ? "Could not copy" : "Copy page"}
      </button>
      <a className="pill" href={`/md/${slug}`} rel="noreferrer">
        View as Markdown
      </a>
      <a className="pill" href={`https://claude.ai/new?q=${prompt}`} rel="noreferrer">
        Open in Claude
      </a>
      <a className="pill" href={`https://chatgpt.com/?q=${prompt}`} rel="noreferrer">
        Open in ChatGPT
      </a>
      {source ? (
        <a className="pill" href={`${GITHUB_URL}/edit/main/${source}`} rel="noreferrer">
          Edit on GitHub
        </a>
      ) : null}
    </div>
  );
}
