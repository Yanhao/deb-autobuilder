package main

import (
	"fmt"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type config struct {
	LogFile      string
	IncomeDir    string
	TimeInterval uint
	DB           string
}

var (
	c      config
	logger *log.Logger
)

func main() {
	var configFile string
	var verbose bool
	var logFile string

	flag.StringVar(&configFile, "c",
		"/etc/deb-go/config.toml", "Specify configuration file location")
	flag.BoolVar(&verbose, "v", false, "Be more verbose")
	flag.StringVar(&logFile, "l", "", "Specify log file")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		fmt.Println(err)
	}

	if c.LogFile != "" {
		logFile = c.LogFile
	}

	if logFile == "" {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Faeild to open log file ", logFile)
		}
		defer f.Close()
		logger = log.New(f, "", log.LstdFlags)
	}

	checkEnv()

	go build()
	fileMonit()
}
