package main

import (
	"fmt"
	"net/http"
)

var logPrefix string = "[FoodVendor] "

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/findvendor", findVendor)

	portString := "8082"

	fmt.Println(logPrefix + "Starting FoodVendor server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)
}

func findVendor(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}
