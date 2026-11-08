#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

for icon in about error history open preferences quit reading; do
  ./convert-icon.sh "src/$icon.svg" "dist/$icon.png"
  ./convert-icon.sh "src/$icon.svg" "dist/$icon.ico"
  ./convert-icon.sh "src/$icon.svg" "dist/$icon-white.ico" 8 '#fff'
done

./convert-icon.sh src/nightscout.svg dist/nightscout.png 0
./convert-icon.sh src/nightscout.svg dist/nightscout.ico 0
./convert-icon.sh src/nightscout-transparent.svg dist/nightscout-transparent.png 4
