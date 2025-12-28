package models

// RacePreview represents pre-race information and analysis
type RacePreview struct {
	Date     string        `json:"date"`
	Races    []PreviewInfo `json:"races"`
	Metadata PreviewMetadata `json:"metadata"`
}

// PreviewInfo represents preview information for a single race
type PreviewInfo struct {
	Number          int              `json:"number"`
	RaceID          string           `json:"raceId"`
	Venue           string           `json:"venue"`
	Distance        int              `json:"distance"`
	CurrentOdds     CurrentOddsInfo  `json:"currentOdds"`
	Predictions     []PredictionRecord `json:"predictions"`
	Analysis        RaceAnalysis     `json:"analysis"`
	LatestNews      []NewsItem       `json:"latestNews"`
}

// CurrentOddsInfo represents current betting odds
type CurrentOddsInfo struct {
	Win        map[int]float64   `json:"win"`
	Place      map[int]float64   `json:"place"`
	ExactaOdds map[string]float64 `json:"exactaOdds"`
	UpdatedAt  string            `json:"updatedAt"`
}

// PredictionRecord represents a prediction for a race
type PredictionRecord struct {
	Predictor       string  `json:"predictor"`
	PredictedWinner int     `json:"predictedWinner"`
	Confidence      float64 `json:"confidence"`
	Rationale       string  `json:"rationale"`
}

// RaceAnalysis represents analysis of a race
type RaceAnalysis struct {
	Favorites      []int    `json:"favorites"`
	Upsets         []int    `json:"upsets"`
	KeyFactors     []string `json:"keyFactors"`
	WeatherImpact  string   `json:"weatherImpact"`
	TrackCondition string   `json:"trackCondition"`
}

// NewsItem represents a news item related to a race
type NewsItem struct {
	Timestamp      string `json:"timestamp"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AffectedRacers []int  `json:"affectedRacers"`
}

// PreviewMetadata represents metadata about previews
type PreviewMetadata struct {
	CreatedAt   string `json:"createdAt"`
	PublishTime string `json:"publishTime"`
	LastUpdated string `json:"lastUpdated"`
}
