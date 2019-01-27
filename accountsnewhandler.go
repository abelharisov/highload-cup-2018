package main

import routing "github.com/qiangxue/fasthttp-routing"

type AccountsNewHandler struct {
	storage Storage
}

func (handler *AccountsNewHandler) ServeHTTP(c *routing.Context) error {
	return &Error{400, "not implemented"}
}
