// SPDX-License-Identifier: MIT OR Apache-2.0
import type { MetadataRoute } from "next";
import { SITE_URL } from "../lib/site";

export default function sitemap(): MetadataRoute.Sitemap {
  return [
    { url: SITE_URL, changeFrequency: "daily", priority: 1 },
    { url: `${SITE_URL}/earn`, changeFrequency: "daily", priority: 0.8 },
    { url: `${SITE_URL}/risk`, changeFrequency: "daily", priority: 0.8 },
    { url: `${SITE_URL}/simulator`, changeFrequency: "weekly", priority: 0.6 },
  ];
}
