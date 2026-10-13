// SPDX-License-Identifier: MIT OR Apache-2.0
import { hasDoc, readDoc } from "./content";
import { endpoints } from "./openapi";
import { slugHref, tabs, type Tab } from "./nav";

export type Item = { href: string; title: string; slug: string; method?: string };
export type Section = { title: string; items: Item[] };

const ACRONYMS: Record<string, string> = { getRedStonePayload: "Get RedStone payload", getOpenAPI: "Get OpenAPI document", stream: "Stream" };

export function humanize(id: string): string {
  if (ACRONYMS[id]) return ACRONYMS[id];
  const words = id.replace(/([a-z0-9])([A-Z])/g, "$1 $2").toLowerCase();
  return words.charAt(0).toUpperCase() + words.slice(1);
}

/** A tab's sidebar: its pages in nav order, and for the API tab every endpoint of the specification. */
export function tree(tab: Tab): Section[] {
  const sections: Section[] = tab.groups.map((g) => ({
    title: g.title,
    items: g.pages.filter(hasDoc).map((slug) => ({ slug, href: slugHref(slug), title: readDoc(slug).title })),
  }));
  if (tab.id === "api") {
    sections.push({
      title: "Endpoints",
      items: endpoints().map((e) => ({
        slug: `api/${e.slug}`,
        href: `/api/${e.slug}`,
        title: humanize(e.id),
        method: e.stream ? "WS" : e.method,
      })),
    });
  }
  return sections.filter((s) => s.items.length > 0);
}

export function flat(tab: Tab): Item[] {
  return tree(tab).flatMap((s) => s.items);
}

export function everyTab() {
  return tabs;
}
