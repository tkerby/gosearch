package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"bufio"
	"strings"
	"net/http"
	"crypto/tls"
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
	fmt.Println(Green+"[+] Response URL:", res.Request.URL.String() + Reset)
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
	fmt.Println(Green+"[+] Response URL:", res.Request.URL.String() + Reset)
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
	fmt.Println(Green+"[+] Response URL:", res.Request.URL.String() + Reset)
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
	fmt.Println(Green+"[+] Response URL:", res.Request.URL.String() + Reset)
	fmt.Println(Green + "[+] Saved response to response.txt" + Reset)
}

func Mode4(url string, errorMsg string) {
	fmt.Println(Yellow+"[*] Testing URL:", url+Reset)
	fmt.Println(Yellow+"[*] Testing error message:", errorMsg+Reset)
	fmt.Println(Yellow + "[*] Mode: 4 (Error Message Check)" + Reset)

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

	bodyStr := string(body)
	fmt.Println(Green+"[+] Response:", res.Status+Reset)
	fmt.Println(Green+"[+] Response URL:", res.Request.URL.String() + Reset)

	if strings.Contains(bodyStr, errorMsg) {
		  fmt.Println(Green+"[+] Error message found in response body: " + errorMsg + "\n[+] This means if a profile does not exist on %s", url, "I can detect it!" + Reset)
    } else {
      fmt.Println(Red+"[-] Error message not found in response body: " + errorMsg + "\n[-] This means if a profile does not exist on %s", url, "I CANNOT detect it!" + Reset)
    }
	}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(Yellow + "Welcome to GoSearch's testing binary." + Reset)
		fmt.Println(Yellow + "First, find a url containing a username." + Red + "Eg. https://instagram.com/zuck" + Reset)
		fmt.Println(Yellow + "Then, provide the mode number you want to test against." + Red + "Eg. ./test https://instagram.com/zuck 0" + Reset)
		fmt.Println(Yellow + "Modes:\n0: Status Code - Manually check if a website throws any status code errors for invalid usernames")
		fmt.Println(Yellow + "1: Response Body - Manually check if the response body contains any errors for invalid usernames (e.g 'username not found')")
		fmt.Println(Yellow + "2: Status Code (No Redirects) - Manually check if a website throws any status code errors for invalid usernames without following redirects")
		fmt.Println(Yellow + "3: Response Body (No Redirects) - Manually check if the response body contains any errors for invalid usernames (e.g 'username not found') without following redirects")
		fmt.Println(Yellow + "4: Error Message Detection - Actively test for and attempt to find any specific error messages in the response body for invalid usernames (e.g. 'user not found' or similar).")		
		os.Exit(1)
	} else if len(os.Args) == 2 {
		fmt.Println(Red + "Mode not provided. Please provide either 0, 1, 2, or 3. Exiting..." + Reset)
		os.Exit(1)
	} else if len(os.Args) > 3 {
		fmt.Println(Red + "Usage: gosearch <url> <mode>\nIssues: https://github.com/ibnaleem/gosearch/issues" + Reset)
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
	} else if mode == "4" {

		scanner := bufio.NewScanner(os.Stdin)

		var errorMsg string
		fmt.Print(Yellow + "[*] Please provide an error message found in the response body for me to check if I can detect it for invalid usernames: " + Reset)
		
		if scanner.Scan() {
			errorMsg = scanner.Text()
			Mode4(url, errorMsg)
		}

		if err:= scanner.Err(); err != nil {
			fmt.Println(Red + "Error reading input: " + err.Error() + Reset)
			os.Exit(1)
		}

	} else {
		fmt.Println(Red + "Invalid mode. Please provide either 0, 1, 2, 3 or 4. Exiting..." + Reset)
		os.Exit(1)
	}
}
