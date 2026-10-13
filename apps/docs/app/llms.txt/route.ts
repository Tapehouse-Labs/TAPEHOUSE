// SPDX-License-Identifier: MIT OR Apache-2.0
import { hasDoc, readDoc } from "@/lib/content";
import { tabs } from "@/lib/nav";
import { endpoints } from "@/lib/openapi";
import { SITE_URL } from "@/lib/site";
import { humanize } from "@/lib/tree";

export const dynamic = "force-static";

export function GET() {
  const lines = [
    "# Tapehouse",
    "",
    "> Portfolio margin for Stock Tokens on Robinhood Chain: a price band that keeps updating while the market is closed, margin accounts that treat a portfolio as one position, and SDKs, an API and an MCP server to read and build on them. Testnet only, not audited.",
    "",
    "Every page is also served as Markdown at its own address with `.md` appended, and the whole documentation as one file at /llms-full.txt.",
  ];
  for (const tab of tabs) {
    for (const group of tab.groups) {
      const pages = group.pages.filter(hasDoc);
      if (pages.length === 0) continue;
      lines.push("", `## ${group.title}`, "");
      for (const slug of pages) {
        const d = readDoc(slug);
        lines.push(`- [${d.title}](${SITE_URL}/md/${slug}): ${d.description}`);
      }
    }
  }
  lines.push("", "## Indexer API endpoints", "");
  for (const e of endpoints()) lines.push(`- [${humanize(e.id)}](${SITE_URL}/md/api/${e.slug}): ${e.summary}`);
  return new Response(`${lines.join("\n")}\n`, { headers: { "content-type": "text/plain; charset=utf-8" } });
}
