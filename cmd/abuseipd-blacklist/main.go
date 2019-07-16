package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const checkUrl = "https://api.abuseipdb.com/api/v2/check"

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

	client := &http.Client{}
	// Populate Check func

	// Do not check for non ip items
	if net.ParseIP(hostname).IsUnspecified() {
		ExitState(100, "net.ParseIP")
	}

	// Call API endpoint
	req, err := http.NewRequest("GET", checkUrl, nil)
	if err != nil {
		ExitState(100, "Unmarshal %s", err)
	}

	// Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Key", apikey)

	// Query Params
	query := url.Values{}
	query.Add("verbose", "true")
	query.Add("ipAddress", hostname)
	query.Add("maxAgeInDays", "1")

	// Add params
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		ExitState(100, "client.Do %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ExitState(100, "ioutil.ReadAll %s", err)
	}
	response := struct {
		Data struct {
			IpAddress      string `json:"ipAddress"`
			TotalReports   int    `json:"totalReports"`
			LastReportedAt string `json:"lastReportedAt"`
			Reports        []struct {
				Categories []int `json:"categories"`
			} `json:"reports,omitempty"`
		} `json:"data"`
		Errors []struct {
			Detail string `json:"detail"`
			Status int    `json:"status"`
		} `json:"errors"`
	}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		ExitState(100, "json.Unmarshal %s", err)
	}

	if response.Data.TotalReports > 0 {
		cat := []string{}
		for _, rep := range response.Data.Reports {
			for _, c := range rep.Categories {
				switch c {
				case 3:
					cat = append(cat, "(3) Fraud Orders")
				case 4:
					cat = append(cat, "(4) DDoS Attack")
				case 5:
					cat = append(cat, "(5) DDoS Attack")
				case 6:
					cat = append(cat, "(6) Ping of Death")
				case 7:
					cat = append(cat, "(7) Phishing")
				case 8:
					cat = append(cat, "(8) Fraud VoIP")
				case 9:
					cat = append(cat, "(9) Open Proxy")
				case 10:
					cat = append(cat, "(10) Web Spam")
				case 11:
					cat = append(cat, "(11) Email Spam")
				case 12:
					cat = append(cat, "(12) Blog Spam")
				case 13:
					cat = append(cat, "(13) VPN IP")
				case 14:
					cat = append(cat, "(14) Port Scan")
				case 15:
					cat = append(cat, "(15) Hacking")
				case 16:
					cat = append(cat, "(16) SQL Injection")
				case 17:
					cat = append(cat, "(17) Spoofing")
				case 18:
					cat = append(cat, "(18) Brute-Force")
				case 19:
					cat = append(cat, "(19) Bad Web Bot")
				case 20:
					cat = append(cat, "(20) Exploited Host")
				case 21:
					cat = append(cat, "(21) Web App Attack")
				case 22:
					cat = append(cat, "(22) SSH	Abuse")
				case 23:
					cat = append(cat, "IoT Targeted")
				default:
					cat = append(cat, fmt.Sprintf("(%d) Uncagegorised", c))
				}
			}
		}

		ExitState(1, "Listed %d in %s", response.Data.TotalReports, strings.Join(cat, ", "))
	}
	ExitState(0, "Resuming normal operation")

}

func ExitState(code int, format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", params...)
	os.Exit(code)
}
