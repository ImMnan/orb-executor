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
		//run, err := cmd.Flags().GetBool("run")
		if err != nil {
			log.Fatalf("Error with the flag: %v", err)
		}
		testRun(fileName)
	},
}

func init() {
	RunCmd.Flags().StringP("filename", "f", "", "Enter the file name or complete filepath")
	RunCmd.Flags().BoolP("run", "r", false, "testing")
}

var Config ExecutionConfig
var err error

// Define the structs
type ExecutionConfig struct {
	Execution []Execution         `yaml:"execution"`
	Scenarios map[string]Scenario `yaml:"scenarios"`
}

type Execution struct {
	Scenario          string         `yaml:"scenario"`
	Executor          string         `yaml:"executor"`
	Concurrency       int            `yaml:"concurrency"`
	HoldFor           int            `yaml:"holdfor"`
	RampUp            int            `yaml:"rampup"`
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
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		return Config, err
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return Config, err
	}
	return Config, nil
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

func (c *ExecutionConfig) GetExecutions() {
	for i := range c.Execution {
		fmt.Println("Execution: \n", i, c.Execution[i])
	}
}

func (e *Execution) GetConcurrency() int {
	return e.Concurrency
}

func (e *Execution) GetHoldFor() int {
	return e.HoldFor
}

func (e *Execution) GetRampUp() (rampUp int, increment []int, err error) {
	rampUp = e.RampUp
	vu := e.Concurrency
	step := vu / rampUp
	if step < 1 {
		err = fmt.Errorf("rampup value is too high, please set a value less than or equal to the concurrency value")
		return rampUp, nil, err
	}
	for i := 0; i < rampUp; i++ {
		increment = append(increment, (i+1)*step)
	}
	return rampUp, increment, nil
}
