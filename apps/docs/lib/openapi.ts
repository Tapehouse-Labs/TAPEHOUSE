// SPDX-License-Identifier: MIT OR Apache-2.0
import { readFileSync } from "node:fs";
import path from "node:path";
import { parse } from "yaml";

export const SPEC_PATH = path.join(process.cwd(), "..", "..", "services", "internal", "api", "openapi.yaml");

type Schema = {
  type?: string | string[];
  description?: string;
  enum?: unknown[];
  pattern?: string;
  minimum?: number;
  maximum?: number;
  default?: unknown;
  required?: string[];
  properties?: Record<string, Schema>;
  items?: Schema;
  additionalProperties?: boolean | Schema;
  oneOf?: Schema[];
  $ref?: string;
};

type Parameter = { name: string; in: string; required?: boolean; description?: string; style?: string; schema?: Schema };
type Response = { description?: string; headers?: Record<string, { schema?: Schema }>; content?: Record<string, { schema?: Schema }>; $ref?: string };
type Operation = {
  operationId: string;
  summary: string;
  description?: string;
  parameters?: Parameter[];
  responses: Record<string, Response>;
};
type Spec = {
  info: { title: string; version: string; description: string };
  paths: Record<string, { get: Operation }>;
  components: { schemas: Record<string, Schema>; responses: Record<string, Response> };
};

export type Endpoint = {
  slug: string;
  id: string;
  path: string;
  method: "GET";
  summary: string;
  description: string;
  pathParams: Parameter[];
  queryParams: Parameter[];
  responses: { status: string; description: string; fields: Field[]; headers: string[] }[];
  example: unknown;
  stream: boolean;
  sold: boolean;
};

export type Field = { name: string; type: string; required: boolean; description: string; depth: number };

let cached: { spec: Spec; endpoints: Endpoint[] } | undefined;

function deref<T extends { $ref?: string }>(spec: Spec, node: T): T {
  if (!node.$ref) return node;
  const parts = node.$ref.replace(/^#\//, "").split("/");
  let cur: unknown = spec;
  for (const part of parts) cur = (cur as Record<string, unknown>)[part];
  return cur as T;
}

function typeOf(spec: Spec, schema: Schema): string {
  const s = deref(spec, schema);
  if (schema.$ref) {
    const name = schema.$ref.split("/").pop() ?? "";
    if (["Integer", "Address", "Bytes32", "Hex"].includes(name)) return name === "Integer" ? "integer string" : name.toLowerCase();
    if (name === "Event") return "Event";
    return name;
  }
  if (s.oneOf) return s.oneOf.map((o) => typeOf(spec, o)).join(" | ");
  if (s.enum) return s.enum.map((v) => (v === null ? "null" : JSON.stringify(v))).join(" | ");
  const t = Array.isArray(s.type) ? s.type.join(" | ") : (s.type ?? "object");
  if (t === "array" && s.items) return `${typeOf(spec, s.items)}[]`;
  return t;
}

function fields(spec: Spec, schema: Schema, depth = 0, out: Field[] = []): Field[] {
  const s = deref(spec, schema);
  if (s.type === "object" || s.properties) {
    for (const [name, prop] of Object.entries(s.properties ?? {})) {
      const resolved = deref(spec, prop);
      out.push({
        name,
        type: typeOf(spec, prop),
        required: (s.required ?? []).includes(name),
        description: resolved.description ?? "",
        depth,
      });
      if (depth < 2 && !prop.$ref) {
        if (resolved.properties) fields(spec, resolved, depth + 1, out);
        else if (resolved.type === "array" && resolved.items && !resolved.items.$ref) fields(spec, resolved.items, depth + 1, out);
      }
    }
  }
  return out;
}

const SAMPLE = {
  address: "0x7E955252E15c84f5768B83c41a71F9eba181802F",
  bytes32: "0x" + "ab".repeat(32),
  hex: "0x",
  integer: "1000000",
};

function sample(spec: Spec, schema: Schema, name = ""): unknown {
  const ref = schema.$ref?.split("/").pop();
  if (ref === "Address") return SAMPLE.address;
  if (ref === "Bytes32") return SAMPLE.bytes32;
  if (ref === "Hex") return SAMPLE.hex;
  if (ref === "Integer") return SAMPLE.integer;
  const s = deref(spec, schema);
  if (s.oneOf) return sample(spec, s.oneOf.find((o) => deref(spec, o).type !== "null") ?? s.oneOf[0], name);
  if (s.enum) return s.enum.find((v) => v !== null) ?? null;
  const t = Array.isArray(s.type) ? s.type.find((x) => x !== "null") : s.type;
  if (t === "integer") return name === "head" || name === "block" ? 134_000_000 : 1;
  if (t === "string") return s.pattern ? "1" : name;
  if (t === "boolean") return false;
  if (t === "array") return s.items ? [sample(spec, s.items)] : [];
  if (t === "object" || s.properties) {
    const out: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(s.properties ?? {})) out[k] = sample(spec, v, k);
    return out;
  }
  return null;
}

export function slugOf(operationId: string): string {
  return operationId.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
}

function load() {
  if (cached) return cached;
  const spec = parse(readFileSync(SPEC_PATH, "utf8")) as Spec;
  const endpoints: Endpoint[] = Object.entries(spec.paths).map(([route, item]) => {
    const op = item.get;
    const params = op.parameters ?? [];
    const responses = Object.entries(op.responses).map(([status, raw]) => {
      const res = deref(spec, raw);
      const schema = res.content?.["application/json"]?.schema;
      return {
        status,
        description: res.description?.replace(/\s+/g, " ").trim() ?? "",
        fields: schema ? fields(spec, schema) : [],
        headers: Object.keys(res.headers ?? {}),
      };
    });
    const ok = op.responses["200"] ? deref(spec, op.responses["200"]).content?.["application/json"]?.schema : undefined;
    return {
      slug: slugOf(op.operationId),
      id: op.operationId,
      path: route,
      method: "GET" as const,
      summary: op.summary,
      description: (op.description ?? "").replace(/\s+/g, " ").trim(),
      pathParams: params.filter((p) => p.in === "path"),
      queryParams: params.filter((p) => p.in === "query"),
      responses,
      example: ok ? sample(spec, ok) : undefined,
      stream: op.operationId === "stream",
      sold: !["stream", "getOpenAPI"].includes(op.operationId),
    };
  });
  cached = { spec, endpoints };
  return cached;
}

export function endpoints(): Endpoint[] {
  return load().endpoints;
}

export function endpointBySlug(slug: string): Endpoint | undefined {
  return load().endpoints.find((e) => e.slug === slug);
}

export function apiInfo() {
  return load().spec.info;
}
