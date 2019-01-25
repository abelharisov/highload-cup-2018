package main

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

var projection = bson.M{
	"sname":   1,
	"id":      1,
	"fname":   1,
	"birth":   1,
	"email":   1,
	"status":  1,
	"premium": 1,
}

func (storage *MongoStorage) Recommend(q *AccountsRecommendQuery) (result []map[string]interface{}, err error) {
	res := storage.accounts.FindOne(context.Background(), bson.M{"id": q.Id})
	if res.Err() != nil {
		err = &Error{404, "account not found"}
		return
	}

	var account Account
	if decodeErr := res.Decode(&account); decodeErr != nil {
		err = &Error{404, "account not found"}
		return
	}

	ids := storage.recIndex.Recommend(account, q.Country, q.City, q.Limit)

	options := options.Find()
	options.SetProjection(projection)

	filter := bson.M{
		"id": bson.M{
			"$in": ids,
		},
	}

	notOrderedResult := map[int64](map[string]interface{}){}

	if cursor, findEerr := storage.accounts.Find(context.Background(), filter, options); findEerr == nil {
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			data := make(map[string]interface{})
			cursor.Decode(&data)
			delete(data, "_id")
			notOrderedResult[data["id"].(int64)] = data
		}
	} else {
		err = &Error{500, findEerr.Error()}
		return
	}

	result = make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		result = append(result, notOrderedResult[int64(id)])
	}

	return
}
