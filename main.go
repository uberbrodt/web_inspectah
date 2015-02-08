package main

import (
	"net/http"
	"net/url"
	"log"
	"strings"
	"bufio"
	"os"
	"regexp"
)

var userAgentStrings []string

var subnet10 = regexp.MustCompile(`^10\..*`)
var subnet172 = regexp.MustCompile(`^172\.[1-3][678901].*`)
var subnet192 = regexp.MustCompile(`^192\.168\..*`)
var localhostRegex = regexp.MustCompile(`^http://localhost`)
var loopbackRegex = regexp.MustCompile(`127.0.0.1`)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readUserAgentFile() error {
	file, err := os.Open("./valid_user_agents.txt")
	check(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	userAgentStrings = lines
	return  scanner.Err()
}

func checkUserAgent(userAgent string) bool {
	for _,value := range userAgentStrings {
		if userAgent == value {
			return true
		}
	}
	return false
}

func checkIpAndReferrer(ipaddr string, referrer string) bool {
	if subnet10.MatchString(ipaddr) || subnet172.MatchString(ipaddr) || subnet192.MatchString(ipaddr) {
		return false
	}

	if localhostRegex.MatchString(referrer) || loopbackRegex.MatchString(referrer) {
		return false
	}
	return true
}

func main() {
	err := readUserAgentFile()
	check(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var x = r.URL.RawQuery
		if len(x) > 0 {
			var keyValMap = make(map[string]string)
			var foo = strings.Split(x,"&")
			for _,value := range foo {
				var keyValuePair = strings.Split(value,"=")
				keyValMap[keyValuePair[0]] = keyValuePair[1]
			}
			var decodedUserAgent, err = url.QueryUnescape(keyValMap["user_agent"])
			check(err)
			var result = checkUserAgent(decodedUserAgent)
			if (result == false) {
				http.Error(w, "Not Valid", 403)
				return
			}
			var decodedReferrer, errReferrer = url.QueryUnescape(keyValMap["referrer"])
			check(errReferrer)
			var referrerResult = checkIpAndReferrer(keyValMap["ip"], decodedReferrer)
			if referrerResult == false {
				http.Error(w, "Not Valid", 403)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
