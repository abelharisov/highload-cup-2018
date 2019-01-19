package main

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
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

func (storage *MongoStorage) Find(query *AccountsQuery) (result []map[string]interface{}) {
	context, _ := context.WithTimeout(context.Background(), time.Minute)

	filters := make(map[string]interface{})
	projection := bson.M{
		"id":    1,
		"email": 1,
	}
	for _, filter := range query.Filters {
		projection[filter.Field] = 1
		if filter.Operation == "eq" {
			filters[filter.Field] = filter.Argument
		} else if filter.Operation == "domain" {
			filters[filter.Field] = bson.M{
				"$regex": regexp.QuoteMeta(filter.Argument),
			}
		}
	}
	options := options.Find()
	options.SetSort(bson.M{"id": -1})
	options.SetLimit(query.Limit)
	options.SetProjection(projection)

	cursor, err := storage.accounts.Find(context, filters, options)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		account := make(map[string]interface{})
		cursor.Decode(&account)
		delete(account, "_id")
		result = append(result, account)
	}

	return
}
