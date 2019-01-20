package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("started!")

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
