package main

import (
	"encoding/json"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/qiangxue/fasthttp-routing"
)

type AccountsFilterHandler struct {
	storage Storage
}

func (handler *AccountsFilterHandler) ServeHTTP(c *routing.Context) error {
	accountsQuery, err := CreateAccountsQuery(ArgsToMap(c.URI().QueryArgs()))
	if err != nil {
		return err
	}

	accounts, err := handler.storage.Find(&accountsQuery)
	if err != nil {
		return err
	}

	var body []byte = []byte("{\"accounts\":[]}")
	if accounts != nil {
		data := bson.M{
			"accounts": accounts,
		}
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return &Error{500, err.Error()}
		}
	}

	c.Response.Header.Set("Content-Length", fmt.Sprint(len(body)))
	_, err = c.Write(body)
	if err != nil {
		return &Error{500, err.Error()}
	}

	return nil
}
