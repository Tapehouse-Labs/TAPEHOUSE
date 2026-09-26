#!/usr/bin/env python3
"""Print a RedStone payload built from the latest redstone-primary-prod data packages.

Usage: redstone-payload.py DATA_PACKAGE_ID [DATA_PACKAGE_ID ...]

A data package ID is a single feed such as NVDA---24_7, or a group such as NY_MARKET_STATUS.

Reads the keyless public gateways, or the main gateway when REDSTONE_API_KEY is set.
"""

import base64
import json
import os
import sys
import urllib.request
from decimal import Decimal

PATH = "/data-packages/latest/redstone-primary-prod"
MAIN_GATEWAY = "https://oracle-gateway.a.redstone.finance"
PUBLIC_GATEWAYS = (
    "https://oracle-gateway-1.a.redstone.finance",
    "https://oracle-gateway-2.a.redstone.finance",
)
REDSTONE_MARKER = bytes.fromhex("000002ed57011e0000")
DEFAULT_DECIMALS = 8


def latest_packages() -> dict:
    key = os.environ.get("REDSTONE_API_KEY")
    requests = (
        [urllib.request.Request(MAIN_GATEWAY + PATH, headers={"x-api-key": key})]
        if key
        else []
    )
    requests += [urllib.request.Request(gateway + PATH) for gateway in PUBLIC_GATEWAYS]
    error = None
    for request in requests:
        try:
            with urllib.request.urlopen(request, timeout=30) as response:
                return json.load(response, parse_float=Decimal)
        except OSError as exc:
            error = exc
    raise SystemExit(f"No RedStone gateway answered: {error}")


def feed_id(feed: str) -> bytes:
    if len(feed.encode()) > 31:
        raise SystemExit(f"{feed}: IDs longer than 31 bytes are not supported")
    return feed.encode().ljust(32, b"\0")


def data_point(point: dict) -> bytes:
    value = Decimal(point["value"]).scaleb(point.get("decimals", DEFAULT_DECIMALS))
    if value != value.to_integral_value():
        raise SystemExit(
            f"{point['dataFeedId']}: value {point['value']} has too many decimals"
        )
    return feed_id(point["dataFeedId"]) + int(value).to_bytes(32, "big")


def signed_package(package: dict) -> bytes:
    points = sorted(
        package["dataPoints"], key=lambda point: feed_id(point["dataFeedId"])
    )
    return (
        b"".join(data_point(point) for point in points)
        + package["timestampMilliseconds"].to_bytes(6, "big")
        + (32).to_bytes(4, "big")
        + len(points).to_bytes(3, "big")
        + base64.b64decode(package["signature"])
    )


def payload(data_package_ids: list) -> bytes:
    latest = latest_packages()
    packages = [
        package for package_id in data_package_ids for package in latest[package_id]
    ]
    if len({package["timestampMilliseconds"] for package in packages}) != 1:
        raise SystemExit("The latest packages carry different timestamps. Retry.")
    return (
        b"".join(signed_package(package) for package in packages)
        + len(packages).to_bytes(2, "big")
        + (0).to_bytes(3, "big")
        + REDSTONE_MARKER
    )


if __name__ == "__main__":
    if len(sys.argv) < 2:
        raise SystemExit(__doc__)
    print("0x" + payload(sys.argv[1:]).hex())
