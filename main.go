package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/internal/command"
	"github.com/prokosna/dementor/lib"
)

func init() {
	// initialize config
	err := dementor.InitConf()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	// CLI
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	cui := &cli.ColoredUi{
		Ui:          ui,
		OutputColor: cli.UiColorGreen,
		InfoColor:   cli.UiColorCyan,
		WarnColor:   cli.UiColorYellow,
		ErrorColor:  cli.UiColorRed,
	}
	c := cli.NewCLI("dementor", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"authenticate": func() (cli.Command, error) {
			return &command.AuthenticateCommand{
				Ui: cui,
			}, nil
		},
		"createProject": func() (cli.Command, error) {
			return &command.CreateProjectCommand{
				Ui: cui,
			}, nil
		},
		"deleteProject": func() (cli.Command, error) {
			return &command.DeleteProjectCommand{
				Ui: cui,
			}, nil
		},
		"fetchFlowsProject": func() (cli.Command, error) {
			return &command.FetchFlowsCommand{
				Ui: cui,
			}, nil
		},
		"fetchJobsFlow": func() (cli.Command, error) {
			return &command.FetchJobsCommand{
				Ui: cui,
			}, nil
		},
		"fetchSchedule": func() (cli.Command, error) {
			return &command.FetchScheduleCommand{
				Ui: cui,
			}, nil
		},
		"scheduleFlow": func() (cli.Command, error) {
			return &command.ScheduleFlowCommand{
				Ui: cui,
			}, nil
		},
		"unscheduleFlow": func() (cli.Command, error) {
			return &command.UnscheduleFlowCommand{
				Ui: cui,
			}, nil
		},
		"uploadProjectZip": func() (cli.Command, error) {
			return &command.UploadProjectZipCommand{
				Ui: cui,
			}, nil
		},
		"kiss": func() (cli.Command, error) {
			return &command.KissCommand{
				Ui: cui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
