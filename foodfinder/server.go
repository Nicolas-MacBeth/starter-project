package main

import (
	"log"
	"net/http"
)

var logPrefix string = "[FoodFinder]"

const portString string = "8080"

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/findfood", findFood)

	log.Printf("%v Starting FoodFinder server to listen for requests on port %v", logPrefix, portString)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}

func findFood(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL) // Log the request for now, this will be updated with actual logic
}

func homepage(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL)
	http.ServeFile(response, request, "public/homepage.html")
}
