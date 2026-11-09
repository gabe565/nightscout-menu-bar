#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_ID='com.gabe565.nightscout-menu-bar'
VERSION="${VERSION:-}"
OUT=dist/linux

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

if [[ ! -f internal/settings/frontend/dist/index.html ]]; then
  go generate ./internal/settings/frontend
fi

ARCH="$(go env GOARCH)"
rm -rf "$OUT" dist/"${BINARY_NAME}_linux_$ARCH".*
mkdir -p "$OUT"

# Build binary
export GOOS=linux CGO_ENABLED=1
echo Build "$BINARY_NAME ($ARCH)"
go build -ldflags="-w -s -X main.version=$VERSION" -trimpath -o "$OUT/$BINARY_NAME" .

# Desktop entry and icons
install -Dm644 "assets/linux/$APP_ID.desktop" "$OUT/share/applications/$APP_ID.desktop"
for ICON in assets/png/icon_*.png; do
  SIZE="${ICON#assets/png/icon_}"
  install -Dm644 "$ICON" "$OUT/share/icons/hicolor/${SIZE%.png}/apps/$APP_ID.png"
done
cp LICENSE README.md "$OUT"

echo Compress "$BINARY_NAME ($ARCH)"
tar -czvf "dist/${BINARY_NAME}_linux_$ARCH.tar.gz" -C "$OUT" .

# Packages
if [[ "$VERSION" =~ ^v[0-9] ]]; then
  PKG_VERSION="${VERSION#v}"
else
  PKG_VERSION="0.0.0-g$(git rev-parse --short HEAD)"
fi
export GOARCH="$ARCH" PKG_VERSION
for FORMAT in deb rpm; do
  echo Package "$BINARY_NAME ($ARCH $FORMAT)"
  go tool nfpm package --config assets/linux/nfpm.yaml --packager "$FORMAT" --target "dist/${BINARY_NAME}_linux_$ARCH.$FORMAT"
done
