package query

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jimtrung/music-streaming-mock-data/util"
)

func GenerateProfile() {
	ids := util.ReadFileTxt("./assets/user_ids.txt")

	var sb strings.Builder

	// Generate the query
	sb.WriteString("INSERT INTO profiles (user_id)\nVALUES")

	for i, id := range ids {
		comma := ","
		if i == len(ids)-1 {
			comma = ""
		}

		sb.WriteString(fmt.Sprintf("\n('%s')%s", id, comma))
	}

	// File name template: timestamp_profiles_generate_query.sql
	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	queryFile, err := os.Create(fmt.Sprintf("./query/%s_profiles_generated_query.sql", timestamp))
	if err != nil {
		log.Fatalf("Failed to create generated query: %v", err)
	}

	queryFile.WriteString(sb.String())
}
