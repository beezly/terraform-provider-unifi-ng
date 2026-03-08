#!/usr/bin/env bash
# Fetch the latest UniFi Network OpenAPI spec from beez.ly/unifi-apis/
# Writes to specs/network-<version>.json and specs/latest-version.txt
set -euo pipefail

SPEC_BASE_URL="https://beez.ly/unifi-apis"
INDEX_URL="$SPEC_BASE_URL/"
SPECS_DIR="$(dirname "$0")/../specs"
mkdir -p "$SPECS_DIR"

echo "Checking latest UniFi API spec version..."

# Scrape the index page to find the latest version
LATEST=$(curl -s "$INDEX_URL" | grep -oP 'network-\K[0-9]+\.[0-9]+\.[0-9]+(?=\.json)' | sort -V | tail -1)

if [ -z "$LATEST" ]; then
    echo "ERROR: Could not determine latest spec version from $INDEX_URL" >&2
    exit 1
fi

echo "Latest version: $LATEST"

# Check if we already have this version
CURRENT=""
if [ -f "$SPECS_DIR/latest-version.txt" ]; then
    CURRENT=$(cat "$SPECS_DIR/latest-version.txt")
fi

if [ "$CURRENT" = "$LATEST" ] && [ -f "$SPECS_DIR/network-$LATEST.json" ]; then
    echo "Already at latest version ($LATEST), nothing to do."
    exit 0
fi

# Download the new spec
SPEC_URL="$SPEC_BASE_URL/network-$LATEST.json"
echo "Downloading $SPEC_URL..."
curl -sf "$SPEC_URL" -o "$SPECS_DIR/network-$LATEST.json"
echo "$LATEST" > "$SPECS_DIR/latest-version.txt"
echo "Downloaded network-$LATEST.json"
echo "NEW_VERSION=$LATEST"
