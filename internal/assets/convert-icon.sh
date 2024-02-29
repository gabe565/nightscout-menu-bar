#!/usr/bin/env bash

set -euo pipefail

HEIGHT=32
PAD="${3:-4}"
DENSITY=144

BASEHEIGHT="$(bc <<< "$HEIGHT - $PAD")"

SRC="$1"
DEST="$2"

inkscape "$SRC" \
  --export-height="$BASEHEIGHT" \
  --export-type=png \
  --export-filename=- \
| convert - \
  -strip \
  -background transparent \
  -gravity center \
  -extent "${HEIGHT}x${HEIGHT}" \
  -units PixelsPerInch \
  -density "$DENSITY" \
  "$DEST"
