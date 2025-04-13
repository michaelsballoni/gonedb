package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func GetTestDb(name string) *sql.DB {
	os.Remove(name)
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		panic(fmt.Sprintf("GetTestDb fails: %v", err))
	}
	gonedb.Setup(db)
	return db
}

func AssertEqual[T comparable](t *testing.T, expected T, got T) {
	if got != expected {
		t.Errorf("AssertEqual: expected %v - got %v", expected, got)
		t.Fail()
	}
}

func AssertTrue(t *testing.T, check bool) {
	if !check {
		t.Errorf("AssertTrue fails")
		t.Fail()
	}
}

func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("AssertError: %v", err)
		t.Fail()
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("AssertNoError: %v", err)
		t.Fail()
	}
}
