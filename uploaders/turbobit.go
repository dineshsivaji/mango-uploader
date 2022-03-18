package uploaders

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type TurboBitHost struct {
	Host
}
type TurboBitUploader struct {
	Uploader
}

func (u *TurboBitHost) Init() {
	u.httpClient = &http.Client{}
}

func (u *TurboBitHost) ParsePage() bool {

	resp, err := u.httpClient.Get(u.InitUrl)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return false
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
		return false
	}
	fmt.Println(fileDomain[1])
	appType := appTypePattern.FindStringSubmatch(sb)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	fmt.Println(appType[1])
	formActionURL := formActionPattern.FindStringSubmatch(sb)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(formActionURL[1])
	u.UploadUrl = formActionURL[1]
	u.uploadParams = map[string]string{
		"apptype": appType[1],
	}
	return true
}

func (u *TurboBitHost) UploadFile() bool {
	mulipartReq, err := MakeMultipartFormData(
		u.UploadUrl, u.uploadParams,
		"Filedata", u.FilePath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	uploadResp, err := u.httpClient.Do(mulipartReq)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(uploadResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		uploadResp.Body.Close()
		fmt.Println(uploadResp.StatusCode)
		fmt.Println(body)
		// log.Print(sb)
		return true
	}
}
