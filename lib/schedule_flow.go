package dementor

import (
	"fmt"
	"net/url"
	"strings"

	"io"
	"io/ioutil"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type ScheduleFlowReq struct {
	ProjectName    string
	FlowName       string
	CronExpression string
	CommonConf
}

// Create project.
func ScheduleFlow(sessionId string, sq *ScheduleFlowReq) error {
	u, err := url.Parse(sq.HTTP.Url)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "scheduleCronFlow")
	q.Set("projectName", sq.ProjectName)
	q.Set("flowName", sq.FlowName)
	q.Set("cronExpression", sq.CronExpression)
	u.Path = strings.Trim(u.Path, "/") + "/schedule"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "POST",
		Uri:      u.String(),
		Insecure: sq.HTTP.Insecure,
	}.Do()
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	body, _ := res.Body.ToString()

	// check error
	if status := gjson.Get(body, "status"); status.Str != "success" {
		errMsg := gjson.Get(body, "message")
		return fmt.Errorf("ERROR: %s\n%s", status, errMsg)
	}

	return nil
}
