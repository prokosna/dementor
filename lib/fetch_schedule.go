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
	u, err := url.Parse(fq.HTTP.Url)
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
		Method:   "GET",
		Uri:      u.String(),
		Insecure: fq.HTTP.Insecure,
	}.Do()
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	body, _ := res.Body.ToString()

	// check error
	if errName := gjson.Get(body, "error"); errName.Exists() {
		errMsg := gjson.Get(body, "message")
		return nil, fmt.Errorf("ERROR: %s\n%s", errName, errMsg)
	}

	// parse body
	var fs FetchScheduleRes
	err = gjson.Unmarshal([]byte(gjson.Get(body, "schedule").Raw), &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
