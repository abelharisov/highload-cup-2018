package main

import (
	"encoding/json"

	routing "github.com/qiangxue/fasthttp-routing"
)

type AccountsUpdateHandler struct {
	storage Storage
}

func (handler *AccountsUpdateHandler) ServeHTTP(c *routing.Context) error {
	body := c.PostBody()

	var rawAccount map[string]interface{}
	err := json.Unmarshal(body, &rawAccount)
	if err != nil {
		return &Error{400, "unmarshall error"}
	}
	for key, _ := range rawAccount {
		if _, ok := allowedFields[key]; !ok {
			return &Error{400, "bad field"}
		}
	}

	if v, ok := rawAccount["sex"]; !ok || (v != "f" && v != "m") {
		return &Error{400, "bad sex"}
	}

	if v, ok := rawAccount["status"]; !ok || ParseStatus(v.(string)) == 0 {
		return &Error{400, "bad status"}
	}

	if _, ok := rawAccount["birth"]; !ok {
		return &Error{400, "bad birth"}
	}

	if _, ok := rawAccount["email"]; !ok {
		return &Error{400, "bad email"}
	}

	var account Account
	err = json.Unmarshal(body, &account)
	if err != nil {
		return &Error{400, "unmarshall error"}
	}

	EnrichAccount(&account, handler.storage.GetNow())

	return nil
}
