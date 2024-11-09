// Created by Ibn Aleem (github.com/ibnaleem)
// Updated Saturday 09 November, 2024 @ 23:50 GMT
// Repository: https://github.com/ibnaleem/gosearch
// Issues: https://github.com/ibnaleem/search/issues
// License: https://github.com/ibnaleem/gosearch/blob/main/LICENSE

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/inancgumus/screen"
	"gopkg.in/yaml.v3"
)

var Reset = "\033[0m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var ASCII string = `
 ________  ________  ________  _______   ________  ________  ________  ___  ___     
|\   ____\|\   __  \|\   ____\|\  ___ \ |\   __  \|\   __  \|\   ____\|\  \|\  \    
\ \  \___|\ \  \|\  \ \  \___|\ \   __/|\ \  \|\  \ \  \|\  \ \  \___|\ \  \\\  \   
 \ \  \  __\ \  \\\  \ \_____  \ \  \_|/_\ \   __  \ \   _  _\ \  \    \ \   __  \  
  \ \  \|\  \ \  \\\  \|____|\  \ \  \_|\ \ \  \ \  \ \  \\  \\ \  \____\ \  \ \  \ 
   \ \_______\ \_______\____\_\  \ \_______\ \__\ \__\ \__\\ _\\ \_______\ \__\ \__\
    \|_______|\|_______|\_________\|_______|\|__|\|__|\|__|\|__|\|_______|\|__|\|__|
                       \|_________|
`
var VERSION string = "v1.0.0"

type Website struct {
	Name      string `yaml:"name"`
	BaseURL   string `yaml:"base_url"`
	URLProbe  string `yaml:"url_probe,omitempty"`
	ErrorType string `yaml:"errorType"`
	ErrorMsg  string `yaml:"errorMsg,omitempty"`
	ErrorCode int    `yaml:"errorCode,omitempty"`
}

type Config struct {
	Websites []Website `yaml:"websites"`
}

func UnmarshalYAML() (Config, error) {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return config, nil
}

func WriteToFile(filename string, content string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
    	panic(err)
	}
	
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
    		panic(err)
	}
}

func BuildURL(baseURL, username string) string {
	return strings.Replace(baseURL, "{}", username, 1)
}

func MakeRequestWithoutErrorMsg(url string, WebsiteErrorCode int, wg *sync.WaitGroup) {
	defer wg.Done()

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", url, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != WebsiteErrorCode {
		fmt.Println(Green + "::", url + Reset)
		WriteToFile("results.txt", url + "\n")
	}
}

func MakeRequestWithErrorMsg(website Website, url string, errorMsg string, username string, wg *sync.WaitGroup) {
	defer wg.Done()

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", url, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, errorMsg) {
		if website.URLProbe != "" {
			url = BuildURL(website.BaseURL, username)
		}
		fmt.Println(Green + "::", url + Reset)
		WriteToFile("results.txt", url + "\n")
	}
}

func Search(config Config, username string) {
	var wg sync.WaitGroup
	var url string

	for _, website := range config.Websites {
		url = BuildURL(website.BaseURL, username)

		// Launch goroutines based on the ErrorType
		if website.ErrorType == "errorMsg" {
			if website.URLProbe != "" {
				url = BuildURL(website.URLProbe, username)
			}
			wg.Add(1)
			go MakeRequestWithErrorMsg(website, url, website.ErrorMsg, username, &wg)
		} else if website.ErrorType == "unknown" {
			fmt.Println(Yellow + ":: [?]", url + Reset)
			WriteToFile("results.txt", "[?] " + url + "\n")
		} else {
			wg.Add(1)
			go MakeRequestWithoutErrorMsg(url, website.ErrorCode, &wg)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <username>")
		return
	}
	var username string = os.Args[1]

	config, err := UnmarshalYAML()
	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return
	}

	screen.Clear()
	fmt.Println(ASCII)
	fmt.Println(VERSION)
	fmt.Println(strings.Repeat("⎯", 60))
	fmt.Println(":: Username                              : ", username)
	fmt.Println(":: Websites                              : ", len(config.Websites))
	fmt.Println(strings.Repeat("⎯", 60))
	fmt.Println(":: A yellow link indicates that I was unable to verify whether the username exists on the platform.")
	
	Search(config, username)
}