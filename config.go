package main

import (
	"bufio"
	"fmt"
	"os"
)

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exitErr(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.slackcat"
}

func readConfig() string {
	path := getConfigPath()
	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines[0]
}
