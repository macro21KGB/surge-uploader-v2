#!/bin/sh

set -e

REPO="macro21KGB/surge-uploader-v2"
BINARY_NAME="surge-uploader"
INSTALL_DIR="$HOME/.local/bin"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux*)  OS="linux" ;;
  Darwin*) OS="darwin" ;;
  *)
    echo "Unsupported operating system: $OS" >&2
    exit 1
    ;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

if [ "$OS" = "linux" ] && [ "$ARCH" = "arm64" ]; then
  echo "Linux arm64 builds are not currently available." >&2
  exit 1
fi

ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"

TMP_FILE="$(mktemp)"
trap 'rm -f "$TMP_FILE"' EXIT

echo "Downloading ${ASSET_NAME} from latest release..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

echo "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"
mv -f "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"

echo "${BINARY_NAME} installed successfully! Run '${BINARY_NAME} help' to get started."

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "Note: ${INSTALL_DIR} is not in your PATH."
    echo "Add it by adding the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    ;;
esac
