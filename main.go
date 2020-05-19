package main

import (
	"fmt"
	"starter-project/foodfinder"
	"starter-project/foodsupplier"
	"starter-project/foodvendor"
	"sync"
)

func main() {
	fmt.Println("Welcome to the suite of ingredient finding services")

	// Waitgroup here is needed for main to "wait" on the servers, if not main would complete and shut down the goroutines
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(3)

	go foodfinder.StartFoodFinder(waitGroup)
	go foodsupplier.StartFoodSupplier(waitGroup)
	go foodvendor.StartFoodVendor(waitGroup)

	waitGroup.Wait()
}
