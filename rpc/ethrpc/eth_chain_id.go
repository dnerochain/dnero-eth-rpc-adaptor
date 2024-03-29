package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dnerochain/dnero-eth-rpc-adaptor/common"
	hexutil "github.com/dnerochain/dnero/common/hexutil"
	"github.com/dnerochain/dnero/ledger/types"
	trpc "github.com/dnerochain/dnero/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type chainIDResultWrapper struct {
	chainID string
}

// ------------------------------- eth_chainId -----------------------------------

func (e *EthRPCService) ChainId(ctx context.Context) (result string, err error) {
	logger.Infof("eth_chainId called")

	client := rpcc.NewRPCClient(common.GetDneroRPCEndpoint())
	rpcRes, rpcErr := client.Call("dnero.GetStatus", trpc.GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, err := common.HandleDneroRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	dneroChainIDResult, ok := resultIntf.(chainIDResultWrapper)
	if !ok {
		return "", fmt.Errorf("failed to convert chainIDResultWrapper")
	}

	dneroChainID := dneroChainIDResult.chainID
	ethChainID := types.MapChainID(dneroChainID, blockHeight).Uint64()
	result = hexutil.EncodeUint64(ethChainID)

	return result, nil
}
