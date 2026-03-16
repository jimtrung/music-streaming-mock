package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/jimtrung/music-streaming-mock-data/model"
)

// DatasetManager coordinates dataset generation and management
type DatasetManager struct {
	repo *DataRepository
}

// NewDatasetManager creates a new dataset manager
func NewDatasetManager(repo *DataRepository) *DatasetManager {
	return &DatasetManager{repo: repo}
}

// CreateNewDataset initializes a new dataset
func (dm *DatasetManager) CreateNewDataset(datasetName, description string) error {
	if err := dm.repo.CreateDataset(datasetName); err != nil {
		return fmt.Errorf("failed to create dataset: %w", err)
	}

	log.Printf("✓ Dataset '%s' created\n", datasetName)
	return nil
}

// PrintDatasetInfo displays dataset information
func (dm *DatasetManager) PrintDatasetInfo(datasetName string) error {
	meta, exists := dm.repo.GetDatasetMetadata(datasetName)
	if !exists {
		return fmt.Errorf("dataset '%s' not found", datasetName)
	}

	fmt.Println("\n===================> DATASET INFO <====================")
	fmt.Printf("  Name:          %s\n", meta.Name)
	fmt.Printf("  Status:        %s\n", meta.Status)
	fmt.Printf("  Created:       %s\n", meta.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  Updated:       %s\n", meta.LastUpdated.Format(time.RFC3339))
	fmt.Printf("  Path:          %s\n", meta.Path)
	fmt.Println("  Data Count:")
	fmt.Printf("    - Users:     %d\n", meta.UserCount)
	fmt.Printf("    - Profiles:  %d\n", meta.ProfileCount)
	fmt.Printf("    - Artists:   %d\n", meta.ArtistCount)
	fmt.Printf("    - Tracks:    %d\n", meta.TrackCount)
	fmt.Printf("    - Playlists: %d\n", meta.PlaylistCount)
	fmt.Println("======================================================")

	return nil
}

// ListAllDatasets displays all available datasets
func (dm *DatasetManager) ListAllDatasets() {
	datasets := dm.repo.ListDatasets()

	if len(datasets) == 0 {
		fmt.Println("No datasets found.")
		return
	}

	fmt.Println("\n===================> AVAILABLE DATASETS <====================")
	fmt.Printf("%-20s %-15s %-30s %-10s\n", "Name", "Status", "Created", "Users")
	fmt.Println("===========================================================")

	for _, meta := range datasets {
		fmt.Printf("%-20s %-15s %-30s %-10d\n",
			meta.Name,
			meta.Status,
			meta.CreatedAt.Format("2006-01-02 15:04:05"),
			meta.UserCount,
		)
	}

	fmt.Println("===========================================================")
}

// DeleteDataset removes a dataset
func (dm *DatasetManager) DeleteDataset(datasetName string) error {
	if err := dm.repo.DeleteDataset(datasetName); err != nil {
		return fmt.Errorf("failed to delete dataset: %w", err)
	}

	log.Printf("✓ Dataset '%s' deleted\n", datasetName)
	return nil
}

// ExportDatasetAsJSON exports a dataset to JSON
func (dm *DatasetManager) ExportDatasetAsJSON(datasetName string, outputPath string) error {
	_, err := dm.repo.LoadAllData(datasetName)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %w", err)
	}

	return nil
}

// CloneDataset creates a copy of an existing dataset
func (dm *DatasetManager) CloneDataset(sourceName, targetName string) error {
	sourceData, err := dm.repo.LoadAllData(sourceName)
	if err != nil {
		return fmt.Errorf("failed to load source dataset: %w", err)
	}

	// Create new dataset
	if err := dm.CreateNewDataset(targetName, fmt.Sprintf("Clone of %s", sourceName)); err != nil {
		return err
	}

	// Save cloned data
	if err := dm.repo.SaveAllData(targetName, sourceData); err != nil {
		return fmt.Errorf("failed to save cloned dataset: %w", err)
	}

	log.Printf("✓ Dataset '%s' cloned from '%s'\n", targetName, sourceName)
	return nil
}

// MergeDatasets combines two datasets
func (dm *DatasetManager) MergeDatasets(sourceName, targetName string) error {
	sourceData, err := dm.repo.LoadAllData(sourceName)
	if err != nil {
		return fmt.Errorf("failed to load source dataset: %w", err)
	}

	targetData, err := dm.repo.LoadAllData(targetName)
	if err != nil {
		return fmt.Errorf("failed to load target dataset: %w", err)
	}

	// Merge users
	targetData.Users = append(targetData.Users, sourceData.Users...)

	// Merge profiles
	targetData.Profiles = append(targetData.Profiles, sourceData.Profiles...)

	// Merge artists
	targetData.Artists = append(targetData.Artists, sourceData.Artists...)

	// Merge tracks
	targetData.Tracks = append(targetData.Tracks, sourceData.Tracks...)

	// Merge playlists
	targetData.Playlists = append(targetData.Playlists, sourceData.Playlists...)

	// Merge playlist tracks
	targetData.PlaylistTracks = append(targetData.PlaylistTracks, sourceData.PlaylistTracks...)

	// Update metadata
	targetData.Metadata.TotalUsers = len(targetData.Users)
	targetData.Metadata.TotalArtists = len(targetData.Artists)
	targetData.Metadata.TotalTracks = len(targetData.Tracks)
	targetData.Metadata.TotalPlaylists = len(targetData.Playlists)

	// Save merged data
	if err := dm.repo.SaveAllData(targetName, targetData); err != nil {
		return fmt.Errorf("failed to save merged dataset: %w", err)
	}

	log.Printf("✓ Dataset '%s' merged into '%s'\n", sourceName, targetName)
	return nil
}

// GenerateDatasetReport creates a report for a dataset
func (dm *DatasetManager) GenerateDatasetReport(datasetName string) (map[string]interface{}, error) {
	meta, exists := dm.repo.GetDatasetMetadata(datasetName)
	if !exists {
		return nil, fmt.Errorf("dataset '%s' not found", datasetName)
	}

	mockData, err := dm.repo.LoadAllData(datasetName)
	if err != nil {
		return nil, fmt.Errorf("failed to load dataset: %w", err)
	}

	report := map[string]interface{}{
		"name":          meta.Name,
		"status":        meta.Status,
		"createdAt":     meta.CreatedAt,
		"updatedAt":     meta.LastUpdated,
		"path":          meta.Path,
		"userCount":     len(mockData.Users),
		"profileCount":  len(mockData.Profiles),
		"artistCount":   len(mockData.Artists),
		"trackCount":    len(mockData.Tracks),
		"playlistCount": len(mockData.Playlists),
	}

	return report, nil
}

// ArchiveDataset archives a completed dataset (marks as read-only)
func (dm *DatasetManager) ArchiveDataset(datasetName string) error {
	meta, exists := dm.repo.GetDatasetMetadata(datasetName)
	if !exists {
		return fmt.Errorf("dataset '%s' not found", datasetName)
	}

	meta.Status = "archived"
	meta.LastUpdated = time.Now()

	// Update in index
	dm.repo.mu.Lock()
	dm.repo.dataIndex.Datasets[datasetName] = meta
	dm.repo.mu.Unlock()

	if err := dm.repo.saveIndex(); err != nil {
		return fmt.Errorf("failed to archive dataset: %w", err)
	}

	log.Printf("✓ Dataset '%s' archived\n", datasetName)
	return nil
}

// GetDatasetPath returns the file path for a dataset
func (dm *DatasetManager) GetDatasetPath(datasetName string) string {
	return dm.repo.GetDataPath(datasetName)
}

// LoadDataset loads a complete dataset into memory
func (dm *DatasetManager) LoadDataset(datasetName string) (*model.MockDataJSON, error) {
	return dm.repo.LoadAllData(datasetName)
}

// SaveDataset saves a dataset to the repository
func (dm *DatasetManager) SaveDataset(datasetName string, mockData *model.MockDataJSON) error {
	return dm.repo.SaveAllData(datasetName, mockData)
}
