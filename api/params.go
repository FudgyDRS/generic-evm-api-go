package handler

import utils "api/pkg/utils"

type GetEvmContractExtCodeSizeRequestParams struct {
	ChainId string `query:"chain-id"`
	JsonRpc string `query:"json-rpc" optional:"true"`
	Address string `query:"contract-address"`
}

type GetEvmContractCodeRequestParams struct {
	ChainId string `query:"chain-id"`
	JsonRpc string `query:"json-rpc" optional:"true"`
	Address string `query:"contract-address"`
}

type GetEvmContractDataAtMemoryRequestParams struct {
	ChainId  string `query:"chain-id"`
	JsonRpc  string `query:"json-rpc" optional:"true"`
	Address  string `query:"contract-address"`
	StorgeAt string `query:"storage-at"`
}

type Parameter struct {
	Type  string `query:"type" optional:"true"`
	Value string `query:"value" optional:"true"`
}

type GetEvmContractCallViewRequestParams struct {
	ChainId      string            `query:"chain-id"`
	JsonRpc      string            `query:"json-rpc" optional:"true"`
	Address      string            `query:"contract-address"`
	MethodName   string            `query:"method-name" optional:"true"`
	MethodParams []utils.Parameter `query:"method-inputs" optional:"true"` // {type, data}
}

type GetEvmContractBalanceRequestParams struct {
	ChainId string `query:"chain-id"`
	JsonRpc string `query:"json-rpc" optional:"true"`
	Address string `query:"address"`
}
