package command

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
	"gopkg.in/yaml.v2"
)

type KissCommand struct {
	Ui cli.Ui
}

var kissOpts struct {
	FilePath string `short:"f" long:"filepath" description:"Path to a recipe file" required:"true"`
}

var kissParser *flags.Parser

func init() {
	kissParser = flags.NewParser(&kissOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

type Recipe struct {
	Url      string    `yaml:"url"`
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	FilePath    string `yaml:"filepath"`
	Flows       []Flow `yaml:"flows"`
}

type Flow struct {
	Name string `yaml:"name"`
	Cron string `yaml:"cron"`
}

func (c *KissCommand) askWhetherContinue() {
	a, err := c.Ui.Ask("Continue? [y/n] >")
	if strings.ToLower(a) != "y" || err != nil {
		c.Ui.Error("Terminated. Bye!")
		os.Exit(1)
	}
}

func (c *KissCommand) createProject(id string, commonConf dementor.CommonConf, project Project) error {
	// Create the project
	c.Ui.Output(fmt.Sprintf("The project [%s] will be created...", project.Name))
	creq := &dementor.CreateProjectReq{
		Name:        project.Name,
		Description: project.Description,
		CommonConf:  commonConf,
	}
	_, err := dementor.CreateProject(id, creq)
	return err
}

func (c *KissCommand) processProject(id string, commonConf dementor.CommonConf, project Project) error {
	errorCount := 0
	// Fetch flows
	freq := &dementor.FetchFlowsProjectReq{
		Project:    project.Name,
		CommonConf: commonConf,
	}
	fres, err := dementor.FetchFlowsProject(id, freq)
	if err != nil {
		// The project does not exist
		c.Ui.Output(fmt.Sprintf("The project [%s] does not exist.\n", project.Name))

		// Create the project
		err = c.createProject(id, commonConf, project)
		if err != nil {
			c.Ui.Error(err.Error())
			c.Ui.Error(fmt.Sprintf("Failed to create the project [%s].", project.Name))
			return err
		}
		c.Ui.Info("Done!\n")
	} else {
		// The project exists.
		a, err := c.Ui.Ask(fmt.Sprintf("The project [%s] exists. Replace(r) or Update(u)? [r/u] >", project.Name))
		if (strings.ToLower(a) == "r" || strings.ToLower(a) == "u") && err == nil {
			// First, Unschedule all flows
			for _, flow := range fres.Flows {
				// Fetch a schedule
				sreq := &dementor.FetchScheduleReq{
					ProjectId:  fres.ProjectId,
					FlowId:     flow.FlowId,
					CommonConf: commonConf,
				}
				sres, err := dementor.FetchSchedule(id, sreq)
				if err != nil {
					c.Ui.Error(err.Error())
					c.Ui.Error(fmt.Sprintf("Failed to fetch the schedule of [%s].", sreq.FlowId))
					c.askWhetherContinue()
					errorCount += 1
					continue
				}
				// Unschedule
				if sres != nil && sres.ScheduleId != "" {
					c.Ui.Warn(fmt.Sprintf("The schedule of [%s] will be removed...", flow.FlowId))
					ureq := &dementor.UnscheduleFlowReq{
						ScheduleId: sres.ScheduleId,
						CommonConf: commonConf,
					}
					err = dementor.UnscheduleFlow(id, ureq)
					if err != nil {
						c.Ui.Error(err.Error())
						c.Ui.Error(fmt.Sprintf("Failed to unschedule [%s].", flow.FlowId))
						c.askWhetherContinue()
						errorCount += 1
						continue
					}
					c.Ui.Info(fmt.Sprintf("Schedule of [%s] was removed!", flow.FlowId))
				}
			}

			if strings.ToLower(a) == "r" {
				// If Replace, delete the project and create
				dreq := &dementor.DeleteProjectReq{
					Project:    project.Name,
					CommonConf: commonConf,
				}
				err = dementor.DeleteProject(id, dreq)
				if err != nil {
					c.Ui.Error(err.Error())
					c.Ui.Error(fmt.Sprintf("Failed to remove the project [%s].", project.Name))
					return err
				}
				c.Ui.Info(fmt.Sprintf("The project [%s] was removed!\n", project.Name))

				// Create the project
				err = c.createProject(id, commonConf, project)
				if err != nil {
					c.Ui.Error(err.Error())
					c.Ui.Error(fmt.Sprintf("Failed to create the project [%s].", project.Name))
					return err
				}
				c.Ui.Info("Done!\n")
			}
		} else {
			c.Ui.Error("Terminated. Bye!")
			if err == nil {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
	}

	// Upload a zip file
	c.Ui.Output(fmt.Sprintf("The zip file [%s] will be uploaded...", project.FilePath))
	ureq := &dementor.UploadProjectZipReq{
		Project:    project.Name,
		FilePath:   project.FilePath,
		CommonConf: commonConf,
	}
	_, err = dementor.UploadProjectZip(id, ureq)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Error(fmt.Sprintf("Failed to upload the zip file [%s].", project.FilePath))
		return err
	}
	c.Ui.Info("Done!\n")

	// Process flows
	for _, flow := range project.Flows {
		err = c.processFlow(id, commonConf, project.Name, flow)
		if err != nil {
			c.Ui.Error(err.Error())
			c.Ui.Error(fmt.Sprintf("Failed to schedule the flow [%s].", flow.Name))
			c.Ui.Info("HINT: Please try replace instead of update.")
			c.askWhetherContinue()
			errorCount += 1
			continue
		}
	}

	if errorCount > 0 {
		c.Ui.Warn(fmt.Sprintf("The project [%s] was processed with %d errors.\n", project.Name, errorCount))
	} else {
		c.Ui.Info(fmt.Sprintf("The project [%s] was successfully processed!\n", project.Name))
	}
	return nil
}

func (c *KissCommand) processFlow(id string, commonConf dementor.CommonConf, project string, flow Flow) error {
	// Schedule a flow
	c.Ui.Output(fmt.Sprintf("The flow [%s] will be scheduled...", flow.Name))
	sreq := &dementor.ScheduleFlowReq{
		ProjectName:    project,
		Flow:           flow.Name,
		CronExpression: flow.Cron,
		CommonConf:     commonConf,
	}
	_, err := dementor.ScheduleFlow(id, sreq)
	if err != nil {
		return err
	}
	c.Ui.Info("Done!\n")
	return nil
}

func (c *KissCommand) Run(args []string) int {
	c.Ui.Output("Start to process the recipe...\n")
	_, err := kissParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	// Parse the recipe file
	buf, err := ioutil.ReadFile(kissOpts.FilePath)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	var recipe Recipe
	err = yaml.Unmarshal(buf, &recipe)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	commonConf := dementor.CommonConf{
		Url:      recipe.Url,
		UserName: recipe.Username,
		Password: recipe.Password,
	}

	// Get a session id
	id, err := getSessionId(commonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Process projects
	for _, project := range recipe.Projects {
		c.processProject(id, commonConf, project)
	}

	c.Ui.Info(fmt.Sprintf("This is the end of the recipe file [%s]. Bye!", kissOpts.FilePath))
	return 0
}

func (c *KissCommand) Help() string {
	buf := new(bytes.Buffer)
	kissParser.WriteHelp(buf)
	return buf.String()
}

func (c *KissCommand) Synopsis() string {
	return "Process a recipe file"
}
