// SPDX-License-Identifier: MIT OR Apache-2.0
import Link from "next/link";
import type { Endpoint, Field } from "@/lib/openapi";
import { liveExample } from "@/lib/live";
import { INDEXER_URL } from "@/lib/site";
import { DEFAULTS, examples, requestUrl } from "@/lib/try";
import { humanize } from "@/lib/tree";
import { CopyButton } from "./CopyButton";
import { fence, Md } from "./Md";
import { PageActions } from "./PageActions";
import { Tab, Tabs } from "./Tabs";
import { TryIt } from "./TryIt";

function Inline({ text }: { text: string }) {
  return (
    <>
      {text.replace(/\s+/g, " ").split(/(`[^`]+`)/).map((part, i) =>
        part.startsWith("`") ? <code key={i}>{part.slice(1, -1)}</code> : part,
      )}
    </>
  );
}

function Params({ title, items }: { title: string; items: Endpoint["queryParams"] }) {
  if (items.length === 0) return null;
  return (
    <section>
      <h2>{title}</h2>
      <div>
        {items.map((p) => (
          <div key={p.name} className="field-row">
            <div className="flex flex-wrap items-baseline gap-x-3">
              <span className="field-name">{p.name}</span>
              <span className="field-type">{schemaType(p.schema)}</span>
              {p.required ? <span className="field-flag">required</span> : null}
            </div>
            {p.description ? <div className="text-[14.5px] text-muted"><Inline text={p.description} /></div> : null}
          </div>
        ))}
      </div>
    </section>
  );
}

function schemaType(s: { type?: string | string[]; enum?: unknown[]; minimum?: number; maximum?: number; default?: unknown } | undefined): string {
  if (!s) return "string";
  const t = Array.isArray(s.type) ? s.type.join(" | ") : (s.type ?? "string");
  const bounds = s.minimum !== undefined || s.maximum !== undefined ? ` (${s.minimum ?? ""}–${s.maximum ?? ""})` : "";
  const def = s.default !== undefined ? `, default ${String(s.default)}` : "";
  return `${t}${bounds}${def}`;
}

function FieldList({ fields }: { fields: Field[] }) {
  return (
    <div>
      {fields.map((f, i) => (
        <div key={`${f.name}-${i}`} className="field-row" style={{ paddingLeft: `${f.depth * 1.25}rem` }}>
          <div className="flex flex-wrap items-baseline gap-x-3">
            <span className="field-name">{f.name}</span>
            <span className="field-type">{f.type}</span>
            {f.required ? <span className="field-flag">always present</span> : null}
          </div>
          {f.description ? <div className="text-[14.5px] text-muted"><Inline text={f.description} /></div> : null}
        </div>
      ))}
    </div>
  );
}

export async function ApiEndpoint({ e }: { e: Endpoint }) {
  const defaults = DEFAULTS[e.id] ?? {};
  const url = requestUrl(INDEXER_URL, e, defaults);
  const live = await liveExample(e);
  const ok = e.responses.find((r) => r.status === "200" || r.status === "101");
  const others = e.responses.filter((r) => r !== ok);
  const params = [...e.pathParams, ...e.queryParams].map((p) => ({
    name: p.name,
    in: p.in as "path" | "query",
    required: Boolean(p.required),
    placeholder: p.description?.replace(/\s+/g, " ").slice(0, 60) ?? "",
  }));
  const tryable = !e.stream && e.id !== "getOpenAPI";
  return (
    <div className="api-grid">
      <div className="prose-docs min-w-0">
        <p className="font-mono text-[12px] uppercase tracking-wider text-faint">Indexer API</p>
        <h1 className="mt-2 font-display text-[2.4rem] font-normal leading-[1.08] tracking-[-0.03em] text-strong">{humanize(e.id)}</h1>
        <p className="lede">{e.summary}</p>
        <div className="mt-5 flex flex-wrap items-center gap-3 rounded-lg border border-rule bg-sunk px-3 py-2">
          <span className="method">{e.stream ? "WS" : e.method}</span>
          <code className="!bg-transparent !p-0 text-[13.5px]">{e.path}</code>
          <span className="ml-auto">
            <CopyButton text={`${INDEXER_URL}${e.path}`} label="Copy URL" />
          </span>
        </div>
        <PageActions slug={`api/${e.slug}`} />
        {e.description ? <p><Inline text={e.description} /></p> : null}
        {e.sold ? (
          <p>
            The public tier allows 10 requests a second per client. Beyond it, send an <code>X-API-Key</code> or{" "}
            <Link href="/agents/x402">pay for each request in USDG</Link>. See <Link href="/api/overview">the overview</Link>.
          </p>
        ) : null}
        <Params title="Path parameters" items={e.pathParams} />
        <Params title="Query parameters" items={e.queryParams} />
        {ok ? (
          <section>
            <h2>Response</h2>
            <p>
              <span className="status-chip" data-ok="true">{ok.status}</span> {ok.description}
            </p>
            {ok.fields.length > 0 ? <FieldList fields={ok.fields} /> : null}
          </section>
        ) : null}
        <section>
          <h2>Errors</h2>
          <div className="table-wrap">
            <table>
              <tbody>
                {others.map((r) => (
                  <tr key={r.status}>
                    <td><span className="status-chip">{r.status}</span></td>
                    <td><Inline text={r.description} /></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      </div>
      <aside className="api-aside" aria-label="Request and response">
        <Tabs group="api-lang">
          {examples(url, e).map((x) => (
            <Tab key={x.title} title={x.title}>
              <Md source={fence(x.lang, x.code)} />
            </Tab>
          ))}
        </Tabs>
        {tryable ? <TryIt base={INDEXER_URL} path={e.path} params={params} defaults={defaults} /> : null}
        {live !== undefined && !e.stream ? (
          <div className="mt-5">
            <div className="mb-1 font-display text-[15px] font-medium text-strong">Example response</div>
            <Md source={fence("json", JSON.stringify(live, null, 2))} />
          </div>
        ) : null}
      </aside>
    </div>
  );
}
