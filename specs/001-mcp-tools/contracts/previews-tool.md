# MCP Tool Contract: previews

**Tool Name**: `previews`
**Description**: Fetch race preview data (直前情報) including odds, predictions, and analysis
**Category**: Data Access / Pre-Race Analysis

---

## Signature

```
previews(year: string, date: string) → RacePreview
```

---

## Input Parameters

### `year` (required)
- **Type**: string
- **Format**: YYYY (4-digit year)
- **Valid Range**: 1900-2100
- **Description**: Calendar year for race previews
- **Example**: "2025"
- **Error if**: Not 4 digits, not parseable as integer, outside range [1900, 2100]
- **Error Message**: `"Invalid year 'xyz': expected 4-digit year between 1900 and 2100"`

### `date` (required)
- **Type**: string
- **Format**: YYYYMMDD (8-digit date)
- **Validation**: Must be a valid calendar date
- **Description**: Specific date for which to fetch race previews
- **Example**: "20251222"
- **Error if**: Not 8 digits, not valid calendar date
- **Error Message**: `"Invalid date '20251320': expected valid calendar date in YYYYMMDD format (e.g., 20251222)"`

---

## Processing

1. **Validate Inputs**: Check year and date formats/ranges
   - If invalid, return error with clear message (do not call upstream API)

2. **Check Cache**: Look for cached RacePreview with key `previews:YYYYMMDD`
   - If hit, return cached data immediately
   - If miss, proceed to step 3

3. **Fetch from API**: Call `GET https://boatraceopenapi.github.io/previews/v2/{year}/{date}.json`
   - Timeout: 10 seconds
   - Follow redirects: up to 5
   - Validate SSL certificates

4. **Retry Logic**: If API call fails (connection error, timeout, 5xx), retry up to 2 times
   - Backoff: 100ms initial, 200ms on second retry
   - Do not retry on: 4xx errors, malformed JSON

5. **Parse Response**: Unmarshal JSON into RacePreview struct
   - If JSON malformed, return error with details

6. **Cache Result**: Store response in session cache with key `previews:YYYYMMDD`

7. **Return**: RacePreview object to user

---

## Output

### Success Response

**Type**: `RacePreview` (structured data)

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
      "currentOdds": {
        "win": {
          "1": 2.5,
          "2": 3.2,
          "3": 4.8,
          "4": 6.2,
          "5": 8.5,
          "6": 12.0
        },
        "place": {
          "1": 1.3,
          "2": 1.5,
          "3": 2.1,
          "4": 2.8,
          "5": 3.5,
          "6": 5.2
        },
        "updatedAt": "2025-12-22T10:45:00Z"
      },
      "predictions": [
        {
          "predictor": "Expert_A",
          "predictedWinner": 2,
          "confidence": 0.85,
          "rationale": "Strong form, good motor condition"
        },
        {
          "predictor": "Model_v2",
          "predictedWinner": 1,
          "confidence": 0.72,
          "rationale": "Historical data suggests strong performance"
        }
      ],
      "analysis": {
        "favorites": [1, 2, 3],
        "upsets": [5, 6],
        "keyFactors": ["Wind speed 3.2 m/s", "Water: slightly rough"],
        "weatherImpact": "Wind will favor inside lanes",
        "trackCondition": "Good condition after morning maintenance"
      },
      "latestNews": [
        {
          "timestamp": "2025-12-22T08:30:00Z",
          "title": "Racer 4 reports slight motor issue",
          "description": "Motor adjustment needed; may affect performance",
          "affectedRacers": [4]
        }
      ]
    },
    ... (multiple races)
  ],
  "metadata": {
    "createdAt": "2025-12-22T06:00:00Z",
    "publishTime": "2025-12-22T14:30:00Z",
    "lastUpdated": "2025-12-22T10:45:00Z"
  }
}
```

**Characteristics**:
- Updated odds reflecting latest betting trends
- Multiple prediction sources with confidence levels
- Pre-race analysis with key factors and expert assessment
- Late-breaking news affecting races
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
  "message": "Failed to parse API response: invalid JSON at line 18",
  "code": "INVALID_JSON_RESPONSE"
}
```

---

## Usage Examples

### Example 1: Get Latest Odds and Predictions
```
previews(year="2025", date="20251222")
```
**Result**: Updated odds, expert predictions, and pre-race analysis

### Example 2: Invalid Year Format
```
previews(year="25", date="20251222")
```
**Result**: ValidationError - "Invalid year '25': expected 4-digit year between 1900 and 2100"

### Example 3: Past Date (No Preview Available)
```
previews(year="2025", date="20251215")
```
**Result**: NotFound - "No races found for 20251215" (races already completed)

---

## Caching Behavior

- **Cache Key**: `previews:20251222`
- **Cache Scope**: Session lifetime
- **Expiration**: Session ends / server restarts
- **Hit Frequency**: Medium-High (users may query multiple times for updated odds)
- **Performance**: Cache hit returns data in <100ms
- **Note**: Odds are updated frequently; users may want fresh data periodically

---

## Real-Time Odds Considerations

- **Update Frequency**: Odds update hourly (or more frequently as race approaches)
- **Cache Validity**: Acceptable cache lifetime is 5-10 minutes (typical user session duration)
- **User Notification**: No notification of cache hit; user gets cached data transparently
- **Future Enhancement**: Time-based cache expiration could be added (TBD post-MVP)

---

## Related Tools

- `programs`: Fetch race program with program odds (pre-betting odds)
- `results`: Fetch actual results after races complete

---

## Notes

- Previews available for upcoming/today's races
- Odds subject to change based on betting activity
- Predictions from multiple sources; confidence levels vary
- Late-breaking news critical for informed betting decisions
- All timestamps in UTC
- No authentication required (public API)
- Odds format: decimal (European) notation (e.g., 2.5 = 150% payout)
