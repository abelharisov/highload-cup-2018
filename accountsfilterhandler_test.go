package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/oleiade/reflections.v1"

	"github.com/stretchr/testify/assert"
)

// type ResponseAccount struct {
// 	Id     int32
// 	Email  *string
// 	Sex    *string
// 	Status *string
// 	Fname  *string
// 	Sname  *string
// 	Phone *string
// }

type ResponseAccount Account

type Response struct {
	Accounts []ResponseAccount
}

func (account *ResponseAccount) GetField(field string) *string {
	actualValue, err := reflections.GetField(account, field)
	if err != nil {
		panic(err)
	}

	return actualValue.(*string)
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
		assert.True(t, *prev >= *account.Id)
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
		"fname_bad=aaa",
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
	emailDomainAssert := func(t *testing.T, account ResponseAccount, value string) {
		ok, err := regexp.MatchString(
			fmt.Sprint("^.*@", regexp.QuoteMeta(value), "$"),
			*account.Email,
		)
		assert.NoError(t, err)
		assert.True(t, ok)
	}

	phoneCodeAssert := func(t *testing.T, account ResponseAccount, value string) {
		ok, err := regexp.MatchString(
			regexp.QuoteMeta(fmt.Sprint("(", value, ")")),
			*account.Phone,
		)
		assert.NoError(t, err)
		assert.True(t, ok)
	}

	birthYearAssert := func(t *testing.T, account ResponseAccount, value string) {
		actualTime := time.Unix(int64(*account.Birth), 0)
		expectedYear, _ := strconv.Atoi(value)
		assert.Equal(t, expectedYear, actualTime.Year())
	}

	premiumAssert := func(t *testing.T, account ResponseAccount, value string) {
		now := time.Now()
		start := time.Unix(int64((*account.Premium).Start), 0)
		finish := time.Unix(int64((*account.Premium).Finish), 0)
		assert.True(t, start.Before(now))
		assert.True(t, finish.After(now))
	}

	// interestsAnyAssert := func(t *testing.T, account ResponseAccount, value string) {
	// 	values := strings.Split(value, ",")
	// 	hasIntersect := false
	// loop:
	// 	for _, actualInterest := range *account.Interests {
	// 		for _, value := range values {
	// 			if actualInterest == value {
	// 				hasIntersect = true
	// 				break loop
	// 			}
	// 		}
	// 	}
	// 	assert.True(t, hasIntersect)
	// }

	// interestsContainsAssert := func(t *testing.T, account ResponseAccount, value string) {
	// 	values := strings.Split(value, ",")

	// 	for _, value := range values {
	// 		found := false
	// 		for _, actualInterest := range *account.Interests {
	// 			if actualInterest == value {
	// 				found = true
	// 				break
	// 			}
	// 		}
	// 		assert.True(t, found)
	// 	}
	// }

	// likesContainsAssert := func(t *testing.T, account ResponseAccount, value string) {
	// 	values := strings.Split(value, ",")
	// 	likes := make([]int32, 0)
	// 	for _, value := range values {
	// 		i, err := strconv.Atoi(value)
	// 		assert.NoError(t, err)
	// 		likes = append(likes, int32(i))
	// 	}

	// 	for _, like := range likes {
	// 		found := false
	// 		for _, actualLike := range *account.Likes {
	// 			if actualLike.Id == like {
	// 				found = true
	// 				break
	// 			}
	// 		}
	// 		assert.True(t, found)
	// 	}
	// }

	ltAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			assert.True(t, *account.GetField(field) < value, fmt.Sprint(*account.GetField(field), " < ", value))
		}
	}

	gtAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			assert.True(t, *account.GetField(field) > value)
		}
	}

	ltAssertInt := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			actualValue, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			valueInt, err := strconv.Atoi(value)
			assert.NoError(t, err)
			assert.True(t, *(actualValue.(*int32)) < int32(valueInt))
		}
	}

	gtAssertInt := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			actualValue, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			valueInt, err := strconv.Atoi(value)
			assert.NoError(t, err)
			assert.True(t, *(actualValue.(*int32)) > int32(valueInt))
		}
	}

	eqAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			assert.Equal(t, &value, account.GetField(field))
		}
	}

	neqAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			assert.NotEqual(t, &value, account.GetField(field))
		}
	}

	nullAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			actual, err := reflections.GetField(account, field)
			assert.NoError(t, err)
			if value == "0" {
				assert.NotNil(t, actual)
			} else {
				assert.Nil(t, actual)
			}
		}
	}

	anyAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			values := strings.Split(value, ",")
			assert.Contains(t, values, *account.GetField(field))
		}
	}

	startsAssert := func(field string) func(t *testing.T, account ResponseAccount, value string) {
		return func(t *testing.T, account ResponseAccount, value string) {
			assert.Equal(t, 0, strings.Index(*account.GetField(field), value))
		}
	}

	cases := []struct {
		param  string
		value  string
		assert func(t *testing.T, account ResponseAccount, value string)
	}{
		{"email_domain", "inbox.ru", emailDomainAssert},
		{"email_domain", "mail.ru", emailDomainAssert},
		{"email_lt", "bbbbbbb", ltAssert("Email")},
		{"email_gt", "ttttttt", gtAssert("Email")},
		{"sex_eq", "m", eqAssert("Sex")},
		{"sex_eq", "f", eqAssert("Sex")},
		{"status_eq", "свободны", eqAssert("Status")},
		{"status_eq", "заняты", eqAssert("Status")},
		{"status_neq", "свободны", neqAssert("Status")},
		{"status_neq", "заняты", neqAssert("Status")},
		{"fname_eq", "Алёна", eqAssert("FName")},
		{"fname_eq", "Павел", eqAssert("FName")},
		{"fname_any", "Алёна,Павел", anyAssert("FName")},
		{"fname_any", "Алёнааааааа,Павел", anyAssert("FName")},
		{"fname_null", "0", nullAssert("FName")},
		{"sname_eq", "Терленчан", eqAssert("SName")},
		{"sname_null", "0", nullAssert("SName")},
		{"sname_starts", "Тер", startsAssert("SName")},
		{"phone_code", "983", phoneCodeAssert},
		{"phone_code", "967", phoneCodeAssert},
		{"phone_null", "1", nullAssert("Phone")},
		{"phone_null", "0", nullAssert("Phone")},
		{"country_null", "0", nullAssert("Country")},
		{"country_null", "1", nullAssert("Country")},
		{"country_eq", "Индция", eqAssert("Country")},
		{"city_eq", "Росоград", eqAssert("City")},
		{"city_any", "Росоград,Пррпрпрп", anyAssert("City")},
		{"city_null", "0", nullAssert("City")},
		{"city_null", "1", nullAssert("City")},
		{"birth_lt", "806700770", ltAssertInt("Birth")},
		{"birth_gt", "806700770", gtAssertInt("Birth")},
		{"birth_year", "1995", birthYearAssert},
		// {"interests_any", "Апельсиновый сок,Матрица,YouTube", interestsAnyAssert},
		// {"interests_contains", "Апельсиновый сок,Матрица,YouTube", interestsContainsAssert},
		// {"likes_contains", "28499,2005", likesContainsAssert},
		{"premium_now", "1", premiumAssert},
		{"premium_null", "0", nullAssert("Premium")},
		{"premium_null", "1", nullAssert("Premium")},
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
				testCase.assert(t, account, testCase.value)
			}
		})
	}
}
