package main

import (
	"fmt"
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

var allowed = map[string]int{
	"email_domain":       1,
	"email_lt":           1,
	"email_gt":           1,
	"sex_eq":             1,
	"status_eq":          1,
	"status_neq":         1,
	"fname_eq":           1,
	"fname_any":          1,
	"fname_null":         1,
	"sname_eq":           1,
	"sname_null":         1,
	"sname_starts":       1,
	"phone_code":         1,
	"phone_null":         1,
	"country_null":       1,
	"country_eq":         1,
	"city_eq":            1,
	"city_any":           1,
	"city_null":          1,
	"birth_lt":           1,
	"birth_gt":           1,
	"birth_year":         1,
	"interests_any":      1,
	"interests_contains": 1,
	"likes_contains":     1,
	"premium_now":        1,
	"premium_null":       1,
}

func CreateAccountsQuery(query map[string]string) (accountsQuery AccountsQuery, err error) {
	if limit, ok := query["limit"]; ok {
		if parsed, e := strconv.ParseInt(limit, 10, 64); e == nil {
			if parsed <= 0 {
				err = &Error{400, "limit < 0"}
				return
			}
			accountsQuery.Limit = parsed
			delete(query, "limit")
		} else {
			err = &Error{400, "limit not number"}
			return
		}
	} else {
		err = &Error{400, "no limit"}
		return
	}

	delete(query, "query_id")

	for filter, arg := range query {
		if _, ok := allowed[filter]; !ok {
			err = &Error{400, fmt.Sprint("Bad field: ", filter)}
			return
		}
		splitted := strings.Split(filter, "_")
		accountsQuery.Filters = append(
			accountsQuery.Filters,
			AccountFilter{
				splitted[0],
				splitted[1],
				arg,
			})
	}

	return
}
