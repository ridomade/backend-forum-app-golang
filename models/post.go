package models

// Item merepresentasikan tabel items di database
type Post struct {
	ID   int    `json:"id"`
	Account_id int `json:"account_id"`
	Title string `json:"title"`
	Content string `json:"content"`
}
