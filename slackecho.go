package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/slack"
)

type SlackEcho struct {
	api         *slack.Slack
	opts        *slack.ChatPostMessageOpt
	queue       *StreamQ
	shutdown    chan os.Signal
	channelName string
	channelId   string
}

func newSlackEcho(token, channelName string) (*SlackEcho, error) {
	se := &SlackEcho{
		api:         slack.New(token),
		opts:        &slack.ChatPostMessageOpt{AsUser: true},
		queue:       newStreamQ(),
		shutdown:    make(chan os.Signal, 1),
		channelName: channelName,
	}
	err := se.lookupSlackId()
	if err != nil {
		return nil, err
	}
	signal.Notify(se.shutdown, os.Interrupt)
	return se, nil
}

func (se *SlackEcho) trap() {
	sigcount := 0
	for sig := range se.shutdown {
		if sigcount > 0 {
			exitErr(fmt.Errorf("aborted"))
		}
		output(fmt.Sprintf("got signal: %s", sig.String()))
		output("press ctrl+c again to exit immediately")
		sigcount++
		go se.exit()
	}
}

func (se *SlackEcho) exit() {
	for {
		if se.queue.isEmpty() {
			os.Exit(0)
		} else {
			output("flushing remaining messages to Slack...")
			time.Sleep(3 * time.Second)
		}
	}
}

//Lookup Slack id for channel, group, or im
func (se *SlackEcho) lookupSlackId() error {
	api := se.api
	channel, err := api.FindChannelByName(se.channelName)
	if err == nil {
		se.channelId = channel.Id
		return nil
	}
	group, err := api.FindGroupByName(se.channelName)
	if err == nil {
		se.channelId = group.Id
		return nil
	}
	im, err := api.FindImByName(se.channelName)
	if err == nil {
		se.channelId = im.Id
		return nil
	}
	fmt.Println(err)
	return fmt.Errorf("No such channel, group, or im")
}

func (se *SlackEcho) addToStreamQ(lines chan string) {
	for line := range lines {
		se.queue.add(line)
	}
	se.exit()
}

//TODO: handle messages with length exceeding maximum for Slack chat
func (se *SlackEcho) processStreamQ(noop bool, pre bool) {
	if !(se.queue.isEmpty()) {
		msglines := se.queue.flush()
		if noop {
			output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), se.channelName))
		} else {
			se.postMsg(msglines, pre, "\n")
		}
	}
	time.Sleep(3 * time.Second)
	se.processStreamQ(noop, pre)
}

func (se *SlackEcho) postMsg(msglines []string, pre bool, sep string){
        fmtStr := "%s"
	if pre {
                fmtStr = "```%s```"
	}
	msg := fmt.Sprintf(fmtStr, strings.Join(msglines, sep))
	err := se.api.ChatPostMessage(se.channelId, msg, se.opts)
	failOnError(err, "", true)
	output(fmt.Sprintf("posted %s message lines to %s", strconv.Itoa(strings.Count(msg, "\n") + 1), se.channelName))
}

func (se *SlackEcho) postLines(lines chan string, noop bool, pre bool) {
        msglines := []string{}
        for line := range lines {
                msglines = append(msglines, line)
        }
        if noop {
                output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), se.channelName))
        } else {
                se.postMsg(msglines, pre, "\n")
        }
}
