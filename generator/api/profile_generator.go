package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jimtrung/music-streaming-mock-data/util"
)

// Mock profile data:
// - name: required
// - avatar: required
// - auth token: required
func GenerateProfile() {
	avatars := fetchAvatarNames()
	accounts := util.ReadFileCSV("./assets/accounts.csv")

	for _, account := range accounts {
		avatar := avatars[rand.Intn(len(avatars))]

		accessToken := signIn(SignInRequest{
			Username: account.Username,
			Password: account.Password,
		})

		sendProfileRequest(account.Username, avatar, accessToken)
	}
}

func fetchAvatarNames() []string {
	files, err := os.ReadDir("./assets/avatars")
	if err != nil {
		log.Fatalf("Failed to read avatars directory: %v", err)
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, f.Name())
		}
	}

	return names
}

func sendProfileRequest(name string, avatarPath string, accessToken string) {
	url := "http://127.0.0.1:5255/auth/profile"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// required name
	if err := writer.WriteField("name", name); err != nil {
		log.Fatalf("Failed to write name field: %v", err)
	}

	// required avatar
	file, err := os.Open("./assets/avatars/" + avatarPath)
	if err != nil {
		log.Fatalf("Failed to open avatar file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("avatar", filepath.Base(avatarPath))
	if err != nil {
		log.Fatalf("Failed to create avatar form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Fatalf("Failed to copy avatar file: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Fatalf("Failed to create profile request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send profile request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Create profile:", name, "| Status:", resp.Status)
}
