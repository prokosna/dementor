package command

import "github.com/prokosna/dementor/lib"

func getSessionId(config dementor.CommonConf) (string, error) {
	areq := &dementor.AuthenticateReq{
		CommonConf: config,
	}
	ares, err := dementor.Authenticate(areq)
	if err != nil {
		return "", err
	}
	return ares.SessionId, nil
}
