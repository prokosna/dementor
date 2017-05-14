package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type DeleteProjectCommand struct {
	Ui cli.Ui
}

var deleteProjectOpts struct {
	Project string `short:"p" long:"project" description:"Project name" required:"true"`
	dementor.CommonConf
}

var deleteProjectParser *flags.Parser

func init() {
	deleteProjectParser = flags.NewParser(&deleteProjectOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *DeleteProjectCommand) Run(args []string) int {
	_, err := deleteProjectParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(deleteProjectOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.DeleteProjectReq{
		Project:    deleteProjectOpts.Project,
		CommonConf: deleteProjectOpts.CommonConf,
	}
	err = dementor.DeleteProject(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully deleted the project: %s", req.Project))
	return 0
}

func (c *DeleteProjectCommand) Help() string {
	buf := new(bytes.Buffer)
	deleteProjectParser.WriteHelp(buf)
	return buf.String()
}

func (c *DeleteProjectCommand) Synopsis() string {
	return "Delete the project."
}
