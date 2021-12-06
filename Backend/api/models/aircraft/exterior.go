package aircraft

type Exterior struct {
	YearPainted *uint16 `json:"yearPainted,omitempty" bson:"yearPainted,omitempty"`
	Notes       string  `json:"notes,omitempty" bson:"notes,omitempty"`
}
