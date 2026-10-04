# Deploying the services to Railway

Five services deploy from one image recipe, `services/Dockerfile`, which builds the binary named by the `SERVICE` build argument: `indexer` (also the API, the WebSocket stream and the RedStone relay), `keeper`, `bidder`, `sponsor` and `alerts`. The image is static, runs as a non-root user on a pinned distroless base and carries `deployments/*.json`, so the network is chosen only by variables.

Each service has its own config file in `deploy/railway/<service>/railway.json`. Set it as the service's config file path; the build context is the repository root and the Dockerfile path is `services/Dockerfile`. Railway passes service variables to the build as build arguments, so give every service a `SERVICE` variable naming its binary.

To build locally: `docker build -f services/Dockerfile --build-arg SERVICE=indexer .`

## Network

Nothing in the image names a network. Point `TAPEHOUSE_DEPLOYMENTS` at `/deployments/46630.json` (Robinhood Chain testnet) or `/deployments/4663.json` (mainnet) and `TAPEHOUSE_RPC_URL` at that chain's RPC. A service refuses to start when the registry and the RPC are of different chains.

## Secrets

The services read secrets only from files. Railway holds them as variables, so the image's entrypoint turns each variable `NAME_B64` (the secret, base64-encoded) into a `0600` file under `/run/secrets/NAME` and sets the variable the service reads: `NAME_FILE` when `NAME` ends in `_PASSWORD`, otherwise `NAME` itself, holding the file's path. It refuses to start if the target variable is also set or the value is not base64. Encode with `base64 < file | tr -d '\n'` and store the result as a sealed variable.

## Variables

Shared by `indexer`, `keeper`, `sponsor` and `alerts`: `SERVICE`, `TAPEHOUSE_DEPLOYMENTS`, `TAPEHOUSE_RPC_URL`. Optional on all four: `TAPEHOUSE_RPC_RATE` (requests a second, default 20). `PORT` need not be set; listen addresses are set below.

| Service | Required | Optional |
| --- | --- | --- |
| `indexer` | `TAPEHOUSE_LISTEN=[::]:8080`, `TAPEHOUSE_DB=/data/index.db` | `TAPEHOUSE_API_KEYS` (SHA-256 digests of the keyed tier's keys, comma-separated), `TAPEHOUSE_TRUSTED_PROXIES` (CIDRs whose `X-Forwarded-For` is trusted), `REDSTONE_API_KEY`, `REDSTONE_BACKUP_API_KEY` |
| `keeper` | `TAPEHOUSE_KEEPER_KEYSTORE_B64`, `TAPEHOUSE_KEEPER_PASSWORD_B64` | `TAPEHOUSE_INDEXER_URL`, `TAPEHOUSE_KEEPERS`, `TAPEHOUSE_RPC_FALLBACK_URLS`, `TAPEHOUSE_HALT_KEYSTORE_B64` with `TAPEHOUSE_HALT_PASSWORD_B64` (without them the halt keeper does not run), `TAPEHOUSE_HALTS_URL`, `TAPEHOUSE_KEEPER_BUY`, `TAPEHOUSE_KEEPER_REBALANCE`, `REDSTONE_API_KEY`, `REDSTONE_BACKUP_API_KEY` |
| `bidder` | `PRIVATE_KEY` (a hex key; the bidder has no keystore option, so keep it a sealed variable) | `TAPEHOUSE_INDEXER_URL`, `TAPEHOUSE_INDEXER_KEY`, `TAPEHOUSE_BIDDER_STATE=/data/bidder.json` (needs a volume to survive a restart), `TAPEHOUSE_BIDDER_BUDGET`, `TAPEHOUSE_BIDDER_DISCOUNT_BPS`, `TAPEHOUSE_BIDDER_COMMIT_LEAD`, `TAPEHOUSE_BIDDER_ASSETS`, `TAPEHOUSE_BIDDER_DRY_RUN` |
| `sponsor` | `TAPEHOUSE_SPONSOR_ADDR=[::]:4338`, `TAPEHOUSE_SPONSOR_KEYSTORE_B64`, `TAPEHOUSE_SPONSOR_PASSWORD_B64` | `TAPEHOUSE_SPONSOR_DAILY`, `TAPEHOUSE_TRUSTED_PROXIES` |
| `alerts` | `TAPEHOUSE_INDEXER_URL`, `TAPEHOUSE_ALERTS_URL` (the public https address links in emails point to), `TAPEHOUSE_ALERTS_SECRET` (32 characters or more), `TAPEHOUSE_ALERTS_LISTEN=[::]:8090`, `TAPEHOUSE_ALERTS_DB=/data/alerts.db`, and a channel: `TAPEHOUSE_SMTP_URL` with `TAPEHOUSE_ALERTS_FROM`, or `TELEGRAM_BOT_TOKEN` | `TAPEHOUSE_ALERTS_INTERVAL`, `TAPEHOUSE_RPC_FALLBACK_URLS`, `TAPEHOUSE_TRUSTED_PROXIES`, `TELEGRAM_API_URL` |

Services reach each other over Railway's private network, for example `TAPEHOUSE_INDEXER_URL=http://indexer.railway.internal:8080`.

## Volumes and replicas

The indexer keeps its index in SQLite, so it needs a volume mounted at `/data` and exactly one replica (`numReplicas` is 1 in its config; do not raise it). `alerts` keeps its subscriptions in SQLite and needs the same: a volume at `/data` and one replica. `bidder` needs a volume at `/data` for its state file. Railway mounts volumes as root, so on each of these three set `RAILWAY_RUN_UID=0`. `keeper` and `sponsor` keep no state and need no volume.

## Health checks

Only the indexer has a health endpoint: `GET /v1/status`, which the indexer's config uses. `sponsor` accepts only `POST`, `alerts` serves only `/v1/...` routes, and `keeper` and `bidder` serve no HTTP, so their configs set no health check and rely on the restart policy.

## Order of creation

1. `indexer`, with its volume. Wait until `/v1/status` answers.
2. `keeper`, with `TAPEHOUSE_INDEXER_URL` pointing at the indexer.
3. `sponsor`.
4. `alerts`, with its volume and `TAPEHOUSE_INDEXER_URL`.
5. `bidder`, with its volume.

## Verifying a deployment

```sh
TAPEHOUSE_RPC_URL=<chain rpc> MAX_LAG=20 \
SERVICE_URLS="https://sponsor.example https://alerts.example" \
scripts/verify-deployment.sh https://indexer.example
```

It checks that the API answers, that the indexer's head is within `MAX_LAG` blocks (default 20) of the chain's, and that each URL in `SERVICE_URLS` answers with a status below 500. It exits non-zero on any failure. `scripts/verify-deployment.test.sh` exercises it against a local stub.
