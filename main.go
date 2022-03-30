package main

import (
	"flag"
	"fmt"

	"example.com/mangouploader/uploaders"
	"example.com/mangouploader/userconfig"
)

func start(um uploaders.Uploader) {
	um.Init()
	um.ParsePage()
	um.UploadFile()
}

func main() {
	var loginConfig bool

	flag.BoolVar(&loginConfig, "login", false, "login config generator for sites")
	// textPtr := flag.String("text", "", "Text to parse.")
	// metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};.")
	// uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
	flag.Parse()
	fmt.Println(loginConfig)
	// userconfig.GetSiteInfo()
	userconfig.HandleSiteAccountConfig()
	// userconfig.Manage()
}
