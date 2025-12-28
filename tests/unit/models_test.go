package unit

import (
	"encoding/json"
	"testing"

	"github.com/boatrace/open-api-mcp-server/internal/models"
)

// TestRaceProgramUnmarshal tests unmarshalling race program JSON
func TestRaceProgramUnmarshal(t *testing.T) {
	jsonData := []byte(`{
		"date": "20251222",
		"races": [
			{
				"number": 1,
				"raceId": "race-1",
				"venue": "venue-1",
				"distance": 1800,
				"windSpeed": 2.3,
				"waterCondition": "calm",
				"racerLineup": [
					{
						"number": 1,
						"name": "Racer 1",
						"racerId": "racer-1",
						"motorNumber": 10,
						"boatNumber": 20,
						"classRating": "A1"
					}
				],
				"bettingInfo": {
					"win": {
						"odds": {"1": 2.5, "2": 3.2},
						"favorites": [1, 2]
					},
					"place": {
						"odds": {"1": 1.3, "2": 1.5},
						"favorites": [1]
					}
				}
			}
		],
		"metadata": {
			"createdAt": "2025-12-22T00:00:00Z",
			"lastUpdated": "2025-12-22T06:00:00Z"
		}
	}`)

	var program models.RaceProgram
	err := json.Unmarshal(jsonData, &program)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if program.Date != "20251222" {
		t.Errorf("expected date 20251222, got %s", program.Date)
	}

	if len(program.Races) != 1 {
		t.Errorf("expected 1 race, got %d", len(program.Races))
	}

	if program.Races[0].Number != 1 {
		t.Errorf("expected race number 1, got %d", program.Races[0].Number)
	}

	if len(program.Races[0].RacerLineup) != 1 {
		t.Errorf("expected 1 racer, got %d", len(program.Races[0].RacerLineup))
	}
}

// TestRaceResultUnmarshal tests unmarshalling race result JSON
func TestRaceResultUnmarshal(t *testing.T) {
	jsonData := []byte(`{
		"date": "20251220",
		"races": [
			{
				"number": 1,
				"raceId": "race-1",
				"venue": "venue-1",
				"winnerRacerInfo": {
					"number": 3,
					"name": "Winner",
					"motorNumber": 10,
					"time": 75.23,
					"reaction": 0.12
				},
				"payoutInfo": {
					"win": {
						"amount": 2450.0,
						"probability": 0.35,
						"winningNumbers": [3]
					},
					"place": {
						"amount": 1200.0,
						"probability": 0.60,
						"winningNumbers": [3]
					}
				},
				"timeResult": {
					"startTime": "2025-12-20T14:30:00Z",
					"finishTime": "2025-12-20T14:31:15Z",
					"raceTime": 75.23
				},
				"statistics": {
					"averageSpeed": 47.3,
					"lapRecords": []
				}
			}
		],
		"metadata": {
			"createdAt": "2025-12-20T16:00:00Z",
			"disqualifiedRacers": []
		}
	}`)

	var result models.RaceResult
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Date != "20251220" {
		t.Errorf("expected date 20251220, got %s", result.Date)
	}

	if len(result.Races) != 1 {
		t.Errorf("expected 1 race, got %d", len(result.Races))
	}

	if result.Races[0].WinnerRacerInfo.Number != 3 {
		t.Errorf("expected winner number 3, got %d", result.Races[0].WinnerRacerInfo.Number)
	}

	if result.Races[0].PayoutInfo.Win.Amount != 2450.0 {
		t.Errorf("expected win payout 2450.0, got %f", result.Races[0].PayoutInfo.Win.Amount)
	}
}

// TestRacePreviewUnmarshal tests unmarshalling race preview JSON
func TestRacePreviewUnmarshal(t *testing.T) {
	jsonData := []byte(`{
		"date": "20251222",
		"races": [
			{
				"number": 1,
				"raceId": "race-1",
				"venue": "venue-1",
				"distance": 1800,
				"currentOdds": {
					"win": {"1": 2.5, "2": 3.2},
					"place": {"1": 1.3, "2": 1.5},
					"exactaOdds": {"1-2": 8.5},
					"updatedAt": "2025-12-22T10:45:00Z"
				},
				"predictions": [
					{
						"predictor": "expert",
						"predictedWinner": 2,
						"confidence": 0.85,
						"rationale": "Strong form"
					}
				],
				"analysis": {
					"favorites": [1, 2],
					"upsets": [5, 6],
					"keyFactors": ["wind 3.2 m/s"],
					"weatherImpact": "favorable",
					"trackCondition": "good"
				},
				"latestNews": []
			}
		],
		"metadata": {
			"createdAt": "2025-12-22T06:00:00Z",
			"publishTime": "2025-12-22T14:30:00Z",
			"lastUpdated": "2025-12-22T10:45:00Z"
		}
	}`)

	var preview models.RacePreview
	err := json.Unmarshal(jsonData, &preview)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if preview.Date != "20251222" {
		t.Errorf("expected date 20251222, got %s", preview.Date)
	}

	if len(preview.Races) != 1 {
		t.Errorf("expected 1 race, got %d", len(preview.Races))
	}

	if len(preview.Races[0].Predictions) != 1 {
		t.Errorf("expected 1 prediction, got %d", len(preview.Races[0].Predictions))
	}

	if preview.Races[0].Predictions[0].Confidence != 0.85 {
		t.Errorf("expected confidence 0.85, got %f", preview.Races[0].Predictions[0].Confidence)
	}
}

// TestRaceProgramMarshal tests marshalling race program to JSON
func TestRaceProgramMarshal(t *testing.T) {
	program := models.RaceProgram{
		Date: "20251222",
		Races: []models.RaceInfo{
			{
				Number:  1,
				RaceID:  "race-1",
				Venue:   "venue-1",
				Distance: 1800,
			},
		},
		Metadata: models.ProgramMetadata{
			CreatedAt:   "2025-12-22T00:00:00Z",
			LastUpdated: "2025-12-22T06:00:00Z",
		},
	}

	data, err := json.Marshal(program)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("Marshal returned empty data")
	}

	// Verify it can be unmarshalled back
	var unmarshalled models.RaceProgram
	err = json.Unmarshal(data, &unmarshalled)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshalled.Date != program.Date {
		t.Errorf("roundtrip failed: expected %s, got %s", program.Date, unmarshalled.Date)
	}
}
