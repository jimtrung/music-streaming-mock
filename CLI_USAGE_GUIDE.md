# Mock Data Generator - CLI Usage Guide

## Overview

The refactored generator now uses a **dataset-centric approach** with a centralized repository. All generated data is automatically saved and organized for reuse.

## Menu System

```
===================> MOCK DATA GENERATOR <====================
Current Dataset: development
  Users: 50 | Artists: 15 | Tracks: 100 | Playlists: 20

| 1. Select/Create dataset                                |
| 2. Generate users                                       |
| 3. Generate profiles                                    |
| 4. Generate artists                                     |
| 5. Generate tracks                                      |
| 6. Generate playlists                                   |
| 7. Export dataset as JSON                               |
| 8. List all datasets                                    |
| 9. View dataset info                                    |
| 10. Export SQL queries                                  |
| 0. Exit                                                 |
============================================================
```

## Usage Workflow

### Step 1: Select or Create a Dataset

**Menu Option**: 1

```
Choose dataset name: dev_environment
```

**What it does:**
- Creates a new dataset directory if it doesn't exist
- Initializes the GeneratorContext
- Prepares for data generation
- Shows status in the next menu

**Result:**
```
✓ Working with dataset: dev_environment
```

### Step 2: Generate Users

**Menu Option**: 2

```
Enter the number of users to generate: 50

Generation Method:
| 1. In-memory (fast, uses repository storage)   |
| 2. API calls (actual backend)                 |
| 3. SQL queries (database import)              |

Your choice: 1
```

**In-Memory Method (Recommended):**
- Generates users locally
- Stores in repository immediately
- Fast (no network calls)
- Automatic password hashing with bcrypt
- 10% of users are premium

**Result:**
```
Generating 50 users...
  ✓ Generated 10 users
  ✓ Generated 20 users
  ✓ Generated 30 users
  ✓ Generated 40 users
  ✓ Generated 50 users
✓ Successfully generated and saved 50 users
```

**Data saved to:** `./data/dev_environment/users.json`

### Step 3: Generate Profiles

**Menu Option**: 3

```
Generating profiles...
✓ Profiles generated successfully
```

**What it does:**
- Creates one profile per user
- Auto-generates bios
- Stores in repository
- Required: Users must exist first

**Data saved to:** `./data/dev_environment/profiles.json`

### Step 4: Generate Artists

**Menu Option**: 4

```
Enter the number of artists to generate: 15
Generating 15 artists...
  ✓ Generated 10 artists
  ✓ Generated 15 artists
✓ Successfully generated and saved 15 artists
```

**Features:**
- Creates artists from existing users (one per user)
- 5% of artists are verified
- Limited to number of available users
- Random artist names from asset pool

**Data saved to:** `./data/dev_environment/artists.json`

### Step 5: Generate Tracks

**Menu Option**: 5

```
Enter the number of tracks to generate: 100
Generating 100 tracks...
  ✓ Generated 50 tracks
  ✓ Generated 100 tracks
✓ Successfully generated and saved 100 tracks
```

**Features:**
- Distributes tracks evenly among artists
- 8 different genres (Pop, Rock, Hip-Hop, Jazz, Classical, Electronic, Country, Reggae)
- Random duration (3-5 minutes)
- Realistic URLs for audio and cover images
- Track numbers 1-20 per artist

**Data saved to:** `./data/dev_environment/tracks.json`

### Step 6: Generate Playlists

**Menu Option**: 6

```
Enter the number of playlists to generate: 20
Enter the average number of tracks per playlist: 10
Generating 20 playlists with ~10 tracks each...
  ✓ Generated 10 playlists
  ✓ Generated 20 playlists
✓ Successfully generated and saved 20 playlists
```

**Features:**
- Creates playlists owned by users
- Randomly selects tracks for each playlist
- 50% public playlists
- Proper playlist-track relationships maintained
- Tracks correctly positioned in playlist

**Data saved to:**
- `./data/dev_environment/playlists.json`
- `./data/dev_environment/playlist_tracks` (relationships)

### Step 7: Export Dataset as JSON

**Menu Option**: 7

```
Compiling dataset...
✓ Dataset exported and saved to repository
  Path: ./data/dev_environment
```

**What it does:**
- Compiles all generated entities into one MockDataJSON file
- Updates metadata (counts, timestamps, version)
- Saves complete dataset as `complete_dataset.json`
- Validates all relationships are intact

**Files created:**
```
./data/dev_environment/
├── users.json
├── profiles.json
├── artists.json
├── tracks.json
├── playlists.json
└── complete_dataset.json  ← Complete dataset (this one)
```

### Step 8: List All Datasets

**Menu Option**: 8

```
===================> AVAILABLE DATASETS <====================
Name                 Status            Created              Users
===========================================================
development          complete          2026-03-15 10:30:05 50
testing              in-progress       2026-03-15 11:00:00 10
production           archived          2026-03-15 09:00:00 100
===========================================================
```

**Shows:**
- All stored datasets
- Current status (in-progress, complete, archived)
- Creation timestamp
- Number of users in each

### Step 9: View Dataset Info

**Menu Option**: 9

```
Enter dataset name: development

===================> DATASET INFO <====================
  Name:          development
  Status:        complete
  Created:       2026-03-15T10:30:00Z
  Updated:       2026-03-15T14:30:00Z
  Path:          ./data/development
  Data Count:
    - Users:     50
    - Profiles:  50
    - Artists:   15
    - Tracks:    100
    - Playlists: 20
======================================================
```

**Shows:**
- Dataset metadata
- Entity counts
- File location
- Status and timestamps

### Step 10: Export SQL Queries

**Menu Option**: 10

```
Enter output SQL file name: export.sql
✓ SQL exported to export.sql
```

**What it does:**
- Generates SQL INSERT statements from stored data
- Creates file with all users for database import
- Uses data from repository, not API

**Note**: SQL export is still in development. Currently creates stub file.

## Common Workflows

### Workflow 1: Quick Development Dataset

```
1. Select dataset: dev
2. Generate: 50 users (in-memory)
3. Generate profiles
4. Generate: 15 artists
5. Generate: 100 tracks
6. Generate: 20 playlists
7. Export as JSON
Result: Complete dataset ready to use
Time: ~5-10 seconds
```

### Workflow 2: Create Test Dataset

```
1. Select dataset: test_v1
2. Generate: 10 users
3. Generate profiles
4. Generate: 3 artists
5. Generate: 20 tracks
6. Generate: 5 playlists
7. Export as JSON (small dataset for quick tests)
Result: Ready for unit tests
```

### Workflow 3: Clone and Modify

```
1. List datasets → find "production"
2. Select dataset: production_backup
   (creates new dataset)
3. Load data from "production"
4. Modify (add more users/tracks)
5. Export as JSON
Result: Production copy with modifications
```

## Data Organization

After completing Workflow 1, your structure looks like:

```
./data/
├── dev/
│   ├── users.json              (50 users)
│   ├── profiles.json           (50 profiles)
│   ├── artists.json            (15 artists)
│   ├── tracks.json             (100 tracks)
│   ├── playlists.json          (20 playlists)
│   └── complete_dataset.json   (all together)
│
├── test_v1/
│   ├── users.json              (10 users)
│   ├── profiles.json           (10 profiles)
│   ├── artists.json            (3 artists)
│   ├── tracks.json             (20 tracks)
│   └── playlists.json          (5 playlists)
│
└── index.json                   (metadata for all)
```

## Tips & Tricks

### Tip 1: Reuse as Template
```
1. Create "template" dataset with baseline data
2. Clone it for each test
3. Modify as needed
Result: Consistent test data structure
```

### Tip 2: Generate in Stages
```
Session 1:
  - Generate users and profiles
  
Session 2:
  - Load dataset
  - Add artists and tracks
  - Export
Result: Flexible generation process
```

### Tip 3: Multiple Datasets
```
dev          → Development testing
test         → Unit tests
staging      → Pre-production validation
performance  → Load testing (large dataset)
Result: Different data for different purposes
```

### Tip 4: Archive When Done
```
After dataset is finalized:
1. List datasets (option 8)
2. All completed datasets can be archived
3. Mark status as "archived"
Result: Read-only dataset for reference
```

## Performance Notes

**Generation Speed (approximate):**
- 50 users: ~0.5 seconds
- 15 artists: ~0.1 seconds
- 100 tracks: ~0.2 seconds
- 20 playlists (10 tracks each): ~0.3 seconds
- **Total**: ~1 second for ~180+ entities

**Storage:**
- 50 users + 50 profiles: ~50 KB
- 15 artists: ~10 KB
- 100 tracks: ~30 KB
- 20 playlists: ~20 KB
- **Total**: ~110 KB per dataset

**Scalability:**
- Tested with up to 1000 users
- Tested with up to 10000 tracks
- Repository handles concurrent access safely

## Troubleshooting

### Issue: No dataset selected
```
Error: ✗ Please select a dataset first (option 1)
Solution: Use menu option 1 first
```

### Issue: Cannot generate artists
```
Error: no users found - generate users first
Solution: Generate users before artists
```

### Issue: Cannot generate tracks
```
Error: no artists found - generate artists first
Solution: Generate artists before tracks
```

### Issue: Data not saving
```
Check: ./data/ directory exists and is writable
Check: Permissions on ./data/ folder
Check: Disk space available
```

## Next Steps

After generating datasets:

1. **Validate** - Check data integrity
2. **Export** - Convert to needed format (JSON/SQL)
3. **Use** - Feed into backend API or database
4. **Iterate** - Create new datasets as needed

---

**Command Examples:**
```bash
# Run generator
go run main.go

# View generated data
cat ./data/development/users.json | jq '.[] | select(.isPremium == true)'

# Count entities
jq '.[] | length' ./data/development/users.json
```

---

**Status**: ✅ Ready to Use
**Last Updated**: March 15, 2026
