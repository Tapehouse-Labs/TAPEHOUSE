// SPDX-License-Identifier: MIT OR Apache-2.0
import { docFiles } from "@/lib/content";
import { markdownOf } from "@/lib/markdown";
import { endpoints } from "@/lib/openapi";

export const dynamicParams = false;

export function generateStaticParams() {
  return [...docFiles().map((s) => ({ slug: s.split("/") })), ...endpoints().map((e) => ({ slug: ["api", e.slug] }))];
}

export async function GET(_: Request, { params }: { params: Promise<{ slug: string[] }> }) {
  const body = markdownOf((await params).slug.join("/"));
  if (!body) return new Response("Not found", { status: 404 });
  return new Response(body, { headers: { "content-type": "text/markdown; charset=utf-8" } });
}
