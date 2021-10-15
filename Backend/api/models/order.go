package models

// order model
type Order struct {
	OrderItems  []string `json:"orderItems"` // TODO: orderItems object ids
	User        string   `json:"user"`       // TODO: user object id
	DateOfOrder string   `json:"dateOfOrder"`
	Status      string   `json:"status"`
	TotalPrice  int      `json:"totalPrice"`
}
