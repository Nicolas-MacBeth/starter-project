package foodfinder

import (
	"fmt"
	"net/http"
	"sync"
)

var logPrefix string = "[FoodFinder] "

// StartFoodFinder initiates the server for this service
func StartFoodFinder(waitGroup *sync.WaitGroup) {

	// Handle routes with their respective functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/findfood", findFood)

	portString := "8080"

	fmt.Println(logPrefix + "Starting FoodFinder server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)

	// This line usually won't run, it's simply to keep main alive while the server runs in a goroutine
	waitGroup.Done()
}

func findFood(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}

func homepage(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, will serve the homepage
}
