package flight

import "go.mongodb.org/mongo-driver/bson/primitive"

type Flight struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	FlightNumber        string             `json:"flightNumber,omitempty" bson:"flightNumber,omitempty"`
	Callsign            string             `json:"callsign,omitempty" bson:"callsign,omitempty"`
	Departure           primitive.ObjectID `json:"departure,omitempty" bson:"departure,omitempty"`
	Arrival             primitive.ObjectID `json:"arrival,omitempty" bson:"arrival,omitempty"`
	Distance            string             `json:"distance,omitempty" bson:"distance,omitempty"`
	FlightTime          string             `json:"flightTime,omitempty" bson:"flightTime,omitempty"`
	AverageArrivalDelay string             `json:"averageArrivalDelay,omitempty" bson:"averageArrivalDelay,omitempty"`
	DepartureTime       map[string]string  `json:"departureTime" bson:"departureTime"`
	ArrivalTime         map[string]string  `json:"arrivalTime" bson:"arrivalTime"`
	Airline             primitive.ObjectID `json:"airline,omitempty" bson:"airline,omitempty"`
	Aircraft            primitive.ObjectID `json:"aircraft,omitempty" bson:"aircraft,omitempty"`
}
