package models

type Post struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
