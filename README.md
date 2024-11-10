# GoSearch
OSINT tool for searching usernames across social networks, written in Go. This project heavily relies on contributors, please see [Contributing](#contributing) for more details.

## Installation & Usage
```
$ git clone https://github.com/ibnaleem/gosearch.git && cd gosearch
```
```
$ go build
```
```
$ ./gosearch <username>
```
I recommend adding the `gosearch` binary to your `/bin` for universal use:
```
$ sudo mv gosearch ~/usr/bin
```
## Why GoSearch?
GoSearch is based on [Sherlock](https://github.com/sherlock-project/sherlock), the well-known username search tool. However, Sherlock has several shortcomings:

1. Python-based, slower than Go.
2. Outdated.
3. Reports false positives as true.
4. Fails to report false negatives.

The primary issue with Sherlock is false negatives: it may fail to detect a username on a platform when it does exist. The secondary issue is false positives: it may incorrectly identify a username as available. GoSearch addresses this by colour-coding potential false results (yellow), indicating uncertainty. This helps users quickly filter out irrelevant URLs. If there is enough demand in the future, we could add the functionality to only report full-positives or only report false negatives.

## Contributing
GoSearch relies on the [config.yaml](https://raw.githubusercontent.com/ibnaleem/gosearch/refs/heads/main/config.yaml) file, which lists all the websites to search. Users can add new sites to expand the search scope. The general format is as follows:

```yaml
- name: "Website name"
  base_url: "https://www.website.com/profiles/{username}"
  url_probe: "optional, see below"
  errorType: "errorMsg/status_code/unknown"
  errorMsg/errorCode: "errorMsg" or 404/406/302, etc.
```

Each entry should have a concise website name for easy manual searching. The most important field is `base_url`, where `{}` is a placeholder for the username. For example:

For Twitter/X:
```yaml
  base_url: "https://www.twitter.com/{}"
```
Here, `{}` is where the username is inserted into the URL.

For YouTube:
```yaml
  base_url: "https://www.youtube.com/c/{}"
```
Again, `{}` is inserted after `/c/`.

The `url_probe` field is used for cases where a website returns the same response (e.g., HTTP/200) for all requests, regardless of whether the account exists. For example, Duolingo always returns a `200 OK` response, even for non-existent usernames. Even inspecting the response body, we find that Duolingo returns the same response body for all requests (minus the username inside the response body). In such cases, use `url_probe` to specify a URL or endpoint that helps verify username existence.

In some cases, a platform may not have a `url_probe`, return identical response bodies for all requests to prevent bots, and use the same status code for both existing and non-existing accounts. In these cases, set the `errorType` to `"unknown"`:

```yaml
  - name: "Twitter/X"
    base_url: "https://twitter.com/{}"
    url_probe: ""
    errorType: "unknown"
    errorCode: 403
```

For example, Twitter returns a 403 status code for both existing and non-existing accounts. This means GoSearch must analyse the response body to find a unique phrase (specified in the `errorMsg` field) that only appears for non-existent accounts (e.g., "account does not exist"). Since there's no unique phrase in the response body to distinguish them, we set the `errorType` to `"unknown"`. This indicates that GoSearch cannot definitively determine whether the account exists, reporting possible false positives or negatives.

To contribute, follow the template above, open a PR, and I'll merge it if GoSearch can successfully detect the accounts.

## LICENSE
This project is licensed under the GNU General Public License - see the [LICENSE](https://github.com/ibnaleem/gosearch/blob/main/LICENSE) file for details.
