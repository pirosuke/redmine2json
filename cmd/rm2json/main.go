package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/antonholmquist/jason"
	redmine "github.com/pirosuke/rm2json"
)

/*
AppConfigDirName is the directory name to save settings.
*/
const AppConfigDirName = "redmine2json"

/*
AppConfig describes setting format.
*/
type AppConfig struct {
	Redmine struct {
		APIKey    string `json:"api_key"`
		URLRoot   string `json:"url_root"`
		BasicAuth struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"basic_auth"`
	} `json:"redmine"`
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func main() {
	configName := flag.String("c", "", "Config Name")
	outputDirPath := flag.String("o", ".", "Output Directory Path")
	flag.Parse()

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !fileExists(*outputDirPath) {
		fmt.Println("output dir does not exist: " + *outputDirPath)
		return
	}

	appConfigDirPath := filepath.Join(userConfigDir, AppConfigDirName)
	os.MkdirAll(appConfigDirPath, os.ModePerm)

	appConfigFilePath := filepath.Join(appConfigDirPath, *configName+".json")
	if !fileExists(appConfigFilePath) {
		fmt.Println("config file does not exist: " + appConfigFilePath)
		return
	}

	jsonContent, err := ioutil.ReadFile(appConfigFilePath)
	if err != nil {
		fmt.Println("failed to read config file: " + appConfigFilePath)
		return
	}

	appConfig := new(AppConfig)
	if err := json.Unmarshal(jsonContent, appConfig); err != nil {
		fmt.Println("failed to read config file: " + appConfigFilePath)
		return
	}

	auth := redmine.BasicAuth{
		UserName: appConfig.Redmine.BasicAuth.Username,
		Password: appConfig.Redmine.BasicAuth.Password,
	}

	limit := 100
	offset := 0
	for offset < 10000 {
		params := redmine.TicketFetchParams{
			Offset: offset,
			Limit:  limit,
		}

		result, _ := redmine.FetchTickets(appConfig.Redmine.URLRoot, appConfig.Redmine.APIKey, params, auth)
		resultJSON, _ := jason.NewObjectFromBytes(result)

		if offset == 0 {
			totalCount, _ := resultJSON.GetInt64("total_count")
			fmt.Println("total_count:", totalCount)
		}

		issueList, _ := resultJSON.GetObjectArray("issues")
		issueCount := len(issueList)
		fmt.Println("result count:", issueCount)
		if issueCount == 0 {
			break
		}

		for _, issue := range issueList {
			issueID, _ := issue.GetInt64("id")

			var buf bytes.Buffer
			err = json.Indent(&buf, []byte(issue.String()), "", "  ")
			if err != nil {
				fmt.Println(err)
				break
			}

			outputFilePath := filepath.Join(*outputDirPath, strconv.FormatInt(issueID, 10)+".json")
			ioutil.WriteFile(outputFilePath, buf.Bytes(), 0644)
		}

		time.Sleep(1 * time.Second)
		offset += limit
	}
}
