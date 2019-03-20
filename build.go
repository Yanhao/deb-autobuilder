package main

import (
	"bufio"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/exec"
	"time"
)

func checkDsc(dscFile string) error {
	return nil
}

func buildPackage(dscFile string) {
	os.Chdir(c.IncomeDir)
	if err := checkDsc(dscFile); err != nil {
		logger.Println(err)
		return;
	}

	c := exec.Command("cowbuilder", "build", dscFile)
	output, err := c.StdoutPipe()
	if err != nil {
		logger.Println(err)
		return;
	}

	err = c.Start()
	if err != nil {
		logger.Println(err)
		return;
	}

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := c.Wait(); err != nil {
		logger.Println(err)
		return;
	}

	fmt.Println(dscFile)
}

func removeDscFileFromDB(dscFile string) error {
	_, err := db.Exec("delete from need_to_build where dsc_file = \"" + dscFile + "\";")
	return err
}

func build() {
	ticker := time.NewTicker(time.Duration(c.TimeInterval) * time.Second)
	defer ticker.Stop()

	for db.Ping() != nil {
		time.Sleep(time.Second * 5)
	}

	for {
		select {
		case _ = <-ticker.C:
			rows, err := db.Query("select * from need_to_build;")
			if err != nil {
				logger.Println("Failed to run db.Query")
				break
			}

			filesNeedToBeRemove := []string{}

			for rows.Next() {
				var dscFile string
				if err := rows.Scan(&dscFile); err != nil {
					continue
				}
				go buildPackage(dscFile)

				filesNeedToBeRemove = append(filesNeedToBeRemove, dscFile)

			}
			rows.Close()

			for _, i := range filesNeedToBeRemove {
				if err := removeDscFileFromDB(i); err != nil {
					logger.Println("Failed to remove ", i, "from db")
				}
			}
		}
	}
}
