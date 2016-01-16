package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

func writeTemp(lines chan string) string {
	tmp, err := ioutil.TempFile(os.TempDir(), "slackcat-")
	failOnError(err, "unable to create tmpfile", false)

	w := bufio.NewWriter(tmp)
	for line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp.Name()
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
	app.Usage = "redirect text and files to slack"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to Slack continuously instead of uploading a single snippet",
		},
		cli.BoolFlag{
			Name:  "pre, p",
			Usage: "Write messages as preformatted text instead of plaintext",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to Slack. Useful for testing",
		},
		cli.BoolFlag{
			Name:  "configure",
			Usage: "Configure Slackcat via oauth",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("configure") {
			configureOA()
			os.Exit(0)
		}

		token := readConfig()
		fileName := c.String("filename")

		if c.String("channel") == "" {
			exitErr(fmt.Errorf("no channel provided!"))
		}

		slackcat, err := newSlackCat(token, c.String("channel"))
		failOnError(err, "Slack API Error", true)

		if len(c.Args()) > 0 {
			if c.Bool("stream") {
				output("filepath provided, ignoring stream option")
			}
			filePath := c.Args()[0]
			if fileName == "" {
				fileName = filepath.Base(filePath)
			}
			slackcat.postFile(filePath, fileName, c.Bool("noop"))
			os.Exit(0)
		}

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
                        if c.String("filename") != "" {
                                output("stream provided, ignoring filename option")
                        }
			output("starting stream")
			go slackcat.addToStreamQ(lines)
			go slackcat.processStreamQ(c.Bool("noop"), c.Bool("pre"))
			go slackcat.trap()
			select {}
		} else if c.String("filename") == "" {
                        slackcat.postMsgs(lines, c.Bool("noop"), c.Bool("pre"))
			os.Exit(0)
                } else {
                        if c.Bool("pre") {
                                output("filename provided, ignoring preformat option")
                        }
                        filePath := writeTemp(lines)
                        defer os.Remove(filePath)
                        slackcat.postFile(filePath, fileName, c.Bool("noop"))
                        os.Exit(0)
                }
	}

	app.Run(os.Args)

}
