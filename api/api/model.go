package handler

type GetEvmContractExtCodeSizeRequestResponse struct {
	ChainId string `json:"chain-id"`
	Address string `json:"contract-address"`
	Size    string `json:"contract-size"`
}

type GetEvmContractCodeRequestResponse struct {
	ChainId string `json:"chain-id"`
	Address string `json:"contract-address"`
	Size    string `json:"contract-size"`
	Code    string `json:"contract-code"`
}

type GetEvmContractDataAtMemoryRequestResponse struct {
	ChainId string `json:"chain-id"`
	Address string `json:"contract-address"`
	Bytes   string `json:"bytes"`
}

type GetEvmContractCallViewRequestResponse struct {
	ChainId    string `json:"chain-id"`
	Address    string `json:"contract-address"`
	MethodName string `json:"method-name"`
	Response   string `json:"response"`
}

type GetEvmContractBalanceRequestResponse struct {
	ChainId string `query:"chain-id"`
	Address string `query:"address"`
	Balance string `query:"balance"`
}
