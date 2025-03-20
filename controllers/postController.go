package controllers

import (
	"encoding/json"
	"firstproject/config"
	"firstproject/middleware"
	"firstproject/models"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// CreatePost function to create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	// Decode JSON from body request
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, `{"message": "Gagal membaca data dari request"}`, http.StatusBadRequest)
		return
	}

	// Check if title and content is empty
	if post.Title == "" || post.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Title and Content cannot be empty",
		})
		return
	}

	// Get claims from JWT
	claims, ok := r.Context().Value(middleware.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get email and ID from claims
	email, emailExists := claims["email"].(string)
	idFloat, idExists := claims["id"].(float64) // ID is float64 in JWT claims

	if !emailExists || !idExists {
		http.Error(w, `{"message": "Unauthorized: Missing email or ID"}`, http.StatusUnauthorized)
		return
	}

	fmt.Println("Email:", email)
	fmt.Println("ID:", idFloat)

	// Insert data to database
	query := "INSERT INTO posts (title, content, account_id) VALUES (?, ?, ?)"
	result, err := config.DB.Exec(query, post.Title, post.Content, int(idFloat))
	if err != nil {
		http.Error(w, `{"message": "Gagal menyimpan data"}`, http.StatusInternalServerError)
		return
	}

	// Get ID from last inserted data
	lastInsertID, _ := result.LastInsertId()
	post.ID = int(lastInsertID)
	post.Account_id = int(idFloat)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
		"post":    post,
	})
}


// GetPosts function to get all posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	// Get data from database
	rows, err := config.DB.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Println("Error Querying Data : \n",err)
		http.Error(w, `{"message": "Gagal mengambil data"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over rows
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Account_id, &post.Title, &post.Content , &post.Created_at, &post.Update_at)
		if err != nil {
			fmt.Println("Error Querying Data : \n",err)
			http.Error(w, `{"message": "Gagal mengambil data"}`, http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"posts":   posts,
	})
}