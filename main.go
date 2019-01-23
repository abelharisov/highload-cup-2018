package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

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
		err := Parse(DataFile, storage, false)
		if err != nil {
			panic(err)
		}
		log.Println("Parsed", time.Now().Sub(start).Seconds())
	}()

	router := routing.New()

	afh := &AccountsFilterHandler{storage}
	router.Get("/accounts/filter/", func(c *routing.Context) error {
		afh.ServeHTTP(c)
		return nil
	})
	agh := &AccountsGroupHandler{storage}
	router.Get("/accounts/group/", func(c *routing.Context) error {
		agh.ServeHTTP(c)
		return nil
	})

	fasthttp.ListenAndServe(":80", router.HandleRequest)
}
