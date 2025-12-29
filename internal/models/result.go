package models

// RaceResult represents completed race results
type RaceResult struct {
	Date     string         `json:"date"`
	Races    []RaceOutcome  `json:"races"`
	Metadata ResultMetadata `json:"metadata"`
}

// RaceOutcome represents the outcome of a completed race
type RaceOutcome struct {
	Number          int             `json:"number"`
	RaceID          string          `json:"raceId"`
	Venue           string          `json:"venue"`
	WinnerRacerInfo RacerResult     `json:"winnerRacerInfo"`
	PayoutInfo      PayoutInfo      `json:"payoutInfo"`
	TimeResult      TimeResult      `json:"timeResult"`
	Statistics      RaceStatistics  `json:"statistics"`
}

// RacerResult represents information about a racer result
type RacerResult struct {
	Number      int     `json:"number"`
	Name        string  `json:"name"`
	MotorNumber int     `json:"motorNumber"`
	Time        float64 `json:"time"`
	Reaction    float64 `json:"reaction"`
}

// PayoutInfo represents payout information for a race
type PayoutInfo struct {
	Win    PayoutRecord `json:"win"`
	Place  PayoutRecord `json:"place"`
	Show   PayoutRecord `json:"show"`
	Exacta PayoutRecord `json:"exacta"`
	Trio   PayoutRecord `json:"trio"`
}

// PayoutRecord represents a single payout record
type PayoutRecord struct {
	Amount         float64 `json:"amount"`
	Probability    float64 `json:"probability"`
	WinningNumbers []int   `json:"winningNumbers"`
}

// TimeResult represents timing information for a race
type TimeResult struct {
	StartTime  string  `json:"startTime"`
	FinishTime string  `json:"finishTime"`
	RaceTime   float64 `json:"raceTime"`
}

// RaceStatistics represents statistics about a race
type RaceStatistics struct {
	AverageSpeed float64   `json:"averageSpeed"`
	LapRecords   []LapTime `json:"lapRecords"`
}

// LapTime represents timing for a single lap
type LapTime struct {
	LapNumber int     `json:"lapNumber"`
	Time      float64 `json:"time"`
	Racer     int     `json:"racer"`
}

// ResultMetadata represents metadata about results
type ResultMetadata struct {
	CreatedAt           string `json:"createdAt"`
	DisqualifiedRacers  []int  `json:"disqualifiedRacers"`
}
