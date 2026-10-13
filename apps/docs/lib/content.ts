// SPDX-License-Identifier: MIT OR Apache-2.0
import { readdirSync, readFileSync, existsSync, statSync } from "node:fs";
import path from "node:path";
import matter from "gray-matter";
import { tabs } from "./nav";

const ROOT = path.join(process.cwd(), "content");

export type Frontmatter = { title: string; description: string };
export type Doc = Frontmatter & { slug: string; body: string };

function fileOf(slug: string): string {
  return path.join(ROOT, `${slug}.mdx`);
}

export function hasDoc(slug: string): boolean {
  return existsSync(fileOf(slug));
}

export function readDoc(slug: string): Doc {
  const raw = readFileSync(fileOf(slug), "utf8");
  const { data, content } = matter(raw);
  if (typeof data.title !== "string" || typeof data.description !== "string") {
    throw new Error(`content/${slug}.mdx needs a title and a description in its front matter`);
  }
  return { slug, title: data.title, description: data.description, body: content };
}

export function docFiles(dir = ROOT, prefix = ""): string[] {
  const slugs: string[] = [];
  for (const name of readdirSync(dir)) {
    const full = path.join(dir, name);
    if (statSync(full).isDirectory()) slugs.push(...docFiles(full, `${prefix}${name}/`));
    else if (name.endsWith(".mdx")) slugs.push(`${prefix}${name.slice(0, -4)}`);
  }
  return slugs;
}

export function navSlugs(): string[] {
  return tabs.flatMap((t) => t.groups.flatMap((g) => g.pages));
}
