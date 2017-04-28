package main

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"database/sql/driver"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreatePoll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO polls").WithArgs(AnyTime{}, "testQ", "testStoryID", "testPartID")
	mock.ExpectCommit()
}
