#!/bin/bash

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
  FILE="bb-darwin-amd64"
elif [ "$OS" = "linux" ]; then
  FILE="bb-linux-amd64"
elif [[ "$OS" == *"mingw"* || "$OS" == *"cygwin"* ]]; then
  FILE="bb-windows-amd64.exe"
else
  echo "Unsupported OS"
  exit 1
fi

# Download the binary
echo "Downloading $FILE..."
curl -L -o bb "https://github.com/srijan619/bb-terminal/releases/latest/download/$FILE"

# Make executable and move to a directory in PATH
chmod +x bb
sudo mv bb /usr/local/bin/

echo "Installation complete! Run 'bb' to start."
