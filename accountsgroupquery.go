package main

import (
	"fmt"
	"strconv"
	"strings"
)

var validKeys = map[string]int{
	"sex":       1,
	"status":    1,
	"interests": 1,
	"country":   1,
	"city":      1,
}

// 2 means int should be sent to mongo
var allowedFilters = map[string]int{
	"email":     1,
	"sex":       1,
	"status":    1,
	"fname":     1,
	"sname":     1,
	"phone":     1,
	"country":   1,
	"city":      1,
	"birth":     1, // year
	"interests": 1, // one string
	"likes":     1, // one id
	// "premium":   1,
	"joined": 1, // year
}

type AccountsGroupQuery struct {
	Keys    []string
	Filters map[string]string
	Limit   int
	Order   int
}

func CreateAccountsGroupQuery(query map[string]string) (accountsGroupQuery AccountsGroupQuery, err error) {
	delete(query, "query_id")

	limit, ok := query["limit"]
	if !ok {
		err = &Error{400, "no limit"}
		return
	}
	accountsGroupQuery.Limit, err = strconv.Atoi(limit)
	if err != nil || accountsGroupQuery.Limit < 0 {
		err = &Error{400, fmt.Sprint("Bad limit", limit)}
		return
	}
	delete(query, "limit")

	order, ok := query["order"]
	if !ok {
		err = &Error{400, "No order"}
		return
	}
	accountsGroupQuery.Order, err = strconv.Atoi(order)
	if err != nil {
		err = &Error{400, fmt.Sprint("Bad order", order)}
		return
	}
	delete(query, "order")

	keys, ok := query["keys"]
	if !ok {
		err = &Error{400, "No keys"}
		return
	}
	accountsGroupQuery.Keys = strings.Split(keys, ",")
	for _, key := range accountsGroupQuery.Keys {
		if _, ok := validKeys[key]; !ok {
			err = &Error{400, fmt.Sprint("Bad key", key)}
			return
		}
	}
	delete(query, "keys")

	accountsGroupQuery.Filters = make(map[string]string)
	for field, values := range query {
		if _, ok := allowedFilters[field]; !ok {
			err = &Error{400, fmt.Sprint("Bad filter", field)}
			return
		}
		accountsGroupQuery.Filters[field] = values
	}

	return
}
