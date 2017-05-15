package dementor

import (
	"net/url"
	"strings"

	"io"
	"io/ioutil"

	"fmt"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type UnscheduleFlowReq struct {
	ScheduleId string
	CommonConf
}

// Unschedule a flow
func UnscheduleFlow(sessionId string, req *UnscheduleFlowReq) error {
	u, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("action", "removeSched")
	q.Set("scheduleId", req.ScheduleId)
	u.Path = strings.Trim(u.Path, "/") + "/schedule"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "POST",
		Uri:      u.String(),
		Insecure: req.Insecure,
	}.Do()
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	body, _ := res.Body.ToString()

	// check status
	if res.StatusCode < 200 || res.StatusCode > 399 {
		return fmt.Errorf("ERROR: StatusCode is %d", res.StatusCode)
	}

	// check error
	if status := gjson.Get(body, "status"); status.Str != "success" {
		errMsg := gjson.Get(body, "message")
		return fmt.Errorf("ERROR: %s\n%s", status, errMsg)
	}

	return nil
}
