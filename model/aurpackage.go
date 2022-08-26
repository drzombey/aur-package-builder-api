package model

type AurPackage struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	CreationDate string `json:"creationDate"`
}
