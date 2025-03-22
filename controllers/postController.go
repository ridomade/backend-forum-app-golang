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

	// Check if author and content is empty
	if post.Author == "" || post.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "author and content cannot be empty",
		})
		return
	}

	// Get claims from JWT
	claims, ok := r.Context().Value(middleware.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get  ID from claims
	idFloat, idExists := claims["id"].(float64) // ID is float64 in JWT claims

	if !idExists {
		http.Error(w, `{"message": "Unauthorized: Missing  ID"}`, http.StatusUnauthorized)
		return
	}

	// Insert data to database
	query := "INSERT INTO posts (author, content, account_id) VALUES (?, ?, ?)"
	result, err := config.DB.Exec(query, post.Author, post.Content, int(idFloat))
	if err != nil {
		http.Error(w, `{"message": "Error Adding Data to the Database"}`, http.StatusInternalServerError)
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

func CreatePostReply(w http.ResponseWriter, r *http.Request) {
	var postReply models.PostReply

	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&postReply)
	if err != nil {
		http.Error(w, `{"message": "Error Ready data From the Request"}`, http.StatusBadRequest)
		return
	}

	// Check if author and content and post_parent is empty
	if postReply.Author == "" || postReply.Content == "" || postReply.PostParent == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "author, post_parent, and content cannot be empty",
		})
		return
	}

	// Get claims from JWT
	claims, ok := r.Context().Value(middleware.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get  ID from claims
	idFloat, idExists := claims["id"].(float64) // ID is float64 in JWT claims

	if !idExists {
		http.Error(w, `{"message": "Unauthorized: Missing  ID"}`, http.StatusUnauthorized)
		return
	}

	// Insert data to database
	query := "INSERT INTO post_reply (post_parent, account_id, author, content) VALUES (?, ?, ?, ?)"
	result, err := config.DB.Exec(query, postReply.PostParent, int(idFloat), postReply.Author, postReply.Content)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"message": "Error Adding Data to the Database"}`, http.StatusInternalServerError)
		return
	}

	// Get ID from last inserted data
	lastInsertID, _ := result.LastInsertId()
	postReply.ID = int(lastInsertID)
	postReply.Account_id = int(idFloat)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
		"post":    postReply,
	})
}

func GetPostReplies(w http.ResponseWriter, r *http.Request) {
	var postReplies []models.PostReply

	// Get data from database
	rows, err := config.DB.Query("SELECT * FROM post_reply WHERE post_parent = ?", r.URL.Query().Get("post_parent"))
	if err != nil {
		fmt.Println("Error Querying Data : \n",err)
		http.Error(w, `{"message": "Error getting the data"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over rows
	for rows.Next() {
		var postReply models.PostReply
		err := rows.Scan(&postReply.ID, &postReply.Account_id, &postReply.PostParent, &postReply.Author, &postReply.Content)
		if err != nil {
			fmt.Println("Error Querying Data : \n",err)
			http.Error(w, `{"message": "Error getting the data"}`, http.StatusInternalServerError)
			return
		}
		postReplies = append(postReplies, postReply)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"post_replies": postReplies,
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
		err := rows.Scan(&post.ID, &post.Account_id, &post.Author, &post.Content , &post.Created_at, &post.Update_at)
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