#!/usr/bin/env bash

set -euo pipefail

HEIGHT=64
PAD="${3:-8}"
DENSITY=144

BASEHEIGHT="$(bc <<< "$HEIGHT - $PAD")"

SRC="$1"
DEST="$2"
COLOR="${4:-}"

if [[ -n "$COLOR" ]]; then
  TMP="$(mktemp -t convert-icon).svg"
  trap 'rm -f "$TMP"' EXIT
  sed "s/currentColor/$COLOR/g" "$SRC" > "$TMP"
  SRC="$TMP"
fi

inkscape "$SRC" \
  --export-height="$BASEHEIGHT" \
  --export-type=png \
  --export-filename=- \
| magick - \
  -strip \
  -background transparent \
  -gravity center \
  -extent "${HEIGHT}x${HEIGHT}" \
  -units PixelsPerInch \
  -density "$DENSITY" \
  "$DEST"
