package fleet

// aircraft model
type Aircraft struct {
	Name                string   `json:"name"`
	Icon                string   `json:"icon"`
	Gallery             []string `json:"gallery"`
	General             string   `json:"general"`             // TODO: AircraftGeneral Object id
	Airframe            string   `json:"airframe"`            // TODO: Airframe Object id
	Engines             []string `json:"engines"`             // TODO: Engine Object id
	APU                 string   `json:"apu"`                 // TODO: APU Object id
	Avionics            string   `json:"avionics"`            // TODO: Avionics Object id
	AdditionalEquipment string   `json:"additionalEquipment"` // TODO: AdditionalEquipment Object id
	Exterior            string   `json:"exterior"`            // TODO: Exterior Object id
	Interior            string   `json:"interior"`            // TODO: Interior Object id
	Performance         string   `json:"performance"`         // TODO: Performance Object id
	Price               int      `json:"price"`
	SellerInfo          string   `json:"sellerInfo"` // TODO: SellerInfo Object id
	Tags                []string `json:"tags"`
}
