#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_NAME='Nightscout Menu Bar'
VERSION="${VERSION:-}"
ICONSET=darwin/Nightscout.iconset
ICNS=darwin/Nightscout.icns

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf dist/*.app assets/{"$ICONSET","$ICNS"}
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
  echo Build "$BINARY_NAME-$ARCH"
  GOARCH="$ARCH" go build -ldflags='-w -s' -trimpath -o "dist/$BINARY_NAME-$ARCH" .
done
lipo -create -output "dist/$BINARY_NAME" "dist/$BINARY_NAME-amd64" "dist/$BINARY_NAME-arm64"
rm "dist/$BINARY_NAME-amd64" "dist/$BINARY_NAME-arm64"
echo ...done

echo Generate "$APP_NAME.app"
APP_CONTENTS="dist/$APP_NAME.app/Contents"
mkdir -p "$APP_CONTENTS"
go run ./assets/darwin/info --version="$VERSION" > "$APP_CONTENTS/info.plist"
mkdir "$APP_CONTENTS/Resources"
cp "assets/$ICNS" "$APP_CONTENTS/Resources"
mkdir "$APP_CONTENTS/MacOS"
mv "dist/$BINARY_NAME" "$APP_CONTENTS/MacOS"

echo Compress "$APP_NAME.app"
tar_name="dist/nightscout-menu-bar_darwin.tar.gz"
tar -czvf "$tar_name" -C dist "$APP_NAME.app"
go run ./assets/darwin/cask --path="$tar_name" --version="$VERSION" > dist/nightscout-menu-bar.rb
echo ...done
