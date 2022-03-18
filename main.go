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

}
