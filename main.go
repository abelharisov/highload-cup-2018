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

	Parse(DataFile, storage, false)

	http.Handle("accounts/filter", &AccountsFilterHandler{})
	http.ListenAndServe("localhost:80", nil)
}
