#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_NAME='Nightscout Menu Bar'

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf dist/*.app
mkdir -p dist

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
cp assets/info.plist "$APP_CONTENTS"
mkdir "$APP_CONTENTS/Resources"
cp assets/Nightscout.icns "$APP_CONTENTS/Resources"
mkdir "$APP_CONTENTS/MacOS"
mv "dist/$BINARY_NAME" "$APP_CONTENTS/MacOS"
echo ...done
