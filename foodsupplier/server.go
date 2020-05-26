package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var logPrefix string = "[Supplier]"
var db *sql.DB

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

const portString string = "8081"

func main() {

	// Create connection to db, defer Close()
	var err error
	db, err = sql.Open("mysql", "starter-project:findMeSomeFoodPlz@/starterproject")
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		return
	}
	defer db.Close()

	// Handle routes with their respective functions
	http.HandleFunc("/foodsupplier", findVendor)

	log.Printf("%v Starting FoodSupplier server to listen for requests on port %v", logPrefix, portString)
	log.Printf("%v %v", logPrefix, http.ListenAndServe(":"+portString, nil))
}

func findVendor(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL) // Log the request for now, this will be updated with actual logic

	var ingredients queryParams

	// Try to decode JSON onto ingredients, send badRequest error if it fails
	err := json.NewDecoder(request.Body).Decode(&ingredients)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// DB operation
	vendorsSlice, err := queryDB(ingredients)
	if err != nil {
		http.Error(response, "Unable to query DB", http.StatusInternalServerError)
		return
	}

	// JSON response
	payload, err := json.Marshal(vendorAndIngredientList{
		Vendors:     vendorsSlice,
		Ingredients: ingredients.IngredientsList,
	})
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(payload)
}

func queryDB(ingredients queryParams) ([]vendorInfo, error) {
	// This is needed to *spread* the IngredientsList slice argument in the Query() function for safe mySQL queries (prepared statements)
	spreadable := make([]interface{}, len(ingredients.IngredientsList))
	for i, value := range ingredients.IngredientsList {
		spreadable[i] = value
	}

	// Query the db, with constructed query
	rows, err := db.Query(constructQuery(ingredients), spreadable...)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		return nil, err
	}

	var vendors []vendorInfo
	// Handle returned data, aggregate returned rows into a slice of structs
	for rows.Next() {
		var vendor vendorInfo

		// Assign proper values to vendorInfo struct
		err2 := rows.Scan(&vendor.ID, &vendor.Name)
		if err2 != nil {
			log.Printf("%v %v", logPrefix, err)
			return nil, err
		}
		vendors = append(vendors, vendor)
	}
	return vendors, nil
}

func constructQuery(ingredients queryParams) string {
	// Base of mySQL query
	query := "SELECT * FROM vendor WHERE id IN (SELECT vendor_id FROM vendor_ingredient WHERE"

	// Construct query with all ingredients as placeholders (for safe querying - prepared statements)
	for i := range ingredients.IngredientsList {
		if i == 0 {
			query = fmt.Sprintf("%v ingredient_name = ?", query)
		} else {
			query = fmt.Sprintf("%v OR ingredient_name = ?", query)
		}
	}
	return fmt.Sprintf("%v)", query)
}
