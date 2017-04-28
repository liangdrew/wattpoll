package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type VoteRequest struct {
	StoryID     string `json:"storyId"`
	PartID      string `json:"partId"`
	Choice      string `json:"choice"`
	ChoiceIndex int    `json:"choiceIndex"`
	Username    string `json:"username"`
}

type Request struct {
	Question 	string   `json:"question"`
	StoryID  	string   `json:"storyId"`
	PartID   	string   `json:"partId"`
	DurationDays    int  	 `json:"durationDays"`
	Choices  	[]Choice `json:"choices"`
}

type Response struct {
	Question   	string    `json:"question"`
	TotalVotes 	int       `json:"totalVotes"`
	UserVote   	int       `json:"userVote"`
	Created    	time.Time `json:"created"`
	DurationDays    int   	  `json:"durationDays"`
	PollClosed 	bool	  `json:"pollClosed"`
	Choices    	[]Choice  `json:"choices"`
}

type PostResponse struct {
	Status int `json:"status"`
}

type Choice struct {
	ID     int    `json:"id"`
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
	http.HandleFunc("/polls/vote", c.votePoll)
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
		log.Printf("err: %s", err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Printf("err: %s", err)
	}
	return db
}

func writePostResponse(w http.ResponseWriter, r PostResponse) {
	d, err := json.Marshal(r)
	if err != nil {
		log.Printf("err: %s", err)
	}
	w.Write(d)
}

func (c *controller) createPoll(w http.ResponseWriter, r *http.Request) {
	var resp PostResponse
	defer func() { writePostResponse(w, resp) }()
	// Read and decode request body
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("decode err: %s", err)
		resp.Status = http.StatusBadRequest
		return
	}
	defer r.Body.Close()

	createPollStmt, err := c.db.Prepare("INSERT polls SET created=?, question=?, story_id=?, part_id=?, duration_days=?")
	if err != nil {
		log.Printf("err: %s", err)
	}
	defer createPollStmt.Close()

	dateNow := time.Now().UTC()
	_, err = createPollStmt.Exec(dateNow, req.Question, req.StoryID, req.PartID, req.DurationDays)
	if err != nil {
		log.Printf("err: %s", err)
	}

	createChoicesStmt, err := c.db.Prepare("INSERT INTO choices (choice, choice_index, votes, part_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("err: %s", err)
	}
	defer createChoicesStmt.Close()
	for i, c := range req.Choices {
		_, err = createChoicesStmt.Exec(c.Choice, i+1, 0, req.PartID)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	resp.Status = http.StatusOK
}

func (c *controller) getPoll(w http.ResponseWriter, r *http.Request) {
	partID := r.URL.Query().Get("partId")
	username := r.URL.Query().Get("username")

	var resp Response
	getPollStmt, err := c.db.Prepare("SELECT question, created, duration_days FROM polls WHERE part_id = ?")
	if err != nil {
		log.Printf("err: %s", err)
	}
	defer getPollStmt.Close()

	err = getPollStmt.QueryRow(partID).Scan(&resp.Question, &resp.Created, &resp.DurationDays)
	if err != nil {
		log.Printf("err: %s", err)
	}

	endDate := resp.Created.AddDate(0, 0, resp.DurationDays)
	resp.PollClosed = time.Now().UTC().After(endDate)

	getUserVoteStmt, err := c.db.Prepare("SELECT choice_index FROM votes WHERE part_id = ? AND username = ?")
	if err != nil {
		log.Printf("err: %s", err)
	}
	defer getUserVoteStmt.Close()

	var userVote int
	err = getUserVoteStmt.QueryRow(partID, username).Scan(&userVote)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("err: %s", err)
	}
	resp.UserVote = userVote

	getChoicesStmt, err := c.db.Prepare("SELECT choice, votes FROM choices WHERE part_id = ? AND choice_index = ?")
	if err != nil {
		log.Printf("err: %s", err)
	}
	defer getChoicesStmt.Close()

	var totalVotes int
	for i := 1; i < 5; i++ {
		var choice Choice
		err = getChoicesStmt.QueryRow(partID, i).Scan(&choice.Choice, &choice.Votes)
		if err != nil {
			log.Printf("err: %s", err)
		}
		choice.ID = i
		totalVotes += choice.Votes
		resp.Choices = append(resp.Choices, choice)
	}

	resp.TotalVotes = totalVotes
	d, err := json.Marshal(resp)
	if err != nil {
		log.Printf("err: %s", err)
	}

	w.Write(d)
}

func (c *controller) votePoll(w http.ResponseWriter, r *http.Request) {
	var resp PostResponse
	defer func() { writePostResponse(w, resp) }()
	var req VoteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("decode err: %s", err)
		resp.Status = http.StatusBadRequest
		return
	}
	defer r.Body.Close()

	if !c.alreadyVoted(req) {
		voteStmt, err := c.db.Prepare("UPDATE choices SET votes = votes + 1 WHERE part_id = ? AND choice = ? AND choice_index = ?")
		if err != nil {
			log.Printf("err: %s", err)
		}
		defer voteStmt.Close()

		_, err = voteStmt.Exec(req.PartID, req.Choice, req.ChoiceIndex)
		if err != nil {
			log.Printf("err: %s", err)
		}

		trackVoteStmt, err := c.db.Prepare("INSERT INTO votes (part_id, choice_index, username) VALUES (?, ?, ?)")
		if err != nil {
			log.Printf("err: %s", err)
		}
		defer trackVoteStmt.Close()

		_, err = trackVoteStmt.Exec(req.PartID, req.ChoiceIndex, req.Username)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	resp.Status = http.StatusOK
}

func (c *controller) alreadyVoted(req VoteRequest) bool {
	checkRow, err := c.db.Prepare("SELECT id FROM votes WHERE username = ? AND part_id = ?")
	if err != nil {
		log.Printf("decode err: %s", err)
	}
	defer checkRow.Close()

	var id int
	err = checkRow.QueryRow(req.Username, req.PartID).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Printf("decode err: %s", err)
	}
	return true
}
