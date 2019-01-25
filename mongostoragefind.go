package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	funk "github.com/thoas/go-funk"
)

func (storage *MongoStorage) Find(query *AccountsQuery) (result []map[string]interface{}, err error) {
	context := context.Background()

	filters := make(map[string]interface{})
	projection := bson.M{
		"id":    1,
		"email": 1,
		// "statusId": 1,
		// "phoneCode": 1,
	}
	var preIds *[]int
	for _, filter := range query.Filters {
		projection[filter.Field] = 1
		if filter.Field == "status" {
			if filter.Operation == "eq" {
				filters["statusId"] = ParseStatus(filter.Argument)
			} else {
				filters["statusId"] = bson.M{
					"$ne": ParseStatus(filter.Argument),
				}
			}
		} else if filter.Operation == "eq" {
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
			if filter.Field == "premium" {
				if filter.Argument == "0" {
					filters["premiumStatus"] = bson.M{
						"$ne": PremiumNull,
					}
				} else {
					filters["premiumStatus"] = PremiumNull
				}

			} else {
				filters[filter.Field] = bson.M{
					"$exists": filter.Argument == "0",
				}
			}
		} else if filter.Operation == "any" {
			values := strings.Split(filter.Argument, ",")
			filters[filter.Field] = bson.M{
				"$in": values,
			}
			// if filter.Field == "interests" {
			// 	iii := storage.interests.AccountsWithInterestsAny(values)
			// 	log.Println("aaaaa ", len(iii))
			// }
		} else if filter.Operation == "lt" || filter.Operation == "gt" {
			if intArg, err := strconv.Atoi(filter.Argument); err == nil {
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
			if code, codeErr := strconv.Atoi(filter.Argument); codeErr == nil {
				filters["phoneCode"] = code
			} else {
				err = &Error{400, fmt.Sprint("Bad code", filter.Argument)}
				return
			}
		} else if filter.Operation == "year" && filter.Field == "birth" {
			if intArg, yearErr := strconv.Atoi(filter.Argument); yearErr == nil {
				filters["birthYear"] = intArg
			} else {
				err = &Error{400, fmt.Sprint("Bad year", filter.Argument)}
				return
			}
		} else if filter.Operation == "now" {
			filters["premiumStatus"] = PremiumActive
		} else if filter.Operation == "contains" && filter.Field == "interests" {
			values := strings.Split(filter.Argument, ",")
			ids := storage.interests.AccountsWithInterestsContains(values)
			if preIds == nil {
				preIds = &ids
			} else {
				intersect := funk.Intersect(*preIds, ids).([]int)
				preIds = &intersect
			}
		} else if filter.Operation == "contains" && filter.Field == "likes" {
			values := strings.Split(filter.Argument, ",")
			intValues := make([]int, 0, len(values))
			for _, value := range values {
				if i, valueErr := strconv.Atoi(value); valueErr == nil {
					intValues = append(intValues, i)
				} else {
					err = &Error{400, fmt.Sprint("Bad year", filter.Argument)}
					return
				}
			}
			ids := storage.likees.AccountsWithLikesContains(intValues)
			if preIds == nil {
				preIds = &ids
			} else {
				intersect := funk.Intersect(*preIds, ids).([]int)
				preIds = &intersect
			}
		}
	}

	if preIds != nil {
		if len(*preIds) == 0 {
			return
		}

		filters["id"] = bson.M{
			"$in": preIds,
		}
	}

	delete(projection, "likes")
	delete(projection, "interests")
	options := options.Find()
	options.SetSort(bson.M{"id": -1})
	options.SetLimit(query.Limit)
	options.SetProjection(projection)

	cursor, findErr := storage.accounts.Find(context, filters, options)
	if findErr != nil {
		err = &Error{500, findErr.Error()}
		return
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
