# MCP Tool Contract: results

**Tool Name**: `results`
**Description**: Fetch race results (結果) for a specific date
**Category**: Data Access / Historical Data

---

## Signature

```
results(year: string, date: string) → RaceResult
```

---

## Input Parameters

### `year` (required)
- **Type**: string
- **Format**: YYYY (4-digit year)
- **Valid Range**: 1900-2100
- **Description**: Calendar year for the race results
- **Example**: "2025"
- **Error if**: Not 4 digits, not parseable as integer, outside range [1900, 2100]
- **Error Message**: `"Invalid year 'xyz': expected 4-digit year between 1900 and 2100"`

### `date` (required)
- **Type**: string
- **Format**: YYYYMMDD (8-digit date)
- **Validation**: Must be a valid calendar date
- **Description**: Specific date for which to fetch race results
- **Example**: "20251220"
- **Error if**: Not 8 digits, not valid calendar date
- **Error Message**: `"Invalid date '20251320': expected valid calendar date in YYYYMMDD format (e.g., 20251222)"`

---

## Processing

1. **Validate Inputs**: Check year and date formats/ranges
   - If invalid, return error with clear message (do not call upstream API)

2. **Check Cache**: Look for cached RaceResult with key `results:YYYYMMDD`
   - If hit, return cached data immediately
   - If miss, proceed to step 3

3. **Fetch from API**: Call `GET https://boatraceopenapi.github.io/results/v2/{year}/{date}.json`
   - Timeout: 10 seconds
   - Follow redirects: up to 5
   - Validate SSL certificates

4. **Retry Logic**: If API call fails (connection error, timeout, 5xx), retry up to 2 times
   - Backoff: 100ms initial, 200ms on second retry
   - Do not retry on: 4xx errors, malformed JSON

5. **Parse Response**: Unmarshal JSON into RaceResult struct
   - If JSON malformed, return error with details

6. **Cache Result**: Store response in session cache with key `results:YYYYMMDD`

7. **Return**: RaceResult object to user

---

## Output

### Success Response

**Type**: `RaceResult` (structured data)

**Structure**:
```
{
  "date": "20251220",
  "races": [
    {
      "number": 1,
      "raceID": "...",
      "venue": "...",
      "winnerRacerInfo": {
        "number": 3,
        "name": "Winner Name",
        "motorNumber": 42,
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
        },
        "exacta": {
          "amount": 8900.0,
          "probability": 0.12,
          "winningNumbers": [3, 1]
        }
      },
      "timeResult": {
        "startTime": "2025-12-20T14:30:00Z",
        "finishTime": "2025-12-20T14:31:15Z",
        "raceTime": 75.23
      },
      "statistics": {
        "averageSpeed": 47.3,
        "lapRecords": [
          { "lapNumber": 1, "time": 37.8, "racer": 3 },
          { "lapNumber": 2, "time": 37.43, "racer": 3 }
        ]
      }
    },
    ... (multiple races)
  ],
  "metadata": {
    "createdAt": "2025-12-20T16:00:00Z",
    "disqualifiedRacers": []
  }
}
```

**Characteristics**:
- Complete race results for the specified date
- All completed races included with winners and payouts
- Timing data and statistics for each race
- No missing fields

---

## Error Responses

### Validation Error
**Condition**: Invalid input parameters
**Response**:
```
{
  "error": "ValidationError",
  "message": "Invalid date '20251320': expected valid calendar date in YYYYMMDD format",
  "code": "INVALID_DATE_FORMAT"
}
```

### API Error - Not Found
**Condition**: Date has no races or results API returns 404
**Response**:
```
{
  "error": "NotFound",
  "message": "No races found for 20251220",
  "code": "NO_RACES_FOR_DATE"
}
```

### API Error - Timeout
**Condition**: Network timeout (>10 seconds) after 2 retries
**Response**:
```
{
  "error": "Timeout",
  "message": "Request to Boatrace API timed out after 10 seconds (retried 2 times)",
  "code": "API_TIMEOUT"
}
```

### API Error - Server Error
**Condition**: Boatrace API returns 5xx error
**Response**:
```
{
  "error": "APIError",
  "message": "Boatrace API returned 503 Service Unavailable (retried 2 times)",
  "code": "API_SERVER_ERROR"
}
```

### JSON Parsing Error
**Condition**: API response is not valid JSON
**Response**:
```
{
  "error": "ParseError",
  "message": "Failed to parse API response: invalid JSON at line 12",
  "code": "INVALID_JSON_RESPONSE"
}
```

---

## Usage Examples

### Example 1: Fetch Results for Completed Race Day
```
results(year="2025", date="20251220")
```
**Result**: Race results with winners, payouts, and timing data

### Example 2: Invalid Date Format
```
results(year="2025", date="12/20")
```
**Result**: ValidationError - "Invalid date '12/20': expected valid calendar date in YYYYMMDD format"

### Example 3: Future Date (Results Not Available)
```
results(year="2025", date="20251230")
```
**Result**: NotFound - "No races found for 20251230" (races not yet completed)

---

## Caching Behavior

- **Cache Key**: `results:20251220`
- **Cache Scope**: Session lifetime
- **Expiration**: Session ends / server restarts
- **Hit Frequency**: High (users often query same date multiple times)
- **Performance**: Cache hit returns data in <100ms

---

## Related Tools

- `programs`: Fetch race program for same date (race lineup information)
- `previews`: Fetch race preview/analysis for same date

---

## Notes

- Results are typically available 1-2 hours after race completion
- Historical data available for past dates (no limit)
- All timestamps in UTC
- No authentication required (public API)
- Payouts are in Japanese Yen (¥)
