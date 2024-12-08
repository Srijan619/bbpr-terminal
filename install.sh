#!/bin/bash
# Download the correct binary for the OS
curl -LO https://github.com/srijan619/bbpr/releases/download/latest/bbpr-linux-amd64
chmod +x bbpr-linux-amd64
sudo mv bbpr-linux-amd64 /usr/local/bin/bbpr
