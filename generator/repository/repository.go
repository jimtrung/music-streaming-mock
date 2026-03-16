package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jimtrung/music-streaming-mock-data/model"
)

// DataRepository manages centralized storage of generated mock data
// Data persists across generators and can be reused
type DataRepository struct {
	mu        sync.RWMutex
	repoPath  string
	dataIndex *DataIndex
}

// DataIndex tracks all generated data with metadata
type DataIndex struct {
	Version     string                  `json:"version"`
	LastUpdated time.Time               `json:"lastUpdated"`
	Datasets    map[string]*DatasetMeta `json:"datasets"`
}

// DatasetMeta stores metadata about a generated dataset
type DatasetMeta struct {
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Path          string    `json:"path"`
	UserCount     int       `json:"userCount"`
	ArtistCount   int       `json:"artistCount"`
	TrackCount    int       `json:"trackCount"`
	PlaylistCount int       `json:"playlistCount"`
	ProfileCount  int       `json:"profileCount"`
	Status        string    `json:"status"` // "in-progress", "complete", "partial"
}

// GeneratedData holds the current collection of generated entities
type GeneratedData struct {
	Users          []model.UserJSON          `json:"users"`
	Profiles       []model.ProfileJSON       `json:"profiles"`
	Artists        []model.ArtistJSON        `json:"artists"`
	Tracks         []model.TrackJSON         `json:"tracks"`
	Playlists      []model.PlaylistJSON      `json:"playlists"`
	PlaylistTracks []model.PlaylistTrackJSON `json:"playlistTracks"`
}

// NewDataRepository creates a new repository instance
func NewDataRepository(repoPath string) (*DataRepository, error) {
	// Ensure repository path exists
	if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create repository path: %w", err)
	}

	repo := &DataRepository{
		repoPath: repoPath,
		dataIndex: &DataIndex{
			Version:  "1.0",
			Datasets: make(map[string]*DatasetMeta),
		},
	}

	// Load existing index
	if err := repo.loadIndex(); err != nil {
		log.Printf("Warning: Could not load index, starting fresh: %v\n", err)
	}

	return repo, nil
}

// GetDataPath returns the path for a specific dataset
func (r *DataRepository) GetDataPath(datasetName string) string {
	return filepath.Join(r.repoPath, datasetName)
}

// CreateDataset creates a new dataset with initial metadata
func (r *DataRepository) CreateDataset(datasetName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)

	// Create directory
	if err := os.MkdirAll(datasetPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dataset directory: %w", err)
	}

	// Add to index
	now := time.Now()
	r.dataIndex.Datasets[datasetName] = &DatasetMeta{
		Name:        datasetName,
		CreatedAt:   now,
		LastUpdated: now,
		Path:        datasetPath,
		Status:      "in-progress",
	}

	r.dataIndex.LastUpdated = time.Now()

	return r.saveIndex()
}

// SaveUsers saves generated users and updates index
func (r *DataRepository) SaveUsers(datasetName string, users []model.UserJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	usersFile := filepath.Join(datasetPath, "users.json")

	// Marshal users
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(usersFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	// Update index
	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.UserCount = len(users)
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// SaveProfiles saves generated profiles and updates index
func (r *DataRepository) SaveProfiles(datasetName string, profiles []model.ProfileJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	profilesFile := filepath.Join(datasetPath, "profiles.json")

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	if err := ioutil.WriteFile(profilesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write profiles file: %w", err)
	}

	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.ProfileCount = len(profiles)
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// SaveArtists saves generated artists and updates index
func (r *DataRepository) SaveArtists(datasetName string, artists []model.ArtistJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	artistsFile := filepath.Join(datasetPath, "artists.json")

	data, err := json.MarshalIndent(artists, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal artists: %w", err)
	}

	if err := ioutil.WriteFile(artistsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write artists file: %w", err)
	}

	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.ArtistCount = len(artists)
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// SaveTracks saves generated tracks and updates index
func (r *DataRepository) SaveTracks(datasetName string, tracks []model.TrackJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	tracksFile := filepath.Join(datasetPath, "tracks.json")

	data, err := json.MarshalIndent(tracks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tracks: %w", err)
	}

	if err := ioutil.WriteFile(tracksFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write tracks file: %w", err)
	}

	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.TrackCount = len(tracks)
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// SavePlaylists saves generated playlists and updates index
func (r *DataRepository) SavePlaylists(datasetName string, playlists []model.PlaylistJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	playlistsFile := filepath.Join(datasetPath, "playlists.json")

	data, err := json.MarshalIndent(playlists, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal playlists: %w", err)
	}

	if err := ioutil.WriteFile(playlistsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write playlists file: %w", err)
	}

	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.PlaylistCount = len(playlists)
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// SaveAllData saves entire generated dataset and marks as complete
func (r *DataRepository) SaveAllData(datasetName string, mockData *model.MockDataJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	completeFile := filepath.Join(datasetPath, "complete_dataset.json")

	data, err := json.MarshalIndent(mockData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal complete dataset: %w", err)
	}

	if err := ioutil.WriteFile(completeFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write complete dataset file: %w", err)
	}

	// Update index
	if meta, exists := r.dataIndex.Datasets[datasetName]; exists {
		meta.UserCount = len(mockData.Users)
		meta.ArtistCount = len(mockData.Artists)
		meta.TrackCount = len(mockData.Tracks)
		meta.PlaylistCount = len(mockData.Playlists)
		meta.ProfileCount = len(mockData.Profiles)
		meta.Status = "complete"
		meta.LastUpdated = time.Now()
		r.dataIndex.LastUpdated = time.Now()
	}

	return r.saveIndex()
}

// LoadUsers loads users from a dataset
func (r *DataRepository) LoadUsers(datasetName string) ([]model.UserJSON, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	usersFile := filepath.Join(datasetPath, "users.json")

	data, err := ioutil.ReadFile(usersFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	var users []model.UserJSON
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// LoadArtists loads artists from a dataset
func (r *DataRepository) LoadArtists(datasetName string) ([]model.ArtistJSON, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	artistsFile := filepath.Join(datasetPath, "artists.json")

	data, err := ioutil.ReadFile(artistsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read artists file: %w", err)
	}

	var artists []model.ArtistJSON
	if err := json.Unmarshal(data, &artists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	return artists, nil
}

// LoadTracks loads tracks from a dataset
func (r *DataRepository) LoadTracks(datasetName string) ([]model.TrackJSON, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	tracksFile := filepath.Join(datasetPath, "tracks.json")

	data, err := ioutil.ReadFile(tracksFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tracks file: %w", err)
	}

	var tracks []model.TrackJSON
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
	}

	return tracks, nil
}

// LoadPlaylists loads playlists from a dataset
func (r *DataRepository) LoadPlaylists(datasetName string) ([]model.PlaylistJSON, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	playlistsFile := filepath.Join(datasetPath, "playlists.json")

	data, err := ioutil.ReadFile(playlistsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read playlists file: %w", err)
	}

	var playlists []model.PlaylistJSON
	if err := json.Unmarshal(data, &playlists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlists: %w", err)
	}

	return playlists, nil
}

// LoadAllData loads complete dataset from repository
func (r *DataRepository) LoadAllData(datasetName string) (*model.MockDataJSON, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	completeFile := filepath.Join(datasetPath, "complete_dataset.json")

	data, err := ioutil.ReadFile(completeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read complete dataset file: %w", err)
	}

	var mockData model.MockDataJSON
	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal complete dataset: %w", err)
	}

	return &mockData, nil
}

// GetDatasetMetadata retrieves metadata for a dataset
func (r *DataRepository) GetDatasetMetadata(datasetName string) (*DatasetMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, exists := r.dataIndex.Datasets[datasetName]
	return meta, exists
}

// ListDatasets returns all available datasets
func (r *DataRepository) ListDatasets() []*DatasetMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasets := make([]*DatasetMeta, 0, len(r.dataIndex.Datasets))
	for _, meta := range r.dataIndex.Datasets {
		datasets = append(datasets, meta)
	}

	return datasets
}

// DeleteDataset removes a dataset from repository
func (r *DataRepository) DeleteDataset(datasetName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)

	// Remove directory
	if err := os.RemoveAll(datasetPath); err != nil {
		return fmt.Errorf("failed to delete dataset directory: %w", err)
	}

	// Remove from index
	delete(r.dataIndex.Datasets, datasetName)
	r.dataIndex.LastUpdated = time.Now()

	return r.saveIndex()
}

// AddUserIDs stores a list of generated user IDs
func (r *DataRepository) AddUserIDs(datasetName string, userIDs []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	datasetPath := r.GetDataPath(datasetName)
	idsFile := filepath.Join(datasetPath, "user_ids.txt")

	var content string
	for _, id := range userIDs {
		content += id + "\n"
	}

	if err := ioutil.WriteFile(idsFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write user IDs file: %w", err)
	}

	return nil
}

// GetUserIDs retrieves stored user IDs
func (r *DataRepository) GetUserIDs(datasetName string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	datasetPath := r.GetDataPath(datasetName)
	idsFile := filepath.Join(datasetPath, "user_ids.txt")

	data, err := ioutil.ReadFile(idsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read user IDs file: %w", err)
	}

	// Parse IDs
	var ids []string
	for _, id := range strings.Split(string(data), "\n") {
		if id != "" {
			ids = append(ids, strings.TrimSpace(id))
		}
	}

	return ids, nil
}

// saveIndex persists the data index to disk
func (r *DataRepository) saveIndex() error {
	indexPath := filepath.Join(r.repoPath, "index.json")

	data, err := json.MarshalIndent(r.dataIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := ioutil.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// loadIndex loads the data index from disk
func (r *DataRepository) loadIndex() error {
	indexPath := filepath.Join(r.repoPath, "index.json")

	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("index file not found: %w", err)
	}

	if err := json.Unmarshal(data, r.dataIndex); err != nil {
		return fmt.Errorf("failed to unmarshal index: %w", err)
	}

	return nil
}

// ExportSQLQueries generates SQL INSERT statements from stored data
func (r *DataRepository) ExportSQLQueries(datasetName string, outputPath string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, err := r.LoadUsers(datasetName)
	if err != nil {
		return fmt.Errorf("failed to load users: %w", err)
	}

	// Generate SQL from users
	sqlContent := generateUserSQL(users)

	if err := ioutil.WriteFile(outputPath, []byte(sqlContent), 0644); err != nil {
		return fmt.Errorf("failed to write SQL file: %w", err)
	}

	return nil
}

// Helper function to generate SQL (stub - implement as needed)
func generateUserSQL(users []model.UserJSON) string {
	return "-- Generated SQL for users\n"
}
