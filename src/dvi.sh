#!/usr/bin/env sh
set -eu

DOCS_DIR="${1}"
PKGS_PATTERN="${2}"

mkdir -p "$DOCS_DIR"

OUT="$DOCS_DIR/api.txt"
echo "Generating docs into $OUT" >&2

{
  echo "API docs (go doc)"
  echo "Generated: $(date)"
  echo

  go list "$PKGS_PATTERN" | while IFS= read -r pkg; do
    echo "============================================================"
    echo "PACKAGE: $pkg"
    echo "============================================================"
    echo
    go doc -all "$pkg"
    echo
  done
} > "$OUT"