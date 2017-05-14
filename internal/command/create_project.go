package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type CreateProjectCommand struct {
	Ui cli.Ui
}

var createProjectOpts struct {
	Name        string `short:"p" long:"project" description:"Project name" required:"true"`
	Description string `short:"d" long:"description" description:"Project description" required:"true"`
	dementor.CommonConf
}

var createProjectParser *flags.Parser

func init() {
	createProjectParser = flags.NewParser(&createProjectOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *CreateProjectCommand) Run(args []string) int {
	_, err := createProjectParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(createProjectOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.CreateProjectReq{
		Name:        createProjectOpts.Name,
		Description: createProjectOpts.Description,
		CommonConf:  createProjectOpts.CommonConf,
	}
	_, err = dementor.CreateProject(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully created a project: %s", req.Name))
	return 0
}

func (c *CreateProjectCommand) Help() string {
	buf := new(bytes.Buffer)
	createProjectParser.WriteHelp(buf)
	return buf.String()
}

func (c *CreateProjectCommand) Synopsis() string {
	return "Create a new project."
}
