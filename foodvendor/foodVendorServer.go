package foodvendor

import (
	"fmt"
	"net/http"
	"sync"
)

var logPrefix string = "[FoodVendor] "

// StartFoodVendor initiates the server for this service
func StartFoodVendor(waitGroup *sync.WaitGroup) {

	// Handle routes with their respective functions
	http.HandleFunc("/findvendor", findVendor)

	portString := "8082"

	fmt.Println(logPrefix + "Starting FoodVendor server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)

	// This line usually won't run, it's simply to keep main alive while the server runs in a goroutine
	waitGroup.Done()
}

func findVendor(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}
