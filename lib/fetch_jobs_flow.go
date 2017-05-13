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

type FetchJobsFlowReq struct {
	Project string // The project name to be fetched
	Flow    string // The flow id to be fetched
	CommonConf
}

type FetchJobsFlowRes struct {
	Project   string `json:"project"`
	ProjectId string `json:"projectId"`
	Flow      string `json:"flow"`
	Nodes     []struct {
		Id   string   `json:"id"`
		Type string   `json:"type"`
		In   []string `json:"in"` // Job ids that this job is directory depending upon
	} `json:"nodes"`
}

// Fetch jobs of a flow
func FetchJobsFlow(sessionId string, fq *FetchJobsFlowReq) (*FetchJobsFlowRes, error) {
	u, err := url.Parse(fq.HTTP.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("session.id", sessionId)
	q.Set("ajax", "fetchflowgraph")
	q.Set("project", fq.Project)
	q.Set("flow", fq.Flow)
	u.Path = strings.Trim(u.Path, "/") + "/manager"
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
	var fs FetchJobsFlowRes
	err = gjson.Unmarshal([]byte(body), &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
