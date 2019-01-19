package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"gopkg.in/oleiade/reflections.v1"

	"github.com/stretchr/testify/assert"
)

type ResponseAccount struct {
	Id     int32
	Email  string
	Sex    string
	Status string
}

type Response struct {
	Accounts []ResponseAccount
}

func TestLimit(t *testing.T) {
	cases := []int{10, 50}

	for _, limit := range cases {
		t.Run(fmt.Sprint("Limit=", limit), func(t *testing.T) {
			teardown, storage := SetupTest(t)
			defer teardown()

			handler := AccountsFilterHandler{
				storage: storage,
			}

			request, _ := http.NewRequest("GET", fmt.Sprint("accounts/filter?sex_eq=m&limit=", limit), nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, http.StatusOK, responseRecorder.Code)

			response := Response{}

			err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
			assert.Nil(t, err, err)

			assert.NotEmpty(t, response)
			assert.NotEmpty(t, response.Accounts)
			assert.Equal(t, limit, len(response.Accounts))
		})
	}
}

func TestAreAccountsSorted(t *testing.T) {
	teardown, storage := SetupTest(t)
	defer teardown()

	handler := AccountsFilterHandler{
		storage: storage,
	}

	request, _ := http.NewRequest("GET", "accounts/filter?sex_eq=m&limit=100", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	response := Response{}

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err, err)

	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Accounts)
	prev := response.Accounts[0].Id
	for _, account := range response.Accounts {
		assert.True(t, prev >= account.Id, fmt.Sprint(prev, " >= ", account.Id, " is falsy"))
		prev = account.Id
	}
}

func TestFilterFiledsInResponse(t *testing.T) {
	requiredFields := []string{"id", "email"}
	cases := map[string]([]string){
		"sex_eq=f": append(requiredFields, []string{"sex"}...),
	}

	for query, fields := range cases {
		t.Run(query, func(t *testing.T) {
			teardown, storage := SetupTest(t)
			defer teardown()

			handler := AccountsFilterHandler{
				storage: storage,
			}

			request, _ := http.NewRequest("GET", fmt.Sprint("accounts/filter?limit=10&", query), nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, http.StatusOK, responseRecorder.Code)

			response := map[string]interface{}{}

			err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
			assert.Nil(t, err, err)

			assert.NotEmpty(t, response)
			accounts := response["accounts"].([]interface{})
			account := accounts[0].(map[string]interface{})

			for k := range account {
				assert.Contains(t, fields, k)
			}
		})
	}
}

func TestBadRequests(t *testing.T) {
	cases := []string{
		"",
		"limit=",
		"limit=f",
		"limit=-10",
		"limit=111f",
	}

	for _, query := range cases {
		t.Run(query, func(t *testing.T) {
			teardown, storage := SetupTest(t)
			defer teardown()

			handler := AccountsFilterHandler{
				storage: storage,
			}

			request, _ := http.NewRequest("GET", fmt.Sprint("accounts/filter?", query), nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestFilter(t *testing.T) {
	domainAssert := func(account ResponseAccount, value interface{}) {
		ok, err := regexp.MatchString(
			fmt.Sprint("^.*@", regexp.QuoteMeta(value.(string)), "$"),
			account.Email,
		)
		assert.NoError(t, err)
		assert.True(t, ok)
	}

	ltAssert := func(field string) func(account ResponseAccount, value interface{}) {
		return func(account ResponseAccount, value interface{}) {
			actualValue, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			assert.True(t, actualValue.(string) < value.(string))
		}
	}

	gtAssert := func(field string) func(account ResponseAccount, value interface{}) {
		return func(account ResponseAccount, value interface{}) {
			actualValue, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			assert.True(t, actualValue.(string) > value.(string))
		}
	}

	eqAssert := func(field string) func(account ResponseAccount, value interface{}) {
		return func(account ResponseAccount, value interface{}) {
			actualValue, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			assert.Equal(t, value.(string), actualValue.(string))
		}
	}

	cases := []struct {
		param string
		value string
		check func(account ResponseAccount, value interface{})
	}{
		{"email_domain", "inbox.ru", domainAssert},
		{"email_domain", "mail.ru", domainAssert},
		{"email_lt", "bbbbbbb", gtAssert("Email")},
		{"email_gt", "ttttttt", ltAssert("Email")},
		{"sex_eq", "m", eqAssert("Sex")},
		{"sex_eq", "f", eqAssert("Sex")},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprint(testCase.param, "=", testCase.value), func(t *testing.T) {
			teardown, storage := SetupTest(t)
			defer teardown()

			handler := AccountsFilterHandler{
				storage: storage,
			}

			request, _ := http.NewRequest("GET", fmt.Sprint("accounts/filter?limit=10&", testCase.param, "=", testCase.value), nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, http.StatusOK, responseRecorder.Code)

			response := Response{}

			err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
			assert.Nil(t, err, err)

			assert.NotEmpty(t, response)
			assert.NotEmpty(t, response.Accounts)
			for _, account := range response.Accounts {
				testCase.check(account, testCase.value)
			}
		})
	}
}
