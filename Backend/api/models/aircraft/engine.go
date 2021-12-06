package aircraft

import "go.mongodb.org/mongo-driver/bson/primitive"

// engine model
type Engine struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwningAircraft primitive.ObjectID `json:"owningAircraft,omitempty" bson:"owningAircraft,omitempty"`
	Model          string             `json:"model,omitempty" bson:"model,omitempty"`
	TotalTime      *uint16            `json:"totalTime,omitempty" bson:"totalTime,omitempty"`
	TBO            *uint16            `json:"tbo,omitempty" bson:"tbo,omitempty"`
	HST            *uint16            `json:"hst,omitempty" bson:"hst,omitempty"`
}
