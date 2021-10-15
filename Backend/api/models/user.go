package models

// user model
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	IsAdmin  bool     `json:"isAdmin"`
	Balance  int      `json:"balance"`
	Fleet    []string `json:"fleet"`
}
