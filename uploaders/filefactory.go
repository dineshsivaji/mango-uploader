package uploaders

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type FileFactoryHost struct {
	Host
}
type FileFactoryHostUploader interface {
	Uploader
}

func (u *FileFactoryHost) Init() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("Failed to initalize cookie jar for filefactory")
		return
	}
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
func (u *FileFactoryHost) ParsePage() bool {
	//ParsePage is not applicable for FileFactory, so just returning
	//true with no further action
	return true
}

func (u *FileFactoryHost) Login() bool {
	if len(u.Credentials.Username) == 0 || len(u.Credentials.Password) == 0 {
		log.Println("Filefactory account information is not available")
		return true
	}

	data := url.Values{}
	data.Set("loginEmail", u.Username)
	data.Set("loginPassword", u.Password)
	data.Set("Submit", "Sign In")
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", u.LoginUrl, strings.NewReader(encodedData))
	if err != nil {
		log.Fatalln("Failed to create new post request")
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Accept", "*/*")
	// req.Header.Add("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	response, err := u.httpClient.Do(req)
	if err != nil {
		log.Fatalln("Failed to login to filefactory : ", err)
		return false
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "auth" {
			decodedCookie, _ := url.QueryUnescape(cookie.Value)
			u.uploadParams = map[string]string{
				"cookie": decodedCookie,
			}
			// fmt.Println(decodedCookie)
		}
	}
	return true
}
func (u *FileFactoryHost) UploadFile() bool {
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
		fileName := filepath.Base(u.FilePath)
		uploadResp.Body.Close()
		fmt.Println(uploadResp.StatusCode)
		respStr := body.String()
		u.DownloadUrl = "https://www.filefactory.com/file/" + respStr + "/" + fileName
		log.Print(u.DownloadUrl)
		return true
	}
}

// um := &uploaders.FileFactoryHost{
// 	Host: uploaders.Host{
// 		InitUrl:   "https://www.filefactory.com/member/signin.php",
// 		LoginUrl:  "https://www.filefactory.com/member/signin.php",
// 		UploadUrl: "https://upload.filefactory.com/upload",
// 		Credentials: uploaders.Credentials{
// 			Username: "herexo6590@superyp.com",
// 			Password: "",
// 		},
// 		FilePath: "/Users/dsivaji/Downloads/housing.csv",
// 	},
// }
