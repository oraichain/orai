package filehandling

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

// CreateMultipartFormData create a multiform data type
func CreateMultipartFormData(fieldName, fileName string) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		return bytes.Buffer{}, nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return bytes.Buffer{}, nil, err
	}
	w.Close()
	return b, w, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
