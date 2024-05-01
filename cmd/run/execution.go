package run

import (
	"fmt"
	"sync"
	"time"
)

func testRun(fileName string) {
	Config, err = LoadConfig(fileName)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	for i := 0; i < len(Config.Execution); i++ {
		scenario := Config.Execution[i].Scenario
		holdFor := Config.Execution[i].GetHoldFor()
		vu := Config.Execution[i].GetConcurrency()
		rampUp, steps, err := Config.Execution[i].GetRampUp()
		if err != nil {
			fmt.Printf("Problem with execution %d\n", i+1)
			panic("Exiting the testrun - RampUp cannot be more that concurrency!")
		}
		executor := i
		fmt.Printf("\nThis is Orca Executor: %d, Scenario: %s\n", executor+1, scenario)
		fmt.Printf("Concurrency: %d\nRamp up: %d\nSteps: %v\nHold for: %d\n", vu, rampUp, steps, holdFor)
		concurrentHoldForRamp(holdFor, vu, rampUp, i, steps)
	}
}

func concurrentHoldForRamp(holdFor, vu, rampUp, scenarioItem int, steps []int) {
	// Create a channel to signal when the time period has elapsed
	done := make(chan bool)
	var wgVu sync.WaitGroup
	var wgRu sync.WaitGroup
	// Start a goroutine to execute the function
	go func() {
		// Execute the function
		start := time.Now()
		fmt.Println("This is the increment slice", steps)
		for _, vu := range steps {
			wgRu.Add(1)
			go concurrentVuRamp(&wgRu, vu, scenarioItem)
			if time.Since(start) >= time.Duration(rampUp)*time.Second {
				break
			}
			wgRu.Wait()
		}
		for {
			// Once Rampup complete, the main loop will start
			wgVu.Add(1)
			go concurrentVu(&wgVu, scenarioItem)
			fmt.Println("\n-")
			wgVu.Wait()
			if time.Since(start) >= time.Duration(holdFor)*time.Second {
				break
			}
		}
		// Signal that the goroutine has finished
		done <- true
	}()
	// Wait for the specified time period
	//	time.Sleep(3 * time.Second)
	// Wait for the goroutine to finish
	<-done
}

func concurrentVuRamp(wgRu *sync.WaitGroup, vu, scenarioItem int) {
	// Dummy implementation of GetLoadConfig() for this example
	start := time.Now()
	for {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!, scenario: %d\n", timeString, vu, scenarioItem)
		time.Sleep(100 * time.Millisecond)
		if time.Since(start) >= time.Duration(1)*time.Second {
			break
		}
	}
	wgRu.Done()
}

func concurrentVu(wgVu *sync.WaitGroup, scenarioItem int) {
	// Dummy implementation of GetLoadConfig() for this example
	vu := Config.Execution[scenarioItem].Concurrency

	for i := 0; i < vu; i++ {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!, scenario: %d\n", timeString, vu, scenarioItem)
		time.Sleep(100 * time.Millisecond)
	}
	wgVu.Done()
}
