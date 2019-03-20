package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yanhao/fsnotify"
	"log"
	"strings"
)

func saveDscFile(filename string, db *sql.DB) {
	if !strings.HasSuffix(filename, ".dsc") {
		return
	}

	insertSQL := strings.Builder{}
	fileSplit := strings.Split(filename, "/")
	basename := fileSplit[len(fileSplit)-1]
	fmt.Fprintf(&insertSQL, "insert into need_to_build values(\"%s\")", basename)

	db.Exec(insertSQL.String())
}

func fileMonit() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if db, err = sql.Open("sqlite3", c.DB); err != nil {
		logger.Fatal("Failed to open database")
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					log.Println("event: ", ev)
					saveDscFile(ev.Name, db)
				}
			case err := <-watcher.Error:
				logger.Println("inotify error: ", err)
			}
		}
	}()

	err = watcher.Watch(c.IncomeDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
	watcher.Close()
	db.Close()

}
