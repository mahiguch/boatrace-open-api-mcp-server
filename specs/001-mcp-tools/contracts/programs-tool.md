# MCP Tool Contract: programs

**Tool Name**: `programs`
**Description**: Fetch race program data (出走表) for a specific date
**Category**: Data Access / Race Information

---

## Signature

```
programs(year: string, date: string) → RaceProgram
```

---

## Input Parameters

### `year` (required)
- **Type**: string
- **Format**: YYYY (4-digit year)
- **Valid Range**: 1900-2100
- **Description**: Calendar year for the race program
- **Example**: "2025"
- **Error if**: Not 4 digits, not parseable as integer, outside range [1900, 2100]
- **Error Message**: `"Invalid year 'xyz': expected 4-digit year between 1900 and 2100"`

### `date` (required)
- **Type**: string
- **Format**: YYYYMMDD (8-digit date)
- **Validation**: Must be a valid calendar date
- **Description**: Specific date for which to fetch the race program
- **Example**: "20251222"
- **Error if**: Not 8 digits, not valid calendar date (e.g., 20251320, 20250230)
- **Error Message**: `"Invalid date '20251320': expected valid calendar date in YYYYMMDD format (e.g., 20251222)"`

---

## Processing

1. **Validate Inputs**: Check year and date formats/ranges
   - If invalid, return error with clear message (do not call upstream API)

2. **Check Cache**: Look for cached RaceProgram with key `programs:YYYYMMDD`
   - If hit, return cached data immediately
   - If miss, proceed to step 3

3. **Fetch from API**: Call `GET https://boatraceopenapi.github.io/programs/v2/{year}/{date}.json`
   - Timeout: 10 seconds
   - Follow redirects: up to 5
   - Validate SSL certificates

4. **Retry Logic**: If API call fails (connection error, timeout, 5xx), retry up to 2 times
   - Backoff: 100ms initial, 200ms on second retry
   - Do not retry on: 4xx errors, malformed JSON

5. **Parse Response**: Unmarshal JSON into RaceProgram struct
   - If JSON malformed, return error with details

6. **Cache Result**: Store response in session cache with key `programs:YYYYMMDD`

7. **Return**: RaceProgram object to user

---

## Output

### Success Response

**Type**: `RaceProgram` (structured data)

**Structure**:
```
{
  "date": "20251222",
  "races": [
    {
      "number": 1,
      "raceID": "...",
      "venue": "...",
      "distance": 1800,
      "windSpeed": 2.3,
      "waterCondition": "slightly rough",
      "racerLineup": [
        {
          "number": 1,
          "name": "Racer Name",
          "racerID": "...",
          "motorNumber": 42,
          "boatNumber": 15,
          "classRating": "A1"
        },
        ... (6 racers total)
      ],
      "bettingInfo": {
        "win": { "odds": { "1": 2.5, "2": 3.2, ... }, "favorites": [1, 2] },
        "place": { ... }
      }
    },
    ... (multiple races)
  ],
  "metadata": {
    "createdAt": "2025-12-22T00:00:00Z",
    "lastUpdated": "2025-12-22T06:00:00Z"
  }
}
```

**Characteristics**:
- Complete race program for the specified date
- All races for the day included
- All racer information and odds present
- No missing fields (complete data set)

---

## Error Responses

### Validation Error
**Condition**: Invalid input parameters
**Response**: Error message with specific issue
```
{
  "error": "ValidationError",
  "message": "Invalid date '20251222': expected valid calendar date in YYYYMMDD format",
  "code": "INVALID_DATE_FORMAT"
}
```

### API Error - Not Found
**Condition**: Date has no races or API returns 404
**Response**:
```
{
  "error": "NotFound",
  "message": "No races found for 20251222",
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
  "message": "Failed to parse API response: invalid JSON at line 5",
  "code": "INVALID_JSON_RESPONSE"
}
```

---

## Usage Examples

### Example 1: Fetch Programs for Today
```
programs(year="2025", date="20251222")
```
**Result**: Full race program for December 22, 2025

### Example 2: Invalid Date Format
```
programs(year="2025", date="12-22")
```
**Result**: ValidationError - "Invalid date '12-22': expected valid calendar date in YYYYMMDD format"

### Example 3: Future Date (No Races Yet)
```
programs(year="2099", date="20991231")
```
**Result**: NotFound - "No races found for 20991231" (or empty races array if API supports future queries)

---

## Caching Behavior

- **Cache Key**: `programs:20251222`
- **Cache Scope**: Session lifetime
- **Expiration**: Session ends / server restarts
- **Hit Frequency**: High (users often query same date multiple times in single session)
- **Performance**: Cache hit returns data in <100ms

---

## Related Tools

- `results`: Fetch completed race outcomes for same date
- `previews`: Fetch race preview/analysis for same date

---

## Notes

- Tool is stateless; each call is independent
- Caching is transparent; user sees same fast response on repeated calls
- All timestamps in UTC
- No authentication required (public API)
