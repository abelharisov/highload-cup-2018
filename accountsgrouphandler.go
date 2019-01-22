package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AccountsGroupHandler struct {
	storage Storage
}

func (handler *AccountsGroupHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// log.Println(request.URL.Query())
	query, err := CreateAccountsGroupQuery(request.URL.Query())

	if err == NoLimitError || err == BadQueryError {
		// log.Println(err)
		response.WriteHeader(400)
		response.Write([]byte("{}"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Connection", "Keep-Alive")

	data, err := handler.storage.Group(&query)
	if err != nil {
		log.Println(err)
		response.WriteHeader(500)
		return
	}

	type M = map[string]interface{}
	type A = []interface{}

	formattedResponse := M{
		"groups": A{},
	}
	for _, d := range data {
		group := d["_id"]
		for k, v := range group.(M) {
			if v == nil {
				delete(group.(M), k)
			}
		}
		group.(M)["count"] = d["count"]
		formattedResponse["groups"] = append(formattedResponse["groups"].(A), group)
	}

	body, err := json.Marshal(formattedResponse)
	response.Header().Set("Content-Length", fmt.Sprint(len(body)))
	response.WriteHeader(200)

	if err != nil {
		log.Println(err)
		response.WriteHeader(500)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		log.Println(err)
		response.WriteHeader(500)
		return
	}
}
