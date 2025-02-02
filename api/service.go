package handler

import (
	"api/pkg/utils"
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func GetVersionRequest(r *http.Request, parameters ...interface{}) (interface{}, error) {
	return utils.VersionResponse{
		Version: Version,
	}, nil
}

func GetEvmContractExtCodeSizeRequest(r *http.Request, parameters ...*GetEvmContractExtCodeSizeRequestParams) (interface{}, error) {
	var params *GetEvmContractExtCodeSizeRequestParams
	var client *ethclient.Client

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetEvmContractExtCodeSizeRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}
	if params.JsonRpc != "" {
		client_, err := DialClient(params.JsonRpc)
		if err != nil {
			err_ := fmt.Errorf("dial client %v failed: %v", params.JsonRpc, err.Error())
			logrus.Error(err_)
			return nil, err_
		}
		client = client_
	} else {
		chainInfo, err := GetChainInfo(params.ChainId)
		if err != nil {
			return nil, err
		}
		client_, err := DialClient(chainInfo.RPC)
		if err != nil {
			err_ := fmt.Errorf("dial client %v failed: %v", chainInfo.RPC, err.Error())
			logrus.Error(err_)
			return nil, err_
		}
		client = client_
	}

	if ok := common.IsHexAddress(params.Address); !ok {
		err_ := fmt.Errorf("contract address is not hex")
		logrus.Error(err_)
		return nil, err_
	}

	_, extCodeSize_, err := ExtCodeSize(client, common.HexToAddress(params.Address))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &GetEvmContractExtCodeSizeRequestResponse{
		ChainId: params.ChainId,
		Address: params.Address,
		Size:    fmt.Sprintf("%+v", extCodeSize_),
	}, nil
}

func GetEvmContractCodeRequest(r *http.Request, parameters ...*GetEvmContractCodeRequestParams) (interface{}, error) {
	var params *GetEvmContractCodeRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetEvmContractCodeRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	chainInfo, err := GetChainInfo(params.ChainId)
	if err != nil {
		return nil, err
	}
	client, err := DialClient(chainInfo.RPC)
	if err != nil {
		err_ := fmt.Errorf("dial client %v failed: %v", chainInfo.RPC, err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	extCode_, extCodeSize_, err := ExtCodeSize(client, common.HexToAddress(params.Address))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &GetEvmContractCodeRequestResponse{
		ChainId: params.ChainId,
		Address: params.Address,
		Size:    fmt.Sprintf("%+v", extCodeSize_),
		Code:    hex.EncodeToString(extCode_),
	}, nil
}

func GetEvmContractDataAtMemoryRequest(r *http.Request, parameters ...*GetEvmContractDataAtMemoryRequestParams) (interface{}, error) {
	var params *GetEvmContractDataAtMemoryRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetEvmContractDataAtMemoryRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	chainInfo, err := GetChainInfo(params.ChainId)
	if err != nil {
		return nil, err
	}
	client, err := DialClient(chainInfo.RPC)
	if err != nil {
		err_ := fmt.Errorf("dial client %v failed: %v", chainInfo.RPC, err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	if ok := common.IsHexAddress(params.Address); !ok {
		err_ := fmt.Errorf("contract address is not hex")
		logrus.Error(err_)
		return nil, err_
	}

	slot, err := strconv.ParseInt(params.StorgeAt, 10, 64)
	if err != nil {
		err_ := fmt.Errorf("could not parse slot: %v", err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	data, err := GetStorageAt(client, common.HexToAddress(params.Address), slot)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	return &GetEvmContractDataAtMemoryRequestResponse{
		ChainId: params.ChainId,
		Address: params.Address,
		Bytes:   hex.EncodeToString(data),
	}, nil
}

func GetEvmContractCallViewRequest(r *http.Request, parameters ...*GetEvmContractCallViewRequestParams) (interface{}, error) {
	var params *GetEvmContractCallViewRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetEvmContractCallViewRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	count := 0
	for key := range r.URL.Query() {
		if strings.HasPrefix(key, "method-inputs") {
			count++
		}
	}

	if count > 0 {
		// Loop through and parse the parameters
		for i := 0; i < count; i++ {
			paramType := r.URL.Query().Get(fmt.Sprintf("method-inputs[%d][type]", i))
			paramValue := r.URL.Query().Get(fmt.Sprintf("method-inputs[%d][value]", i))

			if paramType != "" && paramValue != "" {
				params.MethodParams = append(params.MethodParams, utils.Parameter{
					Type:  paramType,
					Value: paramValue,
				})
			}
		}
		// Assign methodParams to your params object
		//params.MethodParams = methodParams
	}

	fmt.Printf("\n paramters input: %+v", params)

	chainInfo, err := GetChainInfo(params.ChainId)
	if err != nil {
		return nil, err
	}
	client, err := DialClient(chainInfo.RPC)
	if err != nil {
		err_ := fmt.Errorf("dial client %v failed: %v", chainInfo.RPC, err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	if ok := common.IsHexAddress(params.Address); !ok {
		err_ := fmt.Errorf("contract address is not hex")
		logrus.Error(err_)
		return nil, err_
	}

	result, err := CallContract(client, common.HexToAddress(params.Address), params.MethodName, params.MethodParams)
	if err != nil {
		err_ := fmt.Errorf("failed to call contract %v: %w", params.Address, err)
		logrus.Error(err_)
		return nil, err_
	}
	return &GetEvmContractCallViewRequestResponse{
		ChainId:    params.ChainId,
		Address:    params.Address,
		MethodName: params.MethodName,
		Response:   hex.EncodeToString(result),
	}, nil
}

func GetEvmContractBalanceRequest(r *http.Request, parameters ...*GetEvmContractBalanceRequestParams) (interface{}, error) {
	var params *GetEvmContractBalanceRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetEvmContractBalanceRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	chainInfo, err := GetChainInfo(params.ChainId)
	if err != nil {
		return nil, err
	}
	client, err := DialClient(chainInfo.RPC)
	if err != nil {
		err_ := fmt.Errorf("dial client %v failed: %v", chainInfo.RPC, err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	address := common.HexToAddress(params.Address)
	balance, err := client.BalanceAt(context.Background(), address, nil) // nil for latest block
	if err != nil {
		err_ := fmt.Errorf("get balance failed: %v", err.Error())
		logrus.Error(err_)
		return nil, err_
	}

	return &GetEvmContractBalanceRequestResponse{
		ChainId: params.ChainId,
		Address: params.Address,
		Balance: balance.String(),
	}, nil
}
