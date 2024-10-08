#!/usr/bin/env bash

BINARY_NAME='Nightscout Menu Bar'
VERSION="${VERSION:-}"
WIN_VERSION="$VERSION"
if [[ -n "$WIN_VERSION" ]]; then
  WIN_VERSION="${WIN_VERSION#v}"
  WIN_VERSION+='.0'
fi

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf -- dist/*.exe assets/windows/*.png *.syso
mkdir -p dist

command -v go-winres &>/dev/null || go install github.com/tc-hib/go-winres@latest

# Generate metadata
cp -a assets/{png/icon_16x16,windows/icon16}.png
cp -a assets/{png/icon_32x32,windows/icon32}.png
cp -a assets/{png/icon_48x48,windows/icon48}.png
cp -a assets/{png/icon_64x64,windows/icon64}.png
cp -a assets/{png/icon_128x128,windows/icon128}.png
cp -a assets/{png/icon_256x256,windows/icon256}.png
go-winres make --arch=amd64,arm64 --in=assets/windows/winres.json --product-version="$WIN_VERSION" --file-version="$WIN_VERSION"

# Build binary
export GOOS=windows CGO_ENABLED=1
echo Build "$BINARY_NAME.exe"
go build -ldflags="-w -s -H=windowsgui -X main.version=$VERSION" -trimpath -o "dist/$BINARY_NAME.exe" .
