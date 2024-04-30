package run

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a test",
	Long:  ``,
	// Uncomment the following line if your bare application
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running a test")
		fileName, err := cmd.Flags().GetString("filename")
		if err != nil {
			log.Fatalf("Error declaring config file/path: %v", err)
		}
		testRun(fileName)
	},
}

// Define the structs
type ExecutionConfig struct {
	Execution []Execution         `yaml:"execution"`
	Scenarios map[string]Scenario `yaml:"scenarios"`
}

func init() {
	RunCmd.Flags().StringP("filename", "f", "", "Enter the file name or complete filepath")
}

type Execution struct {
	Scenario          string         `yaml:"scenario"`
	Executor          string         `yaml:"executor"`
	Concurrency       int            `yaml:"concurrency"`
	HoldFor           string         `yaml:"hold-for"`
	RampUp            string         `yaml:"ramp-up"`
	Locations         map[string]int `yaml:"locations"`
	LocationsWeighted bool           `yaml:"locations-weighted"`
	Provisioning      string         `yaml:"provisioning"`
}

type Scenario struct {
	Requests []Request `yaml:"requests"`
}

type Request struct {
	URL             string              `yaml:"url"`
	Method          string              `yaml:"method"`
	Label           string              `yaml:"label"`
	Body            string              `yaml:"body,omitempty"`
	BodyFile        string              `yaml:"body-file,omitempty"`
	UploadFiles     []UploadFile        `yaml:"upload-files,omitempty"`
	Headers         map[string]string   `yaml:"headers,omitempty"`
	ThinkTime       string              `yaml:"think-time,omitempty"`
	Timeout         string              `yaml:"timeout,omitempty"`
	ContentEncoding string              `yaml:"content-encoding,omitempty"`
	FollowRedirects bool                `yaml:"follow-redirects,omitempty"`
	RandomSourceIP  bool                `yaml:"random-source-ip,omitempty"`
	Assert          map[string][]string `yaml:"assert,omitempty"`
	AssertJsonPath  []AssertJsonPath    `yaml:"assert-jsonpath,omitempty"`
	AssertXPath     []AssertXPath       `yaml:"assert-xpath,omitempty"`
	ExtractJsonPath ExtractJsonPath     `yaml:"extract-jsonpath,omitempty"`
	ExtractXPath    ExtractXPath        `yaml:"extract-xpath,omitempty"`
}

type UploadFile struct {
	Param    string `yaml:"param"`
	Path     string `yaml:"path"`
	MimeType string `yaml:"mime-type,omitempty"`
}

type Assert struct {
	Contains []string `yaml:"contains,omitempty"`
}

type AssertJsonPath struct {
	JsonPath      string `yaml:"jsonpath,omitempty"`
	Validate      bool   `yaml:"validate,omitempty"`
	ExpectedValue string `yaml:"expected-value,omitempty"`
}

type AssertXPath struct {
	XPath string `yaml:"xpath,omitempty"`
}

type ExtractJsonPath struct {
}

type ExtractXPath struct {
}

func testRun(fileName string) {
	// Set up Viper
	viper.SetConfigFile(fileName)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal YAML into struct
	var config ExecutionConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Access data from the unmarshaled struct
	for _, execution := range config.Execution {
		fmt.Printf("Scenario: %s\n", execution.Scenario)
		fmt.Printf("Executor: %s\n", execution.Executor)
		fmt.Printf("Concurrency: %d\n", execution.Concurrency)
		fmt.Printf("Hold For: %s\n", execution.HoldFor)
		fmt.Printf("Ramp Up: %s\n", execution.RampUp)
		fmt.Printf("Locations: %+v\n", execution.Locations)
		fmt.Printf("Locations Weighted: %t\n", execution.LocationsWeighted)
		fmt.Printf("Provisioning: %s\n", execution.Provisioning)
		fmt.Println()
	}
	// Access scenario requests
	for scenarioName, scenario := range config.Scenarios {
		fmt.Printf("Scenario Name: %s\n", scenarioName)
		for _, request := range scenario.Requests {
			fmt.Printf("URL: %s\n", request.URL)
			fmt.Printf("Method: %s\n", request.Method)
			fmt.Printf("Label: %s\n", request.Label)
			fmt.Printf("Body: %s\n", request.Body)
			fmt.Printf("Body File: %s\n", request.BodyFile)
			fmt.Printf("Upload Files: %+v\n", request.UploadFiles)
			//fmt.Printf("Headers: %+v\n", request.Headers)
			var val []string
			for k := range request.Headers {
				val = append(val, k)
			}
			for _, q := range val {
				fmt.Printf("Header: %s\n", q)
				fmt.Printf("Value: %s\n", request.Headers[q])
			}

			fmt.Printf("Think Time: %s\n", request.ThinkTime)
			fmt.Printf("Timeout: %s\n", request.Timeout)
			fmt.Printf("Content Encoding: %s\n", request.ContentEncoding)
			fmt.Printf("Follow Redirects: %t\n", request.FollowRedirects)
			fmt.Printf("Random Source IP: %t\n", request.RandomSourceIP)
			fmt.Printf("Assert: %+v\n", request.Assert)
			fmt.Printf("Assert JsonPath: %+v\n", request.AssertJsonPath)
			fmt.Printf("Assert XPath: %+v\n", request.AssertXPath)
			fmt.Printf("Extract JsonPath: %+v\n", request.ExtractJsonPath)
			fmt.Printf("Extract XPath: %+v\n", request.ExtractXPath)
			fmt.Println()
		}
	}
}
