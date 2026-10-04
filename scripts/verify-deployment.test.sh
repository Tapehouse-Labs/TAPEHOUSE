#!/bin/sh
# SPDX-License-Identifier: MIT OR Apache-2.0
set -u
here=$(cd "$(dirname "$0")" && pwd)
dir=$(mktemp -d)
trap 'kill $pid 2>/dev/null; wait $pid 2>/dev/null; rm -rf "$dir"' EXIT

cat > "$dir/stub.py" <<'PY'
import http.server, sys

class Handler(http.server.BaseHTTPRequestHandler):
    def reply(self, code, body=b""):
        self.send_response(code)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        if self.path == "/v1/status":
            self.reply(200, b'{"chainId":46630,"start":1,"head":1000,"finalized":990}')
        elif self.path == "/down":
            self.reply(503)
        else:
            self.reply(404)

    def do_POST(self):
        self.rfile.read(int(self.headers.get("Content-Length", 0)))
        self.reply(200, b'{"jsonrpc":"2.0","id":1,"result":"0x3f2"}')

    def log_message(self, *args):
        pass

server = http.server.HTTPServer(("127.0.0.1", 0), Handler)
print(server.server_port, flush=True)
server.serve_forever()
PY
python3 "$dir/stub.py" > "$dir/port" &
pid=$!
while [ ! -s "$dir/port" ]; do sleep 0.1; done
url=http://127.0.0.1:$(cat "$dir/port")
failures=0

expect() {
  want=$1; shift
  "$@" > "$dir/out" 2>&1
  got=$?
  if [ "$got" -ne "$want" ]; then
    echo "not ok: expected exit $want, got $got: $*"; cat "$dir/out"; failures=$((failures + 1))
  else
    echo "ok: $*"
  fi
}

export TAPEHOUSE_RPC_URL=$url
expect 0 env MAX_LAG=10 "$here/verify-deployment.sh" "$url"
expect 0 env MAX_LAG=10 SERVICE_URLS="$url/other" "$here/verify-deployment.sh" "$url"
expect 1 env MAX_LAG=2 "$here/verify-deployment.sh" "$url"
expect 1 env MAX_LAG=10 SERVICE_URLS="$url/down" "$here/verify-deployment.sh" "$url"
expect 1 env MAX_LAG=10 SERVICE_URLS="http://127.0.0.1:1" "$here/verify-deployment.sh" "$url"
expect 1 "$here/verify-deployment.sh" "http://127.0.0.1:1"
expect 2 "$here/verify-deployment.sh"
exit "$failures"
