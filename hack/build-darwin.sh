#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_NAME='Nightscout Menu Bar'
VERSION="${VERSION:-}"
ICONSET=darwin/Nightscout.iconset
ICNS=darwin/Nightscout.icns

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf dist/{amd64,arm64} assets/{"$ICONSET","$ICNS"}
mkdir -p dist

# Generate icns
cp -a assets/{png,"$ICONSET"}
cp -a assets/"$ICONSET"/icon_{32x32,16x16@2x}.png
rm assets/"$ICONSET"/icon_48x48.png
cp -a assets/"$ICONSET"/icon_{64x64,32x32@2x}.png
cp -a assets/"$ICONSET"/icon_{128x128,64x64@2x}.png
cp -a assets/"$ICONSET"/icon_{256x256,128x128@2x}.png
cp -a assets/"$ICONSET"/icon_{512x512,256x256@2x}.png
iconutil --convert icns --output "assets/$ICNS" "assets/$ICONSET"

export GOOS=darwin CGO_ENABLED=1
for ARCH in amd64 arm64; do
  echo Build "$BINARY_NAME ($ARCH)"
  APP_CONTENTS="dist/$ARCH/$APP_NAME.app/Contents"
  mkdir -p "$APP_CONTENTS/MacOS" "$APP_CONTENTS/Resources"
  GOARCH="$ARCH" go build -ldflags="-w -s -X main.version=$VERSION" -trimpath -o "$APP_CONTENTS/MacOS/$BINARY_NAME" .
  go run ./assets/darwin/info --version="$VERSION" > "$APP_CONTENTS/info.plist"
  cp "assets/$ICNS" "$APP_CONTENTS/Resources"

  echo Compress "$APP_NAME.app ($ARCH)"
  tar -czvf "dist/${BINARY_NAME}_darwin_$ARCH.tar.gz" -C "dist/$ARCH" "$APP_NAME.app"
  echo ...done
done
