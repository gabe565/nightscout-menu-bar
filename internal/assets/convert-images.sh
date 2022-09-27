#!/usr/bin/env bash

set -euo pipefail

HEIGHT=32
PAD=4
DENSITY=144

BASEHEIGHT="$(bc <<< "$HEIGHT - $PAD")"

for SRC in "$@"; do (
  DEST="$(basename "${SRC%.*}.png")"

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
) done
