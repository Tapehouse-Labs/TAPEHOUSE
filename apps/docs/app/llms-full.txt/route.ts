// SPDX-License-Identifier: MIT OR Apache-2.0
import { hasDoc } from "@/lib/content";
import { markdownOf } from "@/lib/markdown";
import { tabs } from "@/lib/nav";
import { endpoints } from "@/lib/openapi";

export const dynamic = "force-static";

export function GET() {
  const slugs = [
    ...tabs.flatMap((t) => t.groups.flatMap((g) => g.pages)).filter(hasDoc),
    ...endpoints().map((e) => `api/${e.slug}`),
  ];
  const body = slugs.map((s) => markdownOf(s)).filter(Boolean).join("\n---\n\n");
  return new Response(body, { headers: { "content-type": "text/plain; charset=utf-8" } });
}
