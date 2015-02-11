package main

import . "bayesian"
import (
	"net/http"
	"net/url"
	"log"
	"strings"
	"bufio"
	"regexp"
	"encoding/csv"
	"os"
	"fmt"
)

var userAgentStrings []string

var subnet10 = regexp.MustCompile(`^10\..*`)
var subnet172 = regexp.MustCompile(`^172\.(1[6-9]|2[0-9]|3[0-1])\..*`)
var subnet192 = regexp.MustCompile(`^192\.168\..*`)
var localhostRegex = regexp.MustCompile(`^http://localhost`)
var loopbackRegex = regexp.MustCompile(`127.0.0.1`)
var trained_classifier = return_classifier()

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

	if subnet10.MatchString(ipaddr) || subnet172.MatchString(ipaddr) || subnet192.MatchString(ipaddr) || len(ipaddr) <= 0 {
		return false
	}

	if localhostRegex.MatchString(referrer) || loopbackRegex.MatchString(referrer) || len(referrer) <= 0 {
		return false
	}
	return true
}

const (
	Good Class = "Good"
	Bad Class = "Bad"
)

func get_bad_data() []string {
	bad_data, bad_data_err := os.Open("bad_data.csv")
	check(bad_data_err)

	defer bad_data.Close()
	bad_data_reader := csv.NewReader(bad_data)
	bad_data_reader.FieldsPerRecord = 4
	rawBadDataCSV, badDataErr := bad_data_reader.ReadAll()

	check(badDataErr)

	bad_host_names := make([]string, 0)

	for _, each := range rawBadDataCSV {
		bad_host_names = append(bad_host_names,each[2])
	}

	return bad_host_names
}

func get_good_data() []string {
	good_data, good_data_err := os.Open("good_data.csv") 
	check(good_data_err)

	defer good_data.Close()
	data_reader := csv.NewReader(good_data)
	data_reader.FieldsPerRecord = 4
	rawCsv, dataErr := data_reader.ReadAll()

	check(dataErr)
	
	host_names := make([]string, 0)

	for _, each := range rawCsv {
		host_names = append(host_names,each[2])
	}

	return host_names
}

func return_classifier() *Classifier {
	classifier := NewClassifier(Good, Bad)
	classifier.Learn(get_bad_data(),Bad)
	classifier.Learn(get_good_data(),Good)

	return classifier
}

func host_names_bayesian_score(hostname string) int {
	scores, likely, _ := trained_classifier.LogScores(
		                        []string{hostname})
	fmt.Println(scores)
	return likely
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
			var decodedReferrer, errReferrer = url.QueryUnescape(keyValMap["referer"])
			check(errReferrer)
			var referrerResult = checkIpAndReferrer(keyValMap["ip"], decodedReferrer)
			if referrerResult == false {
				http.Error(w, "Not Valid", 403)
				return
			}

			var bayes_result = host_names_bayesian_score(decodedReferrer)
			if bayes_result > 0 {
				http.Error(w, "Not Valid", 403)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "No Request Data", 403)
		}
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
