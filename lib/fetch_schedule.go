package dementor

import (
	"net/url"
	"strings"

	"fmt"
	"io"
	"io/ioutil"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type FetchScheduleReq struct {
	ProjectId string
	FlowId    string
	CommonConf
}

type FetchScheduleRes struct {
	CronExpression string `json:"cronExpression"`
	ScheduleId     string `json:"scheduleId"`
	SubmitUser     string `json:"submitUser"`
}

// Fetch a schedule of flow
func FetchSchedule(sessionId string, fq *FetchScheduleReq) (*FetchScheduleRes, error) {
	u, err := url.Parse(fq.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "fetchSchedule")
	q.Set("projectId", fq.ProjectId)
	q.Set("flowId", fq.FlowId)
	u.Path = strings.Trim(u.Path, "/") + "/schedule"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method: "GET",
		Uri:    u.String(),
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
	if errName := gjson.Get(body, "error"); errName.Exists() {
		return nil, fmt.Errorf("ERROR: %s", errName)
	}

	// check empty
	v := gjson.Parse(body)
	if !v.Exists() {
		return nil, fmt.Errorf("ERROR: Project or flow does not exist\nProjectId -> %s\nFlowId -> %s",
			fq.ProjectId,
			fq.FlowId)
	}

	// parse body
	var fs FetchScheduleRes
	err = gjson.Unmarshal([]byte(gjson.Get(body, "schedule").Raw), &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
