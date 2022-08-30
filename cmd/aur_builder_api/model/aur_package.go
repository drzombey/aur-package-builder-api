package model

type AurPackage struct {
	ObjectId       string   `json:"_objectId" bson:"_id"`
	AurID          string   `json:"ID"`
	Name           string   `json:"Name"`
	PackageBaseID  string   `json:"PackageBaseID"`
	PackageBase    string   `json:"PackageBase"`
	Version        string   `json:"Version"`
	Description    string   `json:"Description"`
	URL            string   `json:"URL"`
	NumVotes       int      `json:"NumVotes"`
	Popularity     int      `json:"Popularity"`
	OutOfDate      string   `json:"OutOfDate"`
	Maintainer     string   `json:"Maintainer"`
	FirstSubmitted int      `json:"FirstSubmitted"`
	LastModified   int      `json:"LastModified"`
	URLPath        string   `json:"URLPath"`
	Depends        []string `json:"Depends"`
	License        []string `json:"License"`
	MakeDepends    []string `json:"MakeDepends"`
}
