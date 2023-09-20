package ethrpc

import (
	"context"
	"encoding/json"
	"math"
	"math/big"

	"github.com/dnerochain/dnero-eth-rpc-adaptor/common"
	"github.com/dnerochain/dnero/ledger/types"
	trpc "github.com/dnerochain/dnero/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBalance -----------------------------------

func (e *EthRPCService) GetBalance(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getBalance called")

	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the dnero.GetAccount method
	}

	client := rpcc.NewRPCClient(common.GetDneroRPCEndpoint())
	rpcRes, rpcErr := client.Call("dnero.GetAccount", trpc.GetAccountArgs{Address: address, Height: height})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Balance.DTokenWei, nil
	}

	resultIntf, err := common.HandleDneroRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return "0x0", nil
	}

	// result = fmt.Sprintf("0x%x", resultIntf.(*big.Int))
	result = "0x" + (resultIntf.(*big.Int)).Text(16)

	return result, nil
}
