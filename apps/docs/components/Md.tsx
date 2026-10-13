// SPDX-License-Identifier: MIT OR Apache-2.0
import { compile } from "@/lib/mdx";
import { components } from "./mdx";

/** Renders a string of Markdown through the same pipeline as a page. */
export async function Md({ source }: { source: string }) {
  const { Content } = await compile(source);
  return <Content components={components} />;
}

export function fence(lang: string, code: string): string {
  return `\`\`\`${lang}\n${code}\n\`\`\``;
}
