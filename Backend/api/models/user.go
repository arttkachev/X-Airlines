package models

// user model
type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	IsAdmin  bool     `json:"isAdmin"`
	Balance  int      `json:"balance"`
	Airlines []string `json:"airlines"`
}
