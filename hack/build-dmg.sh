#!/usr/bin/env bash

BINARY_NAME='nightscout-menu-bar'
APP_NAME='Nightscout Menu Bar'
BUNDLE_ID='com.gabe565.nightscout-menu-bar'
AUTHOR='gabe565'

set -euo pipefail

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
DIST="$SCRIPT_DIR/../dist"

rm -rf "$DIST"/*

echo Build "$BINARY_NAME"
go build -ldflags='-w -s' -o "$DIST/$BINARY_NAME" "$(git rev-parse --show-toplevel)"
echo ...done

echo Generate icons
"$SCRIPT_DIR/build-icns.sh"
echo ...done

echo Generate "$APP_NAME.app"
cd "$DIST"
appify \
			-id "$BUNDLE_ID" \
			-author "$AUTHOR" \
			-name "$APP_NAME" \
			-icon "$APP_NAME.icns" \
			"$BINARY_NAME"
echo ...done

echo Zip "$APP_NAME.app"
zip -r "$APP_NAME.zip" "$APP_NAME.app"
echo ...done