package main

import (
	"fmt"
	"net/http"
)

var logPrefix string = "[FoodSupplier] "

func main() {

	// Handle routes with their respective functions
	http.HandleFunc("/findsupplier", findSupplier)

	portString := "8081"

	fmt.Println(logPrefix + "Starting FoodSupplier server to listen for requests on port " + portString)
	http.ListenAndServe(":"+portString, nil)
}

func findSupplier(response http.ResponseWriter, req *http.Request) {
	fmt.Print(logPrefix)
	fmt.Println(*req) // Log the request for now, this will be updated with actual logic
}
