package listener

import (
	"fmt"
	"time"
	"math"
	"golang.org/x/exp/constraints"
)

// Define the Listener interface
type Listener interface {
	ProcessMessage(id int, message string, max int64)
	Max() int64
}

// DefaultListener implements the Listener interface with default behavior
type DefaultListener struct{
	max int64
}

func (d DefaultListener) Max() int64 {
	return d.max
}

// This goes at the top level — not inside any other type
func NewDefaultListener() DefaultListener {
	return DefaultListener{max: 999999}
}

// This goes at the top level — not inside any other type
func NewDefaultListenerWith(max int64) DefaultListener {
	return DefaultListener{max: max}
}


// execute work with id, the message triggering the work and the max/limit of the computational load
func (d DefaultListener) ProcessMessage(id int, message string, max int64) {
	// Default: Print the received message with a timestamp

	// load
	start := time.Now()
	primes := findPrimes(d.max)
	elapsed := time.Since(start)
	// end of load
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Printf("%s\tListener %d received: %s\tFound %d primes in %s\n", formattedTime, id, message, len(primes), elapsed)
}

// isPrime checks if a number is prime (CPU-intensive)
func isPrime[T constraints.Integer](n T) bool {
	if n <= 1 {
		return false
	}
	// Convert T to float64 for square root calculation
	limit := int(math.Sqrt(float64(n)))
	for i := 2; i <= limit; i++ {
		if n%T(i) == 0 {
			return false
		}
	}
	return true
}

// findPrimes computes all prime numbers up to max
func findPrimes[T constraints.Integer] (max T) []T {
	var primes []T
	for i := T(2); i <= max; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}
