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

	values := url.Values{}
	values.Add("action", "login")
	values.Add("username", aq.UserName)
	values.Add("password", aq.Password)

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
