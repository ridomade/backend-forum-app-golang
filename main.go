package main

import (
	"fmt"
	"log"
	"net/http"

	"firstproject/config"
	"firstproject/routes"
)

func main() {
	// Connect to database
	config.ConnectDB()

	// Routes
	routes.PostRoutes()
	routes.AccountRoutes()

	var port = "8080"

	// Listen and serve
	fmt.Println("Server Running on Port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	
}
