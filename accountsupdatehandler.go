package main

import routing "github.com/qiangxue/fasthttp-routing"

type AccountsUpdateHandler struct {
	storage Storage
}

func (handler *AccountsUpdateHandler) ServeHTTP(c *routing.Context) error {
	return &Error{400, "not implemented"}
}
