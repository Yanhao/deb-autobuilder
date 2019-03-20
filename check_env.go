package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func checkEnv() {
	if os.Geteuid() != 0 {
		fmt.Println("Only root can run this program")
		os.Exit(1)
	}

	if _, err := os.Stat(c.DB); os.IsNotExist(err) {
		logger.Println(c.DB, " doesn't exist, creating it ...")
		db, err := sql.Open("sqlite3", c.DB)
		if err != nil {
			logger.Fatal("Failed to create db file")
		}

		cTableSQL := `
		create table need_to_build (dsc_file text primary key not null);
		`

		if _, err := db.Exec(cTableSQL); err != nil {
			logger.Fatal("Failed to create ")
		}
		db.Close()
	}
}
