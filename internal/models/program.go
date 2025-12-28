package models

// RaceProgram represents a single day's race program information
type RaceProgram struct {
	Date     string        `json:"date"`
	Races    []RaceInfo    `json:"races"`
	Metadata ProgramMetadata `json:"metadata"`
}

// RaceInfo represents information about a single race
type RaceInfo struct {
	Number          int          `json:"number"`
	RaceID          string       `json:"raceId"`
	Venue           string       `json:"venue"`
	Distance        int          `json:"distance"`
	WindSpeed       float64      `json:"windSpeed"`
	WaterCondition  string       `json:"waterCondition"`
	RacerLineup     []RacerInfo  `json:"racerLineup"`
	BettingInfo     BettingInfo  `json:"bettingInfo"`
}

// RacerInfo represents information about a racer in a race
type RacerInfo struct {
	Number        int    `json:"number"`
	Name          string `json:"name"`
	RacerID       string `json:"racerId"`
	MotorNumber   int    `json:"motorNumber"`
	BoatNumber    int    `json:"boatNumber"`
	ClassRating   string `json:"classRating"`
}

// BettingInfo represents betting odds and information for a race
type BettingInfo struct {
	Win        OddsInfo   `json:"win"`
	Place      OddsInfo   `json:"place"`
	ExactaInfo []string   `json:"exactaInfo"`
	TrioInfo   []string   `json:"trioInfo"`
}

// OddsInfo represents odds information for a betting type
type OddsInfo struct {
	Odds      map[int]float64 `json:"odds"`
	Favorites []int           `json:"favorites"`
}

// ProgramMetadata represents metadata about the program
type ProgramMetadata struct {
	CreatedAt   string `json:"createdAt"`
	LastUpdated string `json:"lastUpdated"`
}
