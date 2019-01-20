package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("started!")

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	storage := &MongoStorage{
		Uri:      MongoUri,
		Database: "hl",
	}
	storage.Init()

	// Parse("/Users/rrabelkharisov/highloadcup/test_accounts_291218/data/data.zip", storage, true)
	Parse(DataFile, storage, false)

	http.Handle("/accounts/filter", &AccountsFilterHandler{storage})
	// http.ListenAndServe(":8000", nil)
	http.ListenAndServe(":80", nil)
}
