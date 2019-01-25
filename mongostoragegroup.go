package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
)

func (storage *MongoStorage) Group(query *AccountsGroupQuery) (result []map[string]interface{}, err error) {
	context := context.Background()

	match := bson.M{}
	for field, value := range query.Filters {
		switch field {
		case "birth", "joined":
			year, convErr := strconv.Atoi(value)
			if convErr != nil {
				err = &Error{400, fmt.Sprint("Bad", field, value)}
				return
			}
			match[fmt.Sprint(field, "Year")] = year
		case "likes":
			likeeId, convErr := strconv.Atoi(value)
			if convErr != nil {
				err = &Error{400, fmt.Sprint("Bad like", value)}
				return
			}
			if likers, ok := storage.likees[likeeId]; ok {
				match["id"] = bson.M{
					"$in": likers,
				}
			} else {
				return // ???
			}
			// match["likeIds"] = likeeId
		default:
			match[field] = value
		}
	}

	groupBy := bson.M{}
	sortStage := bson.D{
		{"count", query.Order},
	}
	for _, field := range query.Keys {
		groupBy[field] = fmt.Sprint("$", field)
	}

	sort.Strings(query.Keys)

	var unwindInterests = false
	for _, field := range query.Keys {
		sortStage = append(
			sortStage,
			bson.E{
				fmt.Sprint("_id.", field),
				query.Order,
			},
		)

		if field == "interests" {
			unwindInterests = true
		}
	}

	pipeline := bson.A{
		bson.M{
			"$match": match,
		},
	}

	if unwindInterests {
		pipeline = append(
			pipeline,
			bson.M{
				"$unwind": "$interests",
			},
		)
	}

	pipeline = append(
		pipeline,
		bson.M{
			"$group": bson.M{
				"_id": groupBy,
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$sort": sortStage,
		},
		bson.M{
			"$limit": query.Limit,
		},
	)

	cursor, findErr := storage.accounts.Aggregate(context, pipeline)
	if findErr != nil {
		err = &Error{500, findErr.Error()}
		return
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		data := make(map[string]interface{})
		cursor.Decode(&data)
		result = append(result, data)
	}

	return
}
