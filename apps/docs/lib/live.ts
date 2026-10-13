// SPDX-License-Identifier: MIT OR Apache-2.0
import { INDEXER_URL } from "./site";
import type { Endpoint } from "./openapi";
import { DEFAULTS, requestUrl } from "./try";

function trim(value: unknown, depth = 0): unknown {
  if (Array.isArray(value)) return value.slice(0, 2).map((v) => trim(v, depth + 1));
  if (value && typeof value === "object") {
    return Object.fromEntries(Object.entries(value).map(([k, v]) => [k, trim(v, depth + 1)]));
  }
  return value;
}

/** A response of the testnet's index, read at build time and again each hour; the specification's sample where it has none. */
export async function liveExample(e: Endpoint): Promise<unknown> {
  if (e.stream || e.id === "getOpenAPI") return e.example;
  if (e.pathParams.some((p) => !DEFAULTS[e.id]?.[p.name])) return e.example;
  try {
    const res = await fetch(requestUrl(INDEXER_URL, e, DEFAULTS[e.id] ?? {}), {
      next: { revalidate: 3600 },
      signal: AbortSignal.timeout(8000),
    });
    if (!res.ok) return e.example;
    const body = (await res.json()) as unknown;
    const text = JSON.stringify(body);
    if (text.length > 60_000) return e.example;
    return trim(body);
  } catch {
    return e.example;
  }
}
