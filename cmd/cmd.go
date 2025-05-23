package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func file_exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <op> <db file path>\n", os.Args[0])
		return
	}

	op := os.Args[1]

	db_file := os.Args[2]
	fmt.Println("op:", op, "db_file:", db_file, "file_exists:", file_exists(db_file))

	fmt.Println("Opening database...")
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Printf("Opening database file failed: %s\n", err)
		os.Exit(1)
		return
	}
	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		fmt.Printf("Getting SQLite version failed: %s\n", err)
		os.Exit(1)
		return
	}
	fmt.Println("SQLite version:", version)

	if op == "setup" {
		fmt.Println("Setting up database...")
		gonedb.Setup(db)
		fmt.Println("Database created!")
		return
	}

	fmt.Println("Unknown op:", op)
	os.Exit(1)
}
