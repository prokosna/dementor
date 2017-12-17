package dementor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type CreateProjectReq struct {
	Name        string
	Description string
	CommonConf
}

type CreateProjectRes struct {
	Status string `json:"status"`
	Path   string `json:"path"`
	Action string `json:"action"`
}

// Create a project
func CreateProject(sessionId string, cq *CreateProjectReq) (*CreateProjectRes, error) {
	u, err := url.Parse(cq.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("action", "create")

	u.Path = strings.Trim(u.Path, "/") + "/manager"
	u.RawQuery = q.Encode()

	values := url.Values{}
	values.Add("session.id", sessionId)
	values.Add("name", cq.Name)
	values.Add("description", cq.Description)

	res, err := goreq.Request{
		Method:      "POST",
		Uri:         u.String(),
		Body:        values.Encode(),
		ContentType: "application/x-www-form-urlencoded",
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
		errMsg := gjson.Get(body, "message")
		return nil, fmt.Errorf("ERROR: %s\n%s", errName, errMsg)
	}

	// parse body
	var cs CreateProjectRes
	err = gjson.Unmarshal([]byte(body), &cs)
	if err != nil {
		return nil, err
	}
	if cs.Status != "success" {
		return nil, fmt.Errorf("ERROR: Response status is not success: %+v", cs)
	}

	return &cs, nil
}
