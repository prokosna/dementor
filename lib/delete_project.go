package dementor

import (
	"net/url"
	"strings"

	"io"
	"io/ioutil"

	"fmt"

	"github.com/franela/goreq"
)

type DeleteProjectReq struct {
	Project string
	CommonConf
}

// Delete a project
func DeleteProject(sessionId string, dq *DeleteProjectReq) error {
	u, err := url.Parse(dq.Url)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("delete", "true")
	q.Set("project", dq.Project)
	u.Path = strings.Trim(u.Path, "/") + "/manager"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "GET",
		Uri:      u.String(),
		Insecure: dq.Insecure,
	}.Do()
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	// check status
	if res.StatusCode < 200 || res.StatusCode > 399 {
		return fmt.Errorf("ERROR: StatusCode is not 2xx: %d", res.StatusCode)
	}

	return nil
}
