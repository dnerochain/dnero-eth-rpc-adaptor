package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/dnerochain/dnero/cmd/dnerocli/cmd/utils"
	tcommon "github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/crypto"
	"github.com/dnerochain/dnero/ledger/types"
	trpc "github.com/dnerochain/dnero/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "common"})

type AddressBook map[string]*crypto.PrivateKey

var TestWallets AddressBook = make(AddressBook)
var TestWalletArr []string

func GetDneroRPCEndpoint() string {
	dneroRPCEndpoint := viper.GetString(CfgDneroRPCEndpoint)
	return dneroRPCEndpoint
}

func HandleDneroRPCResponse(rpcRes *rpcc.RPCResponse, rpcErr error, parse func(jsonBytes []byte) (interface{}, error)) (result interface{}, err error) {
	if rpcErr != nil {
		return nil, fmt.Errorf("failed to get dnero RPC response: %v", rpcErr)
	}
	if rpcRes.Error != nil {
		return nil, fmt.Errorf("dnero RPC returns an error: %v", rpcRes.Error)
	}

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse dnero RPC response: %v, %s", err, string(jsonBytes))
	}

	//logger.Infof("HandleDneroRPCResponse, jsonBytes: %v", string(jsonBytes))
	result, err = parse(jsonBytes)
	if err != nil {
		logger.Warnf("Failed to parse dnero RPC response: %v, %s", err, string(jsonBytes))
	}
	return
}

func GetHeightByTag(tag string) (height tcommon.JSONUint64) {
	switch tag {
	case "latest":
		height = tcommon.JSONUint64(math.MaxUint64)
	case "earliest":
		height = tcommon.JSONUint64(1)
	case "pending":
		height = tcommon.JSONUint64(math.MaxUint64)
	default:
		height = tcommon.JSONUint64(Str2hex2unit(tag))
	}
	return height
}

func HexStrToBigInt(intHexStr string) *big.Int {
	// remove 0x suffix if found in the input string
	if strings.HasPrefix(intHexStr, "0x") {
		intHexStr = strings.TrimPrefix(intHexStr, "0x")
	}

	val := new(big.Int)
	val.SetString(intHexStr, 16)

	return val
}

func Str2hex2unit(str string) uint64 {
	// remove 0x suffix if found in the input string
	if strings.HasPrefix(str, "0x") {
		str = strings.TrimPrefix(str, "0x")
	}

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(str, 16, 64)
	return uint64(result)
}

func Int2hex2str(num int) string {
	return "0x" + strconv.FormatInt(int64(num), 16)
}

func HexStr2Uint64(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1) // remove 0x suffix if found in the input string
	result, _ := strconv.ParseUint(cleaned, 16, 64)  // base 16 for hexadecimal
	return uint64(result)
}

func HexToBytes(hexStr string) ([]byte, error) {
	trimmedHexStr := strings.TrimPrefix(hexStr, "0x")
	data, err := hex.DecodeString(trimmedHexStr)
	return data, err
}

func GenerateSctx(arg EthSmartContractArgObj) (result *types.SmartContractTx, err error) {
	sequence, seqErr := GetSeqByAddress(arg.From)
	if seqErr != nil {
		logger.Errorf("Failed to get sequence by address: %v\n", arg.From)
		if arg.From.String() != "0x0000000000000000000000000000000000000000" {
			return nil, seqErr
		}
		sequence = 1
	}

	from := types.TxInput{
		Address: arg.From, //tcommon.HexToAddress(arg.From.String()),
		Coins: types.Coins{
			DneroWei: new(big.Int).SetUint64(0),
			//DTokenWei: new(big.Int).SetUint64(Str2hex2unit(arg.Value)),
			DTokenWei: HexStrToBigInt(arg.Value),
		},
		Sequence: sequence,
	}

	to := types.TxOutput{
		Address: arg.To, //tcommon.HexToAddress(arg.To.String()),
	}

	gasPriceStr := "0wei"
	if arg.GasPrice != "" {
		gasPriceStr = arg.GasPrice + "wei"
	}

	gasPrice, ok := types.ParseCoinAmount(gasPriceStr)
	if !ok {
		err = errors.New("failed to parse gas price")
		logger.Errorf(fmt.Sprintf("%v", err))
		return nil, err
	}

	data, err := HexToBytes(arg.Data)
	if err != nil {
		logger.Errorf("Failed to decode data: %v, err: %v\n", arg.Data, err)
		return nil, err
	}

	gas := uint64(10000000)
	if arg.Gas != "" {
		gas = Str2hex2unit(arg.Gas)
	}
	fmt.Printf("gas: %v\n", gas)

	result = &types.SmartContractTx{
		From:     from,
		To:       to,
		GasLimit: gas,
		GasPrice: gasPrice,
		Data:     data,
	}
	return result, nil
}

func SignRawBytes(address string, rawBytes []byte) (result *crypto.Signature, err error) {
	priKey, ok := TestWallets[address]
	if !ok {
		return nil, fmt.Errorf("testAddress not found : %s, TestWallets %+v", address, TestWallets)
	}
	result, err = priKey.Sign(rawBytes)
	return
}

func GetSignedBytes(arg EthSmartContractArgObj, chainID string, blockNumber string) (string, error) {
	fromAddress := tcommon.HexToAddress(arg.From.String())
	sctx, _ := GenerateSctx(arg)
	sctxSignBytes := sctx.SignBytes(MapChainID(chainID, blockNumber))
	signature, err := SignRawBytes(strings.ToLower(arg.From.String()), sctxSignBytes)
	if err != nil {
		logger.Errorf("Failed to sign transaction: %v, err is %v\n", sctx, err)
		return "", err
	}

	sctx.SetSignature(fromAddress, signature)
	raw, err := types.TxToBytes(sctx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTXstr := hex.EncodeToString(raw)
	return signedTXstr, nil
}

func GetSctxBytes(arg EthSmartContractArgObj) (sctxBytes []byte, err error) {
	sctx, err := GenerateSctx(arg)
	if err != nil {
		logger.Errorf("Failed to generate smart contract transaction: %v\n", sctx)
		return sctxBytes, err
	}
	sctxBytes, err = types.TxToBytes(sctx)
	if err != nil {
		logger.Errorf("Failed to encode smart contract transaction: %v\n", sctx)
		return sctxBytes, err
	}
	return sctxBytes, nil
}

func GetSeqByAddress(address tcommon.Address) (sequence uint64, err error) {
	client := rpcc.NewRPCClient(GetDneroRPCEndpoint())

	rpcRes, rpcErr := client.Call("dnero.GetAccount", trpc.GetAccountArgs{Address: address.String()})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Sequence, nil
	}

	resultIntf, err := HandleDneroRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		//return sequence, err
		return 1, nil
	}
	sequence = resultIntf.(uint64) + 1

	return sequence, nil
}

func GetCurrentHeight() (height tcommon.JSONUint64, err error) {
	client := rpcc.NewRPCClient(GetDneroRPCEndpoint())
	rpcRes, rpcErr := client.Call("dnero.GetStatus", trpc.GetStatusArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.LatestFinalizedBlockHeight, nil
	}

	resultIntf, err := HandleDneroRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return height, err
	}
	height = resultIntf.(tcommon.JSONUint64)
	return height, nil
}

func MapChainID(chainIDStr string, blockNumber string) string {
	if HexStr2Uint64(blockNumber) < tcommon.HeightRPCCompatibility {
		return mapChainIDWithoutOffset(chainIDStr)
	}
	if chainIDStr == "0x155a" { // ChainID: 5466
		return "mainnet"
	} else if chainIDStr == "0x155b" { // MetaMask mainnet ChainID: 5467 
		return "testnet" // Assigned To MetaMask mainnet - Support DTOKEN Only
	} else if chainIDStr == "0x155c" { // correspond to Ropsten
		return "testnet_sapphire"
	} else if chainIDStr == "0x155d" { // correspond to Rinkeby
		return "testnet_amber"
	} else if chainIDStr == "0x155e" {
		return "privatenet"
	}
	return "" // all other chainIDs
}

func mapChainIDWithoutOffset(chainIDStr string) string {
	if chainIDStr == "0x1" { // correspond to the Ethereum mainnet
		return "mainnet"
	} else if chainIDStr == "0x2" { 
		return "testnet"
	} else if chainIDStr == "0x3" { // correspond to Ropsten
		return "testnet_sapphire"
	} else if chainIDStr == "0x4" { // correspond to Rinkeby
		return "testnet_amber"
	} else if chainIDStr == "0x5" {
		return "privatenet"
	}
	return "" // all other chainIDs
}
