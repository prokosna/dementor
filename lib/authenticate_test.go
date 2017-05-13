// +build unit

package dementor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	var handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, `{"session.id":"test","status":"success"}`)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	conf := Config
	conf.HTTP.Url = ts.URL

	req := &AuthenticateReq{
		CommonConf: conf,
	}
	res, err := Authenticate(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.SessionId == "" {
		t.Fatal("No session id")
	}
	t.Logf("%+v", *res)
}
