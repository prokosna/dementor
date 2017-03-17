package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"github.com/antonholmquist/jason"
	"strings"
	"errors"
)

// Create project.
func UnscheduleFlow(uri string, sessionId string, projectName string, flowName string, cronExpression string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "scheduleCronFlow")
	q.Set("projectName", projectName)
	q.Set("flowName", flowName)
	q.Set("cronExpression", cronExpression)
	u.Path = strings.Trim(u.Path, "/") + "/schedule"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method: "POST",
		Uri: u.String(),
	}.Do()
	if err != nil {
		return err
	}

	defer res.Body.Close()

	v, err := jason.NewObjectFromReader(res.Body)
	if err != nil {
		return err
	}
	status, err := v.GetString("status")
	if err != nil {
		return err
	}
	if status != "success" {
		message, _ := v.GetString("message")
		return errors.New(message)
	}
	return nil
}
