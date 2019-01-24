package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type AccountsFilterHandler struct {
	storage Storage
}

func ArgsToMap(args *fasthttp.Args) (query map[string]string) {
	query = map[string]string{}
	args.VisitAll(func(k, v []byte) {
		query[string(k)] = string(v)
	})

	return
}

func (handler *AccountsFilterHandler) ServeHTTP(c *routing.Context) {
	// log.Println(request.URL.Query())
	accountsQuery, err := CreateAccountsQuery(ArgsToMap(c.URI().QueryArgs()))

	if err == NoLimitError || err == BadQueryError {
		// log.Println(err)
		c.SetStatusCode(400)
		c.Write([]byte("{}"))
		return
	}

	c.SetContentType("application/json")
	c.Response.Header.Set("Connection", "Keep-Alive")

	accounts, err := handler.storage.Find(&accountsQuery)
	if err != nil {
		log.Println(err)
		c.SetStatusCode(500)
		return
	}
	data := struct {
		Accounts []map[string]interface{} `json:"accounts"`
	}{
		accounts,
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
