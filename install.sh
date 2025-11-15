#!/bin/bash

set -e

REPO="lucasepe/t2i"
BINARY="t2i"

usage() {
  echo "Usage: $0 [version]"
  exit 1
}

if [[ $# -gt 1 ]]; then
  usage
elif [[ $# -eq 1 ]]; then
  VERSION="$1"
  LATEST_TAG="v$VERSION"
else
  # Try to fetch the latest release
  JSON=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
  LATEST_TAG=$(echo "$JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

  # Fallback: if no official "latest" release exists, use the most recent tag
  if [[ -z "$LATEST_TAG" || "$LATEST_TAG" == "null" ]]; then
    echo "⚠️  No 'latest' release found, falling back to tags..."
    LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/tags" \
      | grep '"name":' \
      | sed -E 's/.*"([^"]+)".*/\1/' \
      | head -n 1)
  fi

  if [[ -z "$LATEST_TAG" || "$LATEST_TAG" == "null" ]]; then
    echo "❌ Could not determine latest version."
    exit 1
  fi

  VERSION="${LATEST_TAG#v}"
fi

# Detect OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux|darwin) ;;
  msys*|cygwin*|mingw*) OS="windows" ;;
  *) echo "❌ Unsupported OS: $OS" && exit 1 ;;
esac

# Detect arch
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Always zip — this matches your release assets
EXT="zip"
ASSET="${BINARY}-${OS}-${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET}"

echo "📦 Downloading $ASSET from $LATEST_TAG..."
echo "🔗 $URL"
curl -sL "$URL" -o "$ASSET"

TMP_DIR=$(mktemp -d)

echo "📂 Extracting to $TMP_DIR..."
unzip -o "$ASSET" -d "$TMP_DIR" >/dev/null

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  echo "⚠️  No permission for /usr/local/bin, falling back to $HOME/.local/bin"
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "👉 Make sure $INSTALL_DIR is in your PATH"
fi

BIN_PATH=$(find "$TMP_DIR" -type f -name "$BINARY" -perm -111 | head -n 1)
if [[ -z "$BIN_PATH" ]]; then
  echo "❌ Could not find the '$BINARY' binary inside ZIP"
  exit 1
fi

echo "🚀 Installing $BINARY to $INSTALL_DIR..."
chmod +x "$BIN_PATH"
mv "$BIN_PATH" "$INSTALL_DIR/$BINARY"

rm -rf "$ASSET" "$TMP_DIR"

echo "✅ $BINARY $VERSION installed successfully!"
