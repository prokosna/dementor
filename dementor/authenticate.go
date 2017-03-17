package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"encoding/json"
	"errors"
	"github.com/antonholmquist/jason"
	"fmt"
)

type AuthenticateReq struct {
	UserName string
	Password string
}

type AuthenticateRes struct {
	SessionId string `json:"session.id"`
	Status    string `json:"status"`
}

// Authenticate with username and password, then obtain session id.
func Authenticate(uri string, ap AuthenticateReq) (AuthenticateRes, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("action", "login")
	q.Set("username", ap.UserName)
	q.Set("password", ap.Password)
	u.RawQuery = q.Encode()

	res, err := goreq.Request{
		Method: "POST",
		Uri: u.String(),
		Insecure: true,
	}.Do()
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := res.Body.ToString()

	// check error
	// if "error" exists, that's an error
	v, err := jason.NewObjectFromBytes([]byte(body))
	if err != nil {
		return nil, err
	}
	if msg, err := v.GetString("error"); err == nil {
		return nil, errors.New(msg)
	}

	var ars AuthenticateRes
	err = json.Unmarshal([]byte(body), &ars)
	if err != nil {
		return nil, err
	}
	if ars.Status != "success" {
		return nil, fmt.Errorf("ERROR: Response status is not success.: %s", ars.Status)
	}

	return ars, nil
}
