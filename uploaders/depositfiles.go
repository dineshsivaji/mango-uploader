package uploaders

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var loginURL string = "https://depositfiles.com/api/user/login"

type DepositFilesHost struct {
	Host
}
type DepositFilesUploader interface {
	Uploader
}

func (u *DepositFilesHost) Init() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("Failed to initalize cookie jar for depositfiles")
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

func (u *DepositFilesHost) ParsePage() bool {
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
	fmt.Println(sb)

	uploadServerPattern := regexp.MustCompile(`var\s+dfUploadPath\s+=\s+["']([^"']+)`)
	uploadServerMatch := uploadServerPattern.FindStringSubmatch(sb)
	if len(uploadServerMatch) != 2 {
		log.Fatalln("Failed to match password hash in depostfiles page ")
		return false
	}
	uploadServerStr := uploadServerMatch[1]
	u.UploadUrl = uploadServerStr

	return true
}

func (u *DepositFilesHost) Login() bool {

	data := url.Values{}
	data.Set("login", u.Username)
	data.Set("password", u.Password)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", u.LoginUrl, strings.NewReader(encodedData))
	if err != nil {
		log.Fatalln("Failed to create new post request")
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	uploadResp, err := u.httpClient.Do(req)
	if err != nil {
		log.Fatalln("Failed to login to uploaded.net : ", err)
		return false
	}
	fmt.Println(uploadResp.Cookies())
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(uploadResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	uploadResp.Body.Close()
	fmt.Println(uploadResp.StatusCode)
	respStr := body.String()
	fmt.Println(respStr)
	// "member_passkey":"phw1912zwsowgdgf"
	memberKeyPattern := regexp.MustCompile(`member_passkey"\s*:\s*"([^"]+)`)
	memberKeyMatch := memberKeyPattern.FindStringSubmatch(respStr)
	if len(memberKeyMatch) != 2 {
		u.LoginSuccess = true
		log.Println("Failed to login to depositfiles. Check credentials")
		return false
	}
	memberKey := memberKeyMatch[1]
	u.uploadParams = map[string]string{
		"member_passkey": memberKey,
		"format":         "html5",
		"fm":             "_root",
	}
	u.ParsePage()
	return false
}

func preUploadAction(u DepositFilesHost) {
	preUploadUrl := "https://depositfiles.com/api/upload/regular"
	resp, err := u.httpClient.Get(preUploadUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	// for _, v := range resp.Cookies() {
	// 	fmt.Println(v)
	// }
	//Convert the body to type string
	sb := string(body)
	fmt.Println(sb)
}
func (u *DepositFilesHost) UploadFile() bool {
	if !u.LoginSuccess {
		log.Println("Depositfiles allows file upload only after login")
		return false
	}
	preUploadAction(*u)
	mulipartReq, err := MakeMultipartFormData(
		u.UploadUrl, u.uploadParams,
		"files", u.FilePath)
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
		fmt.Println(respStr)
		// u.DownloadUrl = "https://www.filefactory.com/file/" + respStr + "/" + fileName
		// log.Print(u.DownloadUrl)
		return true
	}
}

// u := uploaders.DepositFilesHost{
// 	Host: uploaders.Host{
// 		Credentials: uploaders.Credentials{
// 			Username: "vacela3529@songsign.com",
// 			Password: "",
// 		},
// 		LoginUrl: "https://depositfiles.com/api/user/login",
// 		InitUrl:  "https://depositfiles.com/",
// 		FilePath: "/Users/dsivaji/Downloads/housing.csv",
// 	},
// }
// u.Init()
// u.Login()
// u.UploadFile()
