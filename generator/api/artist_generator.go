package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/util"
)

type Artist struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NOTE: Artist no longer has avatar
// We only send JSON payload to backend API
func GenerateArtist(quantity int) {
	artistNames := util.ReadFileTxt("./assets/artist_names.txt")
	accounts := util.ReadFileCSV("./assets/accounts.csv")
	artistIds := []uuid.UUID{}

	for i := 1; i <= quantity; i++ {
		// Random artist name (remove after pick to avoid duplicates)
		j := rand.Intn(len(artistNames))
		artistName := artistNames[j]
		artistNames = append(artistNames[:j], artistNames[j+1:]...)

		artistName = strings.Trim(artistName, " ")

		// Get access token
		accessToken := signIn(SignInRequest{
			Username: accounts[i-1].Username,
			Password: accounts[i-1].Password,
		})

		// Send create artist request
		artistId := sendRequest(artistName, accessToken)
		artistIds = append(artistIds, artistId)
	}

	util.WriteFileTxt("assets/artist_ids.txt", util.UUIDToString(artistIds))
}

func signIn(request SignInRequest) string {
	url := "http://localhost:8080/auth/signin"

	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Failed to marshal signin request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Fatalf("Failed to create signin request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send signin request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf(
			"Signin failed. Status: %d, Body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode token response: %v", err)
	}

	return result.AccessToken
}

func sendRequest(artistName string, accessToken string) uuid.UUID {
	url := "http://localhost:8080/artist"

	payload := map[string]string{
		"name": artistName,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal artist payload: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Fatalf("Failed to create artist request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send artist request: %v", err)
	}
	defer resp.Body.Close()

	artist := Artist{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&artist); err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	fmt.Println("Create artist:", artist.Name, "| Status:", resp.Status)

	return artist.Id
}
