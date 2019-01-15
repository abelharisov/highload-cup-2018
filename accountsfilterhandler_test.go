package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	Accounts []struct {
		Id    int32
		Email string
		Sex   string
	}
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

func TestFilterSexEq(t *testing.T) {
	cases := []string{"m", "f"}

	for _, sex := range cases {
		t.Run(fmt.Sprint("Sex=", sex), func(t *testing.T) {
			teardown, storage := SetupTest(t)
			defer teardown()

			handler := AccountsFilterHandler{
				storage: storage,
			}

			request, _ := http.NewRequest("GET", fmt.Sprint("accounts/filter?sex_eq=", sex, "&limit=10"), nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			assert.Equal(t, http.StatusOK, responseRecorder.Code)

			response := Response{}

			err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
			assert.Nil(t, err, err)

			assert.NotEmpty(t, response)
			assert.NotEmpty(t, response.Accounts)
			for _, account := range response.Accounts {
				assert.Equal(t, account.Sex, sex)
			}
		})
	}
}
