package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/k0kubun/pp"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type FetchScheduleCommand struct {
	Ui cli.Ui
}

var fetchScheduleOpts struct {
	ProjectId string `short:"p" long:"projectId" description:"Project ID" required:"true"`
	FlowId    string `short:"f" long:"flowId" description:"Flow ID" required:"true"`
	dementor.CommonConf
}

var fetchScheduleParser *flags.Parser

func init() {
	fetchScheduleParser = flags.NewParser(&fetchScheduleOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *FetchScheduleCommand) Run(args []string) int {
	_, err := fetchScheduleParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(fetchScheduleOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.FetchScheduleReq{
		ProjectId:  fetchScheduleOpts.ProjectId,
		FlowId:     fetchScheduleOpts.FlowId,
		CommonConf: fetchScheduleOpts.CommonConf,
	}
	res, err := dementor.FetchSchedule(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	pretty := pp.Sprint(*res)
	c.Ui.Output(fmt.Sprintf("Successfully fetched the schedule: \n%s", pretty))
	return 0
}

func (c *FetchScheduleCommand) Help() string {
	buf := new(bytes.Buffer)
	fetchScheduleParser.WriteHelp(buf)
	return buf.String()
}

func (c *FetchScheduleCommand) Synopsis() string {
	return "Fetch the schedule of the flow."
}
