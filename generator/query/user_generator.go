package query

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimtrung/music-streaming-mock-data/model"
	"github.com/jimtrung/music-streaming-mock-data/util"
	"golang.org/x/crypto/bcrypt"
)

// This will create a new sql query file in query/
func GenerateUser(quantity int) {
	usernames := util.ReadFileTxt("./assets/usernames.txt")
	passwords := util.ReadFileTxt("./assets/passwords.txt")
	booleans := []string{"TRUE", "FALSE"}
	accounts := []model.Account{}

	var querySb strings.Builder
	var idSb strings.Builder

	// Generate the query
	querySb.WriteString("INSERT INTO users (id, username, email, password, role, provider, is_verified, is_premium)\nVALUES")

	for i := 1; i <= quantity; i++ {
		id := uuid.New()
		fmt.Fprintf(&idSb, "%s\n", id.String())

		// Delete after take the username (Make sure there will be no duplicate)
		j := rand.Intn(len(usernames))
		username := usernames[j]
		usernames = append(usernames[:j], usernames[j+1:]...)

		password := passwords[rand.Intn(len(passwords))]
		email := fmt.Sprintf("%s@gmail.com", username)
		isPremium := booleans[rand.Intn(len(booleans))]

		account := model.Account{
			Username: username,
			Password: password,
		}

		accounts = append(accounts, account)

		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		row := fmt.Sprintf(
			"'%s', '%s', '%s', '%s', 'listener'::role_type, 'local'::provider_type, true, %s",
			id, username, email, passwordHash, isPremium,
		)

		comma := ","
		if i == quantity {
			comma = ""
		}

		fmt.Fprintf(&querySb, "\n(%s)%s", row, comma)
	}

	// File name template: timestamp_users_generate_query.sql
	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	queryFilePath := fmt.Sprintf("./query/%s_users_generated_query.sql", timestamp)

	util.WriteFileTxt(queryFilePath, querySb.String())
	util.WriteFileTxt("./assets/user_ids.txt", idSb.String())
	util.WriteFileCSV("./assets/accounts.csv", accounts)
}
