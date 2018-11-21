package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const base = "https://mxtoolbox.com/api/v1/lookup/blacklist/"

type (
	Lookup struct {
		UUID            string
		Command         string
		IsTransitioned  bool
		CommandArgument string
		//"TimeRecorded": "2018-11-21T14:39:54.4383364-06:00",
		//ReportingNameServer": null,
		TimeToComplete   string
		IsEndpoint       bool
		HasSubscriptions bool
		Failed           []Check
		Passed           []Check
		Timeouts         []Check
		//Transcript       map[string]string
	}
	Check struct {
		ID                         int
		Name                       string
		Info                       string
		Url                        string
		BlacklistTTL               string
		BlacklistResponseTime      string
		BlacklistReasonCode        string
		BlacklistReasonDescription string
		DelistUrl                  string
		IsExcludedByUser           bool
	}
)

var hostname, apikey string

func init() {
	flag.StringVar(&hostname, "host", "", "Host or IP to check")
	flag.StringVar(&apikey, "apikey", "", "Client Api Key")

	flag.Parse()

}

func main() {

	if apikey == "" {
		ExitState(10, "Invalid apiKey!")
	}
	client := http.Client{}

	req, err := http.NewRequest("GET", base+hostname, nil)
	req.Header.Add("Authorization", apikey)

	res, err := client.Do(req)
	if err != nil {
		ExitState(100, "Client:Get %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ExitState(100, "Client:Get %s", err)
	}
	response := &Lookup{}
	err = json.Unmarshal(body, response)
	if err != nil {
		ExitState(100, "Unmarshal %s", err)
	}
	if len(response.Failed) > 0 {
		lists := []string{}
		for _, b := range response.Failed {
			lists = append(lists, b.Name)
		}
		ExitState(1, "Blacklisted %s !", strings.Join(lists, ", "))
	}
	ExitState(0, "Resuming normal operation")

}

func ExitState(code int, format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", params...)
	os.Exit(code)
}
