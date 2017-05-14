package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type ScheduleFlowCommand struct {
	Ui cli.Ui
}

var scheduleFlowOpts struct {
	ProjectName    string `short:"p" long:"project" description:"Project name" required:"true"`
	Flow           string `short:"f" long:"flowId" description:"Flow ID" required:"true"`
	CronExpression string `short:"c" long:"cron" description:"Cron expression such as 0 23/30 5,7-10 ? * 6#3" required:"true"`
	dementor.CommonConf
}

var scheduleFlowParser *flags.Parser

func init() {
	scheduleFlowParser = flags.NewParser(&scheduleFlowOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *ScheduleFlowCommand) Run(args []string) int {
	_, err := scheduleFlowParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(scheduleFlowOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.ScheduleFlowReq{
		ProjectName:    scheduleFlowOpts.ProjectName,
		Flow:           scheduleFlowOpts.Flow,
		CronExpression: scheduleFlowOpts.CronExpression,
		CommonConf:     scheduleFlowOpts.CommonConf,
	}
	res, err := dementor.ScheduleFlow(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully scheduled the flow: %s", res.ScheduleId))
	return 0
}

func (c *ScheduleFlowCommand) Help() string {
	buf := new(bytes.Buffer)
	scheduleFlowParser.WriteHelp(buf)
	return buf.String()
}

func (c *ScheduleFlowCommand) Synopsis() string {
	return "Schedule the flow."
}
