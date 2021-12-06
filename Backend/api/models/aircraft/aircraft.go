package aircraft

import (
	trackerdata "github.com/arttkachev/X-Airlines/Backend/api/models/trackerData"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// aircraft model
type Aircraft struct {
	ID          primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	General     *General                 `json:"general,omitempty" bson:"general,omitempty"`
	Airframe    *Airframe                `json:"airframe,omitempty" bson:"airframe,omitempty"`
	Engines     []primitive.ObjectID     `json:"engines" bson:"engines"`
	Exterior    *Exterior                `json:"exterior,omitempty" bson:"exterior,omitempty"`
	Interior    *Interior                `json:"interior,omitempty" bson:"interior,omitempty"`
	Cockpit     *Cockpit                 `json:"cockpit,omitempty" bson:"cockpit,omitempty"`
	Performance *Performance             `json:"performance,omitempty" bson:"performance,omitempty"`
	TrackerData *trackerdata.TrackerData `json:"trackerData,omitempty" bson:"trackerData,omitempty"`
	Owner       primitive.ObjectID       `json:"owner,omitempty" bson:"owner,omitempty"`
	Tags        []string                 `json:"tags" bson:"tags"`
}
