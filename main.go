package main

import (
	"log"
	"net/http"
	"os"

	"simple-api/database"
	"simple-api/routes"
)

func main() {
	dbDns := os.Getenv("DB_DNS")
	if dbDns == "" {
		log.Fatal(1)
	}
	database.SetupConnection(dbDns)

	routes.Setup()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
