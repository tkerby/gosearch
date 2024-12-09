<p align='center'>
<img src='img/gosearch-logo.png' height=50% width=50%><br>
<i>This project heavily relies on contributors, please see <a href="#contributing">Contributing</a> for more details.</i><br>
<code>git clone https://github.com/ibnaleem/gosearch.git && cd gosearch && go build && ./gosearch [username]</code>
</p>

<p align="center">
  <img src="https://github.com/ibnaleem/gosearch/actions/workflows/go.yml/badge.svg?event=push" alt="GitHub Actions Badge"> <img src="https://img.shields.io/github/last-commit/ibnaleem/gosearch"> <img src="https://img.shields.io/github/commit-activity/w/ibnaleem/gosearch"> <img src="https://img.shields.io/github/contributors/ibnaleem/gosearch"> <img alt="GitHub forks" src="https://img.shields.io/github/forks/ibnaleem/gosearch"> <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/ibnaleem/gosearch"> <img alt="GitHub License" src="https://img.shields.io/github/license/ibnaleem/gosearch">
</p>
<hr>

## Installation & Usage
```
$ git clone https://github.com/ibnaleem/gosearch.git && cd gosearch
```
```
$ go build
```
For Unix:
```
$ ./gosearch <username>
```
I recommend adding the `gosearch` binary to your `/bin` for universal use:
```
$ sudo mv gosearch ~/usr/bin
```
For Windows:
```
C:\Users\***\gosearch> gosearch.exe <username>
```
## Why GoSearch?
GoSearch is based on [Sherlock](https://github.com/sherlock-project/sherlock), the well-known username search tool. However, Sherlock has several shortcomings:

1. Python-based, slower than Go.
2. Outdated.
3. Reports false positives as true.
4. Fails to report false negatives.

The primary issue with Sherlock is false negatives: it may fail to detect a username on a platform when it does exist. The secondary issue is false positives: it may incorrectly identify a username as available. GoSearch addresses this by colour-coding potential false results (yellow), indicating uncertainty. This helps users quickly filter out irrelevant URLs. If there is enough demand in the future, we could add the functionality to only report full-positives or only report false negatives.

## Contributing
GoSearch relies on the [config.yaml](https://raw.githubusercontent.com/ibnaleem/gosearch/refs/heads/main/config.yaml) file, which lists all the websites to search. Users can add new sites to expand the search scope. This is where most of the contribution is needed. The general format is as follows:

```yaml
- name: "Website name"
  base_url: "https://www.website.com/profiles/{username}"
  url_probe: "optional, see below"
  errorType: "errorMsg/status_code/unknown"
  errorMsg/errorCode: "errorMsg" or 404/406/302, etc.
  cookies:
  - name: cookie name
    value: cookie value
```
Each entry should have a concise website name for easy manual searching. This avoids any duplicate submissions.
### `base_url`
The `base_url` is the URL GoSearch uses to query for usernames, unless a `url_probe` is provided (see [`url_probe`](#url_probe)). Your first task is to determine *where* user profiles are located on a website. For example, Twitter profiles are found at the root path `/`, so you would set `base_url: "https://twitter.com/{}`. The `{}` at the end of the path is a *placeholder* that GoSearch will replace with the username. 

For instance, if you query `./gosearch ibnaleem`, GoSearch will replace `{}` with "ibnaleem", resulting in the URL `https://twitter.com/privacy/ibnaleem`, assuming the query was made with `https://twitter.com/privacy/{}`.
### `url_probe`
Sometimes, websites block certain requests for security reasons but offer an API or service that can be used to retrieve the same information. The `url_probe` field is used for this purpose. It allows you to specify an API or service URL that can check the availability of a username. It's not the same as the `base_url` because GoSearch will print the API URL to the terminal, even though you’re typically looking for the profile URL.

For example, Duolingo profiles are accessible at `https://duolingo.com/profile/{}`. However, to check username availability, Duolingo provides a `url_probe` URL: `https://www.duolingo.com/2017-06-30/users?username={}`. If we used the `url_probe` as the `base_url`, the terminal would show something like `https://www.duolingo.com/2017-06-30/users?username=ibnaleem` rather than `https://duolingo.com/profile/ibnaleem`, which would be confusing for the user. GoSearch is designed with less experienced programmers in mind, so this distinction helps keep things clear and intuitive.
### `errorType`
There are 3 error types
1. `status_code` - a specific status code that is returned if a username does not exist (typically `404`)
2. `errorMsg` - a custom error message the website displays that is unique to usernames that do not exist
3. `unknown` - when there is no way of ascertaining the difference between a username that exists and does not exist on the website
#### `status_code`
The easiest to contribute, simply find an existing profile and make a request with the following code:
```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func MakeRequest(url string) {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    fmt.Println("Response:", resp.Status)
}

func main() {
    var url string = os.Args[1]
    MakeRequest(url)
}
```
```
$ go build
```
```
$ ./request https://yourwebsite.com/username
Response: 200 OK
```
Where username is the existing username on the website. Then, make the same request with a username that does not exist on the website:
```
$ ./request https://yourwebsite.com/usernamedoesnotexist
Response: 404 Not Found
```
Copy and set `errorCode`, the field under `errorType`, as the code that's printed to the terminal (in this case it's `404`).
#### `errorMsg`
This is more tricky, so what you must do is download the response body to a file. Luckily I've already written the code for you:
```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func MakeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	os.WriteFile("response.txt", []byte(body), 0644)
}

func main() {
	url := os.Args[1] // Take URL as argument from command line
	MakeRequest(url)
}
```
```
$ go build
```
```
./test https://website.com/username
```
Once again, the first username corresponds to an existing account, while the second username is for an account that does not exist. Be sure to rename `response.txt` to avoid having my code overwrite it.
```
$ mv response.txt username_found.txt
```
```
$ ./test https://website.com/username_does_not_exist
```
```
$ mv response.txt username_not_found.txt
```
You’ll need to analyse the response body of `username_not_found.txt` and compare it with `username_found.txt`. Look for any word, phrase, HTML tag, or other unique element that appears only in `username_not_found.txt`. Once you've identified something distinct, add it to the `errorMsg` field under the `errorType` field. Keep in mind that `errorType` can only have one field below it: either `errorCode` or `errorMsg`, **but not both**. Below is *incorrect*:
```yaml
  errorType: "status_code"
  errorCode: 404
  errorMsg: "<title>Hello World"</title>
```
#### `cookies`

Some websites may require cookies to retrieve specific data, such as error codes or session information. For example, the website `dzen.ru` requires the cookie `zen_sso_checked=1`, which is included in the request headers when making a browser request. To test for cookies and analyze the response, you can use the following Go code:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func MakeRequest(url string) {
	client := &http.Client{}

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Create the cookie
	cookie := &http.Cookie{
		Name:  "cookie_name",
		Value: "cookie_value",
	}

	// Add the cookie to the request
	req.AddCookie(cookie)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Output the response status
	fmt.Println("Response Status:", resp.Status)
}

func main() {
	// Ensure URL is provided as the first argument
	if len(os.Args) < 2 {
		log.Fatal("URL is required as the first argument.")
	}
	url := os.Args[1]
	MakeRequest(url)
}
```

When testing cookies, check the response status and body. For example, if you always receive a `200 OK` response, try adding `www.` before the URL, as some websites redirect based on this:

```
$ curl -I https://pinterest.com/username
HTTP/2 308
...
location: https://www.pinterest.com/username
```

```
$ curl -I https://www.pinterest.com/username
HTTP/2 200
```

Additionally, make sure to use the above code to analyse the response body when including the `www.` subdomain and relevant cookies.
#### `"unknown"`
Occasionally, the response body may be empty or lack any unique content in both the `username_not_found.txt` and `username_found.txt` files. After trying cookies, using the `www.` subdomain, you are left with no answers. In these cases, set the `errorType` to `"unknown"` (as a string) and include a `404` `errorCode` field underneath it.

To contribute, follow the template above, open a PR, and I'll merge it if GoSearch can successfully detect the accounts.

## LICENSE
This project is licensed under the GNU General Public License - see the [LICENSE](https://github.com/ibnaleem/gosearch/blob/main/LICENSE) file for details.
