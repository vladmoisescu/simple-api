package main

import (
	"net/http"
	"os"
	"simple-api/database"
	"simple-api/routes"
)

func main() {
	database.SetupConnection()

	routes.Setup()

	//p1 := models.Product{
	//	Name:  "Product 1",
	//	Price: 100,
	//	Stock: 5,
	//}
	//p2 := models.Product{
	//	Name:  "Product 2",
	//	Price: 200,
	//	Stock: 5,
	//}
	//
	//database.DB.Create(&p1)
	//database.DB.Create(&p2)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
