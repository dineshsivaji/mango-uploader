package uploaders

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Host struct {
	Credentials
	InitUrl      string
	LoginUrl     string
	UploadUrl    string
	DownloadUrl  string
	FilePath     string
	LoginSuccess bool
	uploadParams map[string]string
	httpClient   *http.Client
}
type Credentials struct {
	Username string
	Password string
	ApiKey   string
}
type Uploader interface {
	Init()
	Login() bool
	ParsePage() bool
	UploadFile() bool
}

func MakeMultipartFormData(uri string, params map[string]string, paramName string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
