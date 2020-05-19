package foodsupplier

import (
	"fmt"
	"net/http"
	"sync"
)

var logPrefix string = "[FoodSupplier] "

// StartFoodSupplier initiates the server for this service
func StartFoodSupplier(waitGroup *sync.WaitGroup) {

	// Handle routes with their respective functions
	http.HandleFunc("/findsupplier", findSupplier)

	portString := "8081"

	fmt.Println(logPrefix + "Starting FoodSupplier server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)

	// This line usually won't run, it's simply to keep main alive while the server runs in a goroutine
	waitGroup.Done()
}

func findSupplier(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}