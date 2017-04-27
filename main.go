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
	Question string   `json:"question"`
	StoryID  string   `json:"storyId"`
	PartID   string   `json:"partId"`
	Choices  []Choice `json:"choices"`
}

type Response struct {
	Question string    `json:"question"`
	Created  time.Time `json:"created"`
	Choices  []Choice  `json:"choices"`
}

type Choice struct {
	Choice string `json:"choice"`
	Votes  int    `json:"votes"`
}

type controller struct {
	db *sql.DB
}

func newController(db *sql.DB) *controller {
	return &controller{
		db: db,
	}
}

func main() {
	c := newController(getDB())
	defer c.db.Close()
	http.HandleFunc("/health", c.healthCheck)
	http.HandleFunc("/polls/create", c.createPoll)
	http.HandleFunc("/polls/get", c.getPoll)
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (c *controller) healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func getDB() *sql.DB {
	// Test connection to db
	db, err := sql.Open("mysql", "root:root@/poll_service?parseTime=true")
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	return db
}

func (c *controller) createPoll(w http.ResponseWriter, r *http.Request) {
	// Read and decode request body
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Fatalf("decode err: %s", err)
	}
	defer r.Body.Close()

	createPollStmt, err := c.db.Prepare("INSERT polls SET created=?, question=?, story_id=?, part_id=?")
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer createPollStmt.Close()

	datetime := time.Now().UTC()
	_, err = createPollStmt.Exec(datetime, req.Question, req.StoryID, req.PartID)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	createChoicesStmt, err := c.db.Prepare("INSERT INTO choices (choice, choice_index, votes, part_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer createChoicesStmt.Close()
	for i, c := range req.Choices {
		_, err = createChoicesStmt.Exec(c.Choice, i+1, 0, req.PartID)
		if err != nil {
			log.Fatalf("err: %s", err)
		}
	}

	fmt.Fprint(w, "OK")
}

func (c *controller) getPoll(w http.ResponseWriter, r *http.Request) {
	partID := r.URL.Query().Get("partId")

	var resp Response
	getPollStmt, err := c.db.Prepare("SELECT question, created FROM polls WHERE part_id = ?")
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	err = getPollStmt.QueryRow(partID).Scan(&resp.Question, &resp.Created)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	getChoicesStmt, err := c.db.Prepare("SELECT choice, votes FROM choices WHERE part_id = ? AND choice_index = ?")
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	for i := 1; i < 5; i++ {
		var choice Choice
		err = getChoicesStmt.QueryRow(partID, i).Scan(&choice.Choice, &choice.Votes)
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		resp.Choices = append(resp.Choices, choice)
	}

	d, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	w.Write(d)

}
