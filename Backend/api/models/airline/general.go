package airline

type General struct {
	Name   string  `json:"name,omitempty" bson:"name,omitempty"`
	Logo   string  `json:"logo,omitempty" bson:"logo,omitempty"`
	IATA   string  `json:"iata,omitempty" bson:"iata,omitempty"`
	ICAO   string  `json:"icao,omitempty" bson:"icao,omitempty"`
	Fleet  *uint16 `json:"fleet,omitempty" bson:"fleet,omitempty"`
	Rating *uint8  `json:"rating,omitempty" bson:"rating,omitempty"`
}
