package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jimtrung/music-streaming-mock-data/generator/api"
	"github.com/jimtrung/music-streaming-mock-data/generator/generators"
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
			handleExportJSON()
		case 8:
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
	fmt.Println("║                                                      ║")
	fmt.Println("║ [KHO LƯU TRỮ & XUẤT]                                 ║")
	fmt.Println("║  7. Xuất tập dữ liệu dưới dạng JSON                  ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║ [API & HAI CHIỀU]                                    ║")
	fmt.Println("║  8. Xuất dữ liệu đã tải lên API                      ║")
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

	fmt.Printf("Sinh %d người dùng (trong bộ nhớ)...\n", quantity)
	userGen := generators.NewUserGenerator(currentCtx)
	if err := userGen.Generate(quantity); err != nil {
		log.Printf("Lỗi sinh người dùng: %v\n", err)
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

	fmt.Print("\nXuất dữ liệu lên API? (y/n): ")
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
	fmt.Print("Lưu xuất dữ liệu dưới dạng tầp tin JSON? (y/n): ")
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
