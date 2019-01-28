package main

import (
	"context"
	"math"
	"sort"

	"github.com/thoas/go-funk"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

func (s *MongoStorage) Suggset(q *AccountsRecommendQuery) (result []map[string]interface{}, err error) {
	likes := s.likesIndex.GetLikes(q.Id)
	// log.Println("target likes", *likes)
	if likes == nil {
		err = &Error{404, "account not found or no likes"}
		return
	}

	ids := make([]int, 0, likes.Len())
	for like, _ := range likes.ids {
		ids = append(ids, like)
	}

	similars := s.likeeToLikerIndex.AccountsWithLikesAny(ids)

	filterErr := s.filterByCity(q, &similars)
	if filterErr != nil {
		err = filterErr
		return
	}

	sort.Sort(&bySimilarity{
		source:     likes,
		target:     &similars,
		similarity: map[int]float32{},
		likesIndex: s.likesIndex},
	)

	// log.Println("similars", similars)

	suggests := []int{}

	// filter by sex
	// uniq for suggests

	for _, similarId := range similars {
		if similarId == q.Id {
			continue
		}
		slikes := s.likesIndex.GetLikes(similarId)
		// log.Println("slikes", slikes)
		ids := []int{}
		for id, _ := range slikes.ids {
			if !s.likesIndex.HasLike(q.Id, id) {
				ids = append(ids, id)
			}
		}
		sort.Ints(ids)
		ids = funk.ReverseInt(ids)
		suggests = append(suggests, ids...)
		if len(suggests) >= q.Limit {
			break
		}
	}

	if len(suggests) > q.Limit {
		suggests = suggests[0:q.Limit]
	}

	projection := bson.M{
		"id":     1,
		"status": 1,
		"fname":  1,
		"sname":  1,
		"email":  1,
	}
	o := options.Find()
	o.SetProjection(projection)
	cursor, findErr := s.accounts.Find(
		context.Background(),
		bson.M{"id": bson.M{"$in": suggests}},
		o,
	)
	if findErr != nil {
		err = &Error{500, findErr.Error()}
		return
	}

	// log.Println("ids", suggests)

	data := map[int]bson.M{}
	for cursor.Next(context.Background()) {
		var account bson.M
		cursor.Decode(&account)
		delete(account, "_id")
		data[int(account["id"].(int64))] = account
	}

	sortedData := make([]map[string]interface{}, 0, len(suggests))
	for _, suggestId := range suggests {
		sortedData = append(sortedData, data[suggestId])
	}

	result = sortedData

	return
}

func (s *MongoStorage) filterByCity(q *AccountsRecommendQuery, similars *[]int) (err error) {
	if len(q.City) > 0 || len(q.Country) > 0 {
		o := options.Find()
		o.SetProjection(bson.M{"id": 1})
		filter := bson.M{
			"id": bson.M{"$in": similars},
		}
		if len(q.City) > 0 {
			filter["city"] = q.City
		}
		if len(q.Country) > 0 {
			filter["country"] = q.Country
		}
		cursor, findErr := s.accounts.Find(
			context.Background(),
			filter,
			o,
		)
		if findErr != nil {
			err = &Error{500, findErr.Error()}
		}
		defer cursor.Close(context.Background())

		filteredSimilars := []int{}
		for cursor.Next(context.Background()) {
			id := map[string]interface{}{}
			cursor.Decode(&id)
			filteredSimilars = append(filteredSimilars, int(id["id"].(int64)))
		}
		*similars = filteredSimilars
	}

	return nil
}

type bySimilarity struct {
	source     *Likes
	target     *[]int
	similarity map[int]float32
	likesIndex *LikesIndex
}

func (a *bySimilarity) GetSimilarity(i int) float32 {
	id := (*a.target)[i]
	if value, ok := a.similarity[id]; ok {
		return value
	}
	likes := a.likesIndex.GetLikes(id)
	if likes == nil {
		return 0
	}

	similarity := 0.0
	for like, ts := range a.source.ids {
		if sameLikeTs, ok := likes.ids[like]; ok {
			sum := math.Abs(float64(ts - sameLikeTs))
			if sum == 0 {
				sum = 1.0
			}
			similarity += 1.0 / sum
		}
	}

	return float32(similarity)

}

func (a *bySimilarity) Len() int {
	return len(*a.target)
}

func (a *bySimilarity) Less(i, j int) bool {
	similarityI := a.GetSimilarity(i)
	similarityJ := a.GetSimilarity(j)

	return similarityI > similarityJ
}

func (a *bySimilarity) Swap(i, j int) {
	(*a.target)[i], (*a.target)[j] = (*a.target)[j], (*a.target)[i]
}
