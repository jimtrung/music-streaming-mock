# Quick Reference - Generator System

## 🚀 Quick Start (2 minutes)

```bash
# 1. Run generator
go run main.go

# 2. Create dataset (Menu Option 1)
Select dataset: dev

# 3. Generate users (Menu Option 2)
Enter quantity: 50
Choose method: 1 (in-memory)

# 4. Generate supporting entities (Menu Options 3-6)
- Profiles
- Artists (15)
- Tracks (100)
- Playlists (20, 10 tracks each)

# 5. Export (Menu Option 7)
Export dataset as JSON

# Result: ./data/dev/complete_dataset.json
```

## 📁 Data Structure

```
./data/
├── dev/
│   ├── users.json           (50 users)
│   ├── profiles.json        (50 profiles)
│   ├── artists.json         (15 artists)
│   ├── tracks.json          (100 tracks)
│   ├── playlists.json       (20 playlists)
│   └── complete_dataset.json
├── test/
│   └── ...
└── index.json (metadata)
```

## 🎯 Menu Options

| Option | Purpose | What It Does |
|--------|---------|-------------|
| 1 | Select Dataset | Create/switch to dataset |
| 2 | Generate Users | Create user entities |
| 3 | Generate Profiles | Create user profiles |
| 4 | Generate Artists | Create artist entities |
| 5 | Generate Tracks | Create track entities |
| 6 | Generate Playlists | Create playlist entities |
| 7 | Export JSON | Save complete dataset |
| 8 | List Datasets | Show all datasets |
| 9 | View Dataset Info | Show dataset details |
| 10 | Export SQL | Generate SQL queries |
| 0 | Exit | Close program |

## 🔧 Generation Methods

### Method 1: In-Memory (Recommended ⭐)
- **Speed**: 50-100 users in 0.5 seconds
- **Location**: Auto-saves to `./data/{name}/`
- **Use**: Use this for testing
- **Example**: 
  ```
  Choose method: 1
  Generating 50 users...
  ✓ Successfully generated and saved 50 users
  ```

### Method 2: API Calls (Legacy)
- **Speed**: ~30 seconds for 50 users
- **Requires**: Backend running on localhost:5255
- **Use**: For real API testing
- **Example**:
  ```
  Choose method: 2
  Generating 50 users via API...
  ```

### Method 3: SQL Queries (Legacy)
- **Speed**: Instant
- **Requires**: Database schema ready
- **Use**: For database import
- **Example**:
  ```
  Choose method: 3
  Generating 50 users as SQL...
  ```

## 📊 Generated Entity Counts

**Default Generation:**
- Users: 50 (10% premium, bcrypt hashed)
- Profiles: 1 per user (50)
- Artists: 15 (5% verified)
- Tracks: 100 (distributed among artists)
- Playlists: 20 (50% public, ~10 tracks each)

## 💾 Data Loading

```go
// Load repository
repo, _ := repository.NewDataRepository("./data")

// Load dataset
mockData, _ := repo.LoadAllData("dev")

// Access data
users := mockData.Users          // []UserJSON
artists := mockData.Artists      // []ArtistJSON
tracks := mockData.Tracks        // []TrackJSON
playlists := mockData.Playlists  // []PlaylistJSON

// Use in tests
for _, user := range users {
    // Test with this user
}
```

## 🗂️ File Locations

| What | Where |
|------|-------|
| Generated Data | `./data/{dataset_name}/` |
| Metadata Index | `./data/index.json` |
| Documentation | Root folder (*.md files) |
| Source Code | `generator/` folder |
| Legacy Generators | `generator/api/` and `generator/query/` |
| Asset Lists | `assets/` folder |

## ⚡ Performance

| Operation | Time |
|-----------|------|
| Generate 50 users | 0.5s |
| Generate 15 artists | 0.1s |
| Generate 100 tracks | 0.2s |
| Generate 20 playlists | 0.3s |
| **Total** | ~1 second |
| Reuse existing dataset | <1ms |

## 🔄 Dataset Operations

### Clone Dataset
```golang
mgr := repository.NewDatasetManager(repo)
mgr.ClonDataset("dev", "dev_backup")
```

### List All Datasets
```golang
mgr.ListAllDatasets()
// Shows: Name, Status, Created, User Count
```

### View Dataset Info
```golang
mgr.PrintDatasetInfo("dev")
// Shows: Name, Path, Status, Entity Counts
```

### Archive Dataset
```golang
mgr.ArchiveDataset("dev")
// Marks as finalized/read-only
```

## 🎨 Data Features

### User Data
- ID (UUID)
- Username (from asset pool)
- Email (auto-generated)
- Password (bcrypt hashed)
- Role (listener)
- Provider (local)
- IsVerified (true)
- IsPremium (10% of users)

### Artist Data
- ID (UUID)
- UserId (linked to user)
- Name (from asset pool)
- IsVerified (5% of artists)
- CreatedAt/UpdatedAt (timestamps)

### Track Data
- ID (UUID)
- ArtistId (linked to artist)
- Title (from asset pool)
- Genre (8 types)
- Duration (3-5 minutes)
- AudioUrl (realistic URL)
- CoverUrl (realistic URL)

### Playlist Data
- ID (UUID)
- OwnerId (linked to user)
- Name (from asset pool)
- IsPublic (50%)
- TrackIds (array)
- CreatedAt/UpdatedAt

## 🔐 Thread Safety

✅ All operations thread-safe using `sync.RWMutex`
✅ Multiple concurrent generators supported
✅ Safe for parallel dataset creation

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| "No dataset selected" | Use menu option 1 first |
| "No users found" | Generate users before artists |
| "No artists found" | Generate artists before tracks |
| Data not saving | Check `./data/` permissions |
| Wrong data format | Use in-memory method (option 1) |

## 📖 Documentation

| File | Contains |
|------|----------|
| `SUMMARY.md` | Executive overview |
| `GENERATOR_REFACTORING_COMPLETE.md` | Complete refactoring details |
| `generator/REPOSITORY_GUIDE.md` | Architecture & API reference |
| `CLI_USAGE_GUIDE.md` | Step-by-step workflows |
| `FILE_STRUCTURE_REFERENCE.md` | File descriptions |

## 🎓 Common Workflows

### Workflow 1: Create Dev Dataset
```
Option 1: dev
Option 2: 50 users
Option 3: profiles
Option 4: 15 artists
Option 5: 100 tracks
Option 6: 20 playlists
Option 7: export JSON
```

### Workflow 2: Create Multiple Datasets
```
Option 1: dev       (50 users, 100 tracks)
Option 1: test      (10 users, 20 tracks)
Option 1: perf      (500 users, 1000 tracks)
```

### Workflow 3: Reuse & Modify
```
Load: ./data/dev/complete_dataset.json
Edit: Add more users
Save: ./data/dev_modified/
```

## 🚀 Best Practices

1. **Use in-memory generation** (50-60x faster)
2. **Create separate datasets** for different purposes
3. **Archive datasets** when finalized
4. **Backup important** datasets
5. **Document** what each dataset is for
6. **Use metadata** to find the right dataset
7. **Clone** datasets for variations
8. **Reuse** instead of regenerating

## 📋 System Requirements

- Go 1.18+
- Disk space: ~100KB per standard dataset
- Memory: ~50MB for large datasets
- Permissions: Write access to `./data/`

## ✅ Verification

```bash
# Check data created
ls -la ./data/dev/

# View metadata
cat ./data/index.json | jq '.'

# Count entities
jq '.[] | length' ./data/dev/users.json

# View sample user
jq '.[0]' ./data/dev/users.json
```

## 🎯 Next Steps

1. Run `go run main.go`
2. Create a dataset
3. Generate entities
4. Check `./data/` folder
5. Load data in tests
6. Repeat for other datasets

---

**Status**: ✅ Ready to Use  
**Version**: 1.0  
**Last Updated**: March 15, 2026
