## dnero-eth-rpc-adaptor

The `dnero-eth-rpc-adaptor` project is aiming to provide an adaptor which translates the Dnero RPC interface to the Ethereum RPC interface. Please find the currently supported Ethereum RPC APIs [here](https://github.com/dnerochain/dnero-eth-rpc-adaptor#rpc-apis).

### Setup

First, install **Go 1.14.2** and set environment variables `GOPATH` , `GOBIN`, and `PATH`. Next, clone the Dnero blockchain repo and install Dnero following the steps below:

```
mkdir -p $GOPATH/src/github.com/dnerochain 
cd $GOPATH/src/github.com/dnerochain
git clone https://github.com/dnerochain/dnero-protocol-ledger.git $GOPATH/src/github.com/dnerochain/dnero
cd dnero
git checkout privatenet
export GO111MODULE=on
make install
```

Next, clone the `dnero-eth-rpc-adaptor` repo:

```
cd $GOPATH/src/github.com/dnerochain
git clone https://github.com/dnerochain/dnero-eth-rpc-adaptor
```

### Build and Install

#### Build the binary under macOS or Linux
Following the steps below to build the `dnero-eth-rpc-adaptor` binary and copy it into your `$GOPATH/bin`.

```
export DNERO_ETH_RPC_ADAPTOR_HOME=$GOPATH/src/github.com/dnerochain/dnero-eth-rpc-adaptor
cd $DNERO_ETH_RPC_ADAPTOR_HOME
export GO111MODULE=on
make install
```

#### Cross compilation for Windows
On a macOS machine, the following command should build the `dnero-eth-rpc-adaptor.exe` binary under `build/windows/`

```
make windows
```

#### Run the Adaptor with a local Dnero private testnet

First, run a private testnet Dnero node with its RPC port opened at 15511:

```
cd $DNERO_HOME
cp -r ./integration/privatenet ../privatenet
mkdir ~/.dnerocli
cp -r ./integration/privatenet/dnerocli/* ~/.dnerocli/
chmod 700 ~/.dnerocli/keys/encrypted

dnero start --config=../privatenet/node_eth_rpc --password=qwertyuiop
```

Then, open another terminal, create the config folder for the RPC adaptor

```
export DNERO_ETH_RPC_ADAPTOR_HOME=$GOPATH/src/github.com/dnerochain/dnero-eth-rpc-adaptor
cd $DNERO_ETH_RPC_ADAPTOR_HOME
mkdir -p ../privatenet/eth-rpc-adaptor
```

Use your favorite editor to open file `../privatenet/eth-rpc-adaptor/config.yaml`, paste in the follow content, save and close the file:

```
dnero:
  rpcEndpoint: "http://127.0.0.1:15511/rpc"
rpc:
  enabled: true
  httpAddress: "0.0.0.0"
  httpPort: 15444
  wsAddress: "0.0.0.0"
  wsPort: 15445
  timeoutSecs: 600 
  maxConnections: 2048
log:
  levels: "*:debug"
```

Then, launch the adaptor binary with the following command:

```
cd $DNERO_ETH_RPC_ADAPTOR_HOME
dnero-eth-rpc-adaptor start --config=../privatenet/eth-rpc-adaptor
```

The RPC adaptor will first create 10 test wallets, which will be useful for running tests with dev tools like Truffle, Hardhat. After the test wallets are created, the ETH RPC APIs will be ready for use.

### RPC APIs

The RPC APIs should conform to the Ethereum JSON RPC API standard: https://eth.wiki/json-rpc/API. We currently support the following Ethereum RPC APIs:

```
eth_chainId
eth_syncing
eth_accounts
eth_protocolVersion
eth_getBlockByHash
eth_getBlockByNumber
eth_blockNumber
eth_getUncleByBlockHashAndIndex
eth_getTransactionByHash
eth_getTransactionByBlockNumberAndIndex
eth_getTransactionByBlockHashAndIndex
eth_getBlockTransactionCountByHash
eth_getTransactionReceipt
eth_getBalance
eth_getStorageAt
eth_getCode
eth_getTransactionCount
eth_getLogs
eth_getBlockTransactionCountByNumber
eth_call
eth_gasPrice
eth_estimateGas
eth_sendRawTransaction
eth_sendTransaction
net_version
web3_clientVersion
```

The following examples demonstrate how to interact with the RPC APIs using the `curl` command:

```
# Query Chain ID
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":67}' http://localhost:18888/rpc

# Query synchronization status
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' http://localhost:18888/rpc

# Query block number
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:18888/rpc

# Query account DToken balance (should return an integer which represents the current DToken balance in wei)
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2E833968E5bB786Ae419c4d13189fB081Cc43bab", "latest"],"id":1}' http://localhost:18888/rpc
```
