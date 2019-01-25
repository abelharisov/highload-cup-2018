package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func handle(handler Handler, c *routing.Context) {
	c.SetContentType("application/json")
	c.Response.Header.Set("Connection", "Keep-Alive")

	err := handler.ServeHTTP(c)

	if err != nil {
		code := err.(*Error).Code
		c.SetStatusCode(code)

		if code == 400 || code == 404 {
			c.WriteString("{}")
		}

		if code == 500 {
			log.Println(err.Error())
		}
	} else {
		c.SetStatusCode(200)
	}
}

func main() {
	start := time.Now()
	log.Println("started!")

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Println("caught sig: ", sig)
		os.Exit(0)
	}()

	storage := &MongoStorage{
		Uri:      MongoUri,
		Database: "hl",
	}
	storage.Init()

	go func() {
		err := Parse(DataFile, OptionsFile, storage, false)
		if err != nil {
			panic(err)
		}
		log.Println("Parsed", time.Now().Sub(start).Seconds())
	}()

	router := routing.New()

	afh := &AccountsFilterHandler{storage}
	router.Get("/accounts/filter/", func(c *routing.Context) error {
		handle(afh, c)
		return nil
	})

	agh := &AccountsGroupHandler{storage}
	router.Get("/accounts/group/", func(c *routing.Context) error {
		handle(agh, c)
		return nil
	})

	arh := &AccountsRecommendHandler{storage}
	router.Get("/accounts/<id>/recommend/", func(c *routing.Context) error {
		handle(arh, c)
		return nil
	})

	fasthttp.ListenAndServe(":80", router.HandleRequest)
}
