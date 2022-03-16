package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
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

func main() {
	client := &http.Client{}

	resp, err := client.Get("https://turbobit.net/")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// for _, v := range resp.Cookies() {
	// 	fmt.Println(v)
	// }
	//Convert the body to type string
	sb := string(body)
	fileDomainPattern := regexp.MustCompile(`downloadFileDomain:\s+'([^']+)'`)
	formActionPattern := regexp.MustCompile(`form\s+action=["']([^"']+)"`)
	appTypePattern := regexp.MustCompile(`input\s+name=["']apptype["']\s+value=["']([^"]+)`)
	fileDomain := fileDomainPattern.FindStringSubmatch(sb)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(fileDomain[1])
	appType := appTypePattern.FindStringSubmatch(sb)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(appType[1])
	formActionURL := formActionPattern.FindStringSubmatch(sb)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(formActionURL[1])
	extraParams := map[string]string{
		"apptype": appType[1],
	}
	req, err := newfileUploadRequest(formActionURL[1], extraParams, "Filedata", "/Users/dsivaji/Downloads/housing.csv")
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(body)
		// log.Print(sb)
	}
}
