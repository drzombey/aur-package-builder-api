package aur

type Package struct {
	ID             int64       `json:"ID"`
	Name           string      `json:"Name"`
	PackageBaseID  int64       `json:"PackageBaseID"`
	PackageBase    string      `json:"PackageBase"`
	Version        string      `json:"Version"`
	Description    string      `json:"Description"`
	URL            string      `json:"URL"`
	NumVotes       int64       `json:"NumVotes"`
	Popularity     int64       `json:"Popularity"`
	OutOfDate      interface{} `json:"OutOfDate"`
	Maintainer     string      `json:"Maintainer"`
	FirstSubmitted int64       `json:"FirstSubmitted"`
	LastModified   int64       `json:"LastModified"`
	URLPath        string      `json:"URLpath"`
}

type ResponseInfo struct {
	Version     int       `json:"version"`
	Type        string    `json:"type"`
	ResultCount int       `json:"resultcount"`
	Packages    []Package `json:"results"`
}

type SearchType int

const (
	Search SearchType = iota
	Multiinfo
	Error
)

func (s SearchType) String() string {
	switch s {
	case Search:
		return "search"
	case Multiinfo:
		return "multiinfo"
	case Error:
		return "error"
	}

	return "unknown"
}
