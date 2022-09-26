#!/usr/bin/env bash

set -euo pipefail

HEIGHT=26
DENSITY=144

for SRC in "$@"; do (
    DEST="internal/assets/$(basename "${SRC%.*}.png")"

    set -x

    inkscape \
        --export-height="$HEIGHT" \
        --export-filename="$DEST" \
        "$SRC"

    mogrify \
        -units PixelsPerInch \
        -density "$DENSITY" \
        "$DEST"
) done
