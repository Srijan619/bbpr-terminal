#!/bin/bash

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
  FILE="bbpr-darwin-amd64"
elif [ "$OS" = "linux" ]; then
  FILE="bbpr-linux-amd64"
elif [[ "$OS" == *"mingw"* || "$OS" == *"cygwin"* ]]; then
  FILE="bbpr-windows-amd64.exe"
else
  echo "Unsupported OS"
  exit 1
fi

# Download the binary
echo "Downloading $FILE..."
curl -L -o bbpr "https://github.com/srijan619/bbpr-terminal/releases/latest/download/$FILE"

# Make executable and move to a directory in PATH
chmod +x bbpr
sudo mv bbpr /usr/local/bin/

echo "Installation complete! Run 'bbpr' to start."
