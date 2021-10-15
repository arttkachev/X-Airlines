package fleet

type Interior struct {
	YearInterior          uint16 `json:"yearInterior"`
	NumberOfSeats         uint16 `json:"numberOfSeats"`
	LavatoryConfiguration string `json:"lavatoryConfiguration"`
	Notes                 string `json:"notes"`
}
