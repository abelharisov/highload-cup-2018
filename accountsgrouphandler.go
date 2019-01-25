package main

import (
	"encoding/json"
	"fmt"

	"github.com/qiangxue/fasthttp-routing"
)

type AccountsGroupHandler struct {
	storage Storage
}

func (handler *AccountsGroupHandler) ServeHTTP(c *routing.Context) error {
	query, err := CreateAccountsGroupQuery(ArgsToMap(c.URI().QueryArgs()))
	if err != nil {
		return err
	}

	data, err := handler.storage.Group(&query)
	if err != nil {
		return err
	}

	var body []byte = []byte("{\"groups\":[]}")
	if data != nil {
		type M = map[string]interface{}
		type A = []interface{}

		formattedResponse := M{
			"groups": make(A, 0, len(data)),
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

		var err error
		body, err = json.Marshal(formattedResponse)
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
