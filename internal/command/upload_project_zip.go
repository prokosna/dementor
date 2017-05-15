package command

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/prokosna/dementor/lib"
)

type UploadProjectZipCommand struct {
	Ui cli.Ui
}

var uploadProjectZipOpts struct {
	Project  string `short:"p" long:"project" description:"Project name" required:"true"`
	FilePath string `short:"f" long:"filepath" description:"Path to a project zip file" required:"true"`
	dementor.CommonConf
}

var uploadProjectZipParser *flags.Parser

func init() {
	uploadProjectZipParser = flags.NewParser(&uploadProjectZipOpts, flags.IgnoreUnknown|flags.HelpFlag|flags.PassDoubleDash)
}

func (c *UploadProjectZipCommand) Run(args []string) int {
	_, err := uploadProjectZipParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Output(c.Help())
		return 1
	}

	id, err := getSessionId(uploadProjectZipOpts.CommonConf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	req := &dementor.UploadProjectZipReq{
		Project:    uploadProjectZipOpts.Project,
		FilePath:   uploadProjectZipOpts.FilePath,
		CommonConf: uploadProjectZipOpts.CommonConf,
	}
	res, err := dementor.UploadProjectZip(
		id,
		req,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Successfully uploaded the file [%s]", res.ProjectId))
	return 0
}

func (c *UploadProjectZipCommand) Help() string {
	buf := new(bytes.Buffer)
	uploadProjectZipParser.WriteHelp(buf)
	return buf.String()
}

func (c *UploadProjectZipCommand) Synopsis() string {
	return "Upload the project zip file."
}
