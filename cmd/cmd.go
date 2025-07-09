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
		fmt.Println("Usage: <db file path>")
		return
	}

	op := os.Args[1]

	db_file := os.Args[2]
	fmt.Println("db_file:", db_file, "file_exists:", file_exists(db_file))

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

	was_unknown := false
	for {
		if op == "help" || was_unknown {
			fmt.Println("Commands:")
			fmt.Println("setup:", "Set up database for initial use.  You do this once.  Any info in the file is lost.")
			fmt.Println("exit or quit:", "Quit this program")
			fmt.Println("help:", "Display this help ;)")
			was_unknown = false
			continue
		}

		if op == "setup" {
			fmt.Println("Setting up database...")
			gonedb.Setup(db)
			fmt.Println("Database created!")
			continue
		}

		if op == "exit" || op == "quit" {
			return
		}

		fmt.Println("Unknown op:", op)
		was_unknown = true
	}
}
