package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/k0kubun/pp"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type FetchFlowsCommand struct {
	Ui cli.Ui
}

var fetchFlowsOpts struct {
	Project string `short:"p" long:"project" description:"Project name" required:"true"`
	dementor.CommonConf
}

var fetchFlowsParser *flags.Parser

func init() {
	fetchFlowsParser = flags.NewParser(&fetchFlowsOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *FetchFlowsCommand) Run(args []string) int {
	_, err := fetchFlowsParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(fetchFlowsOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.FetchFlowsProjectReq{
		Project:    fetchFlowsOpts.Project,
		CommonConf: fetchFlowsOpts.CommonConf,
	}
	res, err := dementor.FetchFlowsProject(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	pretty := pp.Sprint(*res)
	c.Ui.Output(fmt.Sprintf("Successfully fetched the flows: \n%s", pretty))
	return 0
}

func (c *FetchFlowsCommand) Help() string {
	buf := new(bytes.Buffer)
	fetchFlowsParser.WriteHelp(buf)
	return buf.String()
}

func (c *FetchFlowsCommand) Synopsis() string {
	return "Fetch flows of the project"
}
