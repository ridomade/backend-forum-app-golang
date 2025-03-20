package controllers

import (
	"encoding/json"
	"firstproject/config"
	"firstproject/models"
	"fmt"
	"net/http"
)

// CreateItem menambahkan item baru ke database
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	// Decode JSON dari body request
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Gagal membaca data dari request", http.StatusBadRequest)
		return
	}

	// Simpan data ke database
	post.Account_id = 3
	query	:= "INSERT INTO posts (title, content,account_id) VALUES (?, ?, ?)"
	result, err := config.DB.Exec(query, post.Title, post.Content, 3)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
		return
	}

	// Ambil ID dari item yang baru dimasukkan
	id, _ := result.LastInsertId()

	post.ID = int(id)

	// Kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}