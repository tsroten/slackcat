package main

import (
	"bufio"
	"fmt"
	"os"
        "strings"
)

type Config struct {
        teams map[string]string
        defaultTeam string
        defaultChannel string
}

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exitErr(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.slackcat"
}

func (c *Config) parseChannelOpt(channel string) (string, string, error) {
        // Use the default channel if none is provided.
        if channel == "" {
                if c.defaultChannel == "" {
                        return "", "", fmt.Errorf("no channel provided!")
                } else {
                        return c.defaultTeam, c.defaultChannel, nil
                }
        }

        // If the channel is prefixed with a team.
        if strings.Contains(channel, ":") {
                s := strings.Split(channel, ":")
                return s[0], s[1], nil
        }

        // Use the default team with the provided channel.
        return c.defaultTeam, channel, nil
}

func readConfig() *Config {
        config := &Config{
                teams: make(map[string]string),
                defaultTeam: "",
                defaultChannel: "",
        }
        lines := readLines(getConfigPath())

        if len(lines) == 1 {
                config.teams["default"] = lines[0]
                config.defaultTeam = "default"
                return config
        }

        for _, line := range lines {
                s := strings.Split(line, "=")
                if len(s) != 2 {
                        exitErr(fmt.Errorf("failed to parse config at: %s", line))
                }

                key := strip(s[0])
                switch key {
                case "default_team":
                        config.defaultTeam = strip(s[1])
                case "default_channel":
                        config.defaultChannel = strip(s[1])
                default:
                        config.teams[key] = strip(s[1])
                }
        }

        return config
}

func strip(s string) string {
        return strings.Replace(s, " ", "", -1)
}

func readLines(path string) []string {
	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
                if scanner.Text() != "" {
                        lines = append(lines, scanner.Text())
                }
	}

	return lines
}
