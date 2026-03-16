# Music Streaming Mock Data Generator - Improved

This is the improved mock data generator for the Music Streaming platform, aligned with the API documentation and best practices.

## What's New (Improvements)

### ✨ Key Features

1. **JSON Export** - Generate mock data as structured JSON files
   - Aligned with `MOCK_DATA_SCHEMA.json`
   - Metadata tracking and versioning
   - Reusable across different tools

2. **Schema Validation** - Validate mock data against JSON schema
   - Ensures data integrity
   - Prevents invalid datasets
   - Detailed error reporting

3. **Improved Code Organization**
   - Better separation of concerns
   - Enhanced error handling
   - Type-safe data structures
   - Well-documented functions

4. **Multiple Export Formats**
   - JSON for API testing and documentation
   - SQL for database seeding
   - Automatic timestamps and organization

5. **Dataset Presets**
   - **Minimal**: 10 users, 3 artists, 10 tracks (quick testing)
   - **Standard**: 50 users, 15 artists, 100 tracks (feature development)
   - **Stress**: 500 users, 100 artists, 1000 tracks (performance testing)
   - **E2E**: 30 users, 10 artists, 75 tracks (integration testing)

6. **Enhanced Menu System**
   - More options and better organization
   - Input validation
   - Clear feedback and status messages
   - New: Export JSON, Validate, View dataset info

## Installation

```bash
cd music-streaming/mock-data

# Install dependencies (if needed)
go mod download

# Build the generator
go build -o mock-data-generator

# Run the generator
./mock-data-generator
```

## Usage

### Interactive Menu

```bash
./mock-data-generator
```

**Available Options:**
1. **Mock users** - Generate user accounts
   - Choose API calls or SQL queries
   - Input desired quantity

2. **Mock profiles** - Generate user profiles
   - Includes names, avatars, bios
   - Links to users

3. **Mock artists** - Generate artist profiles
   - Creates artist entities
   - Links to user accounts
   - Verification status

4. **Mock tracks** - Generate audio tracks
   - Track metadata
   - Audio URLs
   - Cover images

5. **Mock playlists** - Generate playlists
   - Playlist metadata
   - Track associations
   - Public/private settings

6. **Export as JSON** [NEW]
   - Generate complete dataset as JSON
   - Includes all entities and relationships
   - Ready for documentation/testing

7. **Validate data** [NEW]
   - Validate JSON files
   - Check data structure
   - View dataset statistics

8. **View dataset info** [NEW]
   - Display dataset metadata
   - Show entity counts
   - Review generation timestamp

## Data Structure

The generated mock data follows this structure:

```json
{
  "metadata": {
    "version": "1.0",
    "generatedAt": "2026-03-15T10:30:00Z",
    "datasetName": "development",
    "totalUsers": 50,
    "totalArtists": 15,
    "totalTracks": 100,
    "totalPlaylists": 25
  },
  "users": [...],
  "profiles": [...],
  "artists": [...],
  "tracks": [...],
  "playlists": [...],
  "playlistTracks": [...]
}
```

See [MOCK_DATA_SCHEMA.json](../docs/MOCK_DATA_SCHEMA.json) for complete schema.

## File Organization

```
mock-data/
├── main.go                          # Entry point with improved menu
│
├── model/
│   ├── csv.go                       # Legacy CSV model
│   ├── mock_data.go                 # NEW: Complete mock data structures
│   └── config.go                    # NEW: Configuration and presets
│
├── generator/
│   ├── api/
│   │   ├── dto.go                   # UPDATED: Enhanced API DTOs
│   │   ├── user_generator.go        # Users and profiles
│   │   ├── artist_generator.go      # Artist generation
│   │   ├── track_generator.go       # Track generation
│   │   └── playlist_generator.go    # Playlist generation
│   │
│   └── query/
│       ├── user_generator.go        # SQL-based user generation
│       └── profile_generator.go     # SQL-based profile generation
│
├── export/
│   └── export.go                    # NEW: JSON and SQL export utilities
│
├── validation/
│   └── validation.go                # NEW: Schema validation
│
├── util/
│   └── file.go                      # IMPROVED: File I/O utilities
│
├── assets/
│   ├── usernames.txt                # Seed data for generation
│   ├── passwords.txt
│   ├── artist_names.txt
│   └── track_titles.txt
│
└── output/
    ├── json/                        # NEW: Generated JSON files
    └── sql/                         # NEW: Generated SQL files
```

## Examples

### Generate Users via API
```
1. Select "1. Mock users"
2. Enter quantity: 50
3. Select method: "1. API calls"
✓ 50 users created and saved to accounts.csv
```

### Export Dataset as JSON
```
1. Select "6. Export as JSON"
2. Enter dataset name: "development"
✓ JSON exported: output/json/mock_data_development_2026-03-15.json
  - Users: 50 | Artists: 15 | Tracks: 100 | Playlists: 25
```

### Validate Dataset
```
1. Select "7. Validate data"
2. Enter path: "output/json/mock_data_development_2026-03-15.json"
✓ JSON structure is valid
  - Users: 50
  - Artists: 15
  - Tracks: 100
  - Playlists: 25
```

## Generation Workflow

### Typical Flow
```
1. Start mock-data-generator
2. Generate users (API or SQL method)
3. Generate artists (API)
4. Generate tracks (API)
5. Generate playlists (API)
6. Export as JSON
7. Validate data
8. Use generated JSON in tests or documentation
```

### Advanced: Using Presets
```go
// In Go code
config := model.GetPresetConfig("standard")
// Use config.UserCount, config.ArtistCount, etc.
// for automated dataset generation
```

## Output Files

### JSON Files
```
output/json/
├── mock_data_development_2026-03-15.json
├── mock_data_testing_2026-03-15.json
└── mock_data_production_2026-03-15.json
```

Each file contains:
- Complete mock data structure
- Metadata with generation info
- All entities and relationships
- Ready for import or testing

### SQL Files
```
output/sql/
├── mock_data_2026-03-15_10-30-45.sql
└── mock_data_2026-03-16_14-20-30.sql
```

Each file contains:
- INSERT statements for users
- INSERT statements for artists
- INSERT statements for playlists
- Ready for database import

## Integration with Tests

### Go Tests
```go
import "github.com/jimtrung/music-streaming-mock-data/util"

func TestUserAPI(t *testing.T) {
    // Load mock data
    mockData, _ := util.ReadFileJSON("output/json/mock_data_test.json")
    
    // Use in tests
    for _, user := range mockData.Users {
        // Run tests with mock data
    }
}
```

### API Integration
```bash
# Load data into database via SQL
psql -U postgres -d music-streaming < output/sql/mock_data_*.sql

# Or via API (see MOCK_DATA_MANAGEMENT.md for scripts)
```

## Configuration

### Default API Endpoint
```
http://127.0.0.1:5255
```

Modify in generator files if your backend runs on different port.

### Seed Data Files
Customize generation by editing:
- `assets/usernames.txt` - List of usernames
- `assets/passwords.txt` - List of passwords
- `assets/artist_names.txt` - List of artist names
- `assets/track_titles.txt` - List of track titles

## Troubleshooting

### "Failed to connect to API"
- Ensure backend is running on `http://127.0.0.1:5255`
- Check firewall settings
- Verify CORS is enabled

### "Invalid file format"
- Ensure CSV seed files have newline-separated values
- No headers in seed files (raw data only)

### "JSON validation failed"
- Verify schema file path
- Check JSON structure matches MOCK_DATA_SCHEMA.json
- Review validation error messages

## Next Steps

### Planned Enhancements (Phase 2)
1. **Automated Dataset Composition**
   - Combine multiple presets
   - Generate complex relationships

2. **CI/CD Integration**
   - Auto-generate test data
   - Run in pipelines

3. **Performance Analysis**
   - Measure generation time
   - Database import benchmarks

4. **Advanced Validation**
   - Referential integrity checks
   - Duplicate detection
   - Relationship validation

## Related Documentation

- **API Documentation**: [../docs/API_DOCUMENTATION.md](../docs/API_DOCUMENTATION.md)
- **Mock Data Schema**: [../docs/MOCK_DATA_SCHEMA.json](../docs/MOCK_DATA_SCHEMA.json)
- **Management Guide**: [../docs/MOCK_DATA_MANAGEMENT.md](../docs/MOCK_DATA_MANAGEMENT.md)
- **Enhancement Roadmap**: [../docs/MOCK_DATA_GENERATOR_ENHANCEMENT.md](../docs/MOCK_DATA_GENERATOR_ENHANCEMENT.md)

## Development

### Adding New Generators

To add a new entity generator:

1. Create `generator/api/new_entity_generator.go`
2. Implement generation function
3. Update `main.go` menu
4. Add to MockDataJSON model
5. Export in export.go

### Code style
- Follow Go conventions
- Use error returns, not panics
- Add documentation comments
- Use meaningful variable names

## License

Part of the Music Streaming project.

---

**Version**: 1.0  
**Last Updated**: March 15, 2026  
**Maintainer**: Backend Team
"# music-streaming-mock" 
