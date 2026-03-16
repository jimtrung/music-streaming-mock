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

	client := &http.Client{}
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

// SendSignUpRequest sends a user signup request to the API
func SendSignUpRequest(data SignUpRequest) {
	url := apiBase + "/auth/signup"

	jsonData, err := json.Marshal(&data)
	if err != nil {
		fmt.Printf("  âœ— Failed to marshal user data: %v\n", err)
		return
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		fmt.Printf("  âœ— Failed to create signup request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  âœ— Failed to send signup request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  âœ— Error: %s\n", string(body))
		return
	}

	fmt.Printf("  âœ“ User created: %s\n", data.Username)
}

// SendProfileRequest updates a profile with name and avatar (multipart PUT)
func SendProfileRequest(name string, avatarPath string, accessToken string) uuid.UUID {
	url := apiBase + "/auth/profile"

	file, err := os.Open(avatarPath)
	if err != nil {
		fmt.Printf("  âœ— Failed to open avatar %s: %v\n", avatarPath, err)
		return uuid.UUID{}
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("name", name); err != nil {
		fmt.Printf("  âœ— Failed to write name field: %v\n", err)
		return uuid.UUID{}
	}

	part, err := writer.CreateFormFile("avatar", filepath.Base(avatarPath))
	if err != nil {
		fmt.Printf("  âœ— Failed to create avatar form file: %v\n", err)
		return uuid.UUID{}
	}

	if _, err := io.Copy(part, file); err != nil {
		fmt.Printf("  âœ— Failed to copy avatar file: %v\n", err)
		return uuid.UUID{}
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("  âœ— Failed to finalize payload: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		fmt.Printf("  âœ— Failed to create profile request: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  âœ— Failed to send profile request: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("  âœ— Error: %s\n", string(bodyBytes))
		return uuid.UUID{}
	}

	fmt.Printf("  âœ“ Profile updated: %s\n", name)
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
		fmt.Printf("  âœ— Failed to marshal artist data: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("  âœ— Failed to create artist request: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  âœ— Failed to send artist request: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  âœ— Error: %s\n", string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  âœ“ Artist created: %s\n", artistName)
	return uuid.UUID{}
}

// SendTrackRequest sends a track creation request to the API
func SendTrackRequest(track model.TrackJSON, accessToken string) uuid.UUID {
	url := apiBase + "/track"

	payload := map[string]interface{}{
		"title":    track.Title,
		"artistId": track.ArtistId,
		"genre":    track.Genre,
		"duration": track.Duration,
		"audioUrl": track.AudioUrl,
		"coverUrl": track.CoverUrl,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("  âœ— Failed to marshal track data: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("  âœ— Failed to create track request: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  âœ— Failed to send track request: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  âœ— Error: %s\n", string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  âœ“ Track created: %s\n", track.Title)
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
		fmt.Printf("  âœ— Failed to marshal playlist data: %v\n", err)
		return uuid.UUID{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("  âœ— Failed to create playlist request: %v\n", err)
		return uuid.UUID{}
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  âœ— Failed to send playlist request: %v\n", err)
		return uuid.UUID{}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("  âœ— Error: %s\n", string(body))
		return uuid.UUID{}
	}

	fmt.Printf("  âœ“ Playlist created: %s\n", playlist.Name)
	return uuid.UUID{}
}
