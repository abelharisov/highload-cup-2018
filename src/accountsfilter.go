package main

import (
	"errors"
	"strconv"
	"strings"
)

type AccountFilter struct {
	Field     string
	Operation string
	Argument  string
}

type AccountsQuery struct {
	Limit   int64
	Filters []AccountFilter
}

var NoLimitError = errors.New("No limit")

func CreateAccountsQuery(query map[string][]string) (accountsQuery AccountsQuery, err error) {
	if limit, ok := query["limit"]; ok {
		if parsed, e := strconv.ParseInt(limit[0], 10, 64); e == nil {
			if parsed == 0 {
				err = NoLimitError
				return
			}
			accountsQuery.Limit = parsed
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
