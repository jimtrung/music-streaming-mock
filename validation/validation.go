package validation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateAgainstSchema validates mock data against the JSON schema
func ValidateAgainstSchema(dataJSON []byte, schemaPath string) error {
	// Load schema from file
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	// Load the data to validate
	dataLoader := gojsonschema.NewBytesLoader(dataJSON)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return fmt.Errorf("failed to validate against schema: %v", err)
	}

	// Check results
	if !result.Valid() {
		errMsg := "Validation failed with errors:\n"
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("  - %s\n", desc.String())
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}

// ValidateJSONFile validates a JSON file against the schema
func ValidateJSONFile(jsonPath, schemaPath string) error {
	// Read the JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Validate
	if err := ValidateAgainstSchema(jsonData, schemaPath); err != nil {
		return err
	}

	fmt.Printf("✓ JSON file '%s' is valid against schema\n", jsonPath)
	return nil
}

// ValidateDataStructure validates the structure of mock data
func ValidateDataStructure(data map[string]interface{}) error {
	// Check required fields
	requiredFields := []string{"metadata", "users", "artists", "tracks", "playlists"}
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate metadata
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		requiredMetadata := []string{"version", "generatedAt", "datasetName"}
		for _, field := range requiredMetadata {
			if _, exists := metadata[field]; !exists {
				return fmt.Errorf("missing required metadata field: %s", field)
			}
		}
	}

	return nil
}

// PrintValidationSummary prints a summary of validation results
func PrintValidationSummary(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var mockData map[string]interface{}
	if err := json.Unmarshal(data, &mockData); err != nil {
		return err
	}

	if meta, ok := mockData["metadata"].(map[string]interface{}); ok {
		fmt.Println("\n=== Mock Data Summary ===")
		if v, exists := meta["version"]; exists {
			fmt.Printf("Version: %v\n", v)
		}
		if name, exists := meta["datasetName"]; exists {
			fmt.Printf("Dataset: %v\n", name)
		}
		if date, exists := meta["generatedAt"]; exists {
			fmt.Printf("Generated: %v\n", date)
		}
		if desc, exists := meta["description"]; exists {
			fmt.Printf("Description: %v\n", desc)
		}
	}

	if users, ok := mockData["users"].([]interface{}); ok {
		fmt.Printf("Users: %d\n", len(users))
	}
	if artists, ok := mockData["artists"].([]interface{}); ok {
		fmt.Printf("Artists: %d\n", len(artists))
	}
	if tracks, ok := mockData["tracks"].([]interface{}); ok {
		fmt.Printf("Tracks: %d\n", len(tracks))
	}
	if playlists, ok := mockData["playlists"].([]interface{}); ok {
		fmt.Printf("Playlists: %d\n", len(playlists))
	}

	return nil
}
