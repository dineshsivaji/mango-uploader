package userconfig

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var siteConfigFile string = "test_account.json"

// type siteAccounts struct {
// 	Username string `json:"username`
// 	password string `json:"password`
// }
// type siteAccount struct {
// 	Sites []siteAccount `json:"sites`
// }

var SupportedSites = []string{
	"turbobit", "zippshare", "filefactory", "uploaded", "radpigator", "depositfiles",
}

func HandleSiteAccountConfig() {
	// Ref : https://stackoverflow.com/a/12518877
	if _, err := os.Stat(siteConfigFile); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist

	}
	selectedSite := SupportedSites[getUserInput()]
	fmt.Println(selectedSite)
}
func getUserInput() int {
	sort.Strings(SupportedSites)
	reader := bufio.NewReader(os.Stdin)
	userChoice := -1
	for {
		fmt.Println("Please select the site to which you want to add credentials")
		for i, site := range SupportedSites {
			fmt.Printf(" %d. %s\n", i+1, site)
		}
		fmt.Println("0. Exit")
		fmt.Print("Your choice : ")
		choice, err := reader.ReadString('\n')
		choice = strings.TrimSuffix(choice, "\n")
		if err != nil {
			fmt.Println("Please give a valid input")
			continue
		}
		userChoice, err = strconv.Atoi(choice)
		if err != nil {
			fmt.Println("Please give a valid input")
			continue
		}

		if userChoice == 0 {
			os.Exit(0)
		} else if userChoice < 0 || userChoice > len(SupportedSites) {
			fmt.Println("Please give a valid input")
			continue
		} else if userChoice != -1 {
			break
		}
	}
	// Returning after subtracting with 1, as slice/array index starts at 0
	return userChoice - 1
}
func GetSiteInfo() {

	var result map[string]interface{}
	// Open our jsonFile
	jsonFile, err := os.Open("test_account1.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened test_account.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// decoder := json.NewDecoder(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// var result map[string]interface{}

	// Read the array open bracket
	// decoder.Token()

	// for decoder.More() {
	// 	decoder.Decode(&result)
	// 	fmt.Println(result)
	// }

	json.Unmarshal([]byte(byteValue), &result)

	for site, cred := range result {
		// cred := result[site].(map[string]interface{})
		myCred := cred.(map[string]interface{})
		fmt.Println(site, myCred["username"])
		// for k, v := range myCred {
		// 	fmt.Println(site, k, v)
		// }
	}

}

// Ref : https://stackoverflow.com/a/32768479
func getCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
