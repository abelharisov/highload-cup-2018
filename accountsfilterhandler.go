package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountsFilterHandler struct {
	storage Storage
}

func (handler *AccountsFilterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	accountsQuery, err := CreateAccountsQuery(request.URL.Query())

	if err == NoLimitError {
		response.WriteHeader(400)
		return
	}

	fmt.Fprint(response, "{\"accounts\":")
	encoder := json.NewEncoder(response)
	accounts := handler.storage.Find(&accountsQuery)
	encoder.Encode(accounts)
	fmt.Fprint(response, "}")
}
