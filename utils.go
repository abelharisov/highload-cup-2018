package main

import "github.com/valyala/fasthttp"

func ArgsToMap(args *fasthttp.Args) (query map[string]string) {
	query = map[string]string{}
	args.VisitAll(func(k, v []byte) {
		query[string(k)] = string(v)
	})

	return
}

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
