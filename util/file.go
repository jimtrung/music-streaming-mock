package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/model"
)

// WriteFileTxt writes text content to a file
func WriteFileTxt(name string, content string) error {
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", name, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// ReadFileTxt reads lines from a text file
func ReadFileTxt(name string) []string {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", name, err)
	}
	defer file.Close()

	content := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	return content
}

// WriteFileCSV writes CSV data to a file (legacy)
func WriteFileCSV(name string, in []model.Account) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", name, err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&in, file); err != nil {
		log.Fatalf("Failed to write to file %s: %v", name, err)
	}
}

// ReadFileCSV reads CSV data from a file (legacy)
func ReadFileCSV(name string) []model.Account {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", name, err)
	}
	defer file.Close()

	out := []model.Account{}
	if err := gocsv.UnmarshalFile(file, &out); err != nil {
		log.Fatalf("Failed to write data to file: %v", err)
	}

	return out
}

// WriteFileJSON writes mock data as JSON to a file
func WriteFileJSON(name string, data *model.MockDataJSON) error {
	// Ensure directory exists
	dir := "output/json"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filepath := dir + "/" + name

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filepath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// ReadFileJSON reads mock data from a JSON file
func ReadFileJSON(name string) (*model.MockDataJSON, error) {
	jsonBytes, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var data model.MockDataJSON
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}

// UUIDToString converts UUID slice to newline-separated string
func UUIDToString(slice []uuid.UUID) string {
	var content string

	for _, element := range slice {
		content += fmt.Sprintf("%s\n", element.String())
	}

	return content
}
