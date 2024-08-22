package api

import (
	"bytes"
	"mime/multipart"
)

// POST "/api/upload"
func (c *Connection) Upload(files []UploadFile) error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for _, file := range files {
		formWriter, err := w.CreateFormFile(file.Name, file.Name)
		if err != nil {
			return err
		}
		_, err = formWriter.Write(file.File)
		if err != nil {
			return err
		}
	}
	err := w.Close()
	if err != nil {
		return err
	}

	c.headers["Content-Type"] = w.FormDataContentType()
	_, err = c.Request(POST, "upload", []string{}, buf)
	c.headers["Content-Type"] = ""

	return err
}
