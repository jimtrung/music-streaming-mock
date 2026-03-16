package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jimtrung/music-streaming-mock-data/generator/api"
	"github.com/jimtrung/music-streaming-mock-data/generator/generators"
	"github.com/jimtrung/music-streaming-mock-data/generator/query"
	"github.com/jimtrung/music-streaming-mock-data/generator/repository"
	"github.com/jimtrung/music-streaming-mock-data/model"
)

var (
	dataRepo   *repository.DataRepository
	dataMgr    *repository.DatasetManager
	currentCtx *generators.GeneratorContext
)

// clearScreen clears the terminal screen
func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// pauseForContinue waits for the user to press Enter
func pauseForContinue() {
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
}

func main() {
	// Initialize repository
	var err error
	dataRepo, err = repository.NewDataRepository("./data")
	if err != nil {
		log.Fatalf("Failed to initialize data repository: %v\n", err)
	}

	dataMgr = repository.NewDatasetManager(dataRepo)

	decision := -1

	for decision != 0 {
		clearScreen()
		printMainMenu()
		fmt.Print("Your decision: ")
		fmt.Scan(&decision)

		if decision == 0 {
			fmt.Println("Exiting...")
			break
		}

		switch decision {
		case 1:
			handleSelectDataset()
		case 2:
			handleMockUser()
		case 3:
			handleMockProfile()
		case 4:
			handleMockArtist()
		case 5:
			handleMockTrack()
		case 6:
			handleMockPlaylist()
		case 7:
			handleExportJSON()
		case 8:
			handleListDatasets()
		case 9:
			handleViewDatasetInfo()
		case 10:
			handleCoreExportSQL()
		case 11:
			handleImportJSON()
		case 12:
			handleGenerateFromAPI()
		case 13:
			handleExportToAPI()
		case 0:
			fmt.Println("Exiting...")
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func printMainMenu() {
	fmt.Println("\n╔════════════════ MOCK DATA GENERATOR ═════════════════╗")
	if currentCtx != nil {
		stats := currentCtx.GetStats()
		fmt.Printf("║ Dataset: %-42s  ║\n", currentCtx.DatasetName)
		fmt.Printf("║ Users: %d | Artists: %d | Tracks: %d | Playlists: %d  ║\n",
			stats["users"], stats["artists"], stats["tracks"], stats["playlists"])
		fmt.Println("╟──────────────────────────────────────────────────────╢")
	} else {
		fmt.Println("║ No dataset selected                                  ║")
		fmt.Println("╟──────────────────────────────────────────────────────╢")
	}
	fmt.Println("║ [DATA GENERATION]                                    ║")
	fmt.Println("║  1. Select/Create dataset                            ║")
	fmt.Println("║  2. Generate users (in-memory)                       ║")
	fmt.Println("║  3. Generate profiles                                ║")
	fmt.Println("║  4. Generate artists                                 ║")
	fmt.Println("║  5. Generate tracks                                  ║")
	fmt.Println("║  6. Generate playlists                               ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║ [REPOSITORY & EXPORT]                                ║")
	fmt.Println("║  7. Export dataset as JSON                           ║")
	fmt.Println("║  8. List all datasets                                ║")
	fmt.Println("║  9. View dataset info                                ║")
	fmt.Println("║ 10. Export SQL queries                               ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║ [API & BIDIRECTIONAL]                                ║")
	fmt.Println("║ 11. Import JSON dataset                              ║")
	fmt.Println("║ 12. Generate from API (new data)                     ║")
	fmt.Println("║ 13. Export loaded data to API                        ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  0. Exit                                             ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func handleSelectDataset() {
	fmt.Print("Enter dataset name (create new or select existing): ")
	var datasetName string
	fmt.Scan(&datasetName)

	if datasetName == "" {
		fmt.Println("Dataset name cannot be empty.")
		return
	}

	// Check if dataset exists
	_, exists := dataRepo.GetDatasetMetadata(datasetName)
	if !exists {
		// Create new dataset
		if err := dataMgr.CreateNewDataset(datasetName, "Generated mock data"); err != nil {
			log.Printf("Error creating dataset: %v\n", err)
			return
		}
	}

	// Initialize generator context
	var err error
	currentCtx, err = generators.NewGeneratorContext(dataRepo, datasetName)
	if err != nil {
		log.Printf("Error initializing generator context: %v\n", err)
		return
	}

	// Load existing data from repository if dataset exists
	if exists {
		allData, err := dataRepo.LoadAllData(datasetName)
		if err == nil && allData != nil {
			currentCtx.Mutex.Lock()
			currentCtx.Users = allData.Users
			currentCtx.Profiles = allData.Profiles
			currentCtx.Artists = allData.Artists
			currentCtx.Tracks = allData.Tracks
			currentCtx.Playlists = allData.Playlists
			currentCtx.PlaylistTracks = allData.PlaylistTracks
			currentCtx.Mutex.Unlock()
		}
	}

	fmt.Printf("✓ Working with dataset: %s\n", datasetName)
	stats := currentCtx.GetStats()
	fmt.Printf("  Loaded - Users: %d | Artists: %d | Tracks: %d | Playlists: %d\n",
		stats["users"], stats["artists"], stats["tracks"], stats["playlists"])
	pauseForContinue()
}

func handleMockUser() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		pauseForContinue()
		return
	}

	fmt.Print("Enter the number of users to generate: ")
	var quantity int
	fmt.Scan(&quantity)

	if quantity <= 0 {
		fmt.Println("Invalid quantity. Must be greater than 0.")
		return
	}

	fmt.Println("\n===================> Generation Method <====================")
	fmt.Println("| 1. In-memory (fast, uses repository storage)   |")
	fmt.Println("| 2. API calls (actual backend)                 |")
	fmt.Println("| 3. SQL queries (database import)              |")
	fmt.Println("=========================================================")

	fmt.Print("Your choice: ")
	var method int
	fmt.Scan(&method)

	switch method {
	case 1:
		fmt.Printf("Generating %d users (in-memory)...\n", quantity)
		userGen := generators.NewUserGenerator(currentCtx)
		if err := userGen.Generate(quantity); err != nil {
			log.Printf("Error generating users: %v\n", err)
		}
	case 2:
		fmt.Printf("Generating %d users via API...\n", quantity)
		api.GenerateUser(quantity)
	case 3:
		fmt.Printf("Generating %d users as SQL...\n", quantity)
		query.GenerateUser(quantity)
	default:
		fmt.Println("Invalid method choice.")
	}
	pauseForContinue()
}

func handleMockProfile() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	fmt.Println("Generating profiles...")
	profileGen := generators.NewProfileGenerator(currentCtx)
	if err := profileGen.Generate(); err != nil {
		log.Printf("Error generating profiles: %v\n", err)
		return
	}

	fmt.Println("✓ Profiles generated successfully")
	pauseForContinue()
}

func handleMockArtist() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	fmt.Print("Enter the number of artists to generate: ")
	var quantity int
	fmt.Scan(&quantity)

	if quantity <= 0 {
		fmt.Println("Invalid quantity. Must be greater than 0.")
		return
	}

	fmt.Printf("Generating %d artists...\n", quantity)
	artistGen := generators.NewArtistGenerator(currentCtx)
	if err := artistGen.Generate(quantity); err != nil {
		log.Printf("Error generating artists: %v\n", err)
	}
	pauseForContinue()
}

func handleMockTrack() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	fmt.Print("Enter the number of tracks to generate: ")
	var quantity int
	fmt.Scan(&quantity)

	if quantity <= 0 {
		fmt.Println("Invalid quantity. Must be greater than 0.")
		return
	}

	fmt.Printf("Generating %d tracks...\n", quantity)
	trackGen := generators.NewTrackGenerator(currentCtx)
	if err := trackGen.Generate(quantity); err != nil {
		log.Printf("Error generating tracks: %v\n", err)
	}
	pauseForContinue()
}

func handleMockPlaylist() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	fmt.Print("Enter the number of playlists to generate: ")
	var quantity int
	fmt.Scan(&quantity)

	if quantity <= 0 {
		fmt.Println("Invalid quantity. Must be greater than 0.")
		return
	}

	fmt.Print("Enter the average number of tracks per playlist: ")
	var tracksPerPlaylist int
	fmt.Scan(&tracksPerPlaylist)

	if tracksPerPlaylist < 0 {
		fmt.Println("Invalid track count. Must be 0 or greater.")
		return
	}

	fmt.Printf("Generating %d playlists with ~%d tracks each...\n", quantity, tracksPerPlaylist)
	playlistGen := generators.NewPlaylistGenerator(currentCtx)
	if err := playlistGen.Generate(quantity, tracksPerPlaylist); err != nil {
		log.Printf("Error generating playlists: %v\n", err)
	}
	pauseForContinue()
}

func handleExportJSON() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		pauseForContinue()
		return
	}
	mockData := currentCtx.CompileToMockData()

	// Update metadata
	mockData.Metadata.Version = "1.0"
	mockData.Metadata.DatasetName = currentCtx.DatasetName
	mockData.Metadata.GeneratedAt = currentCtx.CompileToMockData().Metadata.GeneratedAt
	stats := currentCtx.GetStats()
	mockData.Metadata.TotalUsers = stats["users"]
	mockData.Metadata.TotalArtists = stats["artists"]
	mockData.Metadata.TotalTracks = stats["tracks"]
	mockData.Metadata.TotalPlaylists = stats["playlists"]

	// Save to repository
	if err := dataRepo.SaveAllData(currentCtx.DatasetName, mockData); err != nil {
		log.Printf("Error saving dataset: %v\n", err)
		return
	}

	fmt.Printf("✓ Dataset exported and saved to repository\n")
	fmt.Printf("  Path: %s\n", dataRepo.GetDataPath(currentCtx.DatasetName))
	pauseForContinue()
}

func handleListDatasets() {
	fmt.Println("\n")
	dataMgr.ListAllDatasets()
	pauseForContinue()
}

func handleViewDatasetInfo() {
	fmt.Print("Enter dataset name: ")
	var datasetName string
	fmt.Scan(&datasetName)

	if datasetName == "" {
		fmt.Println("Dataset name cannot be empty.")
		return
	}

	if err := dataMgr.PrintDatasetInfo(datasetName); err != nil {
		log.Printf("Error: %v\n", err)
	}
	pauseForContinue()
}

func handleCoreExportSQL() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	fmt.Print("Enter output SQL file name: ")
	var filename string
	fmt.Scan(&filename)

	if filename == "" {
		filename = "export.sql"
	}

	if err := dataRepo.ExportSQLQueries(currentCtx.DatasetName, filename); err != nil {
		log.Printf("Error exporting SQL: %v\n", err)
		return
	}

	fmt.Printf("✓ SQL exported to %s\n", filename)
	pauseForContinue()
}

func handleImportJSON() {
	fmt.Print("Enter the path to the JSON file to import: ")
	var filePath string
	fmt.Scan(&filePath)

	if filePath == "" {
		fmt.Println("✗ File path cannot be empty.")
		pauseForContinue()
		return
	}

	// Read the JSON file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("✗ Error reading file: %v\n", err)
		return
	}

	// Parse JSON
	var mockData model.MockDataJSON
	if err := json.Unmarshal(fileContent, &mockData); err != nil {
		fmt.Printf("✗ Error parsing JSON: %v\n", err)
		return
	}

	// Create or select dataset
	datasetName := mockData.Metadata.DatasetName
	if datasetName == "" {
		fmt.Print("Enter dataset name for imported data: ")
		fmt.Scan(&datasetName)
	}

	if datasetName == "" {
		fmt.Println("✗ Dataset name cannot be empty.")
		return
	}

	// Check if dataset exists
	_, exists := dataRepo.GetDatasetMetadata(datasetName)
	if !exists {
		if err := dataMgr.CreateNewDataset(datasetName, "Imported from JSON"); err != nil {
			log.Printf("✗ Error creating dataset: %v\n", err)
			return
		}
	}

	// Initialize generator context
	currentCtx, err = generators.NewGeneratorContext(dataRepo, datasetName)
	if err != nil {
		log.Printf("✗ Error initializing generator context: %v\n", err)
		return
	}

	// Load data into context
	currentCtx.Mutex.Lock()
	currentCtx.Users = mockData.Users
	currentCtx.Profiles = mockData.Profiles
	currentCtx.Artists = mockData.Artists
	currentCtx.Tracks = mockData.Tracks
	currentCtx.Playlists = mockData.Playlists
	currentCtx.PlaylistTracks = mockData.PlaylistTracks
	currentCtx.Mutex.Unlock()

	// Save to repository
	if err := dataRepo.SaveAllData(datasetName, &mockData); err != nil {
		log.Printf("✗ Error saving dataset: %v\n", err)
		return
	}

	stats := currentCtx.GetStats()
	fmt.Printf("\n✓ JSON imported successfully!\n")
	fmt.Printf("  Dataset: %s\n", datasetName)
	fmt.Printf("  Users: %d | Artists: %d | Tracks: %d | Playlists: %d\n",
		stats["users"], stats["artists"], stats["tracks"], stats["playlists"])
	pauseForContinue()
}

func handleGenerateFromAPI() {
	fmt.Print("Enter dataset name for API generation: ")
	var datasetName string
	fmt.Scan(&datasetName)

	if datasetName == "" {
		fmt.Println("✗ Dataset name cannot be empty.")
		return
	}

	// Check if dataset exists
	_, exists := dataRepo.GetDatasetMetadata(datasetName)
	if !exists {
		if err := dataMgr.CreateNewDataset(datasetName, "Generated from API"); err != nil {
			log.Printf("✗ Error creating dataset: %v\n", err)
			return
		}
	}

	// Initialize generator context
	var err error
	currentCtx, err = generators.NewGeneratorContext(dataRepo, datasetName)
	if err != nil {
		log.Printf("✗ Error initializing generator context: %v\n", err)
		return
	}

	fmt.Print("Enter the number of users to generate via API: ")
	var userCount int
	fmt.Scan(&userCount)

	if userCount <= 0 {
		fmt.Println("✗ User count must be greater than 0.")
		return
	}

	fmt.Printf("\n[GENERATING FROM API - %d USERS]\n", userCount)
	fmt.Println("Note: API generation may take longer due to backend processing.")

	// Call API to generate
	api.GenerateUser(userCount)

	fmt.Printf("\n✓ API generation completed!\n")
	fmt.Printf("  Dataset: %s\n", datasetName)
	fmt.Println("  (API-generated data has been sent to the backend)")
	fmt.Println("  Tip: Use 'Export to API' (option 13) to push loaded data to endpoints.")
	pauseForContinue()
}

func handleExportToAPI() {
	if currentCtx == nil {
		fmt.Println("✗ Please load a dataset first (import JSON with option 11 or generate with option 12)")
		pauseForContinue()
		return
	}

	stats := currentCtx.GetStats()
	if stats["users"] == 0 {
		fmt.Println("✗ No data loaded in current context. Please import or generate data first.")
		pauseForContinue()
		return
	}

	fmt.Printf("Starting export of loaded data to API endpoints...\n")
	fmt.Printf("  Users: %d | Artists: %d | Tracks: %d | Playlists: %d\n",
		stats["users"], stats["artists"], stats["tracks"], stats["playlists"])

	fmt.Print("\nExport data to API? (y/n): ")
	var confirm string
	fmt.Scan(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("✗ Export cancelled.")
		pauseForContinue()
		return
	}

	fmt.Println("\n[EXPORTING TO API]")

	// Get the compiled mock data
	mockData := currentCtx.CompileToMockData()
	userLookup := make(map[string]model.UserJSON)
	for _, user := range mockData.Users {
		userLookup[user.Id] = user
	}

	avatars, avatarErr := api.ListAvatarFiles(filepath.Join("assets", "avatars"))
	if avatarErr != nil {
		fmt.Printf("✗ Avatar loading issue: %v\n", avatarErr)
	}
	avatarIdx := 0

	// Send users to API
	fmt.Printf("\n📤 Sending %d users to API...\n", len(mockData.Users))
	for idx, user := range mockData.Users {
		fmt.Printf("  [%d/%d] Creating user: %s\n", idx+1, len(mockData.Users), user.Username)
		signUpReq := api.SignUpRequest{
			Username: user.Username,
			Email:    user.Email,
			Password: user.PlainPasswd, // Use plain password for API
		}
		api.SendSignUpRequest(signUpReq)
	}

	// Send profiles to API
	fmt.Printf("\n📤 Sending %d profiles to API...\n", len(mockData.Profiles))
	for idx, profile := range mockData.Profiles {
		user, ok := userLookup[profile.UserId]
		if !ok {
			fmt.Printf("  [%d/%d] Skipped profile (user missing): %s\n", idx+1, len(mockData.Profiles), profile.Name)
			continue
		}

		if user.PlainPasswd == "" {
			fmt.Printf("  [%d/%d] Skipped profile (missing password): %s\n", idx+1, len(mockData.Profiles), user.Username)
			continue
		}

		if len(avatars) == 0 {
			fmt.Printf("  [%d/%d] Skipped profile (no avatars available): %s\n", idx+1, len(mockData.Profiles), user.Username)
			continue
		}

		token, err := api.SignIn(user.Username, user.PlainPasswd)
		if err != nil {
			fmt.Printf("  [%d/%d] Skipped profile (signin failed): %v\n", idx+1, len(mockData.Profiles), err)
			continue
		}

		avatarFile := avatars[avatarIdx%len(avatars)]
		avatarIdx++

		displayName := profile.Name
		if displayName == "" {
			displayName = user.Username
		}

		avatarPath := filepath.Join("assets", "avatars", avatarFile)
		fmt.Printf("  [%d/%d] Updating profile for %s (avatar: %s)\n", idx+1, len(mockData.Profiles), user.Username, avatarFile)
		api.SendProfileRequest(displayName, avatarPath, token)
	}

	// Send artists to API
	fmt.Printf("\n📤 Sending %d artists to API...\n", len(mockData.Artists))
	for idx, artist := range mockData.Artists {
		fmt.Printf("  [%d/%d] Creating artist: %s\n", idx+1, len(mockData.Artists), artist.Name)
		api.SendArtistRequest(artist.Name, "")
	}

	// Send tracks to API
	fmt.Printf("\n📤 Sending %d tracks to API...\n", len(mockData.Tracks))
	for idx, track := range mockData.Tracks {
		fmt.Printf("  [%d/%d] Creating track: %s\n", idx+1, len(mockData.Tracks), track.Title)
		api.SendTrackRequest(track, "")
	}

	// Send playlists to API
	fmt.Printf("\n📤 Sending %d playlists to API...\n", len(mockData.Playlists))
	for idx, playlist := range mockData.Playlists {
		fmt.Printf("  [%d/%d] Creating playlist: %s\n", idx+1, len(mockData.Playlists), playlist.Name)
		api.SendPlaylistRequest(playlist, "")
	}

	fmt.Printf("\n✓ Export to API completed!\n")

	// Save compiled data as JSON for reference
	fmt.Print("Save export as JSON file? (y/n): ")
	fmt.Scan(&confirm)

	if confirm == "y" || confirm == "Y" {
		exportDir := filepath.Join("output", "api")
		if err := os.MkdirAll(exportDir, 0755); err != nil {
			fmt.Printf("✗ Error ensuring export directory: %v\n", err)
		} else {
			filename := filepath.Join(exportDir, fmt.Sprintf("api_export_%s.json", currentCtx.DatasetName))
			data, _ := json.MarshalIndent(mockData, "", "  ")
			if err := os.WriteFile(filename, data, 0644); err != nil {
				fmt.Printf("✗ Error saving file: %v\n", err)
			} else {
				fmt.Printf("✓ Export JSON saved to: %s\n", filename)
			}
		}
	}
	pauseForContinue()
}
