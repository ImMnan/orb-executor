package run

import (
	"fmt"
	"sync"
	"time"
)

func testRun(fileName string) {
	config, err := LoadConfig(fileName)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	for i := 0; i < len(config.Execution); i++ {
		scenario := config.Execution[i].Scenario
		holdFor := config.Execution[i].GetHoldFor()
		vu := config.Execution[i].GetConcurrency()
		rampUp, steps, err := config.Execution[i].GetRampUp()
		if err != nil {
			fmt.Printf("Problem with execution %d\n", i+1)
			panic("Exiting the testrun - RampUp cannot be more that concurrency!")
		}
		executor := i
		fmt.Printf("\nThis is Orca Executor: %d, Scenario: %s\n", executor+1, scenario)
		fmt.Printf("Concurrency: %d\nRamp up: %d\nSteps: %v\nHold for: %d\n", vu, rampUp, steps, holdFor)
		concurrentHoldForRamp(holdFor, vu, rampUp, steps)
	}
}

func concurrentHoldForRamp(holdFor, vu, rampUp int, steps []int) {
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
			go concurrentVuRamp(&wgRu, vu)
			if time.Since(start) >= time.Duration(rampUp)*time.Second {
				break
			}
			wgRu.Wait()
		}
		for {
			// Once Rampup complete, the main loop will start
			wgVu.Add(1)
			go concurrentVu(&wgVu, vu)
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

func concurrentVuRamp(wgRu *sync.WaitGroup, vu int) {
	// Dummy implementation of GetLoadConfig() for this example
	start := time.Now()
	for {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!\n", timeString, vu)
		time.Sleep(100 * time.Millisecond)
		if time.Since(start) >= time.Duration(1)*time.Second {
			break
		}
	}
	wgRu.Done()
}

func concurrentVu(wgVu *sync.WaitGroup, vu int) {
	// Dummy implementation of GetLoadConfig() for this example
	//vu := l.GetThreads()
	for i := 0; i < vu; i++ {
		timePost := time.Now()
		timeString := timePost.Format("2006-01-02 15:04:05")
		fmt.Printf("%v Concurrency %v Status 200, success!\n", timeString, vu)
		time.Sleep(100 * time.Millisecond)
	}
	wgVu.Done()
}
