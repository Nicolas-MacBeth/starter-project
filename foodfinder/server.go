package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var logPrefix string = "[Finder]"

type queryParams struct {
	IngredientsList []string
}
type vendorInfo struct {
	ID   int
	Name string
}

type vendorAndIngredientList struct {
	Vendors     []vendorInfo
	Ingredients []string
}

type inventoryAndPrice struct {
	VendorID   int
	VendorName string
	Ingredient string
	Price      int
	Inventory  int
}

type finalAnswers struct {
	ListOfResults []inventoryAndPrice
}

const portString string = "8080"
const supplierPort string = "8081"
const vendorPort string = "8082"

func main() {
	router := http.NewServeMux()

	// Handle routes with their respective functions
	homepageHandler := http.HandlerFunc(homepage)
	router.Handle("/", homepageHandler)
	findfoodHandler := http.HandlerFunc(findFood)
	router.Handle("/findfood", findfoodHandler)

	fs := http.FileServer(http.Dir("./public"))
	router.Handle("/public/", http.StripPrefix("/public/", fs))

	log.Printf("%v Starting FoodFinder server to listen for requests on port %v", logPrefix, portString)
	log.Printf("%v %v", logPrefix, http.ListenAndServe(":"+portString, router))
}

func findFood(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL)

	var ingredients queryParams

	// Try to decode JSON onto ingredients struct, send badRequest error if it fails
	err := json.NewDecoder(request.Body).Decode(&ingredients)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// Get list of vendors
	payload, err := findVendors(ingredients)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to contact food supplier API", http.StatusInternalServerError)
		return
	}

	// early 404 if no vendors were found
	if len(payload.Vendors) == 0 {
		log.Printf("%v No vendors found for this ingredient list: %v", logPrefix, payload.Ingredients)
		http.Error(response, "No vendors found for this ingredient", http.StatusNotFound)
		return
	}

	// Get vendor details
	result, err := queryVendorsPriceAndInventory(payload)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to contact food vendor API", http.StatusInternalServerError)
		return
	}

	// JSON response
	dataBody, err := json.Marshal(result)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(dataBody)
}

func homepage(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL)
	http.ServeFile(response, request, "public/index.html")
}

func findVendors(ingredients queryParams) (*vendorAndIngredientList, error) {
	body, err := json.Marshal(ingredients)
	if err != nil {
		return nil, err
	}

	// Make request to food supplier service
	response, err := http.Post(fmt.Sprintf("http://localhost:%v/foodsupplier", supplierPort), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Try to decode JSON onto payload
	var payload vendorAndIngredientList

	jsonErr := json.NewDecoder(response.Body).Decode(&payload)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &payload, nil
}

func queryVendorsPriceAndInventory(payload *vendorAndIngredientList) (*finalAnswers, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Make request to food vendor service
	response, err := http.Post(fmt.Sprintf("http://localhost:%v/foodvendor", vendorPort), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Try to decode JSON onto finalResult
	var finalResult finalAnswers

	jsonErr := json.NewDecoder(response.Body).Decode(&finalResult)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &finalResult, nil
}
