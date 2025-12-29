package api

import (
	"fmt"

	"github.com/boatrace/open-api-mcp-server/internal/cache"
	"github.com/boatrace/open-api-mcp-server/internal/models"
	"github.com/boatrace/open-api-mcp-server/internal/validation"
)

// FetchPreviews fetches race preview data for a specific date
func (c *HTTPClient) FetchPreviews(year, date string) (*models.RacePreview, error) {
	// Validate inputs
	if err := validation.ValidateDateRange(year, date); err != nil {
		return nil, err
	}

	// Generate cache key
	cacheKey := cache.GenerateKey("previews", date)

	// Fetch with cache and retry logic
	url := fmt.Sprintf("%s/previews/v2/%s/%s.json", baseURL, year, date)
	data, err := c.FetchWithCache(url, cacheKey)
	if err != nil {
		return nil, err
	}

	// Parse response into RacePreview struct
	var preview models.RacePreview
	if err := UnmarshalJSON(data, &preview); err != nil {
		return nil, err
	}

	return &preview, nil
}
