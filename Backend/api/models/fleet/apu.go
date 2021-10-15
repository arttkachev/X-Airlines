package fleet

type APU struct {
	TotalTime          uint16 `json:"totalTime"`
	MaintenanceProgram string `json:"maintenanceProgram"`
	Notes              string `json:"notes"`
}
