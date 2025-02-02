package handler

import "fmt"

const Version string = "Example API v0"

type ChainInfo struct {
	RPC  string
	ID   string
	Name string
}

var SupportedChains = map[string]ChainInfo{
	"1": {
		RPC:  "https://eth.llamarpc.com",
		ID:   "01",
		Name: "Ethereum Mainnet",
	},
	"137": {
		RPC:  "https",
		ID:   "89",
		Name: "Polygon Mainnet",
	},
	"56": {
		RPC:  "https://bsc-rpc.publicnode.com",
		ID:   "38",
		Name: "Binance Smart Chain",
	},
}

func GetChainInfo(chainId string) (ChainInfo, error) {
	chain, exists := SupportedChains[chainId]
	if !exists {
		return ChainInfo{}, fmt.Errorf("chain ID %v not supported", chainId)
	}
	return chain, nil
}
