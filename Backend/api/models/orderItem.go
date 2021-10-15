package models

// orderItem model
type OrderItem struct {
	Item     string `json:"item"` // TODO: item object id
	Quantity uint8  `json:"quantity"`
}
