package controllers

import (
	"encoding/json"
	"firstproject/config"
	"firstproject/middleware"
	"firstproject/models"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

// CreatePost function to create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	// Decode JSON from body request
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		// Check if author and content is empty
		if post.Author == "" || post.Content == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "author and content cannot be empty",
			})
			return
		}
		fmt.Println("Error Reading Data from Request :", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Reading Data from Request",
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
		// Check if author and content and post_parent is empty
		if postReply.Author == "" || postReply.Content == "" || postReply.PostParent == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "author, post_parent, and content cannot be empty",
			})
			return
		}
		fmt.Println("Error Reading Data from Request :", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Reading Data from Request",
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

func GetPostsWithReplies(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 5

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	// Dapatkan total jumlah post
	var totalPosts int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalPosts)
	if err != nil {
		http.Error(w, `{"message": "Failed to get total post count"}`, http.StatusInternalServerError)
		return
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalPosts) / float64(limit)))

	// Ambil data post
	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?"
	postRows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		http.Error(w, `{"message": "Failed to fetch posts"}`, http.StatusInternalServerError)
		return
	}
	defer postRows.Close()

	var posts []models.Post

	for postRows.Next() {
		var post models.Post
		err := postRows.Scan(&post.ID, &post.Account_id, &post.Author, &post.Content, &post.Created_at, &post.Update_at)
		if err != nil {
			continue
		}

		var replies []models.PostReply
		replyRows, _ := config.DB.Query("SELECT * FROM post_reply WHERE post_parent = ?", post.ID)
		for replyRows.Next() {
			var reply models.PostReply
			_ = replyRows.Scan(&reply.ID, &reply.Account_id, &reply.PostParent, &reply.Author, &reply.Content)
			replies = append(replies, reply)
		}
		replyRows.Close()

		post.Replies = replies
		posts = append(posts, post)
	}

	// Kirim response dengan pagination info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Success",
		"page":       page,
		"limit":      limit,
		"totalPosts": totalPosts,
		"totalPages": totalPages,
		"posts":      posts,
	})
}
