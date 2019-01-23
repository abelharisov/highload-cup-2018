package main

import (
	"encoding/json"
	"fmt"
	"log"

	routing "github.com/qiangxue/fasthttp-routing"
)

type AccountsGroupHandler struct {
	storage Storage
}

func (handler *AccountsGroupHandler) ServeHTTP(c *routing.Context) {
	// log.Println(request.URL.Query())
	query, err := CreateAccountsGroupQuery(ArgsToMap(c.URI().QueryArgs()))

	if err == NoLimitError || err == BadQueryError {
		// log.Println(err)
		c.SetStatusCode(400)
		c.Write([]byte("{}"))
		return
	}

	c.SetContentType("application/json")
	c.Response.Header.Set("Connection", "Keep-Alive")

	data, err := handler.storage.Group(&query)
	if err != nil {
		log.Println(err)
		c.SetStatusCode(500)
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
	c.Response.Header.Set("Content-Length", fmt.Sprint(len(body)))
	c.SetStatusCode(200)

	if err != nil {
		log.Println(err)
		c.SetStatusCode(500)
		return
	}

	_, err = c.Write(body)
	if err != nil {
		log.Println(err)
		c.SetStatusCode(500)
		return
	}
}
