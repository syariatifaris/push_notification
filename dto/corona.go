package dto

type CoronaRegion struct {
	CountryID    string    `json:"country_id" bson:"country_id"`
	UrbanVillage string    `json:"urban_village" bson:"urban_village"`
	City         string    `json:"city" bson:"city"`
	Type         string    `json:"type" bson:"type"`
	Coordinates  []float64 `json:"coordinates" bson:"coordinates"`
}
type Corona struct {
	ID        string       `json:"id" bson:"_id"`
	Zona      string       `json:"zona" bson:"zona"`
	Region    CoronaRegion `json:"region" bson:"region"`
	Confirmed int          `json:"confirmed" bson:"confirmed"`
	Recovered int          `json:"recovered" bson:"recovered"`
	Deaths    int          `json:"deaths" bson:"deaths"`
	Distance  float64      `json:"distance" bson:"distance,omitempty"`
}
