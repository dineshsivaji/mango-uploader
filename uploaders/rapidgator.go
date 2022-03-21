package uploaders

import (
	"log"
	"net/http"
	"net/http/cookiejar"
)

type RapidGatorHost struct {
	Host
}
type RapidGatorUploader interface {
	Uploader
}

func (u *RapidGatorHost) Init() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("Failed to initalize cookie jar for rapidgator")
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
