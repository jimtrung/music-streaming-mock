# Generator Refactoring - Data Repository Architecture

## Overview

The generator system has been refactored to implement a **centralized data repository** that persists generated mock data between runs, making it reusable and well-organized.

### Key Improvements

#### 1. **Centralized Data Repository** (`generator/repository/repository.go`)
- **Purpose**: Manages all generated mock data in an organized structure
- **Location**: `./data/` directory (auto-created)
- **Features**:
  - Persistent JSON storage for users, artists, tracks, playlists, profiles
  - Automatic indexing of all datasets
  - Thread-safe operations (RWMutex)
  - Metadata tracking (creation time, entity counts, status)

#### 2. **Generator Context** (`generator/generators/generator_context.go`)
- **Purpose**: Shared context for all generators, manages in-memory data
- **Features**:
  - Holds all generated entities (users, artists, tracks, playlists)
  - Manages asset pools (usernames, passwords, artist names, track titles)
  - Provides random selection helpers (GetRandomUsername, etc.)
  - Automatic saving to repository after each generation
  - Statistics tracking (GetStats)
  - Compilation to MockDataJSON format

#### 3. **Refactored Generators** (`generator/generators/generators.go`)
- **Generators**: UserGenerator, ProfileGenerator, ArtistGenerator, TrackGenerator, PlaylistGenerator
- **Features**:
  - Use GeneratorContext for data management
  - Automatic persistence to repository
  - Proper error handling
  - Progress reporting
  - UUID generation with bcrypt password hashing

#### 4. **Dataset Manager** (`generator/repository/dataset_manager.go`)
- **Purpose**: High-level operations on datasets
- **Features**:
  - Create/list/delete datasets
  - View dataset metadata and info
  - Clone datasets
  - Merge datasets
  - Archive datasets
  - Generate reports

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     main.go (CLI)                       │
│  - Menu system                                          │
│  - User input handling                                  │
└──────────────┬──────────────────────────────────────────┘
               │
┌──────────────▼──────────────────────────────────────────┐
│            GeneratorContext                             │
│  - Holds in-memory data (users, artists, tracks)        │
│  - Manages asset pools                                  │
│  - Provides random selection                            │
└──────────┬──────────────────────┬─────────────────────┬─┘
           │                      │                     │
    ┌──────▼────────┐    ┌────────▼──────────┐   ┌─────▼──────────┐
    │ UserGenerator │    │ ArtistGenerator   │   │ TrackGenerator │
    │ ProfileGen    │    │ PlaylistGenerator │   │ + others       │
    │ etc.          │    │                   │   │                │
    └──────┬────────┘    └────────┬──────────┘   └─────┬──────────┘
           │                      │                     │
           └──────────────────────┼─────────────────────┘
                                  │ AddUsers, AddArtists, etc.
                    ┌─────────────▼──────────────┐
                    │  DataRepository            │
                    │  - Manage datasets         │
                    │  - Save/Load data (JSON)   │
                    │  - Thread-safe storage     │
                    │  - Index management        │
                    └─────────────┬──────────────┘
                                  │
                    ┌─────────────▼──────────────┐
                    │  ./data/ (File System)     │
                    │  ├── dataset1/             │
                    │  │   ├── users.json        │
                    │  │   ├── artists.json      │
                    │  │   ├── tracks.json       │
                    │  │   └── playlists.json    │
                    │  ├── dataset2/             │
                    │  │   └── ...               │
                    │  └── index.json            │
                    └────────────────────────────┘
```

## File Structure

```
mock-data/
├── data/                           # ← REPOSITORY STORAGE (auto-created)
│   ├── development/
│   │   ├── users.json
│   │   ├── profiles.json
│   │   ├── artists.json
│   │   ├── tracks.json
│   │   ├── playlists.json
│   │   └── complete_dataset.json
│   ├── testing/
│   │   ├── users.json
│   │   └── ...
│   └── index.json
│
├── generator/
│   ├── repository/
│   │   ├── repository.go          # ← Core repository (NEW)
│   │   └── dataset_manager.go     # ← Dataset operations (NEW)
│   ├── generators/
│   │   ├── generator_context.go   # ← Shared context (NEW/REFACTORED)
│   │   └── generators.go          # ← All generators (REFACTORED)
│   ├── api/
│   │   ├── user_generator.go
│   │   ├── artist_generator.go
│   │   └── ... (legacy API generators)
│   └── query/
│       └── ... (legacy SQL generators)
│
└── main.go                         # ← Enhanced with repository support
```

## Usage Workflow

### Scenario 1: Generate Complete Dataset and Save

```go
// 1. Initialize repository
dataRepo, _ := repository.NewDataRepository("./data")
dataMgr := repository.NewDatasetManager(dataRepo)

// 2. Create dataset
dataMgr.CreateNewDataset("dev", "Development dataset")

// 3. Initialize generator context
ctx, _ := generators.NewGeneratorContext(dataRepo, "dev")

// 4. Generate entities
userGen := generators.NewUserGenerator(ctx)
userGen.Generate(50)  // ← Auto-saves to repo

artistGen := generators.NewArtistGenerator(ctx)
artistGen.Generate(15)  // ← Auto-saves to repo

trackGen := generators.NewTrackGenerator(ctx)
trackGen.Generate(100)  // ← Auto-saves to repo

profileGen := generators.NewProfileGenerator(ctx)
profileGen.Generate()  // ← Auto-saves to repo

playlistGen := generators.NewPlaylistGenerator(ctx)
playlistGen.Generate(20, 10)  // ← Auto-saves to repo

// 5. Compile and save complete dataset
mockData := ctx.CompileToMockData()
dataRepo.SaveAllData("dev", mockData)
// Result: ./data/dev/complete_dataset.json
```

### Scenario 2: Reuse Existing Dataset

```go
// 1. List available datasets
datasets := dataRepo.ListDatasets()
for _, dataset := range datasets {
    fmt.Printf("Dataset: %s (%d users)\n", dataset.Name, dataset.UserCount)
}

// 2. Load existing dataset
mockData, _ := dataRepo.LoadAllData("dev")
fmt.Printf("Loaded: %d users, %d artists\n", len(mockData.Users), len(mockData.Artists))

// 3. Use in tests/API calls
for _, user := range mockData.Users {
    // Use user data for testing...
}
```

### Scenario 3: Clone and Modify Dataset

```go
// Clone existing dataset
dataMgr.CloneDataset("dev", "dev_backup")

// Load cloned dataset, modify it
mockData, _ := dataRepo.LoadAllData("dev_backup")
mockData.Users = append(mockData.Users, newUsers...)

// Save modified dataset
dataRepo.SaveAllData("dev_backup", mockData)
```

## Data Storage Format

### Repository Index (`./data/index.json`)

```json
{
  "version": "1.0",
  "lastUpdated": "2026-03-15T14:30:00Z",
  "datasets": {
    "development": {
      "name": "development",
      "createdAt": "2026-03-15T10:30:00Z",
      "path": "./data/development",
      "userCount": 50,
      "artistCount": 15,
      "trackCount": 100,
      "playlistCount": 20,
      "profileCount": 50,
      "status": "complete"
    },
    "testing": {
      "name": "testing",
      "createdAt": "2026-03-15T12:00:00Z",
      "path": "./data/testing",
      "userCount": 10,
      "artistCount": 3,
      "trackCount": 20,
      "playlistCount": 5,
      "profileCount": 10,
      "status": "complete"
    }
  }
}
```

### Dataset Structure (`./data/{name}/`)

**Files generated:**
- `users.json` - All generated users with hashed passwords
- `profiles.json` - User profiles
- `artists.json` - Artists linked to users
- `tracks.json` - Tracks linked to artists
- `playlists.json` - Playlists with track lists
- `complete_dataset.json` - Complete MockDataJSON with all entities

## Refactored Generator Functions

### UserGenerator

```go
userGen := generators.NewUserGenerator(ctx)
userGen.Generate(50)  // Generates 50 users, auto-saves to repo
```

Features:
- Uses bcrypt for password hashing
- 10% premium users automatically
- Auto-saves progress
- Progress reporting every 10 users

### ArtistGenerator

```go
artistGen := generators.NewArtistGenerator(ctx)
artistGen.Generate(15)  // Generates artists from existing users
```

Features:
- Requires users to exist first
- 5% verified artists
- One artist per user (quantity limited to user count)
- Auto-saves

### TrackGenerator

```go
trackGen := generators.NewTrackGenerator(ctx)
trackGen.Generate(100)  // Generates 100 tracks
```

Features:
- Distributes tracks among artists
- Random genres assigned
- Duration between 3-5 minutes
- Realistic URLs for audio/covers
- Auto-saves

### PlaylistGenerator

```go
playlistGen := generators.NewPlaylistGenerator(ctx)
playlistGen.Generate(20, 10)  // 20 playlists with ~10 tracks each
```

Features:
- Creates from existing users and tracks
- Random track selection per playlist
- 50% public playlists
- Tracks numbered correctly
- Auto-saves

### ProfileGenerator

```go
profileGen := generators.NewProfileGenerator(ctx)
profileGen.Generate()  // Creates profiles for all users
```

Features:
- One profile per user
- Auto-generated bio
- Auto-saves

## API (DataRepository)

### Core Methods

```go
// Create new dataset
repo.CreateDataset(name string)

// Save data
repo.SaveUsers(datasetName, users)
repo.SaveArtists(datasetName, artists)
repo.SaveTracks(datasetName, tracks)
repo.SavePlaylists(datasetName, playlists)
repo.SaveAllData(datasetName, mockData)

// Load data
repo.LoadUsers(datasetName)
repo.LoadArtists(datasetName)
repo.LoadTracks(datasetName)
repo.LoadPlaylists(datasetName)
repo.LoadAllData(datasetName)

// Manage datasets
repo.ListDatasets()
repo.GetDatasetMetadata(datasetName)
repo.DeleteDataset(datasetName)
repo.GetDataPath(datasetName)
```

## API (DatasetManager)

```go
// Dataset operations
mgr.CreateNewDataset(name, description)
mgr.DeleteDataset(name)
mgr.ListAllDatasets()
mgr.CloneDataset(source, target)
mgr.MergeDatasets(source, target)
mgr.ArchiveDataset(name)

// Information
mgr.PrintDatasetInfo(name)
mgr.GenerateDatasetReport(name)
mgr.GetDatasetPath(name)

// Data access
mgr.LoadDataset(name)
mgr.SaveDataset(name, mockData)
```

## Thread Safety

- All repository operations are **thread-safe** using `sync.RWMutex`
- GeneratorContext uses mutex for concurrent access
- Multiple datasets can be generated in parallel
- Safe for concurrent reads/writes

## Backward Compatibility

✅ **Fully maintained** - Old API generators still work:
- `api.GenerateUser()` - API-based generation
- `api.GenerateArtist()` - API-based generation
- `query.GenerateUser()` - SQL generation
- Legacy asset files still used

New in-memory generators are **additional**, not replacements.

## Benefits

| Feature | Benefit |
|---------|---------|
| **Centralized Storage** | All data in one organized location, easy to find and reuse |
| **Persistence** | Data survives between program runs |
| **Metadata Tracking** | Know what's in each dataset without loading it |
| **Automatic Saving** | No manual save calls needed |
| **Reusability** | Generate once, use multiple times |
| **Cloning/Merging** | Create variations and combine datasets |
| **Thread-Safe** | Safe for concurrent operations |
| **JSON Format** | Human-readable, tool-friendly |
| **Complete Data** | Includes all relationships (users→artists, tracks→playlists) |

## Migration from Old System

Old system:
```
- User data scattered in assets/accounts.csv
- Artist IDs in assets/artist_ids.txt
- Track IDs in assets/track_ids.txt
- SQL queries in ./query/ with timestamps
```

New system:
```
- All data organized in ./data/{dataset_name}/
- Complete relationships maintained
- Single index for all datasets
- JSON + SQL export capabilities
```

**Transition**: Both systems work together. New generators use repository, old generators still work.

## Next Steps

1. **Use in-memory generators** for fast data generation
2. **Archive completed datasets** to mark them as finalized
3. **Clone datasets** for variations (dev, testing, staging)
4. **Merge datasets** to combine data
5. **Export to backend** using complete JSON files

---

**Status**: ✅ Production Ready
**Version**: 1.0
**Last Updated**: March 15, 2026
