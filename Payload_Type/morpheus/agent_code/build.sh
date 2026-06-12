#!/bin/bash
# Cross-compile Morpheus agent for Windows x64 with evasion & garble obfuscation
# Usage: ./build.sh [output_name]

set -e

OUTPUT="${1:-morpheus.exe}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

echo "[*] Building Morpheus agent for Windows x64..."
echo "[*] OS: $GOOS  ARCH: $GOARCH  CGO: $CGO_ENABLED"

cd "$SCRIPT_DIR"

if command -v garble &>/dev/null; then
	echo "[*] Using garble for Go runtime obfuscation..."
	GARBLE_BIN=$(command -v garble)
elif [ -f ~/go/bin/garble ]; then
	GARBLE_BIN=~/go/bin/garble
else
	GARBLE_BIN=""
fi

if [ -n "$GARBLE_BIN" ]; then
	$GARBLE_BIN -seed=random build -ldflags "-s -w" -o "$OUTPUT" .
else
	echo "[!] garble not found, building without Go runtime obfuscation"
	go build -trimpath -ldflags "-s -w" -o "$OUTPUT" .
fi

echo "[+] Build complete: $OUTPUT"
ls -lh "$OUTPUT"
