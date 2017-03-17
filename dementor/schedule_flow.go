package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"github.com/antonholmquist/jason"
	"strings"
	"errors"
	"fmt"
)

// Create project.
func ScheduleFlow(uri string, sessionId string, projectName string, flowName string, cronExpression string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "scheduleFetch")
	q.Set("projectName", projectName)
	q.Set("flow", flowName)
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
	body, _ := res.Body.ToString()
	println(body)
	v, err := jason.NewObjectFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	status, err := v.GetString("status")
	if err != nil {
		println(err)
		return err
	}
	if status != "success" {
		message, _ := v.GetString("message")
		return errors.New(message)
	}
	println("hoge")
	return nil
}
