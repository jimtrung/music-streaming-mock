package model

// GenerationConfig defines configuration for data generation
type GenerationConfig struct {
	// API Configuration
	APIBaseURL string
	Timeout    int // seconds

	// Generation Presets
	PresetType string // "minimal", "standard", "stress", "e2e"

	// Data Generation Parameters
	UserCount     int
	ArtistCount   int
	TrackCount    int
	PlaylistCount int
	TracksPerList int

	// Export Configuration
	ExportFormat string // "json", "sql", "both"
	OutputDir    string

	// Data Seed Files
	UsernamesFile   string
	PasswordsFile   string
	ArtistNamesFile string
	TrackTitlesFile string
}

// GetPresetConfig returns configuration for a named preset
func GetPresetConfig(preset string) *GenerationConfig {
	switch preset {
	case "minimal":
		return &GenerationConfig{
			PresetType:      "minimal",
			UserCount:       10,
			ArtistCount:     3,
			TrackCount:      10,
			PlaylistCount:   2,
			TracksPerList:   3,
			ExportFormat:    "both",
			OutputDir:       "output",
			UsernamesFile:   "./assets/usernames.txt",
			PasswordsFile:   "./assets/passwords.txt",
			ArtistNamesFile: "./assets/artist_names.txt",
			TrackTitlesFile: "./assets/track_titles.txt",
		}
	case "standard":
		return &GenerationConfig{
			PresetType:      "standard",
			UserCount:       50,
			ArtistCount:     15,
			TrackCount:      100,
			PlaylistCount:   25,
			TracksPerList:   5,
			ExportFormat:    "both",
			OutputDir:       "output",
			UsernamesFile:   "./assets/usernames.txt",
			PasswordsFile:   "./assets/passwords.txt",
			ArtistNamesFile: "./assets/artist_names.txt",
			TrackTitlesFile: "./assets/track_titles.txt",
		}
	case "stress":
		return &GenerationConfig{
			PresetType:      "stress",
			UserCount:       500,
			ArtistCount:     100,
			TrackCount:      1000,
			PlaylistCount:   250,
			TracksPerList:   8,
			ExportFormat:    "json",
			OutputDir:       "output",
			UsernamesFile:   "./assets/usernames.txt",
			PasswordsFile:   "./assets/passwords.txt",
			ArtistNamesFile: "./assets/artist_names.txt",
			TrackTitlesFile: "./assets/track_titles.txt",
		}
	case "e2e":
		return &GenerationConfig{
			PresetType:      "e2e",
			UserCount:       30,
			ArtistCount:     10,
			TrackCount:      75,
			PlaylistCount:   20,
			TracksPerList:   5,
			ExportFormat:    "both",
			OutputDir:       "output",
			UsernamesFile:   "./assets/usernames.txt",
			PasswordsFile:   "./assets/passwords.txt",
			ArtistNamesFile: "./assets/artist_names.txt",
			TrackTitlesFile: "./assets/track_titles.txt",
		}
	default:
		// Return standard config as default
		return GetPresetConfig("standard")
	}
}

// DatasetInfo contains metadata about a generated dataset
type DatasetInfo struct {
	Name        string
	Version     string
	Description string
	GeneratedAt string
	Preset      string
	Stats       DatasetStats
}

// DatasetStats contains statistics about a dataset
type DatasetStats struct {
	TotalUsers        int
	TotalArtists      int
	TotalTracks       int
	TotalPlaylists    int
	RelationshipCount int
}

// Mock datasets can be composed from multiple presets
type DatasetComposer struct {
	Name        string
	Description string
	Datasets    []*MockDataJSON
}

// Compose merges multiple datasets into one
func (dc *DatasetComposer) Compose() *MockDataJSON {
	composed := NewMockData(dc.Name, dc.Description)

	for _, dataset := range dc.Datasets {
		composed.Users = append(composed.Users, dataset.Users...)
		composed.Artists = append(composed.Artists, dataset.Artists...)
		composed.Tracks = append(composed.Tracks, dataset.Tracks...)
		composed.Playlists = append(composed.Playlists, dataset.Playlists...)
		composed.PlaylistTracks = append(composed.PlaylistTracks, dataset.PlaylistTracks...)
		composed.Profiles = append(composed.Profiles, dataset.Profiles...)
	}

	return composed
}
