package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var logPrefix string = "[Vendor]"
var db *sql.DB

const portString = "8082"

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
	http.HandleFunc("/foodvendor", foodVendor)

	log.Printf("%v Starting FoodVendor server to listen for requests on port %v", logPrefix, portString)
	log.Printf("%v %v", logPrefix, http.ListenAndServe(":"+portString, nil))
}

func foodVendor(response http.ResponseWriter, request *http.Request) {
	log.Printf("%v %v %v", logPrefix, request.Method, request.URL) // Log the request for now, this will be updated with actual logic

	var payload vendorAndIngredientList

	// Try to decode, send badRequest error if it fails
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// DB operation
	finalInfo, err := queryDB(payload)

	// JSON response
	dataBody, err := json.Marshal(finalAnswers{ListOfResults: finalInfo})
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		http.Error(response, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(dataBody)
}

func queryDB(payload vendorAndIngredientList) ([]inventoryAndPrice, error) {
	var list []inventoryAndPrice

	// Query DB with constructed query and spreadable arguments
	rows, err := db.Query(constructQuery(payload), makeSpreadable(payload)...)
	if err != nil {
		log.Printf("%v %v", logPrefix, err)
		return nil, err
	}

	// Handle returned data, aggregate returned DB rows into a slice of inventoryAndPrice structs
	for rows.Next() {
		var row inventoryAndPrice

		// Assign proper values to vendorInfo struct
		err := rows.Scan(&row.VendorID, &row.Ingredient, &row.Price, &row.Inventory, &row.VendorName)
		if err != nil {
			log.Printf("%v %v", logPrefix, err)
			return nil, err
		}
		list = append(list, row)
	}
	return list, nil
}

func constructQuery(payload vendorAndIngredientList) string {
	// Base of mySQL query
	query := "SELECT * FROM vendor_ingredient WHERE ("

	// Construct query with all vendors as placeholders
	for i := range payload.Vendors {
		if i == 0 {
			query = fmt.Sprintf("%v vendor_id = ?", query)
		} else {
			query = fmt.Sprintf("%v OR vendor_id = ?", query)
		}
	} // This vendor part of the query is actually useless with how I built the database (actually this entire service/server is),
	// but I made it this way to follow the requirements for the starter project

	// Construct query with all ingredients as placeholders
	query = fmt.Sprintf("%v) AND (", query)
	for i := range payload.Ingredients {
		if i == 0 {
			query = fmt.Sprintf("%v ingredient_name = ?", query)
		} else {
			query = fmt.Sprintf("%v OR ingredient_name = ?", query)
		}
	}
	return fmt.Sprintf("%v) ORDER BY ingredient_name ASC, price ASC", query)
}

func makeSpreadable(payload vendorAndIngredientList) []interface{} {
	// This is needed to spread the IngredientsList slice in the Query() function for safe mySQL queries (prepared statements)
	// Convert all the values I need in the mySQl query into a single slice of interfaces
	spreadable := make([]interface{}, len(payload.Ingredients)+len(payload.Vendors))

	for i, value := range payload.Vendors {
		spreadable[i] = value.ID
	}
	for i, value := range payload.Ingredients {
		spreadable[len(payload.Vendors)+i] = value
	}
	return spreadable
}
