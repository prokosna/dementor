package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type UnscheduleFlowCommand struct {
	Ui cli.Ui
}

var unscheduleFlowOpts struct {
	ScheduleId string `short:"s" long:"scheduleId" description:"Schedule ID" required:"true"`
	dementor.CommonConf
}

var unscheduleFlowParser *flags.Parser

func init() {
	unscheduleFlowParser = flags.NewParser(&unscheduleFlowOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *UnscheduleFlowCommand) Run(args []string) int {
	_, err := unscheduleFlowParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(unscheduleFlowOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.UnscheduleFlowReq{
		ScheduleId: unscheduleFlowOpts.ScheduleId,
		CommonConf: unscheduleFlowOpts.CommonConf,
	}
	err = dementor.UnscheduleFlow(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Successfully unscheduled the flow [%s]", req.ScheduleId))
	return 0
}

func (c *UnscheduleFlowCommand) Help() string {
	buf := new(bytes.Buffer)
	unscheduleFlowParser.WriteHelp(buf)
	return buf.String()
}

func (c *UnscheduleFlowCommand) Synopsis() string {
	return "Unschedule the flow."
}
