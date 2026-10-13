// SPDX-License-Identifier: MIT OR Apache-2.0
import { records } from "@/lib/search";

export const dynamic = "force-static";

export function GET() {
  return Response.json(records());
}
