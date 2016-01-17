# slackecho
Slackecho is a simple commandline utility to post messages to Slack.

It is forked from, and designed to complement, [Slackcat](https://github.com/vektorlab/slackcat).


## Configuration

If you already use [Slackcat](https://github.com/vektorlab/slackcat), you're already configured. If not, download Slackcat and follow the instructions. Your new config file will work for Slackcat and Slackecho.

## Usage
Echo a string as a message:
```bash
$ slackecho --channel general Good morning!
*slackecho* posted 1 message lines to general
```

Pipe command output as a message:
```bash
$ echo -e "hi\nthere" | slackecho --channel general
*slackecho* posted 2 message lines to general
```

Post a message as preformatted text:
```bash
$ echo -e "print('Hello world!')" | slackecho --pre --channel general
*slackecho* posted 1 message lines to general
```

Stream input continously as preformatted text:
```bash
$ tail -F -n0 /path/to/log | slackecho --channel general --stream --pre
*slackecho* posted 5 message lines to general
*slackecho* posted 2 message lines to general
...
```

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--stream, -s | Stream messages to Slack continuously
--pre, -p | Write messages as preformatted text
--noop | Skip posting message to Slack. Useful for testing
--configure | Configure Slackcat/Slackecho via oauth
--channel, -c | Slack channel, group, or user to post to
