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
	accountsFilter := CreateAccountsFilter(request.URL.Query())

	fmt.Fprint(response, "{\"accounts\":")
	encoder := json.NewEncoder(response)
	accounts := handler.storage.Find(&accountsFilter)
	encoder.Encode(accounts)
	fmt.Fprint(response, "}")
}
