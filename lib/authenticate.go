package dementor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type AuthenticateReq struct {
	CommonConf
}

type AuthenticateRes struct {
	SessionId string `json:"session.id"`
	Status    string `json:"status"`
}

// Authenticate with username and password, then obtain session id.
func Authenticate(aq *AuthenticateReq) (*AuthenticateRes, error) {
	u, err := url.Parse(aq.Url)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("action", "login")
	q.Set("username", aq.UserName)
	q.Set("password", aq.Password)
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method:   "POST",
		Uri:      u.String(),
		Insecure: aq.Insecure,
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
		return nil, fmt.Errorf("ERROR: StatusCode is not 2xx: %d", res.StatusCode)
	}

	// check error
	// if "error" exists, that's an error
	if errMsg := gjson.Get(body, "error"); errMsg.Exists() {
		return nil, fmt.Errorf("ERROR: %s", errMsg)
	}

	// parse body
	var as AuthenticateRes
	err = gjson.Unmarshal([]byte(body), &as)
	if err != nil {
		return nil, err
	}
	if as.Status != "success" {
		return nil, fmt.Errorf("ERROR: Response status is not success: %+v", as)
	}

	return &as, nil
}
