package main

import (
	"log"
	"net/http"
)

var logPrefix string = "[Supplier]"

const portString string = "8081"

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/findsupplier", findSupplier)

	log.Printf("%v Starting FoodSupplier server to listen for requests on port %v", logPrefix, portString)
	log.Printf("%v %v", logPrefix, http.ListenAndServe(":"+portString, nil))
}

func findSupplier(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL) // Log the request for now, this will be updated with actual logic
}
