// SPDX-License-Identifier: MIT OR Apache-2.0
import { NextRequest } from "next/server";
import { describe, expect, test } from "vitest";
import { config, proxy } from "../proxy";

const request = (headers: Record<string, string>) => new NextRequest("https://app.tapehouse.xyz/earn", { headers });

describe("proxy", () => {
  test("a visitor from an embargoed country is shown the page that says so, with 451", () => {
    const response = proxy(request({ "x-vercel-ip-country": "IR" }));
    expect(response.status).toBe(451);
    expect(response.headers.get("x-middleware-rewrite")).toBe("https://app.tapehouse.xyz/restricted");
  });

  test("so is one from an embargoed region of Ukraine, and not one from the rest of it", () => {
    expect(proxy(request({ "x-vercel-ip-country": "UA", "x-vercel-ip-country-region": "43" })).status).toBe(451);
    const open = proxy(request({ "x-vercel-ip-country": "UA", "x-vercel-ip-country-region": "30" }));
    expect(open.headers.get("x-middleware-rewrite")).toBeNull();
  });

  test("everyone else, and a request with no location, passes", () => {
    const open: Record<string, string>[] = [{ "x-vercel-ip-country": "CO" }, {}];
    for (const headers of open) {
      const response = proxy(request(headers));
      expect(response.status).toBe(200);
      expect(response.headers.get("x-middleware-rewrite")).toBeNull();
    }
  });

  test("the restricted page and the static files are outside its reach", () => {
    const pattern = new RegExp(`^${config.matcher[0]!}$`);
    expect(pattern.test("/earn")).toBe(true);
    expect(pattern.test("/")).toBe(true);
    expect(pattern.test("/restricted")).toBe(false);
    expect(pattern.test("/_next/static/chunk.js")).toBe(false);
    expect(pattern.test("/icon.svg")).toBe(false);
  });
});
