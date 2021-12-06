package aircraft

// performance model
type Performance struct {
	Range             *float32 `json:"range,omitempty" bson:"range,omitempty"`
	CruiseSpeed       *int16   `json:"cruiseSpeed,omitempty" bson:"cruiseSpeed,omitempty"`
	MaxSpeed          *int16   `json:"maxSpeed,omitempty" bson:"maxSpeed,omitempty"`
	Ceiling           *float32 `json:"ceiling,omitempty" bson:"ceiling,omitempty"`
	MaxTakeoffWeight  *float32 `json:"maxTakeoffWeight,omitempty" bson:"maxTakeoffWeight,omitempty"`
	MaxLandingWeight  *float32 `json:"maxLandingWeight,omitempty" bson:"maxLandingWeight,omitempty"`
	MaxZeroFuelWeight *float32 `json:"maxZeroFuelWeight,omitempty" bson:"maxZeroFuelWeight,omitempty"`
	FuelCapacity      *float32 `json:"fuelCapacity,omitempty" bson:"fuelCapacity,omitempty"`
	TakeoffDistance   *float32 `json:"takeoffDistance,omitempty" bson:"takeoffDistance,omitempty"`
	Wingspan          *float32 `json:"wingspan,omitempty" bson:"wingspan,omitempty"`
}
