# Data Model: Boatrace Open API MCP Server

**Purpose**: Define data entities and relationships for core tools implementation
**Created**: 2025-12-22

---

## Entities

### 1. RaceProgram

**Description**: Represents a single day's race program information from the Boatrace Open API.

**Source**: `GET https://boatraceopenapi.github.io/programs/v2/{YYYY}/{YYYYMMDD}.json`

**Go Struct Definition**:
```go
type RaceProgram struct {
    Date      string            // YYYYMMDD format
    Races     []RaceInfo        // Array of races for the day
    Metadata  ProgramMetadata   // Program-level information
}

type RaceInfo struct {
    Number        int              // Race number (1-12 typically)
    RaceID        string           // Unique race identifier
    Venue         string           // Race venue name
    Distance      int              // Race distance (meters)
    WindSpeed     float64          // Wind speed condition
    WaterCondition string          // Water condition description
    RacerLineup   []RacerInfo      // Array of racers in this race
    BettingInfo   BettingInfo      // Betting odds and information
}

type RacerInfo struct {
    Number        int      // Racer number (1-6 typically)
    Name          string   // Racer name
    RacerID       string   // Unique racer identifier
    MotorNumber   int      // Motor number
    BoatNumber    int      // Boat number
    ClassRating   string   // Racer class/rating
}

type BettingInfo struct {
    Win            OddsInfo    // Win (単勝) odds
    Place          OddsInfo    // Place (複勝) odds
    ExactaInfo     []string    // Exacta (連勝) combinations
    TrioInfo       []string    // Trio (三連勝) combinations
}

type OddsInfo struct {
    Odds           map[int]float64  // Racer number → odds
    Favorites      []int            // Top favorite racer numbers
}

type ProgramMetadata struct {
    CreatedAt      string    // ISO 8601 timestamp
    LastUpdated    string    // ISO 8601 timestamp
}
```

**Relationships**:
- `RaceProgram` contains multiple `RaceInfo`
- `RaceInfo` contains multiple `RacerInfo`
- `RaceInfo` contains `BettingInfo` with odds for each racer

**Validation Rules**:
- `Date`: Must be valid YYYYMMDD format
- `RaceInfo.Number`: Integer between 1-12
- `RacerInfo.Number`: Integer between 1-6
- `BettingInfo.Odds`: All values must be non-negative floats

**State Transitions**: None (immutable response data)

---

### 2. RaceResult

**Description**: Represents a completed race's results and outcomes from the Boatrace Open API.

**Source**: `GET https://boatraceopenapi.github.io/results/v2/{YYYY}/{YYYYMMDD}.json`

**Go Struct Definition**:
```go
type RaceResult struct {
    Date      string            // YYYYMMDD format
    Races     []RaceOutcome     // Array of race results
    Metadata  ResultMetadata    // Result-level information
}

type RaceOutcome struct {
    Number           int              // Race number
    RaceID           string           // Unique race identifier
    Venue            string           // Race venue name
    WinnerRacerInfo  RacerResult      // Winning racer details
    PayoutInfo       PayoutInfo       // Win/place/show payouts
    TimeResult       TimeResult       // Race time and timing data
    Statistics       RaceStatistics   // Performance statistics
}

type RacerResult struct {
    Number        int      // Winning racer number
    Name          string   // Racer name
    MotorNumber   int      // Motor number used
    Time          float64  // Race time in seconds (e.g., 75.23)
    Reaction      float64  // Reaction time in seconds
}

type PayoutInfo struct {
    Win            PayoutRecord   // Win (単勝) payout
    Place          PayoutRecord   // Place (複勝) payout
    Show           PayoutRecord   // Show (3連複) payout
    Exacta         PayoutRecord   // Exacta (連勝) payout
    Trio           PayoutRecord   // Trio (三連勝) payout
}

type PayoutRecord struct {
    Amount         float64        // Payout amount (yen)
    Probability    float64        // Winning probability (0.0-1.0)
    WinningNumbers []int          // Racer numbers in winning combination
}

type TimeResult struct {
    StartTime      string    // ISO 8601 timestamp
    FinishTime     string    // ISO 8601 timestamp
    RaceTime       float64   // Total race time in seconds
}

type RaceStatistics struct {
    AverageSpeed   float64   // Average speed (km/h)
    LapRecords     []LapTime // Individual lap times
}

type LapTime struct {
    LapNumber      int       // Lap number
    Time           float64   // Time for this lap
    Racer          int       // Racer number
}

type ResultMetadata struct {
    CreatedAt      string    // ISO 8601 timestamp
    DisqualifiedRacers []int // Racer numbers disqualified
}
```

**Relationships**:
- `RaceResult` contains multiple `RaceOutcome`
- `RaceOutcome` contains `RacerResult` (winner), `PayoutInfo`, `TimeResult`, `RaceStatistics`
- `PayoutInfo` contains multiple `PayoutRecord` types
- `RaceStatistics` contains multiple `LapTime`

**Validation Rules**:
- `Date`: Must be valid YYYYMMDD format
- `RaceOutcome.Number`: Integer between 1-12
- `RacerResult.Number`: Integer between 1-6
- `PayoutInfo.Amount`: Must be non-negative float
- `TimeResult.RaceTime`: Must be positive float

**State Transitions**: None (immutable historical data)

---

### 3. RacePreview

**Description**: Represents pre-race analysis, predictions, and updated odds from the Boatrace Open API.

**Source**: `GET https://boatraceopenapi.github.io/previews/v2/{YYYY}/{YYYYMMDD}.json`

**Go Struct Definition**:
```go
type RacePreview struct {
    Date      string              // YYYYMMDD format
    Races     []PreviewInfo       // Array of race previews
    Metadata  PreviewMetadata     // Preview-level information
}

type PreviewInfo struct {
    Number           int                   // Race number
    RaceID           string                // Unique race identifier
    Venue            string                // Race venue name
    Distance         int                   // Race distance
    CurrentOdds      CurrentOddsInfo       // Updated betting odds
    Predictions      []PredictionRecord    // Expert predictions
    Analysis         RaceAnalysis          // Pre-race analysis
    LatestNews       []NewsItem            // Late-breaking information
}

type CurrentOddsInfo struct {
    Win            map[int]float64  // Racer number → current odds
    Place          map[int]float64  // Racer number → place odds
    ExactaOdds     map[string]float64    // Exacta combinations → odds
    UpdatedAt      string           // ISO 8601 timestamp of last update
}

type PredictionRecord struct {
    Predictor      string           // Prediction source (expert, model, etc.)
    PredictedWinner int             // Predicted winning racer number
    Confidence     float64          // Confidence level (0.0-1.0)
    Rationale      string           // Prediction explanation
}

type RaceAnalysis struct {
    Favorites      []int            // Top 3 favorite racer numbers
    Upsets         []int            // Underdog racers with potential
    KeyFactors     []string         // Important race conditions
    WeatherImpact  string           // Weather influence on race
    TrackCondition string           // Track/water condition assessment
}

type NewsItem struct {
    Timestamp      string           // ISO 8601 timestamp
    Title          string           // News headline
    Description    string           // News details
    AffectedRacers []int            // Racer numbers affected by news
}

type PreviewMetadata struct {
    CreatedAt      string    // ISO 8601 timestamp
    PublishTime    string    // Scheduled race time
    LastUpdated    string    // Latest odds update time
}
```

**Relationships**:
- `RacePreview` contains multiple `PreviewInfo`
- `PreviewInfo` contains `CurrentOddsInfo`, array of `PredictionRecord`, `RaceAnalysis`, array of `NewsItem`
- Multiple `PredictionRecord` entries (various sources)
- Multiple `NewsItem` entries (time-ordered)

**Validation Rules**:
- `Date`: Must be valid YYYYMMDD format
- `PreviewInfo.Number`: Integer between 1-12
- `CurrentOddsInfo.Win/Place`: All racer number keys must be 1-6
- `PredictionRecord.Confidence`: Must be float between 0.0-1.0
- `PreviewInfo.Favorites`: Racer numbers must be unique, between 1-6

**State Transitions**: None (immutable snapshot data)

---

## Cross-Entity Relationships

```
RaceProgram
├── RaceInfo (for each scheduled race)
│   ├── RacerInfo × 6 (racer lineup)
│   └── BettingInfo (pre-race odds)

RaceResult
├── RaceOutcome (for each completed race)
│   ├── RacerResult (winning racer)
│   ├── PayoutInfo (race payouts)
│   ├── TimeResult (timing data)
│   └── RaceStatistics (performance data)

RacePreview
├── PreviewInfo (for each upcoming race)
│   ├── CurrentOddsInfo (live odds)
│   ├── PredictionRecord × N (expert predictions)
│   ├── RaceAnalysis (race assessment)
│   └── NewsItem × N (latest updates)
```

---

## Shared Concepts

### Race Number
- Integer: 1-12 (typical for Boatrace venues)
- Uniquely identifies a race within a given day/venue

### Racer Number
- Integer: 1-6 (standard boat race lineup)
- Uniquely identifies a racer within a given race

### Date Format
- YYYYMMDD string (e.g., "20251222")
- Must represent valid calendar date
- Validation: Go's `time.Parse("20060102", yyyymmdd)`

### Odds Format
- Float64 in JSON (e.g., 3.5, 12.8)
- Represents betting odds ratio

### Timestamp Format
- ISO 8601 string (e.g., "2025-12-22T14:30:00Z")
- All timestamps in UTC

---

## Validation Rules Summary

| Entity | Field | Rule | Go Type |
|--------|-------|------|---------|
| All | Date | Valid YYYYMMDD | string |
| Race* | Number | 1-12 | int |
| Racer* | Number | 1-6 | int |
| Odds* | Amount | ≥0 | float64 |
| Prediction | Confidence | 0.0-1.0 | float64 |

---

## Implementation Notes

- All entities use exported fields (PascalCase) for JSON marshalling
- JSON tags to be added during implementation (e.g., `json:"raceNumber"`)
- Nil slices are allowed for optional arrays (e.g., empty races, no news)
- Thread-safety not required (models are immutable after unmarshalling)
- No database persistence (models live in memory for session duration)

---

## Next Steps: Phase 1 - API Contracts

Ready to define MCP tool contracts based on these entities:
1. `programs` tool: Input (YYYY, YYYYMMDD) → Output (RaceProgram)
2. `results` tool: Input (YYYY, YYYYMMDD) → Output (RaceResult)
3. `previews` tool: Input (YYYY, YYYYMMDD) → Output (RacePreview)
