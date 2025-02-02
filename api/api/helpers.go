package handler

import (
	"context"
	"encoding/hex"
	"fmt"
	utils "generic-evm-api-go/api/pkg/utils"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func DialClient(jsonrpc string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(jsonrpc)
	if err != nil {
		err_ := fmt.Errorf("client connection failed: %v", err)
		logrus.Error(err_.Error())
		return nil, err_
	}
	return client, nil
}

func ViewFunction(client *ethclient.Client, contractAddress common.Address, parsedABI abi.ABI, methodName string, args ...interface{}) ([]byte, error) {
	data, err := parsedABI.Pack(methodName, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{To: &contractAddress, Data: data}
	result, err := client.CallContract(context.Background(), callMsg, big.NewInt(305965178))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetCallBytes(parsedABI abi.ABI, methodName string, args ...interface{}) ([]byte, error) {
	isArgsEmpty := func(args []interface{}) bool {
		if len(args) == 0 {
			return true
		}

		for _, arg := range args {
			if arg != nil {
				return false
			}
		}
		return true
	}

	var data []byte
	var err error
	if !isArgsEmpty(args) {
		data, err = parsedABI.Pack(methodName, args...)
	} else {
		data, err = parsedABI.Pack(methodName)
	}

	return data, err
}

func ExtCodeSize(client *ethclient.Client, address common.Address) ([]byte, int, error) {
	ctx := context.Background()
	code, err := client.CodeAt(ctx, address, nil) // nil block number for the latest state
	if err != nil {
		return nil, 0, fmt.Errorf("geth client failed to get extcodesize: %v", err)
	}
	return code, len(code), nil
}

func GetStorageAt(client *ethclient.Client, address common.Address, slot int64) ([]byte, error) {
	slotHash := common.BigToHash(common.Big1)
	if slot != 0 {
		slotHash = common.BigToHash(big.NewInt(int64(slot)))
	}

	storage, err := client.StorageAt(context.Background(), address, slotHash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %+v", err.Error())
	}

	return storage, nil
}

func ConstructCallData(methodName string, params []utils.Parameter) ([]byte, error) {
	// Create method signature
	methodSig := createMethodSignature(methodName, params)

	// Calculate function selector (first 4 bytes of the hash of the method signature)
	selector := crypto.Keccak256([]byte(methodSig))[:4]

	// Parse and pack parameters
	packedParams, err := packParameters(params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack parameters: %v", err)
	}

	// Combine selector with packed parameters
	callData := append(selector, packedParams...)

	return callData, nil
}

func createMethodSignature(methodName string, params []utils.Parameter) string {
	var types []string
	for _, param := range params {
		types = append(types, param.Type)
	}
	return fmt.Sprintf("%s(%s)", methodName, strings.Join(types, ","))
}

func packParameters(params []utils.Parameter) ([]byte, error) {
	if len(params) == 0 {
		logrus.Error("param count 0")
		return []byte{}, nil
	}

	var arguments abi.Arguments
	var values []interface{}

	for _, param := range params {
		abiType, err := abi.NewType(param.Type, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %s: %v", param.Type, err)
		}
		arguments = append(arguments, abi.Argument{Type: abiType})

		value, err := parseParameterValue(param.Type, param.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameter value: %v", err)
		}
		values = append(values, value)
	}

	logrus.Error(values)
	packed, err := arguments.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack values: %v", err)
	}

	return packed, nil
}

func parseParameterValue(paramType, paramValue string) (interface{}, error) {
	switch paramType {
	case "uint256":
		val, ok := new(big.Int).SetString(paramValue, 10)
		if !ok {
			return nil, fmt.Errorf("invalid uint256: %s", paramValue)
		}
		return val, nil

	case "address":
		return common.HexToAddress(paramValue), nil

	case "bytes":
		if !strings.HasPrefix(paramValue, "0x") {
			return nil, fmt.Errorf("bytes value must start with 0x")
		}
		data, err := hex.DecodeString(paramValue[2:])
		if err != nil {
			return nil, fmt.Errorf("invalid bytes: %s", paramValue)
		}
		return data, nil

	case "bool":
		return paramValue == "true", nil

	case "string":
		return paramValue, nil

	default:
		return nil, fmt.Errorf("unsupported parameter type: %s", paramType)
	}
}

func CallContract(
	client *ethclient.Client,
	contractAddress common.Address,
	methodName string,
	params []utils.Parameter,
) ([]byte, error) {
	fmt.Printf("\n current params in callcontract: %v", params)
	callData, err := ConstructCallData(methodName, params)
	if err != nil {
		return nil, fmt.Errorf("failed to construct call data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %v", err)
	}

	return result, nil
}
