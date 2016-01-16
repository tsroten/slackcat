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
	fmt.Printf("%s %s\n", cyan("slackcat"), s)
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
	app.Name = "slackcat"
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
		cli.BoolFlag{
			Name:  "configure",
			Usage: "Configure Slackcat via oauth",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("configure") {
			configureOA()
			os.Exit(0)
		}

		token := readConfig()

		if c.String("channel") == "" {
			exitErr(fmt.Errorf("no channel provided!"))
		}

		slackcat, err := newSlackCat(token, c.String("channel"))
		failOnError(err, "Slack API Error", true)

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
			output("starting stream")
			go slackcat.addToStreamQ(lines)
			go slackcat.processStreamQ(c.Bool("noop"), c.Bool("pre"))
			go slackcat.trap()
			select {}
                } else {
                        slackcat.postMsgs(lines, c.Bool("noop"), c.Bool("pre"))
			os.Exit(0)
                }
	}

	app.Run(os.Args)

}
