package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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
		} else if filter.Operation == "neq" {
			filters[filter.Field] = bson.M{
				"$ne": filter.Argument,
			}
		} else if filter.Operation == "domain" {
			filters[filter.Field] = bson.M{
				"$regex": fmt.Sprint(regexp.QuoteMeta(filter.Argument), "$"),
			}
		} else if filter.Operation == "null" {
			filters[filter.Field] = bson.M{
				"$exists": filter.Argument == "0",
			}
		} else if filter.Operation == "any" {
			values := strings.Split(filter.Argument, ",")
			filters[filter.Field] = bson.M{
				"$in": values,
			}
		} else if filter.Operation == "lt" || filter.Operation == "gt" {
			intArg, err := strconv.Atoi(filter.Argument)
			if err == nil {
				filters[filter.Field] = bson.M{
					fmt.Sprint("$", filter.Operation): intArg,
				}
			} else {
				filters[filter.Field] = bson.M{
					fmt.Sprint("$", filter.Operation): filter.Argument,
				}
			}
		} else if filter.Operation == "starts" {
			filters[filter.Field] = bson.M{
				"$regex": fmt.Sprint("^", regexp.QuoteMeta(filter.Argument)),
			}
		} else if filter.Operation == "code" {
			filters[filter.Field] = bson.M{
				"$regex": regexp.QuoteMeta(fmt.Sprint("(", filter.Argument, ")")),
			}
		} else if filter.Operation == "year" && filter.Field == "birth" {
			intArg, err := strconv.Atoi(filter.Argument)
			if err != nil {
				panic(err)
			}
			filters["year"] = intArg
		} else if filter.Operation == "now" && filter.Field == "premium" {
			now := time.Now().Unix()
			filters["premium.start"] = bson.M{
				"$lte": now,
			}
			filters["premium.finish"] = bson.M{
				"$gte": now,
			}
		} else if filter.Operation == "contains" && filter.Field == "interests" {
			values := strings.Split(filter.Argument, ",")
			filters[filter.Field] = bson.M{
				"$all": values,
			}
		} else if filter.Operation == "contains" && filter.Field == "likes" {
			values := strings.Split(filter.Argument, ",")
			ids := make([]int, 0, len(values))
			for _, value := range values {
				id, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				ids = append(ids, id)
			}
			filters["likeIds"] = bson.M{
				"$all": ids,
			}
		}
	}
	delete(projection, "likes")
	delete(projection, "interests")
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

	if result == nil {
		result = make([]map[string]interface{}, 0)
	}

	return
}
