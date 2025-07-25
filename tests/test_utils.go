package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"
	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func GetTestDb(name string) *sql.DB {
	db_file_path := filepath.Join(os.TempDir(), name)
	fmt.Printf("GetTestDb: %s - %s\n", name, db_file_path)

	os.Remove(db_file_path)

	db, err := sql.Open("sqlite3", db_file_path)
	if err != nil {
		panic(fmt.Sprintf("GetTestDb fails Open: %v", err))
	}

	gonedb.Setup(db)
	gonedb.Strings.FlushCaches()

	return db
}

func GetTestStringId(db *sql.DB, val string) int64 {
	id, err := gonedb.Strings.GetId(db, val)
	AssertNoError(err)
	return id
}

func AssertEqual[T comparable](expected T, got T) {
	if got != expected {
		debug.PrintStack()
		panic(fmt.Sprintf("AssertEqual: expected %v - got %v", expected, got))
	}
}

func AssertTrue(check bool) {
	if !check {
		debug.PrintStack()
		panic("AssertTrue failed!")
	}
}

func AssertError(err error) {
	if err == nil {
		debug.PrintStack()
		panic("AssertError failed!")
	}
}

func AssertNoError(err error) {
	if err != nil {
		debug.PrintStack()
		panic(fmt.Sprintf("AssertNoError failed: %v", err))
	}
}
