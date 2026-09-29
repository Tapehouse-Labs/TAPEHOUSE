#!/usr/bin/env bash
# Usage: gas-table.sh [REPORT]
# Prints, as a Markdown table, every call the dev-node suite measured on both the margin program and its Solidity
# reference, from REPORT (stylus/.gas-devnode by default): the L2 gas of each and the reference's over the program's.
# The scenario set comes first by lattice size, then the requirement, then the views; a call measured on both
# without a place in that order is an error, so the table reads the same on every awk.
set -euo pipefail

report=${1:-$(cd "$(dirname "$0")/../.." && pwd)/stylus/.gas-devnode}
awk -v first='scenario(0)|scenarioDigest(32)|scenarioDigest(64)|scenarioDigest(128)|scenarioDigest(256)|scenarioDigest(256,launch)|requirement(3)|requirement(6)|currentRequirement(3)|currentRequirement(1 of 6)|assets()|volatility(NVDA)|correlation(NVDA,SPY)' '
  function commas(n,  s, out) { s = n ""; out = ""; while (length(s) > 3) { out = "," substr(s, length(s) - 2) out; s = substr(s, 1, length(s) - 3) } return s out }
  function row(call) { printf "| `%s` | %s | %s | %.1f× |\n", call, commas(program[call]), commas(reference[call]), reference[call] / program[call] }
  { gas = $NF; name = $0; sub(/ [0-9]+$/, "", name) }
  name ~ /^margin\./ { program[substr(name, 8)] = gas }
  name ~ /^reference\./ { reference[substr(name, 11)] = gas }
  END {
    print "| Call | Stylus | Solidity | Solidity / Stylus |"
    print "|---|---|---|---|"
    n = split(first, calls, "|")
    for (i = 1; i <= n; i++) { listed[calls[i]] = 1; if ((calls[i] in program) && (calls[i] in reference)) row(calls[i]) }
    for (call in reference) if ((call in program) && !(call in listed)) { print "gas-table.sh: no place for " call > "/dev/stderr"; unplaced = 1 }
    exit unplaced
  }' "$report"
