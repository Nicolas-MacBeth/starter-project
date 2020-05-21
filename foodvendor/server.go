package main

import (
	"log"
	"net/http"
)

var logPrefix string = "[Vendor]"

const portString = "8082"

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/findvendor", findVendor)

	log.Printf("%v Starting FoodVendor server to listen for requests on port %v", logPrefix, portString)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}

func findVendor(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL) // Log the request for now, this will be updated with actual logic
}
