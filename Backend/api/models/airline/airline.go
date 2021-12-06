package airline

import "go.mongodb.org/mongo-driver/bson/primitive"

type Airline struct {
	ID      primitive.ObjectID   `json:"id" bson:"_id"`
	General *General             `json:"general,omitempty" bson:"general,omitempty"`
	Fleet   []primitive.ObjectID `json:"fleet" bson:"fleet"`
	Reviews []primitive.ObjectID `json:"reviews" bson:"reviews"`
	Routes  []primitive.ObjectID `json:"routes" bson:"routes"`
	Owner   primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
}
