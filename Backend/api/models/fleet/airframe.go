package fleet

// airframe model
type Airframe struct {
	TotalTime     uint16 `json:"totalTime"`
	TotalLandings uint16 `json:"totalLandings"`
	AirframeNotes string `json:"airframeNotes"`
}
