package fleet

// aircraft general info model
type General struct {
	Year         uint16 `json:"year"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	SerialNumber string `json:"serialNumber"`
	Registration string `json:"registration"`
	Condition    string `json:"condition"`
	Description  string `json:"description"`
	Location     string `json:"aircraftLocation"`
	IsOperating  bool   `json:"isOperating"`
}
