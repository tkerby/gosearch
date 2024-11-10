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
