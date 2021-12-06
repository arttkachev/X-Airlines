package aircraft

type APU struct {
	TotalTime uint16 `json:"totalTime,omitempty" bson:"totalTime,omitempty"`
	Notes     string `json:"notes,omitempty" bson:"notes,omitempty"`
}
