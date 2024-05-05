package run

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func getRequest(executionItem, vu int) {
	scenarioName := Config.Execution[executionItem].Scenario
	//fmt.Println("Scenario: ", scenarioName)
	for i := 0; i < len(Config.Scenarios[scenarioName].Requests); i++ {
		requestItem := Config.Scenarios[scenarioName].Requests
		client := &http.Client{}
		url := requestItem[i].URL
		fmt.Printf("Hitting %s\n", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		trace := &httptrace.ClientTrace{
			DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
				fmt.Printf("DNS Info: %+v\n", dnsInfo)
			},
			GotConn: func(connInfo httptrace.GotConnInfo) {
				fmt.Printf("Got Conn: %+v\n", connInfo)
			},
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		//bodyText, err := io.ReadAll(resp.Body)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("%s\n", bodyText)
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status %v, scenario: %s\n", timeString, vu, resp.Status, scenarioName)
	}
}
