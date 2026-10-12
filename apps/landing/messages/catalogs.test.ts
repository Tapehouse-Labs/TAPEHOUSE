// SPDX-License-Identifier: MIT OR Apache-2.0
import { IntlMessageFormat } from "intl-messageformat";
import { describe, expect, test } from "vitest";
import { LOCALES } from "@/i18n/routing";
import { catalogs } from "./index";

type Tree = { [key: string]: string | Tree };

function leaves(tree: Tree, prefix = ""): Map<string, string> {
  const out = new Map<string, string>();
  for (const [key, value] of Object.entries(tree)) {
    if (typeof value === "string") out.set(`${prefix}${key}`, value);
    else for (const [k, v] of leaves(value, `${prefix}${key}.`)) out.set(k, v);
  }
  return out;
}

type Node = { type: number; value?: string; options?: Record<string, { value: Node[] }>; children?: Node[] };

/** The names of the values a message takes, and how: a plural is not the same use as a plain argument. */
function placeholders(text: string, locale: string): string[] {
  const names: string[] = [];
  const walk = (nodes: Node[]) => {
    for (const node of nodes) {
      if (node.type === 0) continue;
      names.push(`${node.value}:${node.type}`);
      for (const option of Object.values(node.options ?? {})) walk(option.value);
      if (node.children) walk(node.children);
    }
  };
  walk(new IntlMessageFormat(text, locale).getAst() as Node[]);
  return names.sort();
}

describe("the catalogs", () => {
  const english = leaves(catalogs.en as Tree);

  test("every language says everything English says, and nothing more", () => {
    for (const locale of LOCALES) {
      const keys = [...leaves(catalogs[locale] as Tree).keys()].sort();
      expect(keys, locale).toEqual([...english.keys()].sort());
    }
  });

  test("every message is written, and parses as ICU", () => {
    for (const locale of LOCALES) {
      for (const [key, text] of leaves(catalogs[locale] as Tree)) {
        expect(text.trim().length, `${locale} ${key}`).toBeGreaterThan(0);
        expect(() => new IntlMessageFormat(text, locale), `${locale} ${key}`).not.toThrow();
      }
    }
  });

  test("every language fills a message with the values English does, no more and no fewer", () => {
    for (const [key, text] of english) {
      const expected = placeholders(text, "en");
      for (const locale of LOCALES) {
        const translated = leaves(catalogs[locale] as Tree).get(key) ?? "";
        expect(placeholders(translated, locale), `${locale} ${key}`).toEqual(expected);
      }
    }
  });
});
