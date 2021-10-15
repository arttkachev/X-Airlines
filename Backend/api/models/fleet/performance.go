package fleet

// performance model
type Performance struct {
	Range             float32 `json:"range"`
	CruiseSpeed       int16   `json:"cruiseSpeed"`
	MaxSpeed          int16   `json:"maxSpeed"`
	Ceiling           float32 `json:"ceiling"`
	MaxTakeoffWeight  float32 `json:"maxTakeoffWeight"`
	MaxLandingWeight  float32 `json:"maxLandingWeight"`
	MaxZeroFuelWeight float32 `json:"maxZeroFuelWeight"`
	FuelCapacity      float32 `json:"fuelCapacity"`
	TakeoffDistance   float32 `json:"takeoffDistance"`
	Wingspan          float32 `json:"wingspan"`
}
