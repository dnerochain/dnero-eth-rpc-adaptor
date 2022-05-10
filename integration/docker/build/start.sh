#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

$GOBIN build -o ./build/linux/dnero-eth-rpc ./cmd/dnero-eth-rpc

set +x 

echo "Done."



