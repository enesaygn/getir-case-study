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

type mongoHandlerRequestBody struct {
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

type PostHandlerRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostHandlerResponseBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetHandlerRequestBody struct {
	Key string `json:"key"`
}

type GetHandlerResponseBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// getHandler handles GET requests to the root URL
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method == http.MethodGet {
		// Respond to GET requests
		w.Header().Set("Content-Type", "text/plain")
		queryParams := r.URL.Query()

		// Retrieve specific parameter values
		key := queryParams.Get("key")
		res := GetHandlerResponseBody{}

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

	} else {
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// postHandler handles POST requests to the /post URL
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method == http.MethodPost {
		// Respond to POST requests
		w.Header().Set("Content-Type", "application/json")

		req := PostHandlerRequestBody{}
		res := PostHandlerResponseBody{}
		inMemoryInstance := db.GetInMemoryDBInstance()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		inMemoryInstance.Set(req.Key, req.Value)
		responsePostBody(req.Key, req.Value, w, res)

	} else {
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// postHandler handles POST requests to the /post URL
func MongoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Respond to POST requests
		w.Header().Set("Content-Type", "text/plain")

		req := mongoHandlerRequestBody{}
		res := mongoHandlerResponseBody{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			responseBody(err.Error(), 1, w, res)
			return
		}

		layout := "2006-01-02"

		startDate, err := time.Parse(layout, req.StartDate)
		if err != nil {
			responseBody(err.Error(), 1, w, res)
			return
		}
		endDate, err := time.Parse(layout, req.EndDate)
		if err != nil {
			responseBody(err.Error(), 1, w, res)
			return
		}

		collection := db.Connection.Collection("records")

		// Sorgu için filtre oluşturma
		matchStage1 := bson.D{ //first match for range date time
			{Key: "$match", Value: bson.D{
				{Key: "createdAt", Value: bson.D{
					{Key: "$lte", Value: primitive.NewDateTimeFromTime(endDate)},
					{Key: "$gte", Value: primitive.NewDateTimeFromTime(startDate)}},
				},
			},
			},
		}
		projectStage := bson.D{ //sum counts as totalCount
			{Key: "$project", Value: bson.D{
				{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: "$counts"}}},
				{Key: "_id", Value: 1},
				{Key: "key", Value: 1},
				{Key: "createdAt", Value: 1},
			},
			},
		}

		matchStage2 := bson.D{ //second match stage for total count
			{Key: "$match", Value: bson.D{
				{Key: "totalCount", Value: bson.D{
					{Key: "$lte", Value: req.MaxCount},
					{Key: "$gte", Value: req.MinCount}}},
			},
			},
		}
		pipeline := mongo.Pipeline{matchStage1, projectStage, matchStage2}

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

		responseBody("Success", 0, w, res)

	} else {
		// Respond with a method not allowed error for other request methods
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func responseBody(message string, code int, w http.ResponseWriter, res mongoHandlerResponseBody) {
	res.Code = code
	res.Message = message
	json.NewEncoder(w).Encode(res)
}

func responsePostBody(key, value string, w http.ResponseWriter, res PostHandlerResponseBody) {
	res.Key = key
	res.Value = value
	json.NewEncoder(w).Encode(res)
}
