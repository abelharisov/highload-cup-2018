package main

import (
	"context"
	"log"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type MongoStorage struct {
	Uri      string
	Database string
	client   *mongo.Client
	accounts *mongo.Collection
}

func (storage *MongoStorage) Init() {
	var err error
	context, _ := context.WithTimeout(context.Background(), time.Second)
	storage.client, err = mongo.Connect(context, storage.Uri)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	storage.accounts = storage.client.Database(storage.Database).Collection("accounts")
}

func (storage *MongoStorage) LoadAccounts(accounts []Account) {
	context, _ := context.WithTimeout(context.Background(), time.Minute)
	documents := make([]interface{}, 0, len(accounts))
	for _, account := range accounts {
		documents = append(documents, interface{}(account))
	}
	storage.accounts.InsertMany(context, documents)
}

func (storage *MongoStorage) DropDatabase() {
	context, _ := context.WithTimeout(context.Background(), time.Minute)
	storage.client.Database(storage.Database).Drop(context)
}

func (storage *MongoStorage) Find(filters *AccountsFilter) (result []map[string]interface{}) {
	context, _ := context.WithTimeout(context.Background(), time.Minute)
	
	query := make(map[string]interface{})
	for _, filter := range filters.Filters {
		if filter.Operation == "eq" {
			query[filter.Field] = filter.Argument
		}
	}
	options := options.Find()
	options.SetSort(struct{id int32}{-1})
	cursor, err := storage.accounts.Find(context, query, options)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		account := make(map[string]interface{})
		cursor.Decode(&account)
		result = append(result, account)
	}

	return
}
