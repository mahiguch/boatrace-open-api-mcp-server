# Quickstart Guide: Boatrace Open API MCP Server

**Purpose**: Get started with the Boatrace Open API MCP Server in 5 minutes
**Created**: 2025-12-22

---

## Installation

### Prerequisites
- Go 1.21 or later
- Claude desktop or compatible MCP client

### Build from Source
```bash
# Clone repository
git clone https://github.com/yourusername/boatrace-open-api-mcp-server.git
cd boatrace-open-api-mcp-server

# Install dependencies
go mod download

# Build executable
make build
# Output: ./bin/boatrace-mcp (or similar)
```

### Run the Server
```bash
# Start MCP server (uses stdio transport)
./bin/boatrace-mcp

# Server is now ready to receive MCP calls from Claude
```

---

## Using the Tools

### Tool 1: Fetch Race Programs (出走表)

**Get the race lineup and program for a specific date:**

```go
// Call signature
programs(year: "2025", date: "20251222")

// Returns: RaceProgram with lineup, odds, race conditions
```

**In Claude:**
```
User: Get the race program for December 22, 2025

Claude calls: programs(year="2025", date="20251222")

Result:
- Race 1: 1800m, 6 racers, wind 2.3 m/s
  - Racer 1: Tanaka (odds 2.5)
  - Racer 2: Suzuki (odds 3.2)
  - ...
- Race 2: [similar structure]
- ... (remaining races)
```

---

### Tool 2: Fetch Race Results (結果)

**Get winners, payouts, and timing for completed races:**

```go
// Call signature
results(year: "2025", date: "20251220")

// Returns: RaceResult with winners, payouts, statistics
```

**In Claude:**
```
User: What were the results from December 20, 2025?

Claude calls: results(year="2025", date="20251220")

Result:
- Race 1: Winner Racer 3 (time: 75.23s)
  - Win payout: ¥2,450 (odds 2.5)
  - Place payout: ¥1,200 (odds 1.3)
  - Exacta payout: ¥8,900 (1-3)
- Race 2: [similar structure]
- ... (remaining races)
```

---

### Tool 3: Fetch Race Previews (直前情報)

**Get latest odds, predictions, and pre-race analysis:**

```go
// Call signature
previews(year: "2025", date: "20251222")

// Returns: RacePreview with current odds, predictions, news
```

**In Claude:**
```
User: What are the latest odds for today's races?

Claude calls: previews(year="2025", date="20251222")

Result:
- Race 1: Current odds updated at 10:45
  - Favorite: Racer 1 (odds 2.5)
  - Expert prediction: Racer 2 wins (confidence 85%)
  - Key factors: Wind 3.2 m/s, good water conditions
  - News: Racer 4 motor adjustment needed
- Race 2: [similar structure]
- ... (remaining races)
```

---

## Common Usage Patterns

### Pattern 1: Compare Program vs Predictions
```
User: Compare the favorite from the program with expert predictions for race 1

Claude calls:
1. programs(year="2025", date="20251222")
   → Gets program favorite from odds
2. previews(year="2025", date="20251222")
   → Gets expert predictions

Claude compares and reports:
- Program favorite: Racer 1 (odds 2.5)
- Expert prediction: Racer 2 (confidence 85%)
- Discrepancy suggests possible upset
```

### Pattern 2: Analyze Past Results
```
User: How did the actual results compare to predictions for Dec 20?

Claude calls:
1. previews(year="2025", date="20251220")
   → Gets previous day's predictions
2. results(year="2025", date="20251220")
   → Gets actual results

Claude analyzes prediction accuracy
```

### Pattern 3: Current Race Information
```
User: Give me a full overview of today's races

Claude calls:
1. programs(year="2025", date="20251222")
   → Gets lineup and race structure
2. previews(year="2025", date="20251222")
   → Gets latest odds and analysis

Claude combines and presents comprehensive overview
```

---

## Error Handling

### Invalid Date Format
```
User: Get programs for December 2025

Claude attempts: programs(year="2025", date="122025")
                                            ↑ Wrong format!

Server responds with error:
{
  "error": "ValidationError",
  "message": "Invalid date '122025': expected valid calendar date in YYYYMMDD format (e.g., 20251222)",
  "code": "INVALID_DATE_FORMAT"
}

Claude clarifies: "I need the date in YYYYMMDD format (e.g., 20251222)"
```

### Network Error / API Unavailable
```
User: Get the latest race data

Claude calls: programs(year="2025", date="20251222")

Server attempts API call → timeout (retries 2 times) → still times out

Server responds with error:
{
  "error": "Timeout",
  "message": "Request to Boatrace API timed out after 10 seconds (retried 2 times)",
  "code": "API_TIMEOUT"
}

Claude reports: "The Boatrace API is currently unavailable. Please try again in a moment."
```

### No Races for Date
```
User: Get programs for 2099-12-31

Claude calls: programs(year="2099", date="20991231")

Server responds with error:
{
  "error": "NotFound",
  "message": "No races found for 20991231",
  "code": "NO_RACES_FOR_DATE"
}

Claude reports: "No races were found for that date."
```

---

## Performance Tips

### Caching is Automatic
- Tool results are cached for the session duration
- Repeated calls to `programs(year="2025", date="20251222")` use cache (fast)
- Cache cleared when MCP session ends

### Optimal Usage
```
✅ GOOD: Fetch once, reference multiple times
  Claude: "Get programs for Dec 22" → cached
  Claude: "Compare with previews" → uses same cache

❌ INEFFICIENT: Refetch same data repeatedly
  Claude: "Get programs" → API call
  Claude: "Get programs again" → API call (no reason to repeat)
```

### Date Boundaries
```
✅ VALID: Past dates (historical data available)
         Today (live odds, program)
         Near-future dates (program available)

❌ UNAVAILABLE: Far future dates (no program yet)
                Results for future dates (not completed)
```

---

## Date Format Reminders

### Year Format: YYYY
```
✅ "2025" (4-digit year)
❌ "25" (2-digit year)
❌ "2025-01" (includes month)
```

### Date Format: YYYYMMDD
```
✅ "20251222" (December 22, 2025)
❌ "12-22-2025" (wrong format)
❌ "20251232" (invalid date - no 32nd day)
❌ "20250230" (invalid date - Feb 30 doesn't exist)
```

---

## Tool Availability & Limitations

### What's Available
| Tool | P1/P2/P3 | Status | Notes |
|------|----------|--------|-------|
| `programs` | P1 | MVP | Core functionality - race lineup & odds |
| `results` | P2 | MVP | Results & payouts - historical analysis |
| `previews` | P3 | MVP | Predictions & updated odds - pre-race insights |

### What's Not Available (Possible Future Features)
- Live race tracking (in-progress race updates)
- Detailed racer statistics/history
- Advanced betting strategy tools
- Custom odds calculations
- Notifications/alerts

---

## Troubleshooting

### "Invalid date format" error
**Cause**: Date not in YYYYMMDD format
**Solution**: Use 8-digit date (e.g., 20251222 for Dec 22, 2025)

### "Request timed out" error
**Cause**: Boatrace API is slow or unavailable
**Solution**: Retry in a few moments; the API may be temporarily unavailable

### "No races found for [date]"
**Cause**: Either no races scheduled or using future date before program released
**Solution**: Try a different date with known races (recent past dates work best)

### Empty or incomplete data
**Cause**: API returned incomplete response
**Solution**: Retry; incomplete responses are temporary glitches

---

## Next Steps

1. **Try the tools**: Start with `programs` for today's date
2. **Combine them**: Use `programs` + `previews` for complete picture
3. **Analyze results**: Use `results` for past dates to understand patterns
4. **Check back**: Visit this guide as you discover new usage patterns

---

## Support

For issues or feature requests:
- GitHub Issues: https://github.com/yourusername/boatrace-open-api-mcp-server/issues
- Documentation: See `specs/001-mcp-tools/` for detailed specifications

---

## Example Session

```
User: I want to understand today's races before placing bets.

Claude: I'll get the current race program and latest predictions for you.
[calls programs(year="2025", date="20251222")]
[calls previews(year="2025", date="20251222")]

Result summary:
- 12 races scheduled for today at [venue]
- Race 1 (1800m): Favorite is Racer 1 (2.5:1 odds)
  Expert prediction: Racer 2 (85% confidence)
  Key factors: Slight wind, good water conditions

- Race 2 (1600m): Favorite is Racer 3 (3.2:1 odds)
  [similar breakdown for each race]

Would you like a detailed analysis of any specific race?
```

---

## Changelog

### v1.0.0 (2025-12-22)
- Initial release
- Three core tools: programs, results, previews
- Session-based caching
- Explicit input validation
- Comprehensive error handling
