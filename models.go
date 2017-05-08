package main

import (
	"database/sql"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/rubyist/circuitbreaker"
)

type VoteRequest struct {
	StoryID     string `json:"storyId"`
	PartID      string `json:"partId"`
	ChoiceIndex int    `json:"choiceIndex"`
	Username    string `json:"username"`
}

type Request struct {
	Question     string   `json:"question"`
	StoryID      string   `json:"storyId"`
	PartID       string   `json:"partId"`
	DurationDays int      `json:"durationDays"`
	Choices      []Choice `json:"choices"`
}

type Response struct {
	Question     string    `json:"question"`
	TotalVotes   int       `json:"totalVotes"`
	UserVote     int       `json:"userVote"`
	Created      time.Time `json:"created"`
	DurationDays int       `json:"durationDays"`
	PollClosed   bool      `json:"pollClosed"`
	Choices      []Choice  `json:"choices"`
}

type PostResponse struct {
	Status int `json:"status"`
}

type Choice struct {
	ID     int    `json:"id"`
	Choice string `json:"choice"`
	Votes  int    `json:"votes"`
}

type Controller struct {
	db      *sql.DB
	breaker *circuit.Breaker
	logger  log.Logger
}
