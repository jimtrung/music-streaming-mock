package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/jimtrung/music-streaming-mock-data/model"
	"github.com/jimtrung/music-streaming-mock-data/util"
)

func GenerateUser(quantity int) {
	usernames := util.ReadFileTxt("./assets/usernames.txt")
	passwords := util.ReadFileTxt("./assets/passwords.txt")
	accounts := []model.Account{}

	for i := 1; i <= quantity; i++ {
		j := rand.Intn(len(usernames))
		username := usernames[j]
		usernames = append(usernames[:j], usernames[j+1:]...)

		password := passwords[rand.Intn(len(passwords))]
		email := fmt.Sprintf("%s@gmail.com", username)

		account := model.Account{
			Username: username,
			Password: password,
		}

		accounts = append(accounts, account)

		request := SignUpRequest{username, email, password}

		sendSignUpRequest(request)

		util.WriteFileCSV("./assets/accounts.csv", accounts)
	}
}

func sendSignUpRequest(data SignUpRequest) {
	url := "http://127.0.0.1:5255/auth/signup"

	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatalf("Failed to create user request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send sign up request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Create user:", data.Username, "| Status:", resp.Status)
}
