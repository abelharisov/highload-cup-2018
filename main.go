package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	http.Handle("/accounts/filter/", &AccountsFilterHandler{storage})
	http.Handle("/accounts/group/", &AccountsGroupHandler{storage})
	http.ListenAndServe(":80", nil)
}
