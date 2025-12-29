package api

import (
	"fmt"

	"github.com/boatrace/open-api-mcp-server/internal/cache"
	"github.com/boatrace/open-api-mcp-server/internal/models"
	"github.com/boatrace/open-api-mcp-server/internal/validation"
)

// FetchResults fetches race results data for a specific date
func (c *HTTPClient) FetchResults(year, date string) (*models.RaceResult, error) {
	// Validate inputs
	if err := validation.ValidateDateRange(year, date); err != nil {
		return nil, err
	}

	// Generate cache key
	cacheKey := cache.GenerateKey("results", date)

	// Fetch with cache and retry logic
	url := fmt.Sprintf("%s/results/v2/%s/%s.json", baseURL, year, date)
	data, err := c.FetchWithCache(url, cacheKey)
	if err != nil {
		return nil, err
	}

	// Parse response into RaceResult struct
	var result models.RaceResult
	if err := UnmarshalJSON(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
