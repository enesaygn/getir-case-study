package test

import (
	"bytes"
	"encoding/json"
	"getir-case/db"
	"getir-case/endpoints"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoHandler(t *testing.T) {

	db.InitializeMongoDB()

	reqBody := endpoints.MongoHandlerRequestBody{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Errorf("Error in marshalling request body: %v", err)
	}

	handler := http.HandlerFunc(endpoints.MongoHandler)

	server := httptest.NewServer(handler)

	// Send a POST request to our serve
	response, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Errorf("Error in sending POST request: %v", err)
	}

	assert.Equal(t, 200, response.StatusCode)

	expectedResponse := `{
		"code": 0,
		"message": "Success",
		"records": [
			{
				"key": "TAKwGc6Jr4i8Z487",
				"createdAt": "2017-01-28T04:22:14Z",
				"totalCount": 2800
			},
			{
				"key": "NAeQ8eX7e5TEg7oH",
				"createdAt": "2017-01-27T11:19:14Z",
				"totalCount": 2900
			}
		]
	}`

	returnedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error in reading response body: %v", err)
	}

	var expectedResponseI interface{}
	var returnedResponseI interface{}

	err = json.Unmarshal([]byte(expectedResponse), &expectedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling expected response: %v", err)
	}

	err = json.Unmarshal(returnedResponse, &returnedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling returned response: %v", err)
	}

	assert.Equal(t, expectedResponseI, returnedResponseI)

}

func TestInMemoryGet(t *testing.T) {

	// Initialize the in-memory database
	instance := db.GetInMemoryDBInstance()

	// Set a key-value pair
	instance.Set("testKey", "testValue")

	// Create a new request
	handler := http.HandlerFunc(endpoints.InMemoryGetHandler)

	// Create a new server
	server := httptest.NewServer(handler)

	// Set the query parameters
	params := "?key=testKey"

	// Send a POST request to our serve
	response, err := http.Get(server.URL + params)
	if err != nil {
		t.Errorf("Error in sending POST request: %v", err)
	}

	// Check the status code
	assert.Equal(t, 200, response.StatusCode)

	expectedResponse := `{
		"key": "testKey",
		"value": "testValue"
	}`

	var expectedResponseI interface{}
	var returnedResponseI interface{}

	err = json.Unmarshal([]byte(expectedResponse), &expectedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling expected response: %v", err)
	}

	returnedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error in reading response body: %v", err)
	}

	err = json.Unmarshal(returnedResponse, &returnedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling returned response: %v", err)
	}

	assert.Equal(t, expectedResponseI, returnedResponseI)

}

func TestInMemorySet(t *testing.T) {

	reqBody := endpoints.SetHandlerRequestBody{
		Key:   "testKey",
		Value: "testValue",
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Errorf("Error in marshalling request body: %v", err)
	}

	handler := http.HandlerFunc(endpoints.InMemorySetHandler)

	server := httptest.NewServer(handler)

	// Send a POST request to our serve

	response, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Errorf("Error in sending POST request: %v", err)
	}

	assert.Equal(t, 200, response.StatusCode)

	expectedResponse := `{
		"key": "testKey",
		"value": "testValue"
	}`
	returnedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error in reading response body: %v", err)
	}

	var expectedResponseI interface{}
	var returnedResponseI interface{}

	err = json.Unmarshal([]byte(expectedResponse), &expectedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling expected response: %v", err)
	}

	err = json.Unmarshal(returnedResponse, &returnedResponseI)
	if err != nil {
		t.Errorf("Error in unmarshalling returned response: %v", err)
	}

	assert.Equal(t, expectedResponseI, returnedResponseI)

}
