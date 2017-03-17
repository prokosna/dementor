package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"github.com/antonholmquist/jason"
	"strings"
	"errors"
)

type CreateProjectReq struct {
	Name        string
	Description string
}

type CreateProjectRes struct {
	Status string `json:"status"`
	Path   string `json:"path"`
	Action string `json:"action"`
}

// Create project.
func CreateProject(uri string, sessionId string, cq CreateProjectReq) (CreateProjectRes, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("action", "create")
	q.Set("name", cq.Name)
	q.Set("description", cq.Description)
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
