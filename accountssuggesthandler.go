package main

import (
	"encoding/json"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	routing "github.com/qiangxue/fasthttp-routing"
)

type AccountsSuggestHandler struct {
	storage Storage
}

func (handler *AccountsSuggestHandler) ServeHTTP(c *routing.Context) error {
	query, err := CreateAccountsRecommendQuery(c.Param("id"), ArgsToMap(c.QueryArgs()))
	if err != nil {
		return err
	}

	suggests, err := handler.storage.Suggset(&query)
	if err != nil {
		return err
	}

	var body []byte = []byte("{\"accounts\":[]}")
	if suggests != nil {
		data := bson.M{
			"accounts": suggests,
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
