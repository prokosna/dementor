package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"github.com/antonholmquist/jason"
	"strings"
	"errors"
)

type FetchFlowsProjectReq struct {
	Project     string
	Description string
}

type FetchFlowsProjectRes struct {
	Project   string `json:"project"`
	ProjectId string `json:"projectId"`
	Flows     []struct {
		FlowId string `json:"flowId"`
	} `json:"flows"`
}

// Create project.
func FetchFlowsProject(uri string, sessionId string, project string, description string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("action", "create")
	q.Set("name", name)
	q.Set("description", description)
	u.Path = strings.Trim(u.Path, "/") + "/manager"
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
