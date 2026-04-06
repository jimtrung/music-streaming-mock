package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jimtrung/music-streaming-mock-data/model"
)

// ExportAsJSON exports mock data as a JSON file
func ExportAsJSON(data *model.MockDataJSON, filename string) error {
	// Ensure output directory exists
	dir := "output/json"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	filepath := filepath.Join(dir, filename)

	// Marshal to JSON with pretty printing
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	fmt.Printf("✓ JSON exported successfully: %s\n", filepath)
	fmt.Printf("  - Users: %d | Artists: %d | Tracks: %d | Playlists: %d\n",
		len(data.Users), len(data.Artists), len(data.Tracks), len(data.Playlists))

	return nil
}
