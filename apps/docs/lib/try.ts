// SPDX-License-Identifier: MIT OR Apache-2.0
import type { Endpoint } from "./openapi";

/** What each endpoint is tried with: a contract and an event the testnet's index holds. */
export const DEFAULTS: Record<string, Record<string, string>> = {
  listEvents: { contract: "tapehouse.MarginAccounts", event: "Borrow", limit: "2" },
  callView: { contract: "tapehouse.SupplyVault", function: "totalAssets" },
  getGapCoverSeries: { symbol: "TSLA" },
  getRedStonePayload: { feeds: "TSLA---24_7" },
  stream: { contract: "tapehouse.MarginAccounts", event: "Borrow" },
};

export function requestUrl(base: string, e: Endpoint, values: Record<string, string>): string {
  let path = e.path;
  for (const p of e.pathParams) path = path.replace(`{${p.name}}`, encodeURIComponent(values[p.name] || `<${p.name}>`));
  const query = e.queryParams
    .filter((p) => values[p.name])
    .map((p) => `${encodeURIComponent(p.name)}=${encodeURIComponent(values[p.name]).replace(/%2C/gi, ",")}`)
    .join("&");
  return `${base}${path}${query ? `?${query}` : ""}`;
}

export function examples(url: string, e: Endpoint): { title: string; lang: string; code: string }[] {
  if (e.stream) {
    const ws = url.replace(/^http/, "ws");
    return [
      { title: "wscat", lang: "bash", code: `npx wscat --connect "${ws}"` },
      {
        title: "TypeScript",
        lang: "ts",
        code: `const socket = new WebSocket("${ws}")
socket.onmessage = (message) => {
  const update = JSON.parse(message.data)
  console.log(update.type, update)
}`,
      },
    ];
  }
  const list: { title: string; lang: string; code: string }[] = [
    { title: "curl", lang: "bash", code: `curl "${url}"` },
  ];
  if (e.id === "getOpenAPI") return list;
  list.push(
    {
      title: "TypeScript",
      lang: "ts",
      code: `const response = await fetch("${url}")
if (!response.ok) throw new Error(\`\${response.status}\`)
const body = await response.json()`,
    },
    {
      title: "Go",
      lang: "go",
      code: `response, err := http.Get("${url}")
if err != nil {
	return err
}
defer response.Body.Close()
var body map[string]any
err = json.NewDecoder(response.Body).Decode(&body)`,
    },
    {
      title: "Python",
      lang: "python",
      code: `import json
from urllib.request import urlopen

with urlopen("${url}") as response:
    body = json.load(response)`,
    },
  );
  return list;
}
