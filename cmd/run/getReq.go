package run

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "run the requests only a test",
	Long:  ``,
	// Uncomment the following line if your bare application
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatalf("Error declaring config file/path: %v", err)
		}
		getRequestCmd(url)
	},
}

func init() {
	RunCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("url", "u", "", "Enter the url to test")
}

// This command is to call the get request directly, by-passing the config file
func getRequestCmd(url string) {
	client := &http.Client{}
	fmt.Printf("Hitting %s\n", url)
	timeStart := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	var dnsStart, dnsDone, connectStart, connectDone, gotFirstResponseByte time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },

		ConnectStart:         func(_, _ string) { connectStart = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { connectDone = time.Now() },
		GotFirstResponseByte: func() { gotFirstResponseByte = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	timeEnd := time.Now()
	defer resp.Body.Close()

	dnsTime := dnsDone.Sub(dnsStart)
	connectTime := connectDone.Sub(connectStart)
	responseTime := gotFirstResponseByte.Sub(timeStart)
	latency := timeEnd.Sub(timeStart) - responseTime
	timeString := timeStart.Format("2006-01-02 15:04:05")
	fmt.Printf("%v, Concurrency 1, Status %v, DNS: %v, ConnectTime: %v, ResponseTime: %v, Latency: %v\n",
		timeString, resp.Status, dnsTime, connectTime, responseTime, latency)
}

func getRequest(executionItem, vu int) {
	scenarioName := Config.Execution[executionItem].Scenario
	//fmt.Println("Scenario: ", scenarioName)
	for i := 0; i < len(Config.Scenarios[scenarioName].Requests); i++ {
		timeStart := time.Now()
		requestItem := Config.Scenarios[scenarioName].Requests
		client := &http.Client{}
		url := requestItem[i].URL
		//	fmt.Printf("Hitting %s\n", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		var dnsStart, dnsDone, connectStart, connectDone, gotFirstResponseByte time.Time

		trace := &httptrace.ClientTrace{
			DNSStart: func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
			DNSDone:  func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },

			ConnectStart:         func(_, _ string) { connectStart = time.Now() },
			ConnectDone:          func(_, _ string, _ error) { connectDone = time.Now() },
			GotFirstResponseByte: func() { gotFirstResponseByte = time.Now() },
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		timeEnd := time.Now()
		defer resp.Body.Close()
		labelName := requestItem[i].Label
		dnsTime := dnsDone.Sub(dnsStart)
		connectTime := connectDone.Sub(connectStart)
		responseTime := gotFirstResponseByte.Sub(timeStart)
		latency := timeEnd.Sub(timeStart) - responseTime
		timeString := timeStart.Format("2006-01-02 15:04:05")
		fmt.Printf("%v, Concurrency %v, Status %v, DNS: %v, ConnectTime: %v, ResponseTime: %v, Latency: %v, Label: %v\n",
			timeString, vu, resp.Status, dnsTime, connectTime, responseTime, latency, labelName)
	}
}
