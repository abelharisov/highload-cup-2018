package main

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type MongoStorage struct {
	Uri      string
	Database string
	client   *mongo.Client
	accounts *mongo.Collection

	likesIndex *LikesIndex
	likees     LikeeMap
	interests  InterestsMap

	interestDict Dict
	countryDict  Dict
	cityDict     Dict

	recIndex AccountRecIndex

	now int
}

func (storage *MongoStorage) SetNow(now int) {
	storage.now = now
}

func (storage *MongoStorage) GetNow() int {
	return storage.now
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

	storage.likesIndex = CreateLikesIndex()

	storage.accounts = storage.client.Database(storage.Database).Collection("accounts")
	storage.CreateIndexes()
}

func (storage *MongoStorage) CreateIndexes() {
	context := context.Background()

	o := options.Index()
	o.SetUnique(true)
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"id": 1}, Options: o})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"email": 1}, Options: o})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"sex": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"birthYear": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"country": 1}})
	storage.accounts.Indexes().CreateOne(context, mongo.IndexModel{Keys: bson.M{"city": 1}})
}

func (storage *MongoStorage) DropIndexes() {
	log.Println("Todo")
}

func (storage *MongoStorage) LoadAccounts(accounts []Account) (err error) {
	context := context.Background()
	documents := make([]interface{}, 0, len(accounts))

	for _, account := range accounts {
		documents = append(documents, interface{}(account))
		if account.LikeIds != nil {
			for _, likee := range *account.LikeIds {
				storage.likees.AppendLiker(likee, account.Id)
			}
		}
		storage.likesIndex.AddLikes(account)
		if account.Interests != nil {
			storage.interests.Append(account.Id, *account.Interests)
		}
		storage.recIndex.Add(account)
	}
	_, insertErr := storage.accounts.InsertMany(context, documents)
	if insertErr != nil {
		err = &Error{400, insertErr.Error()}
		return
	}

	return nil
}

func (storage *MongoStorage) DropDatabase() {
	context := context.Background()
	storage.client.Database(storage.Database).Drop(context)
}
