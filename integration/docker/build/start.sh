#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

$GOBIN build -o ./build/linux/dnero-eth-rpc-adaptor ./cmd/dnero-eth-rpc-adaptor

set +x 

echo "Done."



