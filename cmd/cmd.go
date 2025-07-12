package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func file_exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: <db file path>")
		return
	}

	db_file := os.Args[1]
	db_existed := file_exists(db_file)
	fmt.Println("db_file:", db_file, "file_exists:", db_existed)

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
	defer db.Close()

	if !db_existed {
		fmt.Println("Setting up gonedb schema...")
		gonedb.Setup(db)
	}

	scanner := bufio.NewScanner(os.Stdin)
	cmd := gonedb.CreateCmd()
	for {
		prompt, err := cmd.GetPrompt(db)
		if err != nil {
			fmt.Printf("Getting prompt failed: %s\n", err)
			os.Exit(1)
			return
		}
		fmt.Printf("%s> ", prompt)

		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		if line == "quit" {
			return
		}

		output, cmd_err := cmd.ProcessCommand(db, line)
		if cmd_err != nil {
			fmt.Printf("ERROR: %s\n", cmd_err)
		} else if len(output) > 0 {
			fmt.Println(output)
		}
	}
}
