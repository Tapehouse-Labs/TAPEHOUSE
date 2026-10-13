// SPDX-License-Identifier: MIT OR Apache-2.0
import { readDoc } from "./content";
import { endpointBySlug, type Endpoint } from "./openapi";
import { INDEXER_URL, SITE_URL } from "./site";
import { DEFAULTS, requestUrl } from "./try";
import { humanize } from "./tree";

/** The page's source without its components: the Markdown a reader or an agent takes away. */
export function plain(body: string): string {
  return body
    .replace(/<Tab title="([^"]*)">/g, "\n**$1**\n")
    .replace(/<\/Tab>/g, "")
    .replace(/<\/?Tabs[^>]*>/g, "")
    .replace(/<Step title="([^"]*)">/g, "\n### $1\n")
    .replace(/<\/Step>/g, "")
    .replace(/<\/?Steps>/g, "")
    .replace(/<\/?Cards>/g, "")
    .replace(/<Card title="([^"]*)" href="([^"]*)">\s*([\s\S]*?)\s*<\/Card>/g, (_, t, h, b) => `- [${t}](${h.startsWith("/") ? `${SITE_URL}${h}` : h}): ${String(b).replace(/\s+/g, " ")}`)
    .replace(/<(Note|Warning|Tip)(?: title="([^"]*)")?>\s*([\s\S]*?)\s*<\/\1>/g, (_, k, t, b) => `> **${t ?? k}.** ${String(b).replace(/\n+/g, "\n> ")}`)
    .replace(/<Callout kind="(\w+)"(?: title="([^"]*)")?>\s*([\s\S]*?)\s*<\/Callout>/g, (_, k, t, b) => `> **${t ?? k}.** ${String(b).replace(/\n+/g, "\n> ")}`)
    .replace(/\]\((\/[^)]*)\)/g, (_, h) => `](${SITE_URL}${h})`)
    .replace(/\n{3,}/g, "\n\n")
    .trim();
}

export function docMarkdown(slug: string): string {
  const doc = readDoc(slug);
  return `# ${doc.title}\n\n> ${doc.description}\n\nSource: ${SITE_URL}${slug === "index" ? "/" : `/${slug}`}\n\n${plain(doc.body)}\n`;
}

export function endpointMarkdown(e: Endpoint): string {
  const url = requestUrl(INDEXER_URL, e, DEFAULTS[e.id] ?? {});
  const rows = (items: Endpoint["queryParams"]) =>
    items.map((p) => `- \`${p.name}\`${p.required ? " (required)" : ""}: ${(p.description ?? "").replace(/\s+/g, " ")}`).join("\n");
  const fields = e.responses
    .filter((r) => r.status === "200")
    .flatMap((r) => r.fields)
    .map((f) => `${"  ".repeat(f.depth)}- \`${f.name}\` (${f.type})${f.description ? `: ${f.description}` : ""}`)
    .join("\n");
  const errors = e.responses
    .filter((r) => r.status !== "200")
    .map((r) => `- ${r.status}: ${r.description}`)
    .join("\n");
  return [
    `# ${humanize(e.id)}`,
    `> ${e.summary}`,
    `Source: ${SITE_URL}/api/${e.slug}`,
    `\`${e.method} ${e.path}\``,
    e.description,
    e.pathParams.length ? `## Path parameters\n\n${rows(e.pathParams)}` : "",
    e.queryParams.length ? `## Query parameters\n\n${rows(e.queryParams)}` : "",
    fields ? `## Response\n\n${fields}` : "",
    `## Example\n\n\`\`\`bash\ncurl "${url}"\n\`\`\``,
    errors ? `## Errors\n\n${errors}` : "",
  ]
    .filter(Boolean)
    .join("\n\n")
    .concat("\n");
}

export function markdownOf(slug: string): string | undefined {
  if (slug.startsWith("api/")) {
    const e = endpointBySlug(slug.slice(4));
    if (e) return endpointMarkdown(e);
  }
  try {
    return docMarkdown(slug);
  } catch {
    return undefined;
  }
}
