// SPDX-License-Identifier: MIT OR Apache-2.0
import GithubSlugger from "github-slugger";
import { docFiles, readDoc } from "./content";
import { slugHref, tabOf } from "./nav";
import { endpoints } from "./openapi";
import { humanize } from "./tree";

export type SearchRecord = { id: string; url: string; page: string; heading: string; area: string; text: string };

function clean(text: string): string {
  return text
    .replace(/<[^>]+>/g, " ")
    .replace(/[*_`>#|]/g, " ")
    .replace(/\[([^\]]+)\]\([^)]*\)/g, "$1")
    .replace(/\s+/g, " ")
    .trim();
}

/** One record per page and one per section of it, so a result opens where the words are. */
export function records(): SearchRecord[] {
  const out: SearchRecord[] = [];
  for (const slug of docFiles()) {
    const doc = readDoc(slug);
    const area = tabOf(slug).title;
    const url = slugHref(slug);
    const slugger = new GithubSlugger();
    const parts = doc.body.split(/^(#{2,3}) +(.+)$/m);
    out.push({ id: slug, url, page: doc.title, heading: "", area, text: clean(`${doc.description} ${parts[0]}`).slice(0, 600) });
    for (let i = 1; i < parts.length; i += 3) {
      const heading = clean(parts[i + 1]);
      out.push({
        id: `${slug}#${i}`,
        url: `${url}#${slugger.slug(heading)}`,
        page: doc.title,
        heading,
        area,
        text: clean(parts[i + 2] ?? "").slice(0, 600),
      });
    }
  }
  for (const e of endpoints()) {
    const params = [...e.pathParams, ...e.queryParams].map((p) => p.name).join(" ");
    out.push({
      id: `api/${e.slug}`,
      url: `/api/${e.slug}`,
      page: humanize(e.id),
      heading: `${e.method} ${e.path}`,
      area: "API reference",
      text: clean(`${e.summary} ${e.description} ${params}`).slice(0, 600),
    });
  }
  return out;
}
