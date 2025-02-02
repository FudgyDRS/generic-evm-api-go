package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "generic-evm-api-go/api/pkg/utils"

	"github.com/sirupsen/logrus"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			logrus.Error(fmt.Sprintf("Recovered from panic: %v", rec))

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	handlerWithCORS := utils.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var response interface{}
		var err error

		w.Header().Set("Content-Type", "application/json")
		switch query.Get("query") {
		case "version":
			response, err = GetVersionRequest(r)
			HandleResponse(w, r, response, err)
			return
		case "evm-contract-ext-code-size":
			response, err = GetEvmContractExtCodeSizeRequest(r)
			HandleResponse(w, r, response, err)
			return
		case "evm-contract-code":
			response, err = GetEvmContractCodeRequest(r)
			HandleResponse(w, r, response, err)
			return
		case "evm-contract-data-at-memory":
			response, err = GetEvmContractDataAtMemoryRequest(r)
			HandleResponse(w, r, response, err)
			return
		case "evm-contract-call-view":
			response, err = GetEvmContractCallViewRequest(r)
			HandleResponse(w, r, response, err)
			return
		case "get-contract-balance":
			response, err = GetEvmContractBalanceRequest(r)
			HandleResponse(w, r, response, err)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(utils.ErrMalformedRequest("Invalid query parameter"))
			return
		}
	}))

	handlerWithCORS.ServeHTTP(w, r)
}

func HandleResponse(w http.ResponseWriter, r *http.Request, response interface{}, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// shoudl json stringify
	logrus.Info(response)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
}
