package ethrpc

import (
	"context"
	"math"
	"math/big"

	"github.com/dnerochain/dnero-eth-rpc/common"

	hexutil "github.com/dnerochain/dnero/common/hexutil"
	trpc "github.com/dnerochain/dnero/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockTransactionCountByNumber -----------------------------------
func (e *EthRPCService) GetBlockTransactionCountByNumber(ctx context.Context, numberStr string) (result hexutil.Uint64, err error) {
	logger.Infof("eth_getBlockTransactionCountByNumber called")
	height := common.GetHeightByTag(numberStr)
	if height == math.MaxUint64 {
		height, err = common.GetCurrentHeight()
		if err != nil {
			return result, err
		}
	}

	chainIDStr, err := e.ChainId(ctx)
	if err != nil {
		logger.Errorf("Failed to get chainID\n")
		return result, nil
	}
	chainID := new(big.Int)
	chainID.SetString(chainIDStr, 16)

	client := rpcc.NewRPCClient(common.GetDneroRPCEndpoint())
	rpcRes, rpcErr := client.Call("dnero.GetBlockByHeight", trpc.GetBlockByHeightArgs{
		Height: height})
	block, err := GetBlockFromTRPCResult(chainID, rpcRes, rpcErr, false)
	return hexutil.Uint64(len(block.Transactions)), err
}
