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
	Flow           string
	CronExpression string
	CommonConf
}

type ScheduleFlowRes struct {
	ScheduleId int64 `json:"scheduleId"`
}

// Create project.
func ScheduleFlow(sessionId string, sq *ScheduleFlowReq) (*ScheduleFlowRes, error) {
	u, err := url.Parse(sq.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "scheduleCronFlow")
	q.Set("projectName", sq.ProjectName)
	q.Set("flow", sq.Flow)
	q.Set("cronExpression", sq.CronExpression)
	u.Path = strings.Trim(u.Path, "/") + "/schedule"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "POST",
		Uri:      u.String(),
		Insecure: sq.Insecure,
	}.Do()
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	body, _ := res.Body.ToString()

	// check status
	if res.StatusCode < 200 || res.StatusCode > 399 {
		return nil, fmt.Errorf("ERROR: StatusCode is %d", res.StatusCode)
	}

	// check error
	if status := gjson.Get(body, "status"); status.Str != "success" {
		errMsg := gjson.Get(body, "message")
		return nil, fmt.Errorf("ERROR: %s\n%s", status, errMsg)
	}

	// parse body
	sr := &ScheduleFlowRes{
		ScheduleId: gjson.Get(body, "scheduleId").Int(),
	}

	return sr, nil
}
