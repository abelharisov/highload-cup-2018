package main

import (
	"strconv"
)

type AccountsRecommendQuery struct {
	Id      int
	Limit   int
	Country string
	City    string
}

func CreateAccountsRecommendQuery(idParam string, query map[string]string) (accountsRecommendQuery AccountsRecommendQuery, err error) {
	limit, ok := query["limit"]
	if !ok {
		err = NoLimitError
		return
	}
	accountsRecommendQuery.Limit, err = strconv.Atoi(limit)
	if err != nil || accountsRecommendQuery.Limit < 0 {
		err = NoLimitError
		return
	}

	if id, err := strconv.Atoi(idParam); err == nil {
		accountsRecommendQuery.Id = id
	}

	if country, ok := query["country"]; ok {
		accountsRecommendQuery.Country = country
	}

	if city, ok := query["city"]; ok {
		accountsRecommendQuery.City = city
	}

	return
}
