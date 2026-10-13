// SPDX-License-Identifier: MIT OR Apache-2.0
"use client";

import { useState } from "react";

type Param = { name: string; in: "path" | "query"; required: boolean; placeholder: string };

/** Sends the request from the reader's browser to the testnet's index and shows the answer as it came. */
export function TryIt({ base, path, params, defaults }: { base: string; path: string; params: Param[]; defaults: Record<string, string> }) {
  const [values, setValues] = useState<Record<string, string>>(defaults);
  const [state, setState] = useState<{ status: number; ms: number; body: string; payment?: string } | { error: string } | null>(null);
  const [busy, setBusy] = useState(false);

  const send = async () => {
    let url = path;
    for (const p of params.filter((x) => x.in === "path")) url = url.replace(`{${p.name}}`, encodeURIComponent(values[p.name] ?? ""));
    const query = new URLSearchParams();
    for (const p of params.filter((x) => x.in === "query")) if (values[p.name]) query.set(p.name, values[p.name]);
    setBusy(true);
    const started = performance.now();
    try {
      const res = await fetch(`${base}${url}${query.size ? `?${query}` : ""}`);
      const text = await res.text();
      let body = text;
      try {
        body = JSON.stringify(JSON.parse(text), null, 2);
      } catch {}
      setState({
        status: res.status,
        ms: Math.round(performance.now() - started),
        body: body.length > 12_000 ? `${body.slice(0, 12_000)}\n…` : body,
        payment: res.headers.get("PAYMENT-REQUIRED") ? "This request needs a payment: see Pay per request." : undefined,
      });
    } catch (e) {
      setState({ error: e instanceof Error ? e.message : "The request failed." });
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="mt-5 rounded-[10px] border border-rule bg-panel">
      <div className="flex items-center justify-between gap-3 border-b border-rule px-4 py-2.5">
        <span className="font-display text-[15px] font-medium text-strong">Try it</span>
        <span className="font-mono text-[11.5px] text-faint">testnet · 46630</span>
      </div>
      <form
        className="space-y-3 p-4"
        onSubmit={(e) => {
          e.preventDefault();
          void send();
        }}
      >
        {params.map((p) => (
          <label key={p.name} className="block">
            <span className="font-mono text-[12.5px] text-strong">
              {p.name}
              {p.required ? <span className="text-warning"> *</span> : null}
            </span>
            <input
              value={values[p.name] ?? ""}
              onChange={(e) => setValues((v) => ({ ...v, [p.name]: e.target.value }))}
              placeholder={p.placeholder}
              spellCheck={false}
              className="mt-1 block w-full rounded-md border border-rule bg-paper px-2.5 py-1.5 font-mono text-[13px] text-strong outline-none focus:border-accent"
            />
          </label>
        ))}
        <button type="submit" disabled={busy} className="pill pill-solid disabled:opacity-60">
          {busy ? "Sending" : "Send request"}
        </button>
      </form>
      {state ? (
        <div className="border-t border-rule">
          {"error" in state ? (
            <p className="p-4 text-[13.5px] text-loss">{state.error}</p>
          ) : (
            <>
              <div className="flex items-center gap-3 px-4 py-2 text-[12.5px]">
                <span className="status-chip" data-ok={state.status < 300}>{state.status}</span>
                <span className="font-mono text-faint">{state.ms} ms</span>
                {state.payment ? <span className="text-muted">{state.payment}</span> : null}
              </div>
              <pre className="max-h-80 overflow-auto border-t border-rule px-4 py-3 font-mono text-[12.5px] leading-[1.6] text-body">{state.body}</pre>
            </>
          )}
        </div>
      ) : null}
    </div>
  );
}
