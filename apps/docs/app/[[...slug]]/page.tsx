// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { ApiEndpoint } from "@/components/ApiEndpoint";
import { components } from "@/components/mdx";
import { PageActions } from "@/components/PageActions";
import { PrevNext } from "@/components/PrevNext";
import { Shell } from "@/components/Shell";
import { Toc } from "@/components/Toc";
import { docFiles, hasDoc, readDoc } from "@/lib/content";
import { compile } from "@/lib/mdx";
import { hrefSlug, slugHref, tabOf } from "@/lib/nav";
import { endpointBySlug, endpoints } from "@/lib/openapi";
import { SITE_URL } from "@/lib/site";
import { flat, humanize } from "@/lib/tree";

type Props = { params: Promise<{ slug?: string[] }> };

export const dynamicParams = false;

export function generateStaticParams() {
  const pages = docFiles().map((slug) => ({ slug: slug === "index" ? [] : slug.split("/") }));
  const api = endpoints().map((e) => ({ slug: ["api", e.slug] }));
  return [...pages, ...api];
}

function endpointOf(slug: string) {
  return slug.startsWith("api/") ? endpointBySlug(slug.slice(4)) : undefined;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const slug = hrefSlug((await params).slug);
  const e = endpointOf(slug);
  const title = e ? humanize(e.id) : hasDoc(slug) ? readDoc(slug).title : undefined;
  if (!title) return {};
  const description = e ? e.summary : readDoc(slug).description;
  const url = `${SITE_URL}${slugHref(slug)}`;
  return {
    title: slug === "index" ? { absolute: "Tapehouse documentation" } : title,
    description,
    alternates: { canonical: url, types: { "text/markdown": `${SITE_URL}/md/${slug}` } },
    openGraph: { type: "article", url, title, description },
  };
}

export default async function Page({ params }: Props) {
  const slug = hrefSlug((await params).slug);
  const endpoint = endpointOf(slug);
  if (endpoint) {
    return (
      <Shell slug={slug} wide>
        <ApiEndpoint e={endpoint} />
      </Shell>
    );
  }
  if (!hasDoc(slug)) notFound();
  const doc = readDoc(slug);
  const { Content, headings } = await compile(doc.body);
  const tab = tabOf(slug);
  const items = flat(tab);
  const at = items.findIndex((i) => i.slug === slug);
  const group = tab.groups.find((g) => g.pages.includes(slug))?.title;
  return (
    <Shell slug={slug} tab={tab} aside={<Toc headings={headings} />}>
      <article>
        <p className="font-mono text-[12px] uppercase tracking-wider text-faint">
          {tab.title}
          {group && group !== tab.title ? ` · ${group}` : ""}
        </p>
        <h1 className="mt-2 font-display text-[2.4rem] font-normal leading-[1.08] tracking-[-0.03em] text-strong">{doc.title}</h1>
        <p className="lede">{doc.description}</p>
        <PageActions slug={slug} source={`apps/docs/content/${slug}.mdx`} />
        <div className="prose-docs mt-8">
          <Content components={components} />
        </div>
        <PrevNext prev={at > 0 ? items[at - 1] : undefined} next={at >= 0 ? items[at + 1] : undefined} />
      </article>
    </Shell>
  );
}
