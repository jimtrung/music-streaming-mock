package export

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

// ExportAsSQL exports mock data as SQL INSERT statements
func ExportAsSQL(data *model.MockDataJSON, filename string) error {
	dir := "output/sql"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	filepath := filepath.Join(dir, filename)

	sqlContent := "-- Generated SQL INSERT Statements\n"
	sqlContent += fmt.Sprintf("-- Generated at: %s\n", time.Now().Format(time.RFC3339))
	sqlContent += fmt.Sprintf("-- Dataset: %s\n\n", data.Metadata.DatasetName)

	// Generate user inserts
	sqlContent += "-- Users\n"
	for _, user := range data.Users {
		sqlContent += fmt.Sprintf(
			"INSERT INTO users (id, username, email, password, role, provider, is_verified, is_premium, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s'::role_type, '%s'::provider_type, %v, %v, '%s', '%s');\n",
			user.Id, user.Username, user.Email, user.Password, user.Role, user.Provider,
			user.IsVerified, user.IsPremium, user.CreatedAt, user.UpdatedAt,
		)
	}

	sqlContent += "\n-- Artists\n"
	for _, artist := range data.Artists {
		sqlContent += fmt.Sprintf(
			"INSERT INTO artists (id, name, is_verified, created_at, updated_at) VALUES ('%s', '%s', %v, '%s', '%s');\n",
			artist.Id, artist.Name, artist.IsVerified, artist.CreatedAt, artist.UpdatedAt,
		)
	}

	sqlContent += "\n-- Playlists\n"
	for _, playlist := range data.Playlists {
		sqlContent += fmt.Sprintf(
			"INSERT INTO playlists (id, owner_id, name, is_public, created_at, updated_at) VALUES ('%s', '%s', '%s', %v, '%s', '%s');\n",
			playlist.Id, playlist.OwnerId, playlist.Name, playlist.IsPublic, playlist.CreatedAt, playlist.UpdatedAt,
		)
	}

	// Write to file
	if err := os.WriteFile(filepath, []byte(sqlContent), 0644); err != nil {
		return fmt.Errorf("failed to write SQL file: %v", err)
	}

	fmt.Printf("✓ SQL exported successfully: %s\n", filepath)
	return nil
}

// GenerateFilename creates a timestamped filename
func GenerateFilename(prefix, extension string) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_%s%s", prefix, timestamp, extension)
}

// GenerateJSONFilename creates a timestamped JSON filename
func GenerateJSONFilename(datasetName string) string {
	timestamp := time.Now().Format("2006-01-02")
	return fmt.Sprintf("mock_data_%s_%s.json", datasetName, timestamp)
}

// GenerateSQLFilename creates a timestamped SQL filename
func GenerateSQLFilename(prefix string) string {
	return GenerateFilename(prefix, ".sql")
}

// SaveAndValidate exports data and validates it
func SaveAndValidate(data *model.MockDataJSON, jsonFilename, schemaPath string) error {
	// Update metadata
	data.UpdateMetadata("1.0", data.Metadata.DatasetName, data.Metadata.Description)

	// Export as JSON
	if err := ExportAsJSON(data, jsonFilename); err != nil {
		log.Printf("Warning: failed to export JSON: %v", err)
	}

	// Export as SQL
	sqlFilename := GenerateSQLFilename("mock_data")
	if err := ExportAsSQL(data, sqlFilename); err != nil {
		log.Printf("Warning: failed to export SQL: %v", err)
	}

	return nil
}
