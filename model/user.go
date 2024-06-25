package model

type User struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	UserType int    `json:"user_type"`
}
