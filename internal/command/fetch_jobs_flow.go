package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/k0kubun/pp"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type FetchJobsCommand struct {
	Ui cli.Ui
}

var fetchJobsOpts struct {
	Project string `short:"p" long:"project" description:"Project name" required:"true"`
	Flow    string `short:"f" long:"flowId" description:"Flow ID" required:"true"`
	dementor.CommonConf
}

var fetchJobsParser *flags.Parser

func init() {
	fetchJobsParser = flags.NewParser(&fetchJobsOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *FetchJobsCommand) Run(args []string) int {
	_, err := fetchJobsParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(fetchJobsOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.FetchJobsFlowReq{
		Project:    fetchJobsOpts.Project,
		Flow:       fetchJobsOpts.Flow,
		CommonConf: fetchJobsOpts.CommonConf,
	}
	res, err := dementor.FetchJobsFlow(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	pretty := pp.Sprint(*res)
	c.Ui.Output(fmt.Sprintf("Successfully fetched the jobs: \n%s", pretty))
	return 0
}

func (c *FetchJobsCommand) Help() string {
	buf := new(bytes.Buffer)
	fetchJobsParser.WriteHelp(buf)
	return buf.String()
}

func (c *FetchJobsCommand) Synopsis() string {
	return "Fetch jobs of the flow."
}
