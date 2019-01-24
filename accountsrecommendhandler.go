package main

import (
	"encoding/json"
	"fmt"
	"log"

	routing "github.com/qiangxue/fasthttp-routing"
)

type AccountsRecommendHandler struct {
	storage Storage
}

func (handler *AccountsRecommendHandler) ServeHTTP(c *routing.Context) {
	// log.Println(request.URL.Query())
	query, err := CreateAccountsRecommendQuery(c.Param("id"), ArgsToMap(c.URI().QueryArgs()))

	if err == NoLimitError || err == BadQueryError {
		// log.Println(err)
		c.SetStatusCode(400)
		c.Write([]byte("{}"))
		return
	}

	c.SetContentType("application/json")
	c.Response.Header.Set("Connection", "Keep-Alive")

	recommends, err := handler.storage.(*MongoStorage).Recommend(&query, )
	if err != nil {
		log.Println(err)
		c.SetStatusCode(500)
		return
	}

	data := struct {
		Accounts []map[string]interface{} `json:"accounts"`
	}{
		recommends,
	}

	body, err := json.Marshal(data)
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
