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

	if err == NoLimitError || err == BadQueryError {
		response.WriteHeader(400)
		response.Write([]byte("{}"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Connection", "Keep-Alive")

	accounts := handler.storage.Find(&accountsQuery)
	data := struct {
		Accounts []map[string]interface{} `json:"accounts"`
	}{
		accounts,
	}

	body, err := json.Marshal(data)
	response.Header().Set("Content-Length", fmt.Sprint(len(body)))
	response.WriteHeader(200)

	if err != nil {
		response.WriteHeader(500)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		response.WriteHeader(500)
		return
	}
}
