package main

import (
	"strings"
	"strconv"
	"errors"
)

type AccountFilter struct {
	Field     string
	Operation string
	Argument  string
}

type AccountsQuery struct {
	Limit int32
	Filters []AccountFilter
}

var NoLimitError = errors.New("No limit")

func CreateAccountsFilter(query map[string][]string) (accountsQuery AccountsQuery, err error) {
	if limit, ok := query["limit"]; ok {
		if parsed, e := strconv.ParseInt(limit[0], 10, 32); e != nil {
			accountsQuery.Limit = int32(parsed)
			delete(query, "limit")
		} else {
			err = NoLimitError
			return
		}
	} else {
		err = NoLimitError
		return
	}

	for filter, arg := range query {
		splitted := strings.Split(filter, "_")
		accountsQuery.Filters = append(
			accountsQuery.Filters,
			AccountFilter{
				splitted[0],
				splitted[1],
				arg[0],
			})
	}
	return
}
