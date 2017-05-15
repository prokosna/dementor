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

type FetchFlowsProjectReq struct {
	Project string
	CommonConf
}

type FetchFlowsProjectRes struct {
	Project   string `json:"project"`
	ProjectId string `json:"projectId"`
	Flows     []struct {
		FlowId string `json:"flowId"`
	} `json:"flows"`
}

// Fetch flows of a project
func FetchFlowsProject(sessionId string, fq *FetchFlowsProjectReq) (*FetchFlowsProjectRes, error) {
	u, err := url.Parse(fq.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "fetchprojectflows")
	q.Set("project", fq.Project)
	u.Path = strings.Trim(u.Path, "/") + "/manager"
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "GET",
		Uri:      u.String(),
		Insecure: fq.Insecure,
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

	// parse body
	var fs FetchFlowsProjectRes
	err = gjson.Unmarshal([]byte(body), &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
