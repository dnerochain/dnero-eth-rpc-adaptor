package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/dnerochain/dnero-eth-rpc/common"

	trpc "github.com/dnerochain/dnero/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_protocolVersion -----------------------------------

func (e *EthRPCService) ProtocolVersion(ctx context.Context) (result string, err error) {
	logger.Infof("eth_protocolVersion called")

	client := rpcc.NewRPCClient(common.GetDneroRPCEndpoint())
	rpcRes, rpcErr := client.Call("dnero.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetVersionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Version, nil
	}

	resultIntf, err := common.HandleDneroRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = resultIntf.(string)

	return result, nil
}
