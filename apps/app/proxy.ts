// SPDX-License-Identifier: MIT OR Apache-2.0
import { type NextRequest, NextResponse } from "next/server";
import { restrictedLocation } from "./lib/compliance";

/**
 * Sends a visitor from an embargoed country or region to the page that says the app is not offered there, as Vercel's
 * geolocation reports them. Without that header, as in development, nothing is restricted.
 */
export function proxy(request: NextRequest) {
  const country = request.headers.get("x-vercel-ip-country");
  const region = request.headers.get("x-vercel-ip-country-region");
  if (!restrictedLocation(country, region)) return NextResponse.next();
  return NextResponse.rewrite(new URL("/restricted", request.url), { status: 451 });
}

export const config = {
  matcher: ["/((?!restricted|_next/static|_next/image|.*\\.(?:png|svg|txt|xml)$).*)"],
};
