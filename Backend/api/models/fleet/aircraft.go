package fleet

import "go.mongodb.org/mongo-driver/bson/primitive"

// aircraft model
type Aircraft struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name`
	//	Icon                string               `json:"icon" bson:"icon`
	//  General primitive.ObjectID `json:"general" bson:"general`
	//	Airframe            primitive.ObjectID   `json:"airframe" bson:"airframe`
	//	Engines             []primitive.ObjectID `json:"engines" bson:"engines`
	//	APU                 primitive.ObjectID   `json:"apu" bson:"apu`
	//	Avionics            primitive.ObjectID   `json:"avionics" bson:"avionics`
	Exterior *Exterior `json:"exterior" bson:"exterior`
	//	Interior            primitive.ObjectID   `json:"interior" bson:"interior`
	//	Performance         primitive.ObjectID   `json:"performance" bson:"performance`
	//	AdditionalEquipment primitive.ObjectID   `json:"additionalEquipment" bson:"additionalEquipment`
	//	Price               float32              `json:"price" bson:"price`
	//	Seller              primitive.ObjectID   `json:"sellerInfo" bson:"seller`
	//	Tags                []string             `json:"tags" bson:"tags`

}
