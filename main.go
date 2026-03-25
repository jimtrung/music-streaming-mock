package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

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
	fmt.Println("\nNhấn Enter để tiếp tục...")
	fmt.Scanln()
}

func main() {
	// Initialize repository
	var err error
	dataRepo, err = repository.NewDataRepository("./data")
	if err != nil {
		log.Fatalf("Không thể khởi tạo kho dữ liệu: %v\n", err)
	}

	dataMgr = repository.NewDatasetManager(dataRepo)

	decision := -1

	for decision != 0 {
		clearScreen()
		printMainMenu()
		fmt.Print("Lựa chọn của bạn: ")
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
			handleAddTrackToPlaylist()
		case 8:
			handleExportJSON()
		case 9:
			handleListDatasets()
		case 10:
			handleViewDatasetInfo()
		case 11:
			handleCoreExportSQL()
		case 12:
			handleImportJSON()
		case 13:
			handleGenerateFromAPI()
		case 14:
			handleExportToAPI()
		case 0:
			fmt.Println("Exiting...")
		default:
			fmt.Println("Lựa chọn không hợp lệ. Vui lòng thử lại.")
		}
	}
}

func printMainMenu() {
	fmt.Println("\n╔════════════════ TRÌNH SINH DỮ LIỆU GIẢ ══════════════╗")
	if currentCtx != nil {
		stats := currentCtx.GetStats()
		fmt.Printf("║ Tập dữ liệu: %-38s  ║\n", currentCtx.DatasetName)
		fmt.Printf("║ Người dùng: %d | Nghệ sĩ: %d | Bài hát: %d | Danh sách: %d  ║\n",
			stats["users"], stats["artists"], stats["tracks"], stats["playlists"])
		fmt.Println("╟──────────────────────────────────────────────────────╢")
	} else {
		fmt.Println("║ Chưa chọn tập dữ liệu                                ║")
		fmt.Println("╟──────────────────────────────────────────────────────╢")
	}
	fmt.Println("║ [SINH DỮ LIỆU]                                       ║")
	fmt.Println("║  1. Chọn/Tạo tập dữ liệu                             ║")
	fmt.Println("║  2. Sinh người dùng (trong bộ nhớ)                   ║")
	fmt.Println("║  3. Sinh hồ sơ                                       ║")
	fmt.Println("║  4. Sinh nghệ sĩ                                     ║")
	fmt.Println("║  5. Sinh bài hát                                     ║")
	fmt.Println("║  6. Sinh danh sách phát                              ║")
	fmt.Println("║  7. Thêm bài hát vào danh sách phát                  ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║ [KHO LƯU TRỮ & XUẤT]                                 ║")
	fmt.Println("║  8. Xuất tập dữ liệu dưới dạng JSON                  ║")
	fmt.Println("║  9. Liệt kê tất cả các tập dữ liệu                   ║")
	fmt.Println("║ 10. Xem thông tin tập dữ liệu                        ║")
	fmt.Println("║ 11. Xuất truy vấn SQL                                ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║ [API & HAI CHIỀU]                                    ║")
	fmt.Println("║ 12. Nhập tập dữ liệu JSON                            ║")
	fmt.Println("║ 13. Sinh từ API (dữ liệu mới)                        ║")
	fmt.Println("║ 14. Xuất dữ liệu đã tải lên API                      ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  0. Thoát                                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func handleSelectDataset() {
	fmt.Print("Nhập tên tập dữ liệu (tạo mới hoặc chọn hiện có): ")
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
		if err := dataMgr.CreateNewDataset(datasetName, "Dữ liệu sinh đời"); err != nil {
			log.Printf("Lỗi tạo tập dữ liệu: %v\n", err)
			return
		}
	}

	// Initialize generator context
	var err error
	currentCtx, err = generators.NewGeneratorContext(dataRepo, datasetName)
	if err != nil {
		log.Printf("Lỗi khởi tạo generator context: %v\n", err)
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

	fmt.Printf("✓ Đang làm việc với tập dữ liệu: %s\n", datasetName)
	stats := currentCtx.GetStats()
	fmt.Printf("  Đã tải - Người dùng: %d | Nghệ sĩ: %d | Bài hát: %d | Danh sách: %d\n",
		stats["users"], stats["artists"], stats["tracks"], stats["playlists"])
	pauseForContinue()
}

func handleMockUser() {
	if currentCtx == nil {
		fmt.Println("✗ Vui lòng chọn tập dữ liệu trước (lựa chọn 1)")
		pauseForContinue()
		return
	}

	fmt.Print("Nhập số lượng người dùng cần sinh: ")
	var quantity int
	fmt.Scan(&quantity)

	if quantity <= 0 {
		fmt.Println("Invalid quantity. Must be greater than 0.")
		return
	}

	fmt.Println("\n===================> Phương pháp sinh dữ liệu <====================")
	fmt.Println("| 1. Trong bộ nhớ (đang nhanh, dùng kho lưu trữ)     |")
	fmt.Println("| 2. Gọi hàng hóa API (backend thật)              |")
	fmt.Println("| 3. Truy vấn SQL (nhập cơ sở dữ liệu)              |")
	fmt.Println("=========================================================")

	fmt.Print("Lựa chọn của bạn: ")
	var method int
	fmt.Scan(&method)

	switch method {
	case 1:
		fmt.Printf("Sinh %d người dùng (trong bộ nhớ)...\n", quantity)
		userGen := generators.NewUserGenerator(currentCtx)
		if err := userGen.Generate(quantity); err != nil {
			log.Printf("Lỗi sinh người dùng: %v\n", err)
		}
	case 2:
		fmt.Printf("Sinh %d người dùng qua API...\n", quantity)
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

func handleAddTrackToPlaylist() {
	if currentCtx == nil {
		fmt.Println("✗ Please select a dataset first (option 1)")
		return
	}

	if len(currentCtx.Playlists) == 0 {
		fmt.Println("✗ No playlists found. Generate playlists first (option 6)")
		pauseForContinue()
		return
	}

	if len(currentCtx.Tracks) == 0 {
		fmt.Println("✗ No tracks found. Generate tracks first (option 5)")
		pauseForContinue()
		return
	}

	// Display available playlists
	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Println("Available Playlists:")
	for i, playlist := range currentCtx.Playlists {
		fmt.Printf("[%d] %s (%d tracks)\n", i+1, playlist.Name, len(playlist.TrackIds))
	}
	fmt.Println("═══════════════════════════════════════════════════════")

	fmt.Print("Select playlist number: ")
	var playlistIndex int
	fmt.Scan(&playlistIndex)

	if playlistIndex < 1 || playlistIndex > len(currentCtx.Playlists) {
		fmt.Println("Invalid playlist selection.")
		return
	}

	playlistIdx := playlistIndex - 1

	fmt.Print("Enter the number of tracks to add to this playlist: ")
	var trackCount int
	fmt.Scan(&trackCount)

	if trackCount <= 0 {
		fmt.Println("Invalid track count. Must be greater than 0.")
		return
	}

	// Display available tracks
	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Println("Available Tracks:")
	for i, track := range currentCtx.Tracks {
		// Find artist name
		artistName := "Unknown"
		for _, artist := range currentCtx.Artists {
			if artist.Id == track.ArtistId {
				artistName = artist.Name
				break
			}
		}
		fmt.Printf("[%d] %s (%s)\n", i+1, track.Title, artistName)
	}
	fmt.Println("═══════════════════════════════════════════════════════")

	var newPlaylistTracks []model.PlaylistTrackJSON
	addedCount := 0

	for i := 0; i < trackCount; i++ {
		fmt.Printf("\n[Track %d/%d]\n", i+1, trackCount)
		fmt.Print("Select track number (or 0 to skip): ")
		var trackIndex int
		fmt.Scan(&trackIndex)

		if trackIndex == 0 {
			continue
		}

		if trackIndex < 1 || trackIndex > len(currentCtx.Tracks) {
			fmt.Println("Invalid track selection.")
			continue
		}

		track := currentCtx.Tracks[trackIndex-1]

		// Check if track already in playlist
		alreadyExists := false
		for _, tid := range currentCtx.Playlists[playlistIdx].TrackIds {
			if tid == track.Id {
				fmt.Printf("✗ Track already in playlist: %s\n", track.Title)
				alreadyExists = true
				break
			}
		}

		if alreadyExists {
			continue
		}

		// Add track to playlist
		currentCtx.Playlists[playlistIdx].TrackIds = append(currentCtx.Playlists[playlistIdx].TrackIds, track.Id)
		newPlaylistTracks = append(newPlaylistTracks, model.PlaylistTrackJSON{
			PlaylistId: currentCtx.Playlists[playlistIdx].Id,
			TrackId:    track.Id,
			Position:   len(currentCtx.Playlists[playlistIdx].TrackIds),
			AddedAt:    time.Now().UTC().Format(time.RFC3339),
		})

		fmt.Printf("✓ Added: %s\n", track.Title)
		addedCount++
	}

	if addedCount == 0 {
		fmt.Println("No tracks were added.")
		pauseForContinue()
		return
	}

	// Add all new playlist-track relationships to context
	currentCtx.Mutex.Lock()
	currentCtx.PlaylistTracks = append(currentCtx.PlaylistTracks, newPlaylistTracks...)
	currentCtx.Mutex.Unlock()

	fmt.Printf("\n✓ Successfully added %d tracks to playlist: %s\n", addedCount, currentCtx.Playlists[playlistIdx].Name)
	fmt.Println("💡 Remember to export (option 8) to save changes!")
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
		fmt.Println("✗ Vui lòng tải tập dữ liệu trước (nhập JSON với lựa chọn 11 hoặc sinh với lựa chọn 12)")
		pauseForContinue()
		return
	}

	stats := currentCtx.GetStats()
	if stats["users"] == 0 {
		fmt.Println("✗ Không có dữ liệu được tải trong ngữ cảnh hiện tại. Vui lòng nhập hoặc sinh dữ liệu trước.")
		pauseForContinue()
		return
	}

	fmt.Printf("Bắt đầu xuất dữ liệu đã tải lên các điểm cuối API...\n")
	fmt.Printf("  Người dùng: %d | Nghệ sĩ: %d | Bài hát: %d | Danh sách: %d\n",
		stats["users"], stats["artists"], stats["tracks"], stats["playlists"])

	fmt.Print("\nXuất dữ liệu lên API? (có/không): ")
	var confirm string
	fmt.Scan(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("✗ Xuất đã bị huỷ bỏ.")
		pauseForContinue()
		return
	}

	fmt.Println("\n[XUẤT LÊN API]")

	// Get the compiled mock data
	mockData := currentCtx.CompileToMockData()
	userLookup := make(map[string]model.UserJSON)
	for _, user := range mockData.Users {
		userLookup[user.Id] = user
	}

	avatars, avatarErr := api.ListAvatarFiles(filepath.Join("assets", "avatars"))
	if avatarErr != nil {
		fmt.Printf("✗ Vấn đề tải ảnh đại diện: %v\n", avatarErr)
	}

	avatarIdx := 0

	// Send users to API
	fmt.Printf("\n📄 Đang gửi %d người dùng lên API...\n", len(mockData.Users))
	successCount := 0
	for idx, user := range mockData.Users {
		fmt.Printf("  [%d/%d] Tạo người dùng: %s\n", idx+1, len(mockData.Users), user.Username)
		signUpReq := api.SignUpRequest{
			Username: user.Username,
			Email:    api.SanitizeEmail(user.Username), // Use sanitized email from username
			Password: user.PlainPasswd,                 // Use plain password for API
		}

		// Check if PlainPasswd is empty
		if user.PlainPasswd == "" {
			fmt.Printf("  ✗ Bỏ qua (không có mật khẩu)\n")
			continue
		}

		if err := api.SendSignUpRequest(signUpReq); err != nil {
			fmt.Printf("  ✗ Thất bại: %v\n", err)
			continue
		}
		successCount++
	}
	fmt.Printf("  ✓ Successfully created %d/%d users\n", successCount, len(mockData.Users))

	// Build user token lookup
	fmt.Printf("\n📤 Authenticating users...\n")
	userTokens := make(map[string]string)
	authSuccess := 0
	for idx, user := range mockData.Users {
		if user.PlainPasswd == "" {
			continue
		}
		token, err := api.SignIn(user.Username, user.PlainPasswd)
		if err != nil {
			fmt.Printf("  [%d] ✗ %s: %v\n", idx+1, user.Username, err)
			continue
		}
		userTokens[user.Id] = token
		authSuccess++
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("  ✓ Authenticated %d users\n", authSuccess)

	// Send profiles to API
	fmt.Printf("\n📤 Sending %d profiles to API...\n", len(mockData.Profiles))
	for idx, profile := range mockData.Profiles {
		user, ok := userLookup[profile.UserId]
		if !ok {
			fmt.Printf("  [%d/%d] Bỏ qua hồ sơ (᨟ười dùng không tìm thấy): %s\n", idx+1, len(mockData.Profiles), profile.Name)
			continue
		}

		if user.PlainPasswd == "" {
			fmt.Printf("  [%d/%d] Bỏ qua hồ sơ (᨟ười dùng không có mật khẩu): %s\n", idx+1, len(mockData.Profiles), user.Username)
			continue
		}

		if len(avatars) == 0 {
			fmt.Printf("  [%d/%d] Bỏ qua hồ sơ (không có ảnh đại diện): %s\n", idx+1, len(mockData.Profiles), user.Username)
			continue
		}

		token, ok := userTokens[user.Id]
		if !ok {
			fmt.Printf("  [%d/%d] Bỏ qua hồ sơ (không có token xác thực): %s\n", idx+1, len(mockData.Profiles), user.Username)
			continue
		}

		avatarFile := avatars[avatarIdx%len(avatars)]
		avatarIdx++

		displayName := profile.Name
		if displayName == "" {
			displayName = user.Username
		}

		avatarPath := filepath.Join("assets", "avatars", avatarFile)
		fmt.Printf("  [%d/%d] Cập nhật hồ sơ cho %s (ảnh đại diện: %s)\n", idx+1, len(mockData.Profiles), user.Username, avatarFile)
		api.SendProfileRequest(displayName, avatarPath, token)
	}

	// Send artists to API
	fmt.Printf("\n📄 Đang gửi %d nghệ sĩ lên API...\n", len(mockData.Artists))
	artistTokens := make(map[string]string) // Map artist ID to owner's token
	for idx, artist := range mockData.Artists {
		fmt.Printf("  [%d/%d] Tạo nghệ sĩ: %s\n", idx+1, len(mockData.Artists), artist.Name)

		token, ok := userTokens[artist.UserId]
		if !ok {
			fmt.Printf("  ✗ Bỏ qua (không có token xác thực cho người sở hữu)\n")
			continue
		}

		api.SendArtistRequest(artist.Name, token)
		artistTokens[artist.Id] = token
	}

	// Send tracks to API
	fmt.Printf("\n📄 Đang gửi %d bài hát lên API...\n", len(mockData.Tracks))
	trackIdMap := make(map[string]string) // Map old ID -> new API ID
	for idx, track := range mockData.Tracks {
		fmt.Printf("  [%d/%d] Tạo bài hát: %s\n", idx+1, len(mockData.Tracks), track.Title)

		token, ok := artistTokens[track.ArtistId]
		if !ok {
			fmt.Printf("  ✗ Bỏ qua (nghệ sĩ chưa được tạo hoặc không có token xác thực)\n")
			continue
		}

		newTrackId := api.SendTrackRequest(track, token)
		if newTrackId.String() != "00000000-0000-0000-0000-000000000000" {
			trackIdMap[track.Id] = newTrackId.String()
		}
	}

	// Send playlists to API
	fmt.Printf("\n📄 Đang gửi %d danh sách phát lên API...\n", len(mockData.Playlists))
	playlistIdMap := make(map[string]string) // Map old ID -> new API ID
	for idx, playlist := range mockData.Playlists {
		fmt.Printf("  [%d/%d] Tạo danh sách phát: %s\n", idx+1, len(mockData.Playlists), playlist.Name)

		token, ok := userTokens[playlist.OwnerId]
		if !ok {
			fmt.Printf("  ✗ Bỏ qua (không có token xác thực cho người sở hữu)\n")
			continue
		}

		newPlaylistId := api.SendPlaylistRequest(playlist, token)
		if newPlaylistId == "" {
			fmt.Printf("  ✗ Không thể tạo danh sách phát\n")
			continue
		}
		playlistIdMap[playlist.Id] = newPlaylistId
	}

	// Send playlist tracks to API
	fmt.Printf("\n📄 Đang thêm %d bài hát vào danh sách phát...\n", len(mockData.PlaylistTracks))
	playlistLookup := make(map[string]model.PlaylistJSON)
	for _, playlist := range mockData.Playlists {
		playlistLookup[playlist.Id] = playlist
	}

	for idx, pt := range mockData.PlaylistTracks {
		playlist, ok := playlistLookup[pt.PlaylistId]
		if !ok {
			fmt.Printf("  [%d/%d] ✗ Danh sách phát không tìm thấy\n", idx+1, len(mockData.PlaylistTracks))
			continue
		}

		token, ok := userTokens[playlist.OwnerId]
		if !ok {
			fmt.Printf("  [%d/%d] ✗ Không có token xác thực cho danh sách phát: %s\n", idx+1, len(mockData.PlaylistTracks), playlist.Name)
			continue
		}

		// Use the new playlist ID returned from API, or fall back to old ID if not found
		apiPlaylistId := pt.PlaylistId
		if newId, exists := playlistIdMap[pt.PlaylistId]; exists {
			apiPlaylistId = newId
		}

		// Use the new track ID returned from API, or fall back to old ID if not found
		apiTrackId := pt.TrackId
		if newId, exists := trackIdMap[pt.TrackId]; exists {
			apiTrackId = newId
		}

		if err := api.SendAddTrackToPlaylistRequest(apiPlaylistId, apiTrackId, pt.Position, token); err != nil {
			fmt.Printf("Playlist ID: %s | Track ID: %s", apiPlaylistId, apiTrackId)
			fmt.Printf("  [%d/%d] ✗ Lỗi: %v\n", idx+1, len(mockData.PlaylistTracks), err)
			continue
		}
	}
	fmt.Printf("  ✓ Hoàn thành thêm bài hát vào danh sách phát\n")

	fmt.Printf("\n✓ Xuất lên API đã hoàn thành!\n")

	// Save compiled data as JSON for reference
	fmt.Print("Lưu xuất dữ liệu dưới dạng tầp tin JSON? (có/không): ")
	fmt.Scan(&confirm)

	if confirm == "y" || confirm == "Y" {
		exportDir := filepath.Join("output", "api")
		if err := os.MkdirAll(exportDir, 0755); err != nil {
			fmt.Printf("✗ Lỗi tạo thư mục xuất: %v\n", err)
		} else {
			filename := filepath.Join(exportDir, fmt.Sprintf("api_export_%s.json", currentCtx.DatasetName))
			data, _ := json.MarshalIndent(mockData, "", "  ")
			if err := os.WriteFile(filename, data, 0644); err != nil {
				fmt.Printf("✗ Lỗi lưu tầp tin: %v\n", err)
			} else {
				fmt.Printf("✓ Tầp tin xuất được lưu tại: %s\n", filename)
			}
		}
	}
	pauseForContinue()
}
