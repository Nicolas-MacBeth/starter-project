package main

import (
	"fmt"
	"net/http"
)

var logPrefix string = "[FoodFinder] "

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/findfood", findFood)

	portString := "8080"

	fmt.Println(logPrefix + "Starting FoodFinder server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)
}

func findFood(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}

func homepage(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, will serve the homepage
}
