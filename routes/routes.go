package routes

import (
	"encoding/json"
	"net/http"

	"firstproject/controllers"
	"firstproject/middleware"
)


func AccountRoutes() {
	http.HandleFunc("/forum/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateAccount(w, r)
		} else {
			http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/forum/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.LoginAccount(w, r)
		} else {
			http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
		}
	})
}

func PostRoutes() {
	// GET request tidak perlu middleware
	http.HandleFunc("/forum/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetPosts(w, r)
			return
		}

		// Jika bukan GET, periksa apakah method adalah POST
		if r.Method == http.MethodPost {
			middleware.JWTMiddleware(http.HandlerFunc(controllers.CreatePost)).ServeHTTP(w, r)
			return
		}

		// Jika method bukan GET atau POST
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method Not Allowed",
		})
	})
}

// func PostRoutes() {
// 	http.Handle("/posts", middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			controllers.CreatePost(w, r)
// 		} else if (r.Method == http.MethodGet) {
// 			controllers.GetPosts(w, r)
// 		} else {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusUnauthorized)
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"message": "Method Not Allowed",
// 			})
// 		}
// 	})))
// }