package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"getir-case/db"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoHandlerRequestBody struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type mongoHandlerResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Records []struct {
		Key        string `json:"key"`
		CreatedAt  string `json:"createdAt"`
		TotalCount int    `json:"totalCount"`
	} `json:"records"`
}

type SetHandlerRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type setHandlerResponseBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetHandlerRequestBody struct {
	Key string `json:"key"`
}

type getHandlerResponseBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// getHandler handles GET requests to the root URL
func InMemoryGetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	// Check the request method
	if r.Method != http.MethodGet {
		log.Println(r.RequestURI, " ", "Method Not Allowed")
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the query parameters from the request
	queryParams := r.URL.Query()

	// Retrieve specific parameter values
	key := queryParams.Get("key")
	res := getHandlerResponseBody{}

	inMemoryInstance := db.GetInMemoryDBInstance()

	// Get the value for the given key from the in-memory database
	val, ok := inMemoryInstance.Get(key)
	if ok {
		// If the key exists, respond with the value
		res.Key = key
		res.Value = val
		json.NewEncoder(w).Encode(res)
	} else {
		// If the key does not exist, respond with a not found error
		res.Key = ""
		res.Value = ""
		json.NewEncoder(w).Encode(res)
	}

}

// setHandler handles POST requests to the /inmemoryset URL
func InMemorySetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Check the request method
	if r.Method != http.MethodPost {
		log.Println(r.RequestURI, " ", "Method Not Allowed")
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a new instance of the in-memory database
	inMemoryInstance := db.GetInMemoryDBInstance()

	// Decode the request body
	req := SetHandlerRequestBody{}
	res := setHandlerResponseBody{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	inMemoryInstance.Set(req.Key, req.Value)

	responseSetBody(req.Key, req.Value, w, res)

}

// MongoHandler handles POST requests to the /mongofetch URL
func MongoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		log.Println(r.RequestURI, " ", "Method Not Allowed")
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	req := MongoHandlerRequestBody{}
	res := mongoHandlerResponseBody{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("MongoHandler Decode Error: ", err.Error())
		responseMongoBody(err.Error(), 1, w, res)
		return
	}

	defer r.Body.Close()

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		log.Println("MongoHandler startDate Parse Error: ", err.Error())
		responseMongoBody(err.Error(), 1, w, res)
		return
	}
	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		log.Println("MongoHandler endDate Parse Error: ", err.Error())
		responseMongoBody(err.Error(), 1, w, res)
		return
	}

	// Access collection
	collection := db.Connection.Collection("records")

	matchedCreatedAt := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "createdAt", Value: bson.D{
				{Key: "$lte", Value: primitive.NewDateTimeFromTime(endDate)},
				{Key: "$gte", Value: primitive.NewDateTimeFromTime(startDate)}},
			},
		},
		},
	}
	projectTotalCount := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: "$counts"}}},
			{Key: "_id", Value: 1},
			{Key: "key", Value: 1},
			{Key: "createdAt", Value: 1},
		},
		},
	}
	matchTotalCount := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "totalCount", Value: bson.D{
				{Key: "$lte", Value: req.MaxCount},
				{Key: "$gte", Value: req.MinCount}}},
		},
		},
	}

	pipeline := mongo.Pipeline{matchedCreatedAt, projectTotalCount, matchTotalCount}

	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println(err.Error())
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			res.Code = 1
			res.Message = "Invalid request"
			json.NewEncoder(w).Encode(res)
			return
		}

		var record struct {
			Key        string `json:"key"`
			CreatedAt  string `json:"createdAt"`
			TotalCount int    `json:"totalCount"`
		}

		record.Key = result.Map()["key"].(string)
		a := result.Map()["createdAt"].(primitive.DateTime).Time()
		record.CreatedAt = a.Format("2006-01-02T15:04:05Z")
		b := int(result.Map()["totalCount"].(int64))
		record.TotalCount = b
		res.Records = append(res.Records, record)

	}

	responseMongoBody("Success", 0, w, res)

}

func responseMongoBody(message string, code int, w http.ResponseWriter, res mongoHandlerResponseBody) {
	res.Code = code
	res.Message = message
	json.NewEncoder(w).Encode(res)
}

func responseSetBody(key, value string, w http.ResponseWriter, res setHandlerResponseBody) {
	res.Key = key
	res.Value = value
	json.NewEncoder(w).Encode(res)
}
