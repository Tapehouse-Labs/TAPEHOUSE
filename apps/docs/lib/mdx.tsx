// SPDX-License-Identifier: MIT OR Apache-2.0
import { evaluate } from "@mdx-js/mdx";
import { Fragment, jsx, jsxs } from "react/jsx-runtime";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import { visit } from "unist-util-visit";
import type { Element, Root, RootContent } from "hast";
import type { MDXContent } from "mdx/types";
import { canonicalLang, highlight } from "./highlight";

export type Heading = { id: string; text: string; depth: 2 | 3 };

function textOf(node: RootContent | Element): string {
  if (node.type === "text") return node.value;
  if ("children" in node) return node.children.map((c) => textOf(c as RootContent)).join("");
  return "";
}

function rehypeHeadings(headings: Heading[]) {
  return () => (tree: Root) => {
    visit(tree, "element", (node: Element) => {
      if ((node.tagName === "h2" || node.tagName === "h3") && typeof node.properties.id === "string") {
        headings.push({ id: node.properties.id, text: textOf(node), depth: node.tagName === "h2" ? 2 : 3 });
      }
    });
  };
}

function rehypeCode() {
  return async (tree: Root) => {
    const hl = await highlight();
    const jobs: Promise<void>[] = [];
    visit(tree, "element", (node: Element, index, parent) => {
      if (node.tagName !== "pre" || !parent || index === undefined) return;
      const code = node.children.find((c): c is Element => c.type === "element" && c.tagName === "code");
      if (!code) return;
      const classes = (code.properties.className as string[] | undefined) ?? [];
      const lang = canonicalLang(classes.find((c) => c.startsWith("language-"))?.slice(9));
      const meta = String((code.data as { meta?: string } | undefined)?.meta ?? "");
      const title = /title="([^"]*)"/.exec(meta)?.[1];
      const raw = textOf(code).replace(/\n$/, "");
      jobs.push(
        Promise.resolve().then(() => {
          const hast = hl.codeToHast(raw, {
            lang,
            themes: { light: "tapehouse-light", dark: "tapehouse-dark" },
            defaultColor: false,
          });
          const pre = hast.children[0] as Element;
          pre.properties = { ...pre.properties, "data-lang": lang, "data-title": title, "data-raw": raw };
          parent.children[index] = pre;
        }),
      );
    });
    await Promise.all(jobs);
  };
}

export async function compile(source: string): Promise<{ Content: MDXContent; headings: Heading[] }> {
  const headings: Heading[] = [];
  const mod = await evaluate(source, {
    Fragment,
    jsx,
    jsxs,
    remarkPlugins: [remarkGfm],
    rehypePlugins: [rehypeSlug, rehypeHeadings(headings), rehypeCode],
  });
  return { Content: mod.default, headings };
}
