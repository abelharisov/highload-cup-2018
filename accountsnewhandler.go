package main

import (
	"encoding/json"

	routing "github.com/qiangxue/fasthttp-routing"
)

type AccountsNewHandler struct {
	storage Storage
}

var allowedFields = map[string]int{
	"id":        1,
	"fname":     1,
	"sname":     1,
	"phone":     1,
	"email":     1,
	"sex":       1,
	"birth":     1,
	"country":   1,
	"city":      1,
	"joined":    1,
	"interests": 1,
	"status":    1,
	"premium":   1,
	"likes":     1,
}

func (handler *AccountsNewHandler) ServeHTTP(c *routing.Context) error {
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

	err = handler.storage.LoadAccounts([]Account{account})
	if err != nil {
		return err
	}

	return nil
}
