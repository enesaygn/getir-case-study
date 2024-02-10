package main

import (
	db "getir-case/db"
	endpoints "getir-case/endpoints"
	"net/http"
)

func main() {

	// Initialize the in-memory database
	db.InitializeMongoDb()

	// Register the handlers
	http.HandleFunc("/get", endpoints.GetHandler)
	http.HandleFunc("/post", endpoints.PostHandler)
	http.HandleFunc("/mongofetch", endpoints.MongoHandler)

	// Start the server
	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
