package dementor

import (
	"net/url"
	"github.com/franela/goreq"
	"bytes"
	"mime/multipart"
	"os"
	"path"
	"net/textproto"
	"fmt"
	"io"
	"strings"
)

// Create project.
func UploadProjectZip(uri string, sessionId string, project string, filePath string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	// prepare a from with zip file
	fileName := path.Base(filePath)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "application/zip")
	part.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	pw, err := w.CreatePart(part)
	if err != nil {
		return nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = io.Copy(pw, f); err != nil {
		return err
	}

	// add another field
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="session.id"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte(sessionId)); err != nil {
		return err
	}
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="ajax"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte("upload")); err != nil {
		return err
	}
	part = make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="project"`)
	pw, err = w.CreatePart(part)
	if _, err = pw.Write([]byte(project)); err != nil {
		return err
	}

	w.Close()

	u.Path = strings.Trim(u.Path, "/") + "/manager"

	res, err := goreq.Request{
		Method: "POST",
		Uri: u.String(),
		Body: &b,
		ContentType: w.FormDataContentType(),
	}.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
