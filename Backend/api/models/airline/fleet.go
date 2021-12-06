package airline

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fleet struct {
	Fleet []primitive.ObjectID `json:"fleet" bson:"fleet"`
}
