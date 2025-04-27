package test

import (
	"database/sql"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func GetTestDb(name string) *sql.DB {
	tmp_file, tmp_err := os.CreateTemp(os.TempDir(), name)
	if tmp_err != nil {
		panic(fmt.Sprintf("GetTestDb fails CreateTemp: %v", tmp_err))
	}
	name = tmp_file.Name()

	os.Remove(name)
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		panic(fmt.Sprintf("GetTestDb fails Open: %v", err))
	}
	gonedb.Setup(db)
	return db
}

func AssertEqual[T comparable](t *testing.T, expected T, got T) {
	if got != expected {
		debug.PrintStack()
		t.Errorf("AssertEqual: expected %v - got %v", expected, got)
		t.Fatal()
	}
}

func AssertTrue(t *testing.T, check bool) {
	if !check {
		debug.PrintStack()
		t.Errorf("AssertTrue fails")
		t.Fatal()
	}
}

func AssertError(t *testing.T, err error) {
	if err == nil {
		debug.PrintStack()
		t.Errorf("AssertError: %v", err)
		t.Fatal()
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Errorf("AssertNoError: %v", err)
		t.Fatal()
	}
}

func GetStringId(t *testing.T, db *sql.DB, val string) int64 {
	id, err := gonedb.Strings.GetId(db, val)
	AssertNoError(t, err)
	return id
}
