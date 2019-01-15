package main

import (
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("started!")
	storage := &MongoStorage{
		Uri: MongoUri,
		Database: "hl",
	}
	storage.Init()

	Parse(DataFile, storage)

	http.Handle("accounts/filter", &AccountsFilterHandler{})
	http.ListenAndServe("localhost:80", nil)
}
