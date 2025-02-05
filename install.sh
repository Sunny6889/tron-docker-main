#!/bin/bash
set -euo pipefail

# Configuration
REPO_OWNER="jesseduffield"
REPO_NAME="tron-docker"
RELEASE_TAG="v0.45.2"
CHECKSUM_FILE="checksums.txt"
RACKAGE_PREFIX="trond"

# Determine the OS and architecture, then set the ASSET_NAME
OS=$(uname -s)
ARCH=$(uname -m)
case "$OS" in
    Linux*)
        case "$ARCH" in
            x86_64) ASSET_NAME="${RACKAGE_PREFIX}_${RELEASE_TAG}_Linux_x86_64.tar.gz";;
            arm64)  ASSET_NAME="${RACKAGE_PREFIX}_${RELEASE_TAG}_Linux_arm64.tar.gz";;
            *)      echo "Unsupported architecture: $ARCH"; exit 1;;
        esac
        ;;
    Darwin*)
        case "$ARCH" in
            x86_64) ASSET_NAME="${RACKAGE_PREFIX}_${RELEASE_TAG}_Darwin_x86_64.tar.gz";;
            arm64)  ASSET_NAME="${RACKAGE_PREFIX}_${RELEASE_TAG}_Darwin_arm64.tar.gz";;
            *)      echo "Unsupported architecture: $ARCH"; exit 1;;
        esac
        ;;
    *) echo "Unsupported OS: $OS"; exit 1;;
esac

# Download URLs
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${RELEASE_TAG}/${ASSET_NAME}"
CHECKSUM_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${RELEASE_TAG}/${CHECKSUM_FILE}"

# Download files
echo "Downloading ${ASSET_NAME}..."
curl -L -O "${DOWNLOAD_URL}" --fail --progress-bar

echo "Downloading checksum file..."
curl -L -O "${CHECKSUM_URL}" --fail --progress-bar

# Verify checksum
echo "Verifying checksum..."
if ! sha256sum --check --ignore-missing "${CHECKSUM_FILE}"; then
  echo "Checksum validation failed!"
  exit 1
fi

echo "✅ Success! ${ASSET_NAME} is valid."