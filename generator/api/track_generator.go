package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/util"
)

func GenerateTracks(quantity int) {
	// Fetch mock data
	titles := util.ReadFileTxt("./assets/track_titles.txt")
	songs := fetchSongNames()
	covers := fetchCoverNames()
	accounts := util.ReadFileCSV("./assets/accounts.csv")
	trackIds := ""

	for i := 1; i <= quantity; i++ {
		for k := 1; k <= rand.Intn(3); k++ {
			// Generate random value
			j := rand.Intn(len(titles))
			title := titles[j]
			titles = append(titles[:j], titles[j+1:]...)
			songPath := songs[rand.Intn(len(songs))]
			coverPath := covers[rand.Intn(len(covers))]

			// Get access_token
			// TODO: Mock 1 track per artist ? Or more ? Random ?
			var accessToken string

			accessToken = signIn(SignInRequest{
				Username: accounts[i-1].Username,
				Password: accounts[i-1].Password,
			})

			trackId := uploadTrack(title, songPath, coverPath, accessToken)
			trackIds += fmt.Sprintf("%s\n", trackId.String())
		}
	}

	util.WriteFileTxt("assets/track_ids.txt", trackIds)
}

func fetchSongNames() []string {
	files, err := os.ReadDir("./assets/songs")
	if err != nil {
		log.Fatalf("Failed to read directory /assets/songs: %v", err)
	}

	names := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	return names
}

func fetchCoverNames() []string {
	files, err := os.ReadDir("./assets/covers")
	if err != nil {
		log.Fatalf("Failed to read directory /assets/covers: %v", err)
	}

	names := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	return names
}

func uploadTrack(title, songPath, coverPath, accessToken string) uuid.UUID {
	url := "http://localhost:8080/track"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ---- audio ----
	audioFile, err := os.Open("./assets/songs/" + songPath)
	if err != nil {
		log.Fatal(err)
	}
	defer audioFile.Close()

	audioPart, err := writer.CreateFormFile("audio", filepath.Base(songPath))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(audioPart, audioFile)

	// ---- cover ----
	coverFile, err := os.Open("./assets/covers/" + coverPath)
	if err != nil {
		log.Fatal(err)
	}
	defer coverFile.Close()

	coverPart, err := writer.CreateFormFile("cover", filepath.Base(coverPath))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(coverPart, coverFile)

	// ---- text field ----
	writer.WriteField("title", title)

	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

	track := CreateTrackResponse{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&track); err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	return track.Id
}
