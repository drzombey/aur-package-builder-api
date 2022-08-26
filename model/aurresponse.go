package model

type AurResponse struct {
	Version     int        `json:"version"`
	Type        string     `json:"type"`
	ResultCount int        `json:"resultcount"`
	Results     AurPackage `json:"results"`
}
