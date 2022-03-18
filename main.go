package main

import (
	"example.com/mangouploader/uploaders"
)

func start(um uploaders.Uploader) {
	um.Init()
	um.ParsePage()
	um.UploadFile()
}

func main() {

	um := &uploaders.ZippyShareHost{
		Host: uploaders.Host{
			InitUrl:  "https://zippyshare.com//",
			LoginUrl: "https://www.zippyshare.com/services/login",
			FilePath: "/Users/dsivaji/Dinesh/Studies/Books/Programming/Learning Python, 5th Edition.pdf",
			Credentials: uploaders.Credentials{
				Username: "mangouploader",
				Password: "",
			},
		},
	}

	um.Init()
	um.Login()
	um.ParsePage()
	um.UploadFile()
}
