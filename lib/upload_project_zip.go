package dementor

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"

	"io/ioutil"

	"github.com/franela/goreq"
	"github.com/tidwall/gjson"
)

type UploadProjectZipReq struct {
	Project  string
	FilePath string
	CommonConf
}

type UploadProjectZipRes struct {
	Error     string `json:"error"`
	ProjectId string `json:"projectId"`
	Version   string `json:"version"`
}

// Upload a project zip file
func UploadProjectZip(sessionId string, uq *UploadProjectZipReq) (*UploadProjectZipRes, error) {
	u, err := url.Parse(uq.Url)
	if err != nil {
		return nil, err
	}

	// prepare a from with zip file
	fileName := path.Base(uq.FilePath)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "application/zip")
	part.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	pw, err := w.CreatePart(part)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(uq.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = io.Copy(pw, f); err != nil {
		return nil, err
	}

	// add another field
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="session.id"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte(sessionId)); err != nil {
		return nil, err
	}
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="ajax"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte("upload")); err != nil {
		return nil, err
	}
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="project"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte(uq.Project)); err != nil {
		return nil, err
	}
	w.Close()
	u.Path = strings.Trim(u.Path, "/") + "/manager"

	res, err := goreq.Request{
		Method:      "POST",
		Uri:         u.String(),
		Body:        &b,
		ContentType: w.FormDataContentType(),
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

	// parse body
	var us UploadProjectZipRes
	err = gjson.Unmarshal([]byte(body), &us)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Invalid upload_project_zip response\nResp -> %s", body)
	}
	if us.Error != "" {
		return nil, fmt.Errorf("ERROR: %+v", us)
	}

	return &us, nil
}
