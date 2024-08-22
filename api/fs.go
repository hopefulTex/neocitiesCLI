package api

import (
	"os"
	"strings"
)

type UploadFile struct {
	Name string
	File []byte
}

func OpenFiles(root string, filePaths []string) ([]UploadFile, []error) {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}
	// multiple files in single request
	var errs []error
	var files []UploadFile
	var tmpFile UploadFile
	for i, filePath := range filePaths {
		if !strings.HasPrefix(filePath, root) {
			filePath = root + filePath
		}
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		tmpFile = UploadFile{
			Name: filePaths[i],
			File: bytes,
		}
		files = append(files, tmpFile)
	}
	return files, errs
}
