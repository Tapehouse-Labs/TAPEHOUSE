// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import { docFiles, hasDoc, navSlugs, readDoc } from "./content";
import { endpoints } from "./openapi";
import { plain } from "./markdown";
import { records } from "./search";

const slugs = docFiles();
const valid = new Set([...slugs.map((s) => (s === "index" ? "/" : `/${s}`)), ...endpoints().map((e) => `/api/${e.slug}`)]);

function prose(body: string): string {
  return body.replace(/```[\s\S]*?```/g, "");
}

describe("content", () => {
  test("every page in the navigation is a file, and every file is in the navigation", () => {
    expect(navSlugs().filter((s) => !hasDoc(s))).toEqual([]);
    expect(slugs.filter((s) => !navSlugs().includes(s))).toEqual([]);
  });

  test("every page has a title and a description", () => {
    for (const slug of slugs) {
      const doc = readDoc(slug);
      expect(doc.title.length, slug).toBeGreaterThan(0);
      expect(doc.description.endsWith("."), `${slug} description is a sentence`).toBe(true);
    }
  });

  test("every internal link names a page", () => {
    const broken: string[] = [];
    for (const slug of slugs) {
      for (const m of prose(readDoc(slug).body).matchAll(/\]\((\/[^)#\s]*)(#[^)]*)?\)|href="(\/[^"#]*)/g)) {
        const target = (m[1] ?? m[3]).replace(/\.md$/, "");
        if (target === "/llms.txt" || target === "/llms-full.txt" || target === "/quickstart.md") continue;
        if (!valid.has(target)) broken.push(`${slug} -> ${target}`);
      }
    }
    expect(broken).toEqual([]);
  });

  test("prose has no em or en dashes", () => {
    const found = slugs.filter((s) => /[—–]/.test(prose(readDoc(s).body)));
    expect(found).toEqual([]);
  });

  test("prose avoids the words the brand bans", () => {
    const banned = /institutional-grade|seamless|revolutionary|the future of|cutting-edge/i;
    expect(slugs.filter((s) => banned.test(readDoc(s).body))).toEqual([]);
  });

  test("a page's Markdown has no components left in it", () => {
    for (const slug of slugs) expect(plain(readDoc(slug).body), slug).not.toMatch(/<\/?(Tabs?|Steps?|Cards?|Note|Warning|Tip|Callout)\b/);
  });

  test("the search index has a record for every page", () => {
    const ids = new Set(records().map((r) => r.id));
    for (const slug of slugs) expect(ids.has(slug), slug).toBe(true);
  });
});
