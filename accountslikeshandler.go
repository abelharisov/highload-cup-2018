package main

import routing "github.com/qiangxue/fasthttp-routing"

type AccountsLikesHandler struct {
	storage Storage
}

func (handler *AccountsLikesHandler) ServeHTTP(c *routing.Context) error {
	return &Error{400, "not implemented"}
}
