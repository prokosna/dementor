package command

import (
	"bytes"

	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type AuthenticateCommand struct {
	Ui cli.Ui
}

var authenticateOpts struct {
	dementor.CommonConf
}

var authenticateParser *flags.Parser

func init() {
	authenticateParser = flags.NewParser(&authenticateOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *AuthenticateCommand) Run(args []string) int {
	_, err := authenticateParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(authenticateOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Session ID: %s", id))
	return 0
}

func (c *AuthenticateCommand) Help() string {
	buf := new(bytes.Buffer)
	authenticateParser.WriteHelp(buf)
	return buf.String()
}

func (c *AuthenticateCommand) Synopsis() string {
	return "Fetch a session id by username and password."
}
