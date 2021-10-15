package fleet

// engine model
type Engine struct {
	Model              string `json:"model"`
	TotalTime          uint16 `json:"totalTime"`
	TBO                uint16 `json:"tbo"`
	HST                uint16 `json:"hst"`
	MaintenanceProgram string `json:"maintenanceProgram"`
	Notes              string `json:"notes"`
}
