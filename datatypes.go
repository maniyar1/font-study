package main

// JSON Structs
type JSONFormat struct {
	Items []Font
}

type Font struct {
	Kind         string
	Family       string
	Category     string
	Variants     []string
	Subsets      []string
	Version      string
	LastModified string
	Files        Files
}

type Files struct {
	Regular string
}

// Option is for HTML/CSS
type Option struct {
	Number  int
	Pangram string
	Font    Font
}

type PageData struct {
	CSS     string
	Options []Option
}

// For storing in leveldb
type FontRatings struct {
	Family                string
	Points                uint64
	AveragePoints         float64
	TotalEntries          uint64
	FirstPlaceOccurances  uint
	SecondPlaceOccurances uint
	ThirdPlaceOccurances  uint
	FourthPlaceOccurances uint
	FifthPlaceOccurances  uint
	SixthPlaceOccurances  uint
}

//for full JSON
type FontRatingsArray []FontRatings

// For thanks page
type Responses struct {
	Responses []Response
}
type Response struct {
	Family        string
	UserPoints    uint
	AveragePoints float64
}
