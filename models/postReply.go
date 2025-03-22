package models

// Item merepresentasikan tabel items di database
type PostReply struct {
	ID   int    `json:"id"`
	Account_id int `json:"account_id"`
	PostParent int `json:"post_parent"`
	Author string `json:"author"`
	Content string `json:"content"`
}
