package aircraft

type Interior struct {
	YearInterior  *uint16 `json:"yearInterior,omitempty" bson:"yearInterior,omitempty"`
	NumberOfSeats *uint16 `json:"numberOfSeats,omitempty" bson:"numberOfSeats,omitempty"`
}
