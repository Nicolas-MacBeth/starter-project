# starter-project

Starter project in Go to increase familiarity with GCP and OpenTelemetry. It is composed of three services: *FoodFinder*, *FoodSupplier*, *FoodVendor*. They will interact together as a suite of products to give a user the availability/vendors of a certain ingredient. Each service has its own server and some services will also need their own mySQL DB table(s). 

It is made on purpose and according to requirements to have the Supplier service return a list of Vendors with the ingredient(s), and then only fetch their information from the Vendor service separately. (I understand this entire project could be implement with one server, however part of the exercice is orchestrating the communication and then setting up telemetry/metrics with OTel for multiple services).

## How to run

1. Clone repo
2. Install [Go](https://golang.org/dl/)
3. Run `cd starter-project` and `go build main.go && ./main`

Created by Nicolas MacBeth, as a starter project for my internship @ Google
