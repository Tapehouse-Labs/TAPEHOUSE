// SPDX-License-Identifier: MIT OR Apache-2.0
import createMiddleware from "next-intl/middleware";
import { routing } from "./i18n/routing";

export default createMiddleware(routing);

export const config = {
  matcher: ["/((?!_next/static|_next/image|.*\\..*).*)"],
};
