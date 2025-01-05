package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Color output constants.
const (
	Red    = "\033[31m"
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

// User-Agent header used in requests.
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0"

func Mode0(url string) {
	fmt.Println(Yellow+"[*] Testing URL:", url+Reset)
	fmt.Println(Yellow + "[*] Mode: 0 (Status Code)" + Reset)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	fmt.Println(Green+"[+] Response:", res.Status+Reset)
	fmt.Println(Green+"[+] Response URL:", res.Request.URL)
}

func Mode1(url string) {
	fmt.Println(Yellow+"[*] Testing URL:", url+Reset)
	fmt.Println(Yellow + "[*] Mode: 1 (Response Body)" + Reset)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("response.txt", body, os.ModePerm)
	fmt.Println(Green+"[+] Response:", res.Status+Reset)
	fmt.Println(Green+"[+] Response URL:", res.Request.URL)
	fmt.Println(Green + "[+] Saved response to response.txt" + Reset)
}

func Mode2(url string) {
	fmt.Println(Yellow+"[*] Testing URL:", url+Reset)
	fmt.Println(Yellow + "[*] Mode: 2 (Status Code Without Following Redirects)" + Reset)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	fmt.Println(Green+"[+] Response:", res.Status+Reset)
	fmt.Println(Green+"[+] Response URL:", res.Request.URL)
}

func Mode3(url string) {
	fmt.Println(Yellow+"[*] Testing URL:", url+Reset)
	fmt.Println(Yellow + "[*] Mode: 3 (Response Body Without Following Redirects)" + Reset)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("response.txt", body, os.ModePerm)
	fmt.Println(Green+"[+] Response:", res.Status+Reset)
	fmt.Println(Green+"[+] Response URL:", res.Request.URL)
	fmt.Println(Green + "[+] Saved response to response.txt" + Reset)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(Red + "Usage: gosearch <url> <mode>\nIssues: https://github.com/ibnaleem/gosearch/issues" + Reset)
		os.Exit(1)
	} else if len(os.Args) == 2 {
		fmt.Println(Red + "Mode not provided. Please provide either 0, 1, 2, or 3. Exiting..." + Reset)
		os.Exit(1)
	} else if len(os.Args) == 1 {
		fmt.Println(Yellow + "Welcome to GoSearch's testing binary." + Reset)
		fmt.Println(Yellow + "First, find a url containing a username." + Red + "Eg. https://instagram.com/zuck" + Reset)
		fmt.Println(Yellow + "Then, provide the mode number you want to test against." + Red + "Eg. ./test https://instagram.com/zuck 0" + Reset)
		fmt.Println(Yellow + "Modes:\n0: Status Code - Check if a website throws any status code errors for invalid usernames")
		fmt.Println(Yellow + "1: Response Body - Check if the response body contains any errors for invalid usernames (e.g 'username not found')")
		fmt.Println(Yellow + "2: Status Code (No Redirects) - Check if a website throws any status code errors for invalid usernames without following redirects")
		fmt.Println(Yellow + "3: Response Body (No Redirects) - Check if the response body contains any errors for invalid usernames (e.g 'username not found') without following redirects")
		os.Exit(1)
	}

	url := os.Args[1]
	mode := os.Args[2]

	if mode == "0" {
		Mode0(url)
	} else if mode == "1" {
		Mode1(url)
	} else if mode == "2" {
		Mode2(url)
	} else if mode == "3" {
		Mode3(url)
	} else {
		fmt.Println(Red + "Invalid mode. Please provide either 0, 1, 2, or 3. Exiting..." + Reset)
		os.Exit(1)
	}
}
