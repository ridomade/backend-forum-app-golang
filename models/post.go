package models

// Item merepresentasikan tabel items di database
type Post struct {
	ID         int         `json:"id"`
	Account_id int         `json:"account_id"`
	Author     string      `json:"author"`
	Content    string      `json:"content"`
	Created_at string      `json:"created_at,omitempty"`
	Update_at  string      `json:"update_at,omitempty"`
	Replies    []PostReply `json:"replies"`
}
