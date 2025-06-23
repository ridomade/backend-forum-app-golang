package routes

import (
	"encoding/json"
	"net/http"

	"firstproject/controllers"
	"firstproject/middleware"
)

// Ubah fungsi agar menerima parameter `mux`
func AccountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/forum/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateAccount(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/forum/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.LoginAccount(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func PostRoutes(mux *http.ServeMux) {

	// GET request tidak perlu middleware
	mux.HandleFunc("/forum/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetPostsWithReplies(w, r)
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

	mux.HandleFunc("/forum/post/reply", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			middleware.JWTMiddleware(http.HandlerFunc(controllers.CreatePostReply)).ServeHTTP(w, r)
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
