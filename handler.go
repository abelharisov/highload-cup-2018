package main

import routing "github.com/qiangxue/fasthttp-routing"

type Handler interface {
	ServeHTTP(c *routing.Context) error
}
