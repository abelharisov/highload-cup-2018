package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"fmt"
)

type Response struct {
	Accounts []struct{
		Id int32
		Email string
		Sex string
	}
}

func TestFilterSexEqM(t *testing.T) {
	teardown, storage := SetupTest(t)
	defer teardown()

	handler := AccountsFilterHandler{
		storage: storage,
	}

	request, _ := http.NewRequest("GET", "accounts/filter?sex_eq=m", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	response := Response{}
	fmt.Println(responseRecorder.Body.String())

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err, err)

	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Accounts)
	for _, account := range response.Accounts {
		assert.Equal(t, account.Sex, "m")
	}
}
// func TestFilterSexEqF() {
	
// }