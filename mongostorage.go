package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type MongoStorage struct {
	Uri      string
	Database string
	client   *mongo.Client
	accounts *mongo.Collection

	likees    LikeeMap
	interests InterestsMap

	interestDict Dict
	countryDict  Dict
	cityDict     Dict

	recIndex AccountRecIndex

	now int
}

func (storage *MongoStorage) Init() {
	context := context.Background()

	var err error
	storage.client, err = mongo.Connect(context, storage.Uri)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	storage.likees = make(LikeeMap)
	storage.interests = make(InterestsMap)

	storage.interestDict = Dict{}
	storage.interestDict.Init()

	storage.countryDict = Dict{}
	storage.countryDict.Init()

	storage.cityDict = Dict{}
	storage.cityDict.Init()

	storage.recIndex = AccountRecIndex{
		InterestDict: &storage.interestDict,
		CountryDict:  &storage.countryDict,
		CityDict:     &storage.cityDict,
	}
	storage.recIndex.Init()

	storage.accounts = storage.client.Database(storage.Database).Collection("accounts")
	storage.CreateIndexes()
}

func (storage *MongoStorage) CreateIndexes() {
	context := context.Background()

	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"sex": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"birthYear": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"country": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"city": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"email": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"id": 1}})
	// storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{
	// 	Keys: bson.D{
	// 		{"sex", 1},
	// 		{"country", 1},
	// 		{"city", 1},
	// 		{"status", 1},
	// 		{"email", 1},
	// 	},
	// })
}

func (storage *MongoStorage) DropIndexes() {
	log.Println("Todo")
}

func (storage *MongoStorage) LoadAccounts(accounts []Account) {
	context := context.Background()
	documents := make([]interface{}, 0, len(accounts))

	for _, account := range accounts {
		documents = append(documents, interface{}(account))
		if account.LikeIds != nil {
			for _, likee := range *account.LikeIds {
				storage.likees.AppendLiker(likee, account.Id)
			}
		}
		if account.Interests != nil {
			storage.interests.Append(account.Id, *account.Interests)
		}
		storage.recIndex.Add(account)
	}
	storage.accounts.InsertMany(context, documents)
}

func (storage *MongoStorage) DropDatabase() {
	context := context.Background()
	storage.client.Database(storage.Database).Drop(context)
}

func (s *MongoStorage) SetNow(now int) {
	s.now = now
}

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
			if code, err := strconv.Atoi(filter.Argument); err == nil {
				filters["phoneCode"] = code
			} else {
				log.Println(err)
			}
		} else if filter.Operation == "year" && filter.Field == "birth" {
			intArg, err := strconv.Atoi(filter.Argument)
			if err != nil {
				panic(err)
			}
			filters["birthYear"] = intArg
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
				if i, err := strconv.Atoi(value); err == nil {
					intValues = append(intValues, i)
				} else {
					panic(err)
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
			result = make([]map[string]interface{}, 0)
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
		err = findErr
		return
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

func (storage *MongoStorage) Group(query *AccountsGroupQuery) (result []map[string]interface{}, err error) {
	context := context.Background()

	match := bson.M{}
	for field, value := range query.Filters {
		switch field {
		case "birth", "joined":
			year, convErr := strconv.Atoi(value)
			if convErr != nil {
				err = convErr
				return
			}
			match[fmt.Sprint(field, "Year")] = year
		case "likes":
			likeeId, convErr := strconv.Atoi(value)
			if convErr != nil {
				err = convErr
				return
			}
			if likers, ok := storage.likees[likeeId]; ok {
				match["id"] = bson.M{
					"$in": likers,
				}
			}
			match["likeIds"] = likeeId
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
	// keys := make([]string, 0)
	// if query.Order == -1 {
	// 	for i := len(query.Keys) - 1; i >= 0; i-- {
	// 		keys = append(keys, query.Keys[i])
	// 	}
	// } else {
	// 	keys = query.Keys
	// }

	var unwindInterests = false
	for _, field := range query.Keys /*keys*/ {
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
		bson.A{
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
		}...,
	)

	cursor, findErr := storage.accounts.Aggregate(context, pipeline)

	if findErr != nil {
		err = findErr
		return
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		data := make(map[string]interface{})
		cursor.Decode(&data)
		result = append(result, data)
	}

	if result == nil {
		result = make([]map[string]interface{}, 0)
	}

	return
}

func (storage *MongoStorage) Recommend(q *AccountsRecommendQuery) (result []map[string]interface{}, err error) {
	// log.Println(*q)
	res := storage.accounts.FindOne(context.Background(), bson.M{"id": q.Id})
	if res == nil || res.Err() != nil {
		log.Println(res.Err())
		err = errors.New("Recommend: account not found")
		return
	}

	var account Account
	if err = res.Decode(&account); err != nil {
		return
	}

	ids := storage.recIndex.Recommend(account, q.Country, q.City, q.Limit)

	options := options.Find()
	options.SetProjection(bson.M{
		"sname":   1,
		"id":      1,
		"fname":   1,
		"birth":   1,
		"email":   1,
		"status":  1,
		"premium": 1,
	})

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
		err = findEerr
		return
	}

	result = make([]map[string]interface{}, 0, len(ids))

	for _, id := range ids {
		result = append(result, notOrderedResult[int64(id)])
	}

	return
}
