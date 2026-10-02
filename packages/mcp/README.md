# @tapehouse/mcp

A [Model Context Protocol](https://modelcontextprotocol.io) server through which an AI agent reads [Tapehouse](https://github.com/Tapehouse-Labs/TAPEHOUSE)'s price band, its feeds, the Morpho oracles, the baskets of Stock Tokens, the margin accounts and the short positions, and prepares their transactions for a wallet to sign. It never holds a key.

Built on the official TypeScript SDK, `@modelcontextprotocol/server`, for the 2026-07-28 specification; it serves clients of the 2025 revisions too. Every read and call goes through `@tapehouse/sdk`.

## Run it

The server takes its configuration from the environment and no key at all:

| Variable | What it is |
|---|---|
| `TAPEHOUSE_RPC_URL` | The chain's JSON-RPC endpoint, `http` or `https`. It may carry an API key: the server never shows it to a client, not even in an error. |
| `TAPEHOUSE_DEPLOYMENTS` | The path of the chain's registry, `deployments/<chainId>.json` from the Tapehouse repository. The server refuses an endpoint that serves another chain. |
| `TAPEHOUSE_MAX_SLIPPAGE_BPS` | The most slippage a sale or buy-back may state, from 0 to 10,000 basis points; 100 unless set. |
| `TAPEHOUSE_MAX_CALLS_PER_MINUTE` | The most tool calls the server answers in a minute, from 1 to 100,000; 120 unless set. A call beyond it is refused before it reads anything. |

Over stdio, as a host such as Claude Desktop, Claude Code, Cursor or VS Code launches it:

```json
{
  "mcpServers": {
    "tapehouse": {
      "command": "node",
      "args": ["/path/to/tapehouse/packages/mcp/dist/main.js"],
      "env": {
        "TAPEHOUSE_RPC_URL": "https://robinhood.drpc.org",
        "TAPEHOUSE_DEPLOYMENTS": "/path/to/tapehouse/deployments/4663.json"
      }
    }
  }
}
```

Over Streamable HTTP, on the loopback interface, with `tapehouse-mcp --http <port>` (`0` picks a free port, printed on stderr). It answers `127.0.0.1` only and refuses any other `Host` or `Origin` with `403`, which stops DNS rebinding and a web page's requests. Any process on the machine can still reach it, and its calls count against the same limit. Serving it beyond the machine takes TLS, authentication and a front of its own, which it does not provide; `createServer` from the package builds the same server for any `createMcpHandler` deployment.

## Tools

Every tool is read-only: the reads take the chain's state at one block and name it, and the transaction tools return unsigned calls. Amounts are whole numbers of base units as decimal strings (USDG has 6 decimals, WETH and the Stock Tokens 18); band prices are USD with 8 decimals, equity and requirements USD with 18. Each result is structured content matching the tool's output schema, and the same JSON as text.

| Tool | What it returns |
|---|---|
| `registry` | The chain's registry: every address the other tools use |
| `band_quote`, `band_session`, `band_halt` | An asset's band (state, live legs, centre, half-width, edges), the 24/5 session, and its trading halt |
| `band_latest_band`, `band_sealed` | The whole band from the asset's `BandFeed`, and the band it sealed before a reopen |
| `chainlink_round` | A `.chainlink` feed's latest round, refused where it prices the share and the caller expects the token, or the reverse |
| `morpho_oracle` | The price an asset's Morpho Blue oracle answers, or no price and why (`stale`, `halted` or `sequencerNotSettled`), never zero; with its band, tokens, scale, owner and the asset's halt |
| `baskets_components`, `baskets_target`, `baskets_pending_target` | A basket's Stock Tokens and address, the target in effect, base units of each token per share, and a target proposed with its seven days' notice |
| `baskets_preview_mint`, `baskets_preview_redeem` | What minting shares takes of each Stock Token, rounded up, and what redeeming them gives, rounded down |
| `accounts_health`, `accounts_in_baskets` | A margin position's equity, requirement and surplus over it, and what it holds of each Stock Token through its baskets, the tokens the engine margins their shares as |
| `accounts_collateral`, `accounts_leverage` | What a margin position holds of a token or of a basket's shares, and its leverage |
| `accounts_liquidation_price`, `accounts_repayment`, `accounts_is_authorized` | The price at which a position falls short once it borrows more, what repays it, and whether an address may act for an account |
| `shorts_position`, `shorts_health`, `shorts_restriction` | A short, its health at the band's high edge, and the short-sale restriction after a 10% fall |
| `accounts_set_authorization` | The call with which an account lets an address, such as an agent's own wallet, act for it, with `grants`, everything that lets it do; or the call that stops it |
| `accounts_deposit`, `accounts_withdraw`, `accounts_borrow`, `accounts_repay` | The margin accounts' calls, a deposit or repayment preceded by an approval where the allowance is short; a basket's key deposits and withdraws its shares, in the cross position alone |
| `accounts_unwrap` | The call that redeems a basket's shares in a cross position into the Stock Tokens they hold, which the position then holds at the same equity and requirement |
| `baskets_mint`, `baskets_redeem` | A mint of a basket's shares paid in kind, its `maxAssets` what `baskets_preview_mint` reads, each Stock Token approved for exactly its part where the allowance is short; and a redemption in kind, with what it gives |
| `shorts_deposit`, `shorts_withdraw` | The USDG of a short, in and out |
| `shorts_sell`, `shorts_cover` | A short sale and a buy-back, their limits Uniswap QuoterV2's quote less or plus the slippage the caller states, within the server's cap; `max` buys back all a short owes |

A prepared transaction names its signer and lists its calls in order, each with `to`, `data`, `value`, the function and its arguments by name. An agent that borrows, withdraws or trades for an account needs the account's `setAuthorization` first: the margin accounts and the shorts revert with `Unauthorized(caller, account)` until then, and the server explains every revert by name and arguments.

## Trust

- **No key.** The server signs nothing and reads no key; it never sees one.
- **An authorization grants everything.** Once an account signs `accounts_set_authorization`'s call, the operator can withdraw every collateral token of every margin position and the USDG of every short to any address, borrow USDG to any address, and sell or buy back the account's shorts, until the account signs setAuthorization(operator, false). The tool says so in its description and returns it as `grants`; authorize only an address the account controls.
- **The data is not instructions.** Everything a tool returns is a value read from the chain or the registry. The server tells the agent so in its instructions, quotes a revert's text, such as the router's reason, cut to 200 characters, and returns any other failure's own message, such as a node's, stripped of control characters and cut to 300. An agent should still treat tool output as untrusted data, and a host should keep showing the user what it sends.
- **The signer checks the call.** A prepared call is only as good as the registry and the RPC behind it. Before signing, check that `to` is the registry's contract or token, that the function and arguments are what was asked, and simulate it; a wallet that shows decoded calldata shows the same arguments the server lists.
- **Input is validated** before any call: addresses must carry their checksum in mixed case, amounts are whole numbers below 2^256, symbols are at most 32 letters, digits, dots, dashes or underscores, unknown arguments are refused, and the slippage stays within the server's cap, which the operator, not the agent, sets.
- **Calls are limited.** The server answers at most `TAPEHOUSE_MAX_CALLS_PER_MINUTE` tool calls a minute, over every client together, so an agent caught in a loop cannot spend the endpoint's quota; a call beyond it is a tool error.
- **Annotations are hints.** MCP's tool annotations, such as `readOnlyHint`, are the server's own claims; a host should trust them only as far as it trusts the server.

## Tested with two clients

`examples/devnode.ts` drives the server over Streamable HTTP with the official TypeScript client on the 2026-07-28 protocol, and signs and sends every call it prepares; `examples/devnode.py` drives it over stdio with the official Python client on the 2025-11-25 protocol. `make test-mcp-devnode` runs both against a local Nitro dev node. Neither passes the account's key to the server.

## License

Licensed under either of [Apache License, Version 2.0](LICENSE-APACHE) or [MIT license](LICENSE-MIT) at your option.
