package main

import (
	"fmt"
	"flag"
	"time"
	"abstraction/listener"
	// "sync"
)

/**
 * EXPLORATION: using channels to go routiners to signal for work to be executed
 *
 * go run main.go --help
 * DEFAULT: go run main.go
 * DEFAULT: go run main.go --duration 15 --numListeners 8 --maxCheckForPrimes 100000
 * go run main.go --duration 30 --numListeners 8 --maxCheckForPrimes 100000
 * go run main.go --duration 60 --numListeners 46 --maxCheckForPrimes 1999999
 *
 *
 */

// Factory function to create listeners dynamically, calling to do some work when the channel recieved a message
func createListener(id int, ch <-chan string, done chan<- bool, listener listener.Listener) {
	go func() {
		for message := range ch {
			// Call the listener's ProcessMessage method
			listener.ProcessMessage(id, message, listener.Max())
		}
		// Indicate that this listener is done processing
		done <- true
	}()
}

func main() {
	duration := flag.Int("duration", 15, "duration of test in seconds")
	// Define a flag named "numListeners", which is an integer of Number of listeners to create
	// Default value is 8, and the description is "Number of listeners"
	numListeners := flag.Int("numListeners", 8, "Number of listeners")
	maxCheckForPrimes := flag.Int64("maxCheckForPrimes", 100000, "maximum int 64 range value to check for prime")

	// Parse the flags from the command line.
	flag.Parse()

	// Create channels for each listener
	channels := make([]chan string, *numListeners)
	done := make(chan bool, *numListeners)

	// SETUP: Create channels and listeners dynamically, using the DefaultListener
	for i := 0; i < *numListeners; i++ {
		channels[i] = make(chan string)
		// Create DefaultListener instance (you can replace this with other implementations)
		defaultListener := listener.NewDefaultListenerWith(*maxCheckForPrimes)
		createListener(i+1, channels[i], done, defaultListener)
		// defer close(channels[i]) // this will cause a panic because <-timer is signaling close already
	}

	// We use a timer to stop broadcasting after X seconds
	timer := time.NewTimer(time.Duration(*duration) * time.Second)

	// Counter for broadcast messages
	index := 1

	// Goroutine to send messages to all listeners
	go func() {
		for {
			select {
			case <-timer.C: // Stop broadcasting after 5 seconds
				// use the defer close so that error and non-errors get a close
				for _, ch := range channels {
					close(ch)
				}
				return
			default:
				// Broadcast message to all listeners
				for _, ch := range channels {
					index++
					message := fmt.Sprintf("Broadcast Message: %d", index)
					ch <- message
				}
				index++
				// on main thread, Sleep to allow for message propagation
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Handle the done signals from listeners using select
	for i := 0; i < *numListeners; i++ {
		<-done
		fmt.Printf("listener %d has finished processing.\n", i)
	}

	// Optional: Wait for all goroutines to finish
	time.Sleep(1 * time.Second) // Wait for any last processing to complete
}
