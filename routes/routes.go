package routes

import (
	"encoding/json"
	"net/http"

	"firstproject/controllers"
	"firstproject/middleware"
)


func AccountRoutes() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateAccount(w, r)
		} else {
			http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.LoginAccount(w, r)
		} else {
			http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
		}
	})
}

func PostRoutes() {
	http.Handle("/posts", middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreatePost(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Method Not Allowed",
			})
		}
	})))
}