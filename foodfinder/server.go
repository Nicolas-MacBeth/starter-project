package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
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

// metrics variables needed
var (
	requests      = stats.Int64("starter-project_requests", "The total requests per ingredient", "reqs")
	ingredient, _ = tag.NewKey("ingredient")
)

const portString string = "8080"
const supplierPort string = "8081"
const vendorPort string = "8082"
const prometheusPort string = "8083"
const counterPort = "8084"

func main() {
	// Create the Prometheus exporter.
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "ocmetrics",
	})
	if err != nil {
		log.Printf("%v Failed to create the Prometheus metrics exporter: %v", logPrefix, err)
	}

	// Run the Prometheus exporter as a scrape endpoint.
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exporter)
		log.Printf("%v Starting Prometheus scrape endpoint on port: %v", logPrefix, prometheusPort)
		log.Println(http.ListenAndServe(":"+prometheusPort, mux))
	}()
	// foodfinder Server
	router := http.NewServeMux()

	// Handler for routes with their respective functions, wrapped in Metrics handler
	homepageHandler := http.HandlerFunc(homepage)
	router.Handle("/", homepageHandler)

	findfoodHandler := http.HandlerFunc(findFood)
	router.Handle("/findfood", findfoodHandler)

	fs := http.FileServer(http.Dir("./public"))
	router.Handle("/public/", http.StripPrefix("/public/", fs))

	// Wrap the router in a ochttp Handler for metrics
	mux := &ochttp.Handler{
		Handler: router,
	}

	v := &view.View{
		Name:        "ingredient_specific_requests",
		Measure:     requests,
		Description: "Requests per ingredient",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{ingredient},
	}
	if err := view.Register(v); err != nil {
		log.Fatalf("%v Failed to register the view: %v", logPrefix, err)
	}

	log.Printf("%v Starting FoodFinder server to listen for requests on port %v", logPrefix, portString)
	log.Printf("%v %v", logPrefix, http.ListenAndServe(":"+portString, mux))
}

func findFood(response http.ResponseWriter, request *http.Request) {
	go notifyJavaServer()

	log.Printf("%v %v %v", logPrefix, request.Method, request.URL)

	var ingredients queryParams

	// Try to decode JSON onto ingredients struct, send badRequest error if it fails
	err := json.NewDecoder(request.Body).Decode(&ingredients)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// record metrics for every ingredient
	for _, ingredient := range ingredients.IngredientsList {
		err := recordIngredient(ingredient)
		if err != nil {
			log.Printf("%v %v", logPrefix, err)
		}
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

func notifyJavaServer() {
	response, err := http.Get(fmt.Sprintf("http://localhost:%v/count", counterPort))
	if err != nil {
		fmt.Printf("Java notifying error, err: %v", err)
	}
	defer response.Body.Close()
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

func recordIngredient(ingredientString string) error {
	ctx, err := tag.New(context.Background(), tag.Upsert(ingredient, ingredientString))
	if err != nil {
		return err
	}
	stats.Record(ctx, requests.M(1))

	return nil
}
