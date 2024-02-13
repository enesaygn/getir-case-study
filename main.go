package main

import (
	db "getir-case/db"
	endpoints "getir-case/endpoints"
	"log"
	"net/http"
	"os"
)

func main() {

	// Initialize the in-memory database
	err := db.InitializeMongoDB()
	if err != nil {
		log.Println("Error initializing the in-memory database: " + err.Error())
		os.Exit(1)
	}

	// Register the handlers
	http.HandleFunc("/inmemoryget", endpoints.InMemoryGetHandler)
	http.HandleFunc("/inmemoryset", endpoints.InMemorySetHandler)
	http.HandleFunc("/mongofetch", endpoints.MongoHandler)

	// Get the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		// Log the error and exit the program
		log.Println("PORT environment variable not set")
		os.Exit(1)
	}

	// Start the server
	println("Server started at PORT: " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		println("Error starting the server: " + err.Error())
		os.Exit(1)
	}

}
