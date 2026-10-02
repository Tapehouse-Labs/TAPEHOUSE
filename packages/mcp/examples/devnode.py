# SPDX-License-Identifier: MIT OR Apache-2.0
# /// script
# requires-python = ">=3.10"
# dependencies = ["mcp==2.2.0"]
# ///
"""Drives the built server over stdio with the official Python MCP client, on the 2025-11-25 protocol.

Usage: uv run --script devnode.py RPC_URL DEPLOYMENTS_JSON SERVER_JS

Reads SPY's band, both Morpho oracles and the basket PAIR on a dev node, refuses a slippage above the server's cap, and
prepares an authorization for the registry's margin accounts, signing nothing.
"""

import asyncio
import json
import os
import sys

from mcp import Client, StdioServerParameters

OPERATOR = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"


def check(condition: bool, message: str) -> None:
    if not condition:
        raise SystemExit(f"FAIL: {message}")


async def main(rpc: str, registry: str, server: str) -> None:
    params = StdioServerParameters(
        command="node",
        args=[server],
        env={"PATH": os.environ["PATH"], "TAPEHOUSE_RPC_URL": rpc, "TAPEHOUSE_DEPLOYMENTS": registry},
    )
    async with Client(params, mode="legacy") as client:
        check(client.protocol_version == "2025-11-25", f"the client negotiated {client.protocol_version}")

        async def use(name: str, arguments: dict) -> dict:
            result = await client.call_tool(name, arguments)
            check(not result.is_error, f"{name} failed: {result.content}")
            check(json.loads(result.content[0].text) == result.structured_content, f"{name}'s text is not its data")
            return result.structured_content

        tools = (await client.list_tools()).tools
        check(len(tools) == 35, f"{len(tools)} tools")
        check(all(tool.annotations.read_only_hint and tool.output_schema for tool in tools), "a tool is not read-only")

        spy = await use("band_quote", {"asset": "SPY"})
        check(spy["state"] != "halted" and int(spy["low"]) < int(spy["mid"]) < int(spy["high"]), f"SPY's band: {spy}")
        priced = await use("morpho_oracle", {"asset": "SPY"})
        check(priced["blockNumber"] == spy["blockNumber"], "the two reads are at different blocks")
        check(int(priced["price"]) == int(spy["low"]) * int(priced["scaleFactor"]), "SPY's oracle is not its low edge")
        nvda = await use("morpho_oracle", {"asset": "NVDA"})
        check(nvda["price"] is None and nvda["noPrice"] in ("stale", "halted"), f"NVDA's oracle: {nvda}")

        pair = await use("baskets_components", {"basket": "PAIR"})
        check([c["asset"] for c in pair["components"]] == ["NVDA", "SPY"], f"PAIR: {pair}")

        refused = await client.call_tool(
            "shorts_sell", {"account": OPERATOR, "asset": "SPY", "amount": "1", "slippageBps": 10_000}
        )
        check(refused.is_error and "slippageBps" in refused.content[0].text, "a slippage above the cap was taken")

        registry_entries = await use("registry", {})
        authorization = await use("accounts_set_authorization", {"operator": OPERATOR, "allowed": True})
        [call] = authorization["calls"]
        check(
            call["to"] == registry_entries["tapehouse"]["MarginAccounts"],
            "the authorization is not for the margin accounts",
        )
        check(call["functionName"] == "setAuthorization" and call["args"]["allowed"] is True, f"{call}")
        check(authorization["grants"].startswith(f"{OPERATOR} can withdraw every collateral token"), "no grants stated")

    print(
        f"python client over stdio, protocol 2025-11-25: SPY band {spy['state']} at {spy['mid']}, "
        f"its Morpho oracle {priced['price']}; NVDA's oracle no price ({nvda['noPrice']}); "
        f"basket PAIR of NVDA and SPY at {pair['address']}; "
        f"setAuthorization prepared for {call['to']}"
    )
    print("PASS")


if __name__ == "__main__":
    asyncio.run(main(*sys.argv[1:4]))
