#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_NAME='Nightscout Menu Bar'

set -euo pipefail

SCRIPT_DIR="$(dirname "$0")"
DIST="$SCRIPT_DIR/../dist"
ASSETS="$SCRIPT_DIR/../assets"

rm -rf "$DIST"/*

echo Build "$BINARY_NAME"
go build -ldflags='-w -s' -o "$DIST/$BINARY_NAME" "$(git rev-parse --show-toplevel)"
echo ...done

echo Generate "$APP_NAME.app"
APP_CONTENTS="$DIST/$APP_NAME.app/Contents"
mkdir -p "$APP_CONTENTS"
cp "$ASSETS/info.plist" "$APP_CONTENTS"
mkdir "$APP_CONTENTS/Resources"
cp "$ASSETS/Nightscout.icns" "$APP_CONTENTS/Resources"
mkdir "$APP_CONTENTS/MacOS"
cp -p "$DIST/$BINARY_NAME" "$APP_CONTENTS/MacOS"
echo ...done

echo Zip "$APP_NAME.app"
cd "$DIST"
zip -r "$APP_NAME.zip" "$APP_NAME.app"
echo ...done
