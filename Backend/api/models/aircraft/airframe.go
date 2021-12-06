package aircraft

// airframe model
type Airframe struct {
	TotalTime     *uint16 `json:"totalTime,omitempty" bson:"totalTime,omitempty"`
	TotalLandings *uint16 `json:"totalLandings,omitempty" bson:"totalLandings,omitempty"`
	AirframeNotes string  `json:"airframeNotes,omitempty" bson:"airframeNotes,omitempty"`
}
