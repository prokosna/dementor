// +build unit

package dementor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadProjectZipSuccess(t *testing.T) {
	var handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, `{"error":"","projectId":"id","version":"1"}`)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	conf := Config
	conf.Url = ts.URL

	req := &UploadProjectZipReq{
		Project:    "TestForUploading",
		FilePath:   "../assets/test.zip",
		CommonConf: conf,
	}
	t.Logf("%+v", *req)
	res, err := UploadProjectZip("sessionId", req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", *res)
}
