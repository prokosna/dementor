// +build unit

package dementor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteProjectSuccess(t *testing.T) {
	var handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, "{}")
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	conf := Config
	conf.Url = ts.URL

	req := &DeleteProjectReq{
		Project:    "DeleteTest",
		CommonConf: conf,
	}
	t.Logf("%+v", *req)
	err := DeleteProject("sessionId", req)
	if err != nil {
		t.Fatal(err)
	}
}
