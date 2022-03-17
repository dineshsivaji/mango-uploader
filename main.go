package main

import (
	"example.com/mangouploader/uploaders"
)

func start(um uploaders.UploadManager) {
	um.Init()
	um.ParsePage()
	um.UploadFile()
}
func main() {

	um := &uploaders.Uploader{
		InitUrl:  "https://turbobit.net/",
		FilePath: "/Users/dsivaji/Downloads/housing.csv",
	}
	// um.Init()
	// um.ParsePage()
	// um.UploadFile()
	start(um)
}
