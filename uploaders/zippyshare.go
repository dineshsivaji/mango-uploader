package uploaders

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ZippyShareHost struct {
	Host
}
type ZippyShareUploader interface {
	Uploader
}

func (u *ZippyShareHost) Init() {
	cookieJar, _ := cookiejar.New(nil)
	u.httpClient = &http.Client{
		//Required cookie present only with login 302 response headers
		//Adding this check to ensure it is not auto-redirected to next page, thereby able to
		//access the auth cookie field sent in this response
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}
}

func (u *ZippyShareHost) Login() bool {

	if len(u.Credentials.Username) == 0 || len(u.Credentials.Password) == 0 {
		log.Println("Zippyshare account information is not available")
		return true
	}
	data := url.Values{}
	data.Set("login", u.Username)
	data.Set("pass", u.Password)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", u.LoginUrl, strings.NewReader(encodedData))
	if err != nil {
		log.Fatalln("Failed to create new post request")
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := u.httpClient.Do(req)
	if err != nil {
		log.Fatalln("Failed to login to zippyshare")
		return false
	}
	defer response.Body.Close()
	// Loop over header names
	// for name, values := range response.Header {
	// 	// Loop over all values for the name.
	// 	for _, value := range values {
	// 		fmt.Println(name, value)
	// 	}
	// }
	if response.StatusCode == 302 {
		reqParamsCount := 0
		_zipName, _zipHash := "", ""
		for _, cookie := range response.Cookies() {
			if cookie.Name == "zipname" {
				_zipName = cookie.Value
				reqParamsCount++
			} else if cookie.Name == "ziphash" {
				_zipHash = cookie.Value
				reqParamsCount++
			}
			if reqParamsCount == 2 {
				u.LoginSuccess = true
				break
			}
		}
		u.uploadParams = map[string]string{
			"private": "true", "embPlayerValues": "false",
			"zipname": _zipName,
			"ziphash": _zipHash,
		}
	}
	return u.LoginSuccess
}

func (u *ZippyShareHost) ParsePage() bool {

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
	// log.Println(sb)
	uploadUrlPattern := regexp.MustCompile(`document\.location\.protocol\+['"]//([^'"]+)`)
	// url:\s+document.location.protocol+'//
	uploadUrl := uploadUrlPattern.FindStringSubmatch(sb)
	u.UploadUrl = "https://" + uploadUrl[1]
	if !u.LoginSuccess {
		u.uploadParams = map[string]string{
			"private": "true", "embPlayerValues": "false",
			"zipname": "", "ziphash": "",
		}
	}
	return false
}

func (u *ZippyShareHost) UploadFile() bool {

	mulipartReq, err := MakeMultipartFormData(
		u.UploadUrl, u.uploadParams,
		"file", u.FilePath)
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
		respStr := body.String()
		fileName := filepath.Base(u.FilePath)
		downloadLinkPattern := regexp.MustCompile("href=['\"]([^\"']+)['\"]><b>" + fileName)
		fileLink := downloadLinkPattern.FindStringSubmatch(respStr)
		u.DownloadUrl = fileLink[1]
		fmt.Println(u.DownloadUrl)
		// log.Print(sb)
		return true
	}
}

// um := &uploaders.ZippyShareHost{
// 	Host: uploaders.Host{
// 		InitUrl:  "https://zippyshare.com//",
// 		LoginUrl: "https://www.zippyshare.com/services/login",
// 		FilePath: "/Users/dsivaji/Dinesh/Studies/Books/Programming/Learning Python, 5th Edition.pdf",
// 		Credentials: uploaders.Credentials{
// 			Username: "mangouploader",
// 			Password: "",
// 		},
// 	},
// }
