// SPDX-License-Identifier: MIT OR Apache-2.0
import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { describe, expect, test } from "vitest";
import { loadRegistry } from "./registry";

describe("loadRegistry", () => {
  test("reads the testnet's registry by default", async () => {
    const { deployments, rpcUrl } = await loadRegistry({});
    expect(deployments.chainId).toBe(46630);
    expect(deployments.tapehouse.Band).toBe("0xa70118d3324D90532E7D2854627b13CacE305641");
    expect(rpcUrl).toBeUndefined();
  });

  test("reads the chain named by TAPEHOUSE_CHAIN_ID, and its RPC URL", async () => {
    const { deployments, rpcUrl } = await loadRegistry({
      TAPEHOUSE_CHAIN_ID: "4663",
      TAPEHOUSE_RPC_URL: "https://rpc.example",
    });
    expect(deployments.chainId).toBe(4663);
    expect(rpcUrl).toBe("https://rpc.example");
  });

  test("TAPEHOUSE_REGISTRY names a registry file, such as the dev node's", async () => {
    const file = join(mkdtempSync(join(tmpdir(), "registry-")), "devnode.json");
    writeFileSync(file, JSON.stringify({ chainId: 412346, tapehouse: { Band: "0x5fbdb2315678afecb367f032d93f642f64180aa3" } }));
    const { deployments } = await loadRegistry({ TAPEHOUSE_REGISTRY: file });
    expect(deployments.tapehouse.Band).toBe("0x5FbDB2315678afecb367f032d93F642f64180aa3");
  });

  test("an address in mixed case must carry its checksum", async () => {
    const file = join(mkdtempSync(join(tmpdir(), "registry-")), "bad.json");
    writeFileSync(file, JSON.stringify({ chainId: 412346, tapehouse: { Band: "0x5FbDB2315678afecb367f032d93F642f64180aA3" } }));
    await expect(loadRegistry({ TAPEHOUSE_REGISTRY: file })).rejects.toThrow(".tapehouse.Band is not an address.");
  });
});
