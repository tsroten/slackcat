package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var version = "dev-build"

func readIn(lines chan string, tee bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
		if tee {
			fmt.Println(scanner.Text())
		}
	}
	close(lines)
}

func output(s string) {
	cyan := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", cyan("slackecho"), s)
}

func failOnError(err error, msg string, appendErr bool) {
	if err != nil {
		if appendErr {
			exitErr(fmt.Errorf("%s: %s", msg, err))
		} else {
			exitErr(fmt.Errorf("%s", msg))
		}
	}
}

func exitErr(err error) {
	output(color.RedString(err.Error()))
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "slackecho"
	app.Usage = "redirect text to slack"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to Slack continuously",
		},
		cli.BoolFlag{
			Name:  "pre, p",
			Usage: "Write messages as preformatted text instead of plaintext",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting message to Slack. Useful for testing",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
	}

	app.Action = func(c *cli.Context) {
		config := readConfig()

                team, channel, err := config.parseChannelOpt(c.String("channel"))
                failOnError(err, "", true)

                token := config.teams[team]
                if token == "" {
                        exitErr(fmt.Errorf("no such team: %s", team))
                }

                slackecho, err := newSlackEcho(token, channel)
		failOnError(err, "Slack API Error", true)

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

                if len(c.Args()) > 0 {
                        if c.Bool("noop") {
                                output(fmt.Sprintf("skipped posting messages to %s", c.String("channel")))
                        } else {
                                slackecho.postMsg(c.Args(), c.Bool("pre"), " ")
                        }
                        os.Exit(0)
                }

		if c.Bool("stream") {
			output("starting stream")
			go slackecho.addToStreamQ(lines)
			go slackecho.processStreamQ(c.Bool("noop"), c.Bool("pre"))
			go slackecho.trap()
			select {}
                } else {
                        slackecho.postLines(lines, c.Bool("noop"), c.Bool("pre"))
			os.Exit(0)
                }
	}

	app.Run(os.Args)

}
