package main

import (
	"encoding/json"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/qiangxue/fasthttp-routing"
)

type AccountsRecommendHandler struct {
	storage Storage
}

func (handler *AccountsRecommendHandler) ServeHTTP(c *routing.Context) error {
	id := c.Param("id")
	query, err := CreateAccountsRecommendQuery(id, ArgsToMap(c.URI().QueryArgs()))
	if err != nil {
		return err
	}

	recommends, err := handler.storage.Recommend(&query)
	if err != nil {
		return err
	}

	var body []byte = []byte("{\"accounts\":[]}")
	if recommends != nil {
		data := bson.M{
			"accounts": recommends,
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
