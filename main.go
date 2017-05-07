package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
)

const (
	address      = "localhost"
	port         = 8081
	errTag       = "err: %s"
	decodeErrTag = "decode err: %s"
)

func newController(l log.Logger) *Controller {
	return &Controller{
		logger: l,
	}
}

func main() {
	l := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	c := newController(l)
	c.getDB()
	defer c.db.Close()
	http.HandleFunc("/health", c.healthCheck)
	http.HandleFunc("/polls/create", c.createPoll)
	http.HandleFunc("/polls/get", c.getPoll)
	http.HandleFunc("/polls/vote", c.votePoll)
	fmt.Printf("Service is running at %s:%d ...\n", address, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
	if err != nil {
		c.logger.Log("ListenAndServe: ", err)
	}
}

func (c *Controller) healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func (c *Controller) getDB() {
	// Test connection to db
	db, err := sql.Open("mysql", "root:root@/poll_service?parseTime=true")
	if err != nil {
		c.logger.Log(errTag, err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		c.logger.Log(errTag, err)
	}
	c.db = db
}

func (c *Controller) writePostResponse(w http.ResponseWriter, r PostResponse) {
	d, err := json.Marshal(r)
	if err != nil {
		c.logger.Log(errTag, err)
	}
	w.Write(d)
}

func (c *Controller) createPoll(w http.ResponseWriter, r *http.Request) {
	var resp PostResponse
	defer func() { c.writePostResponse(w, resp) }()
	// Read and decode request body
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		c.logger.Log(decodeErrTag, err)
		resp.Status = http.StatusBadRequest
		return
	}
	defer r.Body.Close()

	createPollStmt, err := c.db.Prepare("INSERT polls SET created=?, question=?, story_id=?, part_id=?, duration_days=?")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer createPollStmt.Close()

	dateNow := time.Now().UTC()
	_, err = createPollStmt.Exec(dateNow, req.Question, req.StoryID, req.PartID, req.DurationDays)
	if err != nil {
		c.logger.Log(errTag, err)
	}

	createChoicesStmt, err := c.db.Prepare("INSERT INTO choices (choice, choice_index, votes, part_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer createChoicesStmt.Close()
	for i, choice := range req.Choices {
		_, err = createChoicesStmt.Exec(choice.Choice, i+1, 0, req.PartID)
		if err != nil {
			c.logger.Log(errTag, err)
		}
	}
	resp.Status = http.StatusOK
}

func (c *Controller) getPoll(w http.ResponseWriter, r *http.Request) {
	partID := r.URL.Query().Get("partId")
	username := r.URL.Query().Get("username")

	var resp Response
	getPollStmt, err := c.db.Prepare("SELECT question, created, duration_days FROM polls WHERE part_id = ?")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer getPollStmt.Close()

	err = getPollStmt.QueryRow(partID).Scan(&resp.Question, &resp.Created, &resp.DurationDays)
	if err != nil {
		c.logger.Log(errTag, err)
	}

	endDate := resp.Created.AddDate(0, 0, resp.DurationDays)
	resp.PollClosed = time.Now().UTC().After(endDate)

	getUserVoteStmt, err := c.db.Prepare("SELECT choice_index FROM votes WHERE part_id = ? AND username = ?")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer getUserVoteStmt.Close()

	var userVote int
	err = getUserVoteStmt.QueryRow(partID, username).Scan(&userVote)
	if err != nil && err != sql.ErrNoRows {
		c.logger.Log(errTag, err)
	}
	resp.UserVote = userVote

	getChoicesStmt, err := c.db.Prepare("SELECT choice, votes FROM choices WHERE part_id = ? AND choice_index = ?")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer getChoicesStmt.Close()

	var totalVotes int
	for i := 1; i < 5; i++ {
		var choice Choice
		err = getChoicesStmt.QueryRow(partID, i).Scan(&choice.Choice, &choice.Votes)
		if err != nil {
			c.logger.Log(errTag, err)
		}
		choice.ID = i
		totalVotes += choice.Votes
		resp.Choices = append(resp.Choices, choice)
	}

	resp.TotalVotes = totalVotes
	d, err := json.Marshal(resp)
	if err != nil {
		c.logger.Log(errTag, err)
	}

	w.Write(d)
}

func (c *Controller) votePoll(w http.ResponseWriter, r *http.Request) {
	var resp PostResponse
	defer func() { c.writePostResponse(w, resp) }()
	var req VoteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		c.logger.Log(decodeErrTag, err)
		resp.Status = http.StatusBadRequest
		return
	}
	defer r.Body.Close()

	if c.alreadyVoted(req) {
		resp.Status = http.StatusOK
		return
	}

	voteStmt, err := c.db.Prepare("UPDATE choices SET votes = votes + 1 WHERE part_id = ? AND choice_index = ?")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer voteStmt.Close()

	_, err = voteStmt.Exec(req.PartID, req.ChoiceIndex)
	if err != nil {
		c.logger.Log(errTag, err)
	}

	trackVoteStmt, err := c.db.Prepare("INSERT INTO votes (part_id, choice_index, username) VALUES (?, ?, ?)")
	if err != nil {
		c.logger.Log(errTag, err)
	}
	defer trackVoteStmt.Close()

	_, err = trackVoteStmt.Exec(req.PartID, req.ChoiceIndex, req.Username)
	if err != nil {
		c.logger.Log(errTag, err)
	}
	resp.Status = http.StatusOK
}

func (c *Controller) alreadyVoted(req VoteRequest) bool {
	checkRow, err := c.db.Prepare("SELECT id FROM votes WHERE username = ? AND part_id = ?")
	if err != nil {
		c.logger.Log(decodeErrTag, err)
	}
	defer checkRow.Close()

	var id int
	err = checkRow.QueryRow(req.Username, req.PartID).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		c.logger.Log(decodeErrTag, err)
	}
	return true
}
