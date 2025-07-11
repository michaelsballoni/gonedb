package test

import (
	"database/sql"
	"fmt"
	"os"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"
	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func GetTestDb(name string) *sql.DB {
	fmt.Printf("GetTestDb: %s\n", name)
	tmp_file, tmp_err := os.CreateTemp(os.TempDir(), name)
	if tmp_err != nil {
		panic(fmt.Sprintf("GetTestDb fails CreateTemp: %v", tmp_err))
	}
	tmp_file_name := tmp_file.Name()
	fmt.Printf("GetTestDb: %s: %s\n", name, tmp_file_name)

	os.Remove(tmp_file_name)
	db, err := sql.Open("sqlite3", tmp_file_name)
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
