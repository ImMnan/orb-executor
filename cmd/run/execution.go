package run

import (
	"fmt"
	"sync"
	"time"
)

func testRun(Config ExecutionConfig) {
	var wgExecutor sync.WaitGroup
	for i := 0; i < len(Config.Execution); i++ {
		executionItem := i
		wgExecutor.Add(1)
		go concurrentHoldForRamp(&wgExecutor, executionItem)
	}
	wgExecutor.Wait()
}

func concurrentHoldForRamp(wgExecutor *sync.WaitGroup, executionItem int) {
	vu, holdFor, scenario, provisioning := Config.Execution[executionItem].GetExecutionDetails()
	rampUp, steps, err := Config.Execution[executionItem].GetRampUp()
	if err != nil {
		fmt.Printf("Problem with execution %d\n", executionItem+1)
		panic("Exiting the testrun - RampUp cannot be more that concurrency!")
	}

	fmt.Printf("\nThis is Orca Scenario: %s\nProvisioning %s\nVU: %d\n", scenario, provisioning, vu)
	// Create a channel to signal when the time period has elapsed
	done := make(chan bool)
	var wgHu sync.WaitGroup
	var wgRu sync.WaitGroup
	// Start a goroutine to execute the function
	go func() {
		// Execute the function
		start := time.Now()
		//fmt.Println("This is the increment slice", steps)
		for _, vu := range steps {
			// The first one is for the rampUp
			wgRu.Add(1)
			go concurrentVuRamp(&wgRu, vu, executionItem)
			if time.Since(start) >= time.Duration(rampUp)*time.Second {
				break
			}
			wgRu.Wait()
		}
		for {
			// Once Rampup complete, the main loop will start
			wgHu.Add(1)
			go concurrentVu(&wgHu, executionItem)
			fmt.Println("\n-")
			if time.Since(start) >= time.Duration(holdFor)*time.Second {
				break
			}
			wgHu.Wait()
		}
		// Signal that the goroutine has finished
		done <- true
	}()
	// Wait for the specified time period
	//	time.Sleep(3 * time.Second)
	// Wait for the goroutine to finish
	<-done
	wgExecutor.Done()
}

func concurrentVuRamp(wgRu *sync.WaitGroup, vu, executionItem int) {
	// Dummy implementation of GetLoadConfig() for this example
	start := time.Now()
	for {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!, scenario: %d\n", timeString, vu, executionItem)
		time.Sleep(100 * time.Millisecond)
		if time.Since(start) >= time.Duration(1)*time.Second {
			break
		}
	}
	wgRu.Done()
}

func concurrentVu(wgHu *sync.WaitGroup, executionItem int) {
	// Dummy implementation of GetLoadConfig() for this example
	vu := Config.Execution[executionItem].Concurrency
	for i := 0; i < vu; i++ {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!, scenario: %d\n", timeString, vu, executionItem)
		time.Sleep(100 * time.Millisecond)
	}
	wgHu.Done()
}

func getRequest() {
	fmt.Println(Config.Execution[0].Scenario)
}
