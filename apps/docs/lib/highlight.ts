// SPDX-License-Identifier: MIT OR Apache-2.0
import { createHighlighter, type Highlighter, type ThemeRegistrationRaw } from "shiki";

type Palette = { fg: string; comment: string; keyword: string; string: string; number: string; type: string; punct: string; invalid: string };

function theme(name: string, type: "light" | "dark", p: Palette): ThemeRegistrationRaw {
  const rule = (scope: string[], foreground: string, fontStyle?: string) => ({
    scope,
    settings: fontStyle ? { foreground, fontStyle } : { foreground },
  });
  return {
    name,
    type,
    fg: p.fg,
    bg: "transparent",
    settings: [
      { settings: { foreground: p.fg } },
      rule(["comment", "punctuation.definition.comment"], p.comment, "italic"),
      rule(["keyword", "storage", "storage.type", "storage.modifier", "keyword.operator.new", "constant.language"], p.keyword),
      rule(["string", "string.quoted", "punctuation.definition.string", "string.regexp"], p.string),
      rule(["constant.numeric", "constant.character", "constant.other", "support.constant"], p.number),
      rule(["entity.name.type", "entity.name.class", "support.type", "support.class", "entity.other.inherited-class", "meta.type"], p.type),
      rule(["entity.name.function", "support.function", "meta.function-call"], p.fg),
      rule(["variable.other.property", "meta.object-literal.key", "support.type.property-name", "entity.name.tag", "entity.other.attribute-name"], p.type),
      rule(["punctuation", "meta.brace", "keyword.operator"], p.punct),
      rule(["invalid", "markup.deleted"], p.invalid),
      rule(["markup.inserted"], p.string),
    ],
  };
}

const light = theme("tapehouse-light", "light", {
  fg: "#1a1713",
  comment: "#7d7667",
  keyword: "#16326b",
  string: "#1c6b4a",
  number: "#8a5a00",
  type: "#3d5795",
  punct: "#5c564b",
  invalid: "#9e1409",
});

const dark = theme("tapehouse-dark", "dark", {
  fg: "#f2efe8",
  comment: "#8c8576",
  keyword: "#8fa9e8",
  string: "#5fbf92",
  number: "#e0a93c",
  type: "#6e88c8",
  punct: "#b4aea1",
  invalid: "#ff7a66",
});

export const LANGS = ["typescript", "tsx", "javascript", "json", "bash", "rust", "go", "solidity", "yaml", "python", "toml", "http", "diff", "markdown", "ini", "sql"] as const;

const ALIASES: Record<string, string> = {
  ts: "typescript",
  js: "javascript",
  sh: "bash",
  shell: "bash",
  zsh: "bash",
  console: "bash",
  sol: "solidity",
  rs: "rust",
  py: "python",
  yml: "yaml",
  md: "markdown",
  mdx: "markdown",
};

let highlighter: Promise<Highlighter> | undefined;

export function highlight(): Promise<Highlighter> {
  highlighter ??= createHighlighter({ themes: [light, dark], langs: [...LANGS] });
  return highlighter;
}

export function canonicalLang(lang: string | undefined): string {
  if (!lang) return "text";
  const l = ALIASES[lang] ?? lang;
  return (LANGS as readonly string[]).includes(l) ? l : "text";
}

export const LANG_LABEL: Record<string, string> = {
  typescript: "TypeScript",
  tsx: "TSX",
  javascript: "JavaScript",
  json: "JSON",
  bash: "Shell",
  rust: "Rust",
  go: "Go",
  solidity: "Solidity",
  yaml: "YAML",
  python: "Python",
  toml: "TOML",
  http: "HTTP",
  diff: "Diff",
  markdown: "Markdown",
  ini: "INI",
  sql: "SQL",
  text: "Text",
};
