package api

import (
	"fmt"

	"github.com/boatrace/open-api-mcp-server/internal/cache"
	"github.com/boatrace/open-api-mcp-server/internal/models"
	"github.com/boatrace/open-api-mcp-server/internal/validation"
)

// FetchPrograms fetches race program data for a specific date
func (c *HTTPClient) FetchPrograms(year, date string) (*models.RaceProgram, error) {
	// Validate inputs
	if err := validation.ValidateDateRange(year, date); err != nil {
		return nil, err
	}

	// Generate cache key
	cacheKey := cache.GenerateKey("programs", date)

	// Fetch with cache and retry logic
	url := fmt.Sprintf("%s/programs/v2/%s/%s.json", baseURL, year, date)
	data, err := c.FetchWithCache(url, cacheKey)
	if err != nil {
		return nil, err
	}

	// Parse response into RaceProgram struct
	var program models.RaceProgram
	if err := UnmarshalJSON(data, &program); err != nil {
		return nil, err
	}

	return &program, nil
}
