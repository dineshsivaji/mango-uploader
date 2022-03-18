package uploaders

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UploadedDotNetHost struct {
	Host
}
type UploadedDotNetUploader interface {
	Uploader
}

func GenerateUploadedDotNetFileID(length int) string {
	genID := strings.Builder{}
	rand.Seed(time.Now().UnixNano())
	consonants := [20]string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w", "x", "y", "z"}
	vowels := [5]string{"a", "e", "i", "o", "u"}
	for i := 0; i < length/2; i++ {

		c := rand.Int() % 20
		v := rand.Int() % 5
		genID.WriteString(consonants[c])
		genID.WriteString(vowels[v])
	}
	return genID.String()
}

func (u *UploadedDotNetHost) Init() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("Failed to initalize cookie jar for uploded.net")
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

func (u *UploadedDotNetHost) Login() bool {
	if len(u.Credentials.Username) == 0 || len(u.Credentials.Password) == 0 {
		log.Println("Uploaded.net account information is not available")
		return false
	}

	data := url.Values{}
	data.Set("id", u.Username)
	data.Set("pw", u.Password)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", u.LoginUrl, strings.NewReader(encodedData))
	if err != nil {
		log.Fatalln("Failed to create new post request")
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	response, err := u.httpClient.Do(req)
	if err != nil {
		log.Fatalln("Failed to login to uploaded.net : ", err)
		return false
	}
	if len(response.Cookies()) != 0 {
		u.LoginSuccess = true
	}

	return u.LoginSuccess
}

func (u *UploadedDotNetHost) ParsePage() bool {

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
	// http://uploaded.net/js/script.js
	hashPattern := regexp.MustCompile(`id=['"]user_pw['"]\s+value=['"]([^'"]+)`)
	hashMatch := hashPattern.FindStringSubmatch(sb)
	if len(hashMatch) != 2 {
		log.Fatalln("Failed to match password hash in uploaded.net page ")
		return false
	}
	hashStr := hashMatch[1]

	scriptJsUrl := "http://uploaded.net/js/script.js"
	scriptJSResp, _ := u.httpClient.Get(scriptJsUrl)
	//We Read the response body on the line below.
	scriptJSBody, err := ioutil.ReadAll(scriptJSResp.Body)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	scriptJSStr := string(scriptJSBody)
	uploadServerPattern := regexp.MustCompile(`uploadServer\s+=\s+['"]([^'"]+)['"]`)
	uploadServerMatch := uploadServerPattern.FindStringSubmatch(scriptJSStr)
	if len(uploadServerMatch) != 2 {
		log.Fatalln("Failed to match upload server info in uploaded.net page")
		return false
	}
	//Expected form : http://am4-r1f2-stor05.uploaded.net/upload?admincode=damafi&id=17899420&pw=adcb0b4c9e95e75371936f1bd2ec50ee589a860b

	uploadServer := uploadServerMatch[1] + "upload"
	fileID := GenerateUploadedDotNetFileID(6)
	// Ref : https://stackoverflow.com/a/56985985
	//Adding base URL alone first
	baseUrl, _ := url.Parse(uploadServer)
	//Adding the query params
	data := url.Values{}
	data.Set("admincode", fileID)
	data.Set("id", u.Username)
	data.Set("pw", hashStr)
	data.Set("folder", "0")
	fileName := filepath.Base(u.FilePath)
	u.uploadParams = map[string]string{
		"Filename": fileName,
	}
	baseUrl.RawQuery = data.Encode()
	//Converting to string with query params encoded
	u.UploadUrl = baseUrl.String()
	return true
}
func (u *UploadedDotNetHost) UploadFile() bool {

	if !u.LoginSuccess {
		log.Fatalln("You need to login before uploading a file into uploaded.net")
		return false
	}
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
		// for name, values := range uploadResp.Header {
		// 	// Loop over all values for the name.
		// 	for _, value := range values {
		// 		fmt.Println(name, value)
		// 	}
		// }
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(uploadResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(uploadResp.StatusCode)
		respStr := body.String()
		u.DownloadUrl = "http://ul.to/" + strings.Split(respStr, ",")[0]
		fmt.Println(u.DownloadUrl)
		defer uploadResp.Body.Close()
		return true
	}
}

// um := uploaders.UploadedDotNetHost{
// 	Host: uploaders.Host{
// 		Credentials: uploaders.Credentials{
// 			Username: "17899420",
// 			Password: "",
// 		},
// 		LoginUrl: "http://uploaded.net/io/login",
// 		InitUrl:  "http://uploaded.net/me",
// 		FilePath: "/Users/dsivaji/Downloads/housing.csv",
// 		// FilePath: "/Users/dsivaji/get-pip.py",
// 	},
// }
