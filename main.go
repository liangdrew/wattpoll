package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Request struct {
	Question string    `json:"question"`
	StoryID  string    `json:"storyId"`
	PartID   string    `json:"partId"`
	Choices  []Choices `json:"choices"`
}

type Choices struct {
	Choice string `json:"choice"`
}

func main() {
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/create", createPoll)
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	// Read and decode request body
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Fatalf("decode err: %s", err)
	}
	defer r.Body.Close()

	// Test connection to db
	db, err := sql.Open("mysql", "root:root@/poll_service")
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	createPollStmt, err := db.Prepare("INSERT polls SET created=?, question=?, story_id=?, part_id=?")
	defer createPollStmt.Close()

	datetime := time.Now().UTC()
	_, err = createPollStmt.Exec(datetime, req.Question, req.StoryID, req.PartID)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	createChoicesStmt, err := db.Prepare("INSERT INTO choices (choice, choice_index, votes, part_id) VALUES (?, ?, ?, ?)")
	defer createChoicesStmt.Close()
	for i, c := range req.Choices {
		_, err = createChoicesStmt.Exec(c.Choice, i+1, 0, req.PartID)
		if err != nil {
			log.Fatalf("err: %s", err)
		}
	}

	fmt.Fprint(w, "OK")

}
