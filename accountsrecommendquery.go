package main

import (
	"fmt"
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
		err = &Error{400, "no limit"}
		return
	}
	accountsRecommendQuery.Limit, err = strconv.Atoi(limit)
	if err != nil || accountsRecommendQuery.Limit < 0 {
		err = &Error{400, fmt.Sprint("bad limit", limit)}
		return
	}

	if id, idErr := strconv.Atoi(idParam); err == nil {
		accountsRecommendQuery.Id = id
	} else {
		err = &Error{500, idErr.Error()}
		return
	}

	if country, ok := query["country"]; ok {
		if len(country) == 0 {
			err = &Error{400, "empty country"}
			return
		}
		accountsRecommendQuery.Country = country
	}

	if city, ok := query["city"]; ok {
		if len(city) == 0 {
			err = &Error{400, "empty city"}
			return
		}
		accountsRecommendQuery.City = city
	}

	return
}
