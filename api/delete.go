package api

import (
	"net/url"
	"strings"
)

// POST "/api/delete"
func (c *Connection) Delete(filePaths []string) error {
	var files []string
	for _, filePath := range filePaths {
		if filePath == "/index.html" {
			continue
		}
		files = append(files, filePath)

	}
	form := url.Values{
		"filenames[]": files,
	}
	reader := strings.NewReader(form.Encode())
	_, err := c.Request(POST, "delete", []string{}, reader)
	if err != nil {
		return err
	}
	// ignore /index.html
	//"filenames[]="+{filePaths}...
	return nil
}
