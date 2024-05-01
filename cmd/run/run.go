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
		rest(fileName)
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

func LoadConfig(filename string) (ExecutionConfig, error) {
	var config ExecutionConfig
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Example method to print scenario names
func (c *ExecutionConfig) PrintScenarioNames() {
	fmt.Println("Scenarios:")
	for name, scenario := range c.Scenarios {
		fmt.Println("Scenario:", name)
		fmt.Println("Requests:")
		for _, request := range scenario.Requests {
			fmt.Println(" - ", request.Label)
		}
	}
}

// GetRequestsForScenario returns the requests associated with the given scenario name
func (c *ExecutionConfig) GetRequestsForScenario(name string) ([]Request, bool) {
	scenario, ok := c.Scenarios[name]
	if !ok {
		return nil, false
	}
	return scenario.Requests, true
}

func (c *ExecutionConfig) GetScenarios() []string {
	var scenarioNames []string
	for name := range c.Scenarios {
		scenarioNames = append(scenarioNames, name)
	}
	return scenarioNames
}

func rest(fileName string) {
	config, err := LoadConfig(fileName)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Println("This is from config methods")
	// Inspect parsed config
	//	fmt.Printf("%+v\n", config)
	// Use custom methods
	config.PrintScenarioNames()

	//fmt.Println("This is scenario1", config.Execution[0].Scenario)
	scenarioNames := config.GetScenarios()
	// Get requests for scenario
	for _, name := range scenarioNames {
		requests, ok := config.GetRequestsForScenario(name)
		if !ok {
			fmt.Println("Scenario not found")
		} else {
			fmt.Printf("\nRequests for %s scenario:\n", name)
			for _, request := range requests {
				fmt.Printf("Request: %+v\n", request)
			}
		}
	}
}
