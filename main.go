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

	// Mux Router
	mux := http.NewServeMux()

	// Register Routes
	routes.PostRoutes(mux)
	routes.AccountRoutes(mux)

	port := "8080"

	// Wrap dengan CORS middleware
	handler := config.EnableCORS(mux)

	// Start server
	fmt.Println("Server Running on Port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
