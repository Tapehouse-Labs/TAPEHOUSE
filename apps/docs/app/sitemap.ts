// SPDX-License-Identifier: MIT OR Apache-2.0
import type { MetadataRoute } from "next";
import { docFiles } from "@/lib/content";
import { slugHref } from "@/lib/nav";
import { endpoints } from "@/lib/openapi";
import { SITE_URL } from "@/lib/site";

export default function sitemap(): MetadataRoute.Sitemap {
  return [...docFiles().map(slugHref), ...endpoints().map((e) => `/api/${e.slug}`)].map((path) => ({ url: `${SITE_URL}${path}` }));
}
