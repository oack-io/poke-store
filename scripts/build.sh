#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "==> Building Astro frontend..."
cd "$ROOT_DIR/web"
npm run build

echo "==> Copying dist to cmd/server/static..."
rm -rf "$ROOT_DIR/cmd/server/static"
cp -r "$ROOT_DIR/web/dist" "$ROOT_DIR/cmd/server/static"

echo "==> Building Go server..."
cd "$ROOT_DIR"
go build -o bin/poke-store ./cmd/server

echo "==> Build complete! Run with: ./bin/poke-store"
