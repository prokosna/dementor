package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"strings"
)

type DeleteProjectReq struct{
	Project string `json:"project"`
}

// Create project.
func DeleteProject(uri string, sessionId string, project string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("delete", "true")
	q.Set("project", project)
	u.Path = strings.Trim(u.Path, "/") + "/manager"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method: "GET",
		Uri: u.String(),
	}.Do()
	if err != nil {
		return err
	}

	res.Body.Close()
	return nil
}


