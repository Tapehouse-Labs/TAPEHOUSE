// SPDX-License-Identifier: MIT OR Apache-2.0
import { describe, expect, test } from "vitest";
import robots from "../app/robots";
import sitemap from "../app/sitemap";
import { SITE_URL } from "./site";

describe("search engines", () => {
  test("the sitemap lists the sections in the menu, and not the hidden trade surface", () => {
    const urls = sitemap().map((entry) => entry.url);
    expect(urls).toEqual([SITE_URL, `${SITE_URL}/earn`, `${SITE_URL}/risk`, `${SITE_URL}/auction`, `${SITE_URL}/simulator`]);
    expect(urls.some((url) => url.includes("trade"))).toBe(false);
  });

  test("robots lets every page be crawled and names the sitemap", () => {
    expect(robots()).toEqual({ rules: { userAgent: "*", allow: "/" }, sitemap: `${SITE_URL}/sitemap.xml` });
  });
});
