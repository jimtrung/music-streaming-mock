package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/util"
)

type Playlist struct {
	Id       uuid.UUID `json:"id"`
	OwnerId  uuid.UUID `json:"owner_id"`
	Name     string    `json:"name"`
	IsPublic bool      `json:"is_public"`
}

type CreatePlayListRequest struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

func GeneratePlaylist(playlistQuantity, trackQuantity int) {
	// Requirements:
	// - Artist Account (username, password) -> Access Token
	// - Playlist Id - Perharps we can force it 3 playlist / artist so it won't make a big deal ?
	// - Track Id - Another issue, how can we indicate which track belongs to which artist TODO: Use CSV instead of txt?
	// - Position should be automatically increase in the backend

	// Load track ID
	// trackIds := util.ReadFileTxt("assets/track_ids.txt")

	// Load mock names
	names := util.ReadFileTxt("assets/playlist_names.txt")
	var playlistIds string

	// ownerId, err := uuid.Parse("5f027338-f48a-4d3e-854b-db9934d94b60")
	// if err != nil {
	// 	log.Fatalf("Failed to parse uuid: %v", err)
	// }

	req := SignInRequest{
		Username: "jimtrung",
		Password: "Trung@123",
	}
	accessToken := signIn(req)

	// Loop through the quantity
	for i := 0; i < playlistQuantity; i++ {
		playlist := CreatePlayListRequest{
			Name:     names[i],
			IsPublic: true,
		}

		id := sendCreatePlaylistRequest(playlist, accessToken)
		playlistIds += fmt.Sprintf("%s\n", id.String())
	}

	util.WriteFileTxt("assets/playlist_ids.txt", playlistIds)
}

func sendCreatePlaylistRequest(data CreatePlayListRequest, accessToken string) uuid.UUID {
	url := "http://localhost:8080/playlist"

	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatalf("Failed to create playlist request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send create playlist request: %v", err)
	}
	defer resp.Body.Close()

	pl := Playlist{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&pl); err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	fmt.Println("Create playlist:", pl.Name, "| Status:", resp.Status)

	return pl.Id
}
