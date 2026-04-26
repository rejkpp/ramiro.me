#!/usr/bin/env bash
set -euo pipefail

TAILWIND_BIN="./tailwindcss"
TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"

# Download Tailwind CLI if not present.
if [ ! -f "$TAILWIND_BIN" ]; then
    echo "Downloading Tailwind CLI (Linux x64)..."
    curl -sSL -o "$TAILWIND_BIN" "$TAILWIND_URL"
    chmod +x "$TAILWIND_BIN"
fi

echo "Building CSS..."
"$TAILWIND_BIN" -i ./static/css/input.css -o ./static/css/app.css --minify

echo "Generating static site..."
go run ./cmd/gen

echo "Done. Output in ./public/"
