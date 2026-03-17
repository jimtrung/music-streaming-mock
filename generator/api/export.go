package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/model"
)

const (
	apiBase = "http://localhost:8080"
)

// SignIn authenticates a user and returns an access token
func SignIn(username, password string) (string, error) {
	url := apiBase + "/auth/signin"

	payload, err := json.Marshal(SignInRequest{Username: username, Password: password})
	if err != nil {
		return "", fmt.Errorf("failed to marshal signin payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create signin request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send signin request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("signin failed (%d): %s", resp.StatusCode, string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return token.AccessToken, nil
}

// ListAvatarFiles returns avatar file names from the given directory
func ListAvatarFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read avatars directory: %w", err)
	}

	avatars := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			avatars = append(avatars, entry.Name())
		}
	}

	if len(avatars) == 0 {
		return nil, fmt.Errorf("no avatar files found in %s", dir)
	}

	return avatars, nil
}

// SanitizeEmail removes invalid characters from email to make it valid
func SanitizeEmail(username string) string {
	// Remove any non-alphanumeric characters except dots and hyphens
	// Email local part (before @) can only have alphanumeric, dots, hyphens, underscores
	regex := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := regex.ReplaceAllString(username, "")

	// Remove leading/trailing dots and hyphens
	for len(sanitized) > 0 && (sanitized[0] == '.' || sanitized[0] == '-' || sanitized[0] == '_') {
		sanitized = sanitized[1:]
	}
	for len(sanitized) > 0 && (sanitized[len(sanitized)-1] == '.' || sanitized[len(sanitized)-1] == '-') {
		sanitized = sanitized[:len(sanitized)-1]
	}

	if sanitized == "" {
		sanitized = "user"
	}

	return sanitized + "@example.com"
}

// SendSignUpRequest sends a user signup request to the API
func SendSignUpRequest(data SignUpRequest) error {
	url := apiBase + "/auth/signup"

	jsonData, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu người dùng: %w", err)
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("không thể tạo yêu cầu kỷ kẻ đăng ký: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("không thể gửi yêu cầu đăng ký tới %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)

		// Try to parse JSON error response
		var errResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if errMsg, ok := errResp["error"]; ok {
				return fmt.Errorf("kỳ kẻ đăng ký thất bại (mã %d): %v", resp.StatusCode, errMsg)
			}
		}

		return fmt.Errorf("kỳ kẻ đăng ký thất bại (mã %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SendProfileRequest sends a profile update request with avatar
func SendProfileRequest(name string, avatarPath string, accessToken string) uuid.UUID {
	url := apiBase + "/auth/profile"

	file, err := os.Open(avatarPath)
	if err != nil {
		fmt.Printf("  × Không thể mở ảnh đại diện %s: %v\n", avatarPath, err)
		return uuid.UUID{}
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("name", name); err != nil {
		fmt.Printf("  × Không thể ghi trường họ tên: %v\n", err)
		return uuid.UUID{}
	}

	part, err := writer.CreateFormFile("avatar", filepath.Base(avatarPath))
	if err != nil {
		fmt.Printf("  × Không thể tạo trường biểu mẫu ảnh đại diện: %v\n", err)
		return uuid.UUID{}
	}

	if _, err := io.Copy(part, file); err != nil {
		fmt.Printf("  × Không thể sao chép tệp ảnh đại diện: %v\n", err)
		return uuid.UUID{}
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("  × Không thể hoàn tất tải lên: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		fmt.Printf("  × Không thể tạo yêu cầu hồ sơ: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  × Không thể gửi yêu cầu hồ sơ: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("  × Lỗi: %s\n", string(bodyBytes))
		return uuid.UUID{}
	}

	fmt.Printf("  ✓ Hồ sơ được cập nhật: %s\n", name)
	return uuid.UUID{}
}

// SendArtistRequest sends an artist creation request to the API
func SendArtistRequest(artistName string, accessToken string) uuid.UUID {
	url := apiBase + "/artist"

	payload := map[string]string{
		"name": artistName,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("  × Không thể mã hóa dữ liệu nghệ sĩ: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("  × Không thể tạo yêu cầu nghệ sĩ: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  × Không thể gửi yêu cầu nghệ sĩ: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  × Lỗi: %s\n", string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  ✓ Nghệ sĩ được tạo: %s\n", artistName)
	return uuid.UUID{}
}

// SendTrackRequest sends a track creation request to the API using multipart form-data
func SendTrackRequest(track model.TrackJSON, accessToken string) uuid.UUID {
	url := apiBase + "/track"

	// Create multipart writer
	reader, writer := io.Pipe()
	mw := multipart.NewWriter(writer)

	// Goroutine to write form data
	go func() {
		defer writer.Close()
		defer mw.Close()

		// Add form fields
		mw.WriteField("title", track.Title)
		if track.TrackNum > 0 {
			mw.WriteField("trackNumber", fmt.Sprintf("%d", track.TrackNum))
		}

		// Add audio file (required)
		audioFile, err := os.Open(track.AudioUrl)
		if err != nil {
			fmt.Printf("  × Không thể mở tệp âm thanh: %v\n", err)
			return
		}
		defer audioFile.Close()

		fw, err := mw.CreateFormFile("audio", filepath.Base(track.AudioUrl))
		if err != nil {
			fmt.Printf("  × Không thể tạo trường biểu mẫu âm thanh: %v\n", err)
			return
		}
		io.Copy(fw, audioFile)

		// Add cover file (optional)
		if track.CoverUrl != "" {
			coverFile, err := os.Open(track.CoverUrl)
			if err == nil {
				defer coverFile.Close()
				fw, err := mw.CreateFormFile("cover", filepath.Base(track.CoverUrl))
				if err == nil {
					io.Copy(fw, coverFile)
				}
			}
		}
	}()

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		fmt.Printf("  × Không thể tạo yêu cầu bài hát: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  × Không thể gửi yêu cầu bài hát: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		// Try to parse error response
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			if errMsg, ok := errResp["error"]; ok {
				fmt.Printf("  × Lỗi: %v\n", errMsg)
				return uuid.UUID{}
			}
		}
		fmt.Printf("  × Lỗi HTTP %d: %s\n", resp.StatusCode, string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  ✓ Bài hát được tạo: %s\n", track.Title)
	return uuid.UUID{}
}

// SendPlaylistRequest sends a playlist creation request to the API
func SendPlaylistRequest(playlist model.PlaylistJSON, accessToken string) uuid.UUID {
	url := apiBase + "/playlist"

	payload := map[string]interface{}{
		"name":     playlist.Name,
		"isPublic": playlist.IsPublic,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("  × Không thể mã hóa dữ liệu danh sách phát: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("  × Không thể tạo yêu cầu danh sách phát: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  × Không thể gửi yêu cầu danh sách phát: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  × Lỗi: %s\n", string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  ✓ Danh sách phát được tạo: %s\n", playlist.Name)
	return uuid.UUID{}
}

// SendAddTrackToPlaylistRequest sends a request to add a track to a playlist
func SendAddTrackToPlaylistRequest(playlistId, trackId string, position int, accessToken string) error {
	url := apiBase + "/playlist/add"

	payload := map[string]interface{}{
		"playlistId": playlistId,
		"trackId":    trackId,
		"position":   position,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("không thể mã hóa yêu cầu thêm bài hát: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("không thể tạo yêu cầu thêm bài hát: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("không thể gửi yêu cầu thêm bài hát: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("không thể thêm bài hát vào danh sách phát (mã %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListSongFiles returns all MP3 files from the songs directory
func ListSongFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc thư mục bài hát: %w", err)
	}

	songs := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".mp3") {
			songs = append(songs, filepath.Join(dir, entry.Name()))
		}
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("không tìm thấy tệp MP3 trong %s", dir)
	}

	return songs, nil
}

// ListCoverFiles returns all image files from the covers directory
func ListCoverFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc thư mục bìa: %w", err)
	}

	covers := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				covers = append(covers, filepath.Join(dir, entry.Name()))
			}
		}
	}

	if len(covers) == 0 {
		return nil, fmt.Errorf("không tìm thấy tệp hình ảnh trong %s", dir)
	}

	return covers, nil
}
